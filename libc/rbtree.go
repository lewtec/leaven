package libc

// libstdc++ _Rb_tree_node_base (x86_64):
//   _M_color  i32 @0  (_S_red=0, _S_black=1)
//   _M_parent ptr @8
//   _M_left   ptr @16
//   _M_right  ptr @24
// From gcc libstdc++-v3/src/c++98/tree.cc.
// Links are uint64 (LLVM ptr / x86_64 slot), not Go *rbNode: native
// pointers are 4 bytes on 386 and GEP *8 then lands in the wrong field.

const (
	rbColorOff  = 0
	rbParentOff = 8
	rbLeftOff   = 16
	rbRightOff  = 24
	rbCountOff  = 32
	rbRed       = int32(0)
	rbBlack     = int32(1)
)

type rbNode struct {
	color               int32
	_                   int32
	parent, left, right uint64
}

func rbPtr(u uint64) *rbNode {
	if u == 0 {
		return nil
	}
	return As[rbNode](PtrFromBits(u))
}

func rbBits(n *rbNode) uint64 {
	if n == nil {
		return 0
	}
	return PtrBits(Ptr(n))
}

func rb(p *byte) *rbNode {
	if p == nil {
		return nil
	}
	return As[rbNode](Ptr(p))
}

func rbByte(n *rbNode) *byte {
	if n == nil {
		return nil
	}
	return As[byte](Ptr(n))
}

// RbTreeDecrement is std::_Rb_tree_decrement.
func RbTreeDecrement(x *byte) *byte {
	n := rb(x)
	if n == nil {
		return nil
	}
	p := rbPtr(n.parent)
	if n.color == rbRed && p != nil && p.parent == rbBits(n) {
		return rbByte(rbPtr(n.right))
	}
	if n.left != 0 {
		y := rbPtr(n.left)
		for y.right != 0 {
			y = rbPtr(y.right)
		}
		return rbByte(y)
	}
	y := p
	for y != nil && rbBits(n) == y.left {
		n = y
		y = rbPtr(y.parent)
	}
	return rbByte(y)
}

// RbTreeIncrement is std::_Rb_tree_increment.
func RbTreeIncrement(x *byte) *byte {
	n := rb(x)
	if n == nil {
		return nil
	}
	if n.right != 0 {
		n = rbPtr(n.right)
		for n.left != 0 {
			n = rbPtr(n.left)
		}
		return rbByte(n)
	}
	y := rbPtr(n.parent)
	for y != nil && rbBits(n) == y.right {
		n = y
		y = rbPtr(y.parent)
	}
	if y != nil && n.right != rbBits(y) {
		n = y
	}
	return rbByte(n)
}

// rbTreeImplHeaderOff is the byte offset of _Rb_tree_header inside
// libstdc++ _Rb_tree_impl when the comparator is empty (EBO as [8 x i8]
// in clang IR). map/tree default ctors pass `this` (impl start), not the header.
const rbTreeImplHeaderOff = 8

// RbTreeInit is the empty map / _Rb_tree default ctor. `tree` points at the
// _Rb_tree / map object (impl start). Header is at +8; left/right self-ref
// the header node (libstdc++ _Rb_tree_header::_M_reset).
func RbTreeInit(tree *byte) {
	if tree == nil {
		return
	}
	header := As[byte](Off(Ptr(tree), rbTreeImplHeaderOff))
	n := rb(header)
	n.color = rbRed
	self := rbBits(n)
	n.parent = 0
	n.left = self
	n.right = self
	Store[uint64](Ptr(header), rbCountOff, 0)
}

// rbEmptyHeader reports a libstdc++ empty-tree header: parent nil, left
// and right self-ref. A bitwise copy of an empty map leaves dest.left
// pointing at the source header; begin() must not walk that as a value.
func rbEmptyHeader(n *rbNode) bool {
	if n == nil {
		return false
	}
	self := rbBits(n)
	return n.parent == 0 && n.left == self && n.right == self
}

// RbTreeBegin is std::map / _Rb_tree::begin. Returns the leftmost node,
// or the header when the tree is empty — including after a bitwise copy
// of an empty map (left still names the source header).
func RbTreeBegin(tree *byte) *byte {
	if tree == nil {
		return nil
	}
	header := As[byte](Off(Ptr(tree), rbTreeImplHeaderOff))
	h := rb(header)
	if h == nil || h.left == 0 || h.left == rbBits(h) {
		return header
	}
	if rbEmptyHeader(rbPtr(h.left)) {
		return header
	}
	return rbByte(rbPtr(h.left))
}

