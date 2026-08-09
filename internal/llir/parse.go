package llir

import (
	"regexp"

	"github.com/lewtec/leaven/internal/llir/ir"
	v14asm "github.com/lewtec/leaven/internal/llir/v14/asm"
	v22asm "github.com/lewtec/leaven/internal/llir/v22"
)

// opaquePtrTok matches a ptr type token, not %ptr / getelementptr / ptrtoint.
// Dot is an ident char (%add.ptr). Do not treat that as a ptr type.
var opaquePtrTok = regexp.MustCompile(`(?:^|[^%@!.A-Za-z0-9_])ptr(?:[^A-Za-z0-9_.]|$)`)

// ParseString parses LLVM IR assembly into the shared IR model.
// LLVM 15+ (opaque ptr) uses the v22 frontend; typed-pointer IR uses v14.
func ParseString(path, content string) (*ir.Module, error) {
	if opaquePtrTok.MatchString(content) {
		return v22asm.ParseString(path, content)
	}
	return v14asm.ParseString(path, content)
}
