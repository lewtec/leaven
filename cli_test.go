package leaven

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func TestCLIStdinStdout(t *testing.T) {
	ir, err := os.ReadFile("testdata/hello.ll")
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "run", "./cmd/leaven")
	cmd.Stdin = bytes.NewReader(ir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("leaven stdin: %v\n%s", err, stderr.Bytes())
	}
	if !bytes.Contains(stdout.Bytes(), []byte("package main")) {
		t.Fatalf("stdout missing package clause:\n%s\nstderr:\n%s", stdout.Bytes(), stderr.Bytes())
	}
}

func TestCLIDashStdin(t *testing.T) {
	ir, err := os.ReadFile("testdata/hello.ll")
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "run", "./cmd/leaven", "-")
	cmd.Stdin = bytes.NewReader(ir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("leaven -: %v\n%s", err, stderr.Bytes())
	}
	if !bytes.Contains(stdout.Bytes(), []byte("package main")) {
		t.Fatalf("stdout missing package clause:\n%s\nstderr:\n%s", stdout.Bytes(), stderr.Bytes())
	}
}
