package libc

import "unsafe"

// libstdc++ _Rb_tree_node_base (x86_64):
//   _M_color  i32 @0  (_S_red=0, _S_black=1)
//   _M_parent ptr @8
//   _M_left   ptr @16
//   _M_right  ptr @24
// From gcc libstdc++-v3/src/c++98/tree.cc.

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
	parent unsafe.Pointer
	left   unsafe.Pointer
	right  unsafe.Pointer
}

func rb(p *byte) *rbNode {
	return (*rbNode)(unsafe.Pointer(p))
}

func rbPtr(p unsafe.Pointer) *rbNode {
	if p == nil {
		return nil
	}
	return (*rbNode)(p)
}

// RbTreeDecrement is std::_Rb_tree_decrement.
func RbTreeDecrement(x *byte) *byte {
	if x == nil {
		return nil
	}
	n := rb(x)
	if n.color == rbRed && rbPtr(n.parent) != nil && rbPtr(n.parent).parent == unsafe.Pointer(x) {
		return (*byte)(n.right)
	}
	if n.left != nil {
		y := rbPtr(n.left)
		for y.right != nil {
			y = rbPtr(y.right)
		}
		return (*byte)(unsafe.Pointer(y))
	}
	y := rbPtr(n.parent)
	for y != nil && unsafe.Pointer(n) == y.left {
		n = y
		y = rbPtr(y.parent)
	}
	if y == nil {
		return nil
	}
	return (*byte)(unsafe.Pointer(y))
}

// RbTreeIncrement is std::_Rb_tree_increment.
func RbTreeIncrement(x *byte) *byte {
	if x == nil {
		return nil
	}
	n := rb(x)
	if n.right != nil {
		n = rbPtr(n.right)
		for n.left != nil {
			n = rbPtr(n.left)
		}
		return (*byte)(unsafe.Pointer(n))
	}
	y := rbPtr(n.parent)
	for y != nil && unsafe.Pointer(n) == y.right {
		n = y
		y = rbPtr(y.parent)
	}
	if y != nil && n.right != unsafe.Pointer(y) {
		n = y
	}
	return (*byte)(unsafe.Pointer(n))
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
	header := (*byte)(unsafe.Add(unsafe.Pointer(tree), rbTreeImplHeaderOff))
	n := rb(header)
	n.color = rbRed
	n.parent = nil
	self := unsafe.Pointer(header)
	n.left = self
	n.right = self
	*(*uint64)(unsafe.Add(self, rbCountOff)) = 0
}

// RbTreeInsertAndRebalance is std::_Rb_tree_insert_and_rebalance
// (gcc libstdc++-v3/src/c++98/tree.cc).
func RbTreeInsertAndRebalance(insertLeft bool, x, p, header *byte) {
	if x == nil || p == nil || header == nil {
		return
	}
	xn, pn, h := rb(x), rb(p), rb(header)
	xn.parent = unsafe.Pointer(p)
	xn.left = nil
	xn.right = nil
	xn.color = rbRed
	xp := unsafe.Pointer(x)
	if insertLeft {
		pn.left = xp
		if p == header {
			h.parent = xp
			h.right = xp
		} else if unsafe.Pointer(p) == h.left {
			h.left = xp
		}
	} else {
		pn.right = xp
		if unsafe.Pointer(p) == h.right {
			h.right = xp
		}
	}
	for xp != h.parent && rbPtr(xn.parent) != nil && rbPtr(xn.parent).color == rbRed {
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
				xp = unsafe.Pointer(xn)
				continue
			}
			if xp == rbPtr(xn.parent).right {
				xn = rbPtr(xn.parent)
				xp = unsafe.Pointer(xn)
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
				xp = unsafe.Pointer(xn)
				continue
			}
			if xp == rbPtr(xn.parent).left {
				xn = rbPtr(xn.parent)
				xp = unsafe.Pointer(xn)
				rbRotateRight(xn, h)
			}
			rbPtr(xn.parent).color = rbBlack
			xpp.color = rbRed
			rbRotateLeft(xpp, h)
		}
	}
	if root := rbPtr(h.parent); root != nil {
		root.color = rbBlack
	}
}

func rbRotateLeft(x, header *rbNode) {
	y := rbPtr(x.right)
	x.right = y.left
	if y.left != nil {
		rbPtr(y.left).parent = unsafe.Pointer(x)
	}
	rbAttachParent(x, y, header)
	y.left = unsafe.Pointer(x)
	x.parent = unsafe.Pointer(y)
}

func rbRotateRight(x, header *rbNode) {
	y := rbPtr(x.left)
	x.left = y.right
	if y.right != nil {
		rbPtr(y.right).parent = unsafe.Pointer(x)
	}
	rbAttachParent(x, y, header)
	y.right = unsafe.Pointer(x)
	x.parent = unsafe.Pointer(y)
}

