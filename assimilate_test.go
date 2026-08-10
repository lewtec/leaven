package leaven

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestAssimilate builds real upstream trees (cmake+clang++ 22 for csmith,
// cargo+clang 22 for rhai), emits full-program LLVM IR, and compares native
// stdout to the leaven-compiled Go program on the same argv.
func TestAssimilate(t *testing.T) {
	if testing.Short() {
		t.Skip("-short: assimilate is cmake + cargo")
	}
	t.Run("csmith", testAssimilateCsmith)
	t.Run("rhai", testAssimilateRhai)
}

func testAssimilateCsmith(t *testing.T) {
	root := requireProject(t, "csmith", "src/RandomProgramGenerator.cpp")
	clang := miseWhich(t, "clang", "conda:clang@22.1.8")
	clangxx := miseWhich(t, "clang++", "conda:clangxx@22.1.8")
	cmake := miseWhich(t, "cmake", "cmake@4.4.1")
	ninja := miseWhich(t, "ninja", "ninja@1.13.2")
	m4 := miseWhich(t, "m4", "conda:m4@1.4.20")
	link := llvmLink22(t)

	build := t.TempDir()
	cmakeConfigure(t, cmake, ninja, clang, clangxx, m4, root, build, []string{
		"-DCMAKE_CXX_FLAGS=-fno-exceptions",
	})
	cmakeBuild(t, cmake, ninja, build, "csmith")

	native := filepath.Join(build, "src", "csmith")
	if _, err := os.Stat(native); err != nil {
		t.Fatalf("cmake native missing %s: %v", native, err)
	}

	ll := filepath.Join(build, "csmith.ll")
	emitIRFromCompileCommands(t, build, ll, link)
	crossCheck(t, native, ll, []string{"-s", "1"})
}

func testAssimilateRhai(t *testing.T) {
	root := requireProject(t, "rhai", "Cargo.toml")
	clang, sysroot, libdir := clang22LinkEnv(t)
	clangxx := filepath.Join(filepath.Dir(clang), "clang++")
	env := append(os.Environ(),
		"CC="+clang,
		"CXX="+clangxx,
		"CARGO_TARGET_X86_64_UNKNOWN_LINUX_GNU_LINKER="+clang,
		"RUSTFLAGS=-C debuginfo=0 -C linker="+clang+" -C link-arg=--sysroot="+sysroot+" -C link-arg=-L"+libdir,
	)

	build := exec.Command("mise", "exec", "--", "cargo", "build", "--release", "--bin", "rhai-run")
	build.Dir = root
	build.Env = env
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("cargo build rhai-run: %v\n%s", err, tailBytes(out, 4000))
	}
	native := filepath.Join(root, "target", "release", "rhai-run")

	rustc := exec.Command("mise", "exec", "--", "cargo", "rustc", "--release", "--bin", "rhai-run",
		"--", "-C", "panic=abort", "--emit=llvm-ir", "-C", "debuginfo=0")
	rustc.Dir = root
	rustc.Env = env
	if out, err := rustc.CombinedOutput(); err != nil {
		t.Fatalf("cargo rustc rhai-run: %v\n%s", err, tailBytes(out, 4000))
	}
	matches, err := filepath.Glob(filepath.Join(root, "target", "release", "deps", "rhai_run-*.ll"))
	if err != nil || len(matches) == 0 {
		t.Fatalf("no rhai_run-*.ll after cargo rustc (err=%v)", err)
	}

	script := filepath.Join(t.TempDir(), "probe.rhai")
	if err := os.WriteFile(script, []byte("print(40 + 2);\n"), 0644); err != nil {
		t.Fatal(err)
	}
	crossCheck(t, native, matches[0], []string{script})
}

type compileCommand struct {
	Directory string   `json:"directory"`
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
	File      string   `json:"file"`
	Output    string   `json:"output"`
}

func cmakeConfigure(t *testing.T, cmake, ninja, clang, clangxx, m4, src, dst string, extra []string) {
	t.Helper()
	args := []string{
		"-S", src, "-B", dst,
		"-G", "Ninja",
		"-DCMAKE_BUILD_TYPE=Release",
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-DCMAKE_C_COMPILER=" + clang,
		"-DCMAKE_CXX_COMPILER=" + clangxx,
		"-DCMAKE_MAKE_PROGRAM=" + ninja,
		"-DM4=" + m4,
	}
	args = append(args, extra...)
	cmd := exec.Command(cmake, args...)
	cmd.Env = append(os.Environ(), "CC="+clang, "CXX="+clangxx)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cmake configure: %v\n%s", err, tailBytes(out, 4000))
	}
}

