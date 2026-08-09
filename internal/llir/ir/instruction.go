package ir

import "github.com/lewtec/leaven/internal/llir/ir/value"

// === [ Instructions ] ========================================================

// Instruction is an LLVM IR instruction. All instructions (except store and
// fence) implement the value.Named interface and may thus be used directly as
// values.
//
// An Instruction has one of the following underlying types.
//
// Unary instructions
//
// https://llvm.org/docs/LangRef.html#unary-operations
//
//    *ir.InstFNeg   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstFNeg
//
// Binary instructions
//
// https://llvm.org/docs/LangRef.html#binary-operations
//
//    *ir.InstAdd    // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstAdd
//    *ir.InstFAdd   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstFAdd
//    *ir.InstSub    // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstSub
//    *ir.InstFSub   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstFSub
//    *ir.InstMul    // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstMul
//    *ir.InstFMul   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstFMul
//    *ir.InstUDiv   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstUDiv
//    *ir.InstSDiv   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstSDiv
//    *ir.InstFDiv   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstFDiv
//    *ir.InstURem   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstURem
//    *ir.InstSRem   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstSRem
//    *ir.InstFRem   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstFRem
//
// Bitwise instructions
//
// https://llvm.org/docs/LangRef.html#bitwise-binary-operations
//
//    *ir.InstShl    // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstShl
//    *ir.InstLShr   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstLShr
//    *ir.InstAShr   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstAShr
//    *ir.InstAnd    // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstAnd
//    *ir.InstOr     // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstOr
//    *ir.InstXor    // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstXor
//
// Vector instructions
//
// https://llvm.org/docs/LangRef.html#vector-operations
//
//    *ir.InstExtractElement   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstExtractElement
//    *ir.InstInsertElement    // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstInsertElement
//    *ir.InstShuffleVector    // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstShuffleVector
//
// Aggregate instructions
//
// https://llvm.org/docs/LangRef.html#aggregate-operations
//
//    *ir.InstExtractValue   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstExtractValue
//    *ir.InstInsertValue    // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstInsertValue
//
// Memory instructions
//
// https://llvm.org/docs/LangRef.html#memory-access-and-addressing-operations
//
//    *ir.InstAlloca          // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstAlloca
//    *ir.InstLoad            // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstLoad
//    *ir.InstStore           // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstStore
//    *ir.InstFence           // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstFence
//    *ir.InstCmpXchg         // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstCmpXchg
//    *ir.InstAtomicRMW       // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstAtomicRMW
//    *ir.InstGetElementPtr   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstGetElementPtr
//
// Conversion instructions
//
// https://llvm.org/docs/LangRef.html#conversion-operations
//
//    *ir.InstTrunc           // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstTrunc
//    *ir.InstZExt            // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstZExt
//    *ir.InstSExt            // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstSExt
//    *ir.InstFPTrunc         // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstFPTrunc
//    *ir.InstFPExt           // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstFPExt
//    *ir.InstFPToUI          // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstFPToUI
//    *ir.InstFPToSI          // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstFPToSI
//    *ir.InstUIToFP          // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstUIToFP
//    *ir.InstSIToFP          // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstSIToFP
//    *ir.InstPtrToInt        // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstPtrToInt
//    *ir.InstIntToPtr        // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstIntToPtr
//    *ir.InstBitCast         // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstBitCast
//    *ir.InstAddrSpaceCast   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstAddrSpaceCast
//
// Other instructions
//
// https://llvm.org/docs/LangRef.html#other-operations
//
//    *ir.InstICmp         // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstICmp
//    *ir.InstFCmp         // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstFCmp
//    *ir.InstPhi          // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstPhi
//    *ir.InstSelect       // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstSelect
//    *ir.InstFreeze       // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstFreeze
//    *ir.InstCall         // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstCall
//    *ir.InstVAArg        // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstVAArg
//    *ir.InstLandingPad   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstLandingPad
//    *ir.InstCatchPad     // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstCatchPad
//    *ir.InstCleanupPad   // https://pkg.go.dev/github.com/lewtec/leaven/internal/llir/ir#InstCleanupPad
type Instruction interface {
	LLStringer
	// isInstruction ensures that only instructions can be assigned to the
	// instruction.Instruction interface.
	isInstruction()
	// Instruction implements the value.User interface.
	value.User
}
