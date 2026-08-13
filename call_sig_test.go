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
	// Expect [8]byte for packed ptr slot, not uintptr (align 8 → size 16).
	if !strings.Contains(out, "[8]byte") {
		t.Fatalf("expected [8]byte packed ptr field:\n%s", out)
	}
}
