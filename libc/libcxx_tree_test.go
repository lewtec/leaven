package libc

import (
	"testing"
	"unsafe"
)

func TestLibcxxTreeGetValueBadPtr(t *testing.T) {
	if LibcxxTreeGetValue(nil) != nil {
		t.Fatal("nil")
	}
	bad := As[byte](unsafe.Pointer(uintptr(3)))
	if LibcxxTreeGetValue(bad) != nil {
		t.Fatal("0x3")
	}
}

func TestLibcxxTreeGetValueOff(t *testing.T) {
	buf := make([]byte, libcxxTreeNodeSize)
	got := LibcxxTreeGetValue(&buf[0])
	want := &buf[libcxxTreeValueOff]
	if got != want {
		t.Fatalf("got %p want %p", got, want)
	}
}

func TestLibcxxTreeConstructSkipsLowChild(t *testing.T) {
	src := make([]byte, libcxxTreeNodeSize)
	// pair written at +0: left=key, right=unsigned 7.
	Store(Ptr(&src[0]), 0, unsafe.Pointer(uintptr(0x1000)))
	Store(Ptr(&src[0]), 8, uint32(7))
	src[libcxxTreeBlackOff] = 1

	dst := LibcxxTreeConstructFromTree(nil, &src[0], nil)
	if dst == nil {
		t.Fatal("dst")
	}
	if Load[unsafe.Pointer](Ptr(dst), libcxxTreeLeftOff) != nil {
		t.Fatal("must not walk key pointer as a child")
	}
	if Load[unsafe.Pointer](Ptr(dst), libcxxTreeRightOff) != nil {
		t.Fatal("right child 0x7 should be skipped")
	}
	if Load[unsafe.Pointer](Ptr(dst), libcxxTreeValueOff) != unsafe.Pointer(uintptr(0x1000)) {
		t.Fatal("key")
	}
	if Load[uint32](Ptr(dst), libcxxTreeValueOff+8) != 7 {
		t.Fatal("val")
	}
	if Load[byte](Ptr(dst), libcxxTreeBlackOff) != 1 {
		t.Fatal("black")
	}
}

func TestLibcxxTreeConstructDoesNotWalkKeyPtr(t *testing.T) {
	// Overlay with unsigned 0: right is nil, left is a heap object
	// whose parent is not src (Variable* in the pair).
	fake := make([]byte, libcxxTreeNodeSize)
	src := make([]byte, libcxxTreeNodeSize)
	Store(Ptr(&src[0]), libcxxTreeLeftOff, Ptr(&fake[0]))

	dst := LibcxxTreeConstructFromTree(nil, &src[0], nil)
	if Load[unsafe.Pointer](Ptr(dst), libcxxTreeLeftOff) != nil {
		t.Fatal("walked key")
	}
	if Load[unsafe.Pointer](Ptr(dst), libcxxTreeValueOff) != Ptr(&fake[0]) {
		t.Fatal("key")
	}
}

func TestLibcxxTreeConstructCopiesChild(t *testing.T) {
	child := make([]byte, libcxxTreeNodeSize)
	src := make([]byte, libcxxTreeNodeSize)
	Store(Ptr(&child[0]), libcxxTreeValueOff+8, uint32(9))
	Store(Ptr(&child[0]), libcxxTreeParentOff, Ptr(&src[0]))
	Store(Ptr(&src[0]), libcxxTreeLeftOff, Ptr(&child[0]))
	Store(Ptr(&src[0]), libcxxTreeValueOff+8, uint32(1))

	dst := LibcxxTreeConstructFromTree(nil, &src[0], nil)
	left := As[byte](Load[unsafe.Pointer](Ptr(dst), libcxxTreeLeftOff))
	if left == nil {
		t.Fatal("left")
	}
	if left == &child[0] {
		t.Fatal("must copy, not alias")
	}
	if Load[uint32](Ptr(left), libcxxTreeValueOff+8) != 9 {
		t.Fatal("child val")
	}
	if Load[unsafe.Pointer](Ptr(left), libcxxTreeParentOff) != Ptr(dst) {
		t.Fatal("parent")
	}
}
