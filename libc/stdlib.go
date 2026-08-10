package libc

import (
	"reflect"
	"sync"
	"unsafe"
)

// allocRec pins p (so GC cannot collect a block whose only handle is a
// uintptr, e.g. tree-sitter Subtree) and records the requested size for Realloc.
// Freed recs stay in allocs (live=false) so reuse does not allocate map nodes.
// The stack holds the same rec; working set is peak live, not per-parse growth.
type allocRec struct {
	p    any
	n    int64 // logical size for Realloc copy
	cap  int64 // requested size at allocate; 0 means n
	live bool
}

// allocs is uintptr → *allocRec. Keys are distinct heap addresses, so
// sync.Map is the right concurrent map (disjoint writes).
var allocs sync.Map

type poolKey struct {
	t reflect.Type
	n int64
}

type recStack struct {
	mu   sync.Mutex
	recs []*allocRec
}

func (s *recStack) push(r *allocRec) {
	s.mu.Lock()
	s.recs = append(s.recs, r)
	s.mu.Unlock()
}

func (s *recStack) pop() *allocRec {
	s.mu.Lock()
	n := len(s.recs)
	if n == 0 {
		s.mu.Unlock()
		return nil
	}
	r := s.recs[n-1]
	s.recs[n-1] = nil
	s.recs = s.recs[:n-1]
	s.mu.Unlock()
	return r
}

var stacks sync.Map // poolKey → *recStack

func stackFor(k poolKey) *recStack {
	if v, ok := stacks.Load(k); ok {
		return v.(*recStack)
	}
	s := new(recStack)
	actual, loaded := stacks.LoadOrStore(k, s)
	if loaded {
		return actual.(*recStack)
	}
	return s
}

func recCap(rec *allocRec) int64 {
	if rec.cap != 0 {
		return rec.cap
	}
	return rec.n
}

// Retain keeps p reachable until Free. Returns p for chaining.
func Retain[T any](p *T) *T {
	if p != nil {
		allocs.LoadOrStore(uintptr(unsafe.Pointer(p)), &allocRec{p: p, live: true})
	}
	return p
}

func allocSize(p unsafe.Pointer) int64 {
	if v, ok := allocs.Load(uintptr(p)); ok {
		rec := v.(*allocRec)
		if rec.live {
			return rec.n
		}
	}
	return 0
}

func mallocCount[T any](n int64) uintptr {
	var z T
	size := unsafe.Sizeof(z)
	if size == 0 || uintptr(n) == size {
		return 1
	}
	return uintptr(n)/size + 1
}

func zeroMalloc[T any](out *T, n int64) {
	var z T
	s := unsafe.Slice(out, mallocCount[T](n))
	for i := range s {
		s[i] = z
	}
}

// Malloc allocates n bytes of memory. It informs the garbage collector that
// the memory will be used to store objects of type T.
func Malloc[T any](n int64) *T {
	if n <= 0 {
		return nil
	}
	k := poolKey{t: reflect.TypeOf((*T)(nil)), n: n}
	if rec := stackFor(k).pop(); rec != nil {
		out := rec.p.(*T)
		rec.n = n
		rec.live = true
		zeroMalloc(out, n)
		return out
	}
	var z T
	var out *T
	size := unsafe.Sizeof(z)
	if size == 0 || uintptr(n) == size {
		out = new(T)
	} else {
		// One extra element to allow indexing off the end, like C tends to.
		out = &make([]T, mallocCount[T](n))[0]
	}
	allocs.Store(uintptr(unsafe.Pointer(out)), &allocRec{p: out, n: n, cap: n, live: true})
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
// Shrinks (0 < n ≤ old) keep p.
func Realloc(p *byte, n int64) *byte {
	if n <= 0 {
		return nil
	}
	if p == nil {
		return Malloc[byte](n)
	}
	oldN := allocSize(unsafe.Pointer(p))
	if oldN > 0 && n <= oldN {
		if v, ok := allocs.Load(uintptr(unsafe.Pointer(p))); ok {
			v.(*allocRec).n = n
		}
		return p
	}
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
	unpin(unsafe.Pointer(p))
	return out
}

// Free is C free(p). Marks the rec idle and stacks it for reuse.
func Free(p *byte) {
	if p != nil {
		unpin(unsafe.Pointer(p))
	}
}

func unpin(p unsafe.Pointer) {
	v, ok := allocs.Load(uintptr(p))
	if !ok {
		return
	}
	rec := v.(*allocRec)
	if !rec.live {
		return
	}
	rec.live = false
	stackFor(poolKey{t: reflect.TypeOf(rec.p), n: recCap(rec)}).push(rec)
}
