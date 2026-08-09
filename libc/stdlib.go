package libc

import (
	"sync"
	"unsafe"
)

// allocRec pins p (so GC cannot collect a block whose only handle is a
// uintptr, e.g. tree-sitter Subtree) and records the requested size for Realloc.
type allocRec struct {
	p any
	n int64
}

// allocs is uintptr → *allocRec. Keys are distinct heap addresses, so
// sync.Map is the right concurrent map (disjoint writes, grow-only).
var allocs sync.Map

// Retain keeps p reachable for the rest of the process. Returns p for chaining.
func Retain[T any](p *T) *T {
	if p != nil {
		allocs.LoadOrStore(uintptr(unsafe.Pointer(p)), &allocRec{p: p})
	}
	return p
}

func allocSize(p unsafe.Pointer) int64 {
	if v, ok := allocs.Load(uintptr(p)); ok {
		return v.(*allocRec).n
	}
	return 0
}

// Malloc allocates n bytes of memory. It informs the garbage collector that
// the memory will be used to store objects of type T.
func Malloc[T any](n int64) *T {
	var p *T
	var out *T
	if n <= 0 {
		return nil
	}
	if uintptr(n) == unsafe.Sizeof(*p) {
		out = new(T)
	} else {
		// Allocate one extra element to allow indexing off the end, like C tends
		// to do.
		count := uintptr(n)/unsafe.Sizeof(*p) + 1
		out = &make([]T, count)[0]
	}
	allocs.Store(uintptr(unsafe.Pointer(out)), &allocRec{p: out, n: n})
	return out
}

// Calloc allocates a block of memory for count objects of size bytes each.
// The block is zeroed (Go new/make).
func Calloc[T any](count, size int64) *T {
	return Malloc[T](count * size)
}

// Realloc is C realloc(p, n): allocate n bytes, copy min(old,n) from p, return
// the new block. Old size comes from a side table filled by Malloc/Calloc/Realloc.
// If p is unknown (not from our allocator), behaves like malloc(n) with no copy.
func Realloc(p *byte, n int64) *byte {
	if n <= 0 {
		return nil
	}
	if p == nil {
		return Malloc[byte](n)
	}
	oldN := allocSize(unsafe.Pointer(p))
	out := Malloc[byte](n)
	if out == nil {
		return nil
	}
	copyN := oldN
	if copyN > n {
		copyN = n
	}
	if copyN > 0 {
		copy(unsafe.Slice(out, copyN), unsafe.Slice(p, copyN))
	}
	// Leave p in allocs; C free is a no-op under GC.
	return out
}

// Free is C free(p). Go is GC'd; this is a no-op so call sites typecheck.
// (Objects stay in allocs; we do not reclaim.)
func Free(p *byte) {}
