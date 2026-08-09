// Command genir emits testdata/ir/*/input.<LLVM-major>.ll from source.{c,cpp,rs}.
// Compilers run under mise exec (pins in mise.toml, or -c/-cxx / leaven:tool=).
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	handIRMarker = "leaven:hand-ir"
	toolPrefix   = "leaven:tool="
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("genir", flag.ContinueOnError)
	check := fs.Bool("check", false, "regenerate to temp and fail if committed IR is stale")
	cSpec := fs.String("c", "", "mise spec for clang (e.g. conda:clang@22.1.8)")
	cxxSpec := fs.String("cxx", "", "mise spec for clang++")
	if err := fs.Parse(args); err != nil {
		return err
	}

	root, err := moduleRoot()
	if err != nil {
		return err
	}
	irRoot := filepath.Join(root, "testdata", "ir")

	fixtures, err := listFixtures(irRoot)
	if err != nil {
		return err
	}

	need := map[string]bool{}
	for _, f := range fixtures {
		if f.hand {
			continue
		}
		need[f.bin] = true
	}
	cfg := config{root: root, cSpec: *cSpec, cxxSpec: *cxxSpec}
	for _, bin := range []string{"clang", "clang++", "rustc"} {
		if !need[bin] {
			continue
		}
		spec := ""
		switch bin {
		case "clang":
			spec = cfg.cSpec
		case "clang++":
			spec = cfg.cxxSpec
		}
		if err := preflight(cfg, spec, bin); err != nil {
			return err
		}
		ver, err := compilerBanner(cfg, spec, bin)
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n", bin, firstLine(ver))
	}

	var tmp string
	if *check {
		tmp, err = os.MkdirTemp("", "leaven-genir-*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmp)
	}

	regen, skip, stale := 0, 0, 0
	for _, f := range fixtures {
		if f.hand {
			fmt.Printf("skip %s (hand-ir)\n", f.name)
			skip++
			continue
		}
		major, err := irMajor(cfg, f)
		if err != nil {
			return fmt.Errorf("%s: %w", f.name, err)
		}
		dest := filepath.Join(f.dir, "input."+major+".ll")
		if *check {
			out := filepath.Join(tmp, f.name+"."+major+".ll")
			if err := compileOne(cfg, f, out); err != nil {
				return fmt.Errorf("%s: %w", f.name, err)
			}
			same, err := sameFile(dest, out)
			if err != nil {
				return err
			}
			if !same {
				fmt.Printf("stale %s input.%s.ll (go generate .)\n", f.name, major)
				stale++
			} else {
				fmt.Printf("ok   %s input.%s.ll\n", f.name, major)
			}
			continue
		}
		fmt.Printf("gen  %s input.%s.ll\n", f.name, major)
		if err := compileOne(cfg, f, dest); err != nil {
			return fmt.Errorf("%s: %w", f.name, err)
		}
		_ = os.Remove(filepath.Join(f.dir, "input.ll"))
		regen++
	}

	if *check {
		fmt.Printf("stale=%d skipped=%d\n", stale, skip)
		if stale > 0 {
			return fmt.Errorf("%d stale fixture(s)", stale)
		}
		return nil
	}
	fmt.Printf("regenerated=%d skipped=%d\n", regen, skip)
	return nil
}

type config struct {
	root    string
	cSpec   string
	cxxSpec string
}

type fixture struct {
	dir, name, src, bin string
	hand                bool
}

