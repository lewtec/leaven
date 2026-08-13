package libc

import "testing"

func TestI1Pack16(t *testing.T) {
	var v [16]bool
	for i := 0; i < 16; i++ {
		v[i] = true
	}
	if I1Pack16(v) != -1 {
		t.Fatalf("all-true %d", I1Pack16(v))
	}
	v[0] = false
	if got := uint16(I1Pack16(v)); got != 0xfffe {
		t.Fatalf("lane0 clear %#x", got)
	}
	if I1Unpack16(-1) != [16]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true} {
		t.Fatal("unpack all-ones")
	}
}

func TestI1Pack8(t *testing.T) {
	var v [8]bool
	v[0], v[7] = true, true
	if I1Pack8(v) != 0x81 {
		t.Fatalf("got %#x", I1Pack8(v))
	}
	if I1Unpack8(0x81) != v {
		t.Fatal("unpack")
	}
}
