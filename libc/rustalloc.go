package libc

import (
	"math"
	"math/bits"
	"sync/atomic"
	"unsafe"
)

// RustAlloc is __rust_alloc(size, align).
func RustAlloc(size, align int64) unsafe.Pointer {
	if size <= 0 {
		return unsafe.Pointer(uintptr(1))
	}
	b := make([]byte, int(size))
	out := unsafe.Pointer(&b[0])
	allocs.Store(uintptr(out), &allocRec{p: b})
	return out
}

// RustDealloc is __rust_dealloc(ptr, size, align).
func RustDealloc(p unsafe.Pointer, size, align int64) {
	if p != nil {
		allocs.Delete(uintptr(p))
	}
}

// RustRealloc is __rust_realloc(ptr, oldSize, align, newSize).
func RustRealloc(p unsafe.Pointer, oldSize, align, newSize int64) unsafe.Pointer {
	n := RustAlloc(newSize, align)
	if p != nil && n != nil && oldSize > 0 && newSize > 0 {
		m := oldSize
		if newSize < m {
			m = newSize
		}
		copy(unsafe.Slice((*byte)(n), int(newSize)), unsafe.Slice((*byte)(p), int(m)))
	}
	return n
}

// Fence is LLVM `fence <ordering>`. rustc emits `fence acquire` before
// Arc::drop_slow so the last refcount decrement is visible. Same barrier
// as empty asm ~{memory}.
func Fence() {
	atomic.AddUint32(&asmMemFence, 1)
}

// MaximumNumF64 is llvm.maximumnum.f64 (IEEE 754-2019 maximumNumber).
// rustc emits this for f64::max: NaN is ignored; -0.0 is less than +0.0.
func MaximumNumF64(a, b float64) float64 {
	if math.IsNaN(a) {
		return b
	}
	if math.IsNaN(b) {
		return a
	}
	if a == 0 && b == 0 {
		if math.Signbit(a) && !math.Signbit(b) {
			return b
		}
		return a
	}
	if a > b {
		return a
	}
	return b
}

// UMaxU64 is llvm.umax.i64.
func UMaxU64(a, b int64) int64 {
	if uint64(a) > uint64(b) {
		return a
	}
	return b
}

// UMulWithOverflowU64 is llvm.umul.with.overflow.i64.
func UMulWithOverflowU64(a, b int64) struct {
	F0 int64
	F1 bool
} {
	hi, lo := bits.Mul64(uint64(a), uint64(b))
	return struct {
		F0 int64
		F1 bool
	}{int64(lo), hi != 0}
}
