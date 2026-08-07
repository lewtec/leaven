package libc

import (
	"fmt"
	"unsafe"
)

// AssertFail is glibc/clang __assert_fail(expr, file, line, function).
// Noreturn in C; panic in Go so control-flow after it typechecks.
func AssertFail(expr, file *byte, line int32, function *byte) {
	panic(fmt.Sprintf("assertion failed: %s (%s:%d in %s)",
		cStr(expr), cStr(file), line, cStr(function)))
}

func cStr(p *byte) string {
	if p == nil {
		return "<nil>"
	}
	return string(unsafe.Slice(p, Strlen(p)))
}