// rbAttachParent puts y where x sat in the parent / header.
func rbAttachParent(x, y, header *rbNode) {
	y.parent = x.parent
	if unsafe.Pointer(x) == header.parent {
		header.parent = unsafe.Pointer(y)
		return
	}
	p := rbPtr(x.parent)
	if unsafe.Pointer(x) == p.left {
		p.left = unsafe.Pointer(y)
	} else {
		p.right = unsafe.Pointer(y)
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
	var y *rbNode = zn
	var x unsafe.Pointer
	var xParent unsafe.Pointer

	if zn.left == nil {
		x = zn.right
	} else if zn.right == nil {
		x = zn.left
	} else {
		y = rbPtr(zn.right)
		for y.left != nil {
			y = rbPtr(y.left)
		}
		x = y.right
	}

	if y != zn {
		// y is successor; splice y in place of z
		rbPtr(zn.left).parent = unsafe.Pointer(y)
		y.left = zn.left
		if y != rbPtr(zn.right) {
			xParent = y.parent
			if x != nil {
				rbPtr(x).parent = y.parent
			}
			rbPtr(y.parent).left = x
			y.right = zn.right
			rbPtr(zn.right).parent = unsafe.Pointer(y)
		} else {
			xParent = unsafe.Pointer(y)
		}
		if h.parent == unsafe.Pointer(z) {
			h.parent = unsafe.Pointer(y)
		} else if rbPtr(zn.parent).left == unsafe.Pointer(z) {
			rbPtr(zn.parent).left = unsafe.Pointer(y)
		} else {
			rbPtr(zn.parent).right = unsafe.Pointer(y)
		}
		y.parent = zn.parent
		y.color, zn.color = zn.color, y.color
		y = zn
	} else {
		xParent = y.parent
		if x != nil {
			rbPtr(x).parent = y.parent
		}
		if h.parent == unsafe.Pointer(z) {
			h.parent = x
		} else if rbPtr(y.parent).left == unsafe.Pointer(z) {
			rbPtr(y.parent).left = x
		} else {
			rbPtr(y.parent).right = x
		}
		if h.left == unsafe.Pointer(z) {
			if zn.right == nil {
				h.left = zn.parent
			} else {
				h.left = unsafe.Pointer(rbMinimum(rbPtr(x)))
			}
		}
		if h.right == unsafe.Pointer(z) {
			if zn.left == nil {
				h.right = zn.parent
			} else {
				h.right = unsafe.Pointer(rbMaximum(rbPtr(x)))
			}
		}
	}

	if y.color != rbRed {
		for x != h.parent && (x == nil || rbPtr(x).color == rbBlack) {
			if x == rbPtr(xParent).left {
				w := rbPtr(rbPtr(xParent).right)
				if w != nil && w.color == rbRed {
					w.color = rbBlack
					rbPtr(xParent).color = rbRed
					rbRotateLeft(rbPtr(xParent), h)
					w = rbPtr(rbPtr(xParent).right)
				}
				if w == nil {
					x = xParent
					xParent = rbPtr(x).parent
					continue
				}
				if (w.left == nil || rbPtr(w.left).color == rbBlack) &&
					(w.right == nil || rbPtr(w.right).color == rbBlack) {
					w.color = rbRed
					x = xParent
					xParent = rbPtr(x).parent
				} else {
					if w.right == nil || rbPtr(w.right).color == rbBlack {
						if w.left != nil {
							rbPtr(w.left).color = rbBlack
						}
						w.color = rbRed
						rbRotateRight(w, h)
						w = rbPtr(rbPtr(xParent).right)
					}
					if w != nil {
						w.color = rbPtr(xParent).color
					}
					rbPtr(xParent).color = rbBlack
					if w != nil && w.right != nil {
						rbPtr(w.right).color = rbBlack
					}
					rbRotateLeft(rbPtr(xParent), h)
					break
				}
			} else {
				w := rbPtr(rbPtr(xParent).left)
				if w != nil && w.color == rbRed {
					w.color = rbBlack
					rbPtr(xParent).color = rbRed
					rbRotateRight(rbPtr(xParent), h)
					w = rbPtr(rbPtr(xParent).left)
				}
				if w == nil {
					x = xParent
					xParent = rbPtr(x).parent
					continue
				}
				if (w.right == nil || rbPtr(w.right).color == rbBlack) &&
					(w.left == nil || rbPtr(w.left).color == rbBlack) {
					w.color = rbRed
					x = xParent
					xParent = rbPtr(x).parent
				} else {
					if w.left == nil || rbPtr(w.left).color == rbBlack {
						if w.right != nil {
							rbPtr(w.right).color = rbBlack
						}
						w.color = rbRed
						rbRotateLeft(w, h)
						w = rbPtr(rbPtr(xParent).left)
					}
					if w != nil {
						w.color = rbPtr(xParent).color
					}
					rbPtr(xParent).color = rbBlack
					if w != nil && w.left != nil {
						rbPtr(w.left).color = rbBlack
					}
					rbRotateRight(rbPtr(xParent), h)
					break
				}
			}
		}
		if x != nil {
			rbPtr(x).color = rbBlack
		}
	}
	return (*byte)(unsafe.Pointer(y))
}

func rbMinimum(n *rbNode) *rbNode {
	for n != nil && n.left != nil {
		n = rbPtr(n.left)
	}
	return n
}

func rbMaximum(n *rbNode) *rbNode {
	for n != nil && n.right != nil {
		n = rbPtr(n.right)
	}
	return n
}
