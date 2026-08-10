package libc

import (
	"sync"
	"testing"
	"unsafe"
)

func recLive(u uintptr) bool {
	v, ok := allocs.Load(u)
	return ok && v.(*allocRec).live
}

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
	if !recLive(u) {
		t.Fatal("Malloc did not record live pin")
	}
	Free(p)
	if recLive(u) {
		t.Fatal("Free left rec live")
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
	if recLive(old) {
		t.Fatal("Realloc left old block live")
	}
	if !recLive(uintptr(unsafe.Pointer(q))) {
		t.Fatal("Realloc did not record live new pin")
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
	if recLive(u) {
		t.Fatal("Free left rec live after Retain")
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

func TestMallocReusesSameType(t *testing.T) {
	p := Malloc[byte](64)
	*p = 7
	addr := uintptr(unsafe.Pointer(p))
	Free(p)
	hits := 0
	for i := 0; i < 32; i++ {
		q := Malloc[byte](64)
		if uintptr(unsafe.Pointer(q)) == addr {
			hits++
			if *q != 0 {
				t.Fatalf("reused block not zeroed: %d", *q)
			}
		}
		Free(q)
	}
	if hits == 0 {
		t.Fatal("expected at least one reuse of the freed 64-byte block")
	}
}

func TestMallocDoesNotReuseAcrossTypes(t *testing.T) {
	type node struct{ a, b *byte }
	p := Malloc[node](int64(unsafe.Sizeof(node{})))
	addr := uintptr(unsafe.Pointer(p))
	Free((*byte)(unsafe.Pointer(p)))
	for i := 0; i < 32; i++ {
		q := Malloc[byte](int64(unsafe.Sizeof(node{})))
		if uintptr(unsafe.Pointer(q)) == addr {
			t.Fatal("pooled node returned as *byte")
		}
		Free(q)
	}
}

func TestReallocShrinkInPlace(t *testing.T) {
	p := Malloc[byte](16)
	copy(unsafe.Slice(p, 16), []byte("0123456789abcdef"))
	q := Realloc(p, 8)
	if q != p {
		t.Fatal("shrink allocated a new block")
	}
	if got := string(unsafe.Slice(q, 8)); got != "01234567" {
		t.Fatalf("shrink data = %q", got)
	}
	r := Realloc(q, 32)
	if r == nil {
		t.Fatal("grow returned nil")
	}
	if got := string(unsafe.Slice(r, 8)); got != "01234567" {
		t.Fatalf("grow after shrink = %q", got)
	}
	Free(r)
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
				Free(p)
			}
		}()
	}
	wg.Wait()
}
