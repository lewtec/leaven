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
// memcpy the pair ourselves. Child pointers below 4096 are skipped.
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
	copy(
		Bytes(As[byte](Off(Ptr(n), libcxxTreeValueOff)), libcxxTreeValueSize),
		Bytes(As[byte](Off(Ptr(src), libcxxTreeValueOff)), libcxxTreeValueSize),
	)
	left := libcxxTreeCopy(libcxxTreeLoadPtr(src, libcxxTreeLeftOff), depth+1)
	right := libcxxTreeCopy(libcxxTreeLoadPtr(src, libcxxTreeRightOff), depth+1)
	Store(Ptr(n), libcxxTreeLeftOff, Ptr(left))
	Store(Ptr(n), libcxxTreeRightOff, Ptr(right))
	if left != nil {
		Store(Ptr(left), libcxxTreeParentOff, Ptr(n))
	}
	if right != nil {
		Store(Ptr(right), libcxxTreeParentOff, Ptr(n))
	}
	return n
}