func listFixtures(irRoot string) ([]fixture, error) {
	ents, err := os.ReadDir(irRoot)
	if err != nil {
		return nil, err
	}
	var out []fixture
	for _, e := range ents {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(irRoot, e.Name())
		src, err := findSource(dir)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		body, err := os.ReadFile(src)
		if err != nil {
			return nil, err
		}
		bin, err := compilerBin(src)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		out = append(out, fixture{
			dir:  dir,
			name: e.Name(),
			src:  src,
			bin:  bin,
			hand: bytes.Contains(body, []byte(handIRMarker)),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out, nil
}

func findSource(dir string) (string, error) {
	var hits []string
	for _, name := range []string{"source.c", "source.cpp", "source.rs"} {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err == nil {
			hits = append(hits, p)
		}
	}
	if len(hits) != 1 {
		return "", fmt.Errorf("need exactly one of source.c, source.cpp, source.rs (found %d)", len(hits))
	}
	return hits[0], nil
}

func compilerBin(src string) (string, error) {
	switch filepath.Ext(src) {
	case ".c":
		return "clang", nil
	case ".cpp":
		return "clang++", nil
	case ".rs":
		return "rustc", nil
	default:
		return "", fmt.Errorf("unknown source %s", src)
	}
}

func specFor(cfg config, src, bin string) string {
	if s := toolOverride(src); s != "" {
		return s
	}
	switch bin {
	case "clang":
		return cfg.cSpec
	case "clang++":
		return cfg.cxxSpec
	default:
		return ""
	}
}

func toolOverride(src string) string {
	b, err := os.ReadFile(src)
	if err != nil {
		return ""
	}
	for _, line := range bytes.Split(b, []byte("\n")) {
		if i := bytes.Index(line, []byte(toolPrefix)); i >= 0 {
			fields := strings.Fields(string(line[i+len(toolPrefix):]))
			if len(fields) == 0 {
				return ""
			}
			return fields[0]
		}
	}
	return ""
}

func compileOne(cfg config, f fixture, out string) error {
	absSrc, err := filepath.Abs(f.src)
	if err != nil {
		return err
	}
	var args []string
	switch f.bin {
	case "clang":
		args = []string{"clang", "-S", "-emit-llvm", "-fno-discard-value-names", "-std=gnu11", "-O0", "-o", out, absSrc}
	case "clang++":
		args = []string{"clang++", "-S", "-emit-llvm", "-fno-discard-value-names", "-std=c++17", "-O0", "-o", out, absSrc}
	case "rustc":
		args = []string{"rustc", "--emit=llvm-ir", "-C", "opt-level=0", "-C", "debuginfo=0", "--crate-type=lib", "-o", out, absSrc}
	default:
		return fmt.Errorf("unknown compiler %s", f.bin)
	}
	return runMise(cfg.root, specFor(cfg, f.src, f.bin), args...)
}

func irMajor(cfg config, f fixture) (string, error) {
	switch f.bin {
	case "clang", "clang++":
		out, err := miseOutput(cfg.root, specFor(cfg, f.src, f.bin), f.bin, "--version")
		if err != nil {
			return "", err
		}
		return clangMajor(out)
	case "rustc":
		out, err := miseOutput(cfg.root, specFor(cfg, f.src, f.bin), "rustc", "-vV")
		if err != nil {
			return "", err
		}
		return rustcLLVMMajor(out)
	default:
		return "", fmt.Errorf("unknown compiler %s", f.bin)
	}
}

func preflight(cfg config, spec, bin string) error {
	if err := runMise(cfg.root, spec, bin, "--version"); err != nil {
		return fmt.Errorf("mise exec -- %s failed (mise install; conda pin in mise.toml): %w", bin, err)
	}
	return nil
}

func compilerBanner(cfg config, spec, bin string) (string, error) {
	if bin == "rustc" {
		ver, err := miseOutput(cfg.root, spec, "rustc", "--version")
		if err != nil {
			return "", err
		}
		vv, err := miseOutput(cfg.root, spec, "rustc", "-vV")
		if err != nil {
			return "", err
		}
		maj, err := rustcLLVMMajor(vv)
		if err != nil {
			return strings.TrimSpace(ver), nil
		}
		return strings.TrimSpace(ver) + "  LLVM " + maj, nil
	}
	return miseOutput(cfg.root, spec, bin, "--version")
}

func runMise(root, spec string, args ...string) error {
	cmd := miseCmd(root, spec, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func miseOutput(root, spec string, args ...string) (string, error) {
	cmd := miseCmd(root, spec, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %w\n%s", strings.Join(cmd.Args, " "), err, out)
	}
	return string(out), nil
}

func miseCmd(root, spec string, args ...string) *exec.Cmd {
	a := []string{"exec", "-C", root}
	if spec != "" {
		a = append(a, spec)
	}
	a = append(a, "--")
	a = append(a, args...)
	cmd := exec.Command("mise", a...)
	cmd.Dir = root
	return cmd
}

func moduleRoot() (string, error) {
	if dir, err := os.Getwd(); err == nil {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
	}
	cmd := exec.Command("go", "env", "GOMOD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	mod := strings.TrimSpace(string(out))
	if mod == "" || mod == os.DevNull {
		return "", fmt.Errorf("not in a module")
	}
	return filepath.Dir(mod), nil
}

func sameFile(a, b string) (bool, error) {
	da, err := os.ReadFile(a)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	db, err := os.ReadFile(b)
	if err != nil {
		return false, err
	}
	return bytes.Equal(da, db), nil
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return strings.TrimSpace(s)
}

var (
	clangVerRe = regexp.MustCompile(`clang version ([0-9]+)`)
	rustcLLRe  = regexp.MustCompile(`(?m)^LLVM version: ([0-9]+)`)
)

func clangMajor(versionOut string) (string, error) {
	m := clangVerRe.FindStringSubmatch(versionOut)
	if len(m) != 2 {
		return "", fmt.Errorf("could not parse clang major from %q", firstLine(versionOut))
	}
	return m[1], nil
}

func rustcLLVMMajor(versionOut string) (string, error) {
	m := rustcLLRe.FindStringSubmatch(versionOut)
	if len(m) != 2 {
		return "", fmt.Errorf("could not parse rustc LLVM major from version output")
	}
	return m[1], nil
}
