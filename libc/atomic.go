package libc

import (
	"sync/atomic"
	"unsafe"
)

// AtomicSwapI8 is atomicrmw xchg on i8. Go has no SwapInt8; CAS the
// aligned 32-bit word that holds the byte (same approach as LLVM on
// targets without native 8-bit atomics).
func AtomicSwapI8(addr *byte, neu byte) byte {
	if addr == nil {
		return 0
	}
	up := uintptr(unsafe.Pointer(addr))
	base := up &^ 3
	shift := uint((up & 3) * 8)
	mask := uint32(0xff) << shift
	want := uint32(neu) << shift
	p := (*uint32)(unsafe.Pointer(base))
	for {
		oldw := atomic.LoadUint32(p)
		oldb := byte((oldw >> shift) & 0xff)
		neww := (oldw &^ mask) | want
		if atomic.CompareAndSwapUint32(p, oldw, neww) {
			return oldb
		}
	}
}
