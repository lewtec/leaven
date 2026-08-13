package libc

import "testing"

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
