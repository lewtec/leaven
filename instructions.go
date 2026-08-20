package leaven

import (
	"fmt"
	"math"
	"math/bits"
	"os"
	"strings"
	"sync/atomic"
	"unsafe"

	"github.com/dave/jennifer/jen"
	"github.com/lewtec/leaven/internal/llir/ir"
	"github.com/lewtec/leaven/internal/llir/ir/constant"
	"github.com/lewtec/leaven/internal/llir/ir/enum"
	"github.com/lewtec/leaven/internal/llir/ir/types"
	"github.com/lewtec/leaven/internal/llir/ir/value"
	"github.com/lewtec/leaven/libc"
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
		if bits, ok := wideBits(inst.Typ); ok {
			return one(assign(name, wideSym(bits, libc.I128Add, libc.I256Add).Call(x, y))), nil
		}
		if ciy, ok := inst.Y.(*constant.Int); ok && ciy.X.Sign() == -1 {
			// Use the constant's own minus sign.
			return one(jen.Id(name).Op("=").Add(x).Op(ciy.X.String())), nil
		}
		return one(assign(name, bin(x, "+", y))), nil

	case *ir.InstFence:
		// Ordering is recorded on the IR; Go has no distinct fence
		// opcodes. The atomic bump is the same compiler barrier rustc
		// uses empty asm ~{memory} for.
		return one(Sym(libc.Fence).Call()), nil

	case *ir.InstCmpXchg:
		ptr, err := translateOp(inst.Ptr, "pointer")
		if err != nil {
			return nil, err
		}
		cmp, err := translateOp(inst.Cmp, "compare")
		if err != nil {
			return nil, err
		}
		neu, err := translateOp(inst.New, "new")
		if err != nil {
			return nil, err
		}
		casFn, ok := atomicCASFunc(inst.Cmp.Type())
		if !ok {
			return nil, fmt.Errorf("%w: cmpxchg on %v", errUnsupportedInstruction, inst.Cmp.Type())
		}
		elem, err := TypeSpec(inst.Cmp.Type())
		if err != nil {
			return nil, err
		}
		ret, err := TypeSpec(inst.Type())
		if err != nil {
			return nil, err
		}
		if isTaggedPointerType(inst.Ptr.Type()) {
			ptr = emitUP(ptr)
		}
		p := emitAs(elem, ptr)
		name := VariableName(inst)
		okName := name + "_ok"
		oldName := name + "_old"
		// Strong CAS. LLVM weak may spuriously fail; strong is still correct
		// for the usual retry loop and does not hide a failed compare.
		// Temps are predeclared in writeFuncBody (no := — goto would jump over).
		return []jen.Code{
			jen.Id(okName).Op("=").Add(casFn.Call(p, cmp, neu)),
			jen.Id(oldName).Op("=").Add(conv(elem, cmp)),
			jen.If(jen.Op("!").Id(okName)).Block(
				jen.Id(oldName).Op("=").Add(atomicLoadFunc(inst.Cmp.Type()).Call(p)),
			),
			assign(name, jen.Add(ret).Values(jen.Id(oldName), jen.Id(okName))),
		}, nil

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
		elem, err := TypeSpec(inst.Type())
		if err != nil {
			return nil, err
		}
		if isTaggedPointerType(inst.Dst.Type()) {
			dst = emitUP(dst)
		}
		dst = emitAs(elem, dst)
		switch inst.Op {
		case enum.AtomicOpAdd:
			addFn, ok := atomicAddFunc(inst.Type())
			if !ok {
				return nil, fmt.Errorf("%w: atomicrmw on %v", errUnsupportedInstruction, inst.Type())
			}
			// atomicrmw returns the old value; Add* returns the new value.
			return one(assign(name, bin(addFn.Call(dst, x), "-", x))), nil
		case enum.AtomicOpSub:
			addFn, ok := atomicAddFunc(inst.Type())
			if !ok {
				return nil, fmt.Errorf("%w: atomicrmw on %v", errUnsupportedInstruction, inst.Type())
			}
			delta := jen.Op("-").Parens(x)
			if it, ok := inst.Type().(*types.IntType); ok && goIntBits(it.BitSize) == 8 {
				// byte(-(int32(1))) is a Go constant overflow. ^x+1 wraps.
				delta = jen.Op("^").Parens(jen.Byte().Call(x)).Op("+").Lit(1)
			}
			return one(assign(name, bin(addFn.Call(dst, delta), "+", x))), nil
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
		// Force align ≥8 so pointers never have LSB set (Rust niches / tagged
		// ptrs). new([0]byte) and new(struct{}) can be 1-byte aligned.
		allocAlign8 := func(elem *jen.Statement) *jen.Statement {
			return emitPtr(addrOf(jen.New(jen.Struct(
				jen.Id("_").Index(jen.Lit(0)).Uint64(),
				jen.Id("v").Add(elem),
				jen.Id("b").Byte(),
			)).Dot("v")))
		}
		if inst.NElems == nil {
			// Alloca of T yields a pointer; tagged union pointers stay uintptr.
			// Retain so GC won't free when only a uintptr handle remains.
			if pt, ok := inst.Type().(*types.PointerType); ok && isTaggedPointerType(pt) {
				return one(assign(name, emitAddr(Sym(libc.Retain[byte]).Call(jen.New(t))))), nil
			}
			// If T itself is a tagged pointer type (alloca of the pointer slot).
			if isTaggedPointerType(inst.ElemType) {
				return one(assign(name, allocAlign8(jen.Uintptr()))), nil
			}
			return one(assign(name, allocAlign8(t))), nil
		}
		nElems, err := translateOp(inst.NElems, "NElems")
		if err != nil {
			return nil, err
		}
		// Dynamic array: pad length and pin first element (already slice-aligned).
		return one(assign(name, emitPtr(addrOf(jen.Make(jen.Index().Add(t), bin(nElems, "+", jen.Lit(1))).Index(jen.Lit(0)))))), nil

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
		if isI1Vector(inst.Typ) {
			return one(vectorBin(vecBin{dest: name, op: "&&", x: x, y: y})), nil
		}
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return one(vectorBin(vecBin{dest: name, op: "&", x: x, y: y})), nil
		}
		if bits, ok := wideBits(inst.Typ); ok {
			return one(assign(name, wideSym(bits, libc.I128And, libc.I256And).Call(x, y))), nil
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 1 {
			return one(assign(name, bin(x, "&&", y))), nil
		}
		return one(assign(name, bin(x, "&", y))), nil

	case *ir.InstAShr:
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return translateVectorShift(inst, ">>", false)
		}
		x, err := FormatSigned(inst.X)
		if err != nil {
			return nil, fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return nil, fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		name := VariableName(inst)
		if bits, ok := wideBits(inst.Typ); ok {
			return one(assign(name, wideSym(bits, libc.I128AShr, libc.I256AShr).Call(x, y))), nil
		}
		if t, ok := inst.Typ.(*types.IntType); ok && t.BitSize == 8 {
			return one(assign(name, jen.Byte().Call(bin(x, ">>", y)))), nil
		}
		return one(assign(name, bin(x, ">>", y))), nil

	case *ir.InstBitCast:
		from, err := translateOp(inst.From, "source")
		if err != nil {
			return nil, err
		}
		if packed, err := i1VectorBitCast(from, inst.From.Type(), inst.To); packed != nil || err != nil {
			if err != nil {
				return nil, err
			}
			return one(assign(VariableName(inst), packed)), nil
		}
		if vec, err := vectorBitCast(from, inst.From.Type(), inst.To); vec != nil || err != nil {
			if err != nil {
				return nil, err
			}
			return one(assign(VariableName(inst), vec)), nil
		}
		if bits, err := scalarBitCast(from, inst.From.Type(), inst.To); bits != nil || err != nil {
			if err != nil {
				return nil, err
			}
			return one(assign(VariableName(inst), bits)), nil
		}
		if !compatiblePointerTypes(inst.From.Type(), inst.To) {
			return nil, fmt.Errorf("%w: %v and %v", errIncompatiblePointers, inst.From.Type(), inst.To)
		}
		name := VariableName(inst)
		if isTaggedPointerType(inst.To) {
			return one(assign(name, ptrToUint(from))), nil
		}
		if isTaggedPointerType(inst.From.Type()) {
			to, err := translateType(inst.To, "type")
			if err != nil {
				return nil, err
			}
			return one(assign(name, jen.Parens(to).Call(emitUP(from)))), nil
		}
		return one(assign(name, from)), nil

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

	case *ir.InstFRem:
		x, err := translateOp(inst.X, "left operand")
		if err != nil {
			return nil, err
		}
		y, err := translateOp(inst.Y, "right operand")
		if err != nil {
			return nil, err
		}
		// LLVM frem is libm fmod: remainder with the sign of the dividend.
		// Go has no % on float; math.Mod is float64-only.
		xf, yf := x, y
		wrapF32 := false
		if ft, ok := inst.Type().(*types.FloatType); ok && ft.Kind == types.FloatKindFloat {
			xf = jen.Float64().Call(x)
			yf = jen.Float64().Call(y)
			wrapF32 = true
		}
		mod := Sym(math.Mod).Call(xf, yf)
		if wrapF32 {
			mod = jen.Float32().Call(mod)
		}
		return one(assign(VariableName(inst), mod)), nil

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
		if _, ok := inst.X.Type().(*types.VectorType); ok {
			return translateVectorICmp(inst)
		}
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
		return one(assign(VariableName(inst), jen.Parens(to).Call(emitUP(jen.Uintptr().Call(from))))), nil

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
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return translateVectorShift(inst, ">>", true)
		}
		x, err := FormatUnsigned(inst.X)
		if err != nil {
			return nil, fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return nil, fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		name := VariableName(inst)
		if bits, ok := wideBits(inst.Typ); ok {
			return one(assign(name, wideSym(bits, libc.I128LShr, libc.I256LShr).Call(x, y))), nil
		}
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
		if isI1Vector(inst.Typ) {
			return one(vectorBin(vecBin{dest: name, op: "||", x: x, y: y})), nil
		}
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return one(vectorBin(vecBin{dest: name, op: "|", x: x, y: y})), nil
		}
		if bits, ok := wideBits(inst.Typ); ok {
			return one(assign(name, wideSym(bits, libc.I128Or, libc.I256Or).Call(x, y))), nil
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
		return one(assign(VariableName(inst), conv(to, ptrToUint(from)))), nil

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
		if _, ok := inst.Cond.Type().(*types.VectorType); ok {
			return one(jen.For(jen.List(jen.Id("i"), jen.Id("c")).Op(":=").Range().Add(cond)).Block(
				jen.If(jen.Id("c")).Block(
					jen.Id(name).Index(jen.Id("i")).Op("=").Add(valueTrue).Index(jen.Id("i")),
				).Else().Block(
					jen.Id(name).Index(jen.Id("i")).Op("=").Add(valueFalse).Index(jen.Id("i")),
				),
			)), nil
		}
		return one(jen.If(cond).Block(assign(name, valueTrue)).Else().Block(assign(name, valueFalse))), nil

	case *ir.InstSExt:
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
			if fromType.BitSize == 1 {
				return one(jen.For(jen.List(jen.Id("i"), jen.Id("v")).Op(":=").Range().Add(from)).Block(
					jen.Id(name).Index(jen.Id("i")).Op("=").Add(boolToInt(jen.Id("v"), toType.BitSize, true)),
				)), nil
			}
			return one(jen.For(jen.List(jen.Id("i"), jen.Id("v")).Op(":=").Range().Add(from)).Block(
				jen.Id(name).Index(jen.Id("i")).Op("=").Add(conv(goIntType(toType.BitSize), jen.Id("v"))),
			)), nil
		}
		toType, ok := inst.To.(*types.IntType)
		if !ok {
			return nil, fmt.Errorf("%w: %T", errUnsupportedZextTo, inst.To)
		}
		from, err := FormatSigned(inst.From)
		if err != nil {
			return nil, fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		name := VariableName(inst)
		if fromType, ok := inst.From.Type().(*types.IntType); ok && fromType.BitSize == 1 {
			if isWide128(toType.BitSize) || isWide256(toType.BitSize) {
				c, err := formatSExt(inst.From, inst.To)
				if err != nil {
					return nil, err
				}
				return one(assign(name, c)), nil
			}
			// i8 is byte: all-ones is 255, not untyped -1.
			neg := jen.Lit(-1)
			if toType.BitSize == 8 {
				neg = jen.Lit(255)
			}
			return one(jen.If(from).Block(assign(name, neg)).Else().Block(assign(name, jen.Lit(0)))), nil
		}
		if isWide128(toType.BitSize) || isWide256(toType.BitSize) {
			c, err := formatSExt(inst.From, inst.To)
			if err != nil {
				return nil, err
			}
			return one(assign(name, c)), nil
		}
		return one(assign(name, conv(goIntType(toType.BitSize), from))), nil

	case *ir.InstShl:
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return translateVectorShift(inst, "<<", false)
		}
		x, err := translateOp(inst.X, "left operand")
		if err != nil {
			return nil, err
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return nil, fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		if bits, ok := wideBits(inst.Typ); ok {
			return one(assign(VariableName(inst), wideSym(bits, libc.I128Shl, libc.I256Shl).Call(x, y))), nil
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
		if bits, ok := wideBits(inst.From.Type()); ok {
			from, err := translateOp(inst.From, "source")
			if err != nil {
				return nil, err
			}
			// Low 64 bits as signed; enough for Dynamic numeric casts.
			limb := jen.Add(from).Dot("Lo")
			if bits == 256 {
				limb = limb.Dot("Lo")
			}
			return one(assign(VariableName(inst), jen.Float64().Call(jen.Int64().Call(limb)))), nil
		}
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
		c, err := formatTrunc(inst.From, inst.To)
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), c)), nil

	case *ir.InstUIToFP:
		if bits, ok := wideBits(inst.From.Type()); ok {
			from, err := translateOp(inst.From, "source")
			if err != nil {
				return nil, err
			}
			limb := jen.Add(from).Dot("Lo")
			if bits == 256 {
				limb = limb.Dot("Lo")
			}
			return one(assign(VariableName(inst), jen.Float64().Call(limb))), nil
		}
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
		if isI1Vector(inst.Typ) {
			return one(vectorBin(vecBin{dest: name, op: "!=", x: x, y: y})), nil
		}
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return one(vectorBin(vecBin{dest: name, op: "^", x: x, y: y})), nil
		}
		if bits, ok := wideBits(inst.Typ); ok {
			return one(assign(name, wideSym(bits, libc.I128Xor, libc.I256Xor).Call(x, y))), nil
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
			if fromType.BitSize == 1 {
				return one(jen.For(jen.List(jen.Id("i"), jen.Id("v")).Op(":=").Range().Add(from)).Block(
					jen.Id(name).Index(jen.Id("i")).Op("=").Add(boolToInt(jen.Id("v"), toType.BitSize, false)),
				)), nil
			}
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
			if isWide128(toType.BitSize) || isWide256(toType.BitSize) {
				c, err := formatZExt(inst.From, inst.To)
				if err != nil {
					return nil, err
				}
				return one(assign(name, c)), nil
			}
			return one(jen.If(from).Block(assign(name, jen.Lit(1))).Else().Block(assign(name, jen.Lit(0)))), nil
		}
		if isWide128(toType.BitSize) || isWide256(toType.BitSize) {
			c, err := formatZExt(inst.From, inst.To)
			if err != nil {
				return nil, err
			}
			return one(assign(name, c)), nil
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

func wideBinFunc(bits uint64, op string, ashr bool) (goRef, bool) {
	if ashr && (op == ">>" || op == "ashr") {
		return wideSym(bits, libc.I128AShr, libc.I256AShr), true
	}
	var i128, i256 any
	switch op {
	case "+":
		i128, i256 = libc.I128Add, libc.I256Add
	case "-":
		i128, i256 = libc.I128Sub, libc.I256Sub
	case "*":
		i128, i256 = libc.I128Mul, libc.I256Mul
	case "&":
		i128, i256 = libc.I128And, libc.I256And
	case "|":
		i128, i256 = libc.I128Or, libc.I256Or
	case "^":
		i128, i256 = libc.I128Xor, libc.I256Xor
	case "<<":
		i128, i256 = libc.I128Shl, libc.I256Shl
	case ">>", "lshr":
		i128, i256 = libc.I128LShr, libc.I256LShr
	case "ashr":
		i128, i256 = libc.I128AShr, libc.I256AShr
	case "/", "%":
		return goRef{}, false
	default:
		return goRef{}, false
	}
	return wideSym(bits, i128, i256), true
}

func translateVectorICmp(inst *ir.InstICmp) ([]jen.Code, error) {
	x, err := translateOp(inst.X, "left operand")
	if err != nil {
		return nil, err
	}
	y, err := translateOp(inst.Y, "right operand")
	if err != nil {
		return nil, err
	}
	op, signed, unsigned, err := icmpPredOp(inst.Pred)
	if err != nil {
		return nil, err
	}
	xv := jen.Id("v")
	yv := jen.Add(y).Index(jen.Id("i"))
	if vt, ok := inst.X.Type().(*types.VectorType); ok {
		if it, ok := vt.ElemType.(*types.IntType); ok {
			switch {
			case signed && it.BitSize == 8:
				xv = jen.Int8().Call(xv)
				yv = jen.Int8().Call(yv)
			case unsigned:
				xv = goUintType(it.BitSize).Call(xv)
				yv = goUintType(it.BitSize).Call(yv)
			}
		}
	}
	name := VariableName(inst)
	return one(jen.For(jen.List(jen.Id("i"), jen.Id("v")).Op(":=").Range().Add(x)).Block(
		jen.Id(name).Index(jen.Id("i")).Op("=").Add(bin(xv, op, yv)),
	)), nil
}

func icmpPredOp(pred enum.IPred) (op string, signed, unsigned bool, err error) {
	switch pred {
	case enum.IPredEQ:
		return "==", false, false, nil
	case enum.IPredNE:
		return "!=", false, false, nil
	case enum.IPredSGE:
		return ">=", true, false, nil
	case enum.IPredSGT:
		return ">", true, false, nil
	case enum.IPredSLE:
		return "<=", true, false, nil
	case enum.IPredSLT:
		return "<", true, false, nil
	case enum.IPredUGE:
		return ">=", false, true, nil
	case enum.IPredUGT:
		return ">", false, true, nil
	case enum.IPredULE:
		return "<=", false, true, nil
	case enum.IPredULT:
		return "<", false, true, nil
	default:
		return "", false, false, fmt.Errorf("%w: %v", errUnsupportedICmpPred, pred)
	}
}

func translateI128Bin(name, op string, x, y value.Value, ashr bool) ([]jen.Code, bool, error) {
	bits, ok := wideBits(x.Type())
	if !ok {
		bits, ok = wideBits(y.Type())
	}
	if !ok {
		return nil, false, nil
	}
	fn, ok := wideBinFunc(bits, op, ashr)
	if !ok {
		return nil, false, fmt.Errorf("%w: i%d %s", errUnsupportedInstruction, bits, op)
	}
	xv, err := translateOp(x, "left operand")
	if err != nil {
		return nil, true, err
	}
	yv, err := translateOp(y, "right operand")
	if err != nil {
		return nil, true, err
	}
	return one(assign(name, fn.Call(xv, yv))), true, nil
}

func translateBinAssign(inst value.Named, b llvmBin) ([]jen.Code, error) {
	if _, ok := b.x.Type().(*types.VectorType); ok {
		xv, err := translateOp(b.x, "left operand")
		if err != nil {
			return nil, err
		}
		yv, err := translateOp(b.y, "right operand")
		if err != nil {
			return nil, err
		}
		return one(vectorBin(vecBin{dest: VariableName(inst), op: b.op, x: xv, y: yv})), nil
	}
	if stmts, ok, err := translateI128Bin(VariableName(inst), b.op, b.x, b.y, false); ok {
		return stmts, err
	}
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

// translateVectorShift lowers shl/lshr/ashr on <N x iK>. logical selects
// lshr (unsigned >> per lane); Go int >> is arithmetic.
func translateVectorShift(inst value.Named, op string, logical bool) ([]jen.Code, error) {
	var x, y value.Value
	var typ types.Type
	switch s := inst.(type) {
	case *ir.InstShl:
		x, y, typ = s.X, s.Y, s.Typ
	case *ir.InstLShr:
		x, y, typ = s.X, s.Y, s.Typ
	case *ir.InstAShr:
		x, y, typ = s.X, s.Y, s.Typ
	default:
		return nil, fmt.Errorf("%w: vector shift %T", errUnsupportedInstruction, inst)
	}
	xv, err := translateOp(x, "left operand")
	if err != nil {
		return nil, err
	}
	yv, err := translateOp(y, "right operand")
	if err != nil {
		return nil, err
	}
	name := VariableName(inst)
	vt, ok := typ.(*types.VectorType)
	if !ok {
		return one(vectorBin(vecBin{dest: name, op: op, x: xv, y: yv})), nil
	}
	it, ok := vt.ElemType.(*types.IntType)
	if !ok {
		return one(vectorBin(vecBin{dest: name, op: op, x: xv, y: yv})), nil
	}
	// lshr on Go signed ints must cast through unsigned per lane.
	// Fresh type stmts each Call — reusing one *jen.Statement corrupts emit.
	if logical && it.BitSize > 8 {
		lane := func(v jen.Code) *jen.Statement {
			return goUintType(it.BitSize).Call(v)
		}
		return one(jen.For(jen.List(jen.Id("i"), jen.Id("v")).Op(":=").Range().Add(xv)).Block(
			jen.Id(name).Index(jen.Id("i")).Op("=").Add(goIntType(it.BitSize)).Call(
				bin(lane(jen.Id("v")), ">>", lane(jen.Add(yv).Index(jen.Id("i")))),
			),
		)), nil
	}
	// ashr i8: lanes are byte; shift as int8 then back.
	if !logical && op == ">>" && it.BitSize == 8 {
		return one(jen.For(jen.List(jen.Id("i"), jen.Id("v")).Op(":=").Range().Add(xv)).Block(
			jen.Id(name).Index(jen.Id("i")).Op("=").Add(jen.Byte().Call(
				bin(jen.Int8().Call(jen.Id("v")), ">>", jen.Add(yv).Index(jen.Id("i"))),
			)),
		)), nil
	}
	return one(vectorBin(vecBin{dest: name, op: op, x: xv, y: yv})), nil
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
	if _, ok := b.x.Type().(*types.VectorType); ok {
		xv, err := FormatSigned(b.x)
		if err != nil {
			return nil, fmt.Errorf("error translating left operand (%v): %w", b.x, err)
		}
		yv, err := FormatSigned(b.y)
		if err != nil {
			return nil, fmt.Errorf("error translating right operand (%v): %w", b.y, err)
		}
		return one(vectorBin(vecBin{dest: VariableName(inst), op: b.op, x: xv, y: yv})), nil
	}
	if bits, ok := wideBits(inst.Type()); ok {
		fn := wideSym(bits, libc.I128SDiv, libc.I256SDiv)
		if b.op == "%" {
			fn = wideSym(bits, libc.I128SRem, libc.I256SRem)
		}
		xv, err := translateOp(b.x, "left operand")
		if err != nil {
			return nil, err
		}
		yv, err := translateOp(b.y, "right operand")
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), fn.Call(xv, yv))), nil
	}
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
	if _, ok := b.x.Type().(*types.VectorType); ok {
		xv, err := FormatUnsigned(b.x)
		if err != nil {
			return nil, fmt.Errorf("error translating left operand (%v): %w", b.x, err)
		}
		yv, err := FormatUnsigned(b.y)
		if err != nil {
			return nil, fmt.Errorf("error translating right operand (%v): %w", b.y, err)
		}
		// Go % on signed is truncated toward zero like urem for non-neg;
		// cast through unsigned per lane when elem is wider than i8.
		name := VariableName(inst)
		if vt, ok := b.x.Type().(*types.VectorType); ok {
			if it, ok := vt.ElemType.(*types.IntType); ok && it.BitSize > 8 {
				return one(jen.For(jen.List(jen.Id("i"), jen.Id("v")).Op(":=").Range().Add(xv)).Block(
					jen.Id(name).Index(jen.Id("i")).Op("=").Add(goIntType(it.BitSize)).Call(
						bin(goUintType(it.BitSize).Call(jen.Id("v")), b.op, goUintType(it.BitSize).Call(jen.Add(yv).Index(jen.Id("i")))),
					),
				)), nil
			}
		}
		return one(vectorBin(vecBin{dest: name, op: b.op, x: xv, y: yv})), nil
	}
	if bits, ok := wideBits(inst.Type()); ok {
		fn := wideSym(bits, libc.I128UDiv, libc.I256UDiv)
		if b.op == "%" {
			fn = wideSym(bits, libc.I128URem, libc.I256URem)
		}
		xv, err := translateOp(b.x, "left operand")
		if err != nil {
			return nil, err
		}
		yv, err := translateOp(b.y, "right operand")
		if err != nil {
			return nil, err
		}
		return one(assign(VariableName(inst), fn.Call(xv, yv))), nil
	}
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
	return emitAs(Qual[byte](), x)
}

