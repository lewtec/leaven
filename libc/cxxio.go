package libc

import (
	"os"
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
)

// standinCtype is libstdc++ ctype<char>. endl does
//
//	f = *(ios+240); if f==nil throw_bad_cast; if f[56]!=0 { c = f[57+10] }
//
// Identity widen table so '\n' stays 10.
var standinCtype [ctypeSize]byte

func init() {
	standinCtype[ctypeWidenOkOff] = 1
	for i := 0; i < 256; i++ {
		standinCtype[ctypeWidenTabOff+i] = byte(i)
	}
}

// StandinVptr is an Itanium vptr into ifstreamVT. vptr-24 is slot 0 (0).
// Declare-only VTTs store these so inlined dtors do not load nil-24.
func StandinVptr() unsafe.Pointer {
	return unsafe.Add(unsafe.Pointer(&ifstreamVT[0]), 3*8)
}

func setIfstreamABI(this *byte, fail, eof bool) {
	if this == nil {
		return
	}
	base := unsafe.Pointer(this)
	*(*unsafe.Pointer)(base) = StandinVptr()
	st := int32(0)
	if fail {
		st |= iosFailbit
	}
	if eof {
		st |= iosEofbit
	}
	*(*int32)(unsafe.Add(base, iosStateOff)) = st
}

func streamOf(this *byte) *fileStream {
	if this == nil {
		return &fileStream{fail: true}
	}
	if v, ok := streams.Load(uintptr(unsafe.Pointer(this))); ok {
		return v.(*fileStream)
	}
	return &fileStream{fail: true}
}

