package libc

import "unsafe"

// libc++ __tree_node (LLVM 22, 8-byte pointers):
//
//	__left_     ptr  @0
//	__right_    ptr  @8
//	__parent_   ptr  @16
//	__is_black_ i8   @24
//	__value_    …    @32   pair<K,V> for map; CGContext::iv_bounds is
//	                       pair<const Variable*, unsigned> (16 bytes).
//
// Inlined insert can write the pair at +0, so __right_ becomes a small
// integer (Darwin csmith crash: __construct_from_tree → __get_value(0x3)).
const (
	libcxxTreeLeftOff   = 0
	libcxxTreeRightOff  = 8
	libcxxTreeParentOff = 16
	libcxxTreeBlackOff  = 24
	libcxxTreeValueOff  = 32
	libcxxTreeValueSize = 16
	libcxxTreeNodeSize  = 48
	libcxxTreeMinAddr   = 4096
)

func libcxxTreePtr(p *byte) bool {
	a := uintptr(unsafe.Pointer(p))
	return a >= libcxxTreeMinAddr && a%8 == 0
}

func libcxxTreeLoadPtr(n *byte, off int) *byte {
	return As[byte](Load[unsafe.Pointer](Ptr(n), off))
}

// libcxxTreeIntChild is a non-nil pointer that cannot be a node
// (unsigned map value written into __right_ / __left_).
func libcxxTreeIntChild(p *byte) bool {
	return p != nil && !libcxxTreePtr(p)
}

// libcxxTreeRealChild is a node whose __parent_ points at parent.
// A Variable* sitting in __left_ (pair overlaid at +0) fails this.
func libcxxTreeRealChild(parent, child *byte) bool {
	if !libcxxTreePtr(child) {
		return false
	}
	return libcxxTreeLoadPtr(child, libcxxTreeParentOff) == parent
}

// LibcxxTreeGetValue is __tree_node::__get_value. Returns this+32, or nil
// when this is not a user pointer (so a leftover unsigned in a child slot
// does not SEGV).
func LibcxxTreeGetValue(this *byte) *byte {
	if !libcxxTreePtr(this) {
		return nil
	}
	return As[byte](Off(Ptr(this), libcxxTreeValueOff))
}

// LibcxxTreeConstructFromTree is __tree::__construct_from_tree.
// tree and construct are unused (allocator / lambda); we malloc and
// memcpy the pair ourselves.
//
// Inlined insert may write pair<K,V> at +0, so __right_ is a small
// integer and __left_ is the key pointer. Do not walk those; copy the
// pair onto dest+32 and emit a leaf.
func LibcxxTreeConstructFromTree(tree, src, construct *byte) *byte {
	_ = tree
	_ = construct
	return libcxxTreeCopy(src, 0)
}

func libcxxTreeCopy(src *byte, depth int) *byte {
	if depth > 64 || !libcxxTreePtr(src) {
		return nil
	}
	n := Calloc[byte](1, int64(libcxxTreeNodeSize))
	if n == nil {
		return nil
	}
	Store(Ptr(n), libcxxTreeBlackOff, Load[byte](Ptr(src), libcxxTreeBlackOff))
	left := libcxxTreeLoadPtr(src, libcxxTreeLeftOff)
	right := libcxxTreeLoadPtr(src, libcxxTreeRightOff)
	overlay := libcxxTreeIntChild(left) || libcxxTreeIntChild(right) ||
		(left != nil && !libcxxTreeRealChild(src, left)) ||
		(right != nil && !libcxxTreeRealChild(src, right))
	valOff := libcxxTreeValueOff
	if overlay {
		valOff = 0
	}
	copy(
		Bytes(As[byte](Off(Ptr(n), libcxxTreeValueOff)), libcxxTreeValueSize),
		Bytes(As[byte](Off(Ptr(src), valOff)), libcxxTreeValueSize),
	)
	if overlay {
		return n
	}
	nl := libcxxTreeCopy(left, depth+1)
	nr := libcxxTreeCopy(right, depth+1)
	Store(Ptr(n), libcxxTreeLeftOff, Ptr(nl))
	Store(Ptr(n), libcxxTreeRightOff, Ptr(nr))
	if nl != nil {
		Store(Ptr(nl), libcxxTreeParentOff, Ptr(n))
	}
	if nr != nil {
		Store(Ptr(nr), libcxxTreeParentOff, Ptr(n))
	}
	return n
}

func libcxxTreeMin(x *byte) *byte {
	for {
		left := libcxxTreeLoadPtr(x, libcxxTreeLeftOff)
		if !libcxxTreeRealChild(x, left) {
			return x
		}
		x = left
	}
}

func libcxxTreeMax(x *byte) *byte {
	for {
		right := libcxxTreeLoadPtr(x, libcxxTreeRightOff)
		if !libcxxTreeRealChild(x, right) {
			return x
		}
		x = right
	}
}

func libcxxTreeIsLeftChild(x *byte) bool {
	p := libcxxTreeLoadPtr(x, libcxxTreeParentOff)
	if !libcxxTreePtr(p) {
		return false
	}
	return libcxxTreeLoadPtr(p, libcxxTreeLeftOff) == x
}

// LibcxxTreeNext is std::__tree_next. Bad child pointers (overlay
// unsigned / key) are treated as null so increment cannot hang.
func LibcxxTreeNext(x *byte) *byte {
	if !libcxxTreePtr(x) {
		return nil
	}
	right := libcxxTreeLoadPtr(x, libcxxTreeRightOff)
	if libcxxTreeRealChild(x, right) {
		return libcxxTreeMin(right)
	}
	for libcxxTreePtr(x) && !libcxxTreeIsLeftChild(x) {
		x = libcxxTreeLoadPtr(x, libcxxTreeParentOff)
	}
	if !libcxxTreePtr(x) {
		return nil
	}
	return libcxxTreeLoadPtr(x, libcxxTreeParentOff)
}

// LibcxxTreePrev is std::__tree_prev_iter. x may be the end node
// (only __left_ is meaningful).
func LibcxxTreePrev(x *byte) *byte {
	if !libcxxTreePtr(x) {
		return nil
	}
	left := libcxxTreeLoadPtr(x, libcxxTreeLeftOff)
	if libcxxTreeRealChild(x, left) || (libcxxTreePtr(left) && libcxxTreeLoadPtr(left, libcxxTreeParentOff) == x) {
		return libcxxTreeMax(left)
	}
	for libcxxTreePtr(x) && libcxxTreeIsLeftChild(x) {
		x = libcxxTreeLoadPtr(x, libcxxTreeParentOff)
	}
	if !libcxxTreePtr(x) {
		return nil
	}
	return libcxxTreeLoadPtr(x, libcxxTreeParentOff)
}

// LibcxxTreeMin is std::__tree_min.
func LibcxxTreeMin(x *byte) *byte {
	if !libcxxTreePtr(x) {
		return nil
	}
	return libcxxTreeMin(x)
}
