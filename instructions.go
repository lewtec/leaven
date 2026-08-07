package main

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

// TranslateInstruction translates an LLVM instruction to Go.
func TranslateInstruction(inst ir.Instruction) (string, error) {
	switch inst := inst.(type) {
	case *ir.InstAdd:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return fmt.Sprintf("for i, v := range %s { %s[i] = v + %s[i] }", x, VariableName(inst), y), nil
		}
		if ciy, ok := inst.Y.(*constant.Int); ok && ciy.X.Sign() == -1 {
			return fmt.Sprintf("%s = %s %s", VariableName(inst), x, ciy.X), nil // Use the constant's own minus sign.
		}
		return fmt.Sprintf("%s = %s + %s", VariableName(inst), x, y), nil

	case *ir.InstAtomicRMW:
		dst, err := FormatValue(inst.Dst)
		if err != nil {
			return "", fmt.Errorf("error translating destination (%v): %w", inst.Dst, err)
		}
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating operand (%v): %w", inst.X, err)
		}
		name := VariableName(inst)
		addFn, ok := atomicAddFunc(inst.Type())
		if !ok {
			return "", fmt.Errorf("%w: atomicrmw on %v", errUnsupportedInstruction, inst.Type())
		}
		switch inst.Op {
		case enum.AtomicOpAdd:
			// atomicrmw returns the old value; Add* returns the new value.
			return fmt.Sprintf("%s = %s(%s, %s) - %s", name, addFn, dst, x, x), nil
		case enum.AtomicOpSub:
			return fmt.Sprintf("%s = %s(%s, -(%s)) + %s", name, addFn, dst, x, x), nil
		case enum.AtomicOpXChg:
			swapFn, ok := atomicSwapFunc(inst.Type())
			if !ok {
				return "", fmt.Errorf("%w: atomicrmw xchg on %v", errUnsupportedInstruction, inst.Type())
			}
			return fmt.Sprintf("%s = %s(%s, %s)", name, swapFn, dst, x), nil
		default:
			return "", fmt.Errorf("%w: atomicrmw %v", errUnsupportedInstruction, inst.Op)
		}

	case *ir.InstAlloca:
		t, err := TypeSpec(inst.ElemType)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", inst.ElemType, err)
		}
		if inst.NElems == nil {
			if _, ok := inst.ElemType.(*types.ArrayType); ok {
				// If it's an array, allocate an extra byte to allow indexing off the end.
				return fmt.Sprintf("%s = &new(struct{v %s; b byte}).v", VariableName(inst), t), nil
			}
			return fmt.Sprintf("%s = new(%s)", VariableName(inst), t), nil
		}
		nElems, err := FormatValue(inst.NElems)
		if err != nil {
			return "", fmt.Errorf("error translating NElems (%v): %w", inst.NElems, err)
		}
		return fmt.Sprintf("%s = &make([]%s, %s + 1)[0]", VariableName(inst), t, nElems), nil

	case *ir.InstAnd:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return fmt.Sprintf("for i, v := range %s { %s[i] = v & %s[i] }", x, VariableName(inst), y), nil
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 1 {
			return fmt.Sprintf("%s = %s && %s", VariableName(inst), x, y), nil
		}
		return fmt.Sprintf("%s = %s & %s", VariableName(inst), x, y), nil

	case *ir.InstAShr:
		x, err := FormatSigned(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		if t, ok := inst.Typ.(*types.IntType); ok && t.BitSize == 8 {
			return fmt.Sprintf("%s = byte(%s >> %s)", VariableName(inst), x, y), nil
		}
		return fmt.Sprintf("%s = %s >> %s", VariableName(inst), x, y), nil

	case *ir.InstBitCast:
		if !compatiblePointerTypes(inst.From.Type(), inst.To) {
			return "", fmt.Errorf("%w: %v and %v", errIncompatiblePointers, inst.From.Type(), inst.To)
		}
		from, err := FormatValue(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", inst.To, err)
		}
		return fmt.Sprintf("%s = (%s)(unsafe.Pointer(%s))", VariableName(inst), to, from), nil

	case *ir.InstCall:
		callee, err := FormatValue(inst.Callee)
		if err != nil {
			return "", fmt.Errorf("error translating callee (%v): %w", inst.Callee, err)
		}
		args := make([]string, len(inst.Args))
		for i, a := range inst.Args {
			v, err := FormatValue(a)
			if err != nil {
				return "", fmt.Errorf("error translating argument %d (%v): %w", i, a, err)
			}
			args[i] = v
		}
		if renamed, ok := libraryFunctions[callee]; ok {
			callee = renamed
		}
		switch callee {
		case "calloc", "malloc":
			if pt, ok := inst.Typ.(*types.PointerType); ok {
				if et, err := TypeSpec(pt.ElemType); err == nil {
					callee = fmt.Sprintf("libc.%s[%s]", strings.Title(callee), et)
				}
			}
		case "leaven_va_start":
			if len(args) == 1 {
				return fmt.Sprintf("*%s = (*byte)(unsafe.Pointer(&varargs))", args[0]), nil
			}
		case "ldexp":
			if len(args) == 2 {
				return fmt.Sprintf("%s = math.Ldexp(%s, int(%s))", VariableName(inst), args[0], args[1]), nil
			}
		case "llvm_fabs_f32":
			if len(args) == 1 {
				return fmt.Sprintf("%s = float32(math.Abs(float64(%s)))", VariableName(inst), args[0]), nil
			}
		case "llvm_fmuladd_f64":
			if len(args) == 3 {
				return fmt.Sprintf("%s = %s*%s + %s", VariableName(inst), args[0], args[1], args[2]), nil
			}
		case "llvm_fmuladd_f32":
			if len(args) == 3 {
				return fmt.Sprintf("%s = %s*%s + %s", VariableName(inst), args[0], args[1], args[2]), nil
			}
		case "llvm_lifetime_start", "llvm_lifetime_end":
			return ";", nil
		case "llvm_memcpy_p0i8_p0i8_i64":
			return fmt.Sprintf("libc.Memmove(%s, %s, %s)", args[0], args[1], args[2]), nil
		case "llvm_memset_p0i8_i64":
			return fmt.Sprintf("libc.Memset(%s, %s, %s)", args[0], args[1], args[2]), nil
		case "llvm_objectsize_i64_p0i8":
			// Use -1 for unknown size.
			return fmt.Sprintf("%s = -1", VariableName(inst)), nil
		case "llvm_stacksave":
			// Use nil, since we're doing llvm_stackrestore as a no-op.
			return fmt.Sprintf("%s = nil", VariableName(inst)), nil
		case "llvm_stackrestore":
			return ";", nil
		}
		if types.Equal(inst.Type(), types.Void) {
			return fmt.Sprintf("%s(%s)", callee, strings.Join(args, ", ")), nil
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), callee, strings.Join(args, ", ")), nil

	case *ir.InstExtractElement:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating vector (%v): %w", inst.X, err)
		}
		index, err := FormatValue(inst.Index)
		if err != nil {
			return "", fmt.Errorf("error translating index (%v): %w", inst.Index, err)
		}
		return fmt.Sprintf("%s = %s[%s]", VariableName(inst), x, index), nil

	case *ir.InstExtractValue:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating aggregate (%v): %w", inst.X, err)
		}
		expr, err := formatAggregateIndex(x, inst.X.Type(), inst.Indices)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s = %s", VariableName(inst), expr), nil

	case *ir.InstInsertValue:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating aggregate (%v): %w", inst.X, err)
		}
		elem, err := FormatValue(inst.Elem)
		if err != nil {
			return "", fmt.Errorf("error translating element (%v): %w", inst.Elem, err)
		}
		dest, err := formatAggregateIndex(VariableName(inst), inst.Typ, inst.Indices)
		if err != nil {
			return "", err
		}
		// Copy aggregate then assign the field/element.
		return fmt.Sprintf("%s = %s; %s = %s", VariableName(inst), x, dest, elem), nil

	case *ir.InstFAdd:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s + %s", VariableName(inst), x, y), nil

	case *ir.InstFCmp:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}

		var op string
		switch inst.Pred {
		case enum.FPredOEQ:
			op = "=="
		case enum.FPredOGE:
			op = ">="
		case enum.FPredOGT:
			op = ">"
		case enum.FPredOLE:
			op = "<="
		case enum.FPredOLT:
			op = "<"
		case enum.FPredUNE:
			op = "!="
		case enum.FPredORD:
			return fmt.Sprintf("%s = %s == %s && %s == %s", VariableName(inst), x, x, y, y), nil
		case enum.FPredUNO:
			return fmt.Sprintf("%s = %s != %s || %s != %s", VariableName(inst), x, x, y, y), nil
		case enum.FPredUEQ:
			return fmt.Sprintf("%s = %s != %s || %s != %s || %s == %s", VariableName(inst), x, x, y, y, x, y), nil
		case enum.FPredUGT:
			return fmt.Sprintf("%s = %s != %s || %s != %s || %s > %s", VariableName(inst), x, x, y, y, x, y), nil
		case enum.FPredUGE:
			return fmt.Sprintf("%s = %s != %s || %s != %s || %s >= %s", VariableName(inst), x, x, y, y, x, y), nil
		case enum.FPredULT:
			return fmt.Sprintf("%s = %s != %s || %s != %s || %s < %s", VariableName(inst), x, x, y, y, x, y), nil
		case enum.FPredULE:
			return fmt.Sprintf("%s = %s != %s || %s != %s || %s <= %s", VariableName(inst), x, x, y, y, x, y), nil
		case enum.FPredONE:
			return fmt.Sprintf("%s = %s == %s && %s == %s && %s != %s", VariableName(inst), x, x, y, y, x, y), nil
		default:
			return "", fmt.Errorf("%w: %v", errUnsupportedICmpPred, inst.Pred)
		}

		return fmt.Sprintf("%s = %s %s %s", VariableName(inst), x, op, y), nil

	case *ir.InstFDiv:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s / %s", VariableName(inst), x, y), nil

	case *ir.InstFMul:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s * %s", VariableName(inst), x, y), nil

	case *ir.InstFNeg:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating operand (%v): %w", inst.X, err)
		}
		return fmt.Sprintf("%s = - %s", VariableName(inst), x), nil

	case *ir.InstFPExt:
		from, err := FormatValue(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", inst.To, err)
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), to, from), nil

	case *ir.InstFPToSI:
		from, err := FormatValue(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", inst.To, err)
		}
		if to == "byte" {
			return fmt.Sprintf("%s = byte(int8(%s))", VariableName(inst), from), nil
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), to, from), nil

	case *ir.InstFPTrunc:
		from, err := FormatValue(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", inst.To, err)
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), to, from), nil

	case *ir.InstFSub:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s - %s", VariableName(inst), x, y), nil

	case *ir.InstGetElementPtr:
		result, err := GetElementPtr(inst.ElemType, inst.Src, inst.Indices)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s = %s", VariableName(inst), result), nil

	case *ir.InstICmp:
		cmp, err := formatICmp(inst.Pred, inst.X, inst.Y)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s = %s", VariableName(inst), cmp), nil

	case *ir.InstInsertElement:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating initial vector (%v): %w", inst.X, err)
		}
		elem, err := FormatValue(inst.Elem)
		if err != nil {
			return "", fmt.Errorf("error translating new element (%v): %w", inst.Elem, err)
		}
		index, err := FormatValue(inst.Index)
		if err != nil {
			return "", fmt.Errorf("error translating index (%v): %w", inst.Index, err)
		}
		if _, ok := inst.X.(*constant.Undef); ok {
			return fmt.Sprintf("%s[%s] = %s", VariableName(inst), index, elem), nil
		}
		return fmt.Sprintf("%s = %s; %s[%s] = %s", VariableName(inst), x, VariableName(inst), index, elem), nil

	case *ir.InstIntToPtr:
		return "", errIntToPtr

	case *ir.InstLoad:
		src, err := FormatValue(inst.Src)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", inst.Src, err)
		}
		if strings.HasPrefix(src, "&") {
			return fmt.Sprintf("%s = %s", VariableName(inst), strings.TrimPrefix(src, "&")), nil
		}
		return fmt.Sprintf("%s = *%s", VariableName(inst), src), nil

	case *ir.InstLShr:
		x, err := FormatUnsigned(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		if t, ok := inst.Typ.(*types.IntType); ok && t.BitSize > 8 {
			return fmt.Sprintf("%s = int%d(%s >> %s)", VariableName(inst), goIntBits(t.BitSize), x, y), nil
		}
		return fmt.Sprintf("%s = %s >> %s", VariableName(inst), x, y), nil

	case *ir.InstMul:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s * %s", VariableName(inst), x, y), nil

	case *ir.InstOr:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return fmt.Sprintf("for i, v := range %s { %s[i] = v | %s[i] }", x, VariableName(inst), y), nil
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 1 {
			return fmt.Sprintf("%s = %s || %s", VariableName(inst), x, y), nil
		}
		return fmt.Sprintf("%s = %s | %s", VariableName(inst), x, y), nil

	case *ir.InstPtrToInt:
		from, err := FormatValue(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", inst.To, err)
		}
		return fmt.Sprintf("%s = %s(uintptr(unsafe.Pointer(%s)))", VariableName(inst), to, from), nil

	case *ir.InstSDiv:
		x, err := FormatSigned(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatSigned(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 8 {
			return fmt.Sprintf("%s = byte(%s / %s)", VariableName(inst), x, y), nil
		}
		return fmt.Sprintf("%s = %s / %s", VariableName(inst), x, y), nil

	case *ir.InstUDiv:
		x, err := FormatUnsigned(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 8 {
			return fmt.Sprintf("%s = byte(%s / %s)", VariableName(inst), x, y), nil
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize > 8 {
			return fmt.Sprintf("%s = int%d(%s / %s)", VariableName(inst), goIntBits(intType.BitSize), x, y), nil
		}
		return fmt.Sprintf("%s = %s / %s", VariableName(inst), x, y), nil

	case *ir.InstSelect:
		cond, err := FormatValue(inst.Cond)
		if err != nil {
			return "", fmt.Errorf("error translating condition (%v): %w", inst.Cond, err)
		}
		valueTrue, err := FormatValue(inst.ValueTrue)
		if err != nil {
			return "", fmt.Errorf("error translating first operand (%v): %w", inst.ValueTrue, err)
		}
		valueFalse, err := FormatValue(inst.ValueFalse)
		if err != nil {
			return "", fmt.Errorf("error translating second operand (%v): %w", inst.ValueFalse, err)
		}
		name := VariableName(inst)
		return fmt.Sprintf("if %s { %s = %s } else { %s = %s }", cond, name, valueTrue, name, valueFalse), nil

	case *ir.InstSExt:
		toType, ok := inst.To.(*types.IntType)
		if !ok {
			return "", fmt.Errorf("%w: %T", errUnsupportedZextTo, inst.To)
		}
		from, err := FormatSigned(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		return fmt.Sprintf("%s = int%d(%s)", VariableName(inst), goIntBits(toType.BitSize), from), nil

	case *ir.InstShl:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s << %s", VariableName(inst), x, y), nil

	case *ir.InstShuffleVector:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		mask, err := FormatValue(inst.Mask)
		if err != nil {
			return "", fmt.Errorf("error translating mask (%v): %w", inst.Mask, err)
		}
		length := inst.Typ.Len
		return fmt.Sprintf("for i, m := range %s { if m < %d { %s[i] = %s[m] } else { %s[i] = %s[m - %d] } }", mask, length, VariableName(inst), x, VariableName(inst), y, length), nil

	case *ir.InstSIToFP:
		from, err := FormatSigned(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", inst.To, err)
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), to, from), nil

	case *ir.InstSRem:
		x, err := FormatSigned(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatSigned(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 8 {
			return fmt.Sprintf("%s = byte(%s %% %s)", VariableName(inst), x, y), nil
		}
		return fmt.Sprintf("%s = %s %% %s", VariableName(inst), x, y), nil

	case *ir.InstURem:
		x, err := FormatUnsigned(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 8 {
			return fmt.Sprintf("%s = byte(%s %% %s)", VariableName(inst), x, y), nil
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize > 8 {
			return fmt.Sprintf("%s = int%d(%s %% %s)", VariableName(inst), goIntBits(intType.BitSize), x, y), nil
		}
		return fmt.Sprintf("%s = %s %% %s", VariableName(inst), x, y), nil

	case *ir.InstStore:
		dest, err := FormatValue(inst.Dst)
		if err != nil {
			return "", fmt.Errorf("error translating destination (%v): %w", inst.Dst, err)
		}
		src, err := FormatValue(inst.Src)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", inst.Src, err)
		}
		if strings.HasPrefix(dest, "&") {
			return fmt.Sprintf("%s = %s", strings.TrimPrefix(dest, "&"), src), nil
		}
		return fmt.Sprintf("*%s = %s", dest, src), nil

	case *ir.InstSub:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s - %s", VariableName(inst), x, y), nil

	case *ir.InstTrunc:
		if vt, ok := inst.To.(*types.VectorType); ok {
			toType, ok := vt.ElemType.(*types.IntType)
			if !ok {
				return "", fmt.Errorf("%w: %v", errUnsupportedZextTo, inst.To)
			}
			to, err := TypeSpec(toType)
			if err != nil {
				return "", fmt.Errorf("error translating To type (%v): %w", toType, err)
			}
			from, err := FormatValue(inst.From)
			if err != nil {
				return "", fmt.Errorf("error translating source (%v): %w", inst.From, err)
			}
			return fmt.Sprintf("for i, v := range %s { %s[i] = %s(v) }", from, VariableName(inst), to), nil
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating To type (%v): %w", inst.To, err)
		}
		from, err := FormatValue(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		if intType, ok := inst.To.(*types.IntType); ok && intType.BitSize == 1 {
			// trunc … to i1 → bool (not byte(x & 1))
			return fmt.Sprintf("%s = (%s & 1) != 0", VariableName(inst), from), nil
		}
		if intType, ok := inst.To.(*types.IntType); ok && intType.BitSize < 8 {
			return fmt.Sprintf("%s = byte(%s & %d)", VariableName(inst), from, 255>>(8-intType.BitSize)), nil
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), to, from), nil

	case *ir.InstUIToFP:
		from, err := FormatUnsigned(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %w", inst.To, err)
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), to, from), nil

	case *ir.InstXor:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %w", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %w", inst.Y, err)
		}
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return fmt.Sprintf("for i, v := range %s { %s[i] = v ^ %s[i] }", x, VariableName(inst), y), nil
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 1 {
			return fmt.Sprintf("%s = %s != %s", VariableName(inst), x, y), nil
		}
		return fmt.Sprintf("%s = %s ^ %s", VariableName(inst), x, y), nil

	case *ir.InstZExt:
		if vt, ok := inst.To.(*types.VectorType); ok {
			toType, ok := vt.ElemType.(*types.IntType)
			if !ok {
				return "", fmt.Errorf("%w: %v", errUnsupportedZextTo, inst.To)
			}
			ft, ok := inst.From.Type().(*types.VectorType)
			if !ok {
				return "", fmt.Errorf("%w: %v and %v", errMismatchedZextTypes, inst.To, inst.From.Type())
			}
			fromType, ok := ft.ElemType.(*types.IntType)
			if !ok {
				return "", fmt.Errorf("%w: %v", errUnsupportedZextFrom, inst.From.Type())
			}
			from, err := FormatValue(inst.From)
			if err != nil {
				return "", fmt.Errorf("error translating source (%v): %w", inst.From, err)
			}
			tw, fw := goIntBits(toType.BitSize), goIntBits(fromType.BitSize)
			return fmt.Sprintf("for i, v := range %s { %s[i] = int%d(uint%d(uint%d(v))) }", from, VariableName(inst), tw, tw, fw), nil
		}
		toType, ok := inst.To.(*types.IntType)
		if !ok {
			return "", fmt.Errorf("%w: %T", errUnsupportedZextTo, inst.To)
		}
		from, err := FormatUnsigned(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %w", inst.From, err)
		}
		if fromType, ok := inst.From.Type().(*types.IntType); ok && fromType.BitSize == 1 {
			return fmt.Sprintf("if %s { %s = 1 } else { %s = 0 }", from, VariableName(inst), VariableName(inst)), nil
		}
		w := goIntBits(toType.BitSize)
		return fmt.Sprintf("%s = int%d(uint%d(%s))", VariableName(inst), w, w, from), nil

	default:
		return "", fmt.Errorf("%w: %T", errUnsupportedInstruction, inst)
	}
}

// atomicAddFunc returns sync/atomic Add* for the LLVM integer type.
func atomicAddFunc(t types.Type) (string, bool) {
	it, ok := t.(*types.IntType)
	if !ok {
		return "", false
	}
	switch goIntBits(it.BitSize) {
	case 32:
		return "atomic.AddInt32", true
	case 64:
		return "atomic.AddInt64", true
	default:
		return "", false
	}
}

// atomicSwapFunc returns sync/atomic Swap* for the LLVM integer type.
func atomicSwapFunc(t types.Type) (string, bool) {
	it, ok := t.(*types.IntType)
	if !ok {
		return "", false
	}
	switch goIntBits(it.BitSize) {
	case 32:
		return "atomic.SwapInt32", true
	case 64:
		return "atomic.SwapInt64", true
	default:
		return "", false
	}
}

// formatAggregateIndex builds a Go selector/index chain for extractvalue/insertvalue.
// Struct fields are F0, F1, … (see TypeDefinition); arrays use [i].
func formatAggregateIndex(base string, t types.Type, indices []uint64) (string, error) {
	expr := base
	cur := t
	for _, idx := range indices {
		switch ct := cur.(type) {
		case *types.StructType:
			if int(idx) >= len(ct.Fields) {
				return "", fmt.Errorf("%w: index %d (len %d)", errStructIndexRange, idx, len(ct.Fields))
			}
			expr = fmt.Sprintf("%s.F%d", expr, idx)
			cur = ct.Fields[idx]
		case *types.ArrayType:
			expr = fmt.Sprintf("%s[%d]", expr, idx)
			cur = ct.ElemType
		case *types.VectorType:
			expr = fmt.Sprintf("%s[%d]", expr, idx)
			cur = ct.ElemType
		default:
			return "", fmt.Errorf("%w: %T", errUnsupportedAggregate, cur)
		}
	}
	return expr, nil
}

var libraryFunctions = map[string]string{
	"fabs":             "math.Abs",
	"getchar":          "libc.Getchar",
	"leaven_va_arg":    "libc.VAArg",
	"llvm_fabs_f64":    "math.Abs",
	"llvm_fabs_f80":    "math.Abs",
	"llvm_pow_f64":     "math.Pow",
	"memchr":           "libc.Memchr",
	"memcmp":           "libc.Memcmp",
	"__memcpy_chk":     "libc.MemcpyChk",
	"memset_pattern16": "libc.MemsetPattern16",
	"__memset_chk":     "libc.MemsetChk",
	"printf":           "libc.Printf",
	"putc":             "libc.Putc",
	"putchar":          "libc.Putchar",
	"puts":             "libc.Puts",
	"scanf":            "libc.Scanf",
	"sqrt":             "math.Sqrt",
	"__strcat_chk":     "libc.StrcatChk",
	"strchr":           "libc.Strchr",
	"strcmp":           "libc.Strcmp",
	"strcpy":           "libc.Strcpy",
	"strcspn":          "libc.Strcspn",
	"strncat":          "libc.Strncat",
	"strncmp":          "libc.Strncmp",
	"strncpy":          "libc.Strncpy",
	"strrchr":          "libc.Strrchr",
	"strspn":           "libc.Strspn",
	"strstr":           "libc.Strstr",
}
