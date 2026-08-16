package libc

import "testing"

func TestAtomicAddI8(t *testing.T) {
	var b [4]byte
	b[1] = 3
	got := AtomicAddI8(&b[1], 255) // wrap -1
	if got != 2 || b[1] != 2 {
		t.Fatalf("add -1: new=%d cell=%d", got, b[1])
	}
	if b[0] != 0 || b[2] != 0 {
		t.Fatalf("neighbors %v", b)
	}
}

func TestAtomicSwapI8(t *testing.T) {
	var b [8]byte
	// Put the target byte at every alignment offset.
	for off := 0; off < 4; off++ {
		b[off] = 0xab
		old := AtomicSwapI8(&b[off], 0xcd)
		if old != 0xab {
			t.Fatalf("off %d: old=%#x", off, old)
		}
		if b[off] != 0xcd {
			t.Fatalf("off %d: got %#x", off, b[off])
		}
		// Neighbors untouched when they were zero.
		for i := range b {
			if i == off {
				continue
			}
			if b[i] != 0 {
				t.Fatalf("off %d: neighbor %d = %#x", off, i, b[i])
			}
		}
		b[off] = 0
	}
}

func TestAtomicCASI8(t *testing.T) {
	var b [4]byte
	b[1] = 7
	if AtomicCASI8(&b[1], 8, 9) {
		t.Fatal("CAS should fail on mismatch")
	}
	if b[1] != 7 {
		t.Fatalf("value changed on failed CAS: %d", b[1])
	}
	if !AtomicCASI8(&b[1], 7, 9) {
		t.Fatal("CAS should succeed")
	}
	if b[1] != 9 {
		t.Fatalf("after CAS: %d", b[1])
	}
	if AtomicLoadI8(&b[1]) != 9 {
		t.Fatal("load")
	}
}
