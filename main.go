package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func main() {
	packageName := flag.String("package", "main", "Go package name for generated code")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: leaven [flags] input-file.ll\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	if err := validatePackageName(*packageName); err != nil {
		log.Fatal(err)
	}

	inFile := flag.Arg(0)
	m, err := parseIRFile(inFile)
	if err != nil {
		log.Fatal(err)
	}

	outFile := strings.TrimSuffix(inFile, ".ll") + ".go"
	out, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	err = compile(out, m, *packageName)
	if err != nil {
		log.Fatal(err)
	}
}

// validatePackageName checks that name is a legal Go package identifier.
func validatePackageName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: empty package name", errInvalidPackage)
	}
	r, size := utf8.DecodeRuneInString(name)
	if r == utf8.RuneError && size == 1 {
		return fmt.Errorf("%w: %q", errInvalidPackage, name)
	}
	if !unicode.IsLetter(r) && r != '_' {
		return fmt.Errorf("%w: %q", errInvalidPackage, name)
	}
	for _, r := range name[size:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return fmt.Errorf("%w: %q", errInvalidPackage, name)
		}
	}
	return nil
}

func compile(out io.Writer, m *ir.Module, packageName string) error {
	var body bytes.Buffer
	if err := writeModule(&body, m, packageName); err != nil {
		return err
	}
	fmt.Fprintf(out, "package %s\n\n", packageName)
	writeImports(out, body.String())
	_, err := io.WriteString(out, body.String())
	return err
}

// writeImports emits import declarations for packages referenced in body.
func writeImports(out io.Writer, body string) {
	type imp struct {
		path string
		hit  string // substring that indicates the package is needed
	}
	needed := []imp{
		{`"unsafe"`, "unsafe."},
		{`"math"`, "math."},
		{`"os"`, "os."},
		{`"sync/atomic"`, "atomic."},
		{`"github.com/andybalholm/leaven/libc"`, "libc."},
	}
	var paths []string
	for _, n := range needed {
		if strings.Contains(body, n.hit) {
			paths = append(paths, n.path)
		}
	}
	if len(paths) == 0 {
		return
	}
	if len(paths) == 1 {
		fmt.Fprintf(out, "import %s\n\n", paths[0])
		return
	}
	fmt.Fprintln(out, "import (")
	for _, p := range paths {
		fmt.Fprintf(out, "\t%s\n", p)
	}
	fmt.Fprint(out, ")\n\n")
}

