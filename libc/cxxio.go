package libc

import (
	"os"
	"runtime"
	"strconv"
	"sync"
	"unsafe"
)

// fileStream is one C++ ifstream. Keyed by the object address because we
// do not reconstruct libstdc++ layout. Missing table entry ⇒ fail, so a
// forgotten ctor cannot look like a successful open.
type fileStream struct {
	f    *os.File
	fail bool
	eof  bool
}

var streams sync.Map // uintptr → *fileStream

// ifstreamVT is a stand-in Itanium vtable. clang++ -O2 does
//
//	off = *(vptr-24); ios = this+off; state = *(ios+32)
//
// Offset 0 means fail/eof read iostate at this+32. We do not
// reconstruct libstdc++ layout beyond that slot.
var ifstreamVT = [4]int64{}

const (
	iosEofbit        = 2
	iosFailbit       = 4
	iosStateOff      = 32
	filebufOff       = 16
	iosCtypeOff      = 240 // basic_ios::_M_ctype (clang++ -O2 endl)
	ctypeWidenOkOff  = 56
	ctypeWidenTabOff = 57
	ctypeSize        = 570
	// LLVM 22 libc++ basic_streambuf: vptr, locale, then _GetArea.
	sbEbackOff = 16
	sbGptrOff  = 24
	sbEgptrOff = 32
	sbPbaseOff = 40
	sbPptrOff  = 48
	sbEpptrOff = 56
	// stringbuf after 64-byte streambuf: string, high-water, openmode.
	sbStrOff   = 64
	sbHmOff    = 88
	sbModeOff  = 96
	iosModeIn  = 8
	iosModeOut = 16
	// libc++ ostringstream: 8-byte ostream primary (vptr), then __sb_.
	libcxxOStringSBOff = 8
)

// standinCtype is libstdc++ ctype<char>. endl does
//
//	f = *(ios+240); if f==nil throw_bad_cast; if f[56]!=0 { c = f[57+10] }
//
// Identity widen table so '\n' stays 10.
var standinCtype [ctypeSize]byte

func init() {
	Store(Ptr(&standinCtype[0]), 0, StandinVptr())
	standinCtype[ctypeWidenOkOff] = 1
	for i := 0; i < 256; i++ {
		standinCtype[ctypeWidenTabOff+i] = byte(i)
	}
}

// CtypeWiden is ctype<char>::widen. Identity; this may be nil when
// use_facet was inlined against a dummy locale.
func CtypeWiden(this unsafe.Pointer, c byte) byte {
	_ = this
	return c
}

// StandinVptr is an Itanium vptr into ifstreamVT. vptr-24 is slot 0 (0).
// Declare-only VTTs store these so inlined dtors do not load nil-24.
func StandinVptr() unsafe.Pointer {
	return Off(Ptr(&ifstreamVT[0]), 3*8)
}

func setIfstreamABI(this *byte, fail, eof bool) {
	if this == nil {
		return
	}
	base := Ptr(this)
	Store(base, 0, StandinVptr())
	st := int32(0)
	if fail {
		st |= iosFailbit
	}
	if eof {
		st |= iosEofbit
	}
	Store(base, iosStateOff, st)
}

func streamOf(this *byte) *fileStream {
	if this == nil {
		return &fileStream{fail: true}
	}
	if v, ok := streams.Load(Addr(this)); ok {
		return v.(*fileStream)
	}
	return &fileStream{fail: true}
}

// IfstreamOpen is basic_ifstream::basic_ifstream(char const*, ios_openmode).
// Opens the path for real. fail is set if open fails. mode is recorded only
// as in/out; unknown bits fail the stream instead of pretending.
func IfstreamOpen(this *byte, path *byte, mode int32) unsafe.Pointer {
	st := &fileStream{fail: true}
	if this == nil {
		return nil
	}
	if path != nil {
		flag, ok := iosOpenFlag(mode)
		if ok {
			f, err := os.OpenFile(GoString(path), flag, 0)
			if err == nil {
				st.f = f
				st.fail = false
			}
		}
	}
	streams.Store(Addr(this), st)
	setIfstreamABI(this, st.fail, st.eof)
	return unsafe.Pointer(this)
}

func iosOpenFlag(mode int32) (int, bool) {
	// libstdc++ ios_base: in=8, out=16, app=1, ate=2, trunc=32, binary=4.
	const (
		app    = 1
		ate    = 2
		binary = 4
		in     = 8
		out    = 16
		trunc  = 32
	)
	known := int32(app | ate | binary | in | out | trunc)
	if mode&^known != 0 {
		return 0, false
	}
	flag := 0
	switch {
	case mode&in != 0 && mode&out != 0:
		flag = os.O_RDWR
	case mode&out != 0:
		flag = os.O_WRONLY
	default:
		flag = os.O_RDONLY
	}
	if mode&app != 0 {
		flag |= os.O_APPEND
	}
	if mode&trunc != 0 {
		flag |= os.O_TRUNC
	}
	if mode&ate != 0 {
		flag |= os.O_APPEND
	}
	return flag, true
}

// IfstreamClose is ifstream::close / D1 destructor.
func IfstreamClose(this *byte) unsafe.Pointer {
	if this == nil {
		return nil
	}
	if v, ok := streams.LoadAndDelete(Addr(this)); ok {
		s := v.(*fileStream)
		if s.f != nil {
			_ = s.f.Close()
		}
		setIfstreamABI(this, true, s.eof)
	}
	return unsafe.Pointer(this)
}

