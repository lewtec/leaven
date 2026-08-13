package leaven

import (
	"bytes"
	"go/format"
	"strings"
	"testing"
)

// C++ D1/C1 aliases are often forward-referenced from vtables before the
// alias line. The parser plants an i8 @D1 global; emit must still put a
// real function pointer in the vtable (not &var byte).
func TestVtableAliasIsFunctionPointer(t *testing.T) {
	src := `
@_ZTV1X = constant { [3 x ptr] } { [3 x ptr] [
  ptr null,
  ptr null,
  ptr @_ZN1XD1Ev
] }
define void @_ZN1XD2Ev(ptr %this) {
  ret void
}
@_ZN1XD1Ev = alias void (ptr), ptr @_ZN1XD2Ev
`
	m, err := parseIR("alias.ll", strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := Compile(&buf, m, "p"); err != nil {
		t.Fatal(err)
	}
	out, err := format.Source(buf.Bytes())
	if err != nil {
		t.Fatalf("gofmt: %v\n%s", err, buf.Bytes())
	}
	s := string(out)
	if strings.Contains(s, "var _ZN1XD1Ev byte") {
		t.Fatal("alias still emitted as empty byte global")
	}
	if !strings.Contains(s, "_ZN1XD2Ev") {
		t.Fatal("aliasee D2 missing")
	}
	// Vtable slot should reference the function, not address-of a byte var.
	if strings.Contains(s, "unsafe.Pointer(&_ZN1XD1Ev)") {
		t.Fatalf("vtable still takes address of D1 stub:\n%s", s)
	}
}
