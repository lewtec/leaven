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

// RbTreeInit is the empty _Rb_tree / map header: red, no root,
// left/right point at the header. libstdc++ _Rb_tree_header::_M_reset.
func RbTreeInit(header *byte) {
	if header == nil {
		return
	}
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
	y.parent = x.parent
	if unsafe.Pointer(x) == header.parent {
		header.parent = unsafe.Pointer(y)
	} else if unsafe.Pointer(x) == rbPtr(x.parent).left {
		rbPtr(x.parent).left = unsafe.Pointer(y)
	} else {
		rbPtr(x.parent).right = unsafe.Pointer(y)
	}
	y.left = unsafe.Pointer(x)
	x.parent = unsafe.Pointer(y)
}

func rbRotateRight(x, header *rbNode) {
	y := rbPtr(x.left)
	x.left = y.right
	if y.right != nil {
		rbPtr(y.right).parent = unsafe.Pointer(x)
	}
	y.parent = x.parent
	if unsafe.Pointer(x) == header.parent {
		header.parent = unsafe.Pointer(y)
	} else if unsafe.Pointer(x) == rbPtr(x.parent).right {
		rbPtr(x.parent).right = unsafe.Pointer(y)
	} else {
		rbPtr(x.parent).left = unsafe.Pointer(y)
	}
	y.right = unsafe.Pointer(x)
	x.parent = unsafe.Pointer(y)
}
