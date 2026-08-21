package leaven

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/lewtec/leaven/internal/llir/ir"
	"github.com/lewtec/leaven/internal/llir/ir/constant"
	"github.com/lewtec/leaven/internal/llir/ir/enum"
	"github.com/lewtec/leaven/internal/llir/ir/types"
	"github.com/lewtec/leaven/internal/llir/ir/value"
	"github.com/lewtec/leaven/libc"
)

// moduleFuncNames / moduleTypeNames are filled by collectModuleNames so locals
// can be renamed when they would shadow a function or type in Go.
// moduleFuncAliases maps C++ alias names (D1→D2, C1→C2) to the real function:
// vtables forward-ref aliases before the alias line is parsed, so the parser
// plants i8 @D1 globals; emit must still produce a real function pointer.
var (
	moduleFuncNames   map[string]bool
	moduleTypeNames   map[string]bool
	moduleFuncAliases map[string]*ir.Func
)

// collectModuleNames records function and type identifiers used in the module.
func collectModuleNames(m *ir.Module) {
	moduleFuncNames = make(map[string]bool)
	moduleTypeNames = make(map[string]bool)
	moduleFuncAliases = make(map[string]*ir.Func)
	for _, f := range m.Funcs {
		moduleFuncNames[rawIdentName(f)] = true
	}
	for _, t := range m.TypeDefs {
		if n := TypeName(t); n != "" {
			moduleTypeNames[n] = true
		}
	}
	for _, a := range m.Aliases {
		if f, ok := a.Aliasee.(*ir.Func); ok {
			moduleFuncAliases[a.Name()] = f
		}
	}
}

// identPunct maps LLVM/rustc punctuation to `_` so names are Go identifiers.
// rustc on Darwin uses `$` in symbol names (gofmt: illegal character U+0024).
var identPunct = strings.NewReplacer(
	".", "_", "-", "_", "$", "_",
	":", "_", "<", "_", ">", "_",
	",", "_", " ", "_", "*", "p",
	";", "_", "[", "_", "]", "_",
	"(", "_", ")", "_", "'", "_",
	"\"", "_", "/", "_", "\\", "_",
	"=", "_", "+", "_", "&", "_",
	"|", "_", "!", "_", "?", "_",
	"@", "_", "#", "_", "%", "_",
	"^", "_", "~", "_", "`", "_",
)

func sanitizeIdent(name string) string {
	name = identPunct.Replace(name)
	if name == "" {
		return "_"
	}
	if c := name[0]; '0' <= c && c <= '9' {
		name = "v" + name
	}
	if invalidNames[name] {
		name = "_" + name
	}
	return name
}

// rawIdentName sanitizes an LLVM name to a Go identifier without clash renames.
func rawIdentName(v value.Named) string {
	name := v.Name()
	if name == "" {
		return "v" + strings.TrimPrefix(v.Ident(), "%")
	}
	return sanitizeIdent(name)
}

// funcLocalNames disambiguates Go names inside one function (%0 and %v0
// both become v0). Filled by collectFuncLocalNames.
var funcLocalNames map[value.Named]string

func collectFuncLocalNames(fn *ir.Func) {
	funcLocalNames = make(map[value.Named]string)
	used := map[string]bool{rawIdentName(fn): true}
	add := func(v value.Named) {
		if v == nil || types.Equal(v.Type(), types.Void) {
			return
		}
		if _, ok := funcLocalNames[v]; ok {
			return
		}
		base := variableNameBase(v)
		name := base
		for i := 2; used[name]; i++ {
			name = fmt.Sprintf("%s_%d", base, i)
		}
		used[name] = true
		funcLocalNames[v] = name
	}
	for _, p := range fn.Params {
		add(p)
	}
	for _, b := range fn.Blocks {
		for _, inst := range b.Insts {
			if n, ok := inst.(value.Named); ok {
				add(n)
			}
		}
		if n, ok := b.Term.(value.Named); ok {
			add(n)
		}
	}
}

