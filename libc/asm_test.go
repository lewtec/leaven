package libc

import "testing"

func TestInlineAsmMemoryBarrier(t *testing.T) {
	before := asmMemFence
	InlineAsm("", "~{memory}")
	if asmMemFence == before {
		t.Fatal("memory barrier did not touch the fence")
	}
}

func TestInlineAsmUnknownPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	InlineAsm("int $0x80", "")
}
