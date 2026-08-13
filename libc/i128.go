package libc

import (
	"math/big"
	"math/bits"
)

// I128 is LLVM i128 as two little-endian two's-complement limbs.
// rustc TypeId and i128 math use the full width; this is not int64.
type I128 struct {
	Lo uint64
	Hi uint64
}

// I128FromU64 is zext of a 64-bit unsigned value to i128.
func I128FromU64(x uint64) I128 { return I128{Lo: x} }

// I128FromI64 is sext of a 64-bit signed value to i128.
func I128FromI64(x int64) I128 {
	r := I128{Lo: uint64(x)}
	if x < 0 {
		r.Hi = ^uint64(0)
	}
	return r
}

func I128Add(a, b I128) I128 {
	lo, c := bits.Add64(a.Lo, b.Lo, 0)
	hi, _ := bits.Add64(a.Hi, b.Hi, c)
	return I128{Lo: lo, Hi: hi}
}

func I128Sub(a, b I128) I128 {
	lo, br := bits.Sub64(a.Lo, b.Lo, 0)
	hi, _ := bits.Sub64(a.Hi, b.Hi, br)
	return I128{Lo: lo, Hi: hi}
}

func I128Mul(a, b I128) I128 {
	hi, lo := bits.Mul64(a.Lo, b.Lo)
	hi += a.Lo*b.Hi + a.Hi*b.Lo
	return I128{Lo: lo, Hi: hi}
}

func I128And(a, b I128) I128 { return I128{Lo: a.Lo & b.Lo, Hi: a.Hi & b.Hi} }
func I128Or(a, b I128) I128  { return I128{Lo: a.Lo | b.Lo, Hi: a.Hi | b.Hi} }
func I128Xor(a, b I128) I128 { return I128{Lo: a.Lo ^ b.Lo, Hi: a.Hi ^ b.Hi} }

func i128ShAmt(n I128) uint64 { return n.Lo & 127 }

func I128Shl(a, n I128) I128 {
	s := i128ShAmt(n)
	if s == 0 {
		return a
	}
	if s >= 64 {
		return I128{Lo: 0, Hi: a.Lo << (s - 64)}
	}
	return I128{Lo: a.Lo << s, Hi: a.Hi<<s | a.Lo>>(64-s)}
}

func I128LShr(a, n I128) I128 {
	s := i128ShAmt(n)
	if s == 0 {
		return a
	}
	if s >= 64 {
		return I128{Lo: a.Hi >> (s - 64), Hi: 0}
	}
	return I128{Lo: a.Lo>>s | a.Hi<<(64-s), Hi: a.Hi >> s}
}

func I128AShr(a, n I128) I128 {
	s := i128ShAmt(n)
	if s == 0 {
		return a
	}
	sign := int64(a.Hi)
	if s >= 64 {
		return I128{Lo: uint64(sign >> (s - 64)), Hi: uint64(sign >> 63)}
	}
	return I128{Lo: a.Lo>>s | a.Hi<<(64-s), Hi: uint64(sign >> s)}
}

func i128U(a I128) *big.Int {
	u := new(big.Int).SetUint64(a.Hi)
	u.Lsh(u, 64).Or(u, new(big.Int).SetUint64(a.Lo))
	return u
}

func i128S(a I128) *big.Int {
	return bigAsSigned(i128U(a), 128, a.Hi>>63 == 1)
}

func i128FromBig(x *big.Int) I128 {
	mod := new(big.Int).Lsh(big.NewInt(1), 128)
	u := new(big.Int).Mod(x, mod)
	lo := new(big.Int).And(u, new(big.Int).SetUint64(^uint64(0))).Uint64()
	hi := new(big.Int).Rsh(u, 64).Uint64()
	return I128{Lo: lo, Hi: hi}
}

var i128Big = wideBig[I128]{name: "i128", toU: i128U, toS: i128S, from: i128FromBig}

func I128UDiv(a, b I128) I128 { return i128Big.bin(a, b, false, (*big.Int).Div, "udiv") }
func I128SDiv(a, b I128) I128 { return i128Big.bin(a, b, true, (*big.Int).Quo, "sdiv") }
func I128URem(a, b I128) I128 { return i128Big.bin(a, b, false, (*big.Int).Rem, "urem") }
func I128SRem(a, b I128) I128 { return i128Big.bin(a, b, true, (*big.Int).Rem, "srem") }

func I128Eq(a, b I128) bool { return a == b }
func I128Ne(a, b I128) bool { return a != b }

func I128Ult(a, b I128) bool {
	if a.Hi != b.Hi {
		return a.Hi < b.Hi
	}
	return a.Lo < b.Lo
}
func I128Ule(a, b I128) bool { return !I128Ult(b, a) }
func I128Ugt(a, b I128) bool { return I128Ult(b, a) }
func I128Uge(a, b I128) bool { return !I128Ult(a, b) }

func I128Slt(a, b I128) bool {
	if int64(a.Hi) != int64(b.Hi) {
		return int64(a.Hi) < int64(b.Hi)
	}
	return a.Lo < b.Lo
}
func I128Sle(a, b I128) bool { return !I128Slt(b, a) }
func I128Sgt(a, b I128) bool { return I128Slt(b, a) }
func I128Sge(a, b I128) bool { return !I128Slt(a, b) }

func I128TruncI64(a I128) int64 { return int64(a.Lo) }
func I128TruncI32(a I128) int32 { return int32(a.Lo) }
func I128TruncI16(a I128) int16 { return int16(a.Lo) }
func I128TruncI8(a I128) byte   { return byte(a.Lo) }
func I128TruncI1(a I128) bool   { return a.Lo&1 != 0 }
