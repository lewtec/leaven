package leaven

import "testing"

func TestCxxIOCallIfstreamStringCtor(t *testing.T) {
	name := "_ZNSt14basic_ifstreamIcSt11char_traitsIcEEC1ERKNSt7__cxx1112basic_stringIcS2_SaIcEEE"
	fn, args, ok := cxxIOCall(name, nil)
	if !ok {
		t.Fatalf("missed %s", name)
	}
	if fn == nil || len(args) != 3 {
		t.Fatalf("open %v args=%d", fn, len(args))
	}
	if _, ok := cxxIONamed(name); !ok {
		t.Fatal("named miss")
	}
}

func TestCxxIOCallIgnoresUnrelated(t *testing.T) {
	if _, _, ok := cxxIOCall("printf", nil); ok {
		t.Fatal("printf")
	}
}

func TestCxxTreeCall(t *testing.T) {
	fn, args, retPtr, ok := cxxTreeCall("_ZSt18_Rb_tree_decrementPSt18_Rb_tree_node_base", nil)
	if !ok || fn == nil || len(args) != 1 || !retPtr {
		t.Fatalf("dec %v %d ret=%v", fn, len(args), retPtr)
	}
	if _, _, _, ok := cxxTreeCall("_ZSt18_Rb_tree_incrementPKSt18_Rb_tree_node_base", nil); !ok {
		t.Fatal("inc const")
	}
	if _, a, ret, ok := cxxTreeCall("_ZSt29_Rb_tree_insert_and_rebalancebPSt18_Rb_tree_node_baseS0_RS_", nil); !ok || ret || len(a) != 4 {
		t.Fatalf("insert %v %d", ret, len(a))
	}
}

func TestCxxNoopDtor(t *testing.T) {
	if !cxxNoopDtor("_ZNSt6localeD1Ev") {
		t.Fatal("locale")
	}
	if !cxxNoopDtor("_ZNSt8ios_baseD2Ev") {
		t.Fatal("ios_base")
	}
	if cxxNoopDtor("_ZNSt14basic_ifstreamIcSt11char_traitsIcEED1Ev") {
		t.Fatal("ifstream dtor is not a noop")
	}
}
