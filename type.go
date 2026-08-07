package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/llir/llvm/ir/types"
)

// TypeDefinition returns the definition (not just the name) of t.
func TypeDefinition(t types.Type) (string, error) {
	switch t := t.(type) {
	case *types.ArrayType:
		elemType, err := TypeSpec(t.ElemType)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[%d]%s", t.Len, elemType), nil

	case *types.FloatType:
		switch t.Kind {
		case types.FloatKindFloat:
			return "float32", nil
		case types.FloatKindDouble, types.FloatKindX86_FP80:
			return "float64", nil
		default:
			return "", fmt.Errorf("%w: %v", errUnsupportedFloatType, t.Kind)
		}

	case *types.FuncType:
		b := new(bytes.Buffer)
		b.WriteString("func(")
		for i, p := range t.Params {
			if i != 0 {
				b.WriteString(", ")
			}
			pt, err := TypeSpec(p)
			if err != nil {
				return "", fmt.Errorf("error converting type of parameter %d (%v): %w", i, p, err)
			}
			b.WriteString(pt)
		}
		b.WriteString(")")
		if !types.Equal(t.RetType, types.Void) {
			b.WriteString(" ")
			rt, err := TypeSpec(t.RetType)
			if err != nil {
				return "", fmt.Errorf("error converting return type (%v): %w", t.RetType, err)
			}
			b.WriteString(rt)
		}
		return b.String(), nil

	case *types.IntType:
		switch {
		case t.BitSize == 1:
			return "bool", nil
		case t.BitSize <= 8:
			return "byte", nil
		case t.BitSize <= 64:
			// Bitfields and other non-power-of-two widths (e.g. i24) map to the
			// next wider Go integer type (int16/int32/int64).
			return fmt.Sprintf("int%d", goIntBits(t.BitSize)), nil
		default:
			// LLVM bitfield loads can be i104 etc.; Go has no wider fixed ints.
			return "", fmt.Errorf("%w: i%d", errUnsupportedIntWidth, t.BitSize)
		}

	case *types.PointerType:
		if _, ok := t.ElemType.(*types.FuncType); ok {
			// Translate a C function pointer type as a Go function type.
			return TypeDefinition(t.ElemType)
		}
		elemType, err := TypeSpec(t.ElemType)
		if err != nil {
			return "", err
		}
		return "*" + elemType, nil

	case *types.StructType:
		b := new(bytes.Buffer)
		b.WriteString("struct {\n")
		for i, field := range t.Fields {
			fieldType, err := TypeSpec(field)
			if err != nil {
				return "", fmt.Errorf("error converting type of field %d: %w", i, err)
			}
			fmt.Fprintf(b, "\tF%d %s\n", i, fieldType)
		}
		b.WriteString("}")
		return b.String(), nil

	case *types.VectorType:
		elemType, err := TypeSpec(t.ElemType)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[%d]%s", t.Len, elemType), nil

	default:
		return "", fmt.Errorf("%w: %T", errUnsupportedType, t)
	}
}

// TypeSpec returns the name (if it has one) or the definition of t.
func TypeSpec(t types.Type) (string, error) {
	if name := TypeName(t); name != "" {
		return name, nil
	}
	return TypeDefinition(t)
}

// goIntBits rounds an LLVM integer width up to a Go integer width
// (8, 16, 32, or 64). Widths above 64 stay as-is for the caller to reject.
func goIntBits(bits uint64) uint64 {
	switch {
	case bits <= 8:
		return 8
	case bits <= 16:
		return 16
	case bits <= 32:
		return 32
	case bits <= 64:
		return 64
	default:
		return bits
	}
}

// TypeName returns t's name, or the empty string if t is not a named type.
// Clang anonymous structs become names like "anon.1"; those are sanitized to
// legal Go identifiers (anon_1). Library renames (e.g. FILE → os.File) are
// returned as-is so callers can detect the package-qualified form.
func TypeName(t types.Type) string {
	name := t.Name()
	name = strings.TrimPrefix(name, "struct.")
	name = strings.TrimPrefix(name, "union.")

	if name == "" || name == "anon" {
		return ""
	}

	if renamed, ok := libraryTypes[name]; ok {
		return renamed
	}

	// Sanitize remaining punctuation from LLVM/clang type names.
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "-", "_")
	if name == "" {
		return ""
	}
	if c := name[0]; '0' <= c && c <= '9' {
		name = "T" + name
	}
	if invalidNames[name] {
		name = "_" + name
	}
	return name
}

var libraryTypes = map[string]string{
	"FILE":     "os.File",
	"_IO_FILE": "os.File", // glibc / clang struct name for FILE
}

// compatiblePointerTypes returns whether casting t1 to t2 is allowed.
// Any pointer→pointer bitcast is accepted: C (and csmith) freely reinterprets
// unions, arrays, and structs via void*/typed pointers; leaven already models
// those as unsafe.Pointer.
func compatiblePointerTypes(t1, t2 types.Type) bool {
	_, ok1 := t1.(*types.PointerType)
	_, ok2 := t2.(*types.PointerType)
	return ok1 && ok2
}

// hasPointers returns whether t contains pointers.
func hasPointers(t types.Type) bool {
	switch t := t.(type) {
	case *types.ArrayType:
		return hasPointers(t.ElemType)
	case *types.FloatType:
		return false
	case *types.FuncType:
		return true
	case *types.IntType:
		return false
	case *types.PointerType:
		return true
	case *types.StructType:
		for _, f := range t.Fields {
			if hasPointers(f) {
				return true
			}
		}
		return false
	case *types.VectorType:
		return hasPointers(t.ElemType)
	default:
		// We don't know if it contains pointers, so we assume it does,
		// since that means we'll be more careful with it.
		return true
	}
}
