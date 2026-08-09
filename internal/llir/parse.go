package llir

import (
	"github.com/lewtec/leaven/internal/llir/ir"
	v14asm "github.com/lewtec/leaven/internal/llir/v14/asm"
)

// ParseString parses LLVM IR assembly into the shared IR model.
// Today this is the v14 (typed-pointer) frontend only.
func ParseString(path, content string) (*ir.Module, error) {
	return v14asm.ParseString(path, content)
}