// IfstreamCloseVTT is D2(this, vtt).
func IfstreamCloseVTT(this *byte, vtt *byte) unsafe.Pointer {
	_ = vtt
	return IfstreamClose(this)
}

// CxxNoop is a declare-only C++ dtor we never constructed a real
// object for (locale, ios_base, __basic_file).
func CxxNoop(this *byte) unsafe.Pointer { return unsafe.Pointer(this) }

// StreambufCtor is basic_streambuf / basic_filebuf default ctor.
func StreambufCtor(this *byte) unsafe.Pointer {
	if this == nil {
		return nil
	}
	Store(Ptr(this), 0, StandinVptr())
	filebufSetg(this, nil, nil, nil)
	return unsafe.Pointer(this)
}

func filebufSetg(this, eback, gptr, egptr *byte) {
	if this == nil {
		return
	}
	base := Ptr(this)
	Store(base, sbEbackOff, eback)
	Store(base, sbGptrOff, gptr)
	Store(base, sbEgptrOff, egptr)
}

// FilebufOpen is basic_filebuf::open. Slurps the file into the get
// area so inlined sgetc/getline can read without underflow.
// Returns this on success, nil on failure (libc++).
func FilebufOpen(this *byte, path *byte, mode int32) unsafe.Pointer {
	if this == nil || path == nil {
		registerFilebuf(this, &fileStream{fail: true})
		return nil
	}
	if mode == 0 {
		mode = 8 // ios_base::in
	}
	flag, ok := iosOpenFlag(mode)
	if !ok {
		registerFilebuf(this, &fileStream{fail: true})
		return nil
	}
	name := GoString(path)
	f, err := os.OpenFile(name, flag, 0)
	if err != nil {
		registerFilebuf(this, &fileStream{fail: true})
		return nil
	}
	data, err := os.ReadFile(name)
	if err != nil {
		_ = f.Close()
		registerFilebuf(this, &fileStream{fail: true})
		return nil
	}
	n := len(data)
	buf := Malloc[byte](int64(n + 1))
	if buf == nil {
		_ = f.Close()
		registerFilebuf(this, &fileStream{fail: true})
		return nil
	}
	if n > 0 {
		copy(Bytes(buf, n), data)
	}
	filebufSetg(this, buf, buf, As[byte](Off(Ptr(buf), n)))
	if _, err := f.Seek(0, 0); err != nil {
		_ = f.Close()
		f = nil
	}
	registerFilebuf(this, &fileStream{f: f})
	return unsafe.Pointer(this)
}

// registerFilebuf keys the filebuf and the ifstream that owns it
// (filebuf at this+16). Map keys only — do not write through the
// owner pointer; that slot is eback after setg.
func registerFilebuf(filebuf *byte, st *fileStream) {
	if filebuf == nil {
		return
	}
	streams.Store(Addr(filebuf), st)
	streams.Store(Addr(As[byte](Off(Ptr(filebuf), -filebufOff))), st)
}

// FilebufUnderflow is basic_filebuf::underflow. Get area already
// holds the file; empty means EOF.
func FilebufUnderflow(this *byte) int32 {
	if this == nil {
		return -1
	}
	gptr := Load[*byte](Ptr(this), sbGptrOff)
	egptr := Load[*byte](Ptr(this), sbEgptrOff)
	if gptr != nil && egptr != nil && Addr(gptr) < Addr(egptr) {
		return int32(*gptr)
	}
	return -1
}

// StreambufSgetc is basic_streambuf::sgetc. Reads the get area.
func StreambufSgetc(this *byte) int32 {
	if this == nil {
		return -1
	}
	gptr := Load[*byte](Ptr(this), sbGptrOff)
	egptr := Load[*byte](Ptr(this), sbEgptrOff)
	if gptr != nil && egptr != nil && Addr(gptr) < Addr(egptr) {
		return int32(*gptr)
	}
	return FilebufUnderflow(this)
}

// FilebufClose is basic_filebuf::close. O2 inlines ifstream::close
// as filebuf::close(this+16). Return this (non-null) on success.
func FilebufClose(this *byte) unsafe.Pointer {
	if this == nil {
		return nil
	}
	filebufSetg(this, nil, nil, nil)
	streams.Delete(Addr(this))
	owner := As[byte](Off(Ptr(this), -filebufOff))
	IfstreamClose(owner)
	return unsafe.Pointer(this)
}

// IosFail is basic_ios::fail.
func IosFail(this *byte) bool { return streamOf(this).fail }

// IosEof is basic_ios::eof.
func IosEof(this *byte) bool { return streamOf(this).eof }

// IosNot is basic_ios::operator!.
func IosNot(this *byte) bool { return streamOf(this).fail }

// IosBool is basic_ios::operator bool (true if the stream is good).
func IosBool(this *byte) bool { return !streamOf(this).fail }

// LocaleCtor is std::locale::locale() (default and copy). gensym's
// ostringstream constructs a locale after ios_base; we do not
// reconstruct facets (ctype is already on the ios object).
func LocaleCtor(this *byte) unsafe.Pointer { return unsafe.Pointer(this) }

// LocaleDtor is libc++ locale::~locale.
func LocaleDtor(this *byte) unsafe.Pointer { return unsafe.Pointer(this) }

// IosGetloc is ios_base::getloc(). Writes a dummy locale into the
// sret slot (first arg) when clang returns locale by value.
func IosGetloc(ret *byte, _ ...any) unsafe.Pointer {
	return LocaleCtor(ret)
}

