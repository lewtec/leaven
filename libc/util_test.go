package libc

import (
	"os"
	"testing"
	"unsafe"
)

func TestPointerKit(t *testing.T) {
	var buf [16]byte
	p := Ptr(&buf[0])
	Store[int32](p, 4, 0x12345678)
	if Load[int32](p, 4) != 0x12345678 {
		t.Fatal("load/store")
	}
	if As[byte](p) != &buf[0] {
		t.Fatal("as")
	}
	if Addr(&buf[0]) != uintptr(p) {
		t.Fatal("addr")
	}
	if Off(p, 4) != Ptr(&buf[4]) {
		t.Fatal("off")
	}
	if AddPointer(&buf[0], 3) != &buf[3] {
		t.Fatal("addpointer")
	}
	copy(Bytes(&buf[8], 3), []byte("hi"))
	if buf[8] != 'h' || buf[9] != 'i' {
		t.Fatal("bytes")
	}
	if Bytes(nil, 0) != nil || Bytes(&buf[0], 0) != nil {
		t.Fatal("bytes empty")
	}
	if GoString(nil) != "" {
		t.Fatal("gostring nil")
	}
	hit := 0
	fn := func(p unsafe.Pointer) { hit = int(Load[byte](p, 0)) }
	var x byte = 9
	FuncFromCode[func(unsafe.Pointer)](FuncCode(fn))(Ptr(&x))
	if hit != 9 {
		t.Fatalf("funccode hit=%d", hit)
	}
}

func TestArgv(t *testing.T) {
	p := Argv()
	if p == nil {
		t.Fatal("nil argv")
	}
	tab := unsafe.Slice((*uint64)(p), len(os.Args)+1)
	if tab[len(os.Args)] != 0 {
		t.Fatal("missing trailing nil")
	}
	got := GoString(As[byte](PtrFromBits(tab[0])))
	if got != os.Args[0] {
		t.Fatalf("argv[0]=%q want %q", got, os.Args[0])
	}
}
