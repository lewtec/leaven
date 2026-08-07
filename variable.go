package main

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
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
		if n := TypeName(t); n != "" && !strings.Contains(n, ".") {
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
	name := block.Name()
	if name == "" {
		return "block" + strings.TrimPrefix(block.Ident(), "%")
	}
	if c := name[0]; '0' <= c && c <= '9' {
		name = "block" + name
	}
	name = strings.Replace(name, ".", "_", -1)
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

// FormatValue formats a constant or variable as it should appear in an expression.
func FormatValue(v value.Value) (string, error) {
	switch v := v.(type) {
	case *ir.Global:
		name := VariableName(v)
		if types.IsFunc(v.ContentType) {
			// Function used as a value (initializer / assignment), not a call.
			if renamed, ok := libraryFunctions[name]; ok {
				return renamed, nil
			}
			return name, nil
		}
		if renamed, ok := libraryGlobals[name]; ok {
			name = renamed
		}
		return "&" + name, nil

	case value.Named:
		name := VariableName(v)
		// Declarations like @free used as fn values (e.g. current_free = free).
		if renamed, ok := libraryFunctions[name]; ok {
			return renamed, nil
		}
		return name, nil

	case *ir.Arg:
		return FormatValue(v.Value)

	case *constant.Array:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", v.Typ, err)
		}
		b := new(bytes.Buffer)
		if len(v.Elems) < 16 {
			b.WriteString(t)
			b.WriteByte('{')
			for i, c := range v.Elems {
				if i > 0 {
					b.WriteString(", ")
				}
				e, err := FormatValue(c)
				if err != nil {
					return "", fmt.Errorf("error translating element %d (%v): %w", i, c, err)
				}
				fmt.Fprint(b, e)
			}
			b.WriteByte('}')
		} else {
			b.WriteString(t)
			b.WriteString("{\n\t")
			for i, c := range v.Elems {
				if i > 0 {
					if i%16 == 0 {
						b.WriteString(",\n\t")
					} else {
						b.WriteString(", ")
					}
				}
				e, err := FormatValue(c)
				if err != nil {
					return "", fmt.Errorf("error translating element %d (%v): %w", i, c, err)
				}
				fmt.Fprint(b, e)
			}
			b.WriteString(",\n}")
		}
		return b.String(), nil

	case *constant.CharArray:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", v.Typ, err)
		}
		b := new(bytes.Buffer)
		if len(v.X) < 16 {
			b.WriteString(t)
			b.WriteByte('{')
			for i, c := range v.X {
				if i > 0 {
					b.WriteString(", ")
				}
				fmt.Fprintf(b, "%d", c)
			}
			b.WriteByte('}')
		} else {
			b.WriteString(t)
			b.WriteString("{\n\t")
			for i, c := range v.X {
				if i > 0 {
					if i%16 == 0 {
						b.WriteString(",\n\t")
					} else {
						b.WriteString(", ")
					}
				}
				fmt.Fprintf(b, "%d", c)
			}
			b.WriteString(",\n}")
		}
		return b.String(), nil

	case *constant.ExprBitCast:
		from, err := FormatValue(v.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", v.From, err)
		}
		to, err := TypeSpec(v.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", v.To, err)
		}
		// Go forbids unsafe.Pointer(funcValue). Reinterpret via address of a temp.
		if isFuncPointerType(v.To) || isFuncPointerType(v.From.Type()) {
			return fmt.Sprintf("func() %s { tmp := %s; return *(*%s)(unsafe.Pointer(&tmp)) }()", to, from, to), nil
		}
		if isTaggedPointerType(v.To) {
			return fmt.Sprintf("uintptr(unsafe.Pointer(%s))", from), nil
		}
		if isTaggedPointerType(v.From.Type()) {
			return fmt.Sprintf("(%s)(unsafe.Pointer(%s))", to, from), nil
		}
		return fmt.Sprintf("(%s)(unsafe.Pointer(%s))", to, from), nil

	case *constant.ExprIntToPtr:
		from, err := FormatValue(v.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", v.From, err)
		}
		to, err := TypeSpec(v.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", v.To, err)
		}
		return fmt.Sprintf("(%s)(unsafe.Pointer(uintptr(%s)))", to, from), nil

	case *constant.ExprPtrToInt:
		from, err := FormatValue(v.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", v.From, err)
		}
		to, err := TypeSpec(v.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", v.To, err)
		}
		return fmt.Sprintf("%s(uintptr(unsafe.Pointer(%s)))", to, from), nil

	case *constant.ExprGetElementPtr:
		indices := make([]value.Value, len(v.Indices))
		for i, index := range v.Indices {
			indices[i] = index
		}
		return GetElementPtr(v.ElemType, v.Src, indices)

	case *constant.ExprICmp:
		return formatICmp(v.Pred, v.X, v.Y)

	case *constant.ExprZExt:
		return formatZExt(v.From, v.To)

	case *constant.ExprSExt:
		return formatSExt(v.From, v.To)

	case *constant.ExprTrunc:
		return formatTrunc(v.From, v.To)

	case *constant.Float:
		result := v.X.String()
		special := false
		switch result {
		case "+Inf":
			result = "math.Inf(1)"
			special = true
		case "-Inf":
			result = "math.Inf(-1)"
			special = true
		case "NaN":
			result = "math.NaN()"
			special = true
		}
		if special && v.Typ.Kind == types.FloatKindFloat {
			result = fmt.Sprintf("float32(%s)", result)
		}
		return result, nil

	case *constant.Index:
		return FormatValue(v.Constant)

	case *constant.Int:
		if v.Typ.BitSize > 64 {
			return "", fmt.Errorf("%w: i%d constant %v", errUnsupportedIntWidth, v.Typ.BitSize, v.X)
		}
		var value int64
		switch {
		case v.X.IsInt64():
			value = v.X.Int64()
		case v.X.IsUint64():
			// Reinterpret as two's-complement signed for this bit width
			// (e.g. i64 all-ones → -1, not 18446744073709551615).
			value = int64(v.X.Uint64())
		default:
			// Truncate wide math/big values into the declared bit width (≤64).
			truncated, err := intFromBig(v.X, v.Typ.BitSize)
			if err != nil {
				return "", err
			}
			value = truncated
		}

		if v.Typ.BitSize == 1 {
			if value != 0 {
				return "true", nil
			}
			return "false", nil
		}
		switch goIntBits(v.Typ.BitSize) {
		case 8:
			return fmt.Sprint(byte(value)), nil
		case 16:
			return fmt.Sprint(int16(value)), nil
		case 32:
			return fmt.Sprint(int32(value)), nil
		case 64:
			// Typed so large bit patterns never appear as untyped ints.
			return fmt.Sprintf("int64(%d)", value), nil
		default:
			return fmt.Sprint(value), nil
		}

	case *constant.Null:
		// Tagged union pointers are uintptr; use 0 not nil.
		if isTaggedPointerType(v.Typ) {
			return "0", nil
		}
		return "nil", nil

	case *constant.Struct:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", v.Typ, err)
		}
		b := new(bytes.Buffer)
		b.WriteString(t)
		b.WriteByte('{')
		for i, c := range v.Fields {
			if i > 0 {
				b.WriteString(", ")
			}
			e, err := FormatValue(c)
			if err != nil {
				return "", fmt.Errorf("error translating field %d (%v): %w", i, c, err)
			}
			fmt.Fprint(b, e)
		}
		b.WriteByte('}')
		return b.String(), nil

	case *constant.Undef:
		switch v.Typ.(type) {
		case *types.ArrayType, *types.StructType, *types.VectorType:
			t, err := TypeSpec(v.Typ)
			if err != nil {
				return "", fmt.Errorf("error translating type (%v): %w", v.Typ, err)
			}
			return t + "{}", nil
		case *types.IntType, *types.FloatType:
			return "0", nil
		case *types.PointerType:
			return "nil", nil
		default:
			return "", fmt.Errorf("%w: %v", errUnsupportedUndefType, v.Typ)
		}

	case *constant.Vector:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", v.Typ, err)
		}
		b := new(bytes.Buffer)
		b.WriteString(t)
		b.WriteByte('{')
		for i, c := range v.Elems {
			if i > 0 {
				b.WriteString(", ")
			}
			e, err := FormatValue(c)
			if err != nil {
				return "", fmt.Errorf("error translating element %d (%v): %w", i, c, err)
			}
			fmt.Fprint(b, e)
		}
		b.WriteByte('}')
		return b.String(), nil

	case *constant.ZeroInitializer:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", v.Typ, err)
		}
		return t + "{}", nil

	default:
		return "", fmt.Errorf("%w: %T", errUnsupportedValueType, v)
	}
}

