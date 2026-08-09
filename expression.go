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
	if !types.Equal(srcPointerType.ElemType, elemType) {
		return expr{}, errMismatchedSrcElem
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

	if !zeroFirstIndex {
		idx, err := FormatValue(indices[0])
		if err != nil {
			return expr{}, fmt.Errorf("error translating first index (%v): %w", indices[0], err)
		}
		result = libc("AddPointer").Call(result, jen.Int().Call(idx))
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
