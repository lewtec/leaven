package libc

import (
	"testing"
	"unsafe"
)

func TestUMinUMax(t *testing.T) {
	if g := UMinU64(-1, 1); g != 1 {
		t.Fatalf("UMinU64(-1,1)=%d want 1", g)
	}
	if g := UMaxU64(-1, 1); g != -1 {
		t.Fatalf("UMaxU64(-1,1)=%d want -1", g)
	}
	if g := UMinU32(-1, 1); g != 1 {
		t.Fatalf("UMinU32(-1,1)=%d want 1", g)
	}
	if g := SMinI64(-2, 5); g != -2 {
		t.Fatalf("SMinI64=%d", g)
	}
	if g := SMaxI32(-2, 5); g != 5 {
		t.Fatalf("SMaxI32=%d", g)
	}
}

func TestVecReduceAddV4I32(t *testing.T) {
	if g := VecReduceAddV4I32([4]int32{1, 2, 3, 4}); g != 10 {
		t.Fatalf("got %d want 10", g)
	}
}

func TestLoadRelativeI64(t *testing.T) {
	// i32 4 at offset 0 → result is p+4.
	var buf [8]byte
	buf[0] = 4
	p := unsafe.Pointer(&buf[0])
	got := LoadRelativeI64(p, 0)
	want := unsafe.Add(p, 4)
	if got != want {
		t.Fatalf("got %p want %p", got, want)
	}
}
