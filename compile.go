package leaven

import (
	"fmt"
	"io"
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/lewtec/leaven/internal/llir/ir"
	"github.com/lewtec/leaven/internal/llir/ir/types"
	"github.com/lewtec/leaven/internal/llir/ir/value"
)

func Compile(out io.Writer, m *ir.Module, packageName string) error {
	f := newGoFile(packageName)
	if err := writeModule(f, m, packageName); err != nil {
		return err
	}
	return f.Render(out)
}

func writeModule(f *jen.File, m *ir.Module, packageName string) error {
	collectTaggedPointerTypes(m)
	collectModuleNames(m)

	for _, t := range m.TypeDefs {
		name := TypeName(t)
		if name == "" {
			continue
		}

		def, err := TypeDefinition(t)
		if err != nil {
			return fmt.Errorf("error generating type definition for %v: %w", t, err)
		}
		f.Type().Id(name).Add(def)
	}

	type deferredInit struct {
		name string
		val  *jen.Statement
	}
	var deferred []deferredInit
	for _, g := range m.Globals {
		name := VariableName(g)
		if g.Init == nil && hasRuntimeDef(name) {
			continue
		}
		t, err := TypeSpec(g.ContentType)
		if err != nil {
			return fmt.Errorf("error translating type (%v): %w", g.ContentType, err)
		}
		if g.Init == nil {
			f.Var().Id(name).Add(t)
			continue
		}
		val, err := FormatValue(g.Init)
		if err != nil {
			return fmt.Errorf("error translating initializer (%v): %w", g.Init, err)
		}
		// SSO self-ref (`var x = ...&x`) and C++ vtables (`var vt = dtor`
		// when dtor stores vt) are Go initialization cycles. Zero the
		// var and assign in init() after the name exists.
		if refersToGlobal(g.Init, g) || mentionsFunc(g.Init) {
			f.Var().Id(name).Add(t)
			deferred = append(deferred, deferredInit{name, val})
			continue
		}
		f.Var().Id(name).Add(t).Op("=").Add(val)
	}
	if len(deferred) > 0 {
		f.Func().Id("init").Params().BlockFunc(func(g *jen.Group) {
			for _, d := range deferred {
				g.Id(d.name).Op("=").Add(d.val)
			}
		})
	}

	for _, fn := range m.Funcs {
		name := VariableName(fn)
		if fn.Blocks == nil && hasRuntimeDef(name) {
			continue
		}

		if fn.Blocks != nil {
			fixMalloc(fn)
		}

		// Only package main gets a Go program entry point from C main.
		isGoMain := name == "main" && packageName == "main"

		var params []jen.Code
		if !isGoMain {
			for i, p := range fn.Params {
				pt, err := TypeSpec(p.Typ)
				if err != nil {
					return fmt.Errorf("error translating type for parameter %d of %s: %w", i, fn.Name(), err)
				}
				// Declares often have several unnamed ptr params; VariableName
				// would make them all arg_v0.
				pname := VariableName(p)
				if fn.Blocks == nil {
					pname = fmt.Sprintf("a%d", i)
				}
				params = append(params, jen.Id(pname).Add(pt))
			}
			if fn.Sig.Variadic {
				params = append(params, jen.Id("varargs").Op("...").Interface())
			}
		}

		var ret *jen.Statement
		if !isGoMain {
			rt := fn.Sig.RetType
			if !types.Equal(rt, types.Void) {
				retType, err := TypeSpec(rt)
				if err != nil {
					return fmt.Errorf("error translating return type for %s: %w", fn.Name(), err)
				}
				ret = retType
			}
		}

		decl := f.Func().Id(name).Params(params...)
		if ret != nil {
			decl.Add(ret)
		}
		if fn.Blocks == nil {
			decl.Block(jen.Panic(jen.Lit("unsatisfied: " + fn.Name())))
			continue
		}
		var bodyErr error
		decl.BlockFunc(func(g *jen.Group) {
			if err := writeFuncBody(g, fn, isGoMain); err != nil {
				bodyErr = err
			}
		})
		if bodyErr != nil {
			return bodyErr
		}
	}
	return nil
}

