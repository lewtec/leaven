package libc

import "math/big"

// wideBig is i128/i256 math that goes through math/big (div, rem, mul).
type wideBig[T comparable] struct {
	name string
	toU  func(T) *big.Int
	toS  func(T) *big.Int
	from func(*big.Int) T
}

func (w wideBig[T]) bin(a, b T, signed bool, op func(z, x, y *big.Int) *big.Int, what string) T {
	var z T
	if b == z {
		panic(w.name + " " + what + " by zero")
	}
	x, y := w.toU(a), w.toU(b)
	if signed {
		x, y = w.toS(a), w.toS(b)
	}
	return w.from(op(new(big.Int), x, y))
}

func (w wideBig[T]) mul(a, b T) T {
	return w.from(new(big.Int).Mul(w.toU(a), w.toU(b)))
}

func bigAsSigned(u *big.Int, bits uint, neg bool) *big.Int {
	if neg {
		u.Sub(u, new(big.Int).Lsh(big.NewInt(1), bits))
	}
	return u
}
