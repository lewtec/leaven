package leaven

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/lewtec/leaven/internal/llir/ir"
	"github.com/lewtec/leaven/internal/llir/ir/constant"
	"github.com/lewtec/leaven/internal/llir/ir/enum"
	"github.com/lewtec/leaven/internal/llir/ir/types"
	"github.com/lewtec/leaven/internal/llir/ir/value"
)

// moduleFuncNames / moduleTypeNames are filled by collectModuleNames so locals
// can be renamed when they would shadow a function or type in Go.
var (
	moduleFuncNames map[string]bool
	moduleTypeNames map[string]bool
)

// collectModuleNames records function and type identifiers used in the module.
func collectModuleNames(m *ir.Module) {
	moduleFuncNames = make(map[string]bool)
	moduleTypeNames = make(map[string]bool)
	for _, f := range m.Funcs {
		moduleFuncNames[rawIdentName(f)] = true
	}
	for _, t := range m.TypeDefs {
		if n := TypeName(t); n != "" {
			moduleTypeNames[n] = true
		}
	}
}

// rawIdentName sanitizes an LLVM name to a Go identifier without clash renames.
func rawIdentName(v value.Named) string {
	name := v.Name()
	if name == "" {
		return "v" + strings.TrimPrefix(v.Ident(), "%")
	}
	if c := name[0]; '0' <= c && c <= '9' {
		name = "v" + name
	}
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "-", "_")
	if invalidNames[name] {
		name = "_" + name
	}
	return name
}

// VariableName returns the name to use for a local variable or parameter.
func VariableName(v value.Named) string {
	name := rawIdentName(v)
	// Params named like SSA temps (v0, v1, …) collide with anonymous
	// instructions ("%1" → "v1"). Prefix params so both can coexist.
	if _, ok := v.(*ir.Param); ok && ssaTempName(name) {
		name = "arg_" + name
	}
	// Locals/params must not reuse function or type names (Go forbids
	// `var state *state` and `call = is_verbatim(x)` when is_verbatim is *int32).
	if isLocalNamed(v) && (moduleFuncNames[name] || moduleTypeNames[name]) {
		name = "local_" + name
	}
	return name
}

// isLocalNamed reports whether v is a function parameter or instruction result
// (not a global/function/block label).
func isLocalNamed(v value.Named) bool {
	switch v.(type) {
	case *ir.Func, *ir.Global, *ir.Block:
		return false
	default:
		return true
	}
}

// ssaTempName reports whether name looks like an anonymous SSA temp (v0, v12).
func ssaTempName(name string) bool {
	if len(name) < 2 || name[0] != 'v' {
		return false
	}
	for i := 1; i < len(name); i++ {
		if name[i] < '0' || name[i] > '9' {
			return false
		}
	}
	return true
}

func BlockName(v value.Value) string {
	block := v.(*ir.Block)
	// Don't use Name(); numeric labels come back quoted (`"2"`) for LLVM print.
	name := block.LocalName
	if name == "" {
		return "block" + strconv.FormatInt(block.LocalID, 10)
	}
	if _, err := strconv.ParseInt(name, 10, 64); err == nil {
		return "block" + name
	}
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "-", "_")
	if invalidNames[name] {
		name = "_" + name
	}
	return name
}

// invalidNames are Go keywords and predeclared ids that cannot be used as names.
// https://go.dev/ref/spec#Keywords
var invalidNames = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
	// Common predeclared / special
	"init": true, "true": true, "false": true, "iota": true, "nil": true,
	"bool": true, "byte": true, "error": true, "int": true, "int8": true,
	"int16": true, "int32": true, "int64": true, "rune": true, "string": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"uintptr": true, "float32": true, "float64": true, "complex64": true,
	"complex128": true, "any": true, "comparable": true,
}

func namedRef(name string) (*jen.Statement, bool) {
	if ref, ok := libraryFunctions[name]; ok {
		return ref.code(), true
	}
	if c := rustRuntime(name); c != nil {
		return c, true
	}
	return nil, false
}

// hasRuntimeDef reports whether name is provided by libc, a rust shim, or an
// llvm.* intrinsic that translateCall handles. Declare-only IR symbols with
// no runtime def become panic/zero stubs.
func hasRuntimeDef(name string) bool {
	if _, ok := libraryFunctions[name]; ok {
		return true
	}
	if _, ok := libraryGlobals[name]; ok {
		return true
	}
	if _, ok := namedRef(name); ok {
		return true
	}
	return strings.HasPrefix(name, "llvm_")
}

