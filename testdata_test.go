package leaven

import (
	"os"
	"path/filepath"
	"testing"
)

// testdataIR returns committed IR. testdata/*.ll is gitignored (integration build output).
func testdataIR(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "ir", "anon_dot", "input.ll"))
	if err != nil {
		t.Fatal(err)
	}
	return b
}
