package leaven

import (
	"testing"

	"github.com/lewtec/leaven/internal/llir/ir"
	"github.com/lewtec/leaven/internal/llir/ir/types"
)

func TestParseCxxOstreamOps(t *testing.T) {
	cases := []struct {
		name  string
		ident string
		recv  string
		arg   string
		ptr   bool
	}{
		{
			"_ZNSt3__1lsB9nqn220108INS_11char_traitsIcEEEERNS_13basic_ostreamIcT_EES6_c",
			"<<", "", "char", false,
		},
		{
			"_ZNSt3__113basic_ostreamIcNS_11char_traitsIcEEElsEm",
			"<<", "basic_ostream", "unsigned long", false,
		},
		{
			"_ZNSt3__1lsB9nqn220108IcNS_11char_traitsIcEEEERNS_13basic_ostreamIT_T0_EES7_PKc",
			"<<", "", "char", true,
		},
		{
			"_ZNSt3__113basic_ostreamIcNS_11char_traitsIcEEElsB9nqn220108ERKNS_12basic_stringIcS2_NS_9allocatorIcEEEE",
			"<<", "basic_ostream", "basic_string", false,
		},
		{
			"_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_",
			"endl", "", "", false,
		},
		{
			"_ZNSolsEPFRSoS_E",
			"<<", "ostream", "fnptr", true,
		},
		{
			"_ZNSirsERi",
			">>", "istream", "int", false,
		},
	}
	for _, tc := range cases {
		n, ok := parseCxx(tc.name)
		if !ok {
			t.Fatalf("parse %s", tc.name)
		}
		if n.ident != tc.ident || n.recv != tc.recv {
			t.Fatalf("%s: ident=%q recv=%q", tc.name, n.ident, n.recv)
		}
		if tc.arg == "" {
			continue
		}
		arg, ok := n.streamArg()
		if !ok {
			if n.ident == ">>" {
				if len(n.args) != 1 || n.args[0].ident != tc.arg {
					t.Fatalf("%s: >> arg=%v", tc.name, n.args)
				}
				continue
			}
			t.Fatalf("%s: no stream arg", tc.name)
		}
		if arg.ident != tc.arg || arg.ptr != tc.ptr {
			t.Fatalf("%s: arg ident=%q ptr=%v", tc.name, arg.ident, arg.ptr)
		}
	}
}

func TestParseCxxCtorStr(t *testing.T) {
	ctor := "_ZNSt3__119basic_ostringstreamIcNS_11char_traitsIcEENS_9allocatorIcEEEC1B9nqn220108Ev"
	n, ok := parseCxx(ctor)
	if !ok || !n.ctor || n.class() != "basic_ostringstream" || len(n.args) != 0 {
		t.Fatalf("ctor %+v ok=%v", n, ok)
	}
	str := "_ZNKRSt3__119basic_ostringstreamIcNS_11char_traitsIcEENS_9allocatorIcEEE3strB9nqn220108Ev"
	n, ok = parseCxx(str)
	if !ok || n.ident != "str" || n.hasString() || n.class() != "basic_ostringstream" {
		t.Fatalf("str %+v ok=%v", n, ok)
	}
	set := "_ZNSt3__115basic_stringbufIcNS_11char_traitsIcEENS_9allocatorIcEEE3strERKNS_12basic_stringIcS2_S3_EE"
	n, ok = parseCxx(set)
	if !ok || n.ident != "str" || !n.hasString() {
		t.Fatalf("setter %+v ok=%v", n, ok)
	}
}

func TestParseCxxTreeValueType(t *testing.T) {
	small := "_ZNSt3__16__treeINS_12__value_typeIPK8VariablejEENS_19__map_value_compareIS4_NS_4pairIKS4_jEENS_4lessIS4_EEEENS_9allocatorIS9_EEE21__construct_from_treeB9nqn220108IZNSF_21__copy_construct_treeB9nqn220108EPNS_11__tree_nodeIS5_PvEEEUlRKS9_E_EESK_SK_T_"
	n, ok := parseCxx(small)
	if !ok || n.valV != "unsigned int" || !n.trivialTreeValue() {
		t.Fatalf("small %+v ok=%v", n, ok)
	}
	if _, _, _, ok := cxxTreeCall(small, nil); !ok {
		t.Fatal("still map unsigned pair")
	}
	effect := "_ZNSt3__16__treeINS_12__value_typeIPK9Statement6EffectEENS_19__map_value_compareIS4_NS_4pairIKS4_S5_EENS_4lessIS4_EEEENS_9allocatorISA_EEE21__construct_from_treeEv"
	n, ok = parseCxx(effect)
	if ok && n.ident == "__construct_from_tree" && n.trivialTreeValue() {
		t.Fatalf("Effect must not be memcpy'd %+v", n)
	}
	if _, _, _, mapped := cxxTreeCall(effect, nil); mapped {
		t.Fatal("Effect construct_from_tree mapped")
	}
}

func TestParseCxxTreeNext(t *testing.T) {
	n, ok := parseCxx("_ZNSt3__1L11__tree_nextIPNS_16__tree_node_baseIPvEEEET_S6_")
	if !ok || !n.libcxx || n.ident != "__tree_next" {
		t.Fatalf("%+v ok=%v", n, ok)
	}
}

func TestParseCxxNewDelete(t *testing.T) {
	n, ok := parseCxx("_Znwm")
	if !ok || n.ident != "new" || n.recv != "" {
		t.Fatalf("new %+v ok=%v", n, ok)
	}
	n, ok = parseCxx("_ZdlPv")
	if !ok || n.ident != "delete" {
		t.Fatalf("delete %+v ok=%v", n, ok)
	}
}

func TestCxxReplaceBodyGotoMustJump(t *testing.T) {
	gotoFn := ir.NewFunc("_ZNK13StatementGoto9must_jumpEv", types.NewInt(1))
	body, ok := cxxReplaceBody(gotoFn)
	if !ok || len(body) != 1 {
		t.Fatalf("ok=%v n=%d", ok, len(body))
	}
	blockFn := ir.NewFunc("_ZNK5Block9must_jumpEv", types.NewInt(1))
	if _, ok := cxxReplaceBody(blockFn); ok {
		t.Fatal("Block::must_jump stays IR")
	}
	dtor := ir.NewFunc("_ZN12StatementForD1Ev", types.Void)
	if _, ok := cxxReplaceBody(dtor); !ok {
		t.Fatal("StatementFor dtor")
	}
}

func TestParseCxxToString(t *testing.T) {
	n, ok := parseCxx("_ZNSt3__19to_stringEm")
	if !ok || n.ident != "to_string" || n.recv != "" || !n.std {
		t.Fatalf("%+v ok=%v", n, ok)
	}
	signed, ok := stdToStringKind("_ZNSt3__19to_stringEm")
	if !ok || signed {
		t.Fatalf("unsigned long signed=%v ok=%v", signed, ok)
	}
	signed, ok = stdToStringKind("_ZNSt3__19to_stringEl")
	if !ok || !signed {
		t.Fatalf("long signed=%v ok=%v", signed, ok)
	}
	if isStdToString("printf") {
		t.Fatal("printf")
	}
	if isStdToString("_ZSt9to_stringm") {
		t.Fatal("libstdc++ to_string")
	}
}

func TestParseCxxRejectsC(t *testing.T) {
	if _, ok := parseCxx("printf"); ok {
		t.Fatal("printf")
	}
	if _, ok := parseCxx("llvm_fabs_f64"); ok {
		t.Fatal("llvm")
	}
}