func compositeValues(elems []jen.Code) *jen.Statement {
	return jen.Values(elems...)
}

func formatComposite(typ *jen.Statement, elems []jen.Code) *jen.Statement {
	return jen.Add(typ).Add(compositeValues(elems))
}

func fnPtrBitcast(to, from *jen.Statement) *jen.Statement {
	// Go forbids unsafe.Pointer(funcValue). Reinterpret via address of a temp.
	return jen.Func().Params().Add(to).Block(
		jen.Id("tmp").Op(":=").Add(from),
		jen.Return(deref(jen.Parens(ptrTyp(to))).Call(unsafePtr(addrOf(jen.Id("tmp"))))),
	).Call()
}

// FormatValue formats a constant or variable as it should appear in an expression.
func FormatValue(v value.Value) (*jen.Statement, error) {
	e, err := formatExpr(v)
	if err != nil {
		return nil, err
	}
	return e.code, nil
}

func formatExpr(v value.Value) (expr, error) {
	switch v := v.(type) {
	case *ir.Func:
		name := VariableName(v)
		fn := jen.Id(name)
		if c, ok := namedRef(name); ok {
			fn = c
		}
		return val(fnPtrBitcast(jen.Qual("unsafe", "Pointer"), fn)), nil

	case *ir.Alias:
		if v.Aliasee == nil {
			return expr{}, fmt.Errorf("alias %s has no aliasee", v.Name())
		}
		return formatExpr(v.Aliasee)

	case *ir.Global:
		name := VariableName(v)
		if types.IsFunc(v.ContentType) {
			fn := jen.Id(name)
			if c, ok := namedRef(name); ok {
				fn = c
			}
			return val(fnPtrBitcast(jen.Qual("unsafe", "Pointer"), fn)), nil
		}
		if ref, ok := libraryGlobals[name]; ok {
			return addrExpr(ref.code()), nil
		}
		return addrExpr(jen.Id(name)), nil

	case value.Named:
		name := VariableName(v)
		if c, ok := namedRef(name); ok {
			return val(c), nil
		}
		return ident(name), nil

	case *ir.Arg:
		return formatExpr(v.Value)

	case *constant.Array:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return expr{}, fmt.Errorf("error translating type (%v): %w", v.Typ, err)
		}
		elems := make([]jen.Code, len(v.Elems))
		for i, c := range v.Elems {
			e, err := FormatValue(c)
			if err != nil {
				return expr{}, fmt.Errorf("error translating element %d (%v): %w", i, c, err)
			}
			elems[i] = e
		}
		return val(formatComposite(t, elems)), nil

	case *constant.CharArray:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return expr{}, fmt.Errorf("error translating type (%v): %w", v.Typ, err)
		}
		elems := make([]jen.Code, len(v.X))
		for i, c := range v.X {
			elems[i] = litUntyped(int64(c))
		}
		return val(formatComposite(t, elems)), nil

	case *constant.ExprBitCast:
		from, err := FormatValue(v.From)
		if err != nil {
			return expr{}, fmt.Errorf("error translating source (%v): %w", v.From, err)
		}
		if isTaggedPointerType(v.To) {
			return val(uintptrOfPtr(from)), nil
		}
		if isTaggedPointerType(v.From.Type()) {
			to, err := TypeSpec(v.To)
			if err != nil {
				return expr{}, fmt.Errorf("error translating type (%v): %w", v.To, err)
			}
			return val(ptrCast(to, from)), nil
		}
		// Pointer values are already unsafe.Pointer.
		return val(from), nil

	case *constant.ExprIntToPtr:
		from, err := FormatValue(v.From)
		if err != nil {
			return expr{}, fmt.Errorf("error translating source (%v): %w", v.From, err)
		}
		to, err := TypeSpec(v.To)
		if err != nil {
			return expr{}, fmt.Errorf("error translating type (%v): %w", v.To, err)
		}
		return val(jen.Parens(to).Call(unsafePtr(jen.Uintptr().Call(from)))), nil

	case *constant.ExprPtrToInt:
		from, err := FormatValue(v.From)
		if err != nil {
			return expr{}, fmt.Errorf("error translating source (%v): %w", v.From, err)
		}
		to, err := TypeSpec(v.To)
		if err != nil {
			return expr{}, fmt.Errorf("error translating type (%v): %w", v.To, err)
		}
		return val(conv(to, uintptrOfPtr(from))), nil

	case *constant.ExprGetElementPtr:
		indices := make([]value.Value, len(v.Indices))
		for i, index := range v.Indices {
			indices[i] = index
		}
		return GetElementPtr(v.ElemType, v.Src, indices)

	case *constant.ExprICmp:
		c, err := formatICmp(v.Pred, v.X, v.Y)
		if err != nil {
			return expr{}, err
		}
		return val(c), nil

	case *constant.ExprZExt:
		c, err := formatZExt(v.From, v.To)
		if err != nil {
			return expr{}, err
		}
		return val(c), nil

	case *constant.ExprSExt:
		c, err := formatSExt(v.From, v.To)
		if err != nil {
			return expr{}, err
		}
		return val(c), nil

	case *constant.ExprTrunc:
		c, err := formatTrunc(v.From, v.To)
		if err != nil {
			return expr{}, err
		}
		return val(c), nil

	case *constant.ExprAdd:
		return formatBinConst("+", v.X, v.Y)
	case *constant.ExprSub:
		return formatBinConst("-", v.X, v.Y)
	case *constant.ExprMul:
		return formatBinConst("*", v.X, v.Y)
	case *constant.ExprAnd:
		return formatBinConst("&", v.X, v.Y)
	case *constant.ExprOr:
		return formatBinConst("|", v.X, v.Y)
	case *constant.ExprXor:
		return formatBinConst("^", v.X, v.Y)
	case *constant.ExprShl:
		return formatBinConst("<<", v.X, v.Y)
	case *constant.ExprLShr:
		return formatBinConst(">>", v.X, v.Y)
	case *constant.ExprAShr:
		return formatBinConst(">>", v.X, v.Y)

	case *constant.Float:
		result := v.X.String()
		var c *jen.Statement
		switch result {
		case "+Inf":
			c = jen.Qual("math", "Inf").Call(jen.Lit(1))
		case "-Inf":
			c = jen.Qual("math", "Inf").Call(jen.Lit(-1))
		case "NaN":
			c = jen.Qual("math", "NaN").Call()
		default:
			c = jen.Op(result)
		}
		if c != nil && (result == "+Inf" || result == "-Inf" || result == "NaN") && v.Typ.Kind == types.FloatKindFloat {
			c = jen.Float32().Call(c)
		}
		return val(c), nil

	case *constant.Index:
		return formatExpr(v.Constant)

	case *constant.Int:
		if v.Typ.BitSize > 64 {
			return expr{}, fmt.Errorf("%w: i%d constant %v", errUnsupportedIntWidth, v.Typ.BitSize, v.X)
		}
		var n int64
		switch {
		case v.X.IsInt64():
			n = v.X.Int64()
		case v.X.IsUint64():
			// Reinterpret as two's-complement signed for this bit width
			// (e.g. i64 all-ones → -1, not 18446744073709551615).
			n = int64(v.X.Uint64())
		default:
			// Truncate wide math/big values into the declared bit width (≤64).
			truncated, err := intFromBig(v.X, v.Typ.BitSize)
			if err != nil {
				return expr{}, err
			}
			n = truncated
		}

		if v.Typ.BitSize == 1 {
			if n != 0 {
				return val(jen.True()), nil
			}
			return val(jen.False()), nil
		}
		switch goIntBits(v.Typ.BitSize) {
		case 8:
			return val(litUntyped(int64(byte(n)))), nil
		case 16, 32:
			return val(litUntyped(n)), nil
		case 64:
			// Typed so large bit patterns never appear as untyped ints.
			return val(jen.Lit(n)), nil
		default:
			return val(litUntyped(n)), nil
		}

	case *constant.Null:
		// Tagged union pointers are uintptr; use 0 not nil.
		if isTaggedPointerType(v.Typ) {
			return val(jen.Lit(0)), nil
		}
		return val(jen.Nil()), nil

	case *constant.Struct:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return expr{}, fmt.Errorf("error translating type (%v): %w", v.Typ, err)
		}
		elems := make([]jen.Code, len(v.Fields))
		for i, c := range v.Fields {
			e, err := FormatValue(c)
			if err != nil {
				return expr{}, fmt.Errorf("error translating field %d (%v): %w", i, c, err)
			}
			elems[i] = e
		}
		return val(formatComposite(t, elems)), nil

	case *constant.Undef:
		return zeroOf(v.Typ)
	case *constant.Poison:
		return zeroOf(v.Typ)

	case *constant.Vector:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return expr{}, fmt.Errorf("error translating type (%v): %w", v.Typ, err)
		}
		elems := make([]jen.Code, len(v.Elems))
		for i, c := range v.Elems {
			e, err := FormatValue(c)
			if err != nil {
				return expr{}, fmt.Errorf("error translating element %d (%v): %w", i, c, err)
			}
			elems[i] = e
		}
		return val(formatComposite(t, elems)), nil

	case *constant.ZeroInitializer:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return expr{}, fmt.Errorf("error translating type (%v): %w", v.Typ, err)
		}
		return val(jen.Add(t).Values()), nil

	default:
		return expr{}, fmt.Errorf("%w: %T", errUnsupportedValueType, v)
	}
}