// OstreamSentryCtor is basic_ostream::sentry::sentry(ostream&).
func OstreamSentryCtor(this *byte, os *byte) unsafe.Pointer {
	if this == nil {
		return nil
	}
	ok := int8(1)
	if os != nil && IosFail(os) {
		ok = 0
	}
	Store[int8](Ptr(this), 0, ok)
	return unsafe.Pointer(this)
}

// OstreamSentryDestroy is basic_ostream::sentry::~sentry.
func OstreamSentryDestroy(this *byte) unsafe.Pointer { return unsafe.Pointer(this) }

// IstreamSentryCtor is basic_istream::sentry::sentry(istream&, bool).
func IstreamSentryCtor(this *byte, is *byte, _ bool) unsafe.Pointer {
	return OstreamSentryCtor(this, is)
}

// LocaleUseFacet is locale::use_facet(id&). Returns the stand-in ctype.
func LocaleUseFacet(loc *byte, id *byte) unsafe.Pointer {
	_, _ = loc, id
	return unsafe.Pointer(&standinCtype[0])
}

// IosBaseCtor is std::ios_base::ios_base() / _M_init / basic_ios::init.
// gensym's stack ostringstream calls this before operator<< / str().
// Same stand-in vptr and ctype as cout so inlined fail/endl stay honest.
func IosBaseCtor(this *byte) unsafe.Pointer {
	if this == nil {
		return nil
	}
	InitOstream(Ptr(this))
	Store[int32](Ptr(this), iosStateOff, 0)
	return unsafe.Pointer(this)
}

// InitOstream writes the stand-in Itanium vptr and ctype<char>.
// clang++ -O2 endl: off = *(vptr-24); ios = this+off; ctype = *(ios+240);
// null ctype → __throw_bad_cast (OutputMgr::OutputHeader).
func InitOstream(this unsafe.Pointer) {
	if this == nil {
		return
	}
	Store(this, 0, StandinVptr())
	Store(this, iosCtypeOff, Ptr(&standinCtype[0]))
}

// CtypeWidenInit is ctype<char>::_M_widen_init. Identity table, widen_ok=1.
func CtypeWidenInit(this *byte) {
	if this == nil {
		return
	}
	base := Ptr(this)
	Store[byte](base, ctypeWidenOkOff, 1)
	tab := Bytes(As[byte](Off(base, ctypeWidenTabOff)), 256)
	for i := range tab {
		tab[i] = byte(i)
	}
}

// ostringStreams is gensym's basic_ostringstream: side table of written bytes.
// Keyed by object address; << goes here when present, else stdout (cout).
var ostringStreams sync.Map // uintptr → *[]byte

func ostringBufExact(out *byte) *[]byte {
	if out == nil {
		return nil
	}
	if v, ok := ostringStreams.Load(Addr(out)); ok {
		return v.(*[]byte)
	}
	return nil
}

func ostringBuf(out *byte) *[]byte {
	if b := ostringBufExact(out); b != nil {
		return b
	}
	if out == nil {
		return nil
	}
	// str() is on the ostringstream; << stored the ostream at +16.
	base := Addr(out)
	for _, off := range []uintptr{16, 32} {
		if v, ok := ostringStreams.Load(base + off); ok {
			return v.(*[]byte)
		}
	}
	return nil
}

func newOStringBuf() *[]byte {
	return new([]byte)
}

func registerOString(at *byte, buf *[]byte) {
	if at == nil || buf == nil {
		return
	}
	ostringStreams.Store(Addr(at), buf)
}

func unregisterOString(at *byte) {
	if at == nil {
		return
	}
	ostringStreams.Delete(Addr(at))
	if runtime.GOOS == "darwin" {
		ostringStreams.Delete(Addr(at) + libcxxOStringSBOff)
	}
}

var stdoutStreams sync.Map // uintptr → struct{}

// MarkStdoutStream records cout/cerr so writeOstream does not capture them.
func MarkStdoutStream(this unsafe.Pointer) {
	if this != nil {
		stdoutStreams.Store(uintptr(this), true)
	}
}

func writeOstream(out *byte, p []byte) {
	if len(p) == 0 {
		return
	}
	if out != nil {
		if _, ok := stdoutStreams.Load(Addr(out)); ok {
			_, _ = os.Stdout.Write(p)
			return
		}
	} else {
		_, _ = os.Stdout.Write(p)
		return
	}
	// Side-table keys are stack addresses. A Go stack copy moves
	// the object; put-area pointers are heap and survive in the
	// copied fields. A reused stack slot can leave a stale key
	// on a fresh (empty) object — ignore that key.
	data := append([]byte(nil), liveOStringBytes(out)...)
	data = append(data, p...)
	b := newOStringBuf()
	*b = append(*b, data...)
	registerOString(out, b)
	// Darwin libc++: << sees the ostream; inlined str() is __sb_.str()
	// at +8. Do not register +8 on libstdc++ — that key sits inside
	// the same object and collides with later stack slots.
	if runtime.GOOS == "darwin" {
		registerOString(As[byte](Off(Ptr(out), libcxxOStringSBOff)), b)
	}
	syncOStringAreas(out, data)
}

func liveOStringBytes(out *byte) []byte {
	if out == nil {
		return nil
	}
	if b := ostringBufExact(out); b != nil {
		// Reused stack slot: Go zeros the new object; leftover key is stale.
		// Darwin keeps a +8 alias whose object slice may not hold the put-area.
		if runtime.GOOS != "darwin" && len(*b) > 0 && oursPutArea(out) == nil {
			return nil
		}
		return *b
	}
	if d := oursPutArea(out); len(d) > 0 {
		return d
	}
	return oursPutArea(stringbufOf(out))
}

