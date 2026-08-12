package leaven

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/lewtec/leaven/internal/llir/ir"
	"github.com/lewtec/leaven/internal/llir/ir/types"
)

// taggedPointerTypes are LLVM pointer types that appear as the sole field of a
// C union (e.g. tree-sitter Subtree = union { InlineData; HeapData * }).
// Clang collapses those unions to { T* } in IR, but the bit pattern may be a
// non-pointer tagged value (LSB set). In Go we emit them as uintptr so the GC
// does not chase tags and accidental deref paths stay explicit.
var taggedPointerTypes []types.Type

// collectTaggedPointerTypes finds pointer types used as the only field of a
// union typedef (name starts with "union.").
func collectTaggedPointerTypes(m *ir.Module) {
	taggedPointerTypes = nil
	for _, t := range m.TypeDefs {
		name := t.Name()
		if !strings.HasPrefix(name, "union.") {
			continue
		}
		st, ok := t.(*types.StructType)
		if !ok || len(st.Fields) != 1 {
			continue
		}
		pt, ok := st.Fields[0].(*types.PointerType)
		if !ok || pt.IsOpaque() {
			continue
		}
		// Skip function pointers.
		if _, ok := pt.ElemType.(*types.FuncType); ok {
			continue
		}
		taggedPointerTypes = append(taggedPointerTypes, pt)
	}
}

// isTaggedPointerType reports whether t is a pointer type that may hold a
// tagged non-pointer bit pattern (union field).
func isTaggedPointerType(t types.Type) bool {
	for _, tp := range taggedPointerTypes {
		if types.Equal(t, tp) {
			return true
		}
	}
	return false
}

// taggedPointerElem returns the Go type of the pointee for casts from uintptr
// back to a real pointer (e.g. SubtreeHeapData).
func taggedPointerElem(t types.Type) (*jen.Statement, error) {
	pt, ok := t.(*types.PointerType)
	if !ok {
		return nil, fmt.Errorf("%w: %v", errNotPointerType, t)
	}
	return typeSpecIgnoringTagged(pt.ElemType)
}

func typeSpecIgnoringTagged(t types.Type) (*jen.Statement, error) {
	if ref, ok := libraryTypeRef(t); ok {
		return ref.code(), nil
	}
	if name := TypeName(t); name != "" {
		return jen.Id(name), nil
	}
	return typeDefinitionIgnoringTagged(t)
}

func typeDefinitionIgnoringTagged(t types.Type) (*jen.Statement, error) {
	switch t := t.(type) {
	case *types.PointerType:
		if t.IsOpaque() {
			return jen.Qual("unsafe", "Pointer"), nil
		}
		if _, ok := t.ElemType.(*types.FuncType); ok {
			return TypeDefinition(t.ElemType)
		}
		elem, err := typeSpecIgnoringTagged(t.ElemType)
		if err != nil {
			return nil, err
		}
		return ptrTyp(elem), nil
	default:
		return TypeDefinition(t)
	}
}

