package libc

import (
	"testing"
	"unsafe"
)

// Baseline allocator benches. No reuse yet: each Malloc is a new Go alloc plus
// a sync.Map pin; Free only drops the pin.

type benchNode struct {
	F0, F1 *benchNode
}

func BenchmarkMallocFree64(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		p := Malloc[byte](64)
		Free(p)
	}
}

func BenchmarkMallocFreeTyped16(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		p := Malloc[benchNode](int64(unsafe.Sizeof(benchNode{})))
		Free((*byte)(unsafe.Pointer(p)))
	}
}

func BenchmarkReallocGrow32to64(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		p := Malloc[byte](32)
		q := Realloc(p, 64)
		Free(q)
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
