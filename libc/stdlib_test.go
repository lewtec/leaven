package libc

import (
	"sync"
	"testing"
	"unsafe"
)

func TestReallocCopiesOldBytes(t *testing.T) {
	p := Malloc[byte](4)
	copy(unsafe.Slice(p, 4), []byte("abcd"))
	q := Realloc(p, 8)
	got := string(unsafe.Slice(q, 4))
	if got != "abcd" {
		t.Fatalf("Realloc copy = %q, want %q", got, "abcd")
	}
}

func TestFreeUnpins(t *testing.T) {
	p := Malloc[byte](64)
	if p == nil {
		t.Fatal("Malloc returned nil")
	}
	u := uintptr(unsafe.Pointer(p))
	if _, ok := allocs.Load(u); !ok {
		t.Fatal("Malloc did not record pin")
	}
	Free(p)
	if _, ok := allocs.Load(u); ok {
		t.Fatal("Free left pin")
	}
}

func TestReallocUnpinsOld(t *testing.T) {
	p := Malloc[byte](4)
	copy(unsafe.Slice(p, 4), []byte("abcd"))
	old := uintptr(unsafe.Pointer(p))
	q := Realloc(p, 8)
	if q == nil {
		t.Fatal("Realloc returned nil")
	}
	if _, ok := allocs.Load(old); ok {
		t.Fatal("Realloc left pin on old block")
	}
	if _, ok := allocs.Load(uintptr(unsafe.Pointer(q))); !ok {
		t.Fatal("Realloc did not record new pin")
	}
	got := string(unsafe.Slice(q, 4))
	if got != "abcd" {
		t.Fatalf("Realloc copy = %q, want %q", got, "abcd")
	}
	Free(q)
}

func TestRetainThenFree(t *testing.T) {
	p := Malloc[byte](16)
	Retain(p)
	u := uintptr(unsafe.Pointer(p))
	Free(p)
	if _, ok := allocs.Load(u); ok {
		t.Fatal("Free left pin after Retain")
	}
}

func TestFreeNil(t *testing.T) {
	Free(nil)
}

func TestConcurrentMalloc(t *testing.T) {
	const goroutines = 8
	const perG = 1000
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perG; j++ {
				p := Malloc[byte](32)
				if p == nil {
					t.Error("Malloc returned nil")
					return
				}
				*p = byte(j)
				q := Realloc(p, 64)
				if q == nil {
					t.Error("Realloc returned nil")
					return
				}
				if *q != byte(j) {
					t.Errorf("Realloc copy = %d, want %d", *q, byte(j))
					return
				}
				Free(q)
			}
		}()
	}
	wg.Wait()
}

func TestConcurrentMallocFree(t *testing.T) {
	const goroutines = 8
	const perG = 1000
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perG; j++ {
				p := Malloc[byte](32)
				if p == nil {
					t.Error("Malloc returned nil")
					return
				}
				*p = byte(j)
				u := uintptr(unsafe.Pointer(p))
				Free(p)
				if _, ok := allocs.Load(u); ok {
					t.Error("Free left pin")
					return
				}
			}
		}()
	}
	wg.Wait()
}
