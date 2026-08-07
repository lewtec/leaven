package main

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// GetElementPtr translates a getelementptr expression.
func GetElementPtr(elemType types.Type, src value.Value, indices []value.Value) (string, error) {
	srcPointerType, ok := src.Type().(*types.PointerType)
	if !ok {
		return "", fmt.Errorf("%w: %v", errNonPointerSource, src.Type())
	}
	if !types.Equal(srcPointerType.ElemType, elemType) {
		return "", errMismatchedSrcElem
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

	source, err := FormatValue(src)
	if err != nil {
		return "", fmt.Errorf("error translating source pointer (%q): %w", src, err)
	}
	result := source

	if !zeroFirstIndex {
		firstIndex, err := FormatValue(indices[0])
		if err != nil {
			return "", fmt.Errorf("error translating first index (%v): %w", indices[0], err)
		}
		result = fmt.Sprintf("libc.AddPointer(%s, int(%s))", source, firstIndex)
	}
	result = strings.TrimPrefix(result, "&")

	currentType := elemType

	for _, index := range indices[1:] {
		if ind, ok := index.(*constant.Index); ok {
			index = ind.Constant
		}
		switch ct := currentType.(type) {
		case *types.ArrayType:
			v, err := FormatValue(index)
			if err != nil {
				return "", fmt.Errorf("error translating index (%v): %w", index, err)
			}
			result = fmt.Sprintf("%s[%s]", result, v)
			currentType = ct.ElemType
			takeAddress = true

		case *types.StructType:
			ci, ok := index.(*constant.Int)
			if !ok {
				return "", fmt.Errorf("%w: %v %T", errNonConstantStructIdx, index, index)
			}
			result = fmt.Sprintf("%s.F%v", result, ci.X)
			currentType = ct.Fields[ci.X.Int64()]
			takeAddress = true

		default:
			return "", fmt.Errorf("%w: %v", errUnsupportedIndexType, currentType)
		}
	}

	if takeAddress {
		result = "&" + result
	}

	return result, nil
}
