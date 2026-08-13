package leaven

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/lewtec/leaven/internal/llir/ir/types"
)

const libcPath = "github.com/lewtec/leaven/libc"

// goRef is a Go identifier, optionally from another package.
// Empty pkg means a local or builtin name (jen.Id).
type goRef struct {
	pkg  string
	name string
}

func (r goRef) code() *jen.Statement {
	if r.pkg == "" {
		return jen.Id(r.name)
	}
	return jen.Qual(r.pkg, r.name)
}

func libc(name string) *jen.Statement {
	return jen.Qual(libcPath, name)
}

func newGoFile(packageName string) *jen.File {
	f := jen.NewFile(packageName)
	f.ImportName(libcPath, "libc")
	f.ImportName("sync/atomic", "atomic")
	f.ImportName("math", "math")
	f.ImportName("os", "os")
	f.ImportName("unsafe", "unsafe")
	return f
}

func assign(name string, rhs jen.Code) *jen.Statement {
	return jen.Id(name).Op("=").Add(rhs)
}

func bin(lhs jen.Code, op string, rhs jen.Code) *jen.Statement {
	return jen.Add(lhs).Op(op).Add(rhs)
}

func conv(typ, expr jen.Code) *jen.Statement {
	return jen.Add(typ).Call(expr)
}

func ptrTyp(t jen.Code) *jen.Statement {
	return jen.Op("*").Add(t)
}

func addrOf(x jen.Code) *jen.Statement {
	return jen.Op("&").Add(x)
}

func deref(x jen.Code) *jen.Statement {
	return jen.Op("*").Add(x)
}

func unsafePtr(x jen.Code) *jen.Statement {
	return jen.Qual("unsafe", "Pointer").Call(x)
}

func uintptrOfPtr(x jen.Code) *jen.Statement {
	return jen.Uintptr().Call(unsafePtr(x))
}

func ptrCast(typ, expr jen.Code) *jen.Statement {
	return jen.Parens(typ).Call(unsafePtr(expr))
}

func goIntType(bits uint64) *jen.Statement  { return goBitsType(bits, false) }
func goUintType(bits uint64) *jen.Statement { return goBitsType(bits, true) }

func goBitsType(bits uint64, unsigned bool) *jen.Statement {
	switch goIntBits(bits) {
	case 8:
		return jen.Byte()
	case 16:
		if unsigned {
			return jen.Uint16()
		}
		return jen.Int16()
	case 32:
		if unsigned {
			return jen.Uint32()
		}
		return jen.Int32()
	case 64:
		if unsigned {
			return jen.Uint64()
		}
		return jen.Int64()
	case 128:
		return libc("I128")
	case 256:
		return libc("I256")
	default:
		kind := "int"
		if unsigned {
			kind = "uint"
		}
		return jen.Id(fmt.Sprintf("%s%d", kind, bits))
	}
}

func isI128(t types.Type) bool { return isIntBits(t, 128) }
func isI256(t types.Type) bool { return isIntBits(t, 256) }

func isIntBits(t types.Type, n uint64) bool {
	it, ok := t.(*types.IntType)
	return ok && it.BitSize == n
}

// wideBits is 128 or 256 when t is a limb integer we lower via libc.
func wideBits(t types.Type) (uint64, bool) {
	it, ok := t.(*types.IntType)
	if !ok {
		return 0, false
	}
	switch it.BitSize {
	case 128, 256:
		return it.BitSize, true
	default:
		return 0, false
	}
}

func wideFn(bits uint64, op string) string {
	if bits == 256 {
		return "I256" + op
	}
	return "I128" + op
}

func litUntyped(n int64) *jen.Statement {
	return jen.Lit(int(n))
}

func litUint64(n uint64) *jen.Statement {
	return jen.Uint64().Call(jen.Op(fmt.Sprintf("%d", n)))
}

type vecBin struct {
	dest, op string
	x, y     jen.Code
}

func vectorBin(v vecBin) *jen.Statement {
	return jen.For(jen.List(jen.Id("i"), jen.Id("v")).Op(":=").Range().Add(v.x)).Block(
		jen.Id(v.dest).Index(jen.Id("i")).Op("=").Id("v").Op(v.op).Add(v.y).Index(jen.Id("i")),
	)
}

func isStdStream(name string) bool {
	switch name {
	case "_ZSt4cout", "_ZSt4cerr", "_ZSt4clog", "_ZSt3cin":
		return true
	default:
		return false
	}
}

func initStdStream(name string) *jen.Statement {
	return libc("InitOstream").Call(unsafePtr(addrOf(jen.Id(name))))
}

func i1PackFn(n uint64) (string, bool) {
	switch n {
	case 8:
		return "I1Pack8", true
	case 16:
		return "I1Pack16", true
	case 32:
		return "I1Pack32", true
	case 64:
		return "I1Pack64", true
	default:
		return "", false
	}
}

func i1UnpackFn(n uint64) (string, bool) {
	switch n {
	case 8:
		return "I1Unpack8", true
	case 16:
		return "I1Unpack16", true
	case 32:
		return "I1Unpack32", true
	case 64:
		return "I1Unpack64", true
	default:
		return "", false
	}
}

