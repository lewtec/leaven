package libc

import (
	"sync/atomic"
	"unsafe"
)

func i8Word(addr *byte) (p *uint32, shift uint, mask uint32) {
	up := Addr(addr)
	base := up &^ 3
	shift = uint((up & 3) * 8)
	mask = uint32(0xff) << shift
	p = As[uint32](unsafe.Pointer(base))
	return p, shift, mask
}

// AtomicAddI8 is atomicrmw add on i8. Returns the new value (like AddInt32).
func AtomicAddI8(addr *byte, delta byte) byte {
	if addr == nil {
		return 0
	}
	p, shift, mask := i8Word(addr)
	for {
		oldw := atomic.LoadUint32(p)
		oldb := byte((oldw >> shift) & 0xff)
		newb := oldb + delta
		neww := (oldw &^ mask) | (uint32(newb) << shift)
		if atomic.CompareAndSwapUint32(p, oldw, neww) {
			return newb
		}
	}
}

// AtomicSwapI8 is atomicrmw xchg on i8. Go has no SwapInt8; CAS the
// aligned 32-bit word that holds the byte (same approach as LLVM on
// targets without native 8-bit atomics).
func AtomicSwapI8(addr *byte, neu byte) byte {
	if addr == nil {
		return 0
	}
	p, shift, mask := i8Word(addr)
	want := uint32(neu) << shift
	for {
		oldw := atomic.LoadUint32(p)
		oldb := byte((oldw >> shift) & 0xff)
		neww := (oldw &^ mask) | want
		if atomic.CompareAndSwapUint32(p, oldw, neww) {
			return oldb
		}
	}
}

// AtomicCASI8 is cmpxchg on i8 (returns whether the swap happened).
func AtomicCASI8(addr *byte, old, neu byte) bool {
	if addr == nil {
		return false
	}
	p, shift, mask := i8Word(addr)
	want := uint32(neu) << shift
	for {
		oldw := atomic.LoadUint32(p)
		oldb := byte((oldw >> shift) & 0xff)
		if oldb != old {
			return false
		}
		neww := (oldw &^ mask) | want
		if atomic.CompareAndSwapUint32(p, oldw, neww) {
			return true
		}
	}
}

// AtomicLoadI8 is the load half of a failed cmpxchg on i8.
func AtomicLoadI8(addr *byte) byte {
	if addr == nil {
		return 0
	}
	p, shift, _ := i8Word(addr)
	return byte((atomic.LoadUint32(p) >> shift) & 0xff)
}