func writeModule(out io.Writer, m *ir.Module, packageName string) error {
	for _, t := range m.TypeDefs {
		name := TypeName(t)
		if name == "" {
			continue
		}
		if strings.Contains(name, ".") {
			// It's a definition that's beeen replaced by a reference to a standard-library type.
			continue
		}

		def, err := TypeDefinition(t)
		if err != nil {
			return fmt.Errorf("error generating type definition for %v: %w", t, err)
		}

		fmt.Fprintf(out, "type %s %s\n\n", name, def)
	}

	for _, g := range m.Globals {
		if g.Init == nil {
			// Just a declaration; skip it.
			continue
		}
		t, err := TypeSpec(g.ContentType)
		if err != nil {
			return fmt.Errorf("error translating type (%v): %w", g.ContentType, err)
		}
		val, err := FormatValue(g.Init)
		if err != nil {
			return fmt.Errorf("error translating initializer (%v): %w", g.Init, err)
		}
		fmt.Fprintf(out, "var %s %s = %s\n\n", VariableName(g), t, val)
	}

	for _, f := range m.Funcs {
		if f.Blocks == nil {
			// Just a declaration, not a definition; skip it.
			continue
		}

		fixMalloc(f)

		name := VariableName(f)
		// Only package main gets a Go program entry point from C main.
		isGoMain := name == "main" && packageName == "main"

		if isGoMain {
			fmt.Fprintln(out, "func main() {")
		} else {
			fmt.Fprintf(out, "func %s(", name)
			for i, p := range f.Params {
				if i > 0 {
					fmt.Fprint(out, ", ")
				}
				pt, err := TypeSpec(p.Typ)
				if err != nil {
					return fmt.Errorf("error translating type for parameter %d of %s: %w", i, f.Name(), err)
				}
				fmt.Fprintf(out, "%s %s", VariableName(p), pt)
			}
			if f.Sig.Variadic {
				if len(f.Params) > 0 {
					fmt.Fprint(out, ", ")
				}
				fmt.Fprint(out, "varargs ...interface{}")
			}
			fmt.Fprint(out, ") ")
			rt := f.Sig.RetType
			if !types.Equal(rt, types.Void) {
				retType, err := TypeSpec(rt)
				if err != nil {
					return fmt.Errorf("error translating return type for %s: %w", f.Name(), err)
				}
				fmt.Fprintf(out, "%s ", retType)
			}
			fmt.Fprint(out, "{\n")
		}

		// Declare variables.
		vars := make(map[string][]string)
		var allVars []string
		for _, b := range f.Blocks {
			for _, inst := range b.Insts {
				if inst, ok := inst.(value.Named); ok {
					if types.Equal(inst.Type(), types.Void) {
						continue
					}
					t, err := TypeSpec(inst.Type())
					if err != nil {
						return fmt.Errorf("error translating type of %s in %s: %w", inst.Ident(), f.Name(), err)
					}
					vars[t] = append(vars[t], VariableName(inst))
					allVars = append(allVars, VariableName(inst))
				}
			}
		}
		varTypes := make([]string, 0, len(vars))
		for t := range vars {
			varTypes = append(varTypes, t)
		}
		sort.Strings(varTypes)
		for _, t := range varTypes {
			fmt.Fprintf(out, "\tvar %s %s\n", strings.Join(vars[t], ", "), t)
		}
		if len(vars) > 0 {
			fmt.Fprintln(out)
			// Get rid of unused-variable errors.
			for i := range allVars {
				if i == 0 {
					fmt.Fprint(out, "\t_")
				} else {
					fmt.Fprint(out, ", _")
				}
			}
			fmt.Fprintf(out, " = %s\n\n", strings.Join(allVars, ", "))
		}

		// Labels are only legal in Go if something gotos them. Collect jump targets.
		gotoTargets := blockGotoTargets(f)

		// Translate instructions.
		for i, b := range f.Blocks {
			// Entry always emitted; other blocks only if they are branch targets
			// (avoids "label defined and not used" and pure dead code).
			if i != 0 && !gotoTargets[b] {
				continue
			}
			if gotoTargets[b] {
				if i != 0 {
					fmt.Fprintln(out)
				}
				fmt.Fprintf(out, "%s:\n", BlockName(b))
			}
			for _, inst := range b.Insts {
				if _, ok := inst.(*ir.InstPhi); ok {
					continue
				}
				translated, err := TranslateInstruction(inst)
				if err != nil {
					return fmt.Errorf("error translating %q: %w", inst.LLString(), err)
				}
				if translated != "" {
					fmt.Fprintf(out, "\t%s\n", translated)
				}
			}
			switch term := b.Term.(type) {
			case *ir.TermBr:
				phis, err := PhiAssignments(b, term.Target)
				if err != nil {
					return fmt.Errorf("error translating phi nodes: %w", err)
				}
				if phis != "" {
					fmt.Fprintf(out, "\t%s\n", phis)
				}
				fmt.Fprintf(out, "\tgoto %s\n", BlockName(term.Target))

			case *ir.TermCondBr:
				cond, err := FormatValue(term.Cond)
				if err != nil {
					return fmt.Errorf("error translating condition (%v): %w", term.Cond, err)
				}
				fmt.Fprintf(out, "\tif %s {\n", cond)
				phis, err := PhiAssignments(b, term.TargetTrue)
				if err != nil {
					return fmt.Errorf("error translating phi nodes: %w", err)
				}
				if phis != "" {
					fmt.Fprintf(out, "\t\t%s\n", phis)
				}
				fmt.Fprintf(out, "\t\tgoto %s\n", BlockName(term.TargetTrue))
				fmt.Fprintln(out, "\t} else {")
				phis, err = PhiAssignments(b, term.TargetFalse)
				if err != nil {
					return fmt.Errorf("error translating phi nodes: %w", err)
				}
				if phis != "" {
					fmt.Fprintf(out, "\t\t%s\n", phis)
				}
				fmt.Fprintf(out, "\t\tgoto %s\n", BlockName(term.TargetFalse))
				fmt.Fprintln(out, "\t}")

			case *ir.TermRet:
				if term.X == nil {
					// void return
					if i == len(f.Blocks)-1 {
						// Just skip the return statement, since it's the end of the function anyway.
						continue
					}
					fmt.Fprintln(out, "\treturn")
				}
				retVal, err := FormatValue(term.X)
				if err != nil {
					return fmt.Errorf("error translating return value (%v): %w", term.X, err)
				}
				if isGoMain {
					fmt.Fprintf(out, "\tos.Exit(int(%s))\n", retVal)
				} else {
					fmt.Fprintf(out, "\treturn %s\n", retVal)
				}

			case *ir.TermSwitch:
				x, err := FormatValue(term.X)
				if err != nil {
					return fmt.Errorf("error translating control value (%v): %w", term.X, err)
				}
				fmt.Fprintf(out, "\tswitch %s {\n", x)
				for _, c := range term.Cases {
					x, err := FormatValue(c.X)
					if err != nil {
						return fmt.Errorf("error translating case value (%v): %w", c.X, err)
					}
					fmt.Fprintf(out, "\tcase %s:\n", x)
					phis, err := PhiAssignments(b, c.Target)
					if err != nil {
						return fmt.Errorf("error translating phi nodes: %w", err)
					}
					if phis != "" {
						fmt.Fprintf(out, "\t\t%s\n", phis)
					}
					fmt.Fprintf(out, "\t\tgoto %s\n", BlockName(c.Target))
				}
				fmt.Fprint(out, "\tdefault:\n")
				phis, err := PhiAssignments(b, term.TargetDefault)
				if err != nil {
					return fmt.Errorf("error translating phi nodes: %w", err)
				}
				if phis != "" {
					fmt.Fprintf(out, "\t\t%s\n", phis)
				}
				fmt.Fprintf(out, "\t\tgoto %s\n", BlockName(term.TargetDefault))
				fmt.Fprint(out, "\t}\n")

			case *ir.TermUnreachable:
				// LLVM unreachable: UB if executed. panic is a no-return so Go
				// typechecks functions that end only on this path.
				fmt.Fprintln(out, "\tpanic(\"unreachable\")")

			default:
				return fmt.Errorf("%w: %T", errUnsupportedTerminator, term)
			}
		}

		fmt.Fprint(out, "}\n\n")
	}
	return nil
}

