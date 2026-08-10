package leaven

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/lewtec/leaven/internal/llir/ir"
	"github.com/lewtec/leaven/internal/llir/ir/constant"
	"github.com/lewtec/leaven/internal/llir/ir/enum"
	"github.com/lewtec/leaven/internal/llir/ir/types"
	"github.com/lewtec/leaven/internal/llir/ir/value"
)

func translateOp(v value.Value, what string) (*jen.Statement, error) {
	s, err := FormatValue(v)
	if err != nil {
		return nil, fmt.Errorf("error translating %s (%v): %w", what, v, err)
	}
	return s, nil
}

func translateType(t types.Type, what string) (*jen.Statement, error) {
	s, err := TypeSpec(t)
	if err != nil {
		return nil, fmt.Errorf("error translating %s (%v): %w", what, t, err)
	}
	return s, nil
}

// TranslateInstruction translates an LLVM instruction to one or more Go statements.
// A nil slice means the instruction is a no-op in Go.
func TranslateInstruction(inst ir.Instruction) ([]jen.Code, error) {
	switch inst := inst.(type) {
	case *ir.InstAdd:
		x, err := translateOp(inst.X, "left operand")
		if err != nil {
			return nil, err
		}
		y, err := translateOp(inst.Y, "right operand")
		if err != nil {
			return nil, err
		}
		name := VariableName(inst)
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return one(vectorBin(vecBin{dest: name, op: "+", x: x, y: y})), nil
		}
		if ciy, ok := inst.Y.(*constant.Int); ok && ciy.X.Sign() == -1 {
			// Use the constant's own minus sign.
			return one(jen.Id(name).Op("=").Add(x).Op(ciy.X.String())), nil
		}
		return one(assign(name, bin(x, "+", y))), nil

	case *ir.InstAtomicRMW:
		dst, err := translateOp(inst.Dst, "destination")
		if err != nil {
			return nil, err
		}
		x, err := translateOp(inst.X, "operand")
		if err != nil {
			return nil, err
		}
		name := VariableName(inst)
		addFn, ok := atomicAddFunc(inst.Type())
		if !ok {
			return nil, fmt.Errorf("%w: atomicrmw on %v", errUnsupportedInstruction, inst.Type())
		}
		switch inst.Op {
		case enum.AtomicOpAdd:
			// atomicrmw returns the old value; Add* returns the new value.
			return one(assign(name, bin(addFn.Call(dst, x), "-", x))), nil
		case enum.AtomicOpSub:
			return one(assign(name, bin(addFn.Call(dst, jen.Op("-").Parens(x)), "+", x))), nil
		case enum.AtomicOpXChg:
			swapFn, ok := atomicSwapFunc(inst.Type())
			if !ok {
				return nil, fmt.Errorf("%w: atomicrmw xchg on %v", errUnsupportedInstruction, inst.Type())
			}
			return one(assign(name, swapFn.Call(dst, x))), nil
		default:
			return nil, fmt.Errorf("%w: atomicrmw %v", errUnsupportedInstruction, inst.Op)
		}

	case *ir.InstAlloca:
		t, err := translateType(inst.ElemType, "type")
		if err != nil {
			return nil, err
		}
		name := VariableName(inst)
		if inst.NElems == nil {
			if _, ok := inst.ElemType.(*types.ArrayType); ok {
				// If it's an array, allocate an extra byte to allow indexing off the end.
				return one(assign(name, addrOf(jen.New(jen.Struct(
					jen.Id("v").Add(t),
					jen.Id("b").Byte(),
				)).Dot("v")))), nil
			}
			newExpr := jen.New(t)
			// Alloca of T yields *T; if *T is a tagged union pointer, store as uintptr.
			// Retain so GC won't free when only a uintptr handle remains.
			if pt, ok := inst.Type().(*types.PointerType); ok && isTaggedPointerType(pt) {
				return one(assign(name, uintptrOfPtr(libc("Retain").Call(newExpr)))), nil
			}
			// If T itself is a tagged pointer type (alloca of the pointer slot).
			if isTaggedPointerType(inst.ElemType) {
				return one(assign(name, jen.New(jen.Uintptr()))), nil
			}
			return one(assign(name, newExpr)), nil
		}
		nElems, err := translateOp(inst.NElems, "NElems")
		if err != nil {
			return nil, err
		}
		return one(assign(name, addrOf(jen.Make(jen.Index().Add(t), bin(nElems, "+", jen.Lit(1))).Index(jen.Lit(0))))), nil

	case *ir.InstAnd:
		x, err := translateOp(inst.X, "left operand")
		if err != nil {
			return nil, err
		}
		y, err := translateOp(inst.Y, "right operand")
		if err != nil {
			return nil, err
		}
		name := VariableName(inst)
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return one(vectorBin(vecBin{dest: name, op: "&", x: x, y: y})), nil
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 1 {
			return one(assign(name, bin(x, "&&", y))), nil
		}
		return one(assign(name, bin(x, "&", y))), nil

	case *ir.InstAShr:
		x, err := FormatSigned(inst.X)
		if err != nil {
			return nil, fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return nil, fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		name := VariableName(inst)
		if t, ok := inst.Typ.(*types.IntType); ok && t.BitSize == 8 {
			return one(assign(name, jen.Byte().Call(bin(x, ">>", y)))), nil
		}
		return one(assign(name, bin(x, ">>", y))), nil

	case *ir.InstBitCast:
		if !compatiblePointerTypes(inst.From.Type(), inst.To) {
			return nil, fmt.Errorf("%w: %v and %v", errIncompatiblePointers, inst.From.Type(), inst.To)
		}
		from, err := translateOp(inst.From, "source")
		if err != nil {
			return nil, err
		}
		to, err := translateType(inst.To, "type")
		if err != nil {
			return nil, err
		}
		name := VariableName(inst)
		if isFuncPointerType(inst.To) || isFuncPointerType(inst.From.Type()) {
			return one(assign(name, fnPtrBitcast(to, from))), nil
		}
		if isTaggedPointerType(inst.To) {
			return one(assign(name, uintptrOfPtr(from))), nil
		}
		return one(assign(name, ptrCast(to, from))), nil

	case *ir.InstCall:
		return translateCall(inst)

	case *ir.InstFreeze:
		x, err := translateOp(inst.X, "freeze operand")
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), x)), nil

	case *ir.InstLandingPad:
		z, err := zeroOf(inst.ResultType)
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), z.code)), nil

	case *ir.InstExtractElement:
		x, err := translateOp(inst.X, "vector")
		if err != nil {
			return nil, err
		}
		index, err := translateOp(inst.Index, "index")
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), jen.Add(x).Index(index))), nil

	case *ir.InstExtractValue:
		x, err := translateOp(inst.X, "aggregate")
		if err != nil {
			return nil, err
		}
		expr, err := formatAggregateIndex(x, inst.X.Type(), inst.Indices)
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), expr)), nil

	case *ir.InstInsertValue:
		x, err := translateOp(inst.X, "aggregate")
		if err != nil {
			return nil, err
		}
		elem, err := translateOp(inst.Elem, "element")
		if err != nil {
			return nil, err
		}
		dest, err := formatAggregateIndex(jen.Id(VariableName(inst)), inst.Typ, inst.Indices)
		if err != nil {
			return nil, err
		}
		return []jen.Code{
			assign(VariableName(inst), x),
			jen.Add(dest).Op("=").Add(elem),
		}, nil

	case *ir.InstFAdd:
		return translateBinAssign(inst, llvmBin{op: "+", x: inst.X, y: inst.Y})

	case *ir.InstFCmp:
		x, err := translateOp(inst.X, "left operand")
		if err != nil {
			return nil, err
		}
		y, err := translateOp(inst.Y, "right operand")
		if err != nil {
			return nil, err
		}
		name := VariableName(inst)
		switch inst.Pred {
		case enum.FPredOEQ:
			return one(assign(name, bin(x, "==", y))), nil
		case enum.FPredOGE:
			return one(assign(name, bin(x, ">=", y))), nil
		case enum.FPredOGT:
			return one(assign(name, bin(x, ">", y))), nil
		case enum.FPredOLE:
			return one(assign(name, bin(x, "<=", y))), nil
		case enum.FPredOLT:
			return one(assign(name, bin(x, "<", y))), nil
		case enum.FPredUNE:
			return one(assign(name, bin(x, "!=", y))), nil
		case enum.FPredORD:
			return one(assign(name, bin(bin(x, "==", x), "&&", bin(y, "==", y)))), nil
		case enum.FPredUNO:
			return one(assign(name, bin(bin(x, "!=", x), "||", bin(y, "!=", y)))), nil
		case enum.FPredUEQ:
			return one(assign(name, bin(bin(bin(x, "!=", x), "||", bin(y, "!=", y)), "||", bin(x, "==", y)))), nil
		case enum.FPredUGT:
			return one(assign(name, bin(bin(bin(x, "!=", x), "||", bin(y, "!=", y)), "||", bin(x, ">", y)))), nil
		case enum.FPredUGE:
			return one(assign(name, bin(bin(bin(x, "!=", x), "||", bin(y, "!=", y)), "||", bin(x, ">=", y)))), nil
		case enum.FPredULT:
			return one(assign(name, bin(bin(bin(x, "!=", x), "||", bin(y, "!=", y)), "||", bin(x, "<", y)))), nil
		case enum.FPredULE:
			return one(assign(name, bin(bin(bin(x, "!=", x), "||", bin(y, "!=", y)), "||", bin(x, "<=", y)))), nil
		case enum.FPredONE:
			return one(assign(name, bin(bin(bin(x, "==", x), "&&", bin(y, "==", y)), "&&", bin(x, "!=", y)))), nil
		default:
			return nil, fmt.Errorf("%w: %v", errUnsupportedICmpPred, inst.Pred)
		}

	case *ir.InstFDiv:
		return translateBinAssign(inst, llvmBin{op: "/", x: inst.X, y: inst.Y})

	case *ir.InstFMul:
		return translateBinAssign(inst, llvmBin{op: "*", x: inst.X, y: inst.Y})

	case *ir.InstFNeg:
		x, err := translateOp(inst.X, "operand")
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), jen.Op("-").Add(x))), nil

	case *ir.InstFPExt, *ir.InstFPTrunc:
		return translateConvInst(inst)

	case *ir.InstFPToSI:
		from, err := translateOp(inst.From, "source")
		if err != nil {
			return nil, err
		}
		to, err := translateType(inst.To, "type")
		if err != nil {
			return nil, err
		}
		name := VariableName(inst)
		if it, ok := inst.To.(*types.IntType); ok && it.BitSize <= 8 && it.BitSize > 1 {
			return one(assign(name, jen.Byte().Call(jen.Int8().Call(from)))), nil
		}
		return one(assign(name, conv(to, from))), nil

	case *ir.InstFSub:
		return translateBinAssign(inst, llvmBin{op: "-", x: inst.X, y: inst.Y})

	case *ir.InstGetElementPtr:
		result, err := GetElementPtr(inst.ElemType, inst.Src, inst.Indices)
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), result.code)), nil

	case *ir.InstICmp:
		cmp, err := formatICmp(inst.Pred, inst.X, inst.Y)
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), cmp)), nil

	case *ir.InstInsertElement:
		x, err := translateOp(inst.X, "initial vector")
		if err != nil {
			return nil, err
		}
		elem, err := translateOp(inst.Elem, "new element")
		if err != nil {
			return nil, err
		}
		index, err := translateOp(inst.Index, "index")
		if err != nil {
			return nil, err
		}
		name := VariableName(inst)
		if _, ok := inst.X.(*constant.Undef); ok {
			return one(jen.Id(name).Index(index).Op("=").Add(elem)), nil
		}
		return []jen.Code{
			assign(name, x),
			jen.Id(name).Index(index).Op("=").Add(elem),
		}, nil

	case *ir.InstIntToPtr:
		from, err := translateOp(inst.From, "source")
		if err != nil {
			return nil, err
		}
		to, err := translateType(inst.To, "type")
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), jen.Parens(to).Call(unsafePtr(jen.Uintptr().Call(from))))), nil

	case *ir.InstLoad:
		src, err := formatExpr(inst.Src)
		if err != nil {
			return nil, fmt.Errorf("error translating source (%v): %w", inst.Src, err)
		}
		val, err := typedLoad(src, inst.Src, inst.ElemType)
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), val)), nil

	case *ir.InstLShr:
		x, err := FormatUnsigned(inst.X)
		if err != nil {
			return nil, fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return nil, fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		name := VariableName(inst)
		if t, ok := inst.Typ.(*types.IntType); ok && t.BitSize > 8 {
			return one(assign(name, conv(goIntType(t.BitSize), bin(x, ">>", y)))), nil
		}
		return one(assign(name, bin(x, ">>", y))), nil

	case *ir.InstMul:
		return translateBinAssign(inst, llvmBin{op: "*", x: inst.X, y: inst.Y})

	case *ir.InstOr:
		x, err := translateOp(inst.X, "left operand")
		if err != nil {
			return nil, err
		}
		y, err := translateOp(inst.Y, "right operand")
		if err != nil {
			return nil, err
		}
		name := VariableName(inst)
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return one(vectorBin(vecBin{dest: name, op: "|", x: x, y: y})), nil
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 1 {
			return one(assign(name, bin(x, "||", y))), nil
		}
		return one(assign(name, bin(x, "|", y))), nil

	case *ir.InstPtrToInt:
		from, err := translateOp(inst.From, "source")
		if err != nil {
			return nil, err
		}
		to, err := translateType(inst.To, "type")
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), conv(to, uintptrOfPtr(from)))), nil

	case *ir.InstSDiv:
		return translateSignedDivRem(inst, llvmBin{op: "/", x: inst.X, y: inst.Y})

	case *ir.InstUDiv:
		return translateUnsignedDivRem(inst, llvmBin{op: "/", x: inst.X, y: inst.Y})

	case *ir.InstSelect:
		cond, err := translateOp(inst.Cond, "condition")
		if err != nil {
			return nil, err
		}
		valueTrue, err := translateOp(inst.ValueTrue, "first operand")
		if err != nil {
			return nil, err
		}
		valueFalse, err := translateOp(inst.ValueFalse, "second operand")
		if err != nil {
			return nil, err
		}
		name := VariableName(inst)
		return one(jen.If(cond).Block(assign(name, valueTrue)).Else().Block(assign(name, valueFalse))), nil

	case *ir.InstSExt:
		toType, ok := inst.To.(*types.IntType)
		if !ok {
			return nil, fmt.Errorf("%w: %T", errUnsupportedZextTo, inst.To)
		}
		from, err := FormatSigned(inst.From)
		if err != nil {
			return nil, fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		return one(assign(VariableName(inst), conv(goIntType(toType.BitSize), from))), nil

	case *ir.InstShl:
		x, err := translateOp(inst.X, "left operand")
		if err != nil {
			return nil, err
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return nil, fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		return one(assign(VariableName(inst), bin(x, "<<", y))), nil

	case *ir.InstShuffleVector:
		x, err := translateOp(inst.X, "left operand")
		if err != nil {
			return nil, err
		}
		y, err := translateOp(inst.Y, "right operand")
		if err != nil {
			return nil, err
		}
		mask, err := translateOp(inst.Mask, "mask")
		if err != nil {
			return nil, err
		}
		length := int64(inst.Typ.Len)
		name := VariableName(inst)
		return one(jen.For(jen.List(jen.Id("i"), jen.Id("m")).Op(":=").Range().Add(mask)).Block(
			jen.If(jen.Id("m").Op("<").Lit(int(length))).Block(
				jen.Id(name).Index(jen.Id("i")).Op("=").Add(x).Index(jen.Id("m")),
			).Else().Block(
				jen.Id(name).Index(jen.Id("i")).Op("=").Add(y).Index(jen.Id("m").Op("-").Lit(int(length))),
			),
		)), nil

	case *ir.InstSIToFP:
		from, err := FormatSigned(inst.From)
		if err != nil {
			return nil, fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		to, err := translateType(inst.To, "type")
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), conv(to, from))), nil

	case *ir.InstSRem:
		return translateSignedDivRem(inst, llvmBin{op: "%", x: inst.X, y: inst.Y})

	case *ir.InstURem:
		return translateUnsignedDivRem(inst, llvmBin{op: "%", x: inst.X, y: inst.Y})

	case *ir.InstStore:
		dest, err := formatExpr(inst.Dst)
		if err != nil {
			return nil, fmt.Errorf("error translating destination (%v): %w", inst.Dst, err)
		}
		src, err := translateOp(inst.Src, "source")
		if err != nil {
			return nil, err
		}
		st, err := typedStore(dest, inst.Dst, inst.Src.Type(), src)
		if err != nil {
			return nil, err
		}
		return one(st), nil

	case *ir.InstSub:
		return translateBinAssign(inst, llvmBin{op: "-", x: inst.X, y: inst.Y})

	case *ir.InstTrunc:
		if vt, ok := inst.To.(*types.VectorType); ok {
			toType, ok := vt.ElemType.(*types.IntType)
			if !ok {
				return nil, fmt.Errorf("%w: %v", errUnsupportedZextTo, inst.To)
			}
			to, err := translateType(toType, "To type")
			if err != nil {
				return nil, err
			}
			from, err := translateOp(inst.From, "source")
			if err != nil {
				return nil, err
			}
			name := VariableName(inst)
			return one(jen.For(jen.List(jen.Id("i"), jen.Id("v")).Op(":=").Range().Add(from)).Block(
				jen.Id(name).Index(jen.Id("i")).Op("=").Add(to).Call(jen.Id("v")),
			)), nil
		}
		to, err := translateType(inst.To, "To type")
		if err != nil {
			return nil, err
		}
		from, err := translateOp(inst.From, "source")
		if err != nil {
			return nil, err
		}
		name := VariableName(inst)
		if intType, ok := inst.To.(*types.IntType); ok && intType.BitSize == 1 {
			return one(assign(name, jen.Parens(bin(from, "&", jen.Lit(1))).Op("!=").Lit(0))), nil
		}
		if intType, ok := inst.To.(*types.IntType); ok && intType.BitSize < 8 {
			return one(assign(name, jen.Byte().Call(bin(from, "&", litUntyped(int64(255>>(8-intType.BitSize))))))), nil
		}
		return one(assign(name, conv(to, from))), nil

	case *ir.InstUIToFP:
		from, err := FormatUnsigned(inst.From)
		if err != nil {
			return nil, fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		to, err := translateType(inst.To, "type")
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), conv(to, from))), nil

	case *ir.InstXor:
		x, err := translateOp(inst.X, "left operand")
		if err != nil {
			return nil, err
		}
		y, err := translateOp(inst.Y, "right operand")
		if err != nil {
			return nil, err
		}
		name := VariableName(inst)
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return one(vectorBin(vecBin{dest: name, op: "^", x: x, y: y})), nil
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 1 {
			return one(assign(name, bin(x, "!=", y))), nil
		}
		return one(assign(name, bin(x, "^", y))), nil

	case *ir.InstZExt:
		if vt, ok := inst.To.(*types.VectorType); ok {
			toType, ok := vt.ElemType.(*types.IntType)
			if !ok {
				return nil, fmt.Errorf("%w: %v", errUnsupportedZextTo, inst.To)
			}
			ft, ok := inst.From.Type().(*types.VectorType)
			if !ok {
				return nil, fmt.Errorf("%w: %v and %v", errMismatchedZextTypes, inst.To, inst.From.Type())
			}
			fromType, ok := ft.ElemType.(*types.IntType)
			if !ok {
				return nil, fmt.Errorf("%w: %v", errUnsupportedZextFrom, inst.From.Type())
			}
			from, err := translateOp(inst.From, "source")
			if err != nil {
				return nil, err
			}
			name := VariableName(inst)
			tw, fw := goIntBits(toType.BitSize), goIntBits(fromType.BitSize)
			return one(jen.For(jen.List(jen.Id("i"), jen.Id("v")).Op(":=").Range().Add(from)).Block(
				jen.Id(name).Index(jen.Id("i")).Op("=").Add(goIntType(tw)).Call(
					goUintType(tw).Call(goUintType(fw).Call(jen.Id("v"))),
				),
			)), nil
		}
		toType, ok := inst.To.(*types.IntType)
		if !ok {
			return nil, fmt.Errorf("%w: %T", errUnsupportedZextTo, inst.To)
		}
		from, err := FormatUnsigned(inst.From)
		if err != nil {
			return nil, fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		name := VariableName(inst)
		if fromType, ok := inst.From.Type().(*types.IntType); ok && fromType.BitSize == 1 {
			return one(jen.If(from).Block(assign(name, jen.Lit(1))).Else().Block(assign(name, jen.Lit(0)))), nil
		}
		w := goIntBits(toType.BitSize)
		return one(assign(name, conv(goIntType(w), conv(goUintType(w), from)))), nil

	default:
		return nil, fmt.Errorf("%w: %T", errUnsupportedInstruction, inst)
	}
}

