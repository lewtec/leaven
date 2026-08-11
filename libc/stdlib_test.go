package libc

import (
	"runtime"
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
	Free(q)
}

func TestReallocNilIsMalloc(t *testing.T) {
	q := Realloc(nil, 8)
	if q == nil {
		t.Fatal("Realloc(nil, 8) returned nil")
	}
	*q = 1
	Free(q)
}

func TestReallocZeroFrees(t *testing.T) {
	p := Malloc[byte](8)
	if Realloc(p, 0) != nil {
		t.Fatal("Realloc(p, 0) should return nil")
	}
}

func TestFreeNil(t *testing.T) {
	Free(nil)
}

func TestArc4randomBufFills(t *testing.T) {
	var a, b [16]byte
	Arc4randomBuf(unsafe.Pointer(&a[0]), 16)
	Arc4randomBuf(unsafe.Pointer(&b[0]), 16)
	if a == b {
		t.Fatal("two fills were identical")
	}
	Arc4randomBuf(nil, 8)
	Arc4randomBuf(unsafe.Pointer(&a[0]), 0)
}

func TestMallocZeroReturnsUnique(t *testing.T) {
	a := Malloc[byte](0)
	b := Malloc[byte](0)
	if a == nil || b == nil {
		t.Fatal("malloc(0) returned nil")
	}
	if a == b {
		t.Fatal("malloc(0) returned the same pointer twice")
	}
	Free(a)
	Free(b)
}

func TestCallocZeros(t *testing.T) {
	p := Calloc[byte](8, 1)
	for i, b := range unsafe.Slice(p, 8) {
		if b != 0 {
			t.Fatalf("Calloc byte %d = %d", i, b)
		}
	}
	Free(p)
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

func TestRetainGoHeap(t *testing.T) {
	p := new(byte)
	*p = 9
	if Retain(p) != p {
		t.Fatal("Retain did not return p")
	}
	if _, ok := allocs.Load(uintptr(unsafe.Pointer(p))); !ok {
		t.Fatal("Retain did not pin Go heap object")
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
				Free(p)
			}
		}()
	}
	wg.Wait()
}

func TestTypedListSurvivesGC(t *testing.T) {
	type node struct {
		x    int64
		next *node
	}
	var head *node
	for i := 0; i < 100; i++ {
		n := Malloc[node](int64(unsafe.Sizeof(node{})))
		n.x = int64(i)
		n.next = head
		head = n
	}
	runtime.GC()
	n := 0
	for p := head; p != nil; p = p.next {
		n++
	}
	if n != 100 {
		t.Fatalf("list length = %d after GC", n)
	}
	for p := head; p != nil; {
		next := p.next
		Free((*byte)(unsafe.Pointer(p)))
		p = next
	}
}