// VariableName returns the name to use for a local variable or parameter.
func VariableName(v value.Named) string {
	if funcLocalNames != nil {
		if n, ok := funcLocalNames[v]; ok {
			return n
		}
	}
	return variableNameBase(v)
}

func variableNameBase(v value.Named) string {
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
	return sanitizeIdent(name)
}

// invalidNames are Go keywords and predeclared ids that cannot be used as names.
// https://go.dev/ref/spec#Keywords
var invalidNames = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
	// Common predeclared / special. "_" is the blank identifier (not a value).
	"_":    true,
	"init": true, "true": true, "false": true, "iota": true, "nil": true,
	"bool": true, "byte": true, "error": true, "int": true, "int8": true,
	"int16": true, "int32": true, "int64": true, "rune": true, "string": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"uintptr": true, "float32": true, "float64": true, "complex64": true,
	"complex128": true, "any": true, "comparable": true,
}

func namedRef(name string) (*jen.Statement, bool) {
	if ref, ok := libcLookup(name); ok {
		return ref.code(), true
	}
	if c, ok := cxxIONamed(name); ok {
		return c, true
	}
	if c, ok := cxxTreeNamed(name); ok {
		return c, true
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
	if _, ok := libcLookup(name); ok {
		return true
	}
	if _, ok := libraryGlobals[name]; ok {
		return true
	}
	if _, ok := namedRef(name); ok {
		return true
	}
	if isGetline(name) || isLibcxxStringEqCStr(name) || isLibcxxStringCompareCStr(name) ||
		isLibcxxStringErase(name) || isLibcxxStringAppendCStr(name) ||
		isLibcxxStringAssignCStr(name) {
		return true
	}
	return llvmCallHandled(name)
}

func compositeValues(elems []jen.Code) *jen.Statement {
	return jen.Values(elems...)
}

func formatComposite(typ *jen.Statement, elems []jen.Code) *jen.Statement {
	return jen.Add(typ).Add(compositeValues(elems))
}

func fnPtrBitcast(from *jen.Statement) *jen.Statement {
	return Sym(libc.FuncCode[func()]).Call(from)
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
		return val(fnPtrBitcast(fn)), nil

	case *ir.Alias:
		if v.Aliasee == nil {
			return expr{}, fmt.Errorf("alias %s has no aliasee", v.Name())
		}
		return formatExpr(v.Aliasee)

	case *ir.Global:
		// Forward-ref stub for a function alias (see moduleFuncAliases).
		if f := moduleFuncAliases[v.Name()]; f != nil {
			return formatExpr(f)
		}
		name := VariableName(v)
		if types.IsFunc(v.ContentType) {
			fn := jen.Id(name)
			if c, ok := namedRef(name); ok {
				fn = c
			}
			return val(fnPtrBitcast(fn)), nil
		}
		if ref, ok := libraryGlobals[name]; ok {
			return addrExpr(ref.code()), nil
		}
		return addrExpr(jen.Id(name)), nil

	case value.Named:
		// Locals/params keep their VariableName. namedRef is only for
		// globals/funcs (e.g. @read → libc.Read); a local %read must not.
		return ident(VariableName(v)), nil

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
		if packed, err := i1VectorBitCast(from, v.From.Type(), v.To); packed != nil || err != nil {
			if err != nil {
				return expr{}, err
			}
			return val(packed), nil
		}
		if vec, err := vectorBitCast(from, v.From.Type(), v.To); vec != nil || err != nil {
			if err != nil {
				return expr{}, err
			}
			return val(vec), nil
		}
		if bits, err := scalarBitCast(from, v.From.Type(), v.To); bits != nil || err != nil {
			if err != nil {
				return expr{}, err
			}
			return val(bits), nil
		}
		if isTaggedPointerType(v.To) {
			return val(ptrToUint(from)), nil
		}
		if isTaggedPointerType(v.From.Type()) {
			to, err := TypeSpec(v.To)
			if err != nil {
				return expr{}, fmt.Errorf("error translating type (%v): %w", v.To, err)
			}
			return val(jen.Parens(to).Call(emitUP(from))), nil
		}
		// Pointer values are already unsafe.Pointer.
		return val(from), nil

	case *constant.ExprIntToPtr:
		// Unsigned bit pattern: uintptr(negative int64) is a Go constant overflow.
		from, err := FormatUnsigned(v.From)
		if err != nil {
			return expr{}, fmt.Errorf("error translating source (%v): %w", v.From, err)
		}
		to, err := TypeSpec(v.To)
		if err != nil {
			return expr{}, fmt.Errorf("error translating type (%v): %w", v.To, err)
		}
		return val(jen.Parens(to).Call(emitUP(jen.Uintptr().Call(from)))), nil

	case *constant.ExprPtrToInt:
		from, err := FormatValue(v.From)
		if err != nil {
			return expr{}, fmt.Errorf("error translating source (%v): %w", v.From, err)
		}
		to, err := TypeSpec(v.To)
		if err != nil {
			return expr{}, fmt.Errorf("error translating type (%v): %w", v.To, err)
		}
		return val(conv(to, ptrToUint(from))), nil

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
		if _, ok := wideBits(v.X.Type()); ok {
			return formatBinConst("ashr", v.X, v.Y)
		}
		return formatBinConst(">>", v.X, v.Y)

	case *constant.Float:
		result := v.X.String()
		var c *jen.Statement
		switch result {
		case "+Inf":
			c = Sym(math.Inf).Call(jen.Lit(1))
		case "-Inf":
			c = Sym(math.Inf).Call(jen.Lit(-1))
		case "NaN":
			c = Sym(math.NaN).Call()
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
		if isWide128(v.Typ.BitSize) {
			return val(i128Lit(v.X)), nil
		}
		if isWide256(v.Typ.BitSize) {
			return val(i256Lit(v.X)), nil
		}
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
		if it, ok := v.Typ.(*types.IntType); ok && it.BitSize == 1 {
			return val(jen.False()), nil
		}
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
	switch t := typ.(type) {
	case *types.ArrayType, *types.StructType, *types.VectorType:
		ts, err := TypeSpec(typ)
		if err != nil {
			return expr{}, fmt.Errorf("error translating type (%v): %w", typ, err)
		}
		return val(jen.Add(ts).Values()), nil
	case *types.IntType:
		if t.BitSize == 1 {
			return val(jen.False()), nil
		}
		if isWide128(t.BitSize) {
			return val(Qual[libc.I128]().Values()), nil
		}
		if isWide256(t.BitSize) {
			return val(Qual[libc.I256]().Values()), nil
		}
		return val(jen.Lit(0)), nil
	case *types.FloatType:
		return val(jen.Lit(0)), nil
	case *types.PointerType:
		if isTaggedPointerType(t) {
			return val(jen.Lit(0)), nil
		}
		return val(jen.Nil()), nil
	default:
		return expr{}, fmt.Errorf("%w: %v", errUnsupportedUndefType, typ)
	}
}

var libraryGlobals = map[string]goRef{
	"stdin":                                 {pkg: "os", name: "Stdin"},
	"stdout":                                {pkg: "os", name: "Stdout"},
	"stderr":                                {pkg: "os", name: "Stderr"},
	"_ZTVN10__cxxabiv117__class_type_infoE": {pkg: libcPath, name: "ClassTypeInfoVT"},
	"_ZTVN10__cxxabiv120__si_class_type_infoE":  {pkg: libcPath, name: "SIClassTypeInfoVT"},
	"_ZTVN10__cxxabiv121__vmi_class_type_infoE": {pkg: libcPath, name: "VMIClassTypeInfoVT"},
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

func i128Lit(x *big.Int) *jen.Statement {
	lo, hi := i128Limbs(x)
	return Qual[libc.I128]().Values(jen.Dict{
		jen.Id("Lo"): litUint64(lo),
		jen.Id("Hi"): litUint64(hi),
	})
}

func i128Limbs(x *big.Int) (lo, hi uint64) {
	mod := new(big.Int).Lsh(big.NewInt(1), 128)
	u := new(big.Int).Mod(x, mod)
	lo = new(big.Int).And(u, new(big.Int).SetUint64(^uint64(0))).Uint64()
	hi = new(big.Int).Rsh(u, 64).Uint64()
	return lo, hi
}

func i256Lit(x *big.Int) *jen.Statement {
	mod := new(big.Int).Lsh(big.NewInt(1), 256)
	u := new(big.Int).Mod(x, mod)
	loMask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))
	lo := new(big.Int).And(u, loMask)
	hi := new(big.Int).Rsh(u, 128)
	return Qual[libc.I256]().Values(jen.Dict{
		jen.Id("Lo"): i128Lit(lo),
		jen.Id("Hi"): i128Lit(hi),
	})
}

func wideICmp(bits uint64, pred enum.IPred) (goRef, bool) {
	var i128, i256 any
	switch pred {
	case enum.IPredEQ:
		i128, i256 = libc.I128Eq, libc.I256Eq
	case enum.IPredNE:
		i128, i256 = libc.I128Ne, libc.I256Ne
	case enum.IPredSGE:
		i128, i256 = libc.I128Sge, libc.I256Sge
	case enum.IPredSGT:
		i128, i256 = libc.I128Sgt, libc.I256Sgt
	case enum.IPredSLE:
		i128, i256 = libc.I128Sle, libc.I256Sle
	case enum.IPredSLT:
		i128, i256 = libc.I128Slt, libc.I256Slt
	case enum.IPredUGE:
		i128, i256 = libc.I128Uge, libc.I256Uge
	case enum.IPredUGT:
		i128, i256 = libc.I128Ugt, libc.I256Ugt
	case enum.IPredULE:
		i128, i256 = libc.I128Ule, libc.I256Ule
	case enum.IPredULT:
		i128, i256 = libc.I128Ult, libc.I256Ult
	default:
		return goRef{}, false
	}
	return wideSym(bits, i128, i256), true
}

// formatICmp translates an icmp predicate and operands to a Go comparison expr.
func formatICmp(pred enum.IPred, xVal, yVal value.Value) (*jen.Statement, error) {
	if bits, ok := wideBits(xVal.Type()); ok {
		fn, ok := wideICmp(bits, pred)
		if !ok {
			return nil, fmt.Errorf("%w: %v", errUnsupportedICmpPred, pred)
		}
		x, err := FormatValue(xVal)
		if err != nil {
			return nil, fmt.Errorf("error translating left operand (%v): %w", xVal, err)
		}
		y, err := FormatValue(yVal)
		if err != nil {
			return nil, fmt.Errorf("error translating right operand (%v): %w", yVal, err)
		}
		return fn.Call(x, y), nil
	}
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
	// Go forbids ordered compares on unsafe.Pointer; always use uintptr.
	if _, ok := xVal.Type().(*types.PointerType); ok {
		format = FormatUnsigned
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
	if isWide128(toType.BitSize) {
		if fromType, ok := from.Type().(*types.IntType); ok && fromType.BitSize == 1 {
			return Sym(libc.I128FromU64).Call(jen.Map(jen.Bool()).Uint64().Values(jen.Dict{
				jen.True():  jen.Lit(1),
				jen.False(): jen.Lit(0),
			}).Index(src)), nil
		}
		return Sym(libc.I128FromU64).Call(jen.Uint64().Call(src)), nil
	}
	if isWide256(toType.BitSize) {
		if fromType, ok := from.Type().(*types.IntType); ok && fromType.BitSize == 1 {
			return Sym(libc.I256FromU64).Call(jen.Map(jen.Bool()).Uint64().Values(jen.Dict{
				jen.True():  jen.Lit(1),
				jen.False(): jen.Lit(0),
			}).Index(src)), nil
		}
		if isI128(from.Type()) {
			return Sym(libc.I256FromI128).Call(src), nil
		}
		return Sym(libc.I256FromU64).Call(jen.Uint64().Call(src)), nil
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
	if isWide128(toType.BitSize) {
		if fromType, ok := from.Type().(*types.IntType); ok && fromType.BitSize == 1 {
			return Sym(libc.I128FromI64).Call(jen.Map(jen.Bool()).Int64().Values(jen.Dict{
				jen.True():  jen.Lit(-1),
				jen.False(): jen.Lit(0),
			}).Index(src)), nil
		}
		return Sym(libc.I128FromI64).Call(jen.Int64().Call(src)), nil
	}
	if isWide256(toType.BitSize) {
		if fromType, ok := from.Type().(*types.IntType); ok && fromType.BitSize == 1 {
			return Sym(libc.I256FromI64).Call(jen.Map(jen.Bool()).Int64().Values(jen.Dict{
				jen.True():  jen.Lit(-1),
				jen.False(): jen.Lit(0),
			}).Index(src)), nil
		}
		if isI128(from.Type()) {
			return Sym(libc.I256FromI128S).Call(src), nil
		}
		return Sym(libc.I256FromI64).Call(jen.Int64().Call(src)), nil
	}
	if fromType, ok := from.Type().(*types.IntType); ok && fromType.BitSize == 1 {
		// Go has no int32(bool). sext i1: true → -1.
		return boolToInt(src, toType.BitSize, true), nil
	}
	return conv(goIntType(toType.BitSize), src), nil
}

// boolToInt is zext/sext of i1. signed: true→-1; unsigned: true→1.
func boolToInt(src *jen.Statement, bits uint64, signed bool) *jen.Statement {
	t := 1
	if signed {
		t = -1
	}
	// byte cannot hold untyped -1; use 255 for i8 all-ones.
	trueLit := jen.Lit(t)
	if signed && goIntBits(bits) == 8 {
		trueLit = jen.Lit(255)
	}
	return jen.Map(jen.Bool()).Add(goIntType(goIntBits(bits))).Values(jen.Dict{
		jen.True():  trueLit,
		jen.False(): jen.Lit(0),
	}).Index(src)
}

// refersToGlobal reports whether c mentions g (e.g. SSO string: ptr GEP into self).
func refersToGlobal(v value.Value, g *ir.Global) bool {
	return constRefs(v, func(x value.Value) bool {
		gg, ok := x.(*ir.Global)
		return ok && gg == g
	})
}

// mentionsFunc reports whether a constant mentions a function (vtable slots).
// Go treats `var vt = dtor` as depending on dtor; if dtor mentions vt, cycle.
func mentionsFunc(v value.Value) bool {
	return constRefs(v, func(x value.Value) bool {
		_, ok := x.(*ir.Func)
		return ok
	})
}

// constRefs walks a constant. Globals and functions are leaves.
func constRefs(v value.Value, pred func(value.Value) bool) bool {
	if v == nil {
		return false
	}
	if pred(v) {
		return true
	}
	switch x := v.(type) {
	case *ir.Alias:
		return constRefs(x.Aliasee, pred)
	case *ir.Arg:
		return constRefs(x.Value, pred)
	case *ir.IFunc:
		return constRefs(x.Resolver, pred)
	case *constant.Struct:
		for _, f := range x.Fields {
			if constRefs(f, pred) {
				return true
			}
		}
	case *constant.Array:
		for _, e := range x.Elems {
			if constRefs(e, pred) {
				return true
			}
		}
	case *constant.Vector:
		for _, e := range x.Elems {
			if constRefs(e, pred) {
				return true
			}
		}
	case *constant.Index:
		return constRefs(x.Constant, pred)
	case *constant.ExprGetElementPtr:
		if constRefs(x.Src, pred) {
			return true
		}
		for _, idx := range x.Indices {
			if constRefs(idx, pred) {
				return true
			}
		}
	case *constant.ExprBitCast:
		return constRefs(x.From, pred)
	case *constant.ExprIntToPtr:
		return constRefs(x.From, pred)
	case *constant.ExprPtrToInt:
		return constRefs(x.From, pred)
	case *constant.ExprTrunc:
		return constRefs(x.From, pred)
	case *constant.ExprZExt:
		return constRefs(x.From, pred)
	case *constant.ExprSExt:
		return constRefs(x.From, pred)
	case *constant.ExprAdd:
		return constRefs(x.X, pred) || constRefs(x.Y, pred)
	case *constant.ExprSub:
		return constRefs(x.X, pred) || constRefs(x.Y, pred)
	case *constant.ExprMul:
		return constRefs(x.X, pred) || constRefs(x.Y, pred)
	case *constant.ExprAnd:
		return constRefs(x.X, pred) || constRefs(x.Y, pred)
	case *constant.ExprOr:
		return constRefs(x.X, pred) || constRefs(x.Y, pred)
	case *constant.ExprXor:
		return constRefs(x.X, pred) || constRefs(x.Y, pred)
	case *constant.ExprShl:
		return constRefs(x.X, pred) || constRefs(x.Y, pred)
	case *constant.ExprLShr:
		return constRefs(x.X, pred) || constRefs(x.Y, pred)
	case *constant.ExprAShr:
		return constRefs(x.X, pred) || constRefs(x.Y, pred)
	case *constant.ExprICmp:
		return constRefs(x.X, pred) || constRefs(x.Y, pred)
	case *constant.ExprSelect:
		return constRefs(x.Cond, pred) || constRefs(x.X, pred) || constRefs(x.Y, pred)
	case *constant.ExprExtractValue:
		return constRefs(x.X, pred)
	case *constant.ExprInsertValue:
		return constRefs(x.X, pred) || constRefs(x.Elem, pred)
	}
	return false
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
	if bits, ok := wideBits(x.Type()); ok {
		fn, ok := wideBinFunc(bits, op, false)
		if !ok {
			return expr{}, fmt.Errorf("%w: i%d %s", errUnsupportedInstruction, bits, op)
		}
		return val(fn.Call(l, r)), nil
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
	if isI128(from.Type()) {
		it, ok := to.(*types.IntType)
		if !ok {
			return conv(toSpec, src), nil
		}
		switch {
		case it.BitSize == 1:
			return Sym(libc.I128TruncI1).Call(src), nil
		case it.BitSize <= 8:
			return Sym(libc.I128TruncI8).Call(src), nil
		case it.BitSize <= 16:
			return Sym(libc.I128TruncI16).Call(src), nil
		case it.BitSize <= 32:
			return Sym(libc.I128TruncI32).Call(src), nil
		default:
			return Sym(libc.I128TruncI64).Call(src), nil
		}
	}
	if isI256(from.Type()) {
		it, ok := to.(*types.IntType)
		if !ok {
			return conv(toSpec, src), nil
		}
		switch {
		case it.BitSize == 1:
			return Sym(libc.I256TruncI1).Call(src), nil
		case it.BitSize <= 8:
			return Sym(libc.I256TruncI8).Call(src), nil
		case it.BitSize <= 16:
			return Sym(libc.I256TruncI16).Call(src), nil
		case it.BitSize <= 32:
			return Sym(libc.I256TruncI32).Call(src), nil
		case it.BitSize == 128:
			return Sym(libc.I256TruncI128).Call(src), nil
		default:
			return Sym(libc.I256TruncI64).Call(src), nil
		}
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
		if isWide128(ci.Typ.BitSize) || isWide256(ci.Typ.BitSize) {
			return result, nil
		}
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
		if isWide128(t.BitSize) || isWide256(t.BitSize) {
			return result, nil
		}
		if t.BitSize > 8 {
			return conv(goUintType(t.BitSize), result), nil
		}
	case *types.PointerType:
		return ptrToUint(result), nil
	}

	return result, nil
}