type llvmBin struct {
	op   string
	x, y value.Value
}

func translateBinAssign(inst value.Named, b llvmBin) ([]jen.Code, error) {
	xv, err := translateOp(b.x, "left operand")
	if err != nil {
		return nil, err
	}
	yv, err := translateOp(b.y, "right operand")
	if err != nil {
		return nil, err
	}
	return one(assign(VariableName(inst), bin(xv, b.op, yv))), nil
}

func translateConvInst(inst ir.Instruction) ([]jen.Code, error) {
	var fromV value.Value
	var toT types.Type
	var named value.Named
	switch inst := inst.(type) {
	case *ir.InstFPExt:
		fromV, toT, named = inst.From, inst.To, inst
	case *ir.InstFPTrunc:
		fromV, toT, named = inst.From, inst.To, inst
	default:
		return nil, fmt.Errorf("%w: %T", errUnsupportedInstruction, inst)
	}
	from, err := translateOp(fromV, "source")
	if err != nil {
		return nil, err
	}
	to, err := translateType(toT, "type")
	if err != nil {
		return nil, err
	}
	return one(assign(VariableName(named), conv(to, from))), nil
}

func translateSignedDivRem(inst value.Named, b llvmBin) ([]jen.Code, error) {
	xv, err := FormatSigned(b.x)
	if err != nil {
		return nil, fmt.Errorf("error translating left operand (%v): %w", b.x, err)
	}
	yv, err := FormatSigned(b.y)
	if err != nil {
		return nil, fmt.Errorf("error translating right operand (%v): %w", b.y, err)
	}
	name := VariableName(inst)
	if intType, ok := inst.Type().(*types.IntType); ok && intType.BitSize == 8 {
		return one(assign(name, jen.Byte().Call(bin(xv, b.op, yv)))), nil
	}
	return one(assign(name, bin(xv, b.op, yv))), nil
}

