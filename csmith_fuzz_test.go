//go:build csmith

package leaven

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// errRunTimeout is returned when a subprocess exceeds CSMITH_TIMEOUT.
var errRunTimeout = errors.New("run timeout")

// Csmith differential fuzzing of leaven.
//
// Build tag: csmith (normal `go test` skips this file).
//
//	mise install   # clang 14, prebuilt csmith, go, goimports
//	go test -tags=csmith -run TestCsmithFixedSeeds -v -count=1
//	go test -tags=csmith -fuzz=FuzzCsmith -fuzztime=30s
//
// Required tools on PATH (via mise.toml; overrides optional):
//
//	csmith / clang / goimports / go
//
// leaven is invoked with `go run ./cmd/leaven` (no install).
//
// Optional env:
//
//	CSMITH            – absolute path to the csmith binary
//	CSMITH_HOME       – install prefix (bin/csmith + include…/csmith.h)
//	CSMITH_ITERS      – random seeds for TestCsmithRandom (default 5)
//	CSMITH_TIMEOUT    – per-binary run timeout (default 5s)
//	CSMITH_EXTRA_ARGS – extra csmith flags (space-separated)
//	CLANG             – clang binary (default "clang")

func TestCsmithFixedSeeds(t *testing.T) {
	tools := requireCsmithTools(t)

	// Small fixed corpus for regression / smoke. Seeds are just numbers;
	// failures are expected until leaven covers more LLVM.
	// Known interesting seeds. seed=1 currently passes end-to-end; others are
	// retained as regressions when leaven grows. -run uses regex: prefer
	// -run 'TestCsmithFixedSeeds/seed=1$' for a single case.
	seeds := []uint64{1, 42, 100, 12345}
	for _, seed := range seeds {
		t.Run(fmt.Sprintf("seed=%d", seed), func(t *testing.T) {
			runCsmithCase(t, tools, seed)
		})
	}
}

func TestCsmithRandom(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping random csmith under -short")
	}
	tools := requireCsmithTools(t)

	iters := 5
	if v := os.Getenv("CSMITH_ITERS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			t.Fatalf("CSMITH_ITERS must be a positive integer, got %q", v)
		}
		iters = n
	}

	// Deterministic-ish base from time so re-runs cover new ground,
	// while each subtest is still reproducible from its seed name.
	base := uint64(time.Now().UnixNano())
	for i := 0; i < iters; i++ {
		seed := base + uint64(i)*0x9E3779B97F4A7C15
		t.Run(fmt.Sprintf("seed=%d", seed), func(t *testing.T) {
			runCsmithCase(t, tools, seed)
		})
	}
}

// FuzzCsmith feeds seeds into the same pipeline. Run with:
//
//	go test -tags=csmith -fuzz=FuzzCsmith -fuzztime=30s
func FuzzCsmith(f *testing.F) {
	for _, s := range []uint64{1, 2, 3, 7, 42, 99} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, seed uint64) {
		tools := requireCsmithTools(t)
		runCsmithCase(t, tools, seed)
	})
}

type csmithTools struct {
	csmith  string
	include string
	clang   string
	timeout time.Duration
}

func requireCsmithTools(t *testing.T) csmithTools {
	t.Helper()

	csmith, include, err := resolveCsmith()
	if err != nil {
		t.Skip(err.Error())
	}

	clang := os.Getenv("CLANG")
	if clang == "" {
		clang = "clang"
	}
	if _, err := exec.LookPath(clang); err != nil {
		t.Skipf("%s not on PATH (install via mise: conda:clang@14): %v", clang, err)
	}

	if _, err := exec.LookPath("goimports"); err != nil {
		t.Skip("goimports not on PATH")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go not on PATH")
	}

	timeout := 5 * time.Second
	if v := os.Getenv("CSMITH_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			t.Fatalf("CSMITH_TIMEOUT: %v", err)
		}
		timeout = d
	}

	return csmithTools{
		csmith:  csmith,
		include: include,
		clang:   clang,
		timeout: timeout,
	}
}

