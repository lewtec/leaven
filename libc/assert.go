package libc

import "fmt"

// AssertFail is glibc/clang __assert_fail(expr, file, line, function).
// Noreturn in C; panic in Go so control-flow after it typechecks.
func AssertFail(expr, file *byte, line int32, function *byte) {
	panic(fmt.Sprintf("assertion failed: %s (%s:%d in %s)",
		cStr(expr), cStr(file), line, cStr(function)))
}

// AssertRtn is Darwin __assert_rtn(func, file, line, expr).
func AssertRtn(function, file *byte, line int32, expr *byte) {
	AssertFail(expr, file, line, function)
}

func cStr(p *byte) string {
	if p == nil {
		return "<nil>"
	}
	return GoString(p)
}
