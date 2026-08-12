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

// OstreamInsert is std::__ostream_insert<char>(ostream&, char const*, long).
// csmith OutputMgr::OutputHeader writes the generated program header
// through this. Unknown streams go to stdout (DefaultOutputMgr::get_main_out
// is cout).
func OstreamInsert(out *byte, s *byte, n int64) *byte {
	if s != nil && n > 0 {
		_, _ = os.Stdout.Write(unsafe.Slice(s, int(n)))
	}
	return out
}

// OstreamEndl is std::endl. Writes '\n'; flush is a no-op on stdout.
func OstreamEndl(out *byte) *byte {
	_, _ = os.Stdout.Write([]byte{'\n'})
	return out
}

// OstreamLsCStr is operator<<(ostream&, char const*).
func OstreamLsCStr(out *byte, s *byte) *byte {
	if s != nil {
		_, _ = os.Stdout.Write(unsafe.Slice(s, int(Strlen(s))))
	}
	return out
}

// OstreamInsertI64 is ostream::_M_insert<long> / operator<<(long).
func OstreamInsertI64(out *byte, n int64) *byte {
	_, _ = os.Stdout.Write(strconv.AppendInt(nil, n, 10))
	return out
}

// OstreamInsertU64 is ostream::_M_insert<unsigned long>.
func OstreamInsertU64(out *byte, n uint64) *byte {
	_, _ = os.Stdout.Write(strconv.AppendUint(nil, n, 10))
	return out
}

// OstreamPut is ostream::put(char).
func OstreamPut(out *byte, c byte) *byte {
	_, _ = os.Stdout.Write([]byte{c})
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
