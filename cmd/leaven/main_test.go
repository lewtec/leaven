package main

import (
	"errors"
	"testing"

	"github.com/lewtec/lewkit/x/cmd"
)

func TestParseCLIFile(t *testing.T) {
	app, input, err := parseCLI([]string{"--package", "foo", "in.ll"})
	if err != nil {
		t.Fatal(err)
	}
	if input != "in.ll" {
		t.Fatalf("input = %q, want in.ll", input)
	}
	if got := app.Args.Package.Value(); got != "foo" {
		t.Fatalf("package = %q, want foo", got)
	}
	app.Args.file = input
	if got := app.Args.path(); got != "in.ll" {
		t.Fatalf("path() = %q, want in.ll", got)
	}
}

func TestParseCLIInputFlag(t *testing.T) {
	app, input, err := parseCLI([]string{"--input", "in.ll"})
	if err != nil {
		t.Fatal(err)
	}
	if input != "" {
		t.Fatalf("positional = %q, want empty", input)
	}
	if got := app.Args.Input.Value(); got != "in.ll" {
		t.Fatalf("input flag = %q, want in.ll", got)
	}
	if got := app.Args.path(); got != "in.ll" {
		t.Fatalf("path() = %q, want in.ll", got)
	}
}

func TestParseCLIDash(t *testing.T) {
	_, input, err := parseCLI([]string{"-"})
	if err != nil {
		t.Fatal(err)
	}
	if input != "-" {
		t.Fatalf("input = %q, want -", input)
	}
}

func TestParseCLIVersionCommand(t *testing.T) {
	app, input, err := parseCLI([]string{"version"})
	if err != nil {
		t.Fatal(err)
	}
	if input != "" {
		t.Fatalf("input = %q, want empty", input)
	}
	if !app.WantVersion() {
		t.Fatal("WantVersion() = false")
	}
}

func TestParseCLIUnknownFlag(t *testing.T) {
	_, _, err := parseCLI([]string{"--nope"})
	if !errors.Is(err, cmd.ErrUnknownFlag) {
		t.Fatalf("err = %v, want ErrUnknownFlag", err)
	}
}

func TestParseCLITwoFiles(t *testing.T) {
	_, _, err := parseCLI([]string{"a.ll", "b.ll"})
	if !errors.Is(err, cmd.ErrUnknownCommand) {
		t.Fatalf("err = %v, want ErrUnknownCommand", err)
	}
}
