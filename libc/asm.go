package libc

import "sync/atomic"

var asmMemFence uint32

// InlineAsm is LLVM `call asm`. The rustc compiler barrier
// (`""` / `~{memory}`) is an atomic fence so the call is not deleted.
// Any other template panics with the asm text.
func InlineAsm(asm, constraint string) {
	if asm == "" && (constraint == "~{memory}" || constraint == "") {
		atomic.AddUint32(&asmMemFence, 1)
		return
	}
	panic("unsatisfied inline asm: " + asm + " : " + constraint)
}