func translateUnsignedDivRem(inst value.Named, b llvmBin) ([]jen.Code, error) {
	xv, err := FormatUnsigned(b.x)
	if err != nil {
		return nil, fmt.Errorf("error translating left operand (%v): %w", b.x, err)
	}
	yv, err := FormatUnsigned(b.y)
	if err != nil {
		return nil, fmt.Errorf("error translating right operand (%v): %w", b.y, err)
	}
	name := VariableName(inst)
	if intType, ok := inst.Type().(*types.IntType); ok && intType.BitSize == 8 {
		return one(assign(name, jen.Byte().Call(bin(xv, b.op, yv)))), nil
	}
	if intType, ok := inst.Type().(*types.IntType); ok && intType.BitSize > 8 {
		return one(assign(name, conv(goIntType(intType.BitSize), bin(xv, b.op, yv)))), nil
	}
	return one(assign(name, bin(xv, b.op, yv))), nil
}

func asBytePtr(x jen.Code) *jen.Statement {
	return jen.Parens(ptrTyp(jen.Byte())).Call(x)
}

func coerceCallArg(inst *ir.InstCall, i int, a value.Value, got *jen.Statement) *jen.Statement {
	fn, ok := inst.Callee.(*ir.Func)
	if !ok || i >= len(fn.Params) {
		return got
	}
	want, ok := fn.Params[i].Typ.(*types.PointerType)
	if !ok || !want.IsOpaque() {
		return got
	}
	have, ok := a.Type().(*types.PointerType)
	if !ok || have.IsOpaque() {
		return got
	}
	return unsafePtr(got)
}

