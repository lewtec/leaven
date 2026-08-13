package leaven

import (
	"strings"
	"testing"

	"github.com/lewtec/leaven/internal/llir/ir/types"
)

func TestTypeNameRustArrayGeneric(t *testing.T) {
	st := &types.StructType{}
	st.SetName("smallvec::SmallVec<[usize; 2]>")
	got := TypeName(st)
	if got == "" {
		t.Fatal("empty name")
	}
	for _, bad := range []byte{';', '[', ']', ' ', ':', '<', '>', ','} {
		if strings.ContainsRune(got, rune(bad)) {
			t.Fatalf("name %q still has %q", got, string(bad))
		}
	}
}
