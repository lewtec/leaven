package libc

import (
	"testing"
	"unsafe"
)

func TestRbNodeLayout(t *testing.T) {
	if unsafe.Sizeof(uintptr(0)) != 8 {
		t.Skip("rbNode overlay matches x86_64 libstdc++")
	}
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
	hdr.parent = &root
	hdr.left = &root
	hdr.right = &root
	root.color = 1
	root.parent = &hdr

	got := RbTreeDecrement(rbByte(&hdr))
	if got != rbByte(&root) {
		t.Fatalf("dec end = %p, want root", got)
	}
	got = RbTreeIncrement(rbByte(&root))
	if got != rbByte(&hdr) {
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
	RbTreeInit(As[byte](Ptr(&impl)))
	hdr := &impl.hdr
	RbTreeInsertAndRebalance(true, rbByte(&a), rbByte(hdr), rbByte(hdr))
	if hdr.parent != &a || hdr.left != &a || hdr.right != &a {
		t.Fatalf("header %+v", *hdr)
	}
	if a.color != rbBlack || a.parent != hdr {
		t.Fatalf("root %+v", a)
	}
	if p := RbTreeDecrement(rbByte(hdr)); p != rbByte(&a) {
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
	hdr.parent = &b
	hdr.left = &a
	hdr.right = &b
	b.color = 1
	b.parent = &hdr
	b.left = &a
	a.color = rbRed
	a.parent = &b

	if p := RbTreeIncrement(rbByte(&a)); p != rbByte(&b) {
		t.Fatalf("inc a")
	}
	if p := RbTreeDecrement(rbByte(&b)); p != rbByte(&a) {
		t.Fatalf("dec b")
	}
	if p := RbTreeIncrement(rbByte(&b)); p != rbByte(&hdr) {
		t.Fatalf("inc b")
	}
}

func TestRbTreeBeginBitwiseEmptyCopy(t *testing.T) {
	var src, dst struct {
		_   [rbTreeImplHeaderOff]byte
		hdr rbNode
	}
	RbTreeInit(As[byte](Ptr(&src)))
	dst = src
	got := RbTreeBegin(As[byte](Ptr(&dst)))
	want := rbByte(&dst.hdr)
	if got != want {
		t.Fatalf("begin after bitwise empty copy = %p, want dest header %p (src header %p)",
			got, want, rbByte(&src.hdr))
	}
}