func calleeLLVMName(v value.Value) string {
	if n, ok := v.(value.Named); ok {
		return VariableName(n)
	}
	return ""
}

func translateCall(inst *ir.InstCall) ([]jen.Code, error) {
	llvmName := calleeLLVMName(inst.Callee)
	switch llvmName {
	case "llvm_lifetime_start_p0", "llvm_lifetime_end_p0",
		"llvm_experimental_noalias_scope_decl", "llvm_assume",
		"llvm_donothing":
		return nil, nil
	}
	args := make([]jen.Code, len(inst.Args))
	for i, a := range inst.Args {
		v, err := FormatValue(a)
		if err != nil {
			return nil, fmt.Errorf("error translating argument %d (%v): %w", i, a, err)
		}
		args[i] = coerceCallArg(inst, i, a, v)
	}

	var callee *jen.Statement
	switch llvmName {
	case "calloc", "malloc":
		callee = jen.Id(llvmName)
		if pt, ok := inst.Typ.(*types.PointerType); ok && !pt.IsOpaque() && pt.ElemType != nil {
			if et, err := TypeSpec(pt.ElemType); err == nil {
				callee = libc(strings.Title(llvmName)).Types(et)
			}
		}
	case "leaven_va_start":
		if len(args) == 1 {
			return one(deref(args[0]).Op("=").Add(ptrCast(ptrTyp(jen.Byte()), addrOf(jen.Id("varargs"))))), nil
		}
	case "llvm_va_start":
		if len(args) == 1 {
			return one(deref(jen.Parens(ptrTyp(jen.Qual("unsafe", "Pointer")))).Call(
				jen.Qual("unsafe", "Add").Call(unsafePtr(args[0]), jen.Lit(8)),
			).Op("=").Add(unsafePtr(addrOf(jen.Id("varargs"))))), nil
		}
	case "llvm_va_end", "llvm_lifetime_start", "llvm_lifetime_end", "llvm_stackrestore":
		return nil, nil
	case "vsnprintf":
		if len(args) == 4 {
			return one(assign(VariableName(inst), libc("Vsnprintf").Call(
				args[0], args[1], args[2],
				ptrCast(ptrTyp(jen.Byte()), args[3]),
			))), nil
		}
	case "ldexp":
		if len(args) == 2 {
			return one(assign(VariableName(inst), jen.Qual("math", "Ldexp").Call(args[0], jen.Int().Call(args[1])))), nil
		}
	case "llvm_fabs_f32":
		if len(args) == 1 {
			return one(assign(VariableName(inst), jen.Float32().Call(jen.Qual("math", "Abs").Call(jen.Float64().Call(args[0]))))), nil
		}
	case "llvm_fmuladd_f64", "llvm_fmuladd_f32":
		if len(args) == 3 {
			return one(assign(VariableName(inst), bin(bin(args[0], "*", args[1]), "+", args[2]))), nil
		}
	case "llvm_memcpy_p0i8_p0i8_i64", "llvm_memmove_p0i8_p0i8_i64",
		"llvm_memcpy_p0_p0_i64", "llvm_memmove_p0_p0_i64":
		return one(libc("Memmove").Call(asBytePtr(args[0]), asBytePtr(args[1]), args[2])), nil
	case "llvm_memset_p0i8_i64", "llvm_memset_p0_i64":
		return one(libc("Memset").Call(asBytePtr(args[0]), args[1], args[2])), nil
	case "llvm_abs_i32":
		if len(args) >= 1 {
			return one(assign(VariableName(inst), libc("AbsI32").Call(args[0]))), nil
		}
	case "llvm_sadd_with_overflow_i32":
		if len(args) == 2 {
			return one(assign(VariableName(inst), libc("SAddWithOverflowI32").Call(args[0], args[1]))), nil
		}
	case "llvm_umax_i64":
		if len(args) == 2 {
			return one(assign(VariableName(inst), libc("UMaxU64").Call(args[0], args[1]))), nil
		}
	case "llvm_ctpop_i64":
		if len(args) == 1 {
			return one(assign(VariableName(inst), jen.Qual("math/bits", "OnesCount64").Call(jen.Uint64().Call(args[0])))), nil
		}
	case "llvm_umul_with_overflow_i64":
		if len(args) == 2 {
			return one(assign(VariableName(inst), libc("UMulWithOverflowU64").Call(args[0], args[1]))), nil
		}
	case "llvm_objectsize_i64_p0i8":
		return one(assign(VariableName(inst), jen.Op("-").Lit(1))), nil
	case "llvm_stacksave":
		return one(assign(VariableName(inst), jen.Nil())), nil
	}

	if llvmName == "printf" && len(args) > 0 {
		args[0] = asBytePtr(args[0])
	}
	switch llvmName {
	case "strlen", "strcpy", "strchr", "strcmp":
		for i := range args {
			args[i] = asBytePtr(args[i])
		}
	}

	if callee == nil {
		if ref, ok := libraryFunctions[llvmName]; ok {
			callee = ref.code()
		} else if strings.Contains(llvmName, "panicking") ||
			strings.Contains(llvmName, "handle_error") ||
			strings.Contains(llvmName, "throw_logic_error") ||
			strings.Contains(llvmName, "throw_length_error") ||
			strings.Contains(llvmName, "throw_bad_alloc") ||
			strings.Contains(llvmName, "throw_bad_array") {
			return one(jen.Panic(jen.Lit("runtime error"))), nil
		} else if strings.Contains(llvmName, "__rust_alloc") && !strings.Contains(llvmName, "dealloc") && !strings.Contains(llvmName, "realloc") && !strings.Contains(llvmName, "zeroed") {
			return one(assign(VariableName(inst), libc("RustAlloc").Call(args...))), nil
		} else if strings.Contains(llvmName, "__rust_dealloc") {
			return one(libc("RustDealloc").Call(args...)), nil
		} else if strings.Contains(llvmName, "__rust_realloc") {
			return one(assign(VariableName(inst), libc("RustRealloc").Call(args...))), nil
		} else if strings.Contains(llvmName, "__rust_no_alloc_shim") {
			return nil, nil
		} else if llvmName == "_Znwm" || llvmName == "_Znam" {
			return one(assign(VariableName(inst), libc("RustAlloc").Call(args[0], jen.Lit(1)))), nil
		} else if strings.HasPrefix(llvmName, "_Zdl") || strings.HasPrefix(llvmName, "_Zda") {
			return one(libc("RustDealloc").Call(args[0], jen.Lit(0), jen.Lit(1))), nil
		} else {
			var err error
			callee, err = FormatValue(inst.Callee)
			if err != nil {
				return nil, fmt.Errorf("error translating callee (%v): %w", inst.Callee, err)
			}
		}
	}

	callExpr := jen.Add(callee).Call(args...)
	if types.Equal(inst.Type(), types.Void) {
		return one(callExpr), nil
	}
	if isTaggedPointerType(inst.Type()) && (llvmName == "malloc" || llvmName == "calloc" || llvmName == "realloc") {
		callExpr = uintptrOfPtr(callExpr)
	}
	return one(assign(VariableName(inst), callExpr)), nil
}

