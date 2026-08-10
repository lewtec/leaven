package libc

import (
	"os"
	"testing"
	"unsafe"
)

func TestArgv(t *testing.T) {
	p := Argv()
	if p == nil {
		t.Fatal("nil argv")
	}
	tab := unsafe.Slice((**byte)(p), len(os.Args)+1)
	if tab[len(os.Args)] != nil {
		t.Fatal("missing trailing nil")
	}
	if GoString(tab[0]) != os.Args[0] {
		t.Fatalf("argv[0]=%q want %q", GoString(tab[0]), os.Args[0])
	}
}
