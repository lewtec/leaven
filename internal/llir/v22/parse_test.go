package v22

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/lewtec/leaven/internal/llir/ir/types"
)

func TestParseStrlen22(t *testing.T) {
	root := findRepoRoot(t)
	b, err := os.ReadFile(filepath.Join(root, "testdata/ir/c_strlen_map/input.22.ll"))
	if err != nil {
		t.Fatal(err)
	}
	m, err := ParseString("input.22.ll", string(b))
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Funcs) == 0 {
		t.Fatal("no funcs")
	}
	n := m.Funcs[0]
	if len(n.Params) != 1 {
		t.Fatalf("params %d", len(n.Params))
	}
	pt, ok := n.Params[0].Typ.(*types.PointerType)
	if !ok || !pt.IsOpaque() {
		t.Fatalf("param type %v opaque=%v", n.Params[0].Typ, ok && pt.IsOpaque())
	}
}

func TestParseRustAdd22(t *testing.T) {
	root := findRepoRoot(t)
	b, err := os.ReadFile(filepath.Join(root, "testdata/ir/rust_add/input.22.ll"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseString("input.22.ll", string(b)); err != nil {
		t.Fatal(err)
	}
}

func TestParseAll22Fixtures(t *testing.T) {
	root := findRepoRoot(t)
	dir := filepath.Join(root, "testdata", "ir")
	ents, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range ents {
		if !e.IsDir() {
			continue
		}
		path := filepath.Join(dir, e.Name(), "input.22.ll")
		b, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			t.Fatal(err)
		}
		t.Run(e.Name(), func(t *testing.T) {
			m, err := ParseString(path, string(b))
			if err != nil {
				t.Fatal(err)
			}
			if len(m.Funcs) == 0 {
				t.Fatal("no functions")
			}
			for _, f := range m.Funcs {
				if len(f.Blocks) == 0 {
					continue
				}
				for _, b := range f.Blocks {
					if b.Term == nil {
						t.Fatalf("%s: block %s has no terminator", f.Name(), b.Name())
					}
				}
			}
		})
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("no go.mod")
	return ""
}
