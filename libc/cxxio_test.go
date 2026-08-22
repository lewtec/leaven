package libc

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"unsafe"
)

func TestStandinVptrMinus24(t *testing.T) {
	vp := StandinVptr()
	if vp == nil {
		t.Fatal("nil vptr")
	}
	off := Load[int64](vp, -24)
	if off != 0 {
		t.Fatalf("vptr-24=%d", off)
	}
}

func TestIfstreamMissingFileFails(t *testing.T) {
	var obj [256]byte
	IfstreamOpen(&obj[0], &[]byte("no-such-leaven-platform.info\x00")[0], 8)
	if !IosFail(&obj[0]) {
		t.Fatal("missing file did not set fail")
	}
	if IosBool(&obj[0]) {
		t.Fatal("operator bool true on missing file")
	}
	// clang++ -O2: off = *(vptr-24); state = *(this+off+32); fail = state&5
	vp := Load[unsafe.Pointer](Ptr(&obj[0]), 0)
	if vp == nil {
		t.Fatal("vptr not set")
	}
	off := Load[int64](vp, -24)
	state := Load[int32](Ptr(&obj[0]), int(off)+iosStateOff)
	if state&iosFailbit == 0 {
		t.Fatalf("inlined failbit not set, state=%d off=%d", state, off)
	}
	IfstreamClose(&obj[0])
}

func TestIfstreamOpensRealFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "platform.info")
	if err := os.WriteFile(p, []byte("integer size = 4\n"), 0644); err != nil {
		t.Fatal(err)
	}
	var obj [256]byte
	IfstreamOpen(&obj[0], &append([]byte(p), 0)[0], 8)
	if IosFail(&obj[0]) {
		t.Fatal("existing file set fail")
	}
	IfstreamClose(&obj[0])
	if !IosFail(&obj[0]) {
		t.Fatal("close left a live stream")
	}
}

func TestStdStringInitSSO(t *testing.T) {
	var obj [32]byte
	s := []byte("hi")
	p := StdStringInit(&obj[0], &s[0], 2)
	if p != unsafe.Pointer(&obj[0]) {
		t.Fatal("this")
	}
	dp, n := libcxxStringData(&obj[0])
	if n != 2 {
		t.Fatalf("n=%d", n)
	}
	got := Bytes(dp, 2)
	if got[0] != 'h' || got[1] != 'i' {
		t.Fatalf("data=%q", got)
	}
}

func TestStdStringLongRoundTrip(t *testing.T) {
	var obj [32]byte
	s := []byte("abcdefghijklmnopqrstuvwxyz0123") // 30 > SSO
	StdStringInit(&obj[0], &s[0], int64(len(s)))
	if !libcxxIsLong(&obj[0]) {
		t.Fatal("expected long")
	}
	dp, n := libcxxStringData(&obj[0])
	if n != int64(len(s)) || string(Bytes(dp, int(n))) != string(s) {
		t.Fatalf("n=%d data=%q", n, Bytes(dp, int(n)))
	}
	StdStringDestroy(&obj[0])
}

func TestStdStringLongDefaultLayout(t *testing.T) {
	if libcxxAlternate() {
		t.Skip("arm64 alternate")
	}
	var obj [32]byte
	s := []byte("abcdefghijklmnopqrstuvwxyz0123")
	StdStringInit(&obj[0], &s[0], int64(len(s)))
	if Load[uint64](Ptr(&obj[0]), 0)&1 == 0 {
		t.Fatal("long bit in word 0")
	}
	if Load[*byte](Ptr(&obj[0]), 16) == nil {
		t.Fatal("data ptr at +16")
	}
	StdStringDestroy(&obj[0])
}

