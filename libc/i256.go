package libc

import (
	"math/big"
	"math/bits"
)

// I256 is LLVM i256 as two little-endian I128 limbs.
// rustc core::fmt::num::__fmt_inner widens u128 to i256.
type I256 struct {
	Lo I128
	Hi I128
}

// I256FromU64 is zext of a 64-bit unsigned value to i256.
func I256FromU64(x uint64) I256 { return I256{Lo: I128FromU64(x)} }

// I256FromI64 is sext of a 64-bit signed value to i256.
func I256FromI64(x int64) I256 {
	lo := I128FromI64(x)
	r := I256{Lo: lo}
	if x < 0 {
		r.Hi = I128{Lo: ^uint64(0), Hi: ^uint64(0)}
	}
	return r
}

// I256FromI128 is zext of i128 to i256.
func I256FromI128(x I128) I256 { return I256{Lo: x} }

// I256FromI128S is sext of i128 to i256.
func I256FromI128S(x I128) I256 {
	r := I256{Lo: x}
	if int64(x.Hi) < 0 {
		r.Hi = I128{Lo: ^uint64(0), Hi: ^uint64(0)}
	}
	return r
}

func I256Add(a, b I256) I256 {
	l0, c := bits.Add64(a.Lo.Lo, b.Lo.Lo, 0)
	l1, c := bits.Add64(a.Lo.Hi, b.Lo.Hi, c)
	l2, c := bits.Add64(a.Hi.Lo, b.Hi.Lo, c)
	l3, _ := bits.Add64(a.Hi.Hi, b.Hi.Hi, c)
	return I256{Lo: I128{Lo: l0, Hi: l1}, Hi: I128{Lo: l2, Hi: l3}}
}

func I256Sub(a, b I256) I256 {
	l0, br := bits.Sub64(a.Lo.Lo, b.Lo.Lo, 0)
	l1, br := bits.Sub64(a.Lo.Hi, b.Lo.Hi, br)
	l2, br := bits.Sub64(a.Hi.Lo, b.Hi.Lo, br)
	l3, _ := bits.Sub64(a.Hi.Hi, b.Hi.Hi, br)
	return I256{Lo: I128{Lo: l0, Hi: l1}, Hi: I128{Lo: l2, Hi: l3}}
}

func I256Mul(a, b I256) I256 {
	return i256FromBig(new(big.Int).Mul(i256U(a), i256U(b)))
}

func I256And(a, b I256) I256 { return I256{Lo: I128And(a.Lo, b.Lo), Hi: I128And(a.Hi, b.Hi)} }
func I256Or(a, b I256) I256  { return I256{Lo: I128Or(a.Lo, b.Lo), Hi: I128Or(a.Hi, b.Hi)} }
func I256Xor(a, b I256) I256 { return I256{Lo: I128Xor(a.Lo, b.Lo), Hi: I128Xor(a.Hi, b.Hi)} }

func i256ShAmt(n I256) uint64 { return n.Lo.Lo & 255 }

func I256Shl(a, n I256) I256 {
	s := i256ShAmt(n)
	if s == 0 {
		return a
	}
	if s >= 128 {
		return I256{Hi: I128Shl(a.Lo, I128FromU64(s-128))}
	}
	return I256{
		Lo: I128Shl(a.Lo, I128FromU64(s)),
		Hi: I128Or(I128Shl(a.Hi, I128FromU64(s)), I128LShr(a.Lo, I128FromU64(128-s))),
	}
}

func I256LShr(a, n I256) I256 {
	s := i256ShAmt(n)
	if s == 0 {
		return a
	}
	if s >= 128 {
		return I256{Lo: I128LShr(a.Hi, I128FromU64(s-128))}
	}
	return I256{
		Lo: I128Or(I128LShr(a.Lo, I128FromU64(s)), I128Shl(a.Hi, I128FromU64(128-s))),
		Hi: I128LShr(a.Hi, I128FromU64(s)),
	}
}

func I256AShr(a, n I256) I256 {
	s := i256ShAmt(n)
	if s == 0 {
		return a
	}
	sign := I128FromI64(int64(a.Hi.Hi) >> 63)
	if s >= 128 {
		return I256{Lo: I128AShr(a.Hi, I128FromU64(s-128)), Hi: sign}
	}
	return I256{
		Lo: I128Or(I128LShr(a.Lo, I128FromU64(s)), I128Shl(a.Hi, I128FromU64(128-s))),
		Hi: I128AShr(a.Hi, I128FromU64(s)),
	}
}

func i256U(a I256) *big.Int {
	u := i128U(a.Hi)
	u.Lsh(u, 128).Or(u, i128U(a.Lo))
	return u
}

func i256S(a I256) *big.Int {
	u := i256U(a)
	if a.Hi.Hi>>63 == 1 {
		u.Sub(u, new(big.Int).Lsh(big.NewInt(1), 256))
	}
	return u
}

func i256FromBig(x *big.Int) I256 {
	mod := new(big.Int).Lsh(big.NewInt(1), 256)
	u := new(big.Int).Mod(x, mod)
	lo := new(big.Int).And(u, new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1)))
	hi := new(big.Int).Rsh(u, 128)
	return I256{Lo: i128FromBig(lo), Hi: i128FromBig(hi)}
}

func i256ZeroB(b I256, what string) {
	if b.Lo == (I128{}) && b.Hi == (I128{}) {
		panic("i256 " + what + " by zero")
	}
}

func I256UDiv(a, b I256) I256 {
	i256ZeroB(b, "udiv")
	return i256FromBig(new(big.Int).Div(i256U(a), i256U(b)))
}

func I256SDiv(a, b I256) I256 {
	i256ZeroB(b, "sdiv")
	return i256FromBig(new(big.Int).Quo(i256S(a), i256S(b)))
}

func I256URem(a, b I256) I256 {
	i256ZeroB(b, "urem")
	return i256FromBig(new(big.Int).Rem(i256U(a), i256U(b)))
}

func I256SRem(a, b I256) I256 {
	i256ZeroB(b, "srem")
	return i256FromBig(new(big.Int).Rem(i256S(a), i256S(b)))
}

func I256Eq(a, b I256) bool { return a == b }
func I256Ne(a, b I256) bool { return a != b }

func I256Ult(a, b I256) bool {
	if a.Hi != b.Hi {
		return I128Ult(a.Hi, b.Hi)
	}
	return I128Ult(a.Lo, b.Lo)
}
func I256Ule(a, b I256) bool { return !I256Ult(b, a) }
func I256Ugt(a, b I256) bool { return I256Ult(b, a) }
func I256Uge(a, b I256) bool { return !I256Ult(a, b) }

func I256Slt(a, b I256) bool {
	if a.Hi != b.Hi {
		return I128Slt(a.Hi, b.Hi)
	}
	return I128Ult(a.Lo, b.Lo)
}
func I256Sle(a, b I256) bool { return !I256Slt(b, a) }
func I256Sgt(a, b I256) bool { return I256Slt(b, a) }
func I256Sge(a, b I256) bool { return !I256Slt(a, b) }

func I256TruncI128(a I256) I128 { return a.Lo }
func I256TruncI64(a I256) int64 { return int64(a.Lo.Lo) }
func I256TruncI32(a I256) int32 { return int32(a.Lo.Lo) }
func I256TruncI16(a I256) int16 { return int16(a.Lo.Lo) }
func I256TruncI8(a I256) byte   { return byte(a.Lo.Lo) }
func I256TruncI1(a I256) bool   { return a.Lo.Lo&1 != 0 }