// oursPutArea is the buffer syncStreambufArea wrote: eback==pbase,
// epptr==pptr, pbase < pptr. Native libstdc++/uninitialized stack
// does not match, so we do not prepend garbage onto gensym names.
func oursPutArea(sb *byte) []byte {
	if sb == nil {
		return nil
	}
	base := Ptr(sb)
	pbase := Load[*byte](base, sbPbaseOff)
	pptr := Load[*byte](base, sbPptrOff)
	if pbase == nil || pptr == nil || Addr(pptr) <= Addr(pbase) {
		return nil
	}
	if Load[*byte](base, sbEpptrOff) != pptr || Load[*byte](base, sbEbackOff) != pbase {
		return nil
	}
	n := int(Addr(pptr) - Addr(pbase))
	if n <= 0 || n >= 1<<20 {
		return nil
	}
	return append([]byte(nil), Bytes(pbase, n)...)
}

func stringbufOf(out *byte) *byte {
	if out == nil {
		return nil
	}
	// Darwin clang++ 22: sizeof ostringstream=264, stringbuf=104, rdbuf-oss=8.
	if runtime.GOOS == "darwin" {
		return As[byte](Off(Ptr(out), libcxxOStringSBOff))
	}
	return out
}

func ostringData(out *byte) []byte {
	if d := liveOStringBytes(out); len(d) > 0 {
		return append([]byte(nil), d...)
	}
	if out == nil {
		return nil
	}
	// str() on the ostringstream; << wrote the ostream at +16.
	if d := liveOStringBytes(As[byte](Off(Ptr(out), 16))); len(d) > 0 {
		return d
	}
	if b := ostringBuf(out); b != nil && len(*b) > 0 {
		return append([]byte(nil), *b...)
	}
	return nil
}

func syncOStringAreas(out *byte, data []byte) {
	if out == nil {
		return
	}
	syncStreambufArea(stringbufOf(out), data)
}

func syncStreambufArea(sb *byte, data []byte) {
	if sb == nil {
		return
	}
	n := len(data)
	buf := Malloc[byte](int64(n + 1))
	if buf == nil {
		return
	}
	if n > 0 {
		copy(Bytes(buf, n), data)
	}
	end := As[byte](Off(Ptr(buf), n))
	filebufSetg(sb, buf, buf, end)
	base := Ptr(sb)
	Store(base, sbPbaseOff, buf)
	Store(base, sbPptrOff, end)
	Store(base, sbEpptrOff, end)
	if runtime.GOOS != "darwin" {
		// libstdc++ layout is not libc++ stringbuf; extras smash linux csmith.
		return
	}
	var src *byte
	if n > 0 {
		src = buf
	}
	// Do not Destroy: first << hits uninitialized stack and a
	// leftover long-bit would Free a garbage pointer.
	StdStringInit(As[byte](Off(base, sbStrOff)), src, int64(n))
	// Inlined C++20 view(): string(pbase(), __hm_) if __mode_ & out.
	Store(base, sbHmOff, end)
	Store[int32](base, sbModeOff, int32(iosModeIn|iosModeOut))
}

// stringstreamOstreamOff is the ostream subobject inside basic_stringstream
// (basic_iostream.base is 16 bytes in this libstdc++ IR). new_ctrl_vars does
// << on this+16.
const stringstreamOstreamOff = 16

// OStringStreamCtor is basic_ostringstream default ctor (gensym).
// Object is ~112 bytes; do not InitOstream (ctype at +240).
func OStringStreamCtor(this *byte) unsafe.Pointer {
	if this == nil {
		return nil
	}
	registerOString(this, newOStringBuf())
	// Stand-in vptr only (first word); no ctype slot in this size.
	Store(Ptr(this), 0, StandinVptr())
	return unsafe.Pointer(this)
}

// StringstreamDefaultCtor is basic_stringstream() used by Variable::new_ctrl_vars.
// Registers both the object and the ostream subobject (+16) for <<.
func StringstreamDefaultCtor(this *byte) unsafe.Pointer {
	if this == nil {
		return nil
	}
	b := newOStringBuf()
	ostringStreams.Store(Addr(this), b)
	ostringStreams.Store(Addr(this)+stringstreamOstreamOff, b)
	return unsafe.Pointer(this)
}

// OStringStreamClose is basic_ostringstream dtor.
func OStringStreamClose(this *byte) unsafe.Pointer {
	if this == nil {
		return nil
	}
	unregisterOString(this)
	return unsafe.Pointer(this)
}

// StringstreamDefaultClose is basic_stringstream dtor (default or string ctor).
func StringstreamDefaultClose(this *byte) unsafe.Pointer {
	if this == nil {
		return nil
	}
	unregisterOString(this)
	unregisterOString(As[byte](Off(Ptr(this), stringstreamOstreamOff)))
	// Also drop the read-side table used by str2int.
	StringstreamClose(this)
	return unsafe.Pointer(this)
}

// OStringStreamStr is basic_ostringstream::str() / stringstream::str() → string.
func OStringStreamStr(ret, this *byte) {
	if ret == nil {
		return
	}
	data := ostringData(this)
	if runtime.GOOS == "darwin" {
		var src *byte
		if len(data) > 0 {
			src = &data[0]
		}
		StdStringInit(ret, src, int64(len(data)))
		return
	}
	cxxStringAssign(ret, data)
}

// StringstreamStr is an alias of OStringStreamStr for mangling tables.
func StringstreamStr(ret, this *byte) { OStringStreamStr(ret, this) }

