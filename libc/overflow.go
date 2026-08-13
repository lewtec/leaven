package libc

// USubSatU64 is llvm.usub.sat for integer widths ≤64.
func USubSatU64(a, b uint64) uint64 {
	if a < b {
		return 0
	}
	return a - b
}

// UAddSatU64 is llvm.uadd.sat for integer widths ≤64.
func UAddSatU64(a, b uint64) uint64 {
	s := a + b
	if s < a {
		return ^uint64(0)
	}
	return s
}

// AbsI64 is llvm.abs.i64 (INT64_MIN stays).
func AbsI64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// AbsI32 is llvm.abs.i32. INT_MIN stays INT_MIN (two's complement wrap).
func AbsI32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// SAddWithOverflowI32 is llvm.sadd.with.overflow.i32: wrapping sum and signed overflow.
func SAddWithOverflowI32(a, b int32) struct {
	F0 int32
	F1 bool
} {
	s := int32(uint32(a) + uint32(b))
	ov := (a >= 0) == (b >= 0) && (s < 0) != (a < 0)
	return struct {
		F0 int32
		F1 bool
	}{s, ov}
}
