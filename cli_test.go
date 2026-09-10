package leaven

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestCLIStdinStdout(t *testing.T) {
	cmd := exec.Command("go", "run", "./cmd/leaven")
	cmd.Stdin = bytes.NewReader(testdataIR(t))
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
	cmd := exec.Command("go", "run", "./cmd/leaven", "-")
	cmd.Stdin = bytes.NewReader(testdataIR(t))
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

func TestCLIVersion(t *testing.T) {
	for _, args := range [][]string{{"--version"}, {"version"}} {
		cmd := exec.Command("go", append([]string{"run", "./cmd/leaven"}, args...)...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("leaven %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		got := string(bytes.TrimSpace(out))
		if !strings.HasPrefix(got, "dev") {
			t.Fatalf("leaven %s version = %q, want prefix dev", strings.Join(args, " "), got)
		}
	}
}

func TestCLIHelp(t *testing.T) {
	cmd := exec.Command("go", "run", "./cmd/leaven", "--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("leaven --help: %v\n%s", err, out)
	}
	got := string(out)
	for _, want := range []string{"Usage:", "--package", "--input", "--version"} {
		if !strings.Contains(got, want) {
			t.Fatalf("leaven --help missing %q:\n%s", want, got)
		}
	}
}