// i1VectorBitCast packs <N x i1> to iN or unpacks iN to <N x i1>.
// LLVM stores i1 vectors packed; Go keeps [N]bool so extract/and stay lanes.
func i1VectorBitCast(src *jen.Statement, from, to types.Type) (*jen.Statement, error) {
	if n, ok := i1VectorLen(from); ok {
		it, ok := to.(*types.IntType)
		if !ok || it.BitSize != n {
			return nil, fmt.Errorf("%w: %v and %v", errIncompatiblePointers, from, to)
		}
		fn, ok := i1PackFn(n)
		if !ok {
			return nil, fmt.Errorf("%w: <%d x i1>", errUnsupportedIntWidth, n)
		}
		return conv(goIntType(it.BitSize), libc(fn).Call(src)), nil
	}
	if n, ok := i1VectorLen(to); ok {
		it, ok := from.(*types.IntType)
		if !ok || it.BitSize != n {
			return nil, fmt.Errorf("%w: %v and %v", errIncompatiblePointers, from, to)
		}
		fn, ok := i1UnpackFn(n)
		if !ok {
			return nil, fmt.Errorf("%w: <%d x i1>", errUnsupportedIntWidth, n)
		}
		return libc(fn).Call(conv(goIntType(it.BitSize), src)), nil
	}
	return nil, nil
}

// vectorBitCast reinterprets same-size vectors (LLVM bitcast) and
// same-size vector↔scalar (e.g. rustc `bitcast <2 x i64> %x to i128`).
// Not a pointer cast. The IIFE makes a non-addressable operand addressable.
func vectorBitCast(src *jen.Statement, from, to types.Type) (*jen.Statement, error) {
	if !sameSizeVectors(from, to) && !sameSizeVectorScalar(from, to) {
		return nil, nil
	}
	fromT, err := TypeSpec(from)
	if err != nil {
		return nil, err
	}
	toT, err := TypeSpec(to)
	if err != nil {
		return nil, err
	}
	return jen.Func().Params(jen.Id("v").Add(fromT)).Add(toT).Block(
		jen.Return(deref(jen.Parens(ptrTyp(toT)).Call(unsafePtr(addrOf(jen.Id("v")))))),
	).Call(src), nil
}

// scalarBitCast reinterprets same-size int/float bits (LLVM bitcast).
// rustc f32::to_bits is `bitcast float %x to i32`, not a pointer cast.
func scalarBitCast(src *jen.Statement, from, to types.Type) (*jen.Statement, error) {
	if !sameSizeScalars(from, to) {
		if _, ok := scalarBitSize(from); ok {
			if _, ok := scalarBitSize(to); ok {
				return nil, fmt.Errorf("%w: %v and %v", errIncompatiblePointers, from, to)
			}
		}
		return nil, nil
	}
	fromF, _ := from.(*types.FloatType)
	toF, _ := to.(*types.FloatType)
	fromI, _ := from.(*types.IntType)
	toI, _ := to.(*types.IntType)
	bits, _ := scalarBitSize(from)
	switch {
	case fromF != nil && toI != nil:
		switch bits {
		case 32:
			return conv(goIntType(32), jen.Qual("math", "Float32bits").Call(src)), nil
		case 64:
			return conv(goIntType(64), jen.Qual("math", "Float64bits").Call(src)), nil
		default:
			return nil, fmt.Errorf("%w: %v and %v", errUnsupportedFloatType, from, to)
		}
	case fromI != nil && toF != nil:
		switch bits {
		case 32:
			return jen.Qual("math", "Float32frombits").Call(jen.Uint32().Call(src)), nil
		case 64:
			return jen.Qual("math", "Float64frombits").Call(jen.Uint64().Call(src)), nil
		default:
			return nil, fmt.Errorf("%w: %v and %v", errUnsupportedFloatType, from, to)
		}
	default:
		dst, err := TypeSpec(to)
		if err != nil {
			return nil, err
		}
		return conv(dst, src), nil
	}
}

func one(s *jen.Statement) []jen.Code {
	if s == nil {
		return nil
	}
	return []jen.Code{s}
}

func fieldName(i int) string {
	return fmt.Sprintf("F%d", i)
}

func fieldNameU(i uint64) string {
	return fmt.Sprintf("F%d", i)
}

// expr is a Go expression. If base != nil, code is &base (so load/store can
// drop the address-of instead of emitting *&x).
type expr struct {
	code *jen.Statement
	base *jen.Statement
}

func val(c *jen.Statement) expr {
	return expr{code: c}
}

func ident(name string) expr {
	return expr{code: jen.Id(name)}
}

func addrExpr(base *jen.Statement) expr {
	return expr{code: unsafePtr(addrOf(base)), base: base}
}

func (e expr) load() *jen.Statement {
	if e.base != nil {
		return e.base
	}
	return deref(e.code)
}

func (e expr) store(src jen.Code) *jen.Statement {
	if e.base != nil {
		return jen.Add(e.base).Op("=").Add(src)
	}
	return deref(e.code).Op("=").Add(src)
}

func (e expr) dropAddr() *jen.Statement {
	if e.base != nil {
		return e.base
	}
	return e.code
}