func cmakeBuild(t *testing.T, cmake, ninja, dir string, targets ...string) {
	t.Helper()
	args := []string{"--build", dir, "--parallel"}
	for _, tgt := range targets {
		args = append(args, "--target", tgt)
	}
	cmd := exec.Command(cmake, args...)
	cmd.Env = append(os.Environ(), "CMAKE_MAKE_PROGRAM="+ninja)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cmake --build: %v\n%s", err, tailBytes(out, 8000))
	}
}

func emitIRFromCompileCommands(t *testing.T, build, outLL, llvmLink string) {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(build, "compile_commands.json"))
	if err != nil {
		t.Fatalf("compile_commands.json: %v", err)
	}
	var cmds []compileCommand
	if err := json.Unmarshal(raw, &cmds); err != nil {
		t.Fatalf("compile_commands.json: %v", err)
	}
	irDir := filepath.Join(build, "ll")
	if err := os.MkdirAll(irDir, 0755); err != nil {
		t.Fatal(err)
	}
	var lls []string
	for i, c := range cmds {
		if c.Output != "" && !strings.Contains(c.Output, "CMakeFiles/csmith.dir") {
			continue
		}
		ext := strings.ToLower(filepath.Ext(c.File))
		if ext != ".c" && ext != ".cc" && ext != ".cpp" && ext != ".cxx" {
			continue
		}
		if c.Output == "" && strings.Contains(filepath.ToSlash(c.File), "/runtime/") {
			continue
		}
		args := c.Arguments
		if len(args) == 0 {
			args = splitQuoted(c.Command)
		}
		if len(args) == 0 {
			t.Fatalf("empty compile command for %s", c.File)
		}
		ll := filepath.Join(irDir, fmt.Sprintf("%d-%s.ll", i, filepath.Base(c.File)))
		args = withEmitLLVM(args, ll)
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = c.Directory
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("clang -emit-llvm %s: %v\n%s", filepath.Base(c.File), err, tailBytes(out, 4000))
		}
		lls = append(lls, ll)
	}
	if len(lls) == 0 {
		t.Fatalf("compile_commands.json had 0 C/C++ TUs")
	}
	linkArgs := append([]string{"-S", "-o", outLL}, lls...)
	cmd := exec.Command(llvmLink, linkArgs...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("llvm-link %d TUs: %v\n%s", len(lls), err, tailBytes(out, 4000))
	}
	t.Logf("llvm-link %d TUs -> %s", len(lls), filepath.Base(outLL))
}

func withEmitLLVM(args []string, ll string) []string {
	out := make([]string, 0, len(args)+4)
	skip := false
	for i := 0; i < len(args); i++ {
		if skip {
			skip = false
			continue
		}
		a := args[i]
		switch {
		case a == "-o" && i+1 < len(args):
			skip = true
		case strings.HasPrefix(a, "-o") && a != "-o":
			// -ofoo.o
		case a == "-c" || a == "-S":
		default:
			out = append(out, a)
		}
	}
	return append(out, "-S", "-emit-llvm", "-fno-discard-value-names", "-o", ll)
}

func splitQuoted(s string) []string {
	var out []string
	var b strings.Builder
	var q byte
	esc := false
	flush := func() {
		if b.Len() == 0 {
			return
		}
		out = append(out, b.String())
		b.Reset()
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if esc {
			b.WriteByte(c)
			esc = false
			continue
		}
		if c == '\\' && q != '\'' {
			esc = true
			continue
		}
		if q != 0 {
			if c == q {
				q = 0
				continue
			}
			b.WriteByte(c)
			continue
		}
		switch c {
		case '\'', '"':
			q = c
		case ' ', '\t', '\n':
			flush()
		default:
			b.WriteByte(c)
		}
	}
	flush()
	return out
}

func crossCheck(t *testing.T, native, ll string, args []string) {
	t.Helper()
	want := runTimeout(t, 15*time.Second, native, args...)

	m, err := parseIRFile(ll)
	if err != nil {
		t.Fatalf("parse %s: %v\n%s", filepath.Base(ll), err, irSnippet(ll, err))
	}
	var buf bytes.Buffer
	if err := Compile(&buf, m, "main"); err != nil {
		t.Fatalf("compile %s: %v", filepath.Base(ll), err)
	}
	src, err := format.Source(buf.Bytes())
	if err != nil {
		t.Fatalf("gofmt %s: %v", filepath.Base(ll), err)
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), src, 0644); err != nil {
		t.Fatal(err)
	}
	modRoot, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	gomod := fmt.Sprintf("module leavenfixture\n\ngo 1.22\n\nrequire github.com/lewtec/leaven v0.0.0\n\nreplace github.com/lewtec/leaven => %s\n", modRoot)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0644); err != nil {
		t.Fatal(err)
	}
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = dir
	if out, err := tidy.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, tailBytes(out, 2000))
	}
	got := runGoDir(t, dir, args...)
	if !bytes.Equal(want, got) {
		t.Fatalf("native vs leaven mismatch\n---- native (%d bytes) ----\n%s\n---- leaven (%d bytes) ----\n%s",
			len(want), tailBytes(want, 2000), len(got), tailBytes(got, 2000))
	}
}