func zeroOf(typ types.Type) (expr, error) {
	switch typ.(type) {
	case *types.ArrayType, *types.StructType, *types.VectorType:
		t, err := TypeSpec(typ)
		if err != nil {
			return expr{}, fmt.Errorf("error translating type (%v): %w", typ, err)
		}
		return val(jen.Add(t).Values()), nil
	case *types.IntType, *types.FloatType:
		return val(jen.Lit(0)), nil
	case *types.PointerType:
		return val(jen.Nil()), nil
	default:
		return expr{}, fmt.Errorf("%w: %v", errUnsupportedUndefType, typ)
	}
}

var libraryGlobals = map[string]goRef{
	"stdin":  {pkg: "os", name: "Stdin"},
	"stdout": {pkg: "os", name: "Stdout"},
	"stderr": {pkg: "os", name: "Stderr"},
}

// intFromBig truncates x to bitSize bits and returns it as int64 (two's complement).
func intFromBig(x *big.Int, bitSize uint64) (int64, error) {
	if bitSize == 0 || bitSize > 64 {
		return 0, fmt.Errorf("%w: i%d", errUnsupportedIntWidth, bitSize)
	}
	mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(bitSize)), big.NewInt(1))
	t := new(big.Int).And(x, mask)
	// Sign-extend if the high bit of the field is set.
	sign := new(big.Int).Lsh(big.NewInt(1), uint(bitSize-1))
	if t.Cmp(sign) >= 0 {
		t.Sub(t, new(big.Int).Lsh(big.NewInt(1), uint(bitSize)))
	}
	if !t.IsInt64() {
		return 0, fmt.Errorf("%w: %v", errIntConstTooLarge, x)
	}
	return t.Int64(), nil
}

