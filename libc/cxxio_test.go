package libc

import (
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestStandinVptrMinus24(t *testing.T) {
	vp := StandinVptr()
	if vp == nil {
		t.Fatal("nil vptr")
	}
	off := *(*int64)(unsafe.Add(vp, -24))
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
	vp := *(*unsafe.Pointer)(unsafe.Pointer(&obj[0]))
	if vp == nil {
		t.Fatal("vptr not set")
	}
	off := *(*int64)(unsafe.Add(vp, -24))
	state := *(*int32)(unsafe.Add(unsafe.Pointer(&obj[0]), int(off)+iosStateOff))
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

func TestLocaleCtorNilSafe(t *testing.T) {
	LocaleCtor(nil)
	var obj [16]byte
	LocaleCtor(&obj[0])
}

func TestIosBaseCtorVptrAndCtype(t *testing.T) {
	var obj [272]byte
	IosBaseCtor(&obj[0])
	vp := *(*unsafe.Pointer)(unsafe.Pointer(&obj[0]))
	if vp == nil {
		t.Fatal("nil vptr")
	}
	off := *(*int64)(unsafe.Add(vp, -24))
	if off != 0 {
		t.Fatalf("vptr-24=%d", off)
	}
	ct := *(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&obj[0]), iosCtypeOff))
	if ct == nil {
		t.Fatal("nil _M_ctype")
	}
}

func TestInitOstreamVptrMinus24(t *testing.T) {
	var obj [272]byte
	InitOstream(unsafe.Pointer(&obj[0]))
	vp := *(*unsafe.Pointer)(unsafe.Pointer(&obj[0]))
	if vp == nil {
		t.Fatal("nil vptr")
	}
	off := *(*int64)(unsafe.Add(vp, -24))
	if off != 0 {
		t.Fatalf("vptr-24=%d", off)
	}
	ct := *(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&obj[0]), iosCtypeOff))
	if ct == nil {
		t.Fatal("nil _M_ctype")
	}
	if *(*byte)(unsafe.Add(ct, ctypeWidenOkOff)) == 0 {
		t.Fatal("widen_ok")
	}
	if *(*byte)(unsafe.Add(ct, ctypeWidenTabOff+10)) != '\n' {
		t.Fatal("widen \\n")
	}
}

func emptyCxxString() *[32]byte {
	var s [32]byte
	*(*unsafe.Pointer)(unsafe.Pointer(&s[0])) = unsafe.Pointer(&s[16])
	return &s
}

func cxxStringText(s *[32]byte) string {
	p := *(**byte)(unsafe.Pointer(&s[0]))
	if p == nil {
		return ""
	}
	n := *(*int64)(unsafe.Pointer(&s[8]))
	if n < 0 {
		return ""
	}
	return string(unsafe.Slice(p, int(n)))
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

func TestStringstreamStr2intDecAndHex(t *testing.T) {
	// StringUtils::str2int: stringstream(s); [ss>>hex;] ss>>i
	s := emptyCxxString()
	cxxStringAssign(&s[0], []byte("42"))
	var ss [128]byte
	StringstreamCtor(&ss[0], &s[0], 24)
	var out int32 = -1
	IstreamExtractI32(&ss[0], (*byte)(unsafe.Pointer(&out)))
	if out != 42 {
		t.Fatalf("dec got %d", out)
	}
	StringstreamClose(&ss[0])

	s2 := emptyCxxString()
	cxxStringAssign(&s2[0], []byte("0x2a"))
	var ss2 [128]byte
	StringstreamCtor(&ss2[0], &s2[0], 24)
	IstreamApplyIosManip(&ss2[0], nil)
	out = -1
	IstreamExtractI32(&ss2[0], (*byte)(unsafe.Pointer(&out)))
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

func TestOStringStreamGensym(t *testing.T) {
	// gensym: oss << basename << count; return oss.str()
	var oss [112]byte
	OStringStreamCtor(&oss[0])
	prefix := append([]byte("func_"), 0)
	OstreamLsCStr(&oss[0], &prefix[0])
	OstreamInsertI64(&oss[0], 1)
	ret := emptyCxxString()
	OStringStreamStr(&ret[0], &oss[0])
	got := string(cxxStringBytes(&ret[0]))
	if got != "func_1" {
		t.Fatalf("str %q", got)
	}
	OStringStreamClose(&oss[0])
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
	got := string(cxxStringBytes(&ret[0]))
	if got != "i0" {
		t.Fatalf("str %q", got)
	}
	StringstreamDefaultClose(&ss[0])
}
