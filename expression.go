package leaven

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/lewtec/leaven/internal/llir/ir/constant"
	"github.com/lewtec/leaven/internal/llir/ir/types"
	"github.com/lewtec/leaven/internal/llir/ir/value"
)

// GetElementPtr translates a getelementptr expression.
func GetElementPtr(elemType types.Type, src value.Value, indices []value.Value) (expr, error) {
	srcPointerType, ok := src.Type().(*types.PointerType)
	if !ok {
		return expr{}, fmt.Errorf("%w: %v", errNonPointerSource, src.Type())
	}
	// Typed-pointer IR: src must be elemType*. LLVM 15+ also uses
	// `gep i8, T*, i64 n` (byte offset) and `gep %T, ptr, ...` (opaque this).
	byteGEP := types.Equal(elemType, types.I8)
	if !srcPointerType.IsOpaque() && !types.Equal(srcPointerType.ElemType, elemType) && !byteGEP {
		if _, agg := elemType.(*types.StructType); !agg {
			if _, arr := elemType.(*types.ArrayType); !arr {
				return expr{}, errMismatchedSrcElem
			}
		}
	}

	zeroFirstIndex := false
	firstIndex := indices[0]
	if ci, ok := firstIndex.(*constant.Index); ok {
		firstIndex = ci.Constant
	}
	if fi, ok := firstIndex.(*constant.Int); ok {
		switch fi.X.Sign() {
		case 0:
			zeroFirstIndex = true
		}
	}
	takeAddress := false

	srcExpr, err := formatExpr(src)
	if err != nil {
		return expr{}, fmt.Errorf("error translating source pointer (%q): %w", src, err)
	}
	// Keep &src for AddPointer (needs *T). Drop it only when we index into
	// the pointee, matching the old strings.TrimPrefix after that call.
	result := srcExpr.code
	rewrote := false

	// Tagged union pointers are uintptr in Go; cast to real pointer before GEP.
	if pt, ok := src.Type().(*types.PointerType); ok && isTaggedPointerType(pt) {
		elem, err := taggedPointerElem(pt)
		if err != nil {
			return expr{}, err
		}
		result = ptrCast(ptrTyp(elem), result)
		rewrote = true
	}

	// Opaque/mismatched src: cast to *elemType before field/array indexing.
	switch elemType.(type) {
	case *types.StructType, *types.ArrayType:
		if srcPointerType.IsOpaque() || !types.Equal(srcPointerType.ElemType, elemType) {
			et, err := TypeSpec(elemType)
			if err != nil {
				return expr{}, err
			}
			result = ptrCast(ptrTyp(et), result)
			rewrote = true
		}
	}

	if !zeroFirstIndex {
		idx, err := FormatValue(indices[0])
		if err != nil {
			return expr{}, fmt.Errorf("error translating first index (%v): %w", indices[0], err)
		}
		if byteGEP {
			result = libc("AddPointer").Types(jen.Byte()).Call(
				ptrCast(ptrTyp(jen.Byte()), result),
				jen.Int().Call(idx),
			)
		} else {
			result = libc("AddPointer").Call(result, jen.Int().Call(idx))
		}
		rewrote = true
	}
	if !rewrote {
		result = srcExpr.dropAddr()
	}

	currentType := elemType

	for _, index := range indices[1:] {
		if ind, ok := index.(*constant.Index); ok {
			index = ind.Constant
		}
		switch ct := currentType.(type) {
		case *types.ArrayType:
			v, err := FormatValue(index)
			if err != nil {
				return expr{}, fmt.Errorf("error translating index (%v): %w", index, err)
			}
			result = jen.Add(result).Index(v)
			currentType = ct.ElemType
			takeAddress = true

		case *types.StructType:
			ci, ok := index.(*constant.Int)
			if !ok {
				return expr{}, fmt.Errorf("%w: %v %T", errNonConstantStructIdx, index, index)
			}
			result = jen.Add(result).Dot(fieldNameU(ci.X.Uint64()))
			currentType = ct.Fields[ci.X.Int64()]
			takeAddress = true

		default:
			return expr{}, fmt.Errorf("%w: %v", errUnsupportedIndexType, currentType)
		}
	}

	if takeAddress {
		return addrExpr(result), nil
	}
	return val(result), nil
}

func pointerElem(t types.Type) (types.Type, bool) {
	pt, ok := t.(*types.PointerType)
	if !ok {
		return nil, false
	}
	if pt.IsOpaque() {
		return nil, false
	}
	return pt.ElemType, true
}

func typedLoad(src expr, srcVal value.Value, elem types.Type) (*jen.Statement, error) {
	if got, ok := pointerElem(srcVal.Type()); ok && types.Equal(got, elem) {
		return src.load(), nil
	}
	t, err := TypeSpec(elem)
	if err != nil {
		return nil, err
	}
	return deref(ptrCast(ptrTyp(t), src.code)), nil
}

func typedStore(dst expr, dstVal value.Value, elem types.Type, src jen.Code) (*jen.Statement, error) {
	if got, ok := pointerElem(dstVal.Type()); ok && types.Equal(got, elem) {
		return dst.store(src), nil
	}
	t, err := TypeSpec(elem)
	if err != nil {
		return nil, err
	}
	return deref(ptrCast(ptrTyp(t), dst.code)).Op("=").Add(src), nil
}
