package main

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/llir/llvm/asm"
)

// IR fixtures under testdata/ir/ — small LLVM modules that regressed or are
// awkward for full C→native differential tests.
func TestIRSanity(t *testing.T) {
	cases := []struct {
		file         string
		packageName  string
		contains     []string
		notContains  []string
	}{
		{
			file:        "anon_dot.ll",
			packageName: "main",
			contains: []string{
				"type anon_1 struct",
				"var g anon_1 = anon_1{1, 2}",
				"func get() byte",
				"v0 = g.F0",
			},
			notContains: []string{
				"anon.1", // illegal package-selector style type name
			},
		},
		{
			file:        "unreachable.ll",
			packageName: "main",
			contains: []string{
				"func die()",
				`panic("unreachable")`,
				"func f(x int32) int32",
				"goto if_then",
				"goto if_end",
			},
		},
		{
			file:        "anon_dot.ll",
			packageName: "tslib",
			contains: []string{
				"package tslib",
				"type anon_1 struct",
			},
		},
	}

	for _, tc := range cases {
		name := tc.file
		if tc.packageName != "main" {
			name = tc.file + "/pkg=" + tc.packageName
		}
		t.Run(name, func(t *testing.T) {
			path := filepath.Join("testdata", "ir", tc.file)
			m, err := asm.ParseFile(path)
			if err != nil {
				t.Fatalf("parse %s: %v", path, err)
			}
			var buf bytes.Buffer
			if err := compile(&buf, m, tc.packageName); err != nil {
				t.Fatalf("compile: %v", err)
			}
			got := buf.String()

			// Must be valid Go source.
			formatted, err := format.Source(buf.Bytes())
			if err != nil {
				t.Fatalf("go/format: %v\n---- generated ----\n%s", err, got)
			}
			_ = formatted

			for _, s := range tc.contains {
				if !strings.Contains(got, s) {
					t.Errorf("missing %q in:\n%s", s, got)
				}
			}
			for _, s := range tc.notContains {
				if strings.Contains(got, s) {
					t.Errorf("unexpected %q in:\n%s", s, got)
				}
			}
		})
	}
}

// Ensure every .ll under testdata/ir is wired into TestIRSanity (or at least
// parses + compiles). Catches fixtures added without a case.
func TestIRFixturesCompile(t *testing.T) {
	dir := filepath.Join("testdata", "ir")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".ll") {
			continue
		}
		name := e.Name()
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(dir, name)
			m, err := asm.ParseFile(path)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			var buf bytes.Buffer
			if err := compile(&buf, m, "main"); err != nil {
				t.Fatalf("compile: %v\n", err)
			}
			if _, err := format.Source(buf.Bytes()); err != nil {
				t.Fatalf("go/format: %v\n---- generated ----\n%s", err, buf.String())
			}
		})
	}
}
