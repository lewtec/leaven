package libc

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestRustPrintPrintlnI32(t *testing.T) {
	// println!("{n}") → template \xC0\x01\n\x00, one i32 arg.
	tmpl := [4]byte{0xC0, 0x01, '\n', 0x00}
	n := int32(7)
	var args [16]byte
	Store(Ptr(&args[0]), 0, Ptr(&n))
	Store(Ptr(&args[8]), 0, rustFmtFn(RustFmtI32))

	got := captureStdout(t, func() {
		RustPrint(Ptr(&tmpl[0]), Ptr(&args[0]))
	})
	if got != "7\n" {
		t.Fatalf("stdout = %q, want %q", got, "7\n")
	}
}

func TestRustPrintPrintlnUsize(t *testing.T) {
	tmpl := [4]byte{0xC0, 0x01, '\n', 0x00}
	n := uint64(6)
	var args [16]byte
	Store(Ptr(&args[0]), 0, Ptr(&n))
	Store(Ptr(&args[8]), 0, rustFmtFn(RustFmtUsize))

	got := captureStdout(t, func() {
		RustPrint(Ptr(&tmpl[0]), Ptr(&args[0]))
	})
	if got != "6\n" {
		t.Fatalf("stdout = %q, want %q", got, "6\n")
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = w
	fn()
	os.Stdout = old
	w.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	r.Close()
	return buf.String()
}