// RbTreeInsertAndRebalance is std::_Rb_tree_insert_and_rebalance
// (gcc libstdc++-v3/src/c++98/tree.cc).
func RbTreeInsertAndRebalance(insertLeft bool, x, p, header *byte) {
	if x == nil || p == nil || header == nil {
		return
	}
	// Bitwise-copied empty map: leftmost still names the source header.
	if hn := rb(header); hn != nil && hn.left != 0 && hn.left != rbBits(hn) && rbEmptyHeader(rbPtr(hn.left)) {
		self := rbBits(hn)
		hn.left = self
		hn.right = self
		hn.parent = 0
	}
	xn, pn, h := rb(x), rb(p), rb(header)
	xb, pb := rbBits(xn), rbBits(pn)
	xn.parent = pb
	xn.left = 0
	xn.right = 0
	xn.color = rbRed
	xp := xn
	if insertLeft {
		pn.left = xb
		if p == header {
			h.parent = xb
			h.right = xb
		} else if pb == h.left {
			h.left = xb
		}
	} else {
		pn.right = xb
		if pb == h.right {
			h.right = xb
		}
	}
	for rbBits(xp) != h.parent && xn.parent != 0 && rbPtr(xn.parent).color == rbRed {
		xpp := rbPtr(rbPtr(xn.parent).parent)
		if xpp == nil {
			break
		}
		if xn.parent == xpp.left {
			y := rbPtr(xpp.right)
			if y != nil && y.color == rbRed {
				rbPtr(xn.parent).color = rbBlack
				y.color = rbBlack
				xpp.color = rbRed
				xn = xpp
				xp = xn
				continue
			}
			if rbBits(xp) == rbPtr(xn.parent).right {
				xn = rbPtr(xn.parent)
				xp = xn
				rbRotateLeft(xn, h)
			}
			rbPtr(xn.parent).color = rbBlack
			xpp.color = rbRed
			rbRotateRight(xpp, h)
		} else {
			y := rbPtr(xpp.left)
			if y != nil && y.color == rbRed {
				rbPtr(xn.parent).color = rbBlack
				y.color = rbBlack
				xpp.color = rbRed
				xn = xpp
				xp = xn
				continue
			}
			if rbBits(xp) == rbPtr(xn.parent).left {
				xn = rbPtr(xn.parent)
				xp = xn
				rbRotateRight(xn, h)
			}
			rbPtr(xn.parent).color = rbBlack
			xpp.color = rbRed
			rbRotateLeft(xpp, h)
		}
	}
	if hp := rbPtr(h.parent); hp != nil {
		hp.color = rbBlack
	}
}

func rbRotateLeft(x, header *rbNode) {
	y := rbPtr(x.right)
	x.right = y.left
	if y.left != 0 {
		rbPtr(y.left).parent = rbBits(x)
	}
	rbAttachParent(x, y, header)
	y.left = rbBits(x)
	x.parent = rbBits(y)
}

func rbRotateRight(x, header *rbNode) {
	y := rbPtr(x.left)
	x.left = y.right
	if y.right != 0 {
		rbPtr(y.right).parent = rbBits(x)
	}
	rbAttachParent(x, y, header)
	y.right = rbBits(x)
	x.parent = rbBits(y)
}

// rbAttachParent puts y where x sat in the parent / header.
func rbAttachParent(x, y, header *rbNode) {
	y.parent = x.parent
	if rbBits(x) == header.parent {
		header.parent = rbBits(y)
		return
	}
	p := rbPtr(x.parent)
	if rbBits(x) == p.left {
		p.left = rbBits(y)
	} else {
		p.right = rbBits(y)
	}
}

