package leaven

import (
	"bytes"
	"os/exec"
	"testing"
)

// doTestCase runs the specified program in the testdata directory twice,
// once compiled directly with clang, and the other time compiled to Go with
// leaven. It compares the output of the two programs.
func doTestCase(t *testing.T, progName string) {
	p := "testdata/" + progName
	clang := exec.Command("clang", "-o", p+"_c", p+".c", "testdata/main.c", "-lm")
	if out, err := clang.CombinedOutput(); err != nil {
		t.Fatalf("Error in native compilation: %v\n%s", err, out)
	}

	prog := exec.Command(p + "_c")
	nativeOut, err := prog.CombinedOutput()
	if err != nil {
		t.Fatalf("Error running natively-compiled program: %v", err)
	}

	clang2 := exec.Command("clang", "-S", "-emit-llvm", "-fno-builtin", "-nostdinc", "-Iinclude", "-o", p+".ll", p+".c")
	if out, err := clang2.CombinedOutput(); err != nil {
		t.Fatalf("Error compiling to LLVM: %v\n%s", err, out)
	}

	leaven := exec.Command("go", "run", "./cmd/leaven", p+".ll")
	if out, err := leaven.CombinedOutput(); err != nil {
		t.Fatalf("Error running leaven: %v\n%s", err, out)
	}

	goimports := exec.Command("goimports", "-w", p+".go")
	if out, err := goimports.CombinedOutput(); err != nil {
		t.Fatalf("Error running goimports: %v\n%s", err, out)
	}

	goRun := exec.Command("go", "run", p+".go", "testdata/main.go")
	goOut, err := goRun.CombinedOutput()
	if err != nil {
		t.Fatalf("Error running Go program: %v", err)
	}

	if !bytes.Equal(goOut, nativeOut) {
		t.Fatalf("Output does not match. C = %q, Go = %q", nativeOut, goOut)
	}
}

func TestHelloWorld(t *testing.T) {
	doTestCase(t, "hello")
}

func TestHelloWorldPuts(t *testing.T) {
	doTestCase(t, "hello-puts")
}

func TestBinaryTrees(t *testing.T) {
	doTestCase(t, "binarytrees")
}

func TestFannkuch(t *testing.T) {
	doTestCase(t, "fannkuch-redux")
}

func TestNBody(t *testing.T) {
	doTestCase(t, "nbody")
}

func TestSpectralNorm(t *testing.T) {
	doTestCase(t, "spectral-norm")
}

func TestMandelbrot(t *testing.T) {
	doTestCase(t, "mandelbrot")
}
