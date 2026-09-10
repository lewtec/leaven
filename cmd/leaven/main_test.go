package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/lewtec/lewkit/x/cmd"
)

func TestParseCLIFile(t *testing.T) {
	a, err := cmd.Parse[args]("--package", "foo", "in.ll")
	if err != nil {
		t.Fatal(err)
	}
	if got := a.Package.Value(); got != "foo" {
		t.Fatalf("package = %q, want foo", got)
	}
	if got := a.Input.Value(); got != "in.ll" {
		t.Fatalf("input = %q, want in.ll", got)
	}
}

func TestParseCLIDash(t *testing.T) {
	a, err := cmd.Parse[args]("-")
	if err != nil {
		t.Fatal(err)
	}
	if got := a.Input.Value(); got != "-" {
		t.Fatalf("input = %q, want -", got)
	}
}

func TestParseCLIVersionFlag(t *testing.T) {
	a, err := cmd.Parse[args]("--version")
	if err != nil {
		t.Fatal(err)
	}
	if !a.version.Value() {
		t.Fatal("version flag unset")
	}
}

func TestParseCLIPackageDefault(t *testing.T) {
	a, err := cmd.Parse[args]()
	if err != nil {
		t.Fatal(err)
	}
	if got := a.Package.Value(); got != "main" {
		t.Fatalf("package = %q, want main", got)
	}
}

func TestDescription(t *testing.T) {
	got := args{}.Description()
	for _, want := range []string{
		"Transpile LLVM IR to Go.",
		"With no file (or -), read LLVM IR from stdin and write Go to stdout.",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("Description() missing %q:\n%s", want, got)
		}
	}
}

func TestParseCLIUnknownFlag(t *testing.T) {
	_, err := cmd.Parse[args]("--nope")
	if !errors.Is(err, cmd.ErrUnknownFlag) {
		t.Fatalf("err = %v, want ErrUnknownFlag", err)
	}
}

func TestParseCLITwoFiles(t *testing.T) {
	_, err := cmd.Parse[args]("a.ll", "b.ll")
	if !errors.Is(err, cmd.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}