// RbTreeRebalanceForErase is std::_Rb_tree_rebalance_for_erase
// (gcc libstdc++-v3/src/c++98/tree.cc). Returns the node to free (z).
// header is _M_header of the tree.
func RbTreeRebalanceForErase(z, header *byte) *byte {
	if z == nil || header == nil {
		return z
	}
	zn, h := rb(z), rb(header)
	y := zn
	var x, xParent *rbNode

	if zn.left == 0 {
		x = rbPtr(zn.right)
	} else if zn.right == 0 {
		x = rbPtr(zn.left)
	} else {
		y = rbPtr(zn.right)
		for y.left != 0 {
			y = rbPtr(y.left)
		}
		x = rbPtr(y.right)
	}

	yb, zb := rbBits(y), rbBits(zn)
	xb := rbBits(x)
	if y != zn {
		rbPtr(zn.left).parent = yb
		y.left = zn.left
		if yb != zn.right {
			xParent = rbPtr(y.parent)
			if x != nil {
				x.parent = y.parent
			}
			rbPtr(y.parent).left = xb
			y.right = zn.right
			rbPtr(zn.right).parent = yb
		} else {
			xParent = y
		}
		if h.parent == zb {
			h.parent = yb
		} else if zp := rbPtr(zn.parent); zp.left == zb {
			zp.left = yb
		} else {
			zp.right = yb
		}
		y.parent = zn.parent
		y.color, zn.color = zn.color, y.color
		y = zn
	} else {
		xParent = rbPtr(y.parent)
		if x != nil {
			x.parent = y.parent
		}
		if h.parent == zb {
			h.parent = xb
		} else if yp := rbPtr(y.parent); yp.left == zb {
			yp.left = xb
		} else {
			yp.right = xb
		}
		if h.left == zb {
			if zn.right == 0 {
				h.left = zn.parent
			} else {
				h.left = rbBits(rbMinimum(x))
			}
		}
		if h.right == zb {
			if zn.left == 0 {
				h.right = zn.parent
			} else {
				h.right = rbBits(rbMaximum(x))
			}
		}
	}

	if y.color != rbRed {
		for rbBits(x) != h.parent && (x == nil || x.color == rbBlack) {
			if rbBits(x) == xParent.left {
				w := rbPtr(xParent.right)
				if w != nil && w.color == rbRed {
					w.color = rbBlack
					xParent.color = rbRed
					rbRotateLeft(xParent, h)
					w = rbPtr(xParent.right)
				}
				if w == nil {
					x = xParent
					xParent = rbPtr(x.parent)
					continue
				}
				if (w.left == 0 || rbPtr(w.left).color == rbBlack) &&
					(w.right == 0 || rbPtr(w.right).color == rbBlack) {
					w.color = rbRed
					x = xParent
					xParent = rbPtr(x.parent)
				} else {
					if w.right == 0 || rbPtr(w.right).color == rbBlack {
						if w.left != 0 {
							rbPtr(w.left).color = rbBlack
						}
						w.color = rbRed
						rbRotateRight(w, h)
						w = rbPtr(xParent.right)
					}
					if w != nil {
						w.color = xParent.color
					}
					xParent.color = rbBlack
					if w != nil && w.right != 0 {
						rbPtr(w.right).color = rbBlack
					}
					rbRotateLeft(xParent, h)
					break
				}
			} else {
				w := rbPtr(xParent.left)
				if w != nil && w.color == rbRed {
					w.color = rbBlack
					xParent.color = rbRed
					rbRotateRight(xParent, h)
					w = rbPtr(xParent.left)
				}
				if w == nil {
					x = xParent
					xParent = rbPtr(x.parent)
					continue
				}
				if (w.right == 0 || rbPtr(w.right).color == rbBlack) &&
					(w.left == 0 || rbPtr(w.left).color == rbBlack) {
					w.color = rbRed
					x = xParent
					xParent = rbPtr(x.parent)
				} else {
					if w.left == 0 || rbPtr(w.left).color == rbBlack {
						if w.right != 0 {
							rbPtr(w.right).color = rbBlack
						}
						w.color = rbRed
						rbRotateLeft(w, h)
						w = rbPtr(xParent.left)
					}
					if w != nil {
						w.color = xParent.color
					}
					xParent.color = rbBlack
					if w != nil && w.left != 0 {
						rbPtr(w.left).color = rbBlack
					}
					rbRotateRight(xParent, h)
					break
				}
			}
		}
		if x != nil {
			x.color = rbBlack
		}
	}
	return rbByte(y)
}

func rbMinimum(n *rbNode) *rbNode {
	for n != nil && n.left != 0 {
		n = rbPtr(n.left)
	}
	return n
}

func rbMaximum(n *rbNode) *rbNode {
	for n != nil && n.right != 0 {
		n = rbPtr(n.right)
	}
	return n
}
