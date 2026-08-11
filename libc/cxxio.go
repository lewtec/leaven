package libc

import (
	"os"
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

// IfstreamClose is ifstream::close / destructor.
func IfstreamClose(this *byte) {
	if this == nil {
		return
	}
	if v, ok := streams.LoadAndDelete(uintptr(unsafe.Pointer(this))); ok {
		s := v.(*fileStream)
		if s.f != nil {
			_ = s.f.Close()
		}
	}
}

// IosFail is basic_ios::fail.
func IosFail(this *byte) bool { return streamOf(this).fail }

// IosEof is basic_ios::eof.
func IosEof(this *byte) bool { return streamOf(this).eof }

// IosNot is basic_ios::operator!.
func IosNot(this *byte) bool { return streamOf(this).fail }

// IosBool is basic_ios::operator bool (true if the stream is good).
func IosBool(this *byte) bool { return !streamOf(this).fail }
