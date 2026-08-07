package libc

import "unsafe"

// Malloc allocates n bytes of memory. It informs the garbage collector that
// the memory will be used to store objects of type T.
func Malloc[T any](n int64) *T {
	var p *T
	if uintptr(n) == unsafe.Sizeof(*p) {
		return new(T)
	}
	// Allocate one extra element to allow indexing off the end, like C tends
	// to do.
	count := uintptr(n)/unsafe.Sizeof(*p) + 1
	return &make([]T, count)[0]
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
func Free(p *byte) {}
