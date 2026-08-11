package libc

import (
	"crypto/rand"
	"sync"
	"unsafe"

	"modernc.org/memory"
)

// allocator is the process-global heap, same shape as modernc.org/libc
// (one memory.Allocator, one mutex). Zero value is ready.
var (
	allocator   memory.Allocator
	allocatorMu sync.Mutex
)

// allocRec pins a Go-heap object whose only live handle may be a uintptr
// (tagged alloca). Slab blocks are not stored here; the mmap is the pin.
type allocRec struct {
	p any
}

var allocs sync.Map

// Retain keeps p reachable until the process exits or the caller drops it.
// Only for Go-heap objects (alloca). Slab mallocs do not need it.
func Retain[T any](p *T) *T {
	if p != nil {
		allocs.LoadOrStore(uintptr(unsafe.Pointer(p)), &allocRec{p: p})
	}
	return p
}

// Malloc allocates n bytes (C malloc). n==0 allocates 1 byte, like musl/gnulib.
func Malloc[T any](n int64) *T {
	p := xmalloc(n)
	if p == 0 {
		return nil
	}
	return (*T)(unsafe.Pointer(p))
}

// Calloc allocates count*size zeroed bytes (C calloc).
func Calloc[T any](count, size int64) *T {
	n, ok := mulSize(count, size)
	if !ok {
		return nil
	}
	if n == 0 {
		n = 1
	}
	allocatorMu.Lock()
	p, err := allocator.UintptrCalloc(int(n))
	allocatorMu.Unlock()
	if err != nil || p == 0 {
		return nil
	}
	return (*T)(unsafe.Pointer(p))
}

// Realloc is C realloc. n==0 frees p and returns nil.
func Realloc(p *byte, n int64) *byte {
	if n < 0 {
		return nil
	}
	var u uintptr
	if p != nil {
		u = uintptr(unsafe.Pointer(p))
	}
	allocatorMu.Lock()
	q, err := allocator.UintptrRealloc(u, int(n))
	allocatorMu.Unlock()
	if err != nil || q == 0 {
		return nil
	}
	return (*byte)(unsafe.Pointer(q))
}

// Arc4randomBuf is BSD arc4random_buf(buf, n): fill n bytes from the CSPRNG.
func Arc4randomBuf(buf *byte, n int64) {
	if buf == nil || n <= 0 {
		return
	}
	if _, err := rand.Read(unsafe.Slice(buf, int(n))); err != nil {
		panic(err)
	}
}

// Free is C free(p).
func Free(p *byte) {
	if p == nil {
		return
	}
	allocatorMu.Lock()
	err := allocator.UintptrFree(uintptr(unsafe.Pointer(p)))
	allocatorMu.Unlock()
	if err != nil {
		return
	}
}

func xmalloc(n int64) uintptr {
	if n < 0 {
		return 0
	}
	if n == 0 {
		n = 1
	}
	allocatorMu.Lock()
	p, err := allocator.UintptrMalloc(int(n))
	allocatorMu.Unlock()
	if err != nil {
		return 0
	}
	return p
}

func mulSize(a, b int64) (int64, bool) {
	if a < 0 || b < 0 {
		return 0, false
	}
	if a != 0 && b > (1<<63-1)/a {
		return 0, false
	}
	return a * b, true
}
