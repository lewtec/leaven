package libc

import (
	"sync"
	"unsafe"
)

// cxaEnt is one Itanium __cxa_atexit registration.
type cxaEnt struct {
	fn  *byte
	arg *byte
}

var (
	cxaMu   sync.Mutex
	cxaList []cxaEnt
)

// CxaAtexit is __cxa_atexit(fn, arg, dso). C++ static dtors register here
// from llvm.global_ctors. dso is ignored: we never unload.
func CxaAtexit(fn, arg, dso *byte) int32 {
	if fn == nil {
		return 0
	}
	cxaMu.Lock()
	cxaList = append(cxaList, cxaEnt{fn: fn, arg: arg})
	cxaMu.Unlock()
	return 0
}

// runCxaAtExit runs registered dtors last-in first-out.
func runCxaAtExit() {
	cxaMu.Lock()
	ents := cxaList
	cxaList = nil
	cxaMu.Unlock()
	for i := len(ents) - 1; i >= 0; i-- {
		invokeCxa(ents[i].fn, ents[i].arg)
	}
}

func invokeCxa(fn, arg *byte) {
	// Generated calls bitcast a named Go func to unsafe.Pointer (first
	// word) then *byte. Reconstruct the func value with a nil data word.
	var words [2]uintptr
	words[0] = uintptr(unsafe.Pointer(fn))
	f := *(*func(unsafe.Pointer))(unsafe.Pointer(&words[0]))
	f(unsafe.Pointer(arg))
}
