package leaven

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"unsafe"

	"github.com/dave/jennifer/jen"
	"github.com/lewtec/leaven/internal/llir/ir"
	"github.com/lewtec/leaven/internal/llir/ir/types"
	"github.com/lewtec/leaven/libc"
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

// structFieldUintptr reports whether a struct field must be uintptr rather
// than a GC pointer. Opaque ptr is not tagged globally (that would wrap
// every pointer); only proven slots are rewritten.
func structFieldUintptr(st *types.StructType, field types.Type) bool {
	pt, ok := field.(*types.PointerType)
	if !ok {
		return false
	}
	if isTaggedPointerType(pt) {
		return true
	}
	// Clang collapses C union { T* } to { ptr }. The bits may be a tag
	// (LSB set), including on LLVM 15+ opaque ptr.
	if strings.HasPrefix(st.Name(), "union.") && len(st.Fields) == 1 {
		return true
	}
	// Packed { ptr, iN } (libstdc++ _Bit_iterator): the integer sits next
	// to the pointer. Go pads the struct so later i32 stores land in a
	// pointer field (csmith GenerateNewGlobal: 0x1 on the stack).
	if st.Packed && packedMixesPtrAndInt(st) {
		return true
	}
	return false
}

func packedMixesPtrAndInt(st *types.StructType) bool {
	var hasPtr, hasInt bool
	for _, f := range st.Fields {
		switch f.(type) {
		case *types.PointerType:
			hasPtr = true
		case *types.IntType:
			hasInt = true
		}
	}
	return hasPtr && hasInt
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
			return Qual[unsafe.Pointer](), nil
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
			return Qual[libc.I128](), nil
		case t.BitSize == 256:
			// rustc core::fmt::num::__fmt_inner widens u128 to i256.
			return Qual[libc.I256](), nil
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
		return Qual[unsafe.Pointer](), nil

	case *types.StructType:
		var fields []jen.Code
		for i, field := range t.Fields {
			// Drop LLVM ZSTs (empty structs / PhantomData). Go struct{}
			// still takes space in arrays and shifts later fields; LLVM size is 0.
			if isZeroSizeType(field) {
				continue
			}
			var fieldType *jen.Statement
			var err error
			if t.Packed && packedMixesPtrAndInt(t) {
				// <{ ptr, i32 }>: Go {uintptr,int32} is 16 bytes (align 8);
				// LLVM packed is 12. Use [8]byte for the ptr slot so the
				// struct is 12 bytes and parent layouts (vector<bool>) match.
				// Field access is via unsafe address (see _Bit_iterator_base).
				if _, ok := field.(*types.PointerType); ok {
					fieldType = jen.Index(jen.Lit(8)).Byte()
				} else {
					fieldType, err = TypeSpec(field)
				}
			} else if structFieldUintptr(t, field) {
				// This slot may hold a tagged non-pointer (union payload or
				// packed ptr+int). Do not put it in a GC pointer field.
				fieldType = jen.Uintptr()
			} else {
				fieldType, err = TypeSpec(field)
			}
			if err != nil {
				return nil, fmt.Errorf("error converting type of field %d: %w", i, err)
			}
			// Keep LLVM field index in the name (F0, F2, …) so GEP .Fi still works.
			fields = append(fields, jen.Id(fieldName(i)).Add(fieldType))
		}
		if len(fields) == 0 {
			return jen.Struct(), nil
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

	// Sanitize remaining punctuation from LLVM/clang/rustc type names
	// (std::__cxx11::basic_string, smallvec::SmallVec<[usize; 2]>, etc.).
	r := strings.NewReplacer(
		".", "_", "-", "_", ":", "_",
		"<", "_", ">", "_", ",", "_",
		" ", "_", "*", "p",
		";", "_", "[", "_", "]", "_",
		"(", "_", ")", "_",
		"'", "_", "\"", "_",
		"/", "_", "\\", "_",
		"=", "_", "+", "_", "&", "_",
		"|", "_", "!", "_", "?", "_",
		"@", "_", "#", "_", "%", "_",
		"^", "_", "~", "_", "`", "_",
	)
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

var osFileRef = func() goRef {
	t := reflect.TypeFor[os.File]()
	return goRef{pkg: t.PkgPath(), name: t.Name()}
}()

var libraryTypes = map[string]goRef{
	"FILE":     osFileRef,
	"_IO_FILE": osFileRef, // glibc / clang struct name for FILE
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

// scalarBitSize is the LLVM bit width of an integer or IEEE float scalar.
// x86_fp80 is mapped to Go float64, so it is not a same-size bitcast source.
func scalarBitSize(t types.Type) (uint64, bool) {
	switch t := t.(type) {
	case *types.IntType:
		if t.BitSize == 0 {
			return 0, false
		}
		return t.BitSize, true
	case *types.FloatType:
		switch t.Kind {
		case types.FloatKindHalf:
			return 16, true
		case types.FloatKindFloat:
			return 32, true
		case types.FloatKindDouble:
			return 64, true
		case types.FloatKindFP128, types.FloatKindPPC_FP128:
			return 128, true
		default:
			return 0, false
		}
	default:
		return 0, false
	}
}

// sameSizeScalars reports whether a bitcast between t1 and t2 is a
// same-width int/float reinterpret (not a pointer cast).
func sameSizeScalars(t1, t2 types.Type) bool {
	a, ok1 := scalarBitSize(t1)
	b, ok2 := scalarBitSize(t2)
	return ok1 && ok2 && a == b
}

// vectorBitSize is the total LLVM bit width of a non-i1 vector.
// <N x i1> is packed separately (see i1VectorBitCast).
func vectorBitSize(t types.Type) (uint64, bool) {
	vt, ok := t.(*types.VectorType)
	if !ok || isI1Vector(t) {
		return 0, false
	}
	elemBits, ok := scalarBitSize(vt.ElemType)
	if !ok {
		return 0, false
	}
	return elemBits * vt.Len, true
}

// sameSizeVectors reports whether a bitcast between t1 and t2 is a
// same-width vector reinterpret (e.g. <16 x i8> to <2 x i64>).
func sameSizeVectors(t1, t2 types.Type) bool {
	a, ok1 := vectorBitSize(t1)
	b, ok2 := vectorBitSize(t2)
	return ok1 && ok2 && a == b && a > 0
}

// sameSizeVectorScalar reports whether a bitcast between t1 and t2 is a
// same-width vector↔scalar reinterpret (e.g. <2 x i64> to i128).
func sameSizeVectorScalar(t1, t2 types.Type) bool {
	va, vok1 := vectorBitSize(t1)
	sa, sok1 := scalarBitSize(t1)
	vb, vok2 := vectorBitSize(t2)
	sb, sok2 := scalarBitSize(t2)
	switch {
	case vok1 && sok2:
		return va == sb && va > 0
	case sok1 && vok2:
		return sa == vb && sa > 0
	default:
		return false
	}
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

// isZeroSizeType reports whether t has LLVM ABI size 0 (empty struct, [0 x T]).
func isZeroSizeType(t types.Type) bool {
	sz, err := llvmTypeSize(t)
	return err == nil && sz == 0
}

// llvmFieldOffset is the ABI byte offset of field i in st.
func llvmFieldOffset(st *types.StructType, i int64) (int64, error) {
	if i < 0 || int(i) >= len(st.Fields) {
		return 0, fmt.Errorf("%w: field %d", errUnsupportedIndexType, i)
	}
	var off int64
	for j := int64(0); j < i; j++ {
		fs, fa, err := llvmSizeAlign(st.Fields[j])
		if err != nil {
			return 0, err
		}
		if !st.Packed && fa > 1 {
			off = alignUp(off, fa)
		}
		off += fs
	}
	if !st.Packed {
		_, fa, err := llvmSizeAlign(st.Fields[i])
		if err != nil {
			return 0, err
		}
		if fa > 1 {
			off = alignUp(off, fa)
		}
	}
	return off, nil
}

// llvmTypeSize is the ABI size of t in bytes under the default x86_64
// System V layout (matches rustc/clang Linux IR without a custom DL).
func llvmTypeSize(t types.Type) (int64, error) {
	sz, _, err := llvmSizeAlign(t)
	return sz, err
}

func alignUp(off, align int64) int64 {
	if align <= 1 {
		return off
	}
	return (off + align - 1) / align * align
}

func llvmSizeAlign(t types.Type) (size, align int64, err error) {
	switch t := t.(type) {
	case *types.IntType:
		switch {
		case t.BitSize <= 8:
			return 1, 1, nil
		case t.BitSize <= 16:
			return 2, 2, nil
		case t.BitSize <= 32:
			return 4, 4, nil
		case t.BitSize <= 64:
			return 8, 8, nil
		case t.BitSize <= 128:
			return 16, 16, nil
		default:
			b := (int64(t.BitSize) + 7) / 8
			return b, 8, nil
		}
	case *types.FloatType:
		switch t.Kind {
		case types.FloatKindFloat:
			return 4, 4, nil
		case types.FloatKindDouble:
			return 8, 8, nil
		case types.FloatKindX86_FP80:
			return 16, 16, nil
		default:
			return 0, 0, fmt.Errorf("%w: %v", errUnsupportedFloatType, t.Kind)
		}
	case *types.PointerType:
		return 8, 8, nil
	case *types.FuncType:
		return 8, 8, nil
	case *types.ArrayType:
		es, ea, err := llvmSizeAlign(t.ElemType)
		if err != nil {
			return 0, 0, err
		}
		return int64(t.Len) * es, ea, nil
	case *types.VectorType:
		if n, ok := i1VectorLen(t); ok {
			// <N x i1> packs to ceil(N/8) bytes in memory forms.
			b := (int64(n) + 7) / 8
			if b == 0 {
				b = 1
			}
			return b, 1, nil
		}
		es, ea, err := llvmSizeAlign(t.ElemType)
		if err != nil {
			return 0, 0, err
		}
		return int64(t.Len) * es, ea, nil
	case *types.StructType:
		if len(t.Fields) == 0 {
			return 0, 1, nil
		}
		var off, maxA int64 = 0, 1
		for _, f := range t.Fields {
			fs, fa, err := llvmSizeAlign(f)
			if err != nil {
				return 0, 0, err
			}
			if fa > maxA {
				maxA = fa
			}
			if !t.Packed && fa > 1 {
				off = alignUp(off, fa)
			}
			off += fs
		}
		if !t.Packed && maxA > 1 {
			off = alignUp(off, maxA)
		}
		if off == 0 {
			return 0, 1, nil
		}
		return off, maxA, nil
	default:
		return 0, 0, fmt.Errorf("%w: size of %T", errUnsupportedType, t)
	}
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
