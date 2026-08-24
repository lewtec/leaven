package leaven

import "testing"

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

func TestParseCxxRejectsC(t *testing.T) {
	if _, ok := parseCxx("printf"); ok {
		t.Fatal("printf")
	}
	if _, ok := parseCxx("llvm_fabs_f64"); ok {
		t.Fatal("llvm")
	}
}
