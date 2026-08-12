package libc

import "testing"

func TestI128AddMulCmp(t *testing.T) {
	one := I128FromI64(1)
	sum := I128Add(I128FromI64(2), one)
	if sum != I128FromI64(3) {
		t.Fatalf("2+1=%v", sum)
	}
	// wrap at 2^128: max+1 == min signed
	max := I128{Lo: ^uint64(0), Hi: ^uint64(0) >> 1}
	wrapped := I128Add(max, one)
	if wrapped.Hi>>63 == 0 {
		t.Fatalf("max+1 should set sign, got %v", wrapped)
	}
	prod := I128Mul(I128FromU64(^uint64(0)), I128FromU64(2))
	if prod.Lo != ^uint64(0)-1 || prod.Hi != 1 {
		t.Fatalf("umul %v", prod)
	}
	if !I128Ult(I128FromU64(1), I128FromU64(2)) || I128Ult(I128FromU64(2), I128FromU64(1)) {
		t.Fatal("ult")
	}
	neg := I128FromI64(-1)
	if !I128Slt(neg, I128FromI64(0)) || I128Ult(neg, I128FromI64(0)) {
		t.Fatal("slt vs ult on -1")
	}
	if I128TruncI64(I128Add(I128FromI64(-5), I128FromI64(2))) != -3 {
		t.Fatal("sext add")
	}
}

func TestI128ShiftDiv(t *testing.T) {
	x := I128{Lo: 1, Hi: 1}
	sh := I128Shl(x, I128FromU64(64))
	if sh.Lo != 0 || sh.Hi != 1 {
		t.Fatalf("shl64 %v", sh)
	}
	if got := I128LShr(sh, I128FromU64(64)); got != I128FromU64(1) {
		t.Fatalf("lshr64 %v", got)
	}
	neg := I128FromI64(-16)
	if got := I128AShr(neg, I128FromU64(2)); I128TruncI64(got) != -4 {
		t.Fatalf("ashr %v", got)
	}
	if I128UDiv(I128FromU64(100), I128FromU64(7)) != I128FromU64(14) {
		t.Fatal("udiv")
	}
	if I128SDiv(I128FromI64(-100), I128FromI64(7)) != I128FromI64(-14) {
		t.Fatal("sdiv")
	}
}
