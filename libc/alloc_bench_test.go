package libc

import (
	"testing"
	"unsafe"
)

// Allocator benches. modernc.org/memory slab: 0 Go allocs after warmup.

type benchNode struct {
	F0, F1 *benchNode
}

func BenchmarkMallocFree64(b *testing.B) {
	// Keep a page live. Alloc+free of one slot unmaps the slab page.
	const live = 64
	ps := make([]*byte, live)
	for i := range ps {
		ps[i] = Malloc[byte](64)
	}
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		Free(ps[i])
		ps[i] = Malloc[byte](64)
		i++
		if i == live {
			i = 0
		}
	}
	for _, p := range ps {
		Free(p)
	}
}

func BenchmarkMallocFreeTyped16(b *testing.B) {
	const live = 64
	sz := int64(unsafe.Sizeof(benchNode{}))
	ps := make([]*benchNode, live)
	for i := range ps {
		ps[i] = Malloc[benchNode](sz)
	}
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		Free((*byte)(unsafe.Pointer(ps[i])))
		ps[i] = Malloc[benchNode](sz)
		i++
		if i == live {
			i = 0
		}
	}
	for _, p := range ps {
		Free((*byte)(unsafe.Pointer(p)))
	}
}

func BenchmarkReallocGrow32to64(b *testing.B) {
	const live = 64
	ps := make([]*byte, live)
	for i := range ps {
		ps[i] = Malloc[byte](64)
	}
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		Free(ps[i])
		p := Malloc[byte](32)
		ps[i] = Realloc(p, 64)
		i++
		if i == live {
			i = 0
		}
	}
	for _, p := range ps {
		Free(p)
	}
}

func BenchmarkChurn10k64(b *testing.B) {
	const n = 10000
	ps := make([]*byte, n)
	b.ReportAllocs()
	for b.Loop() {
		for i := range ps {
			ps[i] = Malloc[byte](64)
		}
		for i := range ps {
			Free(ps[i])
		}
	}
}

func BenchmarkMallocFree64Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p := Malloc[byte](64)
			Free(p)
		}
	})
}
