package libc

import (
	"fmt"
	"os"
	"unsafe"
)

func Getchar() int32 {
	var buf [1]byte
	_, err := os.Stdin.Read(buf[:])
	if err != nil {
		return -1
	}
	return int32(buf[0])
}

func Putc(c int32, stream *os.File) int32 {
	_, err := stream.Write([]byte{byte(c)})
	if err != nil {
		return -1
	}
	return c
}

func Putchar(c int32) int32 {
	return Putc(c, os.Stdout)
}

// Fprintf is C fprintf(stream, format, ...).
func Fprintf(stream *os.File, format *byte, args ...any) int32 {
	f := fixPrintfFormat(format, args)
	n, err := fmt.Fprintf(stream, f, args...)
	if err != nil {
		return -1
	}
	return int32(n)
}

// Snprintf is C snprintf(buf, n, format, ...).
func Snprintf(buf *byte, n int64, format *byte, args ...any) int32 {
	f := fixPrintfFormat(format, args)
	s := fmt.Sprintf(f, args...)
	if n > 0 && buf != nil {
		dst := unsafe.Slice(buf, n)
		copyLen := len(s)
		if int64(copyLen) >= n {
			copyLen = int(n - 1)
			if copyLen < 0 {
				copyLen = 0
			}
		}
		copy(dst[:copyLen], s)
		dst[copyLen] = 0
	}
	// C returns the number of chars that would have been written (excl. NUL).
	return int32(len(s))
}

// Vsnprintf is C vsnprintf. ap points at system va_list storage
// (__va_list_tag); llvm.va_start stores &varargs at overflow_arg_area (offset 8).
func Vsnprintf(buf *byte, n int64, format *byte, ap *byte) int32 {
	var args []any
	if ap != nil {
		vlPtr := *(**[]interface{})(unsafe.Add(unsafe.Pointer(ap), 8))
		if vlPtr != nil {
			args = append(args, (*vlPtr)...)
		}
	}
	return Snprintf(buf, n, format, args...)
}

// Fdopen is C fdopen(fd, mode). Maps 0/1/2 to stdin/stdout/stderr;
// other FDs use os.NewFile. Mode is accepted but not fully applied.
func Fdopen(fd int32, mode *byte) *os.File {
	_ = mode
	switch fd {
	case 0:
		return os.Stdin
	case 1:
		return os.Stdout
	case 2:
		return os.Stderr
	}
	return os.NewFile(uintptr(fd), fmt.Sprintf("fd%d", fd))
}

// Fclose is C fclose. Does not close stdin/stdout/stderr.
func Fclose(f *os.File) int32 {
	if f == nil {
		return -1
	}
	if f == os.Stdin || f == os.Stdout || f == os.Stderr {
		return 0
	}
	if err := f.Close(); err != nil {
		return -1
	}
	return 0
}

// Abort is C abort(): no-return process abort.
func Abort() {
	panic("abort")
}
