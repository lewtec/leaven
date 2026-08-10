package libc

import "unsafe"

// UMinU64 is llvm.umin.i64 (unsigned compare, i64 bit pattern).
func UMinU64(a, b int64) int64 {
	if uint64(a) < uint64(b) {
		return a
	}
	return b
}

// UMaxU32 is llvm.umax.i32.
func UMaxU32(a, b int32) int32 {
	if uint32(a) > uint32(b) {
		return a
	}
	return b
}

// UMinU32 is llvm.umin.i32.
func UMinU32(a, b int32) int32 {
	if uint32(a) < uint32(b) {
		return a
	}
	return b
}

// SMaxI64 is llvm.smax.i64.
func SMaxI64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// SMinI64 is llvm.smin.i64.
func SMinI64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// SMaxI32 is llvm.smax.i32.
func SMaxI32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// SMinI32 is llvm.smin.i32.
func SMinI32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// VecReduceAddV4I32 is llvm.vector.reduce.add.v4i32.
func VecReduceAddV4I32(v [4]int32) int32 {
	return v[0] + v[1] + v[2] + v[3]
}

// LoadRelativeI64 is llvm.load.relative.i64: load i32 at ptr+offset, return ptr+i32.
func LoadRelativeI64(p unsafe.Pointer, off int64) unsafe.Pointer {
	rel := int64(*(*int32)(unsafe.Add(p, int(off))))
	return unsafe.Add(p, int(rel))
}
