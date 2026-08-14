package leaven

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/lewtec/leaven/internal/llir/ir"
	"github.com/lewtec/leaven/internal/llir/ir/constant"
	"github.com/lewtec/leaven/internal/llir/ir/types"
	"github.com/lewtec/leaven/internal/llir/ir/value"
)

// GetElementPtr translates a getelementptr expression.
// The LLVM pointer value is always unsafe.Pointer. *T exists only while
// indexing into a struct/array or as AddPointer's type argument.
func GetElementPtr(elemType types.Type, src value.Value, indices []value.Value) (expr, error) {
	if vt, ok := src.Type().(*types.VectorType); ok {
		return vectorGEP(vecGEP{elem: elemType, src: src, indices: indices, n: vt.Len})
	}
	// Scalar ptr + vector index → vector of pointers (rustc stride loops).
	if len(indices) > 0 {
		idx0 := indices[0]
		if ci, ok := idx0.(*constant.Index); ok {
			idx0 = ci.Constant
		}
		if ivt, ok := idx0.Type().(*types.VectorType); ok {
			return scalarPtrVectorIndexGEP(vecGEP{elem: elemType, src: src, indices: indices, n: ivt.Len})
		}
	}
	if _, ok := src.Type().(*types.PointerType); !ok {
		return expr{}, fmt.Errorf("%w: %v", errNonPointerSource, src.Type())
	}
	if len(indices) == 0 {
		return expr{}, fmt.Errorf("%w: no indices", errUnsupportedIndexType)
	}

	srcExpr, err := formatExpr(src)
	if err != nil {
		return expr{}, fmt.Errorf("error translating source pointer (%q): %w", src, err)
	}

	firstIndex := indices[0]
	if ci, ok := firstIndex.(*constant.Index); ok {
		firstIndex = ci.Constant
	}
	zeroFirst := false
	if fi, ok := firstIndex.(*constant.Int); ok && fi.X.Sign() == 0 {
		zeroFirst = true
	}

	if len(indices) == 1 && zeroFirst {
		return val(srcExpr.code), nil
	}

	et, err := TypeSpec(elemType)
	if err != nil {
		return expr{}, err
	}
	// Element stride must use LLVM ABI size, not Go sizeof. Nested ZST
	// fields (Rust newtypes / PhantomData as struct{}) inflate Go structs
	// (e.g. OsString nest 40 vs LLVM 24) and break argv/Vec GEPs.
	elemSize, err := llvmTypeSize(elemType)
	if err != nil {
		return expr{}, err
	}

	var result *jen.Statement
	if zeroFirst && srcExpr.base != nil && len(indices) > 1 {
		result = srcExpr.dropAddr()
	} else {
		base := srcExpr.code
		if isTaggedPointerType(src.Type()) {
			telem, err := taggedPointerElem(src.Type())
			if err != nil {
				return expr{}, err
			}
			base = emitAs(telem, emitUP(base))
		} else if len(indices) > 1 {
			// Field GEPs need a typed *T for .F0 chains.
			base = emitAs(et, base)
		} else {
			// Pure element GEP: byte stride with LLVM size.
			base = emitAs(Qual[byte](), base)
		}
		if zeroFirst {
			result = base
		} else {
			idx, err := FormatValue(firstIndex)
			if err != nil {
				return expr{}, fmt.Errorf("error translating first index (%v): %w", firstIndex, err)
			}
			if len(indices) > 1 {
				result = emitAddPtr(et, base, jen.Int().Call(idx))
			} else {
				off := jen.Int().Call(idx).Op("*").Lit(int(elemSize))
				result = emitAddPtr(Qual[byte](), base, off)
			}
		}
	}

	currentType := elemType
	takeAddress := false
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
			fi := ci.X.Int64()
			if fi < 0 || int(fi) >= len(ct.Fields) {
				return expr{}, fmt.Errorf("%w: field %d of %v", errUnsupportedIndexType, fi, ct)
			}
			currentType = ct.Fields[fi]
			if isZeroSizeType(currentType) {
				// ZST omitted from Go struct; GEP address = base + ABI offset.
				off, err := llvmFieldOffset(ct, fi)
				if err != nil {
					return expr{}, err
				}
				var base *jen.Statement
				if takeAddress {
					base = emitAs(Qual[byte](), emitPtr(addrOf(result)))
				} else {
					base = emitAs(Qual[byte](), emitPtr(result))
				}
				result = emitAddPtr(Qual[byte](), base, jen.Lit(int(off)))
				takeAddress = false
			} else {
				result = jen.Add(result).Dot(fieldNameU(ci.X.Uint64()))
				takeAddress = true
			}

		default:
			return expr{}, fmt.Errorf("%w: %v", errUnsupportedIndexType, currentType)
		}
	}

	if takeAddress {
		return addrExpr(result), nil
	}
	return val(emitPtr(result)), nil
}

type vecGEP struct {
	elem    types.Type
	src     value.Value
	indices []value.Value
	n       uint64
}

// scalarPtrVectorIndexGEP is getelementptr T, ptr %p, <N x iK> %idx.
func scalarPtrVectorIndexGEP(g vecGEP) (expr, error) {
	if len(g.indices) != 1 {
		return expr{}, fmt.Errorf("%w: vector index gep extra indices", errUnsupportedIndexType)
	}
	srcExpr, err := FormatValue(g.src)
	if err != nil {
		return expr{}, err
	}
	idx := g.indices[0]
	if ci, ok := idx.(*constant.Index); ok {
		idx = ci.Constant
	}
	idxExpr, err := FormatValue(idx)
	if err != nil {
		return expr{}, err
	}
	elemSize, err := llvmTypeSize(g.elem)
	if err != nil {
		return expr{}, err
	}
	n := int64(g.n)
	elems := make([]jen.Code, n)
	for i := int64(0); i < n; i++ {
		// Fresh src each lane — jennifer Statements are mutable builders.
		off := jen.Int().Call(jen.Add(idxExpr).Index(litUntyped(i))).Op("*").Lit(int(elemSize))
		elems[i] = emitOff(jen.Add(srcExpr), off)
	}
	return val(jen.Index(litUntyped(n)).Qual("unsafe", "Pointer").Values(elems...)), nil
}

