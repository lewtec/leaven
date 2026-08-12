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

func goIntType(bits uint64) *jen.Statement {
	switch goIntBits(bits) {
	case 8:
		return jen.Byte()
	case 16:
		return jen.Int16()
	case 32:
		return jen.Int32()
	case 64:
		return jen.Int64()
	case 128:
		return libc("I128")
	default:
		return jen.Id(fmt.Sprintf("int%d", bits))
	}
}

func goUintType(bits uint64) *jen.Statement {
	switch goIntBits(bits) {
	case 8:
		return jen.Byte()
	case 16:
		return jen.Uint16()
	case 32:
		return jen.Uint32()
	case 64:
		return jen.Uint64()
	case 128:
		return libc("I128")
	default:
		return jen.Id(fmt.Sprintf("uint%d", bits))
	}
}

func isI128(t types.Type) bool {
	it, ok := t.(*types.IntType)
	return ok && it.BitSize == 128
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
