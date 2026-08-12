package libc

import (
	"testing"
	"unsafe"
)

var cxaHit byte

func cxaTestDtor(p unsafe.Pointer) {
	cxaHit = 1
	if p != nil {
		*(*byte)(p) = 7
	}
}

func TestCxaAtexit(t *testing.T) {
	cxaHit = 0
	tmp := cxaTestDtor
	fn := *(**byte)(unsafe.Pointer(&tmp))
	var obj byte
	if CxaAtexit(fn, &obj, nil) != 0 {
		t.Fatal("ret")
	}
	runCxaAtExit()
	if cxaHit != 1 || obj != 7 {
		t.Fatalf("hit=%d obj=%d", cxaHit, obj)
	}
}
