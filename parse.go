package main

import (
	"os"

	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir"
)

// parseIRFile reads path, normalizes clang-14 quirks, and parses with llir.
func parseIRFile(path string) (*ir.Module, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return asm.ParseString(path, normalizeLLVMIR(string(b)))
}
