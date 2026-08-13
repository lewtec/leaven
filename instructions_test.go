package leaven

import "testing"

func TestCxxIOCallIfstreamStringCtor(t *testing.T) {
	name := "_ZNSt14basic_ifstreamIcSt11char_traitsIcEEC1ERKNSt7__cxx1112basic_stringIcS2_SaIcEEE"
	fn, args, retPtr, ok := cxxIOCall(name, nil)
	if !ok {
		t.Fatalf("missed %s", name)
	}
	if fn == nil || len(args) != 3 || retPtr {
		t.Fatalf("open %v args=%d ret=%v", fn, len(args), retPtr)
	}
	if _, ok := cxxIONamed(name); !ok {
		t.Fatal("named miss")
	}
}

func TestCxxIOCallIgnoresUnrelated(t *testing.T) {
	if _, _, _, ok := cxxIOCall("printf", nil); ok {
		t.Fatal("printf")
	}
}

func TestCxxOstreamOpEndl(t *testing.T) {
	name := "_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_"
	fn, args, retPtr, ok := cxxIOCall(name, nil)
	if !ok || fn == nil || !retPtr || len(args) != 1 {
		t.Fatalf("endl %v args=%d ret=%v ok=%v", fn, len(args), retPtr, ok)
	}
	if _, ok := cxxIONamed(name); !ok {
		t.Fatal("named miss")
	}
}

func TestCxxIOCallOstreamInsert(t *testing.T) {
	name := "_ZSt16__ostream_insertIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_PKS3_l"
	fn, args, retPtr, ok := cxxIOCall(name, nil)
	if !ok || fn == nil || !retPtr || len(args) != 3 {
		t.Fatalf("insert %v args=%d ret=%v ok=%v", fn, len(args), retPtr, ok)
	}
	if _, ok := cxxIONamed(name); !ok {
		t.Fatal("named miss")
	}
}

func TestCxxIOCallGetline(t *testing.T) {
	name := "_ZSt7getlineIcSt11char_traitsIcESaIcEERSt13basic_istreamIT_T0_ES7_RNSt7__cxx1112basic_stringIS4_S5_T1_EE"
	fn, args, retPtr, ok := cxxIOCall(name, nil)
	if !ok || fn == nil || !retPtr || len(args) != 2 {
		t.Fatalf("getline %v args=%d ret=%v ok=%v", fn, len(args), retPtr, ok)
	}
	if _, ok := cxxIONamed(name); !ok {
		t.Fatal("named miss")
	}
	if !isGetline(name) {
		t.Fatal("isGetline")
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

func TestCxxIOCallLocaleCtor(t *testing.T) {
	name := "_ZNSt6localeC1Ev"
	fn, args, retPtr, ok := cxxIOCall(name, nil)
	if !ok || fn == nil || retPtr || len(args) != 1 {
		t.Fatalf("locale %v args=%d ret=%v ok=%v", fn, len(args), retPtr, ok)
	}
	if _, ok := cxxIONamed(name); !ok {
		t.Fatal("named miss")
	}
	if !isLocaleCtor("_ZNSt6localeC2ERKS_") {
		t.Fatal("copy ctor")
	}
}

func TestCxxIOCallIosBaseCtor(t *testing.T) {
	name := "_ZNSt8ios_baseC2Ev"
	fn, args, retPtr, ok := cxxIOCall(name, nil)
	if !ok || fn == nil || retPtr || len(args) != 1 {
		t.Fatalf("ctor %v args=%d ret=%v ok=%v", fn, len(args), retPtr, ok)
	}
	if _, ok := cxxIONamed(name); !ok {
		t.Fatal("named miss")
	}
	if _, _, _, ok := cxxIOCall("_ZNSt9basic_iosIcSt11char_traitsIcEE4initEPSt15basic_streambufIcS1_E", nil); !ok {
		t.Fatal("basic_ios::init")
	}
	if !isIosBaseCtor("_ZNSt8ios_base7_M_initEv") {
		t.Fatal("_M_init")
	}
}

func TestCxxIOCallStringstream(t *testing.T) {
	ctor := "_ZNSt7__cxx1118basic_stringstreamIcSt11char_traitsIcESaIcEEC1ERKNS_12basic_stringIcS2_S3_EESt13_Ios_Openmode"
	fn, args, retPtr, ok := cxxIOCall(ctor, nil)
	if !ok || fn == nil || retPtr || len(args) != 3 {
		t.Fatalf("ctor %v args=%d ret=%v ok=%v", fn, len(args), retPtr, ok)
	}
	if _, ok := cxxIONamed(ctor); !ok {
		t.Fatal("named ctor")
	}
	dtor := "_ZNSt7__cxx1118basic_stringstreamIcSt11char_traitsIcESaIcEED1Ev"
	if _, a, ret, ok := cxxIOCall(dtor, nil); !ok || ret || len(a) != 1 {
		t.Fatalf("dtor %v %d", ret, len(a))
	}
	manip := "_ZNSirsEPFRSt8ios_baseS0_E"
	if _, a, ret, ok := cxxIOCall(manip, nil); !ok || !ret || len(a) != 2 {
		t.Fatalf("manip ret=%v n=%d", ret, len(a))
	}
	extract := "_ZNSirsERi"
	if _, a, ret, ok := cxxIOCall(extract, nil); !ok || !ret || len(a) != 2 {
		t.Fatalf("extract ret=%v n=%d", ret, len(a))
	}
}

func TestCxxIOCallOstringstream(t *testing.T) {
	ctor := "_ZNSt7__cxx1119basic_ostringstreamIcSt11char_traitsIcESaIcEEC1Ev"
	fn, args, retPtr, ok := cxxIOCall(ctor, nil)
	if !ok || fn == nil || retPtr || len(args) != 1 {
		t.Fatalf("ctor %v args=%d ret=%v ok=%v", fn, len(args), retPtr, ok)
	}
	str := "_ZNKRSt7__cxx1119basic_ostringstreamIcSt11char_traitsIcESaIcEE3strEv"
	if _, a, ret, ok := cxxIOCall(str, nil); !ok || ret || len(a) != 2 {
		t.Fatalf("str ret=%v n=%d", ret, len(a))
	}
	dtor := "_ZNSt7__cxx1119basic_ostringstreamIcSt11char_traitsIcESaIcEED1Ev"
	if _, a, ret, ok := cxxIOCall(dtor, nil); !ok || ret || len(a) != 1 {
		t.Fatalf("dtor %v %d", ret, len(a))
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
