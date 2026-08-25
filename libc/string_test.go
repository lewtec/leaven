package libc

import (
	"runtime"
	"testing"
	"unsafe"
)

func TestWmemchr(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("wchar_t is 2 bytes")
	}
	ws := []int32{'a', 'b', 'c', 0}
	p := Wmemchr((*byte)(unsafe.Pointer(&ws[0])), 'b', 3)
	if p == nil {
		t.Fatal("miss")
	}
	if Load[int32](Ptr(p), 0) != 'b' {
		t.Fatal("value")
	}
	if Wmemchr((*byte)(unsafe.Pointer(&ws[0])), 'z', 3) != nil {
		t.Fatal("absent")
	}
}
