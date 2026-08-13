package libc

func i1Pack[M ~uint8 | ~uint16 | ~uint32 | ~uint64](v []bool) M {
	var m M
	for i, b := range v {
		if b {
			m |= 1 << uint(i)
		}
	}
	return m
}

func i1Unpack[M ~uint8 | ~uint16 | ~uint32 | ~uint64](m M, v []bool) {
	for i := range v {
		v[i] = m&(1<<uint(i)) != 0
	}
}

// I1Pack8 packs <8 x i1> (lane 0 = low bit) to i8. rustc SIMD icmp masks.
func I1Pack8(v [8]bool) byte { return i1Pack[uint8](v[:]) }

// I1Unpack8 is bitcast i8 to <8 x i1>.
func I1Unpack8(m byte) (v [8]bool) {
	i1Unpack(m, v[:])
	return v
}

// I1Pack16 packs <16 x i1> to i16.
func I1Pack16(v [16]bool) int16 { return int16(i1Pack[uint16](v[:])) }

// I1Unpack16 is bitcast i16 to <16 x i1>.
func I1Unpack16(m int16) (v [16]bool) {
	i1Unpack(uint16(m), v[:])
	return v
}

// I1Pack32 packs <32 x i1> to i32.
func I1Pack32(v [32]bool) int32 { return int32(i1Pack[uint32](v[:])) }

// I1Unpack32 is bitcast i32 to <32 x i1>.
func I1Unpack32(m int32) (v [32]bool) {
	i1Unpack(uint32(m), v[:])
	return v
}

// I1Pack64 packs <64 x i1> to i64.
func I1Pack64(v [64]bool) int64 { return int64(i1Pack[uint64](v[:])) }

// I1Unpack64 is bitcast i64 to <64 x i1>.
func I1Unpack64(m int64) (v [64]bool) {
	i1Unpack(uint64(m), v[:])
	return v
}