// TypeDefinition returns the definition (not just the name) of t.
func TypeDefinition(t types.Type) (*jen.Statement, error) {
	switch t := t.(type) {
	case *types.ArrayType:
		elemType, err := TypeSpec(t.ElemType)
		if err != nil {
			return nil, err
		}
		return jen.Index(litUntyped(int64(t.Len))).Add(elemType), nil

	case *types.FloatType:
		switch t.Kind {
		case types.FloatKindFloat:
			return jen.Float32(), nil
		case types.FloatKindDouble, types.FloatKindX86_FP80:
			return jen.Float64(), nil
		default:
			return nil, fmt.Errorf("%w: %v", errUnsupportedFloatType, t.Kind)
		}

	case *types.FuncType:
		params := make([]jen.Code, 0, len(t.Params)+1)
		for i, p := range t.Params {
			pt, err := TypeSpec(p)
			if err != nil {
				return nil, fmt.Errorf("error converting type of parameter %d (%v): %w", i, p, err)
			}
			params = append(params, pt)
		}
		if t.Variadic {
			params = append(params, jen.Op("...").Interface())
		}
		s := jen.Func().Params(params...)
		if !types.Equal(t.RetType, types.Void) {
			rt, err := TypeSpec(t.RetType)
			if err != nil {
				return nil, fmt.Errorf("error converting return type (%v): %w", t.RetType, err)
			}
			s.Add(rt)
		}
		return s, nil

	case *types.IntType:
		switch {
		case t.BitSize == 1:
			return jen.Bool(), nil
		case t.BitSize <= 8:
			return jen.Byte(), nil
		case t.BitSize <= 64:
			// Bitfields and other non-power-of-two widths (e.g. i24) map to the
			// next wider Go integer type (int16/int32/int64).
			return goIntType(t.BitSize), nil
		case t.BitSize == 128:
			// rustc i128/u128 and TypeId. Two limbs, not int64.
			return libc("I128"), nil
		default:
			// LLVM bitfield loads can be i104 etc.; Go has no wider fixed ints.
			return nil, fmt.Errorf("%w: i%d", errUnsupportedIntWidth, t.BitSize)
		}

	case *types.PointerType:
		// Tagged union pointer field: bag-of-bits, not a GC pointer.
		if isTaggedPointerType(t) {
			return jen.Uintptr(), nil
		}
		// Every LLVM pointer value is unsafe.Pointer. Pointee type lives on
		// load/store/gep, not on the pointer.
		return jen.Qual("unsafe", "Pointer"), nil

	case *types.StructType:
		var fields []jen.Code
		for i, field := range t.Fields {
			fieldType, err := TypeSpec(field)
			if err != nil {
				return nil, fmt.Errorf("error converting type of field %d: %w", i, err)
			}
			fields = append(fields, jen.Id(fieldName(i)).Add(fieldType))
		}
		return jen.Struct(fields...), nil

	case *types.VectorType:
		// <N x i1> stays [N]bool. bitcast to iN packs lanes (see i1VectorBitCast).
		elemType, err := TypeSpec(t.ElemType)
		if err != nil {
			return nil, err
		}
		return jen.Index(litUntyped(int64(t.Len))).Add(elemType), nil

	default:
		return nil, fmt.Errorf("%w: %T", errUnsupportedType, t)
	}
}

// TypeSpec returns the name (if it has one) or the definition of t.
func TypeSpec(t types.Type) (*jen.Statement, error) {
	if t == nil {
		return nil, fmt.Errorf("%w: nil type", errUnsupportedType)
	}
	if ref, ok := libraryTypeRef(t); ok {
		return ref.code(), nil
	}
	if name := TypeName(t); name != "" {
		return jen.Id(name), nil
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

func clangTypeName(t types.Type) string {
	name := t.Name()
	name = strings.TrimPrefix(name, "struct.")
	name = strings.TrimPrefix(name, "union.")
	return name
}

func libraryTypeRef(t types.Type) (goRef, bool) {
	name := clangTypeName(t)
	if name == "" {
		return goRef{}, false
	}
	ref, ok := libraryTypes[name]
	return ref, ok
}

// TypeName returns t's local Go name, or the empty string if t is unnamed or
// mapped to a standard-library type. Clang anonymous structs become names like
// "anon.1"; those are sanitized to legal Go identifiers (anon_1).
func TypeName(t types.Type) string {
	name := clangTypeName(t)

	if name == "" || name == "anon" {
		return ""
	}

	if _, ok := libraryTypes[name]; ok {
		return ""
	}

	// Sanitize remaining punctuation from LLVM/clang type names
	// (std::__cxx11::basic_string, templates, etc.).
	r := strings.NewReplacer(".", "_", "-", "_", ":", "_", "<", "_", ">", "_", ",", "_", " ", "_", "*", "p")
	name = r.Replace(name)
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

var libraryTypes = map[string]goRef{
	"FILE":     {pkg: "os", name: "File"},
	"_IO_FILE": {pkg: "os", name: "File"}, // glibc / clang struct name for FILE
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

// i1VectorLen is N when t is <N x i1>.
func i1VectorLen(t types.Type) (n uint64, ok bool) {
	vt, ok := t.(*types.VectorType)
	if !ok {
		return 0, false
	}
	it, ok := vt.ElemType.(*types.IntType)
	if !ok || it.BitSize != 1 {
		return 0, false
	}
	return vt.Len, true
}

func isI1Vector(t types.Type) bool {
	_, ok := i1VectorLen(t)
	return ok
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
