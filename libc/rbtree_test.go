package libc

import (
	"testing"
	"unsafe"
)

func TestRbNodeLayout(t *testing.T) {
	if unsafe.Offsetof(rbNode{}.color) != rbColorOff {
		t.Fatalf("color")
	}
	if unsafe.Offsetof(rbNode{}.parent) != rbParentOff {
		t.Fatalf("parent")
	}
	if unsafe.Offsetof(rbNode{}.left) != rbLeftOff {
		t.Fatalf("left")
	}
	if unsafe.Offsetof(rbNode{}.right) != rbRightOff {
		t.Fatalf("right")
	}
}

func TestRbTreeIncDecHeader(t *testing.T) {
	// header (end) + one black root. decrement(end) is the root.
	var hdr, root rbNode
	hdr.color = rbRed
	hdr.parent = unsafe.Pointer(&root)
	hdr.left = unsafe.Pointer(&root)
	hdr.right = unsafe.Pointer(&root)
	root.color = 1
	root.parent = unsafe.Pointer(&hdr)

	got := RbTreeDecrement((*byte)(unsafe.Pointer(&hdr)))
	if got != (*byte)(unsafe.Pointer(&root)) {
		t.Fatalf("dec end = %p, want root", got)
	}
	got = RbTreeIncrement((*byte)(unsafe.Pointer(&root)))
	if got != (*byte)(unsafe.Pointer(&hdr)) {
		t.Fatalf("inc root = %p, want end", got)
	}
}

func TestRbTreeInsertEmpty(t *testing.T) {
	// Layout matches clang IR _Rb_tree_impl: [8 x i8] pad + header node + count.
	var impl struct {
		_     [rbTreeImplHeaderOff]byte
		hdr   rbNode
		count uint64
	}
	var a rbNode
	RbTreeInit((*byte)(unsafe.Pointer(&impl)))
	hdr := &impl.hdr
	RbTreeInsertAndRebalance(true, (*byte)(unsafe.Pointer(&a)), (*byte)(unsafe.Pointer(hdr)), (*byte)(unsafe.Pointer(hdr)))
	if hdr.parent != unsafe.Pointer(&a) || hdr.left != unsafe.Pointer(&a) || hdr.right != unsafe.Pointer(&a) {
		t.Fatalf("header %+v", *hdr)
	}
	if a.color != rbBlack || a.parent != unsafe.Pointer(hdr) {
		t.Fatalf("root %+v", a)
	}
	if p := RbTreeDecrement((*byte)(unsafe.Pointer(hdr))); p != (*byte)(unsafe.Pointer(&a)) {
		t.Fatalf("dec end")
	}
	if impl.count != 0 {
		// InsertAndRebalance does not bump count; ctor leaves 0.
	}
}

func TestRbTreeInorderTwo(t *testing.T) {
	//     hdr
	//      |
	//      b
	//     /
	//    a
	var hdr, a, b rbNode
	hdr.color = rbRed
	hdr.parent = unsafe.Pointer(&b)
	hdr.left = unsafe.Pointer(&a)
	hdr.right = unsafe.Pointer(&b)
	b.color = 1
	b.parent = unsafe.Pointer(&hdr)
	b.left = unsafe.Pointer(&a)
	a.color = rbRed
	a.parent = unsafe.Pointer(&b)

	if p := RbTreeIncrement((*byte)(unsafe.Pointer(&a))); p != (*byte)(unsafe.Pointer(&b)) {
		t.Fatalf("inc a")
	}
	if p := RbTreeDecrement((*byte)(unsafe.Pointer(&b))); p != (*byte)(unsafe.Pointer(&a)) {
		t.Fatalf("dec b")
	}
	if p := RbTreeIncrement((*byte)(unsafe.Pointer(&b))); p != (*byte)(unsafe.Pointer(&hdr)) {
		t.Fatalf("inc b")
	}
}