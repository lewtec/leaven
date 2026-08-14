package libc

import (
	"sync"
	"unsafe"
)

// Itanium type_info vtables. User IR stores GEP(@VT, 2) as the type_info
// vptr (byte offset 16). Slots are 8 bytes so that GEP matches x86_64 IR
// on 32-bit GOARCH.
var (
	ClassTypeInfoVT    [4]uint64
	SIClassTypeInfoVT  [4]uint64
	VMIClassTypeInfoVT [4]uint64
)

// abiPtr is the LLVM pointer size in the IR we consume (x86_64).
const abiPtr = 8

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
	// IR stores the Go func's code word as *byte. Rebuild and call.
	f := FuncFromCode[func(unsafe.Pointer)](Ptr(fn))
	f(Ptr(arg))
}

func classTypeInfoVptr() unsafe.Pointer {
	return Off(Ptr(&ClassTypeInfoVT[0]), 2*abiPtr)
}
func siClassTypeInfoVptr() unsafe.Pointer {
	return Off(Ptr(&SIClassTypeInfoVT[0]), 2*abiPtr)
}
func vmiClassTypeInfoVptr() unsafe.Pointer {
	return Off(Ptr(&VMIClassTypeInfoVT[0]), 2*abiPtr)
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
	offToTop := Load[int64](vptr, -2*abiPtr)
	whole := As[byte](Off(Ptr(src), int(offToTop)))
	wholeTI := Load[*byte](vptr, -abiPtr)
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
		base := Load[*byte](Ptr(ti), 2*abiPtr)
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
