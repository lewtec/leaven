package libc

import (
	"math"
	"math/bits"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"unsafe"
)

// RustAlloc is __rust_alloc(size, align) and C++ operator new.
// Uses the same modernc slab as Malloc/Free (no Go-heap make+pin).
// size==0 returns a non-null dangling pointer (rustc / C++ convention).
func RustAlloc(size, align int64) unsafe.Pointer {
	if size <= 0 {
		return unsafe.Pointer(uintptr(1))
	}
	if align < 1 {
		align = 1
	}
	// Round size up so the block is large enough for the alignment demand.
	// modernc returns page/chunk-aligned pointers; for larger aligns we
	// over-allocate and store the raw pointer just before the aligned one.
	if align <= 16 {
		p := xmalloc(size)
		if p == 0 {
			return nil
		}
		return unsafe.Pointer(p)
	}
	raw := xmalloc(size + align + 8)
	if raw == 0 {
		return nil
	}
	// [raw ... raw+8) stores original; payload starts at aligned address.
	base := raw + 8
	aligned := (base + uintptr(align-1)) &^ (uintptr(align) - 1)
	Store[uintptr](unsafe.Pointer(aligned), -8, raw)
	return unsafe.Pointer(aligned)
}

// RustAllocZeroed is __rust_alloc_zeroed: slab alloc then clear.
func RustAllocZeroed(size, align int64) unsafe.Pointer {
	p := RustAlloc(size, align)
	if p != nil && uintptr(p) > 1 && size > 0 {
		clear(Bytes(As[byte](p), int(size)))
	}
	return p
}

// RustDealloc is __rust_dealloc / operator delete. Same slab as Free.
func RustDealloc(p unsafe.Pointer, size, align int64) {
	if p == nil || uintptr(p) <= 1 {
		return
	}
	u := uintptr(p)
	if align > 16 {
		// Over-aligned: raw block pointer is stored 8 bytes before payload.
		u = Load[uintptr](unsafe.Pointer(u), -8)
	}
	slabFree(u)
}

// RustRealloc is __rust_realloc(ptr, oldSize, align, newSize).
func RustRealloc(p unsafe.Pointer, oldSize, align, newSize int64) unsafe.Pointer {
	if newSize <= 0 {
		RustDealloc(p, oldSize, align)
		return unsafe.Pointer(uintptr(1))
	}
	if p == nil || uintptr(p) <= 1 {
		return RustAlloc(newSize, align)
	}
	// Over-aligned blocks always copy; common path uses slab realloc.
	if align > 16 {
		n := RustAlloc(newSize, align)
		if n != nil && oldSize > 0 {
			m := oldSize
			if newSize < m {
				m = newSize
			}
			copy(Bytes(As[byte](n), int(newSize)), Bytes(As[byte](p), int(m)))
		}
		RustDealloc(p, oldSize, align)
		return n
	}
	q := Realloc((*byte)(p), newSize)
	if q == nil {
		return nil
	}
	return unsafe.Pointer(q)
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

// Linux statx(2) dirfd/flags and the stx_size field rustc FileAttr reads.
const (
	statxAtFDCWD     = -100
	statxAtEmptyPath = 0x1000
	statxAtNoFollow  = 0x100
	statxSizeOff     = 40
)

// Statx is Linux statx(2). rustc std FileAttr declares
// `extern_weak ... @statx` and calls it when the symbol exists.
// Fills Linux statx_t (stx_size at offset 40) from os.FileInfo.
func Statx(dirfd int32, pathname *byte, flags int32, mask int32, buf *byte) int32 {
	if buf == nil {
		return -1
	}
	fi, err := statxInfo(dirfd, pathname, flags)
	if err != nil {
		return -1
	}
	Memset(buf, 0, 256)
	Store[uint64](Ptr(buf), statxSizeOff, uint64(fi.Size()))
	return 0
}

func statxInfo(dirfd int32, pathname *byte, flags int32) (os.FileInfo, error) {
	path := ""
	if pathname != nil {
		path = GoString(pathname)
	}
	if flags&statxAtEmptyPath != 0 && path == "" {
		return statxFd(dirfd)
	}
	if path == "" {
		return nil, os.ErrInvalid
	}
	if !filepath.IsAbs(path) && dirfd != statxAtFDCWD {
		if dir, err := statxFdName(dirfd); err == nil {
			path = filepath.Join(dir, path)
		}
	}
	if flags&statxAtNoFollow != 0 {
		return os.Lstat(path)
	}
	return os.Stat(path)
}

func statxFd(fd int32) (os.FileInfo, error) {
	if f := fdGet(fd); f != nil {
		return f.Stat()
	}
	f := os.NewFile(uintptr(fd), "")
	if f == nil {
		return nil, os.ErrInvalid
	}
	// NewFile takes ownership; we only Stat. Drop the finalizer so GC
	// does not close the caller's fd.
	runtime.SetFinalizer(f, nil)
	return f.Stat()
}

func statxFdName(fd int32) (string, error) {
	if f := fdGet(fd); f != nil {
		return f.Name(), nil
	}
	return "", os.ErrInvalid
}