func TestStdStringSubstr(t *testing.T) {
	var src, dst [32]byte
	s := []byte("abcdef")
	StdStringInit(&src[0], &s[0], 6)
	StdStringSubstr(&dst[0], &src[0], 2, 3, nil)
	dp, n := libcxxStringData(&dst[0])
	got := Bytes(dp, int(n))
	if n != 3 || string(got) != "cde" {
		t.Fatalf("n=%d data=%q", n, got)
	}
	c := append([]byte("cde"), 0)
	if !StdStringEqCStr(&dst[0], &c[0]) {
		t.Fatal("eq")
	}
	if StdStringCompareCStr(&src[0], 2, 3, &c[0], 3) != 0 {
		t.Fatal("compare")
	}
	StdStringErase(&src[0], 0, 1)
	if StdStringCompareCStr(&src[0], 0, -1, &[]byte("bcdef\x00")[0], -1) != 0 {
		t.Fatal("erase")
	}
	StdStringAppendCStr(&src[0], &[]byte("x\x00")[0], -1)
	if StdStringCompareCStr(&src[0], 0, -1, &[]byte("bcdefx\x00")[0], -1) != 0 {
		t.Fatal("append")
	}
	StdStringAssignCStr(&src[0], &[]byte("z\x00")[0], -1)
	if StdStringCompareCStr(&src[0], 0, -1, &[]byte("z\x00")[0], -1) != 0 {
		t.Fatal("assign")
	}
	StdStringPushBack(&src[0], 'y')
	if StdStringCompareCStr(&src[0], 0, -1, &[]byte("zy\x00")[0], -1) != 0 {
		t.Fatal("push_back")
	}
	StdStringInsertCStr(&src[0], 1, &[]byte("Q\x00")[0], -1)
	if StdStringCompareCStr(&src[0], 0, -1, &[]byte("zQy\x00")[0], -1) != 0 {
		t.Fatal("insert")
	}
}

func TestLocaleCtorNilSafe(t *testing.T) {
	LocaleCtor(nil)
	var obj [16]byte
	LocaleCtor(&obj[0])
}

func TestIosBaseCtorVptrAndCtype(t *testing.T) {
	var obj [272]byte
	IosBaseCtor(&obj[0])
	vp := Load[unsafe.Pointer](Ptr(&obj[0]), 0)
	if vp == nil {
		t.Fatal("nil vptr")
	}
	off := Load[int64](vp, -24)
	if off != 0 {
		t.Fatalf("vptr-24=%d", off)
	}
	ct := Load[unsafe.Pointer](Ptr(&obj[0]), iosCtypeOff)
	if ct == nil {
		t.Fatal("nil _M_ctype")
	}
	if Load[int32](Ptr(&obj[0]), iosStateOff) != 0 {
		t.Fatal("fresh ios_base must be goodbit")
	}
}

func TestInitOstreamVptrMinus24(t *testing.T) {
	var obj [272]byte
	InitOstream(Ptr(&obj[0]))
	vp := Load[unsafe.Pointer](Ptr(&obj[0]), 0)
	if vp == nil {
		t.Fatal("nil vptr")
	}
	off := Load[int64](vp, -24)
	if off != 0 {
		t.Fatalf("vptr-24=%d", off)
	}
	ct := Load[unsafe.Pointer](Ptr(&obj[0]), iosCtypeOff)
	if ct == nil {
		t.Fatal("nil _M_ctype")
	}
	if Load[byte](ct, ctypeWidenOkOff) == 0 {
		t.Fatal("widen_ok")
	}
	if Load[byte](ct, ctypeWidenTabOff+10) != '\n' {
		t.Fatal("widen \\n")
	}
}

func putTestCxxString(text string) *[32]byte {
	var s [32]byte
	if runtime.GOOS == "darwin" {
		var src *byte
		b := []byte(text)
		if len(b) > 0 {
			src = &b[0]
		}
		StdStringInit(&s[0], src, int64(len(b)))
		return &s
	}
	o := emptyCxxString()
	cxxStringAssign(&o[0], []byte(text))
	return o
}

func emptyCxxString() *[32]byte {
	var s [32]byte
	Store(Ptr(&s[0]), 0, Ptr(&s[16]))
	return &s
}