// writeCMainArgs binds C main(int argc, char **argv) from os.Args.
// argc includes argv[0], same as C.
func writeCMainArgs(g *jen.Group, fn *ir.Func) {
	if len(fn.Params) >= 1 {
		g.Id(VariableName(fn.Params[0])).Op("=").Int32().Call(jen.Len(jen.Qual("os", "Args")))
	}
	if len(fn.Params) >= 2 {
		g.Id(VariableName(fn.Params[1])).Op("=").Add(libc("Argv").Call())
	}
}

func writeFuncBody(g *jen.Group, fn *ir.Func, isGoMain bool) error {
	type varGroup struct {
		typ   *jen.Statement
		names []string
	}
	groups := make(map[string]*varGroup)
	var allVars []string
	addNamed := func(v value.Named) error {
		if types.Equal(v.Type(), types.Void) {
			return nil
		}
		t, err := TypeSpec(v.Type())
		if err != nil {
			return fmt.Errorf("error translating type of %s in %s: %w", v.Ident(), fn.Name(), err)
		}
		key := v.Type().String()
		if groups[key] == nil {
			groups[key] = &varGroup{typ: t}
		}
		groups[key].names = append(groups[key].names, VariableName(v))
		allVars = append(allVars, VariableName(v))
		return nil
	}
	if isGoMain {
		for _, p := range fn.Params {
			if err := addNamed(p); err != nil {
				return err
			}
		}
	}
	for _, b := range fn.Blocks {
		for _, inst := range b.Insts {
			n, ok := inst.(value.Named)
			if !ok {
				continue
			}
			if err := addNamed(n); err != nil {
				return err
			}
		}
		// invoke's result is the terminator, not an inst.
		if n, ok := b.Term.(value.Named); ok {
			if err := addNamed(n); err != nil {
				return err
			}
		}
	}
	varTypes := make([]string, 0, len(groups))
	for t := range groups {
		varTypes = append(varTypes, t)
	}
	sort.Strings(varTypes)
	for _, key := range varTypes {
		grp := groups[key]
		ids := make([]jen.Code, len(grp.names))
		for i, n := range grp.names {
			ids[i] = jen.Id(n)
		}
		g.Var().List(ids...).Add(grp.typ)
	}
	if len(allVars) > 0 {
		lhs := make([]jen.Code, len(allVars))
		rhs := make([]jen.Code, len(allVars))
		for i, n := range allVars {
			lhs[i] = jen.Id("_")
			rhs[i] = jen.Id(n)
		}
		g.List(lhs...).Op("=").List(rhs...)
		g.Line()
	}
	if isGoMain {
		writeCMainArgs(g, fn)
	}

	gotoTargets := blockGotoTargets(fn)

	for i, b := range fn.Blocks {
		if i != 0 && !gotoTargets[b] {
			continue
		}
		if gotoTargets[b] {
			if i != 0 {
				g.Line()
			}
			g.Id(BlockName(b)).Op(":")
		}
		for _, inst := range b.Insts {
			if _, ok := inst.(*ir.InstPhi); ok {
				continue
			}
			translated, err := TranslateInstruction(inst)
			if err != nil {
				return fmt.Errorf("error translating %q: %w", inst.LLString(), err)
			}
			for _, stmt := range translated {
				g.Add(stmt)
			}
		}
		switch term := b.Term.(type) {
		case *ir.TermBr:
			phis, err := PhiAssignments(b, term.Target)
			if err != nil {
				return fmt.Errorf("error translating phi nodes: %w", err)
			}
			if phis != nil {
				g.Add(phis)
			}
			g.Goto().Id(BlockName(term.Target))

		case *ir.TermCondBr:
			cond, err := FormatValue(term.Cond)
			if err != nil {
				return fmt.Errorf("error translating condition (%v): %w", term.Cond, err)
			}
			truePhis, err := PhiAssignments(b, term.TargetTrue)
			if err != nil {
				return fmt.Errorf("error translating phi nodes: %w", err)
			}
			falsePhis, err := PhiAssignments(b, term.TargetFalse)
			if err != nil {
				return fmt.Errorf("error translating phi nodes: %w", err)
			}
			g.If(cond).BlockFunc(func(ig *jen.Group) {
				if truePhis != nil {
					ig.Add(truePhis)
				}
				ig.Goto().Id(BlockName(term.TargetTrue))
			}).Else().BlockFunc(func(ig *jen.Group) {
				if falsePhis != nil {
					ig.Add(falsePhis)
				}
				ig.Goto().Id(BlockName(term.TargetFalse))
			})

		case *ir.TermRet:
			if term.X == nil {
				if i == len(fn.Blocks)-1 {
					continue
				}
				g.Return()
				continue
			}
			retVal, err := FormatValue(term.X)
			if err != nil {
				return fmt.Errorf("error translating return value (%v): %w", term.X, err)
			}
			if isGoMain {
				g.Qual("os", "Exit").Call(jen.Int().Call(retVal))
			} else {
				g.Return(retVal)
			}

		case *ir.TermSwitch:
			x, err := FormatValue(term.X)
			if err != nil {
				return fmt.Errorf("error translating control value (%v): %w", term.X, err)
			}
			type swCase struct {
				x      *jen.Statement
				phis   *jen.Statement
				target string
			}
			cases := make([]swCase, 0, len(term.Cases))
			for _, c := range term.Cases {
				cx, err := FormatValue(c.X)
				if err != nil {
					return fmt.Errorf("error translating case value (%v): %w", c.X, err)
				}
				phis, err := PhiAssignments(b, c.Target)
				if err != nil {
					return fmt.Errorf("error translating phi nodes: %w", err)
				}
				cases = append(cases, swCase{x: cx, phis: phis, target: BlockName(c.Target)})
			}
			defPhis, err := PhiAssignments(b, term.TargetDefault)
			if err != nil {
				return fmt.Errorf("error translating phi nodes: %w", err)
			}
			defTarget := BlockName(term.TargetDefault)
			g.Switch(x).BlockFunc(func(sg *jen.Group) {
				for _, c := range cases {
					c := c
					sg.Case(c.x).BlockFunc(func(cg *jen.Group) {
						if c.phis != nil {
							cg.Add(c.phis)
						}
						cg.Goto().Id(c.target)
					})
				}
				sg.Default().BlockFunc(func(dg *jen.Group) {
					if defPhis != nil {
						dg.Add(defPhis)
					}
					dg.Goto().Id(defTarget)
				})
			})

		case *ir.TermUnreachable:
			g.Panic(jen.Lit("unreachable"))

		case *ir.TermInvoke:
			call := &ir.InstCall{
				LocalIdent: term.LocalIdent,
				Callee:     term.Invokee,
				Args:       term.Args,
				Typ:        term.Typ,
			}
			translated, err := translateCall(call)
			if err != nil {
				return fmt.Errorf("error translating invoke: %w", err)
			}
			for _, stmt := range translated {
				g.Add(stmt)
			}
			normal, ok := term.NormalRetTarget.(*ir.Block)
			if !ok {
				return fmt.Errorf("invoke normal dest is %T", term.NormalRetTarget)
			}
			phis, err := PhiAssignments(b, normal)
			if err != nil {
				return fmt.Errorf("error translating phi nodes: %w", err)
			}
			if phis != nil {
				g.Add(phis)
			}
			g.Goto().Id(BlockName(normal))

		case *ir.TermResume:
			g.Panic(jen.Lit("resume"))

		default:
			return fmt.Errorf("%w: %T", errUnsupportedTerminator, term)
		}
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
		case *ir.TermInvoke:
			add(term.NormalRetTarget)
			add(term.ExceptionRetTarget)
		}
	}
	return targets
}

// PhiAssignments returns an assignment statement expressing the effects of Phi
// nodes on the branch from block a to block b. If block b has no phi nodes,
// it returns nil.
func PhiAssignments(a, b value.Value) (*jen.Statement, error) {
	var dest, src []jen.Code
	for _, inst := range b.(*ir.Block).Insts {
		phi, ok := inst.(*ir.InstPhi)
		if !ok {
			break
		}
		for _, inc := range phi.Incs {
			if inc.Pred == a {
				source, err := FormatValue(inc.X)
				if err != nil {
					return nil, fmt.Errorf("error translating value (%v): %w", inc.X, err)
				}
				src = append(src, source)
				dest = append(dest, jen.Id(VariableName(phi)))
				break
			}
		}
	}
	if len(src) == 0 {
		return nil, nil
	}
	return jen.List(dest...).Op("=").List(src...), nil
}