// IfstreamOpen is basic_ifstream::basic_ifstream(char const*, ios_openmode).
// Opens the path for real. fail is set if open fails. mode is recorded only
// as in/out; unknown bits fail the stream instead of pretending.
func IfstreamOpen(this *byte, path *byte, mode int32) {
	st := &fileStream{fail: true}
	if this == nil {
		return
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
	streams.Store(uintptr(unsafe.Pointer(this)), st)
	setIfstreamABI(this, st.fail, st.eof)
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
func IfstreamClose(this *byte) {
	if this == nil {
		return
	}
	if v, ok := streams.LoadAndDelete(uintptr(unsafe.Pointer(this))); ok {
		s := v.(*fileStream)
		if s.f != nil {
			_ = s.f.Close()
		}
		setIfstreamABI(this, true, s.eof)
	}
}

// IfstreamCloseVTT is D2(this, vtt).
func IfstreamCloseVTT(this *byte, vtt *byte) { IfstreamClose(this) }

// CxxNoop is a declare-only C++ dtor we never constructed a real
// object for (locale, ios_base, __basic_file).
func CxxNoop(this *byte) {}

// FilebufClose is basic_filebuf::close. O2 inlines ifstream::close
// as filebuf::close(this+16). Return this (non-null) on success.
func FilebufClose(this *byte) *byte {
	if this == nil {
		return nil
	}
	owner := (*byte)(unsafe.Add(unsafe.Pointer(this), -filebufOff))
	IfstreamClose(owner)
	return this
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
func LocaleCtor(this *byte) {}

// IosBaseCtor is std::ios_base::ios_base() / _M_init / basic_ios::init.
// gensym's stack ostringstream calls this before operator<< / str().
// Same stand-in vptr and ctype as cout so inlined fail/endl stay honest.
func IosBaseCtor(this *byte) {
	if this == nil {
		return
	}
	InitOstream(unsafe.Pointer(this))
}

// InitOstream writes the stand-in Itanium vptr and ctype<char>.
// clang++ -O2 endl: off = *(vptr-24); ios = this+off; ctype = *(ios+240);
// null ctype → __throw_bad_cast (OutputMgr::OutputHeader).
func InitOstream(this unsafe.Pointer) {
	if this == nil {
		return
	}
	*(*unsafe.Pointer)(this) = StandinVptr()
	*(*unsafe.Pointer)(unsafe.Add(this, iosCtypeOff)) = unsafe.Pointer(&standinCtype[0])
}

// CtypeWidenInit is ctype<char>::_M_widen_init. Identity table, widen_ok=1.
func CtypeWidenInit(this *byte) {
	if this == nil {
		return
	}
	base := unsafe.Pointer(this)
	*(*byte)(unsafe.Add(base, ctypeWidenOkOff)) = 1
	tab := unsafe.Slice((*byte)(unsafe.Add(base, ctypeWidenTabOff)), 256)
	for i := range tab {
		tab[i] = byte(i)
	}
}

// ostringStreams is gensym's basic_ostringstream: side table of written bytes.
// Keyed by object address; << goes here when present, else stdout (cout).
var ostringStreams sync.Map // uintptr → *[]byte

func ostringBuf(out *byte) *[]byte {
	if out == nil {
		return nil
	}
	if v, ok := ostringStreams.Load(uintptr(unsafe.Pointer(out))); ok {
		return v.(*[]byte)
	}
	return nil
}

func writeOstream(out *byte, p []byte) {
	if len(p) == 0 {
		return
	}
	if b := ostringBuf(out); b != nil {
		*b = append(*b, p...)
		return
	}
	_, _ = os.Stdout.Write(p)
}

// stringstreamOstreamOff is the ostream subobject inside basic_stringstream
// (basic_iostream.base is 16 bytes in this libstdc++ IR). new_ctrl_vars does
// << on this+16.
const stringstreamOstreamOff = 16

// OStringStreamCtor is basic_ostringstream default ctor (gensym).
// Object is ~112 bytes; do not InitOstream (ctype at +240).
func OStringStreamCtor(this *byte) {
	if this == nil {
		return
	}
	buf := []byte{}
	ostringStreams.Store(uintptr(unsafe.Pointer(this)), &buf)
	// Stand-in vptr only (first word); no ctype slot in this size.
	*(*unsafe.Pointer)(unsafe.Pointer(this)) = StandinVptr()
}

// StringstreamDefaultCtor is basic_stringstream() used by Variable::new_ctrl_vars.
// Registers both the object and the ostream subobject (+16) for <<.
func StringstreamDefaultCtor(this *byte) {
	if this == nil {
		return
	}
	buf := []byte{}
	base := uintptr(unsafe.Pointer(this))
	ostringStreams.Store(base, &buf)
	ostringStreams.Store(base+stringstreamOstreamOff, &buf)
}

// OStringStreamClose is basic_ostringstream dtor.
func OStringStreamClose(this *byte) {
	if this == nil {
		return
	}
	ostringStreams.Delete(uintptr(unsafe.Pointer(this)))
}

// StringstreamDefaultClose is basic_stringstream dtor (default or string ctor).
func StringstreamDefaultClose(this *byte) {
	if this == nil {
		return
	}
	base := uintptr(unsafe.Pointer(this))
	ostringStreams.Delete(base)
	ostringStreams.Delete(base + stringstreamOstreamOff)
	// Also drop the read-side table used by str2int.
	StringstreamClose(this)
}

// OStringStreamStr is basic_ostringstream::str() / stringstream::str() → string.
func OStringStreamStr(ret, this *byte) {
	if ret == nil {
		return
	}
	var data []byte
	if b := ostringBuf(this); b != nil {
		data = *b
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
		writeOstream(out, unsafe.Slice(s, int(n)))
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
		writeOstream(out, unsafe.Slice(s, int(Strlen(s))))
	}
	return out
}

// OstreamLsString is operator<<(ostream&, basic_string const&).
func OstreamLsString(out *byte, s *byte) *byte {
	writeOstream(out, cxxStringBytes(s))
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

// ostreamPrecision reads ios_base::_M_precision. With StandinVptr, vbase
// offset is 0 so ios is at the ostream address; field 1 is i64 @+8.
// Default precision is 6 (libstdc++). Bookkeeper sets 3 before stats.
func ostreamPrecision(out *byte) int {
	if out == nil {
		return 6
	}
	p := *(*int64)(unsafe.Add(unsafe.Pointer(out), 8))
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
	*(*int32)(unsafe.Add(unsafe.Pointer(this), iosStateOff)) = state
	if v, ok := streams.Load(uintptr(unsafe.Pointer(this))); ok {
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
	v, ok := streams.Load(uintptr(unsafe.Pointer(is)))
	if !ok {
		st := &fileStream{fail: true, eof: true}
		streams.Store(uintptr(unsafe.Pointer(is)), st)
		setIfstreamABI(is, true, true)
		return is
	}
	st := v.(*fileStream)
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
		cxxStringAssign(str, buf)
	}
	setIfstreamABI(is, st.fail, st.eof)
	return is
}

// cxxStringAssign writes data into a libstdc++ __cxx11::basic_string.
// SSO for n<=15; otherwise RustAlloc so the inlined dtor's delete matches.
func cxxStringAssign(s *byte, data []byte) {
	if s == nil {
		return
	}
	base := unsafe.Pointer(s)
	local := (*byte)(unsafe.Add(base, cxxStringLocalOff))
	old := *(**byte)(base)
	if old != nil && old != local {
		RustDealloc(unsafe.Pointer(old), 0, 1)
	}
	n := len(data)
	if n <= cxxStringSSO {
		dst := unsafe.Slice(local, cxxStringSSO+1)
		copy(dst, data)
		dst[n] = 0
		*(**byte)(base) = local
		*(*int64)(unsafe.Add(base, cxxStringLenOff)) = int64(n)
		return
	}
	p := RustAlloc(int64(n+1), 1)
	if p == nil {
		*local = 0
		*(**byte)(base) = local
		*(*int64)(unsafe.Add(base, cxxStringLenOff)) = 0
		return
	}
	dst := unsafe.Slice((*byte)(p), n+1)
	copy(dst, data)
	dst[n] = 0
	*(**byte)(base) = (*byte)(p)
	*(*int64)(unsafe.Add(base, cxxStringLenOff)) = int64(n)
	*(*uint64)(unsafe.Add(base, cxxStringLocalOff)) = uint64(n)
}

// cxxStringBytes reads a libstdc++ __cxx11::basic_string into a Go slice.
func cxxStringBytes(s *byte) []byte {
	if s == nil {
		return nil
	}
	base := unsafe.Pointer(s)
	p := *(**byte)(base)
	n := *(*int64)(unsafe.Add(base, cxxStringLenOff))
	if p == nil || n <= 0 {
		return nil
	}
	return unsafe.Slice(p, int(n))
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
	if v, ok := stringStreams.Load(uintptr(unsafe.Pointer(this))); ok {
		return v.(*stringStream)
	}
	return nil
}

// StringstreamCtor is basic_stringstream(string const&, ios_openmode).
// csmith StringUtils::str2int builds one to parse an int (dec or hex).
// Object is ~128 bytes; do not write ostream ctype at +240.
func StringstreamCtor(this *byte, str *byte, mode int32) {
	_ = mode
	if this == nil {
		return
	}
	data := append([]byte(nil), cxxStringBytes(str)...)
	st := &stringStream{buf: data, base: 10}
	stringStreams.Store(uintptr(unsafe.Pointer(this)), st)
	// Stand-in vptr + clear iostate (fail/eof slots at +32 fit in 128).
	setIfstreamABI(this, false, false)
}

// StringstreamClose is stringstream::~stringstream / D1.
func StringstreamClose(this *byte) {
	if this == nil {
		return
	}
	stringStreams.Delete(uintptr(unsafe.Pointer(this)))
	setIfstreamABI(this, true, false)
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
	dst := (*int32)(unsafe.Pointer(out))
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
