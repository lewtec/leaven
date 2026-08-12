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