// vectorGEP is getelementptr on a vector of pointers. Each lane is offset
// by the (broadcast) index in units of elemType.
func vectorGEP(g vecGEP) (expr, error) {
	if len(g.indices) == 0 {
		return expr{}, fmt.Errorf("%w: no indices", errUnsupportedIndexType)
	}
	srcExpr, err := FormatValue(g.src)
	if err != nil {
		return expr{}, err
	}
	idx := g.indices[0]
	if ci, ok := idx.(*constant.Index); ok {
		idx = ci.Constant
	}
	idxExpr, err := FormatValue(idx)
	if err != nil {
		return expr{}, err
	}
	elemSize, err := llvmTypeSize(g.elem)
	if err != nil {
		return expr{}, err
	}
	_, idxVec := idx.Type().(*types.VectorType)
	n := int64(g.n)
	elems := make([]jen.Code, n)
	for i := int64(0); i < n; i++ {
		lane := jen.Add(srcExpr).Index(litUntyped(i))
		off := idxExpr
		if idxVec {
			off = jen.Add(idxExpr).Index(litUntyped(i))
		}
		byteOff := jen.Int().Call(off).Op("*").Lit(int(elemSize))
		elems[i] = emitOff(lane, byteOff)
	}
	if len(g.indices) > 1 {
		return expr{}, fmt.Errorf("%w: vector gep extra indices", errUnsupportedIndexType)
	}
	return val(jen.Index(litUntyped(n)).Qual("unsafe", "Pointer").Values(elems...)), nil
}

// wholeVarAccess reports whether addrExpr.base is the entire LLVM object of
// type elem. A store T to @g when g is a struct is a first-field overlay
// (LangRef: store T at that address), not `g = T`.
func wholeVarAccess(v value.Value, elem types.Type) bool {
	g, ok := v.(*ir.Global)
	if !ok {
		return true
	}
	return types.Equal(g.ContentType, elem)
}

// isScalarPtrObject reports a load/store of a whole ptr cell (global or
// alloca), not an array/struct slot. Those cells stay unsafe.Pointer.
func isScalarPtrObject(v value.Value) bool {
	switch x := v.(type) {
	case *ir.Global:
		_, ok := x.ContentType.(*types.PointerType)
		return ok
	case *ir.InstAlloca:
		_, ok := x.ElemType.(*types.PointerType)
		return ok
	default:
		return false
	}
}

func overlayMem(addr *jen.Statement, addrTy, elem types.Type) (*jen.Statement, error) {
	t, err := memSlotType(elem)
	if err != nil {
		return nil, err
	}
	if isTaggedPointerType(addrTy) {
		addr = emitUP(addr)
	}
	return deref(emitAs(t, addr)), nil
}

func typedLoad(src expr, srcVal value.Value, elem types.Type) (*jen.Statement, error) {
	// cout/cerr are objects. load ptr from @cout is the vptr at offset 0,
	// not the object as a Go slice or as a pointer value.
	if g, ok := srcVal.(*ir.Global); ok && isStdStream(VariableName(g)) {
		if _, ok := elem.(*types.PointerType); ok {
			slot, err := overlayMem(src.code, srcVal.Type(), elem)
			if err != nil {
				return nil, err
			}
			return fromMemSlotExpr(elem, slot), nil
		}
	}
	if src.base != nil && wholeVarAccess(srcVal, elem) {
		loaded := src.load()
		// Library FILE* globals are *os.File; LLVM pointer values are
		// already unsafe.Pointer.
		if pt, ok := elem.(*types.PointerType); ok && !isTaggedPointerType(pt) {
			if g, ok := srcVal.(*ir.Global); ok {
				if ref, ok := libraryGlobals[VariableName(g)]; ok && ref.pkg == "os" {
					return emitPtr(loaded), nil
				}
			}
			// Scalar ptr globals/allocas are unsafe.Pointer. Array/struct
			// slots are uint64 and need PtrFromBits.
			if !isScalarPtrObject(srcVal) {
				return fromMemSlotExpr(elem, loaded), nil
			}
		}
		return loaded, nil
	}
	slot, err := overlayMem(src.code, srcVal.Type(), elem)
	if err != nil {
		return nil, err
	}
	return fromMemSlotExpr(elem, slot), nil
}

type storeDst struct {
	dst  expr
	val  value.Value
	elem types.Type
}

func typedStore(d storeDst, src jen.Code) (*jen.Statement, error) {
	// Named scalar ptr cell (global/alloca) is unsafe.Pointer. Overlay
	// and aggregate slots are uint64.
	scalarCell := d.dst.base != nil && wholeVarAccess(d.val, d.elem) && isScalarPtrObject(d.val)
	if isTaggedPointerType(d.elem) {
		src = ptrToUint(src)
	} else if isMemPtr(d.elem) && !scalarCell {
		src = asMemSlotCode(d.elem, src)
	}
	if d.dst.base != nil && wholeVarAccess(d.val, d.elem) {
		return d.dst.store(src), nil
	}
	slot, err := overlayMem(d.dst.code, d.val.Type(), d.elem)
	if err != nil {
		return nil, err
	}
	return jen.Add(slot).Op("=").Add(src), nil
}
