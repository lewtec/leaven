package libc

import "testing"

func TestFence(t *testing.T) {
	before := asmMemFence
	Fence()
	if asmMemFence == before {
		t.Fatal("fence did not touch the barrier")
	}
}
