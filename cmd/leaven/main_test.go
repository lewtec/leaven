package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/lewtec/lewkit/x/cmd"
)

func TestParseCLIInputFlag(t *testing.T) {
	app, err := cmd.Parse[cmd.App[args]]("--package", "foo", "--input", "in.ll")
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
	app, err := cmd.Parse[cmd.App[args]]("--input", "-")
	if err != nil {
		t.Fatal(err)
	}
	if got := app.Args.Input.Value(); got != "-" {
		t.Fatalf("input = %q, want -", got)
	}
}

func TestParseCLIVersionCommand(t *testing.T) {
	app, err := cmd.Parse[cmd.App[args]]("version")
	if err != nil {
		t.Fatal(err)
	}
	if !app.WantVersion() {
		t.Fatal("WantVersion() = false")
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

func TestParseCLIUnknownFlag(t *testing.T) {
	_, err := cmd.Parse[cmd.App[args]]("--nope")
	if !errors.Is(err, cmd.ErrUnknownFlag) {
		t.Fatalf("err = %v, want ErrUnknownFlag", err)
	}
}