// blockGotoTargets returns the set of blocks that are targets of an explicit
// branch/switch in f (i.e. need a Go label).
func blockGotoTargets(f *ir.Func) map[*ir.Block]bool {
	targets := make(map[*ir.Block]bool)
	add := func(v value.Value) {
		if b, ok := v.(*ir.Block); ok {
			targets[b] = true
		}
	}
	for _, b := range f.Blocks {
		switch term := b.Term.(type) {
		case *ir.TermBr:
			add(term.Target)
		case *ir.TermCondBr:
			add(term.TargetTrue)
			add(term.TargetFalse)
		case *ir.TermSwitch:
			add(term.TargetDefault)
			for _, c := range term.Cases {
				add(c.Target)
			}
		}
	}
	return targets
}

// PhiAssignments returns an assignment statement expressing the effects of Phi
// nodes on the branch from block a to block b. If block b has no phi nodes,
// it returns the empty string.
func PhiAssignments(a, b value.Value) (string, error) {
	var dest, src []string
	for _, inst := range b.(*ir.Block).Insts {
		phi, ok := inst.(*ir.InstPhi)
		if !ok {
			break
		}
		for _, inc := range phi.Incs {
			if inc.Pred == a {
				source, err := FormatValue(inc.X)
				if err != nil {
					return "", fmt.Errorf("error translating value (%v): %w", inc.X, err)
				}
				src = append(src, source)
				dest = append(dest, VariableName(phi))
				break
			}
		}
	}
	if len(src) == 0 {
		return "", nil
	}
	return strings.Join(dest, ", ") + " = " + strings.Join(src, ", "), nil
}