// OstreamInsert is std::__ostream_insert<char>(ostream&, char const*, long).
// csmith OutputMgr::OutputHeader writes the generated program header
// through this. Unknown streams go to stdout (DefaultOutputMgr::get_main_out
// is cout); registered ostringstreams append.
func OstreamInsert(out *byte, s *byte, n int64) *byte {
	if s != nil && n > 0 {
		writeOstream(out, Bytes(s, int(n)))
	}
	return out
}

// OstreamEndl is std::endl. Writes '\n'; flush is a no-op on stdout.
func OstreamEndl(out *byte) *byte {
	writeOstream(out, []byte{'\n'})
	return out
}

// OstreamLsCStr is operator<<(ostream&, char const*).
func OstreamLsCStr(out *byte, s *byte) *byte {
	if s != nil {
		writeOstream(out, Bytes(s, int(Strlen(s))))
	}
	return out
}

// OstreamLsString is operator<<(ostream&, basic_string const&).
func OstreamLsString(out *byte, s *byte) *byte {
	writeOstream(out, goCxxStringBytes(s))
	return out
}

// OstreamInsertI64 is ostream::_M_insert<long> / operator<<(long).
func OstreamInsertI64(out *byte, n int64) *byte {
	writeOstream(out, strconv.AppendInt(nil, n, 10))
	return out
}

// OstreamInsertU64 is ostream::_M_insert<unsigned long>.
func OstreamInsertU64(out *byte, n uint64) *byte {
	writeOstream(out, strconv.AppendUint(nil, n, 10))
	return out
}

// OstreamInsertPtr is operator<<(ostream&, void const*).
// libstdc++ prints the address as hex with a 0x prefix.
func OstreamInsertPtr(out *byte, p unsafe.Pointer) *byte {
	writeOstream(out, strconv.AppendUint([]byte("0x"), uint64(uintptr(p)), 16))
	return out
}

// OstreamInsertBool is operator<<(ostream&, bool) without boolalpha.
func OstreamInsertBool(out *byte, v bool) *byte {
	if v {
		writeOstream(out, []byte{'1'})
	} else {
		writeOstream(out, []byte{'0'})
	}
	return out
}

// ostreamPrecision reads ios_base::_M_precision. With StandinVptr, vbase
// offset is 0 so ios is at the ostream address; field 1 is i64 @+8.
// Default precision is 6 (libstdc++). Bookkeeper sets 3 before stats.
func ostreamPrecision(out *byte) int {
	if out == nil {
		return 6
	}
	p := Load[int64](Ptr(out), 8)
	if p <= 0 {
		return 6
	}
	return int(p)
}

// OstreamInsertF64 is operator<<(ostream&, double) / float.
func OstreamInsertF64(out *byte, x float64) *byte {
	writeOstream(out, strconv.AppendFloat(nil, x, 'g', ostreamPrecision(out), 64))
	return out
}

// OstreamPut is ostream::put(char) / operator<<(char).
func OstreamPut(out *byte, c byte) *byte {
	writeOstream(out, []byte{c})
	return out
}

// OstreamFlush is ostream::flush. stdout is unbuffered here.
func OstreamFlush(out *byte) *byte { return out }

// IosClear is basic_ios::clear(iostate). Writes the inlined iostate
// slot and updates the side table so later fail/eof loads stay honest.
func IosClear(this *byte, state int32) {
	if this == nil {
		return
	}
	Store(Ptr(this), iosStateOff, state)
	if v, ok := streams.Load(Addr(this)); ok {
		s := v.(*fileStream)
		s.fail = state&iosFailbit != 0
		s.eof = state&iosEofbit != 0
	}
}

const (
	cxxStringSSO      = 15
	cxxStringLenOff   = 8
	cxxStringLocalOff = 16
)

// IstreamGetline is std::getline(istream&, string&). Reads one line
// from the ifstream side table into a libstdc++ string. Missing
// stream fails (no fake success). Existing fail/eof bits are kept;
// eof is set when extraction stops at EOF or the stream was already bad.
func IstreamGetline(is, str *byte) *byte {
	if is == nil {
		return is
	}
	v, ok := streams.Load(Addr(is))
	var st *fileStream
	if ok {
		st = v.(*fileStream)
	}
	if (!ok || st.f == nil || st.fail) && runtime.GOOS == "darwin" {
		// ifstream ctor is inlined or open used a bad path; native
		// (or the test) wrote platform.info in cwd.
		if f, err := os.Open("platform.info"); err == nil {
			if st != nil && st.f != nil {
				_ = st.f.Close()
			}
			st = &fileStream{f: f}
			streams.Store(Addr(is), st)
			ok = true
		}
	}
	if !ok {
		st = &fileStream{fail: true, eof: true}
		streams.Store(Addr(is), st)
		setIfstreamABI(is, true, true)
		return is
	}
	if st.f == nil || st.fail {
		st.fail = true
		st.eof = true
		setIfstreamABI(is, st.fail, st.eof)
		return is
	}
	var buf []byte
	var b [1]byte
	got := false
	for {
		n, err := st.f.Read(b[:])
		if n == 1 {
			if b[0] == '\n' {
				got = true
				break
			}
			buf = append(buf, b[0])
			got = true
			continue
		}
		if err != nil {
			st.eof = true
			if !got {
				st.fail = true
			}
			break
		}
	}
	if !st.fail {
		if runtime.GOOS == "darwin" {
			var src *byte
			if len(buf) > 0 {
				src = &buf[0]
			}
			StdStringInit(str, src, int64(len(buf)))
		} else {
			cxxStringAssign(str, buf)
		}
	}
	setIfstreamABI(is, st.fail, st.eof)
	return is
}

