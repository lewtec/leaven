// The libc package implements various functions from the C standard library
// in Go.
package libc

import (
	"os"
	"unsafe"
)

// GoString returns s converted from a C string to a Go string.
func GoString(s *byte) string {
	return string(unsafe.Slice(s, Strlen(s)))
}

// argvPin keeps Argv's C strings and the pointer table alive for main.
var argvPin []*byte

// Argv is C argv: pointer to a nil-terminated list of *byte (os.Args).
func Argv() unsafe.Pointer {
	args := os.Args
	ptrs := make([]*byte, len(args)+1)
	for i, s := range args {
		b := append([]byte(s), 0)
		ptrs[i] = &b[0]
	}
	argvPin = ptrs
	return unsafe.Pointer(&ptrs[0])
}

// AddPointer does C-style pointer addition: it multiplies offset by
// sizeof(*ptr) and adds it to ptr.
func AddPointer[T any](ptr *T, offset int) *T {
	return (*T)(unsafe.Add(unsafe.Pointer(ptr), offset*int(unsafe.Sizeof(*ptr))))
}
