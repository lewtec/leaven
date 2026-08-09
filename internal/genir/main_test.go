package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClangMajor(t *testing.T) {
	got, err := clangMajor("clang version 22.1.8 (https://github.com/conda-forge/clangdev-feedstock abc)\nTarget: x86_64\n")
	if err != nil {
		t.Fatal(err)
	}
	if got != "22" {
		t.Fatalf("got %q", got)
	}
}

func TestRustcLLVMMajor(t *testing.T) {
	got, err := rustcLLVMMajor("release: 1.97.1\nLLVM version: 22.1.6\n")
	if err != nil {
		t.Fatal(err)
	}
	if got != "22" {
		t.Fatalf("got %q", got)
	}
}

func TestToolOverride(t *testing.T) {
	p := filepath.Join(t.TempDir(), "source.c")
	if err := os.WriteFile(p, []byte("/* leaven:tool=conda:clang@18.1.8 */\nint x;\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if got := toolOverride(p); got != "conda:clang@18.1.8" {
		t.Fatalf("got %q", got)
	}
}