func atomicAddFunc(t types.Type) (*jen.Statement, bool) {
	it, ok := t.(*types.IntType)
	if !ok {
		return nil, false
	}
	switch goIntBits(it.BitSize) {
	case 32:
		return jen.Qual("sync/atomic", "AddInt32"), true
	case 64:
		return jen.Qual("sync/atomic", "AddInt64"), true
	default:
		return nil, false
	}
}

func atomicSwapFunc(t types.Type) (*jen.Statement, bool) {
	it, ok := t.(*types.IntType)
	if !ok {
		return nil, false
	}
	switch goIntBits(it.BitSize) {
	case 32:
		return jen.Qual("sync/atomic", "SwapInt32"), true
	case 64:
		return jen.Qual("sync/atomic", "SwapInt64"), true
	default:
		return nil, false
	}
}

func formatAggregateIndex(base *jen.Statement, t types.Type, indices []uint64) (*jen.Statement, error) {
	curCode := base
	cur := t
	for _, idx := range indices {
		switch ct := cur.(type) {
		case *types.StructType:
			if int(idx) >= len(ct.Fields) {
				return nil, fmt.Errorf("%w: index %d (len %d)", errStructIndexRange, idx, len(ct.Fields))
			}
			curCode = jen.Add(curCode).Dot(fieldNameU(idx))
			cur = ct.Fields[idx]
		case *types.ArrayType:
			curCode = jen.Add(curCode).Index(litUntyped(int64(idx)))
			cur = ct.ElemType
		case *types.VectorType:
			curCode = jen.Add(curCode).Index(litUntyped(int64(idx)))
			cur = ct.ElemType
		default:
			return nil, fmt.Errorf("%w: %T", errUnsupportedAggregate, cur)
		}
	}
	return curCode, nil
}

