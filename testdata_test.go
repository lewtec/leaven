package leaven

import (
	"os"
	"path/filepath"
	"testing"
)

// testdataIR returns committed IR. testdata/*.ll (integration live build) is gitignored.
func testdataIR(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "ir", "c_anon_dot", "input.14.ll"))
	if err != nil {
		t.Fatal(err)
	}
	return b
}