// formatICmp translates an icmp predicate and operands to a Go comparison expr.
func formatICmp(pred enum.IPred, xVal, yVal value.Value) (*jen.Statement, error) {
	var op string
	format := FormatValue
	switch pred {
	case enum.IPredEQ:
		op = "=="
	case enum.IPredNE:
		op = "!="
	case enum.IPredSGE:
		op = ">="
		format = FormatSigned
	case enum.IPredSGT:
		op = ">"
		format = FormatSigned
	case enum.IPredSLE:
		op = "<="
		format = FormatSigned
	case enum.IPredSLT:
		op = "<"
		format = FormatSigned
	case enum.IPredUGE:
		op = ">="
		format = FormatUnsigned
	case enum.IPredUGT:
		op = ">"
		format = FormatUnsigned
	case enum.IPredULE:
		op = "<="
		format = FormatUnsigned
	case enum.IPredULT:
		op = "<"
		format = FormatUnsigned
	default:
		return nil, fmt.Errorf("%w: %v", errUnsupportedICmpPred, pred)
	}
	x, err := format(xVal)
	if err != nil {
		return nil, fmt.Errorf("error translating left operand (%v): %w", xVal, err)
	}
	y, err := format(yVal)
	if err != nil {
		return nil, fmt.Errorf("error translating right operand (%v): %w", yVal, err)
	}
	return bin(x, op, y), nil
}

// formatZExt is the expression form of zext (usable in constant expressions).
func formatZExt(from value.Value, to types.Type) (*jen.Statement, error) {
	toType, ok := to.(*types.IntType)
	if !ok {
		return nil, fmt.Errorf("%w: %T", errUnsupportedZextTo, to)
	}
	src, err := FormatUnsigned(from)
	if err != nil {
		return nil, fmt.Errorf("error translating source (%v): %w", from, err)
	}
	w := goIntBits(toType.BitSize)
	if fromType, ok := from.Type().(*types.IntType); ok && fromType.BitSize == 1 {
		// bool → int expression (no statements in constant-expr context).
		return jen.Map(jen.Bool()).Add(goIntType(w)).Values(jen.Dict{
			jen.True():  jen.Lit(1),
			jen.False(): jen.Lit(0),
		}).Index(src), nil
	}
	return conv(goIntType(w), conv(goUintType(w), src)), nil
}