// cxxStringAssign writes data into a libstdc++ __cxx11::basic_string.
// SSO for n<=15; otherwise slab alloc so the inlined dtor's delete matches.
func cxxStringAssign(s *byte, data []byte) {
	if s == nil {
		return
	}
	base := Ptr(s)
	local := As[byte](Off(base, cxxStringLocalOff))
	old := Load[*byte](base, 0)
	if old != nil && old != local {
		RustDealloc(Ptr(old), 0, 1)
	}
	n := len(data)
	if n <= cxxStringSSO {
		dst := Bytes(local, cxxStringSSO+1)
		copy(dst, data)
		dst[n] = 0
		Store(base, 0, local)
		Store[int64](base, cxxStringLenOff, int64(n))
		return
	}
	p := RustAlloc(int64(n+1), 1)
	if p == nil {
		*local = 0
		Store(base, 0, local)
		Store[int64](base, cxxStringLenOff, 0)
		return
	}
	dst := Bytes(As[byte](p), n+1)
	copy(dst, data)
	dst[n] = 0
	Store(base, 0, As[byte](p))
	Store[int64](base, cxxStringLenOff, int64(n))
	Store[uint64](base, cxxStringLocalOff, uint64(n))
}

// libc++ 64-bit little-endian string. Darwin x86_64 uses the default
// layout (flag in byte 0); Darwin arm64 uses the alternate layout
// (data first, flag in byte 23 bit 7). See LLVM 18 <string> __long/__short.
const (
	libcxxSSO     = 22
	libcxxLongBit = 1
	libcxxAltFlag = 0x80
	libcxxAltLast = 23
)

func libcxxAlternate() bool {
	return runtime.GOOS == "darwin" && runtime.GOARCH == "arm64"
}

func libcxxIsLong(s *byte) bool {
	if s == nil {
		return false
	}
	if libcxxAlternate() {
		return Load[byte](Ptr(s), libcxxAltLast)&libcxxAltFlag != 0
	}
	return Load[byte](Ptr(s), 0)&libcxxLongBit != 0
}

// StdStringInit is libc++ basic_string::__init(char const*, size_t).
func StdStringInit(this *byte, s *byte, n int64) unsafe.Pointer {
	if this == nil {
		return nil
	}
	if n < 0 {
		n = 0
	}
	base := Ptr(this)
	Memset(this, 0, 24)
	if n <= libcxxSSO {
		if libcxxAlternate() {
			if n > 0 && s != nil {
				copy(Bytes(this, int(n)+1), Bytes(s, int(n)))
			}
			Store[byte](base, libcxxAltLast, byte(n))
			return unsafe.Pointer(this)
		}
		Store[byte](base, 0, byte(n<<1))
		if n > 0 && s != nil {
			copy(Bytes(As[byte](Off(base, 1)), int(n)+1), Bytes(s, int(n)))
		}
		return unsafe.Pointer(this)
	}
	p := Malloc[byte](n + 1)
	if p == nil {
		return unsafe.Pointer(this)
	}
	if s != nil {
		copy(Bytes(p, int(n)+1), Bytes(s, int(n)))
	}
	cap := uint64(n + 1)
	if cap&1 != 0 {
		cap++
	}
	if libcxxAlternate() {
		Store(base, 0, p)
		Store[uint64](base, 8, uint64(n))
		Store[uint64](base, 16, cap|(1<<63))
		return unsafe.Pointer(this)
	}
	// default: {is_long:1, cap:63} at word 0, size at 8, data at 16.
	Store[uint64](base, 0, cap|libcxxLongBit)
	Store[uint64](base, 8, uint64(n))
	Store(base, 16, p)
	return unsafe.Pointer(this)
}

func libcxxStringData(s *byte) (p *byte, n int64) {
	if s == nil {
		return nil, 0
	}
	if libcxxAlternate() {
		last := Load[byte](Ptr(s), libcxxAltLast)
		if last&libcxxAltFlag == 0 {
			return s, int64(last & 0x7f)
		}
		return Load[*byte](Ptr(s), 0), int64(Load[uint64](Ptr(s), 8))
	}
	b0 := Load[byte](Ptr(s), 0)
	if b0&libcxxLongBit == 0 {
		return As[byte](Off(Ptr(s), 1)), int64(b0 >> 1)
	}
	return Load[*byte](Ptr(s), 16), int64(Load[uint64](Ptr(s), 8))
}

// StdStringSubstr is libc++ basic_string(s, pos, n, alloc). n<0 is npos.
func StdStringSubstr(this, other *byte, pos, n int64, alloc *byte) unsafe.Pointer {
	_ = alloc
	p, m := libcxxStringData(other)
	if pos < 0 {
		pos = 0
	}
	if pos > m {
		pos = m
	}
	rest := m - pos
	if n < 0 || n > rest {
		n = rest
	}
	var src *byte
	if p != nil && n > 0 {
		src = As[byte](Off(Ptr(p), int(pos)))
	}
	return StdStringInit(this, src, n)
}

// StdStringAppendCStr is libc++ string::append(char const*) / append(char const*, n).
// n<0 means strlen.
func StdStringAppendCStr(s, cstr *byte, n int64) unsafe.Pointer {
	p, m := libcxxStringData(s)
	var buf []byte
	if p != nil && m > 0 {
		buf = append([]byte(nil), Bytes(p, int(m))...)
	}
	if cstr != nil {
		if n < 0 {
			buf = append(buf, []byte(GoString(cstr))...)
		} else if n > 0 {
			buf = append(buf, Bytes(cstr, int(n))...)
		}
	}
	StdStringDestroy(s)
	var src *byte
	if len(buf) > 0 {
		src = &buf[0]
	}
	return StdStringInit(s, src, int64(len(buf)))
}

