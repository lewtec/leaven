package leaven

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strings"
	"unsafe"

	"github.com/dave/jennifer/jen"
	"github.com/lewtec/leaven/internal/llir/ir/types"
	"github.com/lewtec/leaven/libc"
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

func (r goRef) Call(args ...jen.Code) *jen.Statement {
	return r.code().Call(args...)
}

func (r goRef) Types(t ...jen.Code) *jen.Statement {
	return r.code().Types(t...)
}

// Sym is a resolved package-level func (pkg+name). Call/Types emit jen.
// Generic instantiations drop [T]. Maps store this; FuncForPC runs once.
func Sym(fn any) goRef {
	pkg, name := funcPkgName(fn)
	return goRef{pkg: pkg, name: name}
}

func funcPkgName(fn any) (pkg, name string) {
	rv := reflect.ValueOf(fn)
	if rv.Kind() != reflect.Func {
		panic(fmt.Sprintf("Sym: not a func: %T", fn))
	}
	f := runtime.FuncForPC(rv.Pointer())
	if f == nil {
		panic(fmt.Sprintf("Sym: no PC: %T", fn))
	}
	full := f.Name()
	if i := strings.IndexByte(full, '['); i >= 0 {
		full = full[:i]
	}
	lastSlash := strings.LastIndexByte(full, '/')
	rest := full
	if lastSlash >= 0 {
		rest = full[lastSlash+1:]
	}
	dot := strings.IndexByte(rest, '.')
	if dot < 0 {
		panic(fmt.Sprintf("Sym: no pkg: %s", f.Name()))
	}
	return full[:len(full)-len(rest)+dot], rest[dot+1:]
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

// Qual is T as jennifer AST. Named types become jen.Qual(pkg, name).
func Qual[T any]() *jen.Statement {
	return reflectType(reflect.TypeFor[T]())
}

func reflectType(t reflect.Type) *jen.Statement {
	if t == nil {
		return jen.Any()
	}
	if t.Name() != "" && t.PkgPath() != "" {
		return jen.Qual(t.PkgPath(), t.Name())
	}
	switch t.Kind() {
	case reflect.Bool:
		return jen.Bool()
	case reflect.Int:
		return jen.Int()
	case reflect.Int8:
		return jen.Int8()
	case reflect.Int16:
		return jen.Int16()
	case reflect.Int32:
		return jen.Int32()
	case reflect.Int64:
		return jen.Int64()
	case reflect.Uint:
		return jen.Uint()
	case reflect.Uint8:
		return jen.Byte()
	case reflect.Uint16:
		return jen.Uint16()
	case reflect.Uint32:
		return jen.Uint32()
	case reflect.Uint64:
		return jen.Uint64()
	case reflect.Uintptr:
		return jen.Uintptr()
	case reflect.Float32:
		return jen.Float32()
	case reflect.Float64:
		return jen.Float64()
	case reflect.String:
		return jen.String()
	case reflect.UnsafePointer:
		return Qual[unsafe.Pointer]()
	case reflect.Ptr:
		return jen.Op("*").Add(reflectType(t.Elem()))
	case reflect.Slice:
		return jen.Index().Add(reflectType(t.Elem()))
	case reflect.Array:
		return jen.Index(jen.Lit(t.Len())).Add(reflectType(t.Elem()))
	case reflect.Map:
		return jen.Map(reflectType(t.Key())).Add(reflectType(t.Elem()))
	case reflect.Interface:
		return jen.Any()
	default:
		if t.Name() != "" {
			return jen.Id(t.Name())
		}
		return jen.Id(t.String())
	}
}

var (
	symPtr    = Sym(libc.Ptr[byte])
	symAs     = Sym(libc.As[byte])
	symAddr   = Sym(libc.Addr[byte])
	symOff    = Sym(libc.Off)
	symAddPtr = Sym(libc.AddPointer[byte])
)

// emitPtr is libc.Ptr(p) — (void *)p.
func emitPtr(p jen.Code) *jen.Statement {
	return symPtr.Call(p)
}

// emitAs is libc.As[T](p) — (T *)p. t is the element type, not *T.
func emitAs(t, p jen.Code) *jen.Statement {
	return symAs.Types(t).Call(p)
}

// emitAddr is libc.Addr(p) — (uintptr)p.
func emitAddr(p jen.Code) *jen.Statement {
	return symAddr.Call(p)
}

// emitOff is libc.Off(p, n) — (char *)p + n.
func emitOff(p, n jen.Code) *jen.Statement {
	return symOff.Call(p, n)
}

// emitAddPtr is libc.AddPointer[T](p, i).
func emitAddPtr(t, p, i jen.Code) *jen.Statement {
	return symAddPtr.Types(t).Call(p, i)
}

// emitUP is unsafe.Pointer(x) for integer bit patterns (inttoptr),
// not a Go *T. Use emitPtr when x is already a pointer.
func emitUP(x jen.Code) *jen.Statement {
	return Qual[unsafe.Pointer]().Call(x)
}

// ptrToUint is uintptr(unsafe.Pointer(p)). emitUP so untyped nil is valid.
func ptrToUint(p jen.Code) *jen.Statement {
	return jen.Uintptr().Call(emitUP(p))
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
		return Qual[libc.I128]()
	case 256:
		return Qual[libc.I256]()
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

func wideSym(bits uint64, i128, i256 any) goRef {
	if bits == 256 {
		return Sym(i256)
	}
	return Sym(i128)
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
	return Sym(libc.InitOstream).Call(emitPtr(addrOf(jen.Id(name))))
}

func i1PackFn(n uint64) (any, bool) {
	switch n {
	case 8:
		return libc.I1Pack8, true
	case 16:
		return libc.I1Pack16, true
	case 32:
		return libc.I1Pack32, true
	case 64:
		return libc.I1Pack64, true
	default:
		return nil, false
	}
}

func i1UnpackFn(n uint64) (any, bool) {
	switch n {
	case 8:
		return libc.I1Unpack8, true
	case 16:
		return libc.I1Unpack16, true
	case 32:
		return libc.I1Unpack32, true
	case 64:
		return libc.I1Unpack64, true
	default:
		return nil, false
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
		return conv(goIntType(it.BitSize), Sym(fn).Call(src)), nil
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
		return Sym(fn).Call(conv(goIntType(it.BitSize), src)), nil
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
		jen.Return(deref(emitAs(toT, emitPtr(addrOf(jen.Id("v")))))),
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
	bits, ok := scalarBitSize(from)
	if !ok {
		return nil, fmt.Errorf("%w: %v and %v", errIncompatiblePointers, from, to)
	}
	switch {
	case fromF != nil && toI != nil:
		switch bits {
		case 32:
			return conv(goIntType(32), Sym(math.Float32bits).Call(src)), nil
		case 64:
			return conv(goIntType(64), Sym(math.Float64bits).Call(src)), nil
		default:
			return nil, fmt.Errorf("%w: %v and %v", errUnsupportedFloatType, from, to)
		}
	case fromI != nil && toF != nil:
		switch bits {
		case 32:
			return Sym(math.Float32frombits).Call(jen.Uint32().Call(src)), nil
		case 64:
			return Sym(math.Float64frombits).Call(jen.Uint64().Call(src)), nil
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
	return expr{code: emitPtr(addrOf(base)), base: base}
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
