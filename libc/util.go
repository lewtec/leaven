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
func Addr[T any](p *T) uintptr {
	return uintptr(unsafe.Pointer(p))
}

// PtrBits is an LLVM ptr memory slot (uint64) from a Go pointer. nil-safe.
func PtrBits(p unsafe.Pointer) uint64 { return uint64(uintptr(p)) }

// PtrFromBits is a Go pointer from an LLVM ptr memory slot.
func PtrFromBits(u uint64) unsafe.Pointer { return unsafe.Pointer(uintptr(u)) }

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

// funcWord is a Go func value: code pointer, then optional data word.
type funcWord struct {
	code, data unsafe.Pointer
}

// FuncCode is the code word of a Go func. Inverse of FuncFromCode.
// Used for C++ atexit slots and rustc vtable method pointers.
func FuncCode[F any](fn F) unsafe.Pointer {
	return Load[unsafe.Pointer](Ptr(&fn), 0)
}

// FuncFromCode rebuilds a Go func from a code word (nil closure data).
func FuncFromCode[F any](code unsafe.Pointer) F {
	var w funcWord
	w.code = code
	return Load[F](Ptr(&w), 0)
}

// GoString returns s converted from a C string to a Go string.
func GoString(s *byte) string {
	if s == nil {
		return ""
	}
	return string(Bytes(s, int(Strlen(s))))
}

// argvPin keeps Argv's C strings alive for main.
// argvSlots is the IR-width table (8-byte ptr slots, not Go *byte).
var (
	argvPin   []*byte
	argvSlots []uint64
)

// Argv is C argv: pointer to a nil-terminated list of *byte (os.Args).
// Slots are uint64 so x86_64 GEP *8 works on 386.
func Argv() unsafe.Pointer {
	args := os.Args
	pin := make([]*byte, len(args))
	slots := make([]uint64, len(args)+1)
	for i, s := range args {
		b := append([]byte(s), 0)
		pin[i] = &b[0]
		slots[i] = PtrBits(Ptr(&b[0]))
	}
	argvPin = pin
	argvSlots = slots
	return unsafe.Pointer(&slots[0])
}
