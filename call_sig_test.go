package leaven

import (
	"bytes"
	"strings"
	"testing"
)

// select of @f before define used to cache void()*; invoke then cast to func().
func TestInvokeThroughSelectFnptr(t *testing.T) {
	src := `
define void @f(i1 %c, ptr %p) {
  %s = select i1 %c, ptr @pure_a, ptr @pure_b
  %r = invoke i1 %s(ptr %p) to label %n unwind label %e
n:
  ret void
e:
  %lp = landingpad {ptr, i32} cleanup
  resume {ptr, i32} %lp
}
define i1 @pure_a(ptr %x) {
  ret i1 true
}
define i1 @pure_b(ptr %x) {
  ret i1 false
}
`
	m, err := parseIR("t.ll", strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := Compile(&buf, m, "main"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "func() func() {") {
		t.Fatalf("stale void() fnptr cast:\n%s", out)
	}
	if !strings.Contains(out, "func(unsafe.Pointer) bool") {
		t.Fatalf("missing correct cast:\n%s", out)
	}
}

// Local SSA %write / %read must not become libc.Write / libc.Read.
func TestLocalNameNotLibraryRef(t *testing.T) {
	src := `
define void @record(i1 %write) {
  %read = alloca i8, align 1
  %z = zext i1 %write to i8
  store i8 %z, ptr %read, align 1
  %v = load i8, ptr %read, align 1
  ret void
}
declare i64 @read(i32, ptr, i64)
declare i64 @write(i32, ptr, i64)
`
	m, err := parseIR("t.ll", strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := Compile(&buf, m, "main"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "libc.Write") || strings.Contains(out, "libc.Read") {
		t.Fatalf("local write/read remapped to libc:\n%s", out)
	}
	if !strings.Contains(out, "write") {
		t.Fatalf("missing local write:\n%s", out)
	}
}

func TestPackedBitIteratorSize(t *testing.T) {
	// <{ ptr, i32 }> must be 12 bytes in Go so vector<bool> layout matches.
	src := `
%iter = type <{ ptr, i32 }>
%it = type { %iter, [4 x i8] }
@g = global %it zeroinitializer
`
	m, err := parseIR("t.ll", strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := Compile(&buf, m, "main"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	// Packed struct is a byte blob so the parent stays 16 bytes (12+4).
	if !strings.Contains(out, "[12]byte") {
		t.Fatalf("expected [12]byte packed iter:\n%s", out)
	}
}

func TestPackedTreeNodeParentLayout(t *testing.T) {
	src := `
%end = type { ptr }
%base = type <{ %end, ptr, ptr, i8 }>
%uni = type { { ptr, i32 } }
%node = type { %base, [7 x i8], %uni }
@g = global %node zeroinitializer
`
	m, err := parseIR("t.ll", strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := Compile(&buf, m, "main"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "[25]byte") {
		t.Fatalf("packed base should be [25]byte:\n%s", out)
	}
}

func TestTreeNodeValueGEPOffset(t *testing.T) {
	// libc++ __tree_node: packed 25-byte base + 7 pad + union at LLVM +32.
	// Go pads the base to 32 so .F2 would be +40.
	src := `
%end = type { ptr }
%base = type <{ %end, ptr, ptr, i8 }>
%uni = type { { ptr, i32 } }
%node = type { %base, [7 x i8], %uni }
define ptr @getv(ptr %p) {
  %v = getelementptr inbounds nuw %node, ptr %p, i32 0, i32 2
  ret ptr %v
}
`
	m, err := parseIR("t.ll", strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := Compile(&buf, m, "main"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, ".F2") {
		t.Fatalf("value GEP used Go .F2 (offset 40):\n%s", out)
	}
	if !strings.Contains(out, "32") {
		t.Fatalf("value GEP missing LLVM offset 32:\n%s", out)
	}
}