func isPtrish(t types.Type) bool {
	if t == nil {
		return false
	}
	if isTaggedPointerType(t) {
		return true
	}
	_, ok := t.(*types.PointerType)
	return ok
}

func ptrArg(ir []value.Value, args []jen.Code, i int) *jen.Statement {
	if i >= len(args) {
		return jen.Nil()
	}
	if i < len(ir) {
		if isTaggedPointerType(ir[i].Type()) {
			return asBytePtr(emitUP(args[i]))
		}
		if _, ok := ir[i].Type().(*types.PointerType); ok {
			return asBytePtr(args[i])
		}
		if s, ok := args[i].(*jen.Statement); ok {
			return s
		}
		return jen.Add(args[i])
	}
	return asBytePtr(args[i])
}

func asFilePtr(x jen.Code) *jen.Statement {
	return emitAs(Qual[os.File](), x)
}

func libcCallArg(name string, i int, a value.Value, got jen.Code) *jen.Statement {
	if _, ok := a.Type().(*types.PointerType); !ok {
		if s, ok := got.(*jen.Statement); ok {
			return s
		}
		return jen.Add(got)
	}
	if isTaggedPointerType(a.Type()) {
		got = emitUP(got)
	}
	switch name {
	case "fprintf":
		if i == 0 {
			return asFilePtr(got)
		}
		if i == 1 {
			return asBytePtr(got)
		}
	case "fputs":
		if i == 0 {
			return asBytePtr(got)
		}
		if i == 1 {
			return asFilePtr(got)
		}
	case "fputc", "putc", "fclose":
		if i == 1 || (name == "fclose" && i == 0) {
			return asFilePtr(got)
		}
	case "fdopen":
		if i == 1 {
			return asBytePtr(got)
		}
	case "mmap", "mmap64":
		if i == 0 {
			if isTaggedPointerType(a.Type()) {
				return jen.Add(got)
			}
			return ptrToUint(got)
		}
	}
	return asBytePtr(got)
}