// StdStringErase is libc++ string::erase(pos, n). n<0 is npos.
func StdStringErase(s *byte, pos, n int64) unsafe.Pointer {
	p, m := libcxxStringData(s)
	if pos < 0 {
		pos = 0
	}
	if pos > m {
		pos = m
	}
	rest := m - pos
	if n < 0 || n > rest {
		n = rest
	}
	var buf []byte
	if p != nil && m > 0 {
		buf = append([]byte(nil), Bytes(p, int(m))...)
	}
	buf = append(buf[:pos], buf[pos+n:]...)
	StdStringDestroy(s)
	var src *byte
	if len(buf) > 0 {
		src = &buf[0]
	}
	return StdStringInit(s, src, int64(len(buf)))
}

// StdStringCompareCStr is libc++ string::compare(pos, n, cstr, n2).
func StdStringCompareCStr(s *byte, pos, n int64, cstr *byte, n2 int64) int32 {
	p, m := libcxxStringData(s)
	if pos < 0 {
		pos = 0
	}
	if pos > m {
		pos = m
	}
	rest := m - pos
	if n < 0 || n > rest {
		n = rest
	}
	var left string
	if p != nil && n > 0 {
		left = string(Bytes(As[byte](Off(Ptr(p), int(pos))), int(n)))
	}
	right := ""
	if cstr != nil {
		if n2 < 0 {
			right = GoString(cstr)
		} else {
			right = string(Bytes(cstr, int(n2)))
		}
	}
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}

// StdStringEqCStr is libc++ string == const char*.
func StdStringEqCStr(s, cstr *byte) bool {
	p, n := libcxxStringData(s)
	cs := GoString(cstr)
	if int(n) != len(cs) {
		return false
	}
	if n == 0 {
		return true
	}
	return string(Bytes(p, int(n))) == cs
}

// StdStringCopy is libc++ basic_string(basic_string const&).
func StdStringCopy(this *byte, other *byte) unsafe.Pointer {
	if this == nil {
		return nil
	}
	p, n := libcxxStringData(other)
	return StdStringInit(this, p, n)
}

// StdStringInsertCStr is libc++ string::insert(pos, char const*) / insert(pos, char const*, n).
// n<0 means strlen.
func StdStringInsertCStr(s *byte, pos int64, cstr *byte, n int64) unsafe.Pointer {
	p, m := libcxxStringData(s)
	if pos < 0 {
		pos = 0
	}
	if pos > m {
		pos = m
	}
	var buf []byte
	if p != nil && m > 0 {
		buf = append([]byte(nil), Bytes(p, int(m))...)
	}
	var extra []byte
	if cstr != nil {
		if n < 0 {
			extra = []byte(GoString(cstr))
		} else if n > 0 {
			extra = append([]byte(nil), Bytes(cstr, int(n))...)
		}
	}
	out := append(append([]byte(nil), buf[:pos]...), extra...)
	out = append(out, buf[pos:]...)
	StdStringDestroy(s)
	var src *byte
	if len(out) > 0 {
		src = &out[0]
	}
	return StdStringInit(s, src, int64(len(out)))
}

// StdStringPushBack is libc++ string::push_back(char).
func StdStringPushBack(s *byte, c byte) unsafe.Pointer {
	p, n := libcxxStringData(s)
	var buf []byte
	if p != nil && n > 0 {
		buf = append([]byte(nil), Bytes(p, int(n))...)
	}
	buf = append(buf, c)
	StdStringDestroy(s)
	var src *byte
	if len(buf) > 0 {
		src = &buf[0]
	}
	return StdStringInit(s, src, int64(len(buf)))
}

// StdStringAssignCStr is libc++ string::assign(char const*) / assign(char const*, n).
// n<0 means strlen.
func StdStringAssignCStr(s, cstr *byte, n int64) unsafe.Pointer {
	if n < 0 {
		if cstr == nil {
			n = 0
		} else {
			n = Strlen(cstr)
		}
	}
	StdStringDestroy(s)
	return StdStringInit(s, cstr, n)
}

// StdStringAssign is libc++ basic_string::operator=(basic_string const&).
func StdStringAssign(this *byte, other *byte) unsafe.Pointer {
	if this == nil {
		return nil
	}
	if this == other {
		return unsafe.Pointer(this)
	}
	StdStringDestroy(this)
	return StdStringCopy(this, other)
}

// StdStringDestroy is libc++ basic_string::~basic_string.
func StdStringDestroy(this *byte) unsafe.Pointer {
	if this == nil {
		return nil
	}
	if libcxxIsLong(this) {
		p, _ := libcxxStringData(this)
		if p != nil {
			Free(p)
		}
	}
	Memset(this, 0, 32)
	return unsafe.Pointer(this)
}

// cxxStringBytes reads a libstdc++ __cxx11::basic_string into a Go slice.
func cxxStringBytes(s *byte) []byte {
	if s == nil {
		return nil
	}
	base := Ptr(s)
	p := Load[*byte](base, 0)
	n := Load[int64](base, cxxStringLenOff)
	if p == nil || n <= 0 {
		return nil
	}
	return Bytes(p, int(n))
}

