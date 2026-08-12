package libc

// I1Pack8 packs <8 x i1> (lane 0 = low bit) to i8. rustc SIMD icmp masks.
func I1Pack8(v [8]bool) byte {
	var m uint8
	for i, b := range v {
		if b {
			m |= 1 << uint(i)
		}
	}
	return m
}

// I1Unpack8 is bitcast i8 to <8 x i1>.
func I1Unpack8(m byte) (v [8]bool) {
	for i := range v {
		v[i] = m&(1<<uint(i)) != 0
	}
	return v
}

// I1Pack16 packs <16 x i1> to i16.
func I1Pack16(v [16]bool) int16 {
	var m uint16
	for i, b := range v {
		if b {
			m |= 1 << uint(i)
		}
	}
	return int16(m)
}

// I1Unpack16 is bitcast i16 to <16 x i1>.
func I1Unpack16(m int16) (v [16]bool) {
	u := uint16(m)
	for i := range v {
		v[i] = u&(1<<uint(i)) != 0
	}
	return v
}

// I1Pack32 packs <32 x i1> to i32.
func I1Pack32(v [32]bool) int32 {
	var m uint32
	for i, b := range v {
		if b {
			m |= 1 << uint(i)
		}
	}
	return int32(m)
}

// I1Unpack32 is bitcast i32 to <32 x i1>.
func I1Unpack32(m int32) (v [32]bool) {
	u := uint32(m)
	for i := range v {
		v[i] = u&(1<<uint(i)) != 0
	}
	return v
}

// I1Pack64 packs <64 x i1> to i64.
func I1Pack64(v [64]bool) int64 {
	var m uint64
	for i, b := range v {
		if b {
			m |= 1 << uint(i)
		}
	}
	return int64(m)
}

// I1Unpack64 is bitcast i64 to <64 x i1>.
func I1Unpack64(m int64) (v [64]bool) {
	u := uint64(m)
	for i := range v {
		v[i] = u&(1<<uint(i)) != 0
	}
	return v
}
