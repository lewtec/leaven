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