// resolveCsmith finds the binary and the directory that contains csmith.h.
// Order: CSMITH, CSMITH_HOME, then PATH (mise puts github:…/csmith-bin there).
// Include layouts vary: plain include/ vs include/csmith-<ver>/.
func resolveCsmith() (csmithBin, includeDir string, err error) {
	if p := os.Getenv("CSMITH"); p != "" {
		if !filepath.IsAbs(p) {
			return "", "", fmt.Errorf("%w: CSMITH must be an absolute path, got %q", errCsmithConfig, p)
		}
		if st, e := os.Stat(p); e != nil || st.IsDir() {
			return "", "", fmt.Errorf("CSMITH=%q is not a usable binary: %w", p, e)
		}
		home := filepath.Dir(filepath.Dir(p))
		if envHome := os.Getenv("CSMITH_HOME"); envHome != "" {
			home = envHome
		}
		inc, e := findCsmithInclude(home)
		if e != nil {
			return "", "", e
		}
		return p, inc, nil
	}

	if home := os.Getenv("CSMITH_HOME"); home != "" {
		if !filepath.IsAbs(home) {
			return "", "", fmt.Errorf("%w: CSMITH_HOME must be an absolute path, got %q", errCsmithConfig, home)
		}
		bin := filepath.Join(home, "bin", "csmith")
		if st, e := os.Stat(bin); e != nil || st.IsDir() {
			return "", "", fmt.Errorf("csmith binary not found at %s: %w", bin, e)
		}
		inc, e := findCsmithInclude(home)
		if e != nil {
			return "", "", e
		}
		return bin, inc, nil
	}

	p, e := exec.LookPath("csmith")
	if e != nil {
		return "", "", fmt.Errorf("csmith not on PATH (run: mise install): %w", e)
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", "", fmt.Errorf("resolve csmith path %q: %w", p, err)
	}
	home := filepath.Dir(filepath.Dir(abs))
	inc, e := findCsmithInclude(home)
	if e != nil {
		return "", "", fmt.Errorf("csmith on PATH (%s) but %w", abs, e)
	}
	return abs, inc, nil
}

// findCsmithInclude locates the directory containing csmith.h under prefix.
func findCsmithInclude(prefix string) (string, error) {
	candidates := []string{
		filepath.Join(prefix, "include"),
	}
	matches, err := filepath.Glob(filepath.Join(prefix, "include", "csmith*"))
	if err != nil {
		return "", fmt.Errorf("glob csmith include dirs: %w", err)
	}
	candidates = append(candidates, matches...)
	for _, dir := range candidates {
		if st, err := os.Stat(filepath.Join(dir, "csmith.h")); err == nil && !st.IsDir() {
			return dir, nil
		}
	}
	// One-level walk under include/ for versioned layouts.
	incRoot := filepath.Join(prefix, "include")
	entries, err := os.ReadDir(incRoot)
	if err == nil {
		for _, ent := range entries {
			if !ent.IsDir() {
				continue
			}
			dir := filepath.Join(incRoot, ent.Name())
			if st, err := os.Stat(filepath.Join(dir, "csmith.h")); err == nil && !st.IsDir() {
				return dir, nil
			}
		}
	}
	return "", fmt.Errorf("%w: csmith.h not found under %s/include", errCsmithConfig, prefix)
}

// defaultCsmithArgs keeps programs smaller / closer to what leaven can handle.
// Full-featured csmith is still reachable via CSMITH_EXTRA_ARGS (space-separated).
func defaultCsmithArgs(seed uint64, outFile string) []string {
	args := []string{
		"--seed", strconv.FormatUint(seed, 10),
		"--output", outFile,
		// Size caps
		"--max-funcs", "3",
		"--max-block-depth", "3",
		"--max-block-size", "3",
		"--max-expr-complexity", "4",
		"--max-array-dim", "2",
		"--max-array-len-per-dim", "4",
		"--max-struct-fields", "4",
		"--max-union-fields", "3",
		// Features that tend to blow up IR or hit missing leaven ops early
		"--no-bitfields",
		"--no-volatiles",
		"--no-volatile-pointers",
		"--no-packed-struct",
		"--no-builtins",
		"--no-argc",
	}
	if extra := strings.TrimSpace(os.Getenv("CSMITH_EXTRA_ARGS")); extra != "" {
		args = append(args, strings.Fields(extra)...)
	}
	return args
}

