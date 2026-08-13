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

// Rust OsString nest: {{{{i64, ptr, {}}, {}}, i64}} is 24 ABI bytes; Go sizeof is 40.
func TestLLVMTypeSizeNestedZST(t *testing.T) {
	empty := &types.StructType{Fields: nil}
	inner := &types.StructType{Fields: []types.Type{
		types.I64, types.I8Ptr, empty,
	}}
	mid := &types.StructType{Fields: []types.Type{inner, empty}}
	outer := &types.StructType{Fields: []types.Type{mid, types.I64}}
	sz, err := llvmTypeSize(outer)
	if err != nil {
		t.Fatal(err)
	}
	if sz != 24 {
		t.Fatalf("llvm size %d want 24", sz)
	}
}