func cxxStringText(s *[32]byte) string {
	if runtime.GOOS == "darwin" {
		p, n := libcxxStringData(&s[0])
		if p == nil || n <= 0 {
			return ""
		}
		return string(Bytes(p, int(n)))
	}
	p := Load[*byte](Ptr(&s[0]), 0)
	if p == nil {
		return ""
	}
	n := Load[int64](Ptr(&s[0]), 8)
	if n < 0 {
		return ""
	}
	return string(Bytes(p, int(n)))
}

func TestFilebufOpenFillsGetArea(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "platform.info")
	if err := os.WriteFile(p, []byte("integer size = 4\n"), 0644); err != nil {
		t.Fatal(err)
	}
	var obj [256]byte
	fb := &obj[filebufOff]
	got := FilebufOpen(fb, &append([]byte(p), 0)[0], 8)
	if got == nil {
		t.Fatal("open")
	}
	gptr := Load[*byte](Ptr(fb), sbGptrOff)
	egptr := Load[*byte](Ptr(fb), sbEgptrOff)
	if gptr == nil || egptr == nil || Addr(gptr) >= Addr(egptr) {
		t.Fatal("empty get area")
	}
	if StreambufSgetc(fb) != int32('i') {
		t.Fatalf("sgetc=%d", StreambufSgetc(fb))
	}
	if IosFail(fb) {
		t.Fatal("open left filebuf failed")
	}
	if IosFail(&obj[0]) {
		t.Fatal("open left ifstream failed")
	}
	FilebufClose(fb)
}

func TestIstreamGetlineMissingFails(t *testing.T) {
	var obj [256]byte
	var str = emptyCxxString()
	got := IstreamGetline(&obj[0], &str[0])
	if got != &obj[0] {
		t.Fatal("did not return the stream")
	}
	if !IosFail(&obj[0]) {
		t.Fatal("missing stream did not fail")
	}
	if !IosEof(&obj[0]) {
		t.Fatal("missing stream did not set eof")
	}
}

func TestIstreamGetlineKeepsFail(t *testing.T) {
	var obj [256]byte
	IfstreamOpen(&obj[0], &[]byte("no-such-leaven-getline.info\x00")[0], 8)
	if !IosFail(&obj[0]) {
		t.Fatal("open of missing file should fail")
	}
	var str = emptyCxxString()
	IstreamGetline(&obj[0], &str[0])
	if !IosFail(&obj[0]) {
		t.Fatal("getline cleared fail on a bad stream")
	}
	if !IosEof(&obj[0]) {
		t.Fatal("getline on a bad stream left eof clear")
	}
	IfstreamClose(&obj[0])
}

func TestIstreamGetlineReadsLine(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "platform.info")
	if err := os.WriteFile(p, []byte("integer size = 4\npointer size = 8\n"), 0644); err != nil {
		t.Fatal(err)
	}
	var obj [256]byte
	IfstreamOpen(&obj[0], &append([]byte(p), 0)[0], 8)
	if IosFail(&obj[0]) {
		t.Fatal("open failed")
	}
	s1 := emptyCxxString()
	IstreamGetline(&obj[0], &s1[0])
	if IosFail(&obj[0]) {
		t.Fatal("first line failed")
	}
	if got := cxxStringText(s1); got != "integer size = 4" {
		t.Fatalf("line1 %q", got)
	}
	s2 := emptyCxxString()
	IstreamGetline(&obj[0], &s2[0])
	if IosFail(&obj[0]) {
		t.Fatal("second line failed")
	}
	if got := cxxStringText(s2); got != "pointer size = 8" {
		t.Fatalf("line2 %q", got)
	}
	s3 := emptyCxxString()
	IstreamGetline(&obj[0], &s3[0])
	if !IosFail(&obj[0]) || !IosEof(&obj[0]) {
		t.Fatal("EOF getline should set fail+eof")
	}
	IfstreamClose(&obj[0])
}

