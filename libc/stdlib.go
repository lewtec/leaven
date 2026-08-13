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
// (tagged alloca via Retain). Not used for slab malloc.
type allocRec struct {
	p any
}

var allocs sync.Map

// slabLive tracks modernc blocks we still own. Double free is a no-op so
// mismatched C++/Rust drop paths do not corrupt the freelist.
var slabLive sync.Map // uintptr → struct{}

// Retain keeps p reachable until the process exits or the caller drops it.
// Only for Go-heap objects (alloca→uintptr). Slab mallocs need no pin.
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
	return As[T](unsafe.Pointer(p))
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
	p := xmalloc(n)
	if p == 0 {
		return nil
	}
	clear(Bytes(As[byte](unsafe.Pointer(p)), int(n)))
	return As[T](unsafe.Pointer(p))
}

// Realloc is C realloc. n==0 frees p and returns nil.
// Grows via malloc+copy+free so the old block goes through freeQuarantine
// (modernc UintptrRealloc freelist-scribbles the old block immediately).
func Realloc(p *byte, n int64) *byte {
	if n < 0 {
		return nil
	}
	if n == 0 {
		Free(p)
		return nil
	}
	if p == nil {
		return Malloc[byte](n)
	}
	u := Addr(p)
	if u <= 1 {
		return Malloc[byte](n)
	}
	allocatorMu.Lock()
	if _, ok := slabLive.Load(u); !ok {
		allocatorMu.Unlock()
		return Malloc[byte](n)
	}
	oldCap := memory.UintptrUsableSize(u)
	allocatorMu.Unlock()
	if oldCap >= int(n) {
		return p
	}
	q := Malloc[byte](n)
	if q == nil {
		return nil
	}
	copy(Bytes(q, int(n)), Bytes(p, oldCap))
	Free(p)
	return q
}

// Arc4randomBuf is BSD arc4random_buf(buf, n): fill n bytes from the CSPRNG.
func Arc4randomBuf(buf *byte, n int64) {
	if buf == nil || n <= 0 {
		return
	}
	if _, err := rand.Read(Bytes(buf, int(n))); err != nil {
		panic(err)
	}
}

// Free is C free(p). All heap traffic (malloc, RustAlloc, operator new with
// align≤16) shares the modernc slab.
func Free(p *byte) {
	if p == nil {
		return
	}
	slabFree(Addr(p))
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
	if err != nil || p == 0 {
		allocatorMu.Unlock()
		return 0
	}
	slabLive.Store(p, struct{}{})
	allocatorMu.Unlock()
	return p
}

// freeQuarantine delays return of blocks to modernc's freelist.
// modernc writes freelist prev/next into the first 16 bytes of a freed
// block; a use-after-free then reads that as user data (csmith seed
// 99999: expand_struct_union_vars saw a nil Variable* from prev=0).
// Depth 1 already fixes that seed; keep a modest ring for margin.
const freeQuarantineCap = 64

var (
	freeQuarantine [freeQuarantineCap]uintptr
	freeQHead      int
	freeQCount     int
)

func slabFree(u uintptr) {
	if u <= 1 {
		return
	}
	allocatorMu.Lock()
	if _, ok := slabLive.LoadAndDelete(u); !ok {
		allocatorMu.Unlock()
		return
	}
	if freeQCount == freeQuarantineCap {
		old := freeQuarantine[freeQHead]
		_ = allocator.UintptrFree(old)
		freeQuarantine[freeQHead] = u
		freeQHead = (freeQHead + 1) % freeQuarantineCap
	} else {
		freeQuarantine[(freeQHead+freeQCount)%freeQuarantineCap] = u
		freeQCount++
	}
	allocatorMu.Unlock()
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
