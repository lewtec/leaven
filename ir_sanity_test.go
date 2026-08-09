package leaven

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/lewtec/leaven/internal/llir/ir"
)

// expectTable is the JSON table in each fixture folder (expect.json).
type expectTable struct {
	Cases []expectCase `json:"cases"`
	// Run, if set, go-builds and runs the generated program as package main.
	Run *expectRun `json:"run"`
	// Parse, if false, skips llir (e.g. rustc LLVM 22 opaque ptr).
	Parse *bool `json:"parse"`
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

// TestIRSanity walks testdata/ir/<fixture>/{input.<n>.ll,expect.json}.
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
			tablePath := filepath.Join(dir, "expect.json")

			raw, err := os.ReadFile(tablePath)
			if err != nil {
				t.Fatalf("read %s: %v", tablePath, err)
			}
			var table expectTable
			if err := json.Unmarshal(raw, &table); err != nil {
				t.Fatalf("json %s: %v", tablePath, err)
			}
			if table.Parse != nil && !*table.Parse {
				return
			}
			if len(table.Cases) == 0 {
				t.Fatalf("%s: empty cases table", tablePath)
			}

			lls, err := fixtureIRFiles(dir)
			if err != nil {
				t.Fatal(err)
			}
			for _, llPath := range lls {
				llPath := llPath
				ver := irFileVersion(llPath)
				t.Run("ir"+ver, func(t *testing.T) {
					m, err := parseIRFile(llPath)
					if err != nil {
						t.Fatalf("parse %s: %v", llPath, err)
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
		if _, err := os.Stat(filepath.Join(dir, "expect.json")); err != nil {
			t.Errorf("%s: missing expect.json", e.Name())
		}
		src, err := fixtureSource(dir)
		if err != nil {
			t.Errorf("%s: %v", e.Name(), err)
		} else {
			prefix, err := fixtureLangPrefix(src)
			if err != nil {
				t.Errorf("%s: %v", e.Name(), err)
			} else if !strings.HasPrefix(e.Name(), prefix) {
				t.Errorf("%s: dir must start with %s (have %s)", e.Name(), prefix, filepath.Base(src))
			}
		}
		if _, err := fixtureIRFiles(dir); err != nil {
			t.Errorf("%s: %v", e.Name(), err)
		}
	}
	if found == 0 {
		t.Fatal("no fixture folders under testdata/ir")
	}
}

const irHandMarker = "leaven:hand-ir"

var (
	errFixtureSourceCount = errors.New("need exactly one of source.c, source.cpp, source.rs")
	errFixtureIRMissing   = errors.New("need at least one input.<LLVM-major>.ll")
	irFileName            = regexp.MustCompile(`^input\.([0-9]+)\.ll$`)
)

func fixtureIRFiles(dir string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "input.*.ll"))
	if err != nil {
		return nil, err
	}
	var out []string
	for _, p := range matches {
		if irFileName.MatchString(filepath.Base(p)) {
			out = append(out, p)
		}
	}
	sort.Strings(out)
	if len(out) == 0 {
		return nil, errFixtureIRMissing
	}
	return out, nil
}

func irFileVersion(path string) string {
	m := irFileName.FindStringSubmatch(filepath.Base(path))
	if len(m) != 2 {
		return "?"
	}
	return m[1]
}

func fixtureSource(dir string) (string, error) {
	var hits []string
	for _, name := range []string{"source.c", "source.cpp", "source.rs"} {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err == nil {
			hits = append(hits, p)
		}
	}
	if len(hits) != 1 {
		return "", fmt.Errorf("%w (found %d)", errFixtureSourceCount, len(hits))
	}
	return hits[0], nil
}

func fixtureLangPrefix(src string) (string, error) {
	switch filepath.Ext(src) {
	case ".c":
		return "c_", nil
	case ".cpp":
		return "cpp_", nil
	case ".rs":
		return "rust_", nil
	default:
		return "", fmt.Errorf("%w: %s", errFixtureSourceCount, src)
	}
}

func fixtureHandIR(t *testing.T, sourcePath string) bool {
	t.Helper()
	b, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.Contains(b, []byte(irHandMarker))
}

// TestIRFixturesProducer requires generated input.<n>.ll to name clang or rustc.
func TestIRFixturesProducer(t *testing.T) {
	root := filepath.Join("testdata", "ir")
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(root, e.Name())
		src, err := fixtureSource(dir)
		if err != nil {
			t.Errorf("%s: %v", e.Name(), err)
			continue
		}
		if fixtureHandIR(t, src) {
			continue
		}
		lls, err := fixtureIRFiles(dir)
		if err != nil {
			t.Errorf("%s: %v", e.Name(), err)
			continue
		}
		for _, p := range lls {
			ll, err := os.ReadFile(p)
			if err != nil {
				t.Errorf("%s: %v", e.Name(), err)
				continue
			}
			want := "clang version "
			if strings.HasSuffix(src, ".rs") {
				want = "rustc version "
			} else {
				want = "clang version " + irFileVersion(p) + "."
			}
			if !bytes.Contains(ll, []byte(want)) {
				t.Errorf("%s: missing %q (run mise run ir:gen)", filepath.Base(p), want)
			}
		}
	}
}