var libraryGlobals = map[string]string{
	"stdin":  "os.Stdin",
	"stdout": "os.Stdout",
	"stderr": "os.Stderr",
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
func formatICmp(pred enum.IPred, xVal, yVal value.Value) (string, error) {
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
		return "", fmt.Errorf("%w: %v", errUnsupportedICmpPred, pred)
	}
	x, err := format(xVal)
	if err != nil {
		return "", fmt.Errorf("error translating left operand (%v): %w", xVal, err)
	}
	y, err := format(yVal)
	if err != nil {
		return "", fmt.Errorf("error translating right operand (%v): %w", yVal, err)
	}
	return fmt.Sprintf("%s %s %s", x, op, y), nil
}

// formatZExt is the expression form of zext (usable in constant expressions).
func formatZExt(from value.Value, to types.Type) (string, error) {
	toType, ok := to.(*types.IntType)
	if !ok {
		return "", fmt.Errorf("%w: %T", errUnsupportedZextTo, to)
	}
	src, err := FormatUnsigned(from)
	if err != nil {
		return "", fmt.Errorf("error translating source (%v): %w", from, err)
	}
	w := goIntBits(toType.BitSize)
	if fromType, ok := from.Type().(*types.IntType); ok && fromType.BitSize == 1 {
		// bool → int expression (no statements in constant-expr context).
		return fmt.Sprintf("map[bool]int%d{true: 1, false: 0}[%s]", w, src), nil
	}
	return fmt.Sprintf("int%d(uint%d(%s))", w, w, src), nil
}