func TestStringbufStrThenExtract(t *testing.T) {
	// Darwin CGOptions: stringstream(s) inlines as stringbuf::str; >> uses the parent.
	s := putTestCxxString("42")
	var ss [176]byte
	sb := &ss[16]
	StringbufStr(sb, &s[0])
	if StreambufSgetc(sb) != int32('4') {
		t.Fatalf("sgetc=%d", StreambufSgetc(sb))
	}
	var out int32 = -1
	IstreamExtractI32(&ss[0], As[byte](Ptr(&out)))
	if out != 42 {
		t.Fatalf("extract parent got %d", out)
	}
	StringstreamClose(sb)
	StringstreamClose(&ss[0])
}

func TestGoCxxStringBytesLibcxxShort(t *testing.T) {
	var s [32]byte
	StdStringInit(&s[0], &[]byte("42")[0], 2)
	p, n := libcxxStringData(&s[0])
	if n != 2 || string(Bytes(p, 2)) != "42" {
		t.Fatalf("libcxx n=%d", n)
	}
	if runtime.GOOS == "darwin" {
		if got := string(goCxxStringBytes(&s[0])); got != "42" {
			t.Fatalf("darwin go %q", got)
		}
		return
	}
	libstd := emptyCxxString()
	cxxStringAssign(&libstd[0], []byte("42"))
	if got := string(goCxxStringBytes(&libstd[0])); got != "42" {
		t.Fatalf("libstd %q", got)
	}
}

func TestStringstreamStr2intDecAndHex(t *testing.T) {
	// StringUtils::str2int: stringstream(s); [ss>>hex;] ss>>i
	s := putTestCxxString("42")
	var ss [128]byte
	StringstreamCtor(&ss[0], &s[0], 24)
	var out int32 = -1
	IstreamExtractI32(&ss[0], As[byte](Ptr(&out)))
	if out != 42 {
		t.Fatalf("dec got %d", out)
	}
	StringstreamClose(&ss[0])

	s2 := putTestCxxString("0x2a")
	var ss2 [128]byte
	StringstreamCtor(&ss2[0], &s2[0], 24)
	IstreamApplyIosManip(&ss2[0], nil)
	out = -1
	IstreamExtractI32(&ss2[0], As[byte](Ptr(&out)))
	if out != 0x2a {
		t.Fatalf("hex got %d want %d", out, 0x2a)
	}
	StringstreamClose(&ss2[0])
}

func TestOstreamInsertWrites(t *testing.T) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	msg := []byte("/*\n")
	var dummy [8]byte
	MarkStdoutStream(unsafe.Pointer(&dummy[0]))
	defer stdoutStreams.Delete(Addr(&dummy[0]))
	got := OstreamInsert(&dummy[0], &msg[0], int64(len(msg)))
	_ = w.Close()
	os.Stdout = old
	if got != &dummy[0] {
		t.Fatal("did not return the stream")
	}
	buf := make([]byte, 8)
	n, _ := r.Read(buf)
	_ = r.Close()
	if string(buf[:n]) != "/*\n" {
		t.Fatalf("wrote %q", buf[:n])
	}
}

func TestWriteOstreamLazyGensym(t *testing.T) {
	// Inlined ostringstream ctor: no table entry; << must not hit stdout.
	var oss [192]byte
	defer OStringStreamClose(&oss[0])
	OstreamInsertI64(&oss[0], 12)
	OstreamInsertI64(&oss[0], 34)
	ret := emptyCxxString()
	OStringStreamStr(&ret[0], &oss[0])
	if got := string(goCxxStringBytes(&ret[0])); got != "1234" {
		t.Fatalf("lazy str %q", got)
	}
}

