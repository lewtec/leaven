package leaven

import (
	"regexp"
	"strings"
)

// alignOnAtomic strips ", align N" that clang 13+ emits on atomicrmw / cmpxchg.
// github.com/llir/ll (2021) parses atomicrmw without the align operand.
var alignOnAtomic = regexp.MustCompile(`(?i)(atomicrmw\b.*?\b(?:seq_cst|acq_rel|acquire|release|monotonic|unordered))\s*,\s*align\s+\d+`)

// normalizeLLVMIR rewrites LLVM IR so the vendored llir parser accepts clang-14 output.
func normalizeLLVMIR(src string) string {
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
