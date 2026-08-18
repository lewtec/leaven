package libc

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFcntlGetFL(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no host fcntl")
	}
	dir := t.TempDir()
	p := filepath.Join(dir, "f")
	if err := os.WriteFile(p, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	path := append([]byte(p), 0)
	fd := Open(&path[0], 0)
	if fd < 0 {
		t.Fatal("open")
	}
	defer Close(fd)
	const fGetFL = 3
	got := Fcntl(fd, fGetFL)
	if got < 0 {
		t.Fatalf("F_GETFL=%d", got)
	}
	if Fcntl(-1, fGetFL) != -1 {
		t.Fatal("bad fd")
	}
}

func TestGetentropy(t *testing.T) {
	var a, b [16]byte
	if Getentropy(&a[0], 16) != 0 || Getentropy(&b[0], 16) != 0 {
		t.Fatal("getentropy")
	}
	if a == b {
		t.Fatal("two fills were identical")
	}
	if Getentropy(nil, 0) != 0 {
		t.Fatal("n=0")
	}
	if Getentropy(nil, 8) != -1 {
		t.Fatal("nil buf")
	}
	if Getentropy(&a[0], 257) != -1 {
		t.Fatal("n>256")
	}
}

func TestStrerrorR(t *testing.T) {
	var buf [16]byte
	if StrerrorR(1, &buf[0], int64(len(buf))) != 0 {
		t.Fatal("strerror_r")
	}
	if buf[0] == 0 {
		t.Fatal("empty")
	}
	if StrerrorR(1, nil, 8) != -1 {
		t.Fatal("nil buf")
	}
}

func TestNSGetArgc(t *testing.T) {
	p := NSGetArgc()
	if p == nil || *(*int32)(p) < 1 {
		t.Fatal("argc")
	}
	argv := NSGetArgv()
	if argv == nil {
		t.Fatal("argv")
	}
}

func TestPthreadGetStackaddrNp(t *testing.T) {
	p := PthreadGetStackaddrNp(PthreadSelf())
	if p == nil {
		t.Fatal("nil stackaddr")
	}
	if n := PthreadGetStacksizeNp(PthreadSelf()); n != dummyStackSize {
		t.Fatalf("stacksize=%d", n)
	}
}