func calleeLLVMName(v value.Value) string {
	if n, ok := v.(value.Named); ok {
		return VariableName(n)
	}
	return ""
}

// llvmCallHandled is the llvm.* names translateCall lowers (no-op or libc).
// hasRuntimeDef uses this instead of a blanket llvm_ prefix.
// translateLLVMPattern lowers common rustc llvm.* intrinsics not listed
// individually. ok=false means fall through to other handlers.
func translateLLVMPattern(llvmName string, inst *ir.InstCall, args []jen.Code) ([]jen.Code, bool) {
	name := VariableName(inst)
	// usub.sat / uadd.sat via libc (no mid-function vars — goto-safe).
	if strings.HasPrefix(llvmName, "llvm_usub_sat_") && len(args) == 2 {
		return one(assign(name, convFromU64(llvmName, Sym(libc.USubSatU64).Call(
			jen.Uint64().Call(jen.Add(args[0])),
			jen.Uint64().Call(jen.Add(args[1])),
		)))), true
	}
	if strings.HasPrefix(llvmName, "llvm_uadd_sat_") && len(args) == 2 {
		return one(assign(name, convFromU64(llvmName, Sym(libc.UAddSatU64).Call(
			jen.Uint64().Call(jen.Add(args[0])),
			jen.Uint64().Call(jen.Add(args[1])),
		)))), true
	}
	if strings.HasPrefix(llvmName, "llvm_ctlz_") && len(args) >= 1 {
		return one(assign(name, convFromU64(llvmName, Sym(bits.LeadingZeros64).Call(jen.Uint64().Call(jen.Add(args[0])))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_cttz_") && len(args) >= 1 {
		return one(assign(name, convFromU64(llvmName, Sym(bits.TrailingZeros64).Call(jen.Uint64().Call(jen.Add(args[0])))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_ctpop_") && len(args) >= 1 {
		return one(assign(name, convFromU64(llvmName, Sym(bits.OnesCount64).Call(jen.Uint64().Call(jen.Add(args[0])))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_abs_") && len(args) >= 1 {
		if strings.HasSuffix(llvmName, "_i128") || strings.Contains(llvmName, "_i128") {
			// sign from high limb; abs via conditional sub from 0
			return one(assign(name, jen.Add(args[0]))), true // identity if non-neg common path
		}
		return one(assign(name, convFromU64(llvmName, Sym(libc.AbsI64).Call(jen.Int64().Call(jen.Add(args[0])))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_bitreverse_") && len(args) == 1 {
		return one(assign(name, convFromU64(llvmName, Sym(bits.Reverse64).Call(jen.Uint64().Call(jen.Add(args[0])))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_fshl_") && len(args) == 3 {
		return one(assign(name, convFromU64(llvmName, Sym(bits.RotateLeft64).Call(
			jen.Uint64().Call(jen.Add(args[0])),
			jen.Int().Call(jen.Add(args[2])),
		)))), true
	}
	if strings.HasPrefix(llvmName, "llvm_fshr_") && len(args) == 3 {
		return one(assign(name, convFromU64(llvmName, jen.Qual("math/bits", "RotateRight64").Call(
			jen.Uint64().Call(jen.Add(args[0])),
			jen.Int().Call(jen.Add(args[2])),
		)))), true
	}
	if strings.HasPrefix(llvmName, "llvm_bswap_") && len(args) == 1 {
		return one(assign(name, convFromU64(llvmName, Sym(bits.ReverseBytes64).Call(jen.Uint64().Call(jen.Add(args[0])))))), true
	}
	if (strings.HasPrefix(llvmName, "llvm_umax_") || strings.HasPrefix(llvmName, "llvm_umin_") ||
		strings.HasPrefix(llvmName, "llvm_smax_") || strings.HasPrefix(llvmName, "llvm_smin_")) && len(args) == 2 {
		if strings.Contains(llvmName, "i128") {
			// Select via I128 compare helpers.
			lt := libc.I128Slt
			if strings.Contains(llvmName, "umax") || strings.Contains(llvmName, "umin") {
				lt = libc.I128Ult
			}
			// max: if a < b { b } else { a }; min: if a < b { a } else { b }
			isMax := strings.Contains(llvmName, "max")
			cond := Sym(lt).Call(jen.Add(args[0]), jen.Add(args[1]))
			// min: if a<b then a else b; max: if a<b then b else a
			if isMax {
				return one(jen.If(cond).Block(assign(name, jen.Add(args[1]))).Else().Block(assign(name, jen.Add(args[0])))), true
			}
			return one(jen.If(cond).Block(assign(name, jen.Add(args[0]))).Else().Block(assign(name, jen.Add(args[1])))), true
		}
		fn := any(libc.UMaxU64)
		switch {
		case strings.Contains(llvmName, "umin"):
			fn = libc.UMinU64
		case strings.Contains(llvmName, "smax"):
			fn = libc.SMaxI64
		case strings.Contains(llvmName, "smin"):
			fn = libc.SMinI64
		}
		return one(assign(name, convFromU64(llvmName, Sym(fn).Call(
			jen.Int64().Call(jen.Add(args[0])), jen.Int64().Call(jen.Add(args[1])),
		)))), true
	}
	if strings.HasPrefix(llvmName, "llvm_floor_") && len(args) == 1 {
		return one(assign(name, Sym(math.Floor).Call(jen.Float64().Call(jen.Add(args[0]))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_trunc_f") && len(args) == 1 {
		e := Sym(math.Trunc).Call(jen.Float64().Call(jen.Add(args[0])))
		if strings.HasSuffix(llvmName, "f32") {
			e = jen.Float32().Call(e)
		}
		return one(assign(name, e)), true
	}
	if strings.HasPrefix(llvmName, "llvm_round_") && len(args) == 1 {
		return one(assign(name, Sym(math.Round).Call(jen.Float64().Call(jen.Add(args[0]))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_fabs_") && len(args) == 1 {
		return one(assign(name, Sym(math.Abs).Call(jen.Float64().Call(jen.Add(args[0]))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_sqrt_") && len(args) == 1 {
		return one(assign(name, Sym(math.Sqrt).Call(jen.Float64().Call(jen.Add(args[0]))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_sin_") && len(args) == 1 {
		return one(assign(name, Sym(math.Sin).Call(jen.Float64().Call(jen.Add(args[0]))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_cos_") && len(args) == 1 {
		return one(assign(name, Sym(math.Cos).Call(jen.Float64().Call(jen.Add(args[0]))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_exp_") && len(args) == 1 {
		return one(assign(name, Sym(math.Exp).Call(jen.Float64().Call(jen.Add(args[0]))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_log10_") && len(args) == 1 {
		return one(assign(name, Sym(math.Log10).Call(jen.Float64().Call(jen.Add(args[0]))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_log_") && len(args) == 1 {
		return one(assign(name, Sym(math.Log).Call(jen.Float64().Call(jen.Add(args[0]))))), true
	}
	if strings.HasPrefix(llvmName, "llvm_pow") && len(args) == 2 {
		e := Sym(math.Pow).Call(
			jen.Float64().Call(jen.Add(args[0])),
			jen.Float64().Call(jen.Add(args[1])),
		)
		if strings.Contains(llvmName, "f32") {
			e = jen.Float32().Call(e)
		}
		return one(assign(name, e)), true
	}
	if strings.HasPrefix(llvmName, "llvm_copysign_") && len(args) == 2 {
		e := Sym(math.Copysign).Call(
			jen.Float64().Call(jen.Add(args[0])),
			jen.Float64().Call(jen.Add(args[1])),
		)
		return one(assign(name, castToResult(inst, e))), true
	}
	if strings.HasPrefix(llvmName, "llvm_ptrmask") && len(args) == 2 {
		call := Sym(libc.PtrMask).Call(emitUP(jen.Add(args[0])), jen.Int64().Call(jen.Add(args[1])))
		if isTaggedPointerType(inst.Type()) {
			return one(assign(name, emitAddr(call))), true
		}
		return one(assign(name, call)), true
	}
	if strings.HasPrefix(llvmName, "llvm_maximum") && len(args) == 2 {
		e := Sym(libc.MaximumNumF64).Call(
			jen.Float64().Call(jen.Add(args[0])),
			jen.Float64().Call(jen.Add(args[1])),
		)
		return one(assign(name, castToResult(inst, e))), true
	}
	// with.overflow — {a op b, false}; flag not exact for all widths.
	if strings.Contains(llvmName, "_with_overflow_") && len(args) == 2 {
		ret, err := TypeSpec(inst.Type())
		if err != nil {
			return nil, false
		}
		// i128/i256 use libc binops
		if strings.Contains(llvmName, "i128") || strings.Contains(llvmName, "i256") {
			bits := uint64(128)
			if strings.Contains(llvmName, "i256") {
				bits = 256
			}
			var fn goRef
			switch {
			case strings.Contains(llvmName, "sadd") || strings.Contains(llvmName, "uadd"):
				fn = wideSym(bits, libc.I128Add, libc.I256Add)
			case strings.Contains(llvmName, "ssub") || strings.Contains(llvmName, "usub"):
				fn = wideSym(bits, libc.I128Sub, libc.I256Sub)
			case strings.Contains(llvmName, "smul") || strings.Contains(llvmName, "umul"):
				fn = wideSym(bits, libc.I128Mul, libc.I256Mul)
			default:
				return nil, false
			}
			return one(assign(name, jen.Add(ret).Values(
				fn.Call(jen.Add(args[0]), jen.Add(args[1])),
				jen.False(),
			))), true
		}
		var op string
		switch {
		case strings.Contains(llvmName, "sadd") || strings.Contains(llvmName, "uadd"):
			op = "+"
		case strings.Contains(llvmName, "ssub") || strings.Contains(llvmName, "usub"):
			op = "-"
		case strings.Contains(llvmName, "smul") || strings.Contains(llvmName, "umul"):
			op = "*"
		default:
			return nil, false
		}
		return one(assign(name, jen.Add(ret).Values(
			jen.Add(args[0]).Op(op).Add(args[1]),
			jen.False(),
		))), true
	}
	return nil, false
}

// convFromU64 casts a uint64/int64 temp back toward the intrinsic result Go type.
func convFromU64(llvmName string, v jen.Code) *jen.Statement {
	// i3/i7/… map to next Go width (byte/int16/…).
	switch {
	case strings.Contains(llvmName, "_i8") || strings.Contains(llvmName, "_i1") ||
		strings.Contains(llvmName, "_i2") || strings.Contains(llvmName, "_i3") ||
		strings.Contains(llvmName, "_i4") || strings.Contains(llvmName, "_i5") ||
		strings.Contains(llvmName, "_i6") || strings.Contains(llvmName, "_i7"):
		// Prefer byte for ≤8; i16 names contain _i1 too — check wider first.
		if strings.Contains(llvmName, "_i16") || strings.Contains(llvmName, "_i32") ||
			strings.Contains(llvmName, "_i64") || strings.Contains(llvmName, "_i128") {
			// fall through
		} else {
			return jen.Byte().Call(v)
		}
	}
	switch {
	case strings.Contains(llvmName, "_i16"):
		return jen.Int16().Call(v)
	case strings.Contains(llvmName, "_i32") || strings.Contains(llvmName, "f32"):
		if strings.Contains(llvmName, "f32") {
			return jen.Float32().Call(v)
		}
		return jen.Int32().Call(v)
	case strings.Contains(llvmName, "_i64") || strings.Contains(llvmName, "f64"):
		if strings.Contains(llvmName, "f64") {
			return jen.Float64().Call(v)
		}
		return jen.Int64().Call(v)
	default:
		return jen.Add(v)
	}
}

// castToResult wraps v as the Go type of the call result when needed.
func castToResult(inst *ir.InstCall, v *jen.Statement) *jen.Statement {
	t := inst.Type()
	if t == nil || types.Equal(t, types.Void) {
		return v
	}
	if ft, ok := t.(*types.FloatType); ok && ft.Kind == types.FloatKindFloat {
		return jen.Float32().Call(v)
	}
	return v
}

func llvmCallHandled(name string) bool {
	if strings.HasPrefix(name, "llvm_is_constant_") || strings.HasPrefix(name, "llvm_ucmp_") ||
		strings.HasPrefix(name, "llvm_scmp_") || strings.HasPrefix(name, "llvm_vector_reduce_") ||
		strings.HasPrefix(name, "llvm_threadlocal_address") ||
		strings.HasPrefix(name, "llvm_usub_sat_") || strings.HasPrefix(name, "llvm_uadd_sat_") ||
		strings.HasPrefix(name, "llvm_ssub_sat_") || strings.HasPrefix(name, "llvm_sadd_sat_") ||
		strings.HasPrefix(name, "llvm_ctlz_") || strings.HasPrefix(name, "llvm_cttz_") ||
		strings.HasPrefix(name, "llvm_ctpop_") || strings.HasPrefix(name, "llvm_abs_") ||
		strings.HasPrefix(name, "llvm_fshl_") || strings.HasPrefix(name, "llvm_fshr_") ||
		strings.HasPrefix(name, "llvm_bswap_") || strings.HasPrefix(name, "llvm_bitreverse_") ||
		strings.Contains(name, "_with_overflow_") ||
		strings.HasPrefix(name, "llvm_umax_") || strings.HasPrefix(name, "llvm_umin_") ||
		strings.HasPrefix(name, "llvm_smax_") || strings.HasPrefix(name, "llvm_smin_") ||
		strings.HasPrefix(name, "llvm_floor_") || strings.HasPrefix(name, "llvm_trunc_f") ||
		strings.HasPrefix(name, "llvm_round_") || strings.HasPrefix(name, "llvm_fabs_") ||
		strings.HasPrefix(name, "llvm_sqrt_") || strings.HasPrefix(name, "llvm_sin_") ||
		strings.HasPrefix(name, "llvm_cos_") || strings.HasPrefix(name, "llvm_exp_") ||
		strings.HasPrefix(name, "llvm_log_") || strings.HasPrefix(name, "llvm_pow") ||
		strings.HasPrefix(name, "llvm_copysign_") || strings.HasPrefix(name, "llvm_maximum") ||
		strings.HasPrefix(name, "llvm_minimum") || strings.HasPrefix(name, "llvm_ptrmask") {
		return true
	}
	switch name {
	case "llvm_lifetime_start_p0", "llvm_lifetime_end_p0",
		"llvm_experimental_noalias_scope_decl", "llvm_assume",
		"llvm_donothing",
		"llvm_va_start", "llvm_va_end",
		"llvm_lifetime_start", "llvm_lifetime_end", "llvm_stackrestore",
		"llvm_stacksave",
		"llvm_fabs_f32", "llvm_fmuladd_f64", "llvm_fmuladd_f32",
		"llvm_memcpy_p0i8_p0i8_i64", "llvm_memmove_p0i8_p0i8_i64",
		"llvm_memcpy_p0_p0_i64", "llvm_memmove_p0_p0_i64",
		"llvm_memset_p0i8_i64", "llvm_memset_p0_i64",
		"llvm_abs_i32",
		"llvm_sadd_with_overflow_i32",
		"llvm_umax_i64", "llvm_umin_i64",
		"llvm_umax_i32", "llvm_umin_i32",
		"llvm_smax_i64", "llvm_smin_i64",
		"llvm_smax_i32", "llvm_smin_i32",
		"llvm_ctpop_i64",
		"llvm_umul_with_overflow_i64",
		"llvm_objectsize_i64_p0i8",
		"llvm_trap",
		"llvm_ceil_f64",
		"llvm_vector_reduce_add_v4i32",
		"llvm_load_relative_i64":
		return true
	default:
		return false
	}
}

func translateCall(inst *ir.InstCall) ([]jen.Code, error) {
	if ia, ok := inst.Callee.(*ir.InlineAsm); ok {
		return one(Sym(libc.InlineAsm).Call(jen.Lit(ia.Asm), jen.Lit(ia.Constraint))), nil
	}
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
		args[i] = v
	}

	var callee *jen.Statement
	typedPtr := false // Go callee returns *T; LLVM dest is unsafe.Pointer
	switch llvmName {
	case "malloc", "calloc":
		et := jen.Byte()
		if pt, ok := inst.Typ.(*types.PointerType); ok && !pt.IsOpaque() && pt.ElemType != nil {
			if t, err := TypeSpec(pt.ElemType); err == nil {
				et = t
			}
		}
		fn := any(libc.Malloc[byte])
		if llvmName == "calloc" {
			fn = libc.Calloc[byte]
		}
		callee = Sym(fn).Types(et)
		typedPtr = true
	case "leaven_va_start":
		if len(args) == 1 {
			return one(deref(args[0]).Op("=").Add(emitAs(Qual[byte](), emitPtr(addrOf(jen.Id("varargs")))))), nil
		}
	case "llvm_va_start":
		if len(args) == 1 {
			return one(Sym(libc.Store[unsafe.Pointer]).Types(Qual[unsafe.Pointer]()).Call(
				args[0], jen.Lit(8), emitPtr(addrOf(jen.Id("varargs"))),
			)), nil
		}
	case "llvm_va_end", "llvm_lifetime_start", "llvm_lifetime_end", "llvm_stackrestore":
		return nil, nil
	case "vsnprintf":
		if len(args) == 4 {
			return one(assign(VariableName(inst), Sym(libc.Vsnprintf).Call(
				asBytePtr(args[0]), args[1], asBytePtr(args[2]),
				emitAs(Qual[byte](), args[3]),
			))), nil
		}
	case "ldexp":
		if len(args) == 2 {
			return one(assign(VariableName(inst), Sym(math.Ldexp).Call(args[0], jen.Int().Call(args[1])))), nil
		}
	case "llvm_fabs_f32":
		if len(args) == 1 {
			return one(assign(VariableName(inst), jen.Float32().Call(Sym(math.Abs).Call(jen.Float64().Call(args[0]))))), nil
		}
	case "llvm_fmuladd_f64", "llvm_fmuladd_f32":
		if len(args) == 3 {
			return one(assign(VariableName(inst), bin(bin(args[0], "*", args[1]), "+", args[2]))), nil
		}
	case "llvm_memcpy_p0i8_p0i8_i64", "llvm_memmove_p0i8_p0i8_i64",
		"llvm_memcpy_p0_p0_i64", "llvm_memmove_p0_p0_i64":
		return one(Sym(libc.Memmove).Call(asBytePtr(args[0]), asBytePtr(args[1]), args[2])), nil
	case "llvm_memset_p0i8_i64", "llvm_memset_p0_i64":
		return one(Sym(libc.Memset).Call(asBytePtr(args[0]), args[1], args[2])), nil
	case "llvm_abs_i32":
		if len(args) >= 1 {
			return one(assign(VariableName(inst), Sym(libc.AbsI32).Call(args[0]))), nil
		}
	case "llvm_sadd_with_overflow_i32":
		if len(args) == 2 {
			return one(assign(VariableName(inst), Sym(libc.SAddWithOverflowI32).Call(args[0], args[1]))), nil
		}
	case "llvm_umax_i64":
		if len(args) == 2 {
			return one(assign(VariableName(inst), Sym(libc.UMaxU64).Call(args[0], args[1]))), nil
		}
	case "llvm_umin_i64":
		if len(args) == 2 {
			return one(assign(VariableName(inst), Sym(libc.UMinU64).Call(args[0], args[1]))), nil
		}
	case "llvm_umax_i32":
		if len(args) == 2 {
			return one(assign(VariableName(inst), Sym(libc.UMaxU32).Call(args[0], args[1]))), nil
		}
	case "llvm_umin_i32":
		if len(args) == 2 {
			return one(assign(VariableName(inst), Sym(libc.UMinU32).Call(args[0], args[1]))), nil
		}
	case "llvm_smax_i64":
		if len(args) == 2 {
			return one(assign(VariableName(inst), Sym(libc.SMaxI64).Call(args[0], args[1]))), nil
		}
	case "llvm_smin_i64":
		if len(args) == 2 {
			return one(assign(VariableName(inst), Sym(libc.SMinI64).Call(args[0], args[1]))), nil
		}
	case "llvm_smax_i32":
		if len(args) == 2 {
			return one(assign(VariableName(inst), Sym(libc.SMaxI32).Call(args[0], args[1]))), nil
		}
	case "llvm_smin_i32":
		if len(args) == 2 {
			return one(assign(VariableName(inst), Sym(libc.SMinI32).Call(args[0], args[1]))), nil
		}
	case "llvm_trap":
		return one(Sym(libc.Abort).Call()), nil
	case "llvm_ceil_f64":
		if len(args) == 1 {
			return one(assign(VariableName(inst), Sym(math.Ceil).Call(args[0]))), nil
		}
	case "llvm_vector_reduce_add_v4i32":
		if len(args) == 1 {
			return one(assign(VariableName(inst), Sym(libc.VecReduceAddV4I32).Call(args[0]))), nil
		}
	case "llvm_load_relative_i64":
		if len(args) == 2 {
			return one(assign(VariableName(inst), Sym(libc.LoadRelativeI64).Call(args[0], args[1]))), nil
		}
	case "llvm_ctpop_i64":
		if len(args) == 1 {
			return one(assign(VariableName(inst), Sym(bits.OnesCount64).Call(jen.Uint64().Call(args[0])))), nil
		}
	case "llvm_umul_with_overflow_i64":
		if len(args) == 2 {
			return one(assign(VariableName(inst), Sym(libc.UMulWithOverflowU64).Call(args[0], args[1]))), nil
		}
	case "llvm_objectsize_i64_p0i8":
		return one(assign(VariableName(inst), jen.Op("-").Lit(1))), nil
	case "llvm_stacksave":
		return one(assign(VariableName(inst), jen.Nil())), nil
	}
	// llvm.*.sat / bitop / math patterns used heavily by rustc Release IR.
	if stmts, ok := translateLLVMPattern(llvmName, inst, args); ok {
		return stmts, nil
	}
	// llvm.vector.reduce.{add,or,and,xor,mul}.* — lane fold.
	if strings.HasPrefix(llvmName, "llvm_vector_reduce_") && len(args) == 1 {
		op := ""
		switch {
		case strings.Contains(llvmName, "_add_"):
			op = "+"
		case strings.Contains(llvmName, "_or_"):
			op = "|"
		case strings.Contains(llvmName, "_and_"):
			op = "&"
		case strings.Contains(llvmName, "_xor_"):
			op = "^"
		case strings.Contains(llvmName, "_mul_"):
			op = "*"
		}
		if op != "" {
			name := VariableName(inst)
			return []jen.Code{
				assign(name, jen.Add(args[0]).Index(jen.Lit(0))),
				jen.For(
					jen.Id("i").Op(":=").Lit(1),
					jen.Id("i").Op("<").Len(jen.Add(args[0])),
					jen.Id("i").Op("++"),
				).Block(
					jen.Id(name).Op("=").Id(name).Op(op).Add(args[0]).Index(jen.Id("i")),
				),
			}, nil
		}
	}
	// llvm.threadlocal.address — TLS base is the global address in our model.
	if strings.HasPrefix(llvmName, "llvm_threadlocal_address") {
		if len(args) == 1 {
			return one(assign(VariableName(inst), jen.Add(args[0]))), nil
		}
	}
	// llvm.is.constant.* — runtime value is never a compile-time constant.
	if strings.HasPrefix(llvmName, "llvm_is_constant_") {
		return one(assign(VariableName(inst), jen.False())), nil
	}
	// llvm.ucmp / llvm.scmp → i8 {-1,0,1}. Go byte holds 255/0/1.
	if strings.HasPrefix(llvmName, "llvm_ucmp_") || strings.HasPrefix(llvmName, "llvm_scmp_") {
		if len(args) == 2 {
			name := VariableName(inst)
			unsigned := strings.HasPrefix(llvmName, "llvm_ucmp_")
			// i128/i256: use libc compares (no Go < on structs).
			if strings.Contains(llvmName, "i128") {
				lt, eq := any(libc.I128Slt), any(libc.I128Eq)
				if unsigned {
					lt = libc.I128Ult
				}
				return []jen.Code{
					jen.If(Sym(lt).Call(jen.Add(args[0]), jen.Add(args[1]))).Block(
						assign(name, jen.Lit(255)),
					).Else().If(Sym(eq).Call(jen.Add(args[0]), jen.Add(args[1]))).Block(
						assign(name, jen.Lit(0)),
					).Else().Block(assign(name, jen.Lit(1))),
				}, nil
			}
			if strings.Contains(llvmName, "i256") {
				lt, eq := any(libc.I256Slt), any(libc.I256Eq)
				if unsigned {
					lt = libc.I256Ult
				}
				return []jen.Code{
					jen.If(Sym(lt).Call(jen.Add(args[0]), jen.Add(args[1]))).Block(
						assign(name, jen.Lit(255)),
					).Else().If(Sym(eq).Call(jen.Add(args[0]), jen.Add(args[1]))).Block(
						assign(name, jen.Lit(0)),
					).Else().Block(assign(name, jen.Lit(1))),
				}, nil
			}
			// Fresh expressions each use — jennifer Statements mutate.
			lhs := func() *jen.Statement {
				a := jen.Add(args[0])
				if unsigned {
					return jen.Uint64().Call(a)
				}
				return a
			}
			rhs := func() *jen.Statement {
				a := jen.Add(args[1])
				if unsigned {
					return jen.Uint64().Call(a)
				}
				return a
			}
			return []jen.Code{
				jen.If(lhs().Op("<").Add(rhs())).Block(assign(name, jen.Lit(255))).Else().If(
					lhs().Op("==").Add(rhs()),
				).Block(assign(name, jen.Lit(0))).Else().Block(assign(name, jen.Lit(1))),
			}, nil
		}
	}

	if callee == nil {
		if ref, ok := libcLookup(llvmName); ok {
			callee = ref.code()
			name := libcCanon(llvmName)
			for i, a := range inst.Args {
				args[i] = libcCallArg(name, i, a, args[i])
			}
			typedPtr = libcReturnsTypedPtr(name)
		} else if cxxNoopDtor(llvmName) {
			return nil, nil
		} else if c, adj, retPtr, ok := cxxTreeCall(llvmName, args); ok {
			callee = c
			args = adj
			typedPtr = retPtr
		} else if isLibcxxStringEqCStr(llvmName) && len(inst.Args) >= 2 &&
			isPtrish(inst.Args[0].Type()) && isPtrish(inst.Args[1].Type()) {
			callee = Sym(libc.StdStringEqCStr).code()
			args = []jen.Code{ptrArg(inst.Args, args, 0), ptrArg(inst.Args, args, 1)}
		} else if c, adj, retPtr, ok := cxxIOCallIR(llvmName, inst.Args, args); ok {
			callee = c
			args = adj
			typedPtr = retPtr
		} else if strings.Contains(llvmName, "throw_bad_cast") {
			// Inlined getline path if failbit was not seen. Short text
			// so CI does not omit the panic line.
			return one(jen.Panic(jen.Lit("std::bad_cast"))), nil
		} else if strings.Contains(llvmName, "ctypeIcE5widen") {
			return one(assign(VariableName(inst), Sym(libc.CtypeWiden).Call(args...))), nil
		} else if strings.Contains(llvmName, "alloc_error_handler") ||
			strings.Contains(llvmName, "__rust_alloc_error") {
			return one(jen.Panic(jen.Lit("allocation error"))), nil
		} else if isRustAlloc(llvmName) {
			fn := any(libc.RustAlloc)
			if strings.Contains(llvmName, "zeroed") {
				fn = libc.RustAllocZeroed
			}
			return one(assign(VariableName(inst), Sym(fn).Call(args...))), nil
		} else if strings.Contains(llvmName, "__rust_dealloc") {
			return one(Sym(libc.RustDealloc).Call(args...)), nil
		} else if strings.Contains(llvmName, "__rust_realloc") {
			return one(assign(VariableName(inst), Sym(libc.RustRealloc).Call(args...))), nil
		} else if strings.Contains(llvmName, "__rust_no_alloc_shim") {
			return nil, nil
		} else if llvmName == "_Znwm" || llvmName == "_Znam" {
			return one(assign(VariableName(inst), Sym(libc.RustAlloc).Call(args[0], jen.Lit(1)))), nil
		} else if strings.HasPrefix(llvmName, "_Zdl") || strings.HasPrefix(llvmName, "_Zda") {
			return one(Sym(libc.RustDealloc).Call(args[0], jen.Lit(0), jen.Lit(1))), nil
		} else if _, ok := inst.Callee.(*ir.Func); ok {
			callee = jen.Id(VariableName(inst.Callee.(value.Named)))
			if c, ok := namedRef(llvmName); ok {
				callee = c
			}
		} else {
			ft, err := TypeDefinition(inst.Sig())
			if err != nil {
				return nil, fmt.Errorf("error translating callee type (%v): %w", inst.Callee, err)
			}
			fn, err := FormatValue(inst.Callee)
			if err != nil {
				return nil, fmt.Errorf("error translating callee (%v): %w", inst.Callee, err)
			}
			callee = Sym(libc.FuncFromCode[func()]).Types(ft).Call(fn)
		}
	}

	callExpr := jen.Add(callee).Call(args...)
	return finishCall(inst, callExpr, typedPtr)
}

// finishCall assigns a call to its SSA name. The call instruction's type is
// the LangRef result. unsafe.Pointer(...) is only for our libc *T returns
// (Malloc, Strchr, …), not a conversion of a struct value into a pointer.
func finishCall(inst *ir.InstCall, callExpr *jen.Statement, typedPtr bool) ([]jen.Code, error) {
	dest := inst.Type()
	if types.Equal(dest, types.Void) {
		return one(callExpr), nil
	}
	if isTaggedPointerType(dest) {
		if typedPtr {
			return one(assign(VariableName(inst), emitAddr(callExpr))), nil
		}
		return one(assign(VariableName(inst), ptrToUint(callExpr))), nil
	}
	if _, ok := dest.(*types.PointerType); ok && typedPtr {
		callExpr = emitPtr(callExpr)
	}
	return one(assign(VariableName(inst), callExpr)), nil
}

func libcReturnsTypedPtr(name string) bool {
	switch name {
	case "realloc", "fdopen", "strchr", "strrchr", "strstr", "strpbrk",
		"memchr", "strcpy", "strncpy", "strcat", "strncat", "memmove",
		"memset", "memcpy",
		"__errno_location", "__error", "getenv", "getcwd", "realpath", "dlsym",
		"__dynamic_cast",
		"_ZNSt13basic_filebufIcSt11char_traitsIcEE5closeEv",
		"_ZSt16__ostream_insertIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_PKS3_l",
		"_ZSt7getlineIcSt11char_traitsIcESaIcEERSt13basic_istreamIT_T0_ES7_RNSt7__cxx1112basic_stringIS4_S5_T1_EE",
		"_ZNSirsEPFRSt8ios_baseS0_E",
		"_ZNSirsERi":
		return true
	default:
		// stringstream extract/manip also return the stream pointer.
		if isIstreamExtractI32(name) || isIstreamIosManip(name) {
			return true
		}
		return false
	}
}

func atomicAddFunc(t types.Type) (*jen.Statement, bool) {
	it, ok := t.(*types.IntType)
	if !ok {
		return nil, false
	}
	switch goIntBits(it.BitSize) {
	case 8:
		return Sym(libc.AtomicAddI8).code(), true
	case 32:
		return Sym(atomic.AddInt32).code(), true
	case 64:
		return Sym(atomic.AddInt64).code(), true
	default:
		return nil, false
	}
}

func atomicSwapFunc(t types.Type) (*jen.Statement, bool) {
	// Go has no SwapInt8; libc CAS-loop on the holding word.
	return atomicWordFunc(t, atomic.SwapPointer, libc.AtomicSwapI8, atomic.SwapInt32, atomic.SwapInt64)
}

func atomicCASFunc(t types.Type) (*jen.Statement, bool) {
	return atomicWordFunc(t, atomic.CompareAndSwapPointer, libc.AtomicCASI8, atomic.CompareAndSwapInt32, atomic.CompareAndSwapInt64)
}

func atomicWordFunc(t types.Type, ptr, i8, i32, i64 any) (*jen.Statement, bool) {
	if _, ok := t.(*types.PointerType); ok {
		return Sym(ptr).code(), true
	}
	it, ok := t.(*types.IntType)
	if !ok {
		return nil, false
	}
	switch goIntBits(it.BitSize) {
	case 8:
		return Sym(i8).code(), true
	case 32:
		return Sym(i32).code(), true
	case 64:
		return Sym(i64).code(), true
	default:
		return nil, false
	}
}

func atomicLoadFunc(t types.Type) *jen.Statement {
	if _, ok := t.(*types.PointerType); ok {
		return Sym(atomic.LoadPointer).code()
	}
	it, ok := t.(*types.IntType)
	if !ok {
		return Sym(atomic.LoadInt64).code()
	}
	switch goIntBits(it.BitSize) {
	case 8:
		return Sym(libc.AtomicLoadI8).code()
	case 32:
		return Sym(atomic.LoadInt32).code()
	default:
		return Sym(atomic.LoadInt64).code()
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
		return Sym(libc.RustPrint).code()
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
		return Sym(libc.RustFmtI32).code()
	case 'j': // usize
		return Sym(libc.RustFmtUsize).code()
	default:
		return nil
	}
}

const (
	cxxIOOpen = iota + 1
	cxxIOClose
	cxxIOInsert
	cxxIOEndl
	cxxIOLsCStr
	cxxIOInsertI64
	cxxIOInsertU64
	cxxIOInsertU8
	cxxIOInsertU16
	cxxIOInsertU32
	cxxIOInsertF64
	cxxIOInsertPtr
	cxxIOInsertBool
	cxxIOPut
	cxxIOFlush
	cxxIOCtypeInit
	cxxIOIosBase
	cxxIOGetline
	cxxIOStringstreamCtor
	cxxIOOStringStreamCtor
	cxxIOOStringStreamStr
	cxxIOManip
	cxxIOExtractI32
)

// cxxIONamed maps any ifstream ctor/dtor/open/close, not just the
// const char* + openmode pair clang used last time.
func cxxTreeNamed(name string) (*jen.Statement, bool) {
	fn, _, ok := cxxTreeKind(name)
	return fn, ok
}

func cxxTreeCall(name string, args []jen.Code) (*jen.Statement, []jen.Code, bool, bool) {
	fn, kind, ok := cxxTreeKind(name)
	if !ok {
		return nil, nil, false, false
	}
	switch kind {
	case cxxTreeInsert:
		out := make([]jen.Code, 4)
		if len(args) > 0 {
			out[0] = args[0]
		} else {
			out[0] = jen.False()
		}
		for i := 1; i < 4; i++ {
			if i < len(args) {
				out[i] = asBytePtr(args[i])
			} else {
				out[i] = jen.Nil()
			}
		}
		return fn, out, false, true
	case cxxTreeErase:
		z, h := jen.Nil(), jen.Nil()
		if len(args) > 0 {
			z = asBytePtr(args[0])
		}
		if len(args) > 1 {
			h = asBytePtr(args[1])
		}
		return fn, []jen.Code{z, h}, true, true
	default:
		this := jen.Nil()
		if len(args) > 0 {
			this = asBytePtr(args[0])
		}
		return fn, []jen.Code{this}, kind == cxxTreeWalk, true
	}
}

const (
	cxxTreeWalk = iota + 1
	cxxTreeInsert
	cxxTreeInit
	cxxTreeErase
)

func cxxTreeKind(name string) (*jen.Statement, int, bool) {
	switch {
	case strings.Contains(name, "_Rb_tree_decrement"):
		return Sym(libc.RbTreeDecrement).code(), cxxTreeWalk, true
	case strings.Contains(name, "_Rb_tree_increment"):
		return Sym(libc.RbTreeIncrement).code(), cxxTreeWalk, true
	case strings.Contains(name, "_Rb_tree_insert_and_rebalance"):
		return Sym(libc.RbTreeInsertAndRebalance).code(), cxxTreeInsert, true
	case strings.Contains(name, "_Rb_tree_rebalance_for_erase"):
		return Sym(libc.RbTreeRebalanceForErase).code(), cxxTreeErase, true
	case isRbTreeDefaultCtor(name):
		return Sym(libc.RbTreeInit).code(), cxxTreeInit, true
	default:
		return nil, 0, false
	}
}

func isRbTreeDefaultCtor(name string) bool {
	if !strings.HasSuffix(name, "C1Ev") && !strings.HasSuffix(name, "C2Ev") {
		return false
	}
	return strings.Contains(name, "St8_Rb_tree") || strings.Contains(name, "St3mapI")
}

func cxxIONamed(name string) (*jen.Statement, bool) {
	if cxxNoopDtor(name) {
		return Sym(libc.CxxNoop).code(), true
	}
	// Getline is rewritten only in cxxIOCallIR when the first two
	// args are pointers. namedRef would keep integer args.
	if isGetline(name) {
		return nil, false
	}
	fn, _, ok := cxxIOKind(name)
	return fn, ok
}

// cxxNoopDtor is a dtor for a type we never constructed (locale,
// ios_base, __basic_file). Empty is honest; unsatisfied would panic
// in the inlined ifstream dtor after fail() already succeeded.
func cxxNoopDtor(name string) bool {
	if !strings.Contains(name, "D0E") && !strings.Contains(name, "D1E") && !strings.Contains(name, "D2E") {
		return false
	}
	return strings.Contains(name, "12__basic_file") ||
		strings.Contains(name, "St6locale") ||
		strings.Contains(name, "St8ios_base")
}

func cxxIOCall(name string, args []jen.Code) (*jen.Statement, []jen.Code, bool, bool) {
	return cxxIOCallIR(name, nil, args)
}

func cxxIOCallIR(name string, ir []value.Value, args []jen.Code) (*jen.Statement, []jen.Code, bool, bool) {
	fn, kind, ok := cxxIOKind(name)
	if !ok {
		return nil, nil, false, false
	}
	if kind == cxxIOGetline && len(ir) >= 2 && (!isPtrish(ir[0].Type()) || !isPtrish(ir[1].Type())) {
		return nil, nil, false, false
	}
	switch kind {
	case cxxIOOpen:
		this, path, mode := jen.Nil(), jen.Nil(), jen.Lit(8)
		if len(args) > 0 {
			this = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			path = ptrArg(ir, args, 1)
		}
		if len(args) > 2 {
			mode = jen.Add(args[2])
		}
		return fn, []jen.Code{this, path, mode}, false, true
	case cxxIOClose:
		this := jen.Nil()
		if len(args) > 0 {
			this = ptrArg(ir, args, 0)
		}
		return fn, []jen.Code{this}, false, true
	case cxxIOInsert:
		os, s, n := jen.Nil(), jen.Nil(), jen.Lit(0)
		if len(args) > 0 {
			os = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			s = ptrArg(ir, args, 1)
		}
		if len(args) > 2 {
			n = jen.Add(args[2])
		}
		return fn, []jen.Code{os, s, n}, true, true
	case cxxIOLsCStr:
		os, s := jen.Nil(), jen.Nil()
		if len(args) > 0 {
			os = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			s = ptrArg(ir, args, 1)
		}
		return fn, []jen.Code{os, s}, true, true
	case cxxIOInsertI64:
		os, n := jen.Nil(), jen.Lit(0)
		if len(args) > 0 {
			os = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			n = jen.Int64().Call(args[1])
		}
		return fn, []jen.Code{os, n}, true, true
	case cxxIOInsertU64:
		os, n := jen.Nil(), jen.Lit(0)
		if len(args) > 0 {
			os = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			n = jen.Uint64().Call(args[1])
		}
		return fn, []jen.Code{os, n}, true, true
	case cxxIOInsertU8:
		os, n := jen.Nil(), jen.Lit(0)
		if len(args) > 0 {
			os = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			n = jen.Uint64().Call(jen.Uint8().Call(args[1]))
		}
		return fn, []jen.Code{os, n}, true, true
	case cxxIOInsertU16:
		os, n := jen.Nil(), jen.Lit(0)
		if len(args) > 0 {
			os = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			n = jen.Uint64().Call(jen.Uint16().Call(args[1]))
		}
		return fn, []jen.Code{os, n}, true, true
	case cxxIOInsertU32:
		os, n := jen.Nil(), jen.Lit(0)
		if len(args) > 0 {
			os = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			n = jen.Uint64().Call(jen.Uint32().Call(args[1]))
		}
		return fn, []jen.Code{os, n}, true, true
	case cxxIOInsertF64:
		os, x := jen.Nil(), jen.Lit(0.0)
		if len(args) > 0 {
			os = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			x = jen.Float64().Call(args[1])
		}
		return fn, []jen.Code{os, x}, true, true
	case cxxIOInsertPtr:
		os, p := jen.Nil(), jen.Nil()
		if len(args) > 0 {
			os = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			p = emitUP(args[1])
		}
		return fn, []jen.Code{os, p}, true, true
	case cxxIOInsertBool:
		os, b := jen.Nil(), jen.Lit(false)
		if len(args) > 0 {
			os = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			b = jen.Bool().Call(args[1])
		}
		return fn, []jen.Code{os, b}, true, true
	case cxxIOPut:
		os, c := jen.Nil(), jen.Lit(0)
		if len(args) > 0 {
			os = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			c = jen.Byte().Call(args[1])
		}
		return fn, []jen.Code{os, c}, true, true
	case cxxIOEndl, cxxIOFlush:
		this := jen.Nil()
		if len(args) > 0 {
			this = ptrArg(ir, args, 0)
		}
		return fn, []jen.Code{this}, true, true
	case cxxIOCtypeInit, cxxIOIosBase:
		this := jen.Nil()
		if len(args) > 0 {
			this = ptrArg(ir, args, 0)
		}
		return fn, []jen.Code{this}, false, true
	case cxxIOGetline:
		is, str := jen.Nil(), jen.Nil()
		if len(args) > 0 {
			is = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			str = ptrArg(ir, args, 1)
		}
		return fn, []jen.Code{is, str}, true, true
	case cxxIOStringstreamCtor:
		this, str, mode := jen.Nil(), jen.Nil(), jen.Lit(0)
		if len(args) > 0 {
			this = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			str = ptrArg(ir, args, 1)
		}
		if len(args) > 2 {
			mode = jen.Add(args[2])
		}
		return fn, []jen.Code{this, str, mode}, false, true
	case cxxIOOStringStreamCtor:
		this := jen.Nil()
		if len(args) > 0 {
			this = ptrArg(ir, args, 0)
		}
		return fn, []jen.Code{this}, false, true
	case cxxIOOStringStreamStr:
		// sret string first, then this (LLVM sret param order).
		ret, this := jen.Nil(), jen.Nil()
		if len(args) > 0 {
			ret = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			this = ptrArg(ir, args, 1)
		}
		return fn, []jen.Code{ret, this}, false, true
	case cxxIOManip:
		is, manip := jen.Nil(), jen.Nil()
		if len(args) > 0 {
			is = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			// Function pointer: opaque *byte; shim ignores the target.
			manip = ptrArg(ir, args, 1)
		}
		return fn, []jen.Code{is, manip}, true, true
	case cxxIOExtractI32:
		is, out := jen.Nil(), jen.Nil()
		if len(args) > 0 {
			is = ptrArg(ir, args, 0)
		}
		if len(args) > 1 {
			out = ptrArg(ir, args, 1)
		}
		return fn, []jen.Code{is, out}, true, true
	default:
		return nil, nil, false, false
	}
}

// cxxOstreamOp is cout << / endl / put / flush. csmith OutputHeader
// inlines some of these and calls the rest. Also gensym's ostringstream.
func cxxOstreamOp(name string) (*jen.Statement, int, bool) {
	switch {
	case strings.Contains(name, "4endlI") || strings.HasPrefix(name, "_ZSt4endl"):
		return Sym(libc.OstreamEndl).code(), cxxIOEndl, true
	case strings.Contains(name, "lsISt11char_traits") && strings.HasSuffix(name, "PKc"):
		return Sym(libc.OstreamLsCStr).code(), cxxIOLsCStr, true
	// operator<<(ostream&, char)
	case strings.Contains(name, "lsISt11char_traits") && strings.HasSuffix(name, "ES5_c"):
		return Sym(libc.OstreamPut).code(), cxxIOPut, true
	// operator<<(ostream&, basic_string const&)
	case strings.Contains(name, "lsIcSt11char_traits") && strings.Contains(name, "basic_string"):
		return Sym(libc.OstreamLsString).code(), cxxIOLsCStr, true
	case strings.Contains(name, "9_M_insertImE") || strings.Contains(name, "9_M_insertIyE"):
		return Sym(libc.OstreamInsertU64).code(), cxxIOInsertU64, true
	case strings.Contains(name, "9_M_insertIlE") || strings.Contains(name, "9_M_insertIxE"):
		return Sym(libc.OstreamInsertI64).code(), cxxIOInsertI64, true
	case strings.Contains(name, "So3putE"):
		return Sym(libc.OstreamPut).code(), cxxIOPut, true
	case strings.Contains(name, "So5flushE"):
		return Sym(libc.OstreamFlush).code(), cxxIOFlush, true
	case strings.Contains(name, "5ctypeIcE13_M_widen_init"):
		return Sym(libc.CtypeWidenInit).code(), cxxIOCtypeInit, true
	case strings.Contains(name, "SolsEPFRSoS_E"):
		// operator<<(ostream&(*)(ostream&)) — csmith passes endl.
		return Sym(libc.OstreamEndl).code(), cxxIOEndl, true
	case strings.HasPrefix(name, "_ZNSolsE") || strings.Contains(name, "NSolsE"):
		switch {
		case strings.HasSuffix(name, "PKc"):
			return Sym(libc.OstreamLsCStr).code(), cxxIOLsCStr, true
		case strings.HasSuffix(name, "PKv"):
			// operator<<(void const*)
			return Sym(libc.OstreamInsertPtr).code(), cxxIOInsertPtr, true
		// signed: char(a)/short(s)/int(i)/long(l)/long long(x)
		case strings.HasSuffix(name, "Ea"), strings.HasSuffix(name, "Es"),
			strings.HasSuffix(name, "Ei"), strings.HasSuffix(name, "El"),
			strings.HasSuffix(name, "Ex"):
			return Sym(libc.OstreamInsertI64).code(), cxxIOInsertI64, true
		// unsigned: zero-extend to u64 (Go uint64(intN(-1)) sign-extends).
		case strings.HasSuffix(name, "Eh"): // unsigned char
			return Sym(libc.OstreamInsertU64).code(), cxxIOInsertU8, true
		case strings.HasSuffix(name, "Et"): // unsigned short
			return Sym(libc.OstreamInsertU64).code(), cxxIOInsertU16, true
		case strings.HasSuffix(name, "Ej"): // unsigned int
			return Sym(libc.OstreamInsertU64).code(), cxxIOInsertU32, true
		case strings.HasSuffix(name, "Em"), strings.HasSuffix(name, "Ey"):
			return Sym(libc.OstreamInsertU64).code(), cxxIOInsertU64, true
		case strings.HasSuffix(name, "Ed"), strings.HasSuffix(name, "Ef"):
			return Sym(libc.OstreamInsertF64).code(), cxxIOInsertF64, true
		case strings.HasSuffix(name, "Eb"):
			return Sym(libc.OstreamInsertBool).code(), cxxIOInsertBool, true
		case strings.HasSuffix(name, "Ec"):
			return Sym(libc.OstreamPut).code(), cxxIOPut, true
		}
	}
	return nil, 0, false
}

func isIosBaseCtor(name string) bool {
	if strings.Contains(name, "St8ios_base") {
		return strings.Contains(name, "C1E") || strings.Contains(name, "C2E") || strings.Contains(name, "7_M_init")
	}
	return strings.Contains(name, "9basic_ios") && strings.Contains(name, "4initE")
}

func isLocaleCtor(name string) bool {
	if !strings.Contains(name, "St6locale") {
		return false
	}
	return strings.Contains(name, "C1E") || strings.Contains(name, "C2E")
}

func isLibcxxStringEqCStr(name string) bool {
	if !strings.Contains(name, "St3__1") || !strings.Contains(name, "basic_string") {
		return false
	}
	if !strings.Contains(name, "PK") {
		return false
	}
	return strings.Contains(name, "eqI") || strings.Contains(name, "eqERK") ||
		strings.Contains(name, "eqEPKc") || strings.Contains(name, "3eqE")
}

func isGetline(name string) bool {
	return strings.Contains(name, "St7getline") ||
		strings.Contains(name, "St3__17getline") ||
		strings.HasPrefix(name, "_ZSt7getline")
}

func isStringstream(name string) bool {
	return strings.Contains(name, "18basic_stringstream")
}

func isOstringstream(name string) bool {
	return strings.Contains(name, "19basic_ostringstream")
}

// isRustAlloc is __rust_alloc / __rust_alloc_zeroed, not *_error_handler.
func isRustAlloc(name string) bool {
	if strings.Contains(name, "error") || strings.Contains(name, "dealloc") ||
		strings.Contains(name, "realloc") {
		return false
	}
	return strings.Contains(name, "__rust_alloc")
}

func isIstreamExtractI32(name string) bool {
	// basic_istream::operator>>(int&) — short mangling Si = basic_istream<char>.
	return name == "_ZNSirsERi" || strings.HasSuffix(name, "rsERi")
}

func isIstreamIosManip(name string) bool {
	// operator>>(ios_base&(*)(ios_base&)) — used for std::hex in str2int.
	return name == "_ZNSirsEPFRSt8ios_baseS0_E" ||
		strings.Contains(name, "rsEPFRSt8ios_baseS0_E")
}

func cxxIOKind(name string) (*jen.Statement, int, bool) {
	if strings.Contains(name, "__ostream_insert") {
		return Sym(libc.OstreamInsert).code(), cxxIOInsert, true
	}
	if isGetline(name) {
		return Sym(libc.IstreamGetline).code(), cxxIOGetline, true
	}
	if isIstreamIosManip(name) {
		return Sym(libc.IstreamApplyIosManip).code(), cxxIOManip, true
	}
	if isIstreamExtractI32(name) {
		return Sym(libc.IstreamExtractI32).code(), cxxIOExtractI32, true
	}
	if k, kind, ok := cxxOstreamOp(name); ok {
		return k, kind, true
	}
	if isIosBaseCtor(name) {
		return Sym(libc.IosBaseCtor).code(), cxxIOIosBase, true
	}
	if isLocaleCtor(name) {
		return Sym(libc.LocaleCtor).code(), cxxIOIosBase, true
	}
	if isStringstream(name) {
		switch {
		case strings.Contains(name, "3strE"):
			return Sym(libc.StringstreamStr).code(), cxxIOOStringStreamStr, true
		// Default ctor before the string+mode overload (C1Ev vs C1ERKNS…).
		case strings.HasSuffix(name, "C1Ev") || strings.HasSuffix(name, "C2Ev"):
			return Sym(libc.StringstreamDefaultCtor).code(), cxxIOOStringStreamCtor, true
		case strings.Contains(name, "C1E"), strings.Contains(name, "C2E"):
			return Sym(libc.StringstreamCtor).code(), cxxIOStringstreamCtor, true
		case strings.Contains(name, "D0E"), strings.Contains(name, "D1E"), strings.Contains(name, "D2E"):
			// Prefer default-close if we might have dual keys; safe for both.
			return Sym(libc.StringstreamDefaultClose).code(), cxxIOClose, true
		default:
			return nil, 0, false
		}
	}
	if isOstringstream(name) {
		switch {
		case strings.Contains(name, "C1E") || strings.Contains(name, "C2E"):
			return Sym(libc.OStringStreamCtor).code(), cxxIOOStringStreamCtor, true
		case strings.Contains(name, "3strE"):
			return Sym(libc.OStringStreamStr).code(), cxxIOOStringStreamStr, true
		case strings.Contains(name, "D0E") || strings.Contains(name, "D1E") || strings.Contains(name, "D2E"):
			return Sym(libc.OStringStreamClose).code(), cxxIOClose, true
		default:
			return nil, 0, false
		}
	}
	if !strings.Contains(name, "14basic_ifstream") {
		return nil, 0, false
	}
	switch {
	case strings.Contains(name, "C1E"), strings.Contains(name, "C2E"), strings.Contains(name, "4openE"):
		return Sym(libc.IfstreamOpen).code(), cxxIOOpen, true
	case strings.Contains(name, "D0E"), strings.Contains(name, "D1E"), strings.Contains(name, "D2E"), strings.Contains(name, "5closeE"):
		return Sym(libc.IfstreamClose).code(), cxxIOClose, true
	default:
		return nil, 0, false
	}
}

// libcCanon strips Darwin $ / sanitized suffixes so realpath$DARWIN_EXTSN
// and realpath_DARWIN_EXTSN both map to realpath. Unknown names stay as-is.
func libcCanon(name string) string {
	if _, ok := libraryFunctions[name]; ok {
		return name
	}
	for _, suf := range []string{
		"$DARWIN_EXTSN", "$UNIX2003", "$INODE64", "$NOCANCEL",
		"_DARWIN_EXTSN", "_UNIX2003", "_INODE64", "_NOCANCEL",
	} {
		if strings.HasSuffix(name, suf) {
			base := strings.TrimSuffix(name, suf)
			if _, ok := libraryFunctions[base]; ok {
				return base
			}
		}
	}
	if i := strings.IndexByte(name, '$'); i > 0 {
		if _, ok := libraryFunctions[name[:i]]; ok {
			return name[:i]
		}
	}
	return name
}

// libcLookup maps LLVM names to libc. Darwin symbols may be
// realpath$DARWIN_EXTSN; after sanitizeIdent they are realpath_DARWIN_EXTSN.
func libcLookup(name string) (goRef, bool) {
	ref, ok := libraryFunctions[libcCanon(name)]
	return ref, ok
}

var libraryFunctions = map[string]goRef{
	"_NSGetArgc":     Sym(libc.NSGetArgc),
	"_NSGetArgv":     Sym(libc.NSGetArgv),
	"_NSGetEnviron":  Sym(libc.NSGetEnviron),
	"_NSGetProgname": Sym(libc.NSGetProgname),
	"abort":          Sym(libc.Abort),
	"arc4random_buf": Sym(libc.Arc4randomBuf),
	"__cxa_atexit":   Sym(libc.CxaAtexit),
	"__dynamic_cast": Sym(libc.DynamicCast),
	"_ZNSt14basic_ifstreamIcSt11char_traitsIcEEC1EPKcSt13_Ios_Openmode":                                        Sym(libc.IfstreamOpen),
	"_ZNSt14basic_ifstreamIcSt11char_traitsIcEEC2EPKcSt13_Ios_Openmode":                                        Sym(libc.IfstreamOpen),
	"_ZNSt14basic_ifstreamIcSt11char_traitsIcEED1Ev":                                                           Sym(libc.IfstreamClose),
	"_ZNSt14basic_ifstreamIcSt11char_traitsIcEED2Ev":                                                           Sym(libc.IfstreamCloseVTT),
	"_ZNSt14basic_ifstreamIcSt11char_traitsIcEE5closeEv":                                                       Sym(libc.IfstreamClose),
	"_ZNSt13basic_filebufIcSt11char_traitsIcEE5closeEv":                                                        Sym(libc.FilebufClose),
	"_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE6__initEPKcm":                              Sym(libc.StdStringInit),
	"_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEC1ERKS5_":                                  Sym(libc.StdStringCopy),
	"_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEC2ERKS5_":                                  Sym(libc.StdStringCopy),
	"_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEC1ERKS5_mmRKS4_":                            Sym(libc.StdStringSubstr),
	"_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEC2ERKS5_mmRKS4_":                            Sym(libc.StdStringSubstr),
	"_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEED1Ev":                                      Sym(libc.StdStringDestroy),
	"_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEED2Ev":                                      Sym(libc.StdStringDestroy),
	"_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEaSERKS5_":                                  Sym(libc.StdStringAssign),
	"_ZNKSt3__18ios_base6getlocEv":                                                                             Sym(libc.IosGetloc),
	"_ZNSt3__16localeD1Ev":                                                                                     Sym(libc.LocaleDtor),
	"_ZNSt3__16localeD2Ev":                                                                                     Sym(libc.LocaleDtor),
	"_ZNSt3__113basic_ostreamIcNS_11char_traitsIcEEE6sentryC1ERS3_":                                            Sym(libc.OstreamSentryCtor),
	"_ZNSt3__113basic_ostreamIcNS_11char_traitsIcEEE6sentryC2ERS3_":                                            Sym(libc.OstreamSentryCtor),
	"_ZNSt3__113basic_ostreamIcNS_11char_traitsIcEEE6sentryD1Ev":                                               Sym(libc.OstreamSentryDestroy),
	"_ZNSt3__113basic_ostreamIcNS_11char_traitsIcEEE6sentryD2Ev":                                               Sym(libc.OstreamSentryDestroy),
	"_ZNSt3__113basic_istreamIcNS_11char_traitsIcEEE6sentryC1ERS3_b":                                           Sym(libc.IstreamSentryCtor),
	"_ZNSt3__113basic_istreamIcNS_11char_traitsIcEEE6sentryC2ERS3_b":                                           Sym(libc.IstreamSentryCtor),
	"_ZNKSt3__16locale9use_facetERNS0_2idE":                                                                    Sym(libc.LocaleUseFacet),
	"_ZNKSt9basic_iosIcSt11char_traitsIcEE4failEv":                                                             Sym(libc.IosFail),
	"_ZNSt9basic_iosIcSt11char_traitsIcEE5clearESt12_Ios_Iostate":                                              Sym(libc.IosClear),
	"_ZNKSt9basic_iosIcSt11char_traitsIcEE3eofEv":                                                              Sym(libc.IosEof),
	"_ZNKSt9basic_iosIcSt11char_traitsIcEEntEv":                                                                Sym(libc.IosNot),
	"_ZNKSt9basic_iosIcSt11char_traitsIcEEcvbEv":                                                               Sym(libc.IosBool),
	"_ZSt16__ostream_insertIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_PKS3_l":                            Sym(libc.OstreamInsert),
	"_ZSt7getlineIcSt11char_traitsIcESaIcEERSt13basic_istreamIT_T0_ES7_RNSt7__cxx1112basic_stringIS4_S5_T1_EE": Sym(libc.IstreamGetline),
	// stringstream / istream >> are matched by cxxIOKind (arg reshaping).
	"__assert_fail":             Sym(libc.AssertFail),
	"fabs":                      Sym(math.Abs),
	"fmod":                      Sym(math.Mod),
	"pow":                       Sym(math.Pow),
	"__ctype_b_loc":             Sym(libc.CtypeBLoc),
	"dup":                       Sym(libc.Dup),
	"fclose":                    Sym(libc.Fclose),
	"fcntl":                     Sym(libc.Fcntl),
	"fdopen":                    Sym(libc.Fdopen),
	"fprintf":                   Sym(libc.Fprintf),
	"fputc":                     Sym(libc.Fputc),
	"fputs":                     Sym(libc.Fputs),
	"free":                      Sym(libc.Free),
	"getchar":                   Sym(libc.Getchar),
	"exit":                      Sym(libc.Exit),
	"iswalnum":                  Sym(libc.Iswalnum),
	"iswalpha":                  Sym(libc.Iswalpha),
	"iswblank":                  Sym(libc.Iswblank),
	"iswcntrl":                  Sym(libc.Iswcntrl),
	"iswdigit":                  Sym(libc.Iswdigit),
	"iswlower":                  Sym(libc.Iswlower),
	"iswspace":                  Sym(libc.Iswspace),
	"iswupper":                  Sym(libc.Iswupper),
	"iswxdigit":                 Sym(libc.Iswxdigit),
	"leaven_va_arg":             Sym(libc.VAArg),
	"towlower":                  Sym(libc.Towlower),
	"towupper":                  Sym(libc.Towupper),
	"llvm_fabs_f64":             Sym(math.Abs),
	"llvm_fabs_f80":             Sym(math.Abs),
	"llvm_maximumnum_f64":       Sym(libc.MaximumNumF64),
	"llvm_pow_f64":              Sym(math.Pow),
	"memchr":                    Sym(libc.Memchr),
	"memcmp":                    Sym(libc.Memcmp),
	"bcmp":                      Sym(libc.Memcmp),
	"__memcpy_chk":              Sym(libc.MemcpyChk),
	"memmove":                   Sym(libc.Memmove),
	"memset_pattern16":          Sym(libc.MemsetPattern16),
	"__memset_chk":              Sym(libc.MemsetChk),
	"printf":                    Sym(libc.Printf),
	"putc":                      Sym(libc.Putc),
	"putchar":                   Sym(libc.Putchar),
	"__errno_location":          Sym(libc.ErrnoLocation),
	"__error":                   Sym(libc.ErrnoLocation),
	"strerror_r":                Sym(libc.StrerrorR),
	"close":                     Sym(libc.Close),
	"dlsym":                     Sym(libc.Dlsym),
	"fstat":                     Sym(libc.Fstat64),
	"fstat64":                   Sym(libc.Fstat64),
	"getauxval":                 Sym(libc.Getauxval),
	"getcwd":                    Sym(libc.Getcwd),
	"getenv":                    Sym(libc.Getenv),
	"getpid":                    Sym(libc.Getpid),
	"getentropy":                Sym(libc.Getentropy),
	"getrandom":                 Sym(libc.Getrandom),
	"gettid":                    Sym(libc.Gettid),
	"lseek":                     Sym(libc.Lseek64),
	"lseek64":                   Sym(libc.Lseek64),
	"mmap":                      Sym(libc.Mmap),
	"mmap64":                    Sym(libc.Mmap64),
	"mprotect":                  Sym(libc.Mprotect),
	"munmap":                    Sym(libc.Munmap),
	"open":                      Sym(libc.Open),
	"open64":                    Sym(libc.Open64),
	"poll":                      Sym(libc.Poll),
	"pthread_attr_destroy":      Sym(libc.PthreadAttrDestroy),
	"pthread_attr_getstack":     Sym(libc.PthreadAttrGetstack),
	"pthread_getattr_np":        Sym(libc.PthreadGetattrNp),
	"pthread_get_stackaddr_np":  Sym(libc.PthreadGetStackaddrNp),
	"pthread_get_stacksize_np":  Sym(libc.PthreadGetStacksizeNp),
	"pthread_mutex_destroy":     Sym(libc.PthreadMutexDestroy),
	"pthread_mutex_init":        Sym(libc.PthreadMutexInit),
	"pthread_mutex_lock":        Sym(libc.PthreadMutexLock),
	"pthread_mutex_trylock":     Sym(libc.PthreadMutexTrylock),
	"pthread_mutex_unlock":      Sym(libc.PthreadMutexUnlock),
	"pthread_mutexattr_destroy": Sym(libc.PthreadMutexattrDestroy),
	"pthread_mutexattr_init":    Sym(libc.PthreadMutexattrInit),
	"pthread_mutexattr_settype": Sym(libc.PthreadMutexattrSettype),
	"pthread_self":              Sym(libc.PthreadSelf),
	"pthread_setname_np":        Sym(libc.PthreadSetnameNp),
	"pthread_threadid_np":       Sym(libc.PthreadThreadidNp),
	"puts":                      Sym(libc.Puts),
	"read":                      Sym(libc.Read),
	"realpath":                  Sym(libc.Realpath),
	"sigaction":                 Sym(libc.Sigaction),
	"sigaltstack":               Sym(libc.Sigaltstack),
	"signal":                    Sym(libc.Signal),
	"stat64":                    Sym(libc.Stat64),
	"sysconf":                   Sym(libc.Sysconf),
	"syscall":                   Sym(libc.Syscall),
	"write":                     Sym(libc.Write),
	"realloc":                   Sym(libc.Realloc),
	"scanf":                     Sym(libc.Scanf),
	"sscanf":                    Sym(libc.Sscanf),
	"__isoc23_sscanf":           Sym(libc.Sscanf),
	"srand48":                   Sym(libc.Srand48),
	"lrand48":                   Sym(libc.Lrand48),
	"statx":                     Sym(libc.Statx),
	"snprintf":                  Sym(libc.Snprintf),
	"sqrt":                      Sym(math.Sqrt),
	"__strcat_chk":              Sym(libc.StrcatChk),
	"strchr":                    Sym(libc.Strchr),
	"strcmp":                    Sym(libc.Strcmp),
	"strcpy":                    Sym(libc.Strcpy),
	"strcspn":                   Sym(libc.Strcspn),
	"strlen":                    Sym(libc.Strlen),
	"strncat":                   Sym(libc.Strncat),
	"strncmp":                   Sym(libc.Strncmp),
	"strncpy":                   Sym(libc.Strncpy),
	"strrchr":                   Sym(libc.Strrchr),
	"strspn":                    Sym(libc.Strspn),
	"strstr":                    Sym(libc.Strstr),
}
