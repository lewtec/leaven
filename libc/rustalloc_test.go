package libc

import (
	"math"
	"testing"
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
