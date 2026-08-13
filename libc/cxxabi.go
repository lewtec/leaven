package libc

import (
	"sync"
	"unsafe"
)

// Itanium type_info vtables. User IR stores GEP(@VT, 2) as the type_info
// vptr. DynamicCast compares that word to identify class / si / vmi.
var (
	ClassTypeInfoVT    [4]unsafe.Pointer
	SIClassTypeInfoVT  [4]unsafe.Pointer
	VMIClassTypeInfoVT [4]unsafe.Pointer
)

const (
	tiKindUnknown = iota
	tiKindClass
	tiKindSI
	tiKindVMI

	// Itanium __base_class_type_info::__offset_flags
	tiVirtualMask = 1
	tiPublicMask  = 2
	tiOffsetShift = 8
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

func classTypeInfoVptr() unsafe.Pointer   { return unsafe.Pointer(&ClassTypeInfoVT[2]) }
func siClassTypeInfoVptr() unsafe.Pointer { return unsafe.Pointer(&SIClassTypeInfoVT[2]) }
func vmiClassTypeInfoVptr() unsafe.Pointer {
	return unsafe.Pointer(&VMIClassTypeInfoVT[2])
}

func typeInfoKind(ti *byte) int {
	if ti == nil {
		return tiKindUnknown
	}
	vptr := Load[unsafe.Pointer](Ptr(ti), 0)
	switch vptr {
	case classTypeInfoVptr():
		return tiKindClass
	case siClassTypeInfoVptr():
		return tiKindSI
	case vmiClassTypeInfoVptr():
		return tiKindVMI
	default:
		return tiKindUnknown
	}
}

// DynamicCast is Itanium __dynamic_cast(src, src_type, dst_type, src2dst).
// Null src is a null result (ABI). A missing vtable or unknown type_info
// kind panics so the unsatisfied signal is not turned into a silent null.
func DynamicCast(src, srcType, dstType *byte, src2dst int64) *byte {
	if src == nil {
		return nil
	}
	if srcType == nil || dstType == nil {
		panic("unsatisfied: __dynamic_cast")
	}
	vptr := Load[unsafe.Pointer](Ptr(src), 0)
	if vptr == nil {
		panic("unsatisfied: __dynamic_cast")
	}
	ptrSize := int(unsafe.Sizeof(uintptr(0)))
	offToTop := Load[int64](vptr, -2*ptrSize)
	whole := As[byte](Off(Ptr(src), int(offToTop)))
	wholeTI := Load[*byte](vptr, -ptrSize)
	_ = src2dst // hint only; the walk is the ABI result
	return walkType(whole, wholeTI, dstType)
}

func walkType(obj, ti, want *byte) *byte {
	if ti == nil || obj == nil {
		panic("unsatisfied: __dynamic_cast")
	}
	if ti == want {
		return obj
	}
	switch typeInfoKind(ti) {
	case tiKindClass:
		return nil
	case tiKindSI:
		base := Load[*byte](Ptr(ti), 2*int(unsafe.Sizeof(uintptr(0))))
		return walkType(obj, base, want)
	case tiKindVMI:
		return walkVMI(obj, ti, want)
	default:
		panic("unsatisfied: __dynamic_cast")
	}
}

func walkVMI(obj, ti, want *byte) *byte {
	vptr := Load[unsafe.Pointer](Ptr(obj), 0)
	if vptr == nil {
		panic("unsatisfied: __dynamic_cast")
	}
	baseCount := Load[uint32](Ptr(ti), 16+4)
	bases := Off(Ptr(ti), 24)
	var found *byte
	for i := uint32(0); i < baseCount; i++ {
		slot := Off(bases, int(i)*16)
		baseTI := Load[*byte](slot, 0)
		flags := Load[int64](slot, 8)
		if flags&tiPublicMask == 0 {
			continue
		}
		off := flags >> tiOffsetShift
		if flags&tiVirtualMask != 0 {
			off = Load[int64](vptr, int(off))
		}
		sub := As[byte](Off(Ptr(obj), int(off)))
		got := walkType(sub, baseTI, want)
		if got == nil {
			continue
		}
		if found == nil {
			found = got
			continue
		}
		if found != got {
			// Ambiguous public base: ABI returns null.
			return nil
		}
	}
	return found
}
