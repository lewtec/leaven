// The libc package implements various functions from the C standard library
// in Go.
package libc

import (
	"os"
	"unsafe"
)

// Pointer kit (byte offset vs element index):
//
//	Ptr / As / Addr     (void *)p  /  (T *)p  /  (uintptr)p
//	Off(p, n)           (char *)p + n
//	AddPointer(p, i)    p + i          (scaled by sizeof *p)
//	Load / Store        *(T *)((char *)p + off)
//	Bytes(p, n)         (char *, n) as []byte
//	GoString            NUL-terminated C string

// Ptr is (void *)p.
func Ptr[T any](p *T) unsafe.Pointer {
	return unsafe.Pointer(p)
}

// As is (T *)p.
func As[T any](p unsafe.Pointer) *T {
	return (*T)(p)
}

// Addr is (uintptr)p. Map keys, nil-safe.
func Addr(p *byte) uintptr {
	return uintptr(unsafe.Pointer(p))
}

// Off is (char *)p + n.
func Off(p unsafe.Pointer, n int) unsafe.Pointer {
	return unsafe.Add(p, n)
}

// AddPointer is p + i (C pointer arithmetic, scaled).
func AddPointer[T any](ptr *T, offset int) *T {
	return (*T)(unsafe.Add(unsafe.Pointer(ptr), offset*int(unsafe.Sizeof(*ptr))))
}

// Load is *(T *)((char *)p + off).
func Load[T any](p unsafe.Pointer, off int) T {
	return *(*T)(unsafe.Add(p, off))
}

// Store is *(T *)((char *)p + off) = v.
func Store[T any](p unsafe.Pointer, off int, v T) {
	*(*T)(unsafe.Add(p, off)) = v
}

// Bytes is a C buffer (p, n). n<=0 yields nil.
func Bytes(p *byte, n int) []byte {
	if n <= 0 {
		return nil
	}
	return unsafe.Slice(p, n)
}

// GoString returns s converted from a C string to a Go string.
func GoString(s *byte) string {
	if s == nil {
		return ""
	}
	return string(Bytes(s, int(Strlen(s))))
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