func runGoDir(t *testing.T, dir string, args ...string) []byte {
	t.Helper()
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Start(); err != nil {
		t.Fatalf("go run: %v", err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("go run: %v\n%s", err, tailBytes(buf.Bytes(), 4000))
		}
	case <-time.After(30 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatalf("go run timeout\n%s", tailBytes(buf.Bytes(), 4000))
	}
	return buf.Bytes()
}

func runTimeout(t *testing.T, d time.Duration, bin string, args ...string) []byte {
	t.Helper()
	cmd := exec.Command(bin, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Start(); err != nil {
		t.Fatalf("%s: %v", filepath.Base(bin), err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("%s: %v\n%s", filepath.Base(bin), err, tailBytes(buf.Bytes(), 4000))
		}
	case <-time.After(d):
		_ = cmd.Process.Kill()
		t.Fatalf("%s timeout after %s\n%s", filepath.Base(bin), d, tailBytes(buf.Bytes(), 4000))
	}
	return buf.Bytes()
}

func requireProject(t *testing.T, name, marker string) string {
	t.Helper()
	root := filepath.Join("testdata", "projects", name)
	mark := filepath.Join(root, marker)
	if _, err := os.Stat(mark); err != nil {
		syncProjects(t)
	}
	if _, err := os.Stat(mark); err != nil {
		t.Fatalf("%s not placed after workspaced apply: %v", root, err)
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

func syncProjects(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("mise"); err != nil {
		t.Fatal("mise not on PATH")
	}
	cmd := exec.Command("mise", "exec", "--", "workspaced", "codebase", "apply")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("workspaced codebase apply: %v\n%s", err, tailBytes(out, 2000))
	}
}

func llvmLink22(t *testing.T) string {
	t.Helper()
	return miseWhich(t, "llvm-link", "conda:llvm-tools@22.1.8")
}

func clang22LinkEnv(t *testing.T) (clang, sysroot, libdir string) {
	t.Helper()
	clang = miseWhich(t, "clang", "conda:clang@22.1.8")
	prefix := filepath.Dir(filepath.Dir(clang))
	sysroot = filepath.Join(prefix, "x86_64-conda-linux-gnu", "sysroot")
	libdir = filepath.Join(prefix, "lib")
	return clang, sysroot, libdir
}

func miseWhich(t *testing.T, bin, tool string) string {
	t.Helper()
	out, err := exec.Command("mise", "which", "--tool", tool, bin).CombinedOutput()
	if err != nil {
		t.Fatalf("mise which %s (%s): %v\n%s", bin, tool, err, out)
	}
	p := strings.TrimSpace(string(out))
	if p == "" {
		t.Fatalf("mise which %s (%s): empty path", bin, tool)
	}
	return p
}

func tailBytes(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return fmt.Sprintf("…(%d bytes omitted)\n%s", len(b)-n, b[len(b)-n:])
}

func irSnippet(path string, err error) string {
	b, rdErr := os.ReadFile(path)
	if err == nil || rdErr != nil {
		return ""
	}
	base := filepath.Base(path)
	msg := err.Error()
	at := strings.Index(msg, base+":")
	if at < 0 {
		return ""
	}
	var line int
	if _, scanErr := fmt.Sscanf(msg[at+len(base)+1:], "%d:", &line); scanErr != nil || line <= 0 {
		return ""
	}
	lines := strings.Split(string(b), "\n")
	if line > len(lines) {
		return ""
	}
	lo, hi := line-3, line+2
	if lo < 1 {
		lo = 1
	}
	if hi > len(lines) {
		hi = len(lines)
	}
	var sb strings.Builder
	for i := lo; i <= hi; i++ {
		mark := "  "
		if i == line {
			mark = ">>"
		}
		fmt.Fprintf(&sb, "%s %d: %s\n", mark, i, lines[i-1])
	}
	return sb.String()
}
