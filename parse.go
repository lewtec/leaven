package leaven

import (
	"io"
	"os"

	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir"
)

// parseIRFile reads path, normalizes clang-14 quirks, and parses with llir.
func parseIRFile(path string) (*ir.Module, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseIR(path, f)
}

// parseIR reads LLVM IR from r (name is for parse diagnostics).
func parseIR(name string, r io.Reader) (*ir.Module, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return asm.ParseString(name, normalizeLLVMIR(string(b)))
}
