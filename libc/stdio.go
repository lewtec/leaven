package libc

import (
	"fmt"
	"os"
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

// Abort is C abort(): no-return process abort.
func Abort() {
	panic("abort")
}
