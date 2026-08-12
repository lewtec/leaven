package libc

import (
	"testing"
	"unsafe"
)

func TestSscanfLu(t *testing.T) {
	var v uint64
	n := Sscanf(cbyte("1"), cbyte("%lu"), (*byte)(unsafe.Pointer(&v)))
	if n != 1 || v != 1 {
		t.Fatalf("n=%d v=%d", n, v)
	}
	var w uint64
	n = Sscanf(cbyte("42"), cbyte("%lu"), &w)
	if n != 1 || w != 42 {
		t.Fatalf("typed n=%d w=%d", n, w)
	}
	var z uint64 = 7
	n = Sscanf(cbyte("nope"), cbyte("%lu"), &z)
	if n != 0 || z != 7 {
		t.Fatalf("fail n=%d z=%d", n, z)
	}
}

func cbyte(s string) *byte {
	b := append([]byte(s), 0)
	return &b[0]
}
