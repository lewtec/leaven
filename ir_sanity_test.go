package leaven

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/llir/llvm/ir"
)

// expectTable is the JSON table in each fixture folder (expect.json).
type expectTable struct {
	Cases []expectCase `json:"cases"`
	// Run, if set, go-builds and runs the generated program as package main.
	Run *expectRun `json:"run"`
}

type expectCase struct {
	Name        string   `json:"name"`
	Package     string   `json:"package"`
	Contains    []string `json:"contains"`
	NotContains []string `json:"not_contains"`
}

// expectRun executes generated Go and checks stdout/stderr substrings.
type expectRun struct {
	StdoutContains []string `json:"stdout_contains"`
	StderrContains []string `json:"stderr_contains"`
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

			m, err := parseIRFile(llPath)
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
					if err := Compile(&buf, m, pkg); err != nil {
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

			if table.Run != nil {
				t.Run("run", func(t *testing.T) {
					runFixtureProgram(t, m, table.Run)
				})
			}
		})
	}
}

// runFixtureProgram generates package main, builds in a temp module with a
// replace to this repo, and checks output (catches panics / wrong runtime).
func runFixtureProgram(t *testing.T, m *ir.Module, run *expectRun) {
	t.Helper()
	var buf bytes.Buffer
	if err := Compile(&buf, m, "main"); err != nil {
		t.Fatalf("compile: %v", err)
	}
	src, err := format.Source(buf.Bytes())
	if err != nil {
		t.Fatalf("go/format: %v", err)
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), src, 0644); err != nil {
		t.Fatal(err)
	}
	// Module path of this checkout (for replace).
	modRoot, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	gomod := fmt.Sprintf("module leavenfixture\n\ngo 1.22\n\nrequire github.com/lewtec/leaven v0.0.0\n\nreplace github.com/lewtec/leaven => %s\n", modRoot)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0644); err != nil {
		t.Fatal(err)
	}

	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = dir
	if out, err := tidy.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, out)
	}
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run: %v\n---- output ----\n%s", err, out)
	}
	text := string(out)
	for _, s := range run.StdoutContains {
		if !strings.Contains(text, s) {
			t.Errorf("output missing %q\n---- output ----\n%s", s, text)
		}
	}
	for _, s := range run.StderrContains {
		if !strings.Contains(text, s) {
			t.Errorf("output missing stderr %q\n---- output ----\n%s", s, text)
		}
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
			// Allow a top-level README; everything else is a fixture folder.
			if e.Name() == "README.md" {
				continue
			}
			t.Errorf("unexpected file at testdata/ir/%s (use one folder per fixture)", e.Name())
			continue
		}
		found++
		dir := filepath.Join(root, e.Name())
		for _, need := range []string{"input.ll", "expect.json", "source.c"} {
			if _, err := os.Stat(filepath.Join(dir, need)); err != nil {
				t.Errorf("%s: missing %s", e.Name(), need)
			}
		}
	}
	if found == 0 {
		t.Fatal("no fixture folders under testdata/ir")
	}
}
