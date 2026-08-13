package libc

import (
	"math"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestMaximumNumF64(t *testing.T) {
	if g := MaximumNumF64(1, 2); g != 2 {
		t.Fatalf("1,2 = %v", g)
	}
	if g := MaximumNumF64(-3, 1); g != 1 {
		t.Fatalf("-3,1 = %v", g)
	}
	nan := math.NaN()
	if g := MaximumNumF64(nan, 4); g != 4 {
		t.Fatalf("NaN,4 = %v", g)
	}
	if g := MaximumNumF64(4, nan); g != 4 {
		t.Fatalf("4,NaN = %v", g)
	}
	if !math.IsNaN(MaximumNumF64(nan, nan)) {
		t.Fatal("NaN,NaN")
	}
	pz, nz := 0.0, math.Copysign(0, -1)
	if g := MaximumNumF64(nz, pz); math.Signbit(g) || g != 0 {
		t.Fatalf("-0,+0 want +0, got %v sign=%v", g, math.Signbit(g))
	}
	if g := MaximumNumF64(pz, nz); math.Signbit(g) || g != 0 {
		t.Fatalf("+0,-0 want +0, got %v sign=%v", g, math.Signbit(g))
	}
	if g := MaximumNumF64(nz, nz); !math.Signbit(g) || g != 0 {
		t.Fatalf("-0,-0 want -0, got %v sign=%v", g, math.Signbit(g))
	}
}

func TestFence(t *testing.T) {
	before := asmMemFence
	Fence()
	if asmMemFence == before {
		t.Fatal("fence did not touch the barrier")
	}
}

func TestStatxFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "f")
	if err := os.WriteFile(p, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	path := append([]byte(p), 0)
	var buf [256]byte
	if Statx(-100, &path[0], 0, 0x7ff, &buf[0]) != 0 {
		t.Fatal("statx failed")
	}
	size := *(*uint64)(unsafe.Pointer(&buf[40]))
	if size != 5 {
		t.Fatalf("stx_size=%d", size)
	}
	if Statx(-100, &path[0], 0, 0x7ff, nil) != -1 {
		t.Fatal("nil buf")
	}
}

func TestStatxEmptyPath(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "f")
	if err := os.WriteFile(p, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	empty := []byte{0}
	var buf [256]byte
	const atEmptyPath = 0x1000
	if Statx(int32(f.Fd()), &empty[0], atEmptyPath, 0x7ff, &buf[0]) != 0 {
		t.Fatal("statx AT_EMPTY_PATH failed")
	}
	size := *(*uint64)(unsafe.Pointer(&buf[40]))
	if size != 5 {
		t.Fatalf("stx_size=%d", size)
	}
}
