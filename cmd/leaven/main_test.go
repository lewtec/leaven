package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/lewtec/lewkit/x/cmd"
)

func TestParseCLIFile(t *testing.T) {
	app, err := cmd.Parse[cmd.App[args]]("--package", "foo", "in.ll")
	if err != nil {
		t.Fatal(err)
	}
	if got := app.Args.Package.Value(); got != "foo" {
		t.Fatalf("package = %q, want foo", got)
	}
	if got := app.Args.Input.Value(); got != "in.ll" {
		t.Fatalf("input = %q, want in.ll", got)
	}
}

func TestParseCLIDash(t *testing.T) {
	app, err := cmd.Parse[cmd.App[args]]("-")
	if err != nil {
		t.Fatal(err)
	}
	if got := app.Args.Input.Value(); got != "-" {
		t.Fatalf("input = %q, want -", got)
	}
}

func TestParseCLIPackageDefault(t *testing.T) {
	app, err := cmd.Parse[cmd.App[args]]()
	if err != nil {
		t.Fatal(err)
	}
	if got := app.Args.Package.Value(); got != "main" {
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

func TestParseCLITwoFiles(t *testing.T) {
	_, err := cmd.Parse[cmd.App[args]]("a.ll", "b.ll")
	if !errors.Is(err, cmd.ErrInvalidArgument) {
		t.Fatalf("err = %v, want ErrInvalidArgument", err)
	}
}
