package leaven

import (
	"regexp"
	"strings"
)

// alignOnAtomic strips ", align N" that clang 13+ emits on atomicrmw / cmpxchg.
// The vendored v14 parser (llir 0.3.5) parses atomicrmw without the align operand.
var alignOnAtomic = regexp.MustCompile(`(?i)(atomicrmw\b.*?\b(?:seq_cst|acq_rel|acquire|release|monotonic|unordered))\s*,\s*align\s+\d+`)

// mustprogressAttr is a clang-14 C++ function attr unknown to llir 0.3.5.
var mustprogressAttr = regexp.MustCompile(`\bmustprogress\b\s*`)

// normalizeLLVMIR rewrites LLVM IR so the vendored llir parser accepts clang-14 output.
func normalizeLLVMIR(src string) string {
	src = mustprogressAttr.ReplaceAllString(src, "")
	// Per-line: only touch instructions that mention atomicrmw (keep struct align= alone).
	lines := strings.Split(src, "\n")
	for i, line := range lines {
		if strings.Contains(line, "atomicrmw") {
			lines[i] = alignOnAtomic.ReplaceAllString(line, "$1")
			// Also strip a bare ", align N" if still present after ordering.
			lines[i] = stripTrailingAlign(lines[i])
		}
	}
	return strings.Join(lines, "\n")
}

var trailingAlign = regexp.MustCompile(`,\s*align\s+\d+\s*$`)

func stripTrailingAlign(line string) string {
	return trailingAlign.ReplaceAllString(line, "")
}