func runCsmithCase(t *testing.T, tools csmithTools, seed uint64) {
	t.Helper()

	dir := t.TempDir()
	cFile := filepath.Join(dir, "prog.c")
	nativeBin := filepath.Join(dir, "prog_c")
	llFile := filepath.Join(dir, "prog.ll")
	goFile := filepath.Join(dir, "prog.go")

	// 1. Generate
	gen := exec.Command(tools.csmith, defaultCsmithArgs(seed, cFile)...)
	if out, err := gen.CombinedOutput(); err != nil {
		t.Fatalf("csmith seed=%d: %v\n%s", seed, err, out)
	}

	// 2. Native reference binary
	clangNative := exec.Command(tools.clang, append(clangSysrootFlags(),
		"-O0",
		"-I"+tools.include,
		"-o", nativeBin,
		cFile,
	)...)
	if out, err := clangNative.CombinedOutput(); err != nil {
		t.Fatalf("clang native seed=%d: %v\n%s", seed, err, out)
	}

	nativeOut, err := runWithTimeout(nativeBin, tools.timeout)
	if err != nil {
		// Infinite/very long loops in csmith programs are not leaven bugs.
		if errors.Is(err, errRunTimeout) {
			t.Skipf("native run seed=%d: %v (skip; not a leaven failure)", seed, err)
		}
		t.Fatalf("native run seed=%d: %v\n%s", seed, err, nativeOut)
	}

	// 3. Emit LLVM IR
	clangLL := exec.Command(tools.clang, append(clangSysrootFlags(),
		"-O0",
		"-S", "-emit-llvm",
		"-fno-discard-value-names",
		"-I"+tools.include,
		"-o", llFile,
		cFile,
	)...)
	if out, err := clangLL.CombinedOutput(); err != nil {
		t.Fatalf("clang -emit-llvm seed=%d: %v\n%s", seed, err, out)
	}

	// 4. Transpile with leaven
	leaven := exec.Command("go", "run", "./cmd/leaven", llFile)
	if out, err := leaven.CombinedOutput(); err != nil {
		t.Fatalf("leaven seed=%d: %v\n%s\n(source kept under temp dir; re-run with -count=1 -v)", seed, err, out)
	}

	// 5. Format and run Go
	goimports := exec.Command("goimports", "-w", goFile)
	if out, err := goimports.CombinedOutput(); err != nil {
		t.Fatalf("goimports seed=%d: %v\n%s", seed, err, out)
	}

	// Run from module root so imports like github.com/lewtec/leaven/libc resolve.
	goRun := exec.Command("go", "run", goFile)
	goOut, err := runCmdWithTimeout(goRun, tools.timeout)
	if err != nil {
		t.Fatalf("go run seed=%d: %v\n%s", seed, err, goOut)
	}

	if !bytes.Equal(goOut, nativeOut) {
		t.Fatalf("output mismatch seed=%d\nC:  %q\nGo: %q", seed, nativeOut, goOut)
	}
}

func runWithTimeout(bin string, d time.Duration) ([]byte, error) {
	return runCmdWithTimeout(exec.Command(bin), d)
}

func runCmdWithTimeout(cmd *exec.Cmd, d time.Duration) ([]byte, error) {
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	// Unbuffered: buffered "chan error, 1" trips go/err-arg-is-not-last.
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return buf.Bytes(), err
	case <-time.After(d):
		if cmd.Process != nil {
			if killErr := cmd.Process.Kill(); killErr != nil {
				<-done
				return buf.Bytes(), fmt.Errorf("%w after %s: kill: %w", errRunTimeout, d, killErr)
			}
		}
		<-done // drain Wait; avoid leaking the process table entry
		return buf.Bytes(), fmt.Errorf("%w after %s", errRunTimeout, d)
	}
}
