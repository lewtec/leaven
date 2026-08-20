package leaven

import "testing"

func TestLibcLookupDarwinSuffix(t *testing.T) {
	if _, ok := libcLookup("realpath$DARWIN_EXTSN"); !ok {
		t.Fatal("realpath$DARWIN_EXTSN")
	}
	if _, ok := libcLookup("realpath_DARWIN_EXTSN"); !ok {
		t.Fatal("realpath_DARWIN_EXTSN")
	}
	if _, ok := libcLookup("realpath"); !ok {
		t.Fatal("realpath")
	}
	if _, ok := libcLookup("fstat"); !ok {
		t.Fatal("fstat")
	}
	if _, ok := libcLookup("fstat$INODE64"); !ok {
		t.Fatal("fstat$INODE64")
	}
	if _, ok := libcLookup("lseek"); !ok {
		t.Fatal("lseek")
	}
	if _, ok := libcLookup("lseek$UNIX2003"); !ok {
		t.Fatal("lseek$UNIX2003")
	}
	if _, ok := libcLookup("close$NOCANCEL"); !ok {
		t.Fatal("close$NOCANCEL")
	}
	if _, ok := libcLookup("close_NOCANCEL"); !ok {
		t.Fatal("close_NOCANCEL")
	}
}

func TestLibcReturnsTypedPtrDarwinSuffix(t *testing.T) {
	for _, name := range []string{"realpath", "realpath$DARWIN_EXTSN", "realpath_DARWIN_EXTSN"} {
		if !libcReturnsTypedPtr(libcCanon(name)) {
			t.Fatalf("typedPtr %s", name)
		}
	}
	if libcReturnsTypedPtr(libcCanon("_NSGetArgc")) {
		t.Fatal("_NSGetArgc must not wrap")
	}
	if libcReturnsTypedPtr(libcCanon("mmap64_INODE64")) {
		t.Fatal("mmap64 must not wrap")
	}
}

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

func TestCxxIOCallLibcxxGetline(t *testing.T) {
	name := "_ZNSt3__17getlineB9nqn220108IcNS_11char_traitsIcEENS_9allocatorIcEEEERNS_13basic_istreamIT_T0_EES9_RNS_12basic_stringIS6_S7_T1_EE"
	fn, args, retPtr, ok := cxxIOCall(name, nil)
	if !ok || fn == nil || !retPtr || len(args) != 2 {
		t.Fatalf("libcxx getline %v args=%d ret=%v ok=%v", fn, len(args), retPtr, ok)
	}
	if !isGetline(name) {
		t.Fatal("isGetline libcxx")
	}
}

func TestStreambufSgetcGptrMatch(t *testing.T) {
	sgetc := "_ZNSt3__115basic_streambufIcNS_11char_traitsIcEEE5sgetcB9nqn220108Ev"
	gptr := "_ZNKSt3__115basic_streambufIcNS_11char_traitsIcEEE4gptrB9nqn220108Ev"
	if !isStreambufSgetc(sgetc) || isStreambufSgetc(gptr) {
		t.Fatal("sgetc")
	}
	if !isStreambufGptr(gptr) || isStreambufGptr(sgetc) {
		t.Fatal("gptr")
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