func TestWriteOstreamLazyParentStr(t *testing.T) {
	// << on the ostream subobject; str() on the ostringstream.
	var ss [192]byte
	os := &ss[16]
	defer OStringStreamClose(os)
	defer OStringStreamClose(&ss[0])
	OstreamLsCStr(os, &[]byte("g_\x00")[0])
	OstreamInsertI64(os, 1)
	ret := emptyCxxString()
	OStringStreamStr(&ret[0], &ss[0])
	if got := string(goCxxStringBytes(&ret[0])); got != "g_1" {
		t.Fatalf("parent str %q", got)
	}
	sb := stringbufOf(os)
	pbase := Load[*byte](Ptr(sb), sbPbaseOff)
	pptr := Load[*byte](Ptr(sb), sbPptrOff)
	if pbase == nil || pptr == nil || Addr(pbase) >= Addr(pptr) {
		t.Fatal("empty put area")
	}
	if string(Bytes(pbase, int(Addr(pptr)-Addr(pbase)))) != "g_1" {
		t.Fatalf("put area %q", Bytes(pbase, int(Addr(pptr)-Addr(pbase))))
	}
	if runtime.GOOS == "darwin" {
		if Load[*byte](Ptr(sb), sbHmOff) != pptr {
			t.Fatal("hm")
		}
		if Load[int32](Ptr(sb), sbModeOff)&iosModeOut == 0 {
			t.Fatalf("mode=%d", Load[int32](Ptr(sb), sbModeOff))
		}
		sp, sn := libcxxStringData(As[byte](Off(Ptr(sb), sbStrOff)))
		if sn != 3 || string(Bytes(sp, int(sn))) != "g_1" {
			t.Fatalf("__str_ %q", Bytes(sp, int(sn)))
		}
	}
	OStringStreamClose(os)
	OStringStreamClose(&ss[0])
}

func TestOStringStreamGensym(t *testing.T) {
	// gensym: oss << basename << count; return oss.str()
	var oss [192]byte
	OStringStreamCtor(&oss[0])
	defer OStringStreamClose(&oss[0])
	prefix := append([]byte("func_"), 0)
	OstreamLsCStr(&oss[0], &prefix[0])
	OstreamInsertI64(&oss[0], 1)
	ret := emptyCxxString()
	OStringStreamStr(&ret[0], &oss[0])
	got := string(goCxxStringBytes(&ret[0]))
	if got != "func_1" {
		t.Fatalf("str %q", got)
	}
	OStringStreamClose(&oss[0])
}

func TestOStringStreamStrOnStringbufThis(t *testing.T) {
	// Darwin: << on ostringstream, inlined str() on __sb_ at +8.
	if runtime.GOOS != "darwin" {
		t.Skip("libstdc++ str() uses the ostringstream this")
	}
	var oss [192]byte
	OStringStreamCtor(&oss[0])
	defer OStringStreamClose(&oss[0])
	OstreamLsCStr(&oss[0], &[]byte("g_\x00")[0])
	OstreamInsertI64(&oss[0], 3)
	ret := emptyCxxString()
	OStringStreamStr(&ret[0], &oss[libcxxOStringSBOff])
	if got := string(goCxxStringBytes(&ret[0])); got != "g_3" {
		t.Fatalf("sb this str %q", got)
	}
}

func TestOStringStreamStrAfterMovedThis(t *testing.T) {
	// Go stack copy changes &oss; put-area fields move with the object.
	var oss [192]byte
	OStringStreamCtor(&oss[0])
	defer OStringStreamClose(&oss[0])
	OstreamLsCStr(&oss[0], &[]byte("g_\x00")[0])
	OstreamInsertI64(&oss[0], 2)
	var moved [192]byte
	copy(moved[:], oss[:])
	ret := emptyCxxString()
	OStringStreamStr(&ret[0], &moved[0])
	if got := string(goCxxStringBytes(&ret[0])); got != "g_2" {
		t.Fatalf("moved str %q", got)
	}
}

func TestStringstreamDefaultNewCtrlVars(t *testing.T) {
	// Variable::new_ctrl_vars: ss(); << 'i' via this+16; << n; str()
	var ss [128]byte
	StringstreamDefaultCtor(&ss[0])
	os := &ss[stringstreamOstreamOff]
	OstreamPut(os, 'i')
	OstreamInsertU64(os, 0)
	ret := emptyCxxString()
	StringstreamStr(&ret[0], &ss[0])
	got := string(goCxxStringBytes(&ret[0]))
	if got != "i0" {
		t.Fatalf("str %q", got)
	}
	StringstreamDefaultClose(&ss[0])
}
