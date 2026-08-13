package libc

// libstdc++ _Rb_tree_node_base (x86_64):
//   _M_color  i32 @0  (_S_red=0, _S_black=1)
//   _M_parent ptr @8
//   _M_left   ptr @16
//   _M_right  ptr @24
// From gcc libstdc++-v3/src/c++98/tree.cc.
// Links are *rbNode (same size as void*) so the overlay stays ABI-correct.

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
	color  int32
	_      int32
	parent *rbNode
	left   *rbNode
	right  *rbNode
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
	if n.color == rbRed && n.parent != nil && n.parent.parent == n {
		return rbByte(n.right)
	}
	if n.left != nil {
		y := n.left
		for y.right != nil {
			y = y.right
		}
		return rbByte(y)
	}
	y := n.parent
	for y != nil && n == y.left {
		n = y
		y = y.parent
	}
	return rbByte(y)
}

// RbTreeIncrement is std::_Rb_tree_increment.
func RbTreeIncrement(x *byte) *byte {
	n := rb(x)
	if n == nil {
		return nil
	}
	if n.right != nil {
		n = n.right
		for n.left != nil {
			n = n.left
		}
		return rbByte(n)
	}
	y := n.parent
	for y != nil && n == y.right {
		n = y
		y = y.parent
	}
	if y != nil && n.right != y {
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
	n.parent = nil
	n.left = n
	n.right = n
	Store[uint64](Ptr(header), rbCountOff, 0)
}

// RbTreeInsertAndRebalance is std::_Rb_tree_insert_and_rebalance
// (gcc libstdc++-v3/src/c++98/tree.cc).
func RbTreeInsertAndRebalance(insertLeft bool, x, p, header *byte) {
	if x == nil || p == nil || header == nil {
		return
	}
	xn, pn, h := rb(x), rb(p), rb(header)
	xn.parent = pn
	xn.left = nil
	xn.right = nil
	xn.color = rbRed
	xp := xn
	if insertLeft {
		pn.left = xp
		if p == header {
			h.parent = xp
			h.right = xp
		} else if pn == h.left {
			h.left = xp
		}
	} else {
		pn.right = xp
		if pn == h.right {
			h.right = xp
		}
	}
	for xp != h.parent && xn.parent != nil && xn.parent.color == rbRed {
		xpp := xn.parent.parent
		if xpp == nil {
			break
		}
		if xn.parent == xpp.left {
			y := xpp.right
			if y != nil && y.color == rbRed {
				xn.parent.color = rbBlack
				y.color = rbBlack
				xpp.color = rbRed
				xn = xpp
				xp = xn
				continue
			}
			if xp == xn.parent.right {
				xn = xn.parent
				xp = xn
				rbRotateLeft(xn, h)
			}
			xn.parent.color = rbBlack
			xpp.color = rbRed
			rbRotateRight(xpp, h)
		} else {
			y := xpp.left
			if y != nil && y.color == rbRed {
				xn.parent.color = rbBlack
				y.color = rbBlack
				xpp.color = rbRed
				xn = xpp
				xp = xn
				continue
			}
			if xp == xn.parent.left {
				xn = xn.parent
				xp = xn
				rbRotateRight(xn, h)
			}
			xn.parent.color = rbBlack
			xpp.color = rbRed
			rbRotateLeft(xpp, h)
		}
	}
	if h.parent != nil {
		h.parent.color = rbBlack
	}
}

func rbRotateLeft(x, header *rbNode) {
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	rbAttachParent(x, y, header)
	y.left = x
	x.parent = y
}

func rbRotateRight(x, header *rbNode) {
	y := x.left
	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}
	rbAttachParent(x, y, header)
	y.right = x
	x.parent = y
}

// rbAttachParent puts y where x sat in the parent / header.
func rbAttachParent(x, y, header *rbNode) {
	y.parent = x.parent
	if x == header.parent {
		header.parent = y
		return
	}
	p := x.parent
	if x == p.left {
		p.left = y
	} else {
		p.right = y
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

	if zn.left == nil {
		x = zn.right
	} else if zn.right == nil {
		x = zn.left
	} else {
		y = zn.right
		for y.left != nil {
			y = y.left
		}
		x = y.right
	}

	if y != zn {
		zn.left.parent = y
		y.left = zn.left
		if y != zn.right {
			xParent = y.parent
			if x != nil {
				x.parent = y.parent
			}
			y.parent.left = x
			y.right = zn.right
			zn.right.parent = y
		} else {
			xParent = y
		}
		if h.parent == zn {
			h.parent = y
		} else if zn.parent.left == zn {
			zn.parent.left = y
		} else {
			zn.parent.right = y
		}
		y.parent = zn.parent
		y.color, zn.color = zn.color, y.color
		y = zn
	} else {
		xParent = y.parent
		if x != nil {
			x.parent = y.parent
		}
		if h.parent == zn {
			h.parent = x
		} else if y.parent.left == zn {
			y.parent.left = x
		} else {
			y.parent.right = x
		}
		if h.left == zn {
			if zn.right == nil {
				h.left = zn.parent
			} else {
				h.left = rbMinimum(x)
			}
		}
		if h.right == zn {
			if zn.left == nil {
				h.right = zn.parent
			} else {
				h.right = rbMaximum(x)
			}
		}
	}

	if y.color != rbRed {
		for x != h.parent && (x == nil || x.color == rbBlack) {
			if x == xParent.left {
				w := xParent.right
				if w != nil && w.color == rbRed {
					w.color = rbBlack
					xParent.color = rbRed
					rbRotateLeft(xParent, h)
					w = xParent.right
				}
				if w == nil {
					x = xParent
					xParent = x.parent
					continue
				}
				if (w.left == nil || w.left.color == rbBlack) &&
					(w.right == nil || w.right.color == rbBlack) {
					w.color = rbRed
					x = xParent
					xParent = x.parent
				} else {
					if w.right == nil || w.right.color == rbBlack {
						if w.left != nil {
							w.left.color = rbBlack
						}
						w.color = rbRed
						rbRotateRight(w, h)
						w = xParent.right
					}
					if w != nil {
						w.color = xParent.color
					}
					xParent.color = rbBlack
					if w != nil && w.right != nil {
						w.right.color = rbBlack
					}
					rbRotateLeft(xParent, h)
					break
				}
			} else {
				w := xParent.left
				if w != nil && w.color == rbRed {
					w.color = rbBlack
					xParent.color = rbRed
					rbRotateRight(xParent, h)
					w = xParent.left
				}
				if w == nil {
					x = xParent
					xParent = x.parent
					continue
				}
				if (w.right == nil || w.right.color == rbBlack) &&
					(w.left == nil || w.left.color == rbBlack) {
					w.color = rbRed
					x = xParent
					xParent = x.parent
				} else {
					if w.left == nil || w.left.color == rbBlack {
						if w.right != nil {
							w.right.color = rbBlack
						}
						w.color = rbRed
						rbRotateLeft(w, h)
						w = xParent.left
					}
					if w != nil {
						w.color = xParent.color
					}
					xParent.color = rbBlack
					if w != nil && w.left != nil {
						w.left.color = rbBlack
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
	for n != nil && n.left != nil {
		n = n.left
	}
	return n
}

func rbMaximum(n *rbNode) *rbNode {
	for n != nil && n.right != nil {
		n = n.right
	}
	return n
}
