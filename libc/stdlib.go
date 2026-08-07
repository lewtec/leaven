package libc

import "unsafe"

// retained pins heap objects whose only live handle may be a uintptr
// (tagged-pointer unions such as tree-sitter Subtree). Without this the GC
// reclaims the object and later loads see freed memory.
var retained []any

// Retain keeps p reachable for the rest of the process. Returns p for chaining.
func Retain[T any](p *T) *T {
	if p != nil {
		retained = append(retained, p)
	}
	return p
}

// Malloc allocates n bytes of memory. It informs the garbage collector that
// the memory will be used to store objects of type T.
func Malloc[T any](n int64) *T {
	var p *T
	var out *T
	if uintptr(n) == unsafe.Sizeof(*p) {
		out = new(T)
	} else {
		// Allocate one extra element to allow indexing off the end, like C tends
		// to do.
		count := uintptr(n)/unsafe.Sizeof(*p) + 1
		out = &make([]T, count)[0]
	}
	return Retain(out)
}

// Calloc allocates a block of memory for count objects of size bytes each.
func Calloc[T any](count, size int64) *T {
	return Malloc[T](count * size)
}

// Realloc is C realloc(p, n). Without a recorded old size, this allocates a
// new block of n bytes; the previous block is left for the GC. Suitable for
// build/link; not a full C semantics clone.
func Realloc(p *byte, n int64) *byte {
	if n <= 0 {
		return nil
	}
	return Malloc[byte](n)
}

// Free is C free(p). Go is GC'd; this is a no-op so call sites typecheck.
// (Objects stay in retained; we do not reclaim.)
func Free(p *byte) {}
