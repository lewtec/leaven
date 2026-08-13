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

// Free is C free(p). Also accepts RustAlloc / operator new blocks (Go heap
// pinned in allocs); those must not go to the modernc slab free.
func Free(p *byte) {
	if p == nil {
		return
	}
	u := uintptr(unsafe.Pointer(p))
	if _, ok := allocs.Load(u); ok {
		RustDealloc(unsafe.Pointer(p), 0, 1)
		return
	}
	allocatorMu.Lock()
	err := allocator.UintptrFree(u)
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

// POSIX 48-bit LCG (IEEE Std 1003.1). csmith AbsRndNumGenerator
// seeds with srand48 and draws with lrand48.
const (
	rand48Mult = uint64(0x5deece66d)
	rand48Add  = uint64(0xb)
	rand48Mask = uint64(1)<<48 - 1
)

var (
	rand48Mu    sync.Mutex
	rand48State = uint64(0x1234abcd330e)
)

// Srand48 is POSIX srand48(3): Xi = (uint32(seed) << 16) | 0x330e.
func Srand48(seed int64) {
	rand48Mu.Lock()
	rand48State = (uint64(uint32(seed)) << 16) | 0x330e
	rand48Mu.Unlock()
}

// Lrand48 is POSIX lrand48(3): bits 47..1 of the next Xi (31 bits).
func Lrand48() int64 {
	rand48Mu.Lock()
	rand48State = (rand48State*rand48Mult + rand48Add) & rand48Mask
	x := int64(rand48State >> 17)
	rand48Mu.Unlock()
	return x
}
