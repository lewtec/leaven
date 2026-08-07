package main

import (
	"bytes"
	"encoding/json"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/llir/llvm/asm"
)

// expectTable is the JSON table in each fixture folder (expect.json).
type expectTable struct {
	Cases []expectCase `json:"cases"`
}

type expectCase struct {
	Name        string   `json:"name"`
	Package     string   `json:"package"`
	Contains    []string `json:"contains"`
	NotContains []string `json:"not_contains"`
}

// TestIRSanity walks testdata/ir/<fixture>/{input.ll,expect.json}.
func TestIRSanity(t *testing.T) {
	root := filepath.Join("testdata", "ir")
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		fixture := e.Name()
		t.Run(fixture, func(t *testing.T) {
			dir := filepath.Join(root, fixture)
			llPath := filepath.Join(dir, "input.ll")
			tablePath := filepath.Join(dir, "expect.json")

			m, err := asm.ParseFile(llPath)
			if err != nil {
				t.Fatalf("parse %s: %v", llPath, err)
			}

			raw, err := os.ReadFile(tablePath)
			if err != nil {
				t.Fatalf("read %s: %v", tablePath, err)
			}
			var table expectTable
			if err := json.Unmarshal(raw, &table); err != nil {
				t.Fatalf("json %s: %v", tablePath, err)
			}
			if len(table.Cases) == 0 {
				t.Fatalf("%s: empty cases table", tablePath)
			}

			for _, tc := range table.Cases {
				tc := tc
				name := tc.Name
				if name == "" {
					name = tc.Package
				}
				if name == "" {
					name = "default"
				}
				t.Run(name, func(t *testing.T) {
					pkg := tc.Package
					if pkg == "" {
						pkg = "main"
					}
					var buf bytes.Buffer
					if err := compile(&buf, m, pkg); err != nil {
						t.Fatalf("compile package=%s: %v", pkg, err)
					}
					got := buf.String()

					if _, err := format.Source(buf.Bytes()); err != nil {
						t.Fatalf("go/format: %v\n---- generated ----\n%s", err, got)
					}

					for _, s := range tc.Contains {
						if !strings.Contains(got, s) {
							t.Errorf("missing %q in:\n%s", s, got)
						}
					}
					for _, s := range tc.NotContains {
						if strings.Contains(got, s) {
							t.Errorf("unexpected %q in:\n%s", s, got)
						}
					}
				})
			}
		})
	}
}

// TestIRFixturesLayout ensures each fixture folder has the required pair of files.
func TestIRFixturesLayout(t *testing.T) {
	root := filepath.Join("testdata", "ir")
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	var found int
	for _, e := range entries {
		if !e.IsDir() {
			// No loose files at the ir root.
			t.Errorf("unexpected file at testdata/ir/%s (use one folder per fixture)", e.Name())
			continue
		}
		found++
		dir := filepath.Join(root, e.Name())
		for _, need := range []string{"input.ll", "expect.json"} {
			if _, err := os.Stat(filepath.Join(dir, need)); err != nil {
				t.Errorf("%s: missing %s", e.Name(), need)
			}
		}
	}
	if found == 0 {
		t.Fatal("no fixture folders under testdata/ir")
	}
}
