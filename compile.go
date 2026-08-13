package leaven

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/lewtec/leaven/internal/llir/ir"
	"github.com/lewtec/leaven/internal/llir/ir/constant"
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
	var streamInits []string
	for _, g := range m.Globals {
		if isLLVMSpecialGlobal(g.Name()) {
			continue
		}
		name := VariableName(g)
		if g.Init == nil && hasRuntimeDef(name) {
			continue
		}
		t, err := TypeSpec(g.ContentType)
		if err != nil {
			return fmt.Errorf("error translating type (%v): %w", g.ContentType, err)
		}
		if g.Init == nil {
			if init := vttStandin(name, g.ContentType); init != nil {
				f.Var().Id(name).Add(t).Op("=").Add(init)
				continue
			}
			f.Var().Id(name).Add(t)
			if isStdStream(name) {
				streamInits = append(streamInits, name)
			}
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
	ctors := globalCtorFuncs(m)
	if len(deferred) > 0 || len(ctors) > 0 || len(streamInits) > 0 {
		f.Func().Id("init").Params().BlockFunc(func(g *jen.Group) {
			for _, d := range deferred {
				g.Id(d.name).Op("=").Add(d.val)
			}
			for _, name := range streamInits {
				g.Add(initStdStream(name))
			}
			for _, fn := range ctors {
				g.Id(VariableName(fn)).Call()
			}
		})
	}

	for _, fn := range m.Funcs {
		collectFuncLocalNames(fn)
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
			decl.Block(jen.Panic(jen.Lit(unsatisfiedMsg(fn.Name()))))
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

func isLLVMSpecialGlobal(name string) bool {
	switch name {
	case "llvm.global_ctors", "llvm.global_dtors", "llvm.used", "llvm.compiler.used":
		return true
	default:
		return false
	}
}

// globalCtorFuncs is llvm.global_ctors, lowest priority first.
func globalCtorFuncs(m *ir.Module) []*ir.Func {
	var g *ir.Global
	for _, cand := range m.Globals {
		if cand.Name() == "llvm.global_ctors" {
			g = cand
			break
		}
	}
	if g == nil || g.Init == nil {
		return nil
	}
	arr, ok := g.Init.(*constant.Array)
	if !ok {
		return nil
	}
	type ent struct {
		prio int64
		fn   *ir.Func
	}
	var ents []ent
	for _, el := range arr.Elems {
		st, ok := el.(*constant.Struct)
		if !ok || len(st.Fields) < 2 {
			continue
		}
		prio := int64(65535)
		if n, ok := st.Fields[0].(*constant.Int); ok {
			prio = n.X.Int64()
		}
		fn, ok := st.Fields[1].(*ir.Func)
		if !ok || fn.Blocks == nil {
			continue
		}
		ents = append(ents, ent{prio, fn})
	}
	sort.SliceStable(ents, func(i, j int) bool { return ents[i].prio < ents[j].prio })
	out := make([]*ir.Func, len(ents))
	for i, e := range ents {
		out[i] = e.fn
	}
	return out
}

// vttStandin fills a declare-only Itanium VTT with StandinVptr so
// inlined dtors can load *(vptr-24) without faulting on nil.
func vttStandin(name string, t types.Type) *jen.Statement {
	if !strings.HasPrefix(name, "_ZTT") {
		return nil
	}
	n, wrap := vttLen(t)
	if n <= 0 {
		return nil
	}
	elems := make([]jen.Code, n)
	for i := range elems {
		elems[i] = libc("StandinVptr").Call()
	}
	arr := jen.Index(litUntyped(int64(n))).Qual("unsafe", "Pointer").Values(elems...)
	if wrap {
		return jen.Values(arr)
	}
	return arr
}

func vttLen(t types.Type) (n int, wrap bool) {
	switch t := t.(type) {
	case *types.ArrayType:
		if _, ok := t.ElemType.(*types.PointerType); ok {
			return int(t.Len), false
		}
	case *types.StructType:
		if len(t.Fields) == 1 {
			if n, wrap := vttLen(t.Fields[0]); n > 0 && !wrap {
				return n, true
			}
		}
	}
	return 0, false
}

// unsatisfiedMsg keeps the panic line short enough that tailBytes / CI
// still show it. Itanium names alone are hundreds of bytes.
func unsatisfiedMsg(name string) string {
	const max = 80
	if len(name) <= max {
		return "unsatisfied: " + name
	}
	return "unsatisfied: " + name[:max]
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
			// cmpxchg temps (name_ok, name_old) must be function-scoped;
			// mid-function := is illegal when other blocks goto past them.
			if cx, ok := inst.(*ir.InstCmpXchg); ok {
				elem, err := TypeSpec(cx.Cmp.Type())
				if err != nil {
					return err
				}
				base := VariableName(cx)
				addExtra := func(name string, t *jen.Statement) {
					key := "extra:" + name
					if groups[key] == nil {
						groups[key] = &varGroup{typ: t}
					}
					groups[key].names = append(groups[key].names, name)
					allVars = append(allVars, name)
				}
				addExtra(base+"_ok", jen.Bool())
				addExtra(base+"_old", elem)
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

	reachable := reachableBlocks(fn)
	gotoTargets := blockGotoTargetsFrom(fn, reachable)

	for i, b := range fn.Blocks {
		// Skip unreachable blocks so their successors are not left with
		// unused labels (Go rejects "label defined and not used").
		if i != 0 && !reachable[b] {
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
	return blockGotoTargetsFrom(f, reachableBlocks(f))
}

// reachableBlocks is the set of blocks reachable from the entry via
// terminators. Unreachable blocks are not emitted.
func reachableBlocks(f *ir.Func) map[*ir.Block]bool {
	if len(f.Blocks) == 0 {
		return nil
	}
	reach := make(map[*ir.Block]bool)
	queue := []*ir.Block{f.Blocks[0]}
	reach[f.Blocks[0]] = true
	for len(queue) > 0 {
		b := queue[0]
		queue = queue[1:]
		var succ []value.Value
		switch term := b.Term.(type) {
		case *ir.TermBr:
			succ = []value.Value{term.Target}
		case *ir.TermCondBr:
			succ = []value.Value{term.TargetTrue, term.TargetFalse}
		case *ir.TermSwitch:
			succ = []value.Value{term.TargetDefault}
			for _, c := range term.Cases {
				succ = append(succ, c.Target)
			}
		case *ir.TermInvoke:
			// writeFuncBody only emits the normal edge (no landing pads).
			succ = []value.Value{term.NormalRetTarget}
		}
		for _, v := range succ {
			nb, ok := v.(*ir.Block)
			if !ok || reach[nb] {
				continue
			}
			reach[nb] = true
			queue = append(queue, nb)
		}
	}
	return reach
}

func blockGotoTargetsFrom(f *ir.Func, from map[*ir.Block]bool) map[*ir.Block]bool {
	targets := make(map[*ir.Block]bool)
	add := func(v value.Value) {
		if b, ok := v.(*ir.Block); ok {
			targets[b] = true
		}
	}
	for _, b := range f.Blocks {
		if from != nil && !from[b] {
			continue
		}
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
			// Match writeFuncBody: only the normal edge is a real goto.
			add(term.NormalRetTarget)
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
