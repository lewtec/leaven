package libc

import "testing"

func TestI256AddMulCmp(t *testing.T) {
	one := I256FromI64(1)
	sum := I256Add(I256FromI64(2), one)
	if sum != I256FromI64(3) {
		t.Fatalf("2+1=%v", sum)
	}
	wide := I256FromI128(I128{Lo: 0, Hi: 1})
	sum2 := I256Add(wide, I256FromU64(1))
	if sum2.Lo.Lo != 1 || sum2.Lo.Hi != 1 {
		t.Fatalf("2^64+1=%v", sum2)
	}
	prod := I256Mul(I256FromI128(I128{Lo: ^uint64(0), Hi: 0}), I256FromU64(2))
	if prod.Lo.Lo != ^uint64(0)-1 || prod.Lo.Hi != 1 {
		t.Fatalf("umul %v", prod)
	}
	if !I256Ult(I256FromU64(1), I256FromU64(2)) || I256Ult(I256FromU64(2), I256FromU64(1)) {
		t.Fatal("ult")
	}
	neg := I256FromI64(-1)
	if !I256Slt(neg, I256FromI64(0)) || I256Ult(neg, I256FromI64(0)) {
		t.Fatal("slt vs ult on -1")
	}
	if I256TruncI64(I256Add(I256FromI64(-5), I256FromI64(2))) != -3 {
		t.Fatal("sext add")
	}
	if I256TruncI128(I256FromI128(I128FromU64(9))) != I128FromU64(9) {
		t.Fatal("trunc i128")
	}
}

func TestI256ShiftDiv(t *testing.T) {
	x := I256FromI128(I128{Lo: 1, Hi: 1})
	sh := I256Shl(x, I256FromU64(128))
	if sh.Lo != (I128{}) || sh.Hi.Lo != 1 || sh.Hi.Hi != 1 {
		t.Fatalf("shl128 %v", sh)
	}
	if got := I256LShr(sh, I256FromU64(128)); got != x {
		t.Fatalf("lshr128 %v", got)
	}
	neg := I256FromI64(-16)
	if got := I256AShr(neg, I256FromU64(2)); I256TruncI64(got) != -4 {
		t.Fatalf("ashr %v", got)
	}
	if I256UDiv(I256FromU64(100), I256FromU64(7)) != I256FromU64(14) {
		t.Fatal("udiv")
	}
	if I256SDiv(I256FromI64(-100), I256FromI64(7)) != I256FromI64(-14) {
		t.Fatal("sdiv")
	}
	// zext i128 then urem: core::fmt::num path
	n := I256FromI128(I128FromU64(10000))
	if I256URem(n, I256FromU64(100)) != I256FromU64(0) {
		t.Fatal("urem")
	}
}