// formatSExt is the expression form of sext.
func formatSExt(from value.Value, to types.Type) (string, error) {
	toType, ok := to.(*types.IntType)
	if !ok {
		return "", fmt.Errorf("%w: %T", errUnsupportedZextTo, to)
	}
	src, err := FormatSigned(from)
	if err != nil {
		return "", fmt.Errorf("error translating source (%v): %w", from, err)
	}
	return fmt.Sprintf("int%d(%s)", goIntBits(toType.BitSize), src), nil
}

// formatTrunc is the expression form of trunc.
func formatTrunc(from value.Value, to types.Type) (string, error) {
	toSpec, err := TypeSpec(to)
	if err != nil {
		return "", fmt.Errorf("error translating To type (%v): %w", to, err)
	}
	src, err := FormatValue(from)
	if err != nil {
		return "", fmt.Errorf("error translating source (%v): %w", from, err)
	}
	if intType, ok := to.(*types.IntType); ok && intType.BitSize == 1 {
		return fmt.Sprintf("(%s & 1) != 0", src), nil
	}
	if intType, ok := to.(*types.IntType); ok && intType.BitSize < 8 {
		return fmt.Sprintf("byte(%s & %d)", src, 255>>(8-intType.BitSize)), nil
	}
	return fmt.Sprintf("%s(%s)", toSpec, src), nil
}

// FormatSigned is like FormatValue, except that it converts "byte" to "int8".
func FormatSigned(v value.Value) (string, error) {
	result, err := FormatValue(v)
	if err != nil {
		return "", err
	}

	if ci, ok := v.(*constant.Int); ok {
		if ci.Typ.BitSize == 8 {
			return fmt.Sprint(int8(ci.X.Int64())), nil
		}
		return result, nil
	}

	if t, ok := v.Type().(*types.IntType); ok && t.BitSize == 8 {
		return fmt.Sprintf("int8(%s)", result), nil
	}
	return result, nil
}

// FormatUnsigned is like FormatValue, except that it converts integer types to
// unsigned.
func FormatUnsigned(v value.Value) (string, error) {
	result, err := FormatValue(v)
	if err != nil {
		return "", err
	}

	if ci, ok := v.(*constant.Int); ok {
		if ci.Typ.BitSize > 64 {
			return "", fmt.Errorf("%w: i%d constant %v", errUnsupportedIntWidth, ci.Typ.BitSize, ci.X)
		}
		var value uint64
		switch {
		case ci.X.IsUint64():
			value = ci.X.Uint64()
		case ci.X.IsInt64():
			value = uint64(ci.X.Int64())
			switch goIntBits(ci.Typ.BitSize) {
			case 8:
				return fmt.Sprintf("byte(%d)", byte(value)), nil
			case 16:
				return fmt.Sprintf("uint16(%d)", uint16(value)), nil
			case 32:
				return fmt.Sprintf("uint32(%d)", uint32(value)), nil
			default:
				return fmt.Sprintf("uint64(%d)", value), nil
			}
		default:
			// Low bitSize bits as unsigned.
			mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(ci.Typ.BitSize)), big.NewInt(1))
			t := new(big.Int).And(ci.X, mask)
			if !t.IsUint64() {
				return "", fmt.Errorf("%w: %v", errIntConstTooLarge, ci.X)
			}
			value = t.Uint64()
		}

		if ci.Typ.BitSize == 1 {
			if value != 0 {
				return "true", nil
			}
			return "false", nil
		}
		switch goIntBits(ci.Typ.BitSize) {
		case 8:
			return fmt.Sprint(byte(value)), nil
		case 16:
			return fmt.Sprint(uint16(value)), nil
		case 32:
			return fmt.Sprint(uint32(value)), nil
		case 64:
			// Untyped 18446744073709551615 overflows int64 in Go source.
			return fmt.Sprintf("uint64(%d)", value), nil
		default:
			return fmt.Sprint(value), nil
		}
	}

	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize > 8 {
			return fmt.Sprintf("uint%d(%s)", goIntBits(t.BitSize), result), nil
		}
	case *types.PointerType:
		return fmt.Sprintf("uintptr(unsafe.Pointer(%s))", result), nil
	}

	return result, nil
}
