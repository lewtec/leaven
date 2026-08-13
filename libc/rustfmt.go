package libc

import (
	"fmt"
	"os"
	"unsafe"
)

// rustArg is core::fmt::rt::Argument: { value, formatter }.
// rustc 1.97 stores a 16-byte placeholder (no tag) for Display args.
const rustArgSize = 16

type rustFmtFn func(val, formatter unsafe.Pointer) bool

// write_str method: fn(*mut W, *const u8, usize) -> bool (Ok=false).
type rustWriteStrFn func(w, data unsafe.Pointer, n int64) bool

func rustStdoutWriteStr(_ unsafe.Pointer, data unsafe.Pointer, n int64) bool {
	if data != nil && n > 0 {
		os.Stdout.Write(Bytes(As[byte](data), int(n)))
	}
	return false
}

// Minimal Formatter for RustPrint: flags=0 so pad takes the write_str path.
// Layout: { data:ptr, vtable:ptr, flags:i32, … } — first 24 bytes matter.
var (
	rustPrintVtable struct {
		drop                  unsafe.Pointer
		size, align           uint64
		writeStr, writeChar, writeFmt unsafe.Pointer
	}
	rustPrintFmt struct {
		data   unsafe.Pointer
		vtable unsafe.Pointer
		flags  int32
		_pad   [4]byte
	}
	rustPrintFmtInit bool
)

func rustPrintFormatter() unsafe.Pointer {
	if !rustPrintFmtInit {
		rustPrintVtable.size = 8
		rustPrintVtable.align = 8
		rustPrintVtable.writeStr = FuncCode(rustStdoutWriteStr)
		rustPrintFmt.vtable = Ptr(&rustPrintVtable)
		rustPrintFmtInit = true
	}
	return Ptr(&rustPrintFmt)
}

// RustPrint is std::io::stdio::_print(template, args).
// Template encoding: rustc 1.97 fmt::Arguments (byte pieces + 0xC0 placeholders).
func RustPrint(tmpl, args unsafe.Pointer) {
	if tmpl == nil {
		return
	}
	p := tmpl
	next := 0
	formatter := rustPrintFormatter()
	for {
		n := Load[byte](p, 0)
		p = Off(p, 1)
		if n == 0 {
			return
		}
		if n < 128 {
			os.Stdout.Write(Bytes(As[byte](p), int(n)))
			p = Off(p, int(n))
			continue
		}
		if n == 128 {
			lenb := Bytes(As[byte](p), 2)
			nlen := int(lenb[0]) | int(lenb[1])<<8
			p = Off(p, 2)
			os.Stdout.Write(Bytes(As[byte](p), nlen))
			p = Off(p, nlen)
			continue
		}
		if n < 0xC0 {
			panic("invalid rust fmt template")
		}
		skip := 0
		if n&1 != 0 {
			skip += 4
		}
		if n&2 != 0 {
			skip += 2
		}
		if n&4 != 0 {
			skip += 2
		}
		idx := next
		if n&8 != 0 {
			ib := Bytes(As[byte](Off(p, skip)), 2)
			idx = int(ib[0]) | int(ib[1])<<8
			skip += 2
		} else {
			next++
		}
		p = Off(p, skip)
		slot := Off(args, idx*rustArgSize)
		val := Load[unsafe.Pointer](slot, 0)
		fn := Load[rustFmtFn](slot, 8)
		fn(val, formatter)
	}
}

// RustFmtI32 is <i32 as core::fmt::Display>::fmt. Ok is false (i1 0).
func RustFmtI32(val, _ unsafe.Pointer) bool {
	fmt.Fprint(os.Stdout, *(*int32)(val))
	return false
}

// RustFmtUsize is <usize as core::fmt::Display>::fmt.
func RustFmtUsize(val, _ unsafe.Pointer) bool {
	fmt.Fprint(os.Stdout, *(*uint64)(val))
	return false
}