// rustRuntime maps rustc v0-mangled std/core symbols we implement in libc.
// Crate hashes change; match the stable path suffix.
func rustRuntime(name string) *jen.Statement {
	if strings.Contains(name, "stdio6__print") {
		return libc("RustPrint")
	}
	if !strings.Contains(name, "7Display3fmt") {
		return nil
	}
	i := strings.Index(name, "3fmt3num3imp")
	if i < 0 {
		return nil
	}
	rest := name[i+len("3fmt3num3imp"):]
	if rest == "" {
		return nil
	}
	switch rest[0] {
	case 'l': // i32
		return libc("RustFmtI32")
	case 'j': // usize
		return libc("RustFmtUsize")
	default:
		return nil
	}
}

var libraryFunctions = map[string]goRef{
	"abort":            {libcPath, "Abort"},
	"__assert_fail":    {libcPath, "AssertFail"},
	"fabs":             {"math", "Abs"},
	"__ctype_b_loc":    {libcPath, "CtypeBLoc"},
	"dup":              {libcPath, "Dup"},
	"fclose":           {libcPath, "Fclose"},
	"fdopen":           {libcPath, "Fdopen"},
	"fprintf":          {libcPath, "Fprintf"},
	"fputc":            {libcPath, "Fputc"},
	"fputs":            {libcPath, "Fputs"},
	"free":             {libcPath, "Free"},
	"getchar":          {libcPath, "Getchar"},
	"exit":             {libcPath, "Exit"},
	"iswalnum":         {libcPath, "Iswalnum"},
	"iswalpha":         {libcPath, "Iswalpha"},
	"iswblank":         {libcPath, "Iswblank"},
	"iswcntrl":         {libcPath, "Iswcntrl"},
	"iswdigit":         {libcPath, "Iswdigit"},
	"iswlower":         {libcPath, "Iswlower"},
	"iswspace":         {libcPath, "Iswspace"},
	"iswupper":         {libcPath, "Iswupper"},
	"iswxdigit":        {libcPath, "Iswxdigit"},
	"leaven_va_arg":    {libcPath, "VAArg"},
	"towlower":         {libcPath, "Towlower"},
	"towupper":         {libcPath, "Towupper"},
	"llvm_fabs_f64":    {"math", "Abs"},
	"llvm_fabs_f80":    {"math", "Abs"},
	"llvm_pow_f64":     {"math", "Pow"},
	"memchr":           {libcPath, "Memchr"},
	"memcmp":           {libcPath, "Memcmp"},
	"__memcpy_chk":     {libcPath, "MemcpyChk"},
	"memmove":          {libcPath, "Memmove"},
	"memset_pattern16": {libcPath, "MemsetPattern16"},
	"__memset_chk":     {libcPath, "MemsetChk"},
	"printf":           {libcPath, "Printf"},
	"putc":             {libcPath, "Putc"},
	"putchar":          {libcPath, "Putchar"},
	"puts":             {libcPath, "Puts"},
	"realloc":          {libcPath, "Realloc"},
	"scanf":            {libcPath, "Scanf"},
	"snprintf":         {libcPath, "Snprintf"},
	"sqrt":             {"math", "Sqrt"},
	"__strcat_chk":     {libcPath, "StrcatChk"},
	"strchr":           {libcPath, "Strchr"},
	"strcmp":           {libcPath, "Strcmp"},
	"strcpy":           {libcPath, "Strcpy"},
	"strcspn":          {libcPath, "Strcspn"},
	"strlen":           {libcPath, "Strlen"},
	"strncat":          {libcPath, "Strncat"},
	"strncmp":          {libcPath, "Strncmp"},
	"strncpy":          {libcPath, "Strncpy"},
	"strrchr":          {libcPath, "Strrchr"},
	"strspn":           {libcPath, "Strspn"},
	"strstr":           {libcPath, "Strstr"},
}
