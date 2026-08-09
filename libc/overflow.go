package libc

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
