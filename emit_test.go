package leaven

import (
	"bytes"
	"math"
	"os"
	"strings"
	"testing"
	"unsafe"

	"github.com/dave/jennifer/jen"
	"github.com/lewtec/leaven/libc"
)

func renderQual(t *testing.T, c jen.Code) string {
	t.Helper()
	f := jen.NewFile("p")
	f.Var().Id("_").Add(c)
	var buf bytes.Buffer
	if err := f.Render(&buf); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func TestQual(t *testing.T) {
	got := renderQual(t, Qual[byte]())
	if !strings.Contains(got, "var _ byte") {
		t.Fatalf("byte: %s", got)
	}
	got = renderQual(t, Qual[unsafe.Pointer]())
	if !strings.Contains(got, "unsafe.Pointer") {
		t.Fatalf("unsafe.Pointer: %s", got)
	}
	got = renderQual(t, Qual[os.File]())
	if !strings.Contains(got, "os.File") {
		t.Fatalf("os.File: %s", got)
	}
	f := jen.NewFile("p")
	f.Var().Id("x").Op("=").Add(libcT(libc.As[byte], Qual[byte](), jen.Id("p")))
	var buf bytes.Buffer
	if err := f.Render(&buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "libc.As[byte](p)") {
		t.Fatalf("As[byte]: %s", buf.String())
	}
}

func TestSym(t *testing.T) {
	f := jen.NewFile("p")
	f.Var().Id("x").Op("=").Add(Sym(libc.Off).Call(jen.Id("p"), jen.Lit(1)))
	f.Var().Id("y").Op("=").Add(Sym(math.Abs).Call(jen.Lit(1.0)))
	var buf bytes.Buffer
	if err := f.Render(&buf); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !strings.Contains(got, "libc.Off(p, 1)") {
		t.Fatalf("Off: %s", got)
	}
	if !strings.Contains(got, "math.Abs(1") {
		t.Fatalf("Abs: %s", got)
	}
}
