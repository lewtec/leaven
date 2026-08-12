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
	rbRed       = int32(0)
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
