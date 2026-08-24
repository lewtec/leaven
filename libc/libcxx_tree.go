package libc

import "unsafe"

// libc++ __tree_node (LLVM 22, 8-byte pointers):
//
//	__left_     ptr  @0
//	__right_    ptr  @8
//	__parent_   ptr  @16
//	__is_black_ i8   @24
//	__value_    …    @32   pair<K,V> (16 bytes for Variable*→unsigned;
//	                       larger for Statement*→Effect).
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
	return libcxxTreeAddr(uintptr(unsafe.Pointer(p)))
}

func libcxxTreeAddr(a uintptr) bool {
	return a >= libcxxTreeMinAddr && a%8 == 0
}

func libcxxTreeLoadU(n *byte, off int) uintptr {
	return uintptr(Load[unsafe.Pointer](Ptr(n), off))
}

func libcxxTreeNode(a uintptr) *byte {
	if !libcxxTreeAddr(a) {
		return nil
	}
	return As[byte](unsafe.Pointer(a))
}

func libcxxTreeLoadPtr(n *byte, off int) *byte {
	return libcxxTreeNode(libcxxTreeLoadU(n, off))
}

// libcxxTreeRealChildU is a node whose __parent_ points at parent.
// A Variable* or leftover unsigned in __left_/__right_ fails this.
func libcxxTreeRealChildU(parent *byte, child uintptr) bool {
	c := libcxxTreeNode(child)
	if c == nil {
		return false
	}
	return libcxxTreeLoadU(c, libcxxTreeParentOff) == uintptr(unsafe.Pointer(parent))
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

func libcxxTreeOverlay(src *byte, leftU, rightU uintptr) bool {
	// Integer leftover (0x3, 0x7) first — do not deref it.
	if leftU != 0 && !libcxxTreeAddr(leftU) {
		return true
	}
	if rightU != 0 && !libcxxTreeAddr(rightU) {
		return true
	}
	if leftU != 0 && !libcxxTreeRealChildU(src, leftU) {
		return true
	}
	if rightU != 0 && !libcxxTreeRealChildU(src, rightU) {
		return true
	}
	return false
}

func libcxxTreeCopy(src *byte, depth int) *byte {
	if depth > 64 || !libcxxTreePtr(src) {
		return nil
	}
	leftU := libcxxTreeLoadU(src, libcxxTreeLeftOff)
	rightU := libcxxTreeLoadU(src, libcxxTreeRightOff)
	overlay := libcxxTreeOverlay(src, leftU, rightU)
	_, destN, valOff, valN := libcxxTreeCopySpan(src, overlay)
	n := Calloc[byte](1, int64(destN))
	if n == nil {
		return nil
	}
	Store(Ptr(n), libcxxTreeBlackOff, Load[byte](Ptr(src), libcxxTreeBlackOff))
	copy(
		Bytes(As[byte](Off(Ptr(n), libcxxTreeValueOff)), valN),
		Bytes(As[byte](Off(Ptr(src), valOff)), valN),
	)
	if overlay {
		return n
	}
	nl := libcxxTreeCopy(libcxxTreeNode(leftU), depth+1)
	nr := libcxxTreeCopy(libcxxTreeNode(rightU), depth+1)
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

// libcxxTreeCopySpan is how much of src to memcpy onto dest+32.
// Slab nodes use the real malloc size so map<Statement*, Effect>
// (Effect has a vector at +24) is not truncated to 16 bytes.
// Go-slice fixtures and other non-slab pointers keep the 16-byte pair.
func libcxxTreeCopySpan(src *byte, overlay bool) (srcN, destN, valOff, valN int) {
	srcN = slabUsable(src)
	valOff = libcxxTreeValueOff
	if overlay {
		valOff = 0
	}
	if srcN < libcxxTreeNodeSize {
		return srcN, libcxxTreeNodeSize, valOff, libcxxTreeValueSize
	}
	destN = srcN
	if overlay {
		destN = libcxxTreeValueOff + srcN
	}
	valN = destN - libcxxTreeValueOff
	if overlay {
		if valN > srcN {
			valN = srcN
		}
	} else if valN > srcN-libcxxTreeValueOff {
		valN = srcN - libcxxTreeValueOff
	}
	return srcN, destN, valOff, valN
}

func libcxxTreeMin(x *byte) *byte {
	for {
		left := libcxxTreeLoadU(x, libcxxTreeLeftOff)
		if !libcxxTreeRealChildU(x, left) {
			return x
		}
		x = libcxxTreeNode(left)
	}
}

func libcxxTreeMax(x *byte) *byte {
	for {
		right := libcxxTreeLoadU(x, libcxxTreeRightOff)
		if !libcxxTreeRealChildU(x, right) {
			return x
		}
		x = libcxxTreeNode(right)
	}
}

func libcxxTreeIsLeftChild(x *byte) bool {
	p := libcxxTreeLoadPtr(x, libcxxTreeParentOff)
	if p == nil {
		return false
	}
	return libcxxTreeLoadU(p, libcxxTreeLeftOff) == uintptr(unsafe.Pointer(x))
}

// LibcxxTreeNext is std::__tree_next. Bad child pointers (overlay
// unsigned / key) are treated as null so increment cannot hang.
func LibcxxTreeNext(x *byte) *byte {
	if !libcxxTreePtr(x) {
		return nil
	}
	right := libcxxTreeLoadU(x, libcxxTreeRightOff)
	if libcxxTreeRealChildU(x, right) {
		return libcxxTreeMin(libcxxTreeNode(right))
	}
	for i := 0; libcxxTreePtr(x) && !libcxxTreeIsLeftChild(x) && i < 64; i++ {
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
	left := libcxxTreeLoadU(x, libcxxTreeLeftOff)
	c := libcxxTreeNode(left)
	if libcxxTreeRealChildU(x, left) || (c != nil && libcxxTreeLoadU(c, libcxxTreeParentOff) == uintptr(unsafe.Pointer(x))) {
		return libcxxTreeMax(c)
	}
	for i := 0; libcxxTreePtr(x) && libcxxTreeIsLeftChild(x) && i < 64; i++ {
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
