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
			}
		}()
	}
	wg.Wait()
}