// formatSExt is the expression form of sext.
func formatSExt(from value.Value, to types.Type) (*jen.Statement, error) {
	toType, ok := to.(*types.IntType)
	if !ok {
		return nil, fmt.Errorf("%w: %T", errUnsupportedZextTo, to)
	}
	src, err := FormatSigned(from)
	if err != nil {
		return nil, fmt.Errorf("error translating source (%v): %w", from, err)
	}
	return conv(goIntType(toType.BitSize), src), nil
}

func formatBinConst(op string, x, y constant.Constant) (expr, error) {
	l, err := FormatValue(x)
	if err != nil {
		return expr{}, fmt.Errorf("error translating left (%v): %w", x, err)
	}
	r, err := FormatValue(y)
	if err != nil {
		return expr{}, fmt.Errorf("error translating right (%v): %w", y, err)
	}
	return val(jen.Parens(bin(l, op, r))), nil
}

// formatTrunc is the expression form of trunc.
func formatTrunc(from value.Value, to types.Type) (*jen.Statement, error) {
	toSpec, err := TypeSpec(to)
	if err != nil {
		return nil, fmt.Errorf("error translating To type (%v): %w", to, err)
	}
	src, err := FormatValue(from)
	if err != nil {
		return nil, fmt.Errorf("error translating source (%v): %w", from, err)
	}
	if intType, ok := to.(*types.IntType); ok && intType.BitSize == 1 {
		return jen.Parens(bin(src, "&", jen.Lit(1))).Op("!=").Lit(0), nil
	}
	if intType, ok := to.(*types.IntType); ok && intType.BitSize < 8 {
		return jen.Byte().Call(bin(src, "&", litUntyped(int64(255>>(8-intType.BitSize))))), nil
	}
	return conv(toSpec, src), nil
}

// FormatSigned is like FormatValue, except that it converts "byte" to "int8".
func FormatSigned(v value.Value) (*jen.Statement, error) {
	result, err := FormatValue(v)
	if err != nil {
		return nil, err
	}

	if ci, ok := v.(*constant.Int); ok {
		if ci.Typ.BitSize == 8 {
			return litUntyped(int64(int8(ci.X.Int64()))), nil
		}
		return result, nil
	}

	if t, ok := v.Type().(*types.IntType); ok && t.BitSize == 8 {
		return jen.Int8().Call(result), nil
	}
	return result, nil
}

// FormatUnsigned is like FormatValue, except that it converts integer types to
// unsigned.
func FormatUnsigned(v value.Value) (*jen.Statement, error) {
	result, err := FormatValue(v)
	if err != nil {
		return nil, err
	}

	if ci, ok := v.(*constant.Int); ok {
		if ci.Typ.BitSize > 64 {
			return nil, fmt.Errorf("%w: i%d constant %v", errUnsupportedIntWidth, ci.Typ.BitSize, ci.X)
		}
		var n uint64
		switch {
		case ci.X.IsUint64():
			n = ci.X.Uint64()
		case ci.X.IsInt64():
			n = uint64(ci.X.Int64())
			switch goIntBits(ci.Typ.BitSize) {
			case 8:
				return jen.Byte().Call(litUntyped(int64(byte(n)))), nil
			case 16:
				return jen.Uint16().Call(litUntyped(int64(uint16(n)))), nil
			case 32:
				return jen.Uint32().Call(litUntyped(int64(uint32(n)))), nil
			default:
				// Unsigned decimal: uint64(int64(neg)) is a Go constant overflow.
				return litUint64(n), nil
			}
		default:
			// Low bitSize bits as unsigned.
			mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(ci.Typ.BitSize)), big.NewInt(1))
			t := new(big.Int).And(ci.X, mask)
			if !t.IsUint64() {
				return nil, fmt.Errorf("%w: %v", errIntConstTooLarge, ci.X)
			}
			n = t.Uint64()
		}

		if ci.Typ.BitSize == 1 {
			if n != 0 {
				return jen.True(), nil
			}
			return jen.False(), nil
		}
		switch goIntBits(ci.Typ.BitSize) {
		case 8:
			return litUntyped(int64(byte(n))), nil
		case 16:
			return litUntyped(int64(uint16(n))), nil
		case 32:
			return litUntyped(int64(uint32(n))), nil
		case 64:
			// Untyped 18446744073709551615 overflows int64 in Go source.
			return litUint64(n), nil
		default:
			return litUntyped(int64(n)), nil
		}
	}

	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize > 8 {
			return conv(goUintType(t.BitSize), result), nil
		}
	case *types.PointerType:
		return uintptrOfPtr(result), nil
	}

	return result, nil
}