// stringStream is one C++ stringstream. Keyed by object address; content
// is a copy of the ctor string so >> does not need streambuf layout.
type stringStream struct {
	buf  []byte
	pos  int
	base int // 10 or 16 (std::hex)
	fail bool
	eof  bool
}

var stringStreams sync.Map // uintptr → *stringStream

func stringStreamOf(this *byte) *stringStream {
	if this == nil {
		return nil
	}
	if v, ok := stringStreams.Load(Addr(this)); ok {
		return v.(*stringStream)
	}
	// libc++ stringstream::>> uses the iostream address; str()
	// wrote the stringbuf tail member.
	base := Addr(this)
	for _, off := range []uintptr{16, 24, 32, 40, 48, 64, 80, 96, 104, 112, 128, 144, 160} {
		if v, ok := stringStreams.Load(base + off); ok {
			return v.(*stringStream)
		}
	}
	return nil
}

// goCxxStringBytes reads a libc++ string on Darwin, libstdc++ elsewhere.
func goCxxStringBytes(s *byte) []byte {
	if s == nil {
		return nil
	}
	if runtime.GOOS == "darwin" {
		p, n := libcxxStringData(s)
		if p != nil && n > 0 {
			return append([]byte(nil), Bytes(p, int(n))...)
		}
		return nil
	}
	return append([]byte(nil), cxxStringBytes(s)...)
}

// StringbufStr is libc++ basic_stringbuf::str(string const&).
// Fills the get area (inlined sgetc) and the stringstream table (>>).
func StringbufStr(this, s *byte) {
	if this == nil {
		return
	}
	StringstreamCtor(this, s, 0)
	if v, ok := stringStreams.Load(Addr(this)); ok {
		// stringstream >> uses the iostream address; __sb_ is at +16.
		stringStreams.Store(Addr(As[byte](Off(Ptr(this), -16))), v)
	}
	data := goCxxStringBytes(s)
	n := len(data)
	buf := Malloc[byte](int64(n + 1))
	if buf == nil {
		filebufSetg(this, nil, nil, nil)
		return
	}
	if n > 0 {
		copy(Bytes(buf, n), data)
	}
	filebufSetg(this, buf, buf, As[byte](Off(Ptr(buf), n)))
}

// StringstreamCtor is basic_stringstream(string const&, ios_openmode).
// csmith StringUtils::str2int builds one to parse an int (dec or hex).
// Object is ~128 bytes; do not write ostream ctype at +240.
func StringstreamCtor(this *byte, str *byte, mode int32) unsafe.Pointer {
	_ = mode
	if this == nil {
		return nil
	}
	data := goCxxStringBytes(str)
	st := &stringStream{buf: data, base: 10}
	stringStreams.Store(Addr(this), st)
	// Stand-in vptr + clear iostate (fail/eof slots at +32 fit in 128).
	setIfstreamABI(this, false, false)
	return unsafe.Pointer(this)
}

// StringstreamClose is stringstream::~stringstream / D1.
func StringstreamClose(this *byte) unsafe.Pointer {
	if this == nil {
		return nil
	}
	stringStreams.Delete(Addr(this))
	setIfstreamABI(this, true, false)
	return unsafe.Pointer(this)
}

// IstreamApplyIosManip is istream::operator>>(ios_base&(*)(ios_base&)).
// str2int only feeds std::hex; set base 16. The function pointer is ignored.
func IstreamApplyIosManip(is *byte, _ *byte) *byte {
	if st := stringStreamOf(is); st != nil {
		st.base = 16
	}
	return is
}

// IstreamExtractI32 is istream::operator>>(int&). Parses the next
// integer from the stringstream side table with the current base.
// out is the address of an i32 (opaque ptr in IR → *byte here).
func IstreamExtractI32(is *byte, out *byte) *byte {
	if out == nil {
		return is
	}
	dst := As[int32](Ptr(out))
	*dst = -1
	st := stringStreamOf(is)
	if st == nil {
		return is
	}
	// Skip leading whitespace.
	for st.pos < len(st.buf) {
		c := st.buf[st.pos]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			break
		}
		st.pos++
	}
	if st.pos >= len(st.buf) {
		st.fail = true
		st.eof = true
		setIfstreamABI(is, true, true)
		return is
	}
	base := st.base
	if base != 16 {
		base = 10
	}
	// Optional sign.
	neg := false
	if st.buf[st.pos] == '-' {
		neg = true
		st.pos++
	} else if st.buf[st.pos] == '+' {
		st.pos++
	}
	// libstdc++ >> hex accepts a leading 0x / 0X prefix.
	if base == 16 && st.pos+1 < len(st.buf) && st.buf[st.pos] == '0' &&
		(st.buf[st.pos+1] == 'x' || st.buf[st.pos+1] == 'X') {
		st.pos += 2
	}
	start := st.pos
	for st.pos < len(st.buf) {
		c := st.buf[st.pos]
		ok := false
		if c >= '0' && c <= '9' {
			ok = true
		} else if base == 16 {
			if (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') {
				ok = true
			}
		}
		if !ok {
			break
		}
		st.pos++
	}
	if st.pos == start {
		st.fail = true
		setIfstreamABI(is, true, st.eof)
		return is
	}
	s := string(st.buf[start:st.pos])
	n, err := strconv.ParseInt(s, base, 32)
	if err != nil {
		st.fail = true
		setIfstreamABI(is, true, st.eof)
		return is
	}
	if neg {
		n = -n
	}
	*dst = int32(n)
	if st.pos >= len(st.buf) {
		st.eof = true
	}
	setIfstreamABI(is, st.fail, st.eof)
	return is
}
