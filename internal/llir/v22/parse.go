package v22

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lewtec/leaven/internal/llir/ir"
	"github.com/lewtec/leaven/internal/llir/ir/constant"
	"github.com/lewtec/leaven/internal/llir/ir/enum"
	"github.com/lewtec/leaven/internal/llir/ir/types"
	"github.com/lewtec/leaven/internal/llir/ir/value"
)

// ParseString parses LLVM 15+ textual IR (opaque ptr) into the shared IR model.
func ParseString(path, content string) (*ir.Module, error) {
	toks, err := lex(content)
	if err != nil {
		if path != "" {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		return nil, err
	}
	p := &parser{
		path:    path,
		toks:    toks,
		m:       ir.NewModule(),
		types:   make(map[string]types.Type),
		globals: make(map[string]*ir.Global),
		funcs:   make(map[string]*ir.Func),
		aliases: make(map[string]*ir.Alias),
		ptr:     types.NewOpaquePointer(),
	}
	if len(toks) > 0 {
		p.tok = toks[0]
	}
	p.predeclareFuncs()
	p.predeclareGlobals()
	if err := p.parseModule(); err != nil {
		return nil, err
	}
	return p.m, nil
}

type parser struct {
	path    string
	toks    []token
	i       int
	m       *ir.Module
	types   map[string]types.Type
	globals map[string]*ir.Global
	funcs   map[string]*ir.Func
	aliases map[string]*ir.Alias
	ptr     *types.PointerType
	tok     token

	// per-function
	fn        *ir.Func
	locals    map[string]value.Value
	blocks    map[string]*ir.Block
	phis      []pendingPhi
	nextLocal int64
}

type pendingPhi struct {
	inst *ir.InstPhi
	incs []pendingInc
}

type pendingInc struct {
	val string // local name, or "" if const already set
	c   value.Value
	bb  string
}

func (p *parser) predeclareFuncs() {
	for i := 0; i < len(p.toks); i++ {
		if p.toks[i].kind != kIdent {
			continue
		}
		switch p.toks[i].s {
		case "define", "declare":
		default:
			continue
		}
		for j := i + 1; j < len(p.toks) && j < i+120; j++ {
			if p.toks[j].kind == kGlobal {
				name := p.toks[j].s
				if _, ok := p.funcs[name]; !ok {
					f := ir.NewFunc(name, types.Void)
					p.funcs[name] = f
					p.m.Funcs = append(p.m.Funcs, f)
				}
				break
			}
		}
	}
}

func (p *parser) predeclareGlobals() {
	for i := 0; i < len(p.toks); i++ {
		if p.toks[i].kind != kGlobal {
			continue
		}
		if i+1 >= len(p.toks) || p.toks[i+1].kind != kEq {
			continue
		}
		name := p.toks[i].s
		kind := ""
		for j := i + 2; j < len(p.toks) && j < i+24; j++ {
			if p.toks[j].kind != kIdent {
				break
			}
			switch p.toks[j].s {
			case "global", "constant":
				kind = "global"
			case "alias", "ifunc":
				kind = "alias"
			}
			if kind != "" {
				break
			}
		}
		if kind != "global" {
			continue
		}
		if _, ok := p.globals[name]; ok {
			continue
		}
		if _, ok := p.funcs[name]; ok {
			continue
		}
		g := ir.NewGlobal(name, types.I8)
		p.globals[name] = g
		p.m.Globals = append(p.m.Globals, g)
	}
}

func (p *parser) parseModule() error {
	for !p.done() {
		if p.tok.kind == kEOF {
			return nil
		}
		if err := p.parseTop(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseTop() error {
	t := p.tok
	switch t.kind {
	case kIdent:
		switch t.s {
		case "source_filename":
			p.next()
			if err := p.expect(kEq); err != nil {
				return err
			}
			s, err := p.stringLit()
			if err != nil {
				return err
			}
			p.m.SourceFilename = s
			return nil
		case "target":
			p.next()
			which, err := p.ident()
			if err != nil {
				return err
			}
			if err := p.expect(kEq); err != nil {
				return err
			}
			s, err := p.stringLit()
			if err != nil {
				return err
			}
			switch which {
			case "datalayout":
				p.m.DataLayout = s
			case "triple":
				p.m.TargetTriple = s
			}
			return nil
		case "define", "declare":
			return p.parseFunc(t.s == "define")
		case "module":
			p.next()
			if err := p.wantIdent("asm"); err != nil {
				return err
			}
			s, err := p.stringLit()
			if err != nil {
				return err
			}
			p.m.ModuleAsms = append(p.m.ModuleAsms, s)
			return nil
		case "attributes":
			return p.parseAttrGroup()
		default:
			if strings.HasPrefix(t.s, "$") {
				return p.skipComdatDef()
			}
			return p.errorf("unexpected %q", t.s)
		}
	case kLocal:
		return p.parseTypeDef()
	case kGlobal:
		return p.parseGlobal()
	case kMetaID, kMetaName, kBang:
		return p.skipMetadataDef()
	default:
		return p.errorf("unexpected token %s", t)
	}
}

func (p *parser) parseTypeDef() error {
	name := p.tok.s
	p.next()
	if err := p.expect(kEq); err != nil {
		return err
	}
	if err := p.wantIdent("type"); err != nil {
		return err
	}
	st := p.namedStruct(name)
	if p.isIdent("opaque") {
		p.next()
		st.Opaque = true
		return nil
	}
	packed := false
	if p.tok.kind == kLt {
		p.next()
		packed = true
	}
	if err := p.expect(kLBrace); err != nil {
		return err
	}
	var fields []types.Type
	if p.tok.kind != kRBrace {
		for {
			t, err := p.parseType()
			if err != nil {
				return err
			}
			fields = append(fields, t)
			if p.tok.kind != kComma {
				break
			}
			p.next()
		}
	}
	if err := p.expect(kRBrace); err != nil {
		return err
	}
	if packed {
		if err := p.expect(kGt); err != nil {
			return err
		}
		st.Packed = true
	}
	st.Fields = fields
	st.Opaque = false
	return nil
}

func (p *parser) namedStruct(name string) *types.StructType {
	if t, ok := p.types[name]; ok {
		if st, ok := t.(*types.StructType); ok {
			return st
		}
	}
	st := &types.StructType{TypeName: name}
	p.types[name] = st
	p.m.TypeDefs = append(p.m.TypeDefs, st)
	return st
}

func (p *parser) parseGlobal() error {
	name := p.tok.s
	p.next()
	if err := p.expect(kEq); err != nil {
		return err
	}
	var (
		link    enum.Linkage
		pre     enum.Preemption
		vis     enum.Visibility
		unnamed enum.UnnamedAddr
		tls     enum.TLSModel
	)
	for p.tok.kind == kIdent {
		switch p.tok.s {
		case "private":
			link = enum.LinkagePrivate
			p.next()
		case "internal":
			link = enum.LinkageInternal
			p.next()
		case "external":
			link = enum.LinkageExternal
			p.next()
		case "extern_weak":
			link = enum.LinkageExternWeak
			p.next()
		case "common":
			link = enum.LinkageCommon
			p.next()
		case "appending":
			link = enum.LinkageAppending
			p.next()
		case "linkonce_odr":
			link = enum.LinkageLinkOnceODR
			p.next()
		case "weak_odr":
			link = enum.LinkageWeakODR
			p.next()
		case "weak":
			link = enum.LinkageWeak
			p.next()
		case "dso_local":
			pre = enum.PreemptionDSOLocal
			p.next()
		case "dso_preemptable":
			pre = enum.PreemptionDSOPreemptable
			p.next()
		case "hidden":
			vis = enum.VisibilityHidden
			p.next()
		case "protected":
			vis = enum.VisibilityProtected
			p.next()
		case "unnamed_addr":
			unnamed = enum.UnnamedAddrUnnamedAddr
			p.next()
		case "local_unnamed_addr":
			unnamed = enum.UnnamedAddrLocalUnnamedAddr
			p.next()
		case "thread_local":
			model, err := p.parseThreadLocal()
			if err != nil {
				return err
			}
			tls = model
		case "constant", "global":
			g := p.ensureGlobal(name)
			g.Linkage = link
			g.Preemption = pre
			g.Visibility = vis
			g.UnnamedAddr = unnamed
			g.TLSModel = tls
			g.Immutable = p.tok.s == "constant"
			p.next()
			return p.parseGlobalInit(g)
		case "alias":
			p.next()
			return p.parseAlias(name, link, pre, vis, unnamed, tls)
		case "ifunc":
			p.next()
			return p.parseIFunc(name, link, pre, vis, unnamed, tls)
		default:
			return p.errorf("unexpected global prefix %q", p.tok.s)
		}
	}
	if p.isIdent("constant") || p.isIdent("global") {
		g := p.ensureGlobal(name)
		g.Linkage = link
		g.Preemption = pre
		g.Visibility = vis
		g.UnnamedAddr = unnamed
		g.TLSModel = tls
		g.Immutable = p.tok.s == "constant"
		p.next()
		return p.parseGlobalInit(g)
	}
	return p.errorf("expected global, constant, or alias")
}

func (p *parser) ensureGlobal(name string) *ir.Global {
	if g, ok := p.globals[name]; ok {
		return g
	}
	g := ir.NewGlobal(name, types.I8)
	p.globals[name] = g
	p.m.Globals = append(p.m.Globals, g)
	return g
}

func (p *parser) parseGlobalInit(g *ir.Global) error {
	ct, err := p.parseType()
	if err != nil {
		return err
	}
	g.ContentType = ct
	g.Typ = types.NewPointer(ct)
	if p.startsValue() {
		init, err := p.parseValue(ct)
		if err != nil {
			return err
		}
		if c, ok := init.(constant.Constant); ok {
			g.Init = c
		} else {
			return p.errorf("global init is not a constant")
		}
	}
	for p.tok.kind == kComma {
		p.next()
		if err := p.skipGlobalSuffix(g); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseThreadLocal() (enum.TLSModel, error) {
	p.next() // thread_local
	if p.tok.kind != kLParen {
		return enum.TLSModelGeneric, nil
	}
	p.next()
	if p.tok.kind != kIdent {
		return enum.TLSModelNone, p.errorf("expected TLS model")
	}
	var m enum.TLSModel
	switch p.tok.s {
	case "localdynamic":
		m = enum.TLSModelLocalDynamic
	case "initialexec":
		m = enum.TLSModelInitialExec
	case "localexec":
		m = enum.TLSModelLocalExec
	default:
		return enum.TLSModelNone, p.errorf("unknown TLS model %q", p.tok.s)
	}
	p.next()
	if err := p.expect(kRParen); err != nil {
		return enum.TLSModelNone, err
	}
	return m, nil
}

func (p *parser) parseAlias(name string, link enum.Linkage, pre enum.Preemption, vis enum.Visibility, unnamed enum.UnnamedAddr, tls enum.TLSModel) error {
	if _, err := p.parseType(); err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	at, err := p.parseType()
	if err != nil {
		return err
	}
	if p.tok.kind == kGlobal {
		if _, ok := p.funcs[p.tok.s]; !ok {
			if _, ok := p.globals[p.tok.s]; !ok {
				f := ir.NewFunc(p.tok.s, types.Void)
				p.funcs[p.tok.s] = f
				p.m.Funcs = append(p.m.Funcs, f)
			}
		}
	}
	val, err := p.parseValue(at)
	if err != nil {
		return err
	}
	c, ok := val.(constant.Constant)
	if !ok {
		return p.errorf("aliasee is not a constant")
	}
	a := ir.NewAlias(name, c)
	a.Linkage = link
	a.Preemption = pre
	a.Visibility = vis
	a.UnnamedAddr = unnamed
	a.TLSModel = tls
	p.m.Aliases = append(p.m.Aliases, a)
	p.aliases[name] = a
	switch t := val.(type) {
	case *ir.Func:
		p.funcs[name] = t
	case *ir.Global:
		p.globals[name] = t
	}
	for p.tok.kind == kComma {
		p.next()
		if err := p.skipGlobalSuffix(nil); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseIFunc(name string, link enum.Linkage, pre enum.Preemption, vis enum.Visibility, unnamed enum.UnnamedAddr, tls enum.TLSModel) error {
	if _, err := p.parseType(); err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	at, err := p.parseType()
	if err != nil {
		return err
	}
	resolver, err := p.asConst(p.parseValue(at))
	if err != nil {
		return err
	}
	fn := ir.NewIFunc(name, resolver)
	fn.Linkage = link
	fn.Preemption = pre
	fn.Visibility = vis
	fn.UnnamedAddr = unnamed
	fn.TLSModel = tls
	p.m.IFuncs = append(p.m.IFuncs, fn)
	for p.tok.kind == kComma {
		p.next()
		if err := p.skipGlobalSuffix(nil); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) skipGlobalSuffix(g *ir.Global) error {
	if p.isIdent("align") {
		p.next()
		if p.tok.kind != kInt {
			return p.errorf("expected align integer")
		}
		if g != nil {
			g.Align = ir.Align(p.tok.i)
		}
		p.next()
		return nil
	}
	if p.tok.kind == kMetaName || p.tok.kind == kMetaID {
		return p.skipMDNode()
	}
	// section, comdat, etc.
	if p.tok.kind == kIdent {
		p.next()
		if p.tok.kind == kLParen {
			return p.skipBalanced(kLParen, kRParen)
		}
		if p.tok.kind == kString {
			p.next()
		}
		return nil
	}
	p.next()
	return nil
}

func (p *parser) parseFunc(def bool) error {
	p.next() // define/declare
	var (
		link     enum.Linkage
		pre      enum.Preemption
		unnamed  enum.UnnamedAddr
		retAttrs []ir.ReturnAttribute
	)
	for p.tok.kind == kIdent {
		switch p.tok.s {
		case "private":
			link = enum.LinkagePrivate
			p.next()
		case "internal":
			link = enum.LinkageInternal
			p.next()
		case "external":
			link = enum.LinkageExternal
			p.next()
		case "extern_weak":
			link = enum.LinkageExternWeak
			p.next()
		case "linkonce_odr":
			link = enum.LinkageLinkOnceODR
			p.next()
		case "weak_odr":
			link = enum.LinkageWeakODR
			p.next()
		case "weak":
			link = enum.LinkageWeak
			p.next()
		case "dso_local":
			pre = enum.PreemptionDSOLocal
			p.next()
		case "dso_preemptable":
			pre = enum.PreemptionDSOPreemptable
			p.next()
		case "hidden", "protected":
			p.next()
		case "fastcc", "ccc", "coldcc", "tailcc", "webkit_jscc", "anyregcc",
			"preserve_mostcc", "preserve_allcc", "swiftcc", "cxx_fast_tlscc":
			p.next()
		case "unnamed_addr", "local_unnamed_addr":
			// after params usually; stop header prefixes
			goto ret
		default:
			if p.skipRetAttr() {
				continue
			}
			goto ret
		}
	}
ret:
	ret, err := p.parseType()
	if err != nil {
		return err
	}
	if p.tok.kind != kGlobal {
		return p.errorf("expected function name")
	}
	name := p.tok.s
	p.next()
	if err := p.expect(kLParen); err != nil {
		return err
	}
	params, variadic, err := p.parseParams()
	if err != nil {
		return err
	}
	f, ok := p.funcs[name]
	if !ok {
		f = p.m.NewFunc(name, ret, params...)
		p.funcs[name] = f
	} else {
		ptypes := make([]types.Type, len(params))
		for i, param := range params {
			ptypes[i] = param.Type()
		}
		sig := types.NewFunc(ret, ptypes...)
		sig.Variadic = variadic
		f.Sig = sig
		f.Params = params
		f.Typ = types.NewPointer(sig)
	}
	f.Sig.Variadic = variadic
	f.Linkage = link
	f.Preemption = pre
	f.ReturnAttrs = retAttrs

	for !p.done() {
		if p.tok.kind == kIdent {
			switch p.tok.s {
			case "unnamed_addr":
				unnamed = enum.UnnamedAddrUnnamedAddr
				p.next()
				continue
			case "local_unnamed_addr":
				unnamed = enum.UnnamedAddrLocalUnnamedAddr
				p.next()
				continue
			case "comdat":
				p.next()
				if p.tok.kind == kLParen {
					if err := p.skipBalanced(kLParen, kRParen); err != nil {
						return err
					}
				}
				continue
			case "align":
				p.next()
				if p.tok.kind == kInt {
					p.next()
				}
				continue
			case "section", "gc", "prefix", "prologue":
				p.next()
				if p.tok.kind == kLParen {
					if err := p.skipBalanced(kLParen, kRParen); err != nil {
						return err
					}
				} else if p.tok.kind == kString || p.tok.kind == kInt || p.tok.kind == kIdent {
					p.next()
				}
				continue
			case "personality":
				p.next()
				if _, err := p.parseTypedValue(); err != nil {
					return err
				}
				continue
			}
			if isFuncAttr(p.tok.s) || p.tok.s == "memory" || p.tok.kind == kString {
				if err := p.skipFuncAttr(); err != nil {
					return err
				}
				continue
			}
		}
		if p.tok.kind == kAttrID {
			p.next()
			continue
		}
		if p.tok.kind == kString {
			if err := p.skipFuncAttr(); err != nil {
				return err
			}
			continue
		}
		if p.tok.kind == kMetaID || p.tok.kind == kMetaName || p.tok.kind == kBang {
			if err := p.skipMDNode(); err != nil {
				return err
			}
			continue
		}
		break
	}
	f.UnnamedAddr = unnamed

	if !def {
		return nil
	}
	if err := p.expect(kLBrace); err != nil {
		return err
	}
	if err := p.parseBody(f); err != nil {
		return err
	}
	return p.expect(kRBrace)
}

func (p *parser) parseParams() ([]*ir.Param, bool, error) {
	p.nextLocal = 0
	var params []*ir.Param
	variadic := false
	if p.tok.kind == kRParen {
		p.next()
		return params, false, nil
	}
	for {
		if p.tok.kind == kDots {
			p.next()
			variadic = true
			break
		}
		t, err := p.parseType()
		if err != nil {
			return nil, false, err
		}
		if err := p.skipParamAttrs(); err != nil {
			return nil, false, err
		}
		name := ""
		if p.tok.kind == kLocal {
			name = p.tok.s
			p.next()
		}
		param := newParam(name, t)
		if name == "" {
			param.SetID(p.nextLocal)
			p.nextLocal++
		} else if id, err := strconv.ParseInt(name, 10, 64); err == nil {
			if id >= p.nextLocal {
				p.nextLocal = id + 1
			}
		}
		params = append(params, param)
		if p.tok.kind != kComma {
			break
		}
		p.next()
	}
	if err := p.expect(kRParen); err != nil {
		return nil, false, err
	}
	return params, variadic, nil
}

func (p *parser) parseBody(f *ir.Func) error {
	p.fn = f
	p.locals = make(map[string]value.Value)
	p.blocks = make(map[string]*ir.Block)
	p.phis = nil
	for _, param := range f.Params {
		p.locals[localKey(param)] = param
	}
	// First pass: collect labels so br/phi can resolve blocks.
	if err := p.collectBlocks(f); err != nil {
		return err
	}
	first := true
	for !p.done() && p.tok.kind != kRBrace {
		name, ok := p.parseLabel()
		if !ok {
			if !first {
				return p.errorf("expected basic block label")
			}
			name = p.implicitEntryName()
		}
		first = false
		block := p.blocks[name]
		if block == nil {
			return p.errorf("unknown block %%%s", name)
		}
		for !p.done() && p.tok.kind != kRBrace && !p.atLabel() {
			if err := p.parseInst(block); err != nil {
				return err
			}
		}
	}
	if err := p.resolveFwds(); err != nil {
		return err
	}
	for _, phi := range p.phis {
		for _, inc := range phi.incs {
			bb, ok := p.blocks[inc.bb]
			if !ok {
				return p.errorf("unknown block %%%s", inc.bb)
			}
			var v value.Value
			if inc.c != nil {
				v = inc.c
			} else {
				v = p.locals[inc.val]
				if v == nil {
					return p.errorf("unknown value %%%s in phi", inc.val)
				}
			}
			phi.inst.Incs = append(phi.inst.Incs, ir.NewIncoming(v, bb))
		}
	}
	p.fn = nil
	return nil
}

func (p *parser) collectBlocks(f *ir.Func) error {
	saveI, saveTok := p.i, p.tok
	if !p.atLabel() {
		p.addImplicitEntry(f)
	}
	depth := 0
	for !p.done() {
		if depth == 0 && p.tok.kind == kRBrace {
			break
		}
		if depth == 0 && p.atLabel() {
			name := p.tok.s
			if p.tok.kind == kInt {
				name = strconv.FormatInt(p.tok.i, 10)
			}
			if _, exists := p.blocks[name]; !exists {
				p.blocks[name] = f.NewBlock(name)
			}
		}
		switch p.tok.kind {
		case kLParen, kLBrace, kLBrack, kLt:
			depth++
		case kRParen, kRBrace, kRBrack, kGt:
			depth--
		}
		p.next()
	}
	p.i, p.tok = saveI, saveTok
	return nil
}

// implicitEntryName is the LLVM unnamed-value id of a first block with no
// label. Params %0,%1 then the block is %2 (first inst often %3).
func (p *parser) implicitEntryName() string {
	return strconv.FormatInt(p.nextLocal, 10)
}

func (p *parser) addImplicitEntry(f *ir.Func) {
	name := p.implicitEntryName()
	if _, exists := p.blocks[name]; exists {
		return
	}
	b := f.NewBlock("")
	b.SetID(p.nextLocal)
	p.blocks[name] = b
}

func (p *parser) atLabel() bool {
	if p.tok.kind == kIdent || p.tok.kind == kInt {
		return p.peekKind() == kColon
	}
	return false
}

func (p *parser) parseLabel() (string, bool) {
	if p.tok.kind == kIdent && p.peekKind() == kColon {
		name := p.tok.s
		p.next()
		p.next()
		return name, true
	}
	if p.tok.kind == kInt && p.peekKind() == kColon {
		name := strconv.FormatInt(p.tok.i, 10)
		p.next()
		p.next()
		return name, true
	}
	return "", false
}

func (p *parser) parseInst(block *ir.Block) error {
	var ident ir.LocalIdent
	var name string
	if p.tok.kind == kLocal && p.peekKind() == kEq {
		name = p.tok.s
		ident = ir.NewLocalIdent(name)
		p.next()
		if err := p.expect(kEq); err != nil {
			return err
		}
	}
	if p.tok.kind != kIdent {
		return p.errorf("expected instruction, got %s", p.tok)
	}
	if p.tok.s == "tail" || p.tok.s == "musttail" || p.tok.s == "notail" {
		p.next()
	}
	if p.tok.kind != kIdent {
		return p.errorf("expected instruction, got %s", p.tok)
	}
	op := p.tok.s
	p.next()
	if strings.HasPrefix(op, "dbg_") {
		return p.skipDbgRecord()
	}
	switch op {
	case "alloca":
		return p.parseAlloca(block, ident, name)
	case "load":
		return p.parseLoad(block, ident, name)
	case "store":
		return p.parseStore(block)
	case "getelementptr":
		return p.parseGEP(block, ident, name)
	case "atomicrmw":
		return p.parseAtomicRMW(block, ident, name)
	case "cmpxchg":
		return p.parseCmpXchg(block, ident, name)
	case "fence":
		return p.parseFence(block)
	case "call":
		return p.parseCall(block, ident, name)
	case "invoke":
		return p.parseInvoke(block, ident, name)
	case "landingpad":
		return p.parseLandingPad(block, ident, name)
	case "resume":
		return p.parseResume(block)
	case "freeze":
		return p.parseFreeze(block, ident, name)
	case "ret":
		return p.parseRet(block)
	case "br":
		return p.parseBr(block)
	case "unreachable":
		block.Term = ir.NewUnreachable()
		return p.skipInstMD()
	case "icmp":
		return p.parseICmp(block, ident, name)
	case "add", "sub", "mul", "udiv", "sdiv", "urem", "srem", "and", "or", "xor",
		"shl", "lshr", "ashr",
		"fadd", "fsub", "fmul", "fdiv", "frem":
		return p.parseBin(block, ident, name, op)
	case "fcmp":
		return p.parseFCmp(block, ident, name)
	case "fneg":
		return p.parseFNeg(block, ident, name)
	case "zext", "sext", "trunc", "ptrtoint", "inttoptr", "bitcast", "addrspacecast",
		"fptoui", "fptosi", "uitofp", "sitofp", "fpext", "fptrunc":
		return p.parseCast(block, ident, name, op)
	case "phi":
		return p.parsePhi(block, ident, name)
	case "switch":
		return p.parseSwitch(block)
	case "select":
		return p.parseSelect(block, ident, name)
	case "extractvalue":
		return p.parseExtractValue(block, ident, name)
	case "insertvalue":
		return p.parseInsertValue(block, ident, name)
	case "extractelement":
		return p.parseExtractElement(block, ident, name)
	case "insertelement":
		return p.parseInsertElement(block, ident, name)
	case "shufflevector":
		return p.parseShuffleVector(block, ident, name)
	default:
		return p.errorf("unsupported instruction %q", op)
	}
}

// fwdRef is a use of a local before its defining instruction is parsed.
type fwdRef struct {
	name string
	typ  types.Type
}

func (f *fwdRef) String() string {
	if f.typ != nil {
		return fmt.Sprintf("%s %%%s", f.typ, f.name)
	}
	return "%" + f.name
}

func (f *fwdRef) Ident() string { return "%" + f.name }

func (f *fwdRef) Type() types.Type {
	if f.typ != nil {
		return f.typ
	}
	return types.I8
}

func (p *parser) resolveFwds() error {
	rewrite := func(op *value.Value) error {
		f, ok := (*op).(*fwdRef)
		if !ok {
			return nil
		}
		v, ok := p.locals[f.name]
		if !ok {
			return p.errorf("undefined local %%%s", f.name)
		}
		if _, still := v.(*fwdRef); still {
			return p.errorf("undefined local %%%s", f.name)
		}
		*op = v
		return nil
	}
	for _, b := range p.fn.Blocks {
		for _, inst := range b.Insts {
			for _, op := range inst.Operands() {
				if err := rewrite(op); err != nil {
					return err
				}
			}
		}
		if b.Term != nil {
			for _, op := range b.Term.Operands() {
				if err := rewrite(op); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *parser) bind(name string, v value.Value) {
	if name != "" {
		p.locals[name] = v
	}
}

func (p *parser) parseAlloca(block *ir.Block, ident ir.LocalIdent, name string) error {
	elem, err := p.parseType()
	if err != nil {
		return err
	}
	inst := ir.NewAlloca(elem)
	inst.LocalIdent = ident
	if p.tok.kind == kComma {
		p.next()
		if p.isIdent("align") {
			p.next()
			if p.tok.kind != kInt {
				return p.errorf("expected align")
			}
			inst.Align = ir.Align(p.tok.i)
			p.next()
		} else {
			// NElems
			tv, err := p.parseTypedValue()
			if err != nil {
				return err
			}
			inst.NElems = tv
			if p.tok.kind == kComma {
				p.next()
				if p.isIdent("align") {
					p.next()
					if p.tok.kind != kInt {
						return p.errorf("expected align")
					}
					inst.Align = ir.Align(p.tok.i)
					p.next()
				}
			}
		}
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) parseLoad(block *ir.Block, ident ir.LocalIdent, name string) error {
	atomic := false
	if p.isIdent("atomic") {
		atomic = true
		p.next()
	}
	if p.isIdent("volatile") {
		p.next()
	}
	elem, err := p.parseType()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	src, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	inst := ir.NewLoad(elem, src)
	inst.LocalIdent = ident
	inst.Atomic = atomic
	ord, err := p.parseSyncScopeAndOrdering()
	if err != nil {
		return err
	}
	inst.Ordering = ord
	if err := p.parseAlignAndMD(&inst.Align); err != nil {
		return err
	}
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) parseStore(block *ir.Block) error {
	atomic := false
	if p.isIdent("atomic") {
		atomic = true
		p.next()
	}
	if p.isIdent("volatile") {
		p.next()
	}
	src, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	dst, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	inst := ir.NewStore(src, dst)
	inst.Atomic = atomic
	ord, err := p.parseSyncScopeAndOrdering()
	if err != nil {
		return err
	}
	inst.Ordering = ord
	if err := p.parseAlignAndMD(&inst.Align); err != nil {
		return err
	}
	block.Insts = append(block.Insts, inst)
	return nil
}

// parseSyncScopeAndOrdering consumes optional syncscope("...") and an
// atomic ordering. LangRef: load/store atomic … [syncscope] <ordering>, align.
func (p *parser) parseSyncScopeAndOrdering() (enum.AtomicOrdering, error) {
	if p.isIdent("syncscope") {
		p.next()
		if p.tok.kind == kLParen {
			if err := p.skipBalanced(kLParen, kRParen); err != nil {
				return enum.AtomicOrderingNone, err
			}
		}
	}
	if p.tok.kind == kIdent {
		if o, ok := atomicOrdering(p.tok.s); ok {
			p.next()
			return o, nil
		}
	}
	return enum.AtomicOrderingNone, nil
}

func (p *parser) skipGEPFlags() (inbounds bool, err error) {
	for p.tok.kind == kIdent {
		switch p.tok.s {
		case "inbounds":
			inbounds = true
			p.next()
		case "nuw", "nsw":
			p.next()
		case "inrange":
			p.next()
			if p.tok.kind == kLParen {
				if err := p.skipBalanced(kLParen, kRParen); err != nil {
					return false, err
				}
			}
		default:
			return inbounds, nil
		}
	}
	return inbounds, nil
}

func (p *parser) parseGEP(block *ir.Block, ident ir.LocalIdent, name string) error {
	inbounds, err := p.skipGEPFlags()
	if err != nil {
		return err
	}
	elem, err := p.parseType()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	src, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	var idxs []value.Value
	for p.tok.kind == kComma {
		p.next()
		if p.tok.kind == kMetaID || p.tok.kind == kMetaName {
			p.skipMDNode()
			break
		}
		if p.isIdent("inrange") {
			p.next()
			if p.tok.kind == kLParen {
				if err := p.skipBalanced(kLParen, kRParen); err != nil {
					return err
				}
			}
		}
		idx, err := p.parseTypedValue()
		if err != nil {
			return err
		}
		idxs = append(idxs, idx)
	}
	inst := ir.NewGetElementPtr(elem, src, idxs...)
	inst.LocalIdent = ident
	inst.InBounds = inbounds
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) parseAtomicRMW(block *ir.Block, ident ir.LocalIdent, name string) error {
	if p.isIdent("volatile") {
		p.next()
	}
	if p.tok.kind != kIdent {
		return p.errorf("expected atomicrmw op")
	}
	op, err := atomicOp(p.tok.s)
	if err != nil {
		return p.errorf("%v", err)
	}
	p.next()
	dst, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	x, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	ord := enum.AtomicOrderingSeqCst
	if p.tok.kind == kIdent {
		if o, ok := atomicOrdering(p.tok.s); ok {
			ord = o
			p.next()
		}
	}
	inst := ir.NewAtomicRMW(op, dst, x, ord)
	inst.LocalIdent = ident
	inst.Typ = x.Type()
	if err := p.parseAlignAndMD(nil); err != nil {
		return err
	}
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) parseCmpXchg(block *ir.Block, ident ir.LocalIdent, name string) error {
	weak := false
	if p.isIdent("weak") {
		weak = true
		p.next()
	}
	if p.isIdent("volatile") {
		p.next()
	}
	ptr, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	cmp, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	neu, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	succ, err := p.parseSyncScopeAndOrdering()
	if err != nil {
		return err
	}
	fail := enum.AtomicOrderingMonotonic
	if p.tok.kind == kIdent {
		if o, ok := atomicOrdering(p.tok.s); ok {
			fail = o
			p.next()
		}
	}
	if succ == enum.AtomicOrderingNone {
		succ = enum.AtomicOrderingSeqCst
	}
	inst := ir.NewCmpXchg(ptr, cmp, neu, succ, fail)
	inst.LocalIdent = ident
	inst.Weak = weak
	if err := p.parseAlignAndMD(nil); err != nil {
		return err
	}
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

// parseFence is LangRef `fence [syncscope("…")] <ordering>`. rustc emits
// `fence acquire` before Arc::drop_slow. No SSA result.
func (p *parser) parseFence(block *ir.Block) error {
	ord, err := p.parseSyncScopeAndOrdering()
	if err != nil {
		return err
	}
	if ord == enum.AtomicOrderingNone {
		return p.errorf("expected fence ordering")
	}
	inst := ir.NewFence(ord)
	if err := p.skipInstMD(); err != nil {
		return err
	}
	block.Insts = append(block.Insts, inst)
	return nil
}

type callSite struct {
	callee value.Value
	args   []value.Value
	ret    types.Type
	fast   []enum.FastMathFlag
}

func (p *parser) parseCallSite() (callSite, error) {
	var cs callSite
	// LangRef: call FastMathFlag* CallingConvopt ... Typ
	// rustc emits `call nsz double @llvm.maximumnum.f64`.
	cs.fast = p.parseFastMath()
	p.skipCallingConv()
	for p.skipRetAttr() {
	}
	typ, err := p.parseType()
	if err != nil {
		return cs, err
	}
	var fnty *types.FuncType
	if ft, ok := typ.(*types.FuncType); ok {
		fnty = ft
		typ = ft.RetType
	}
	callee, err := p.parseCallee(typ)
	if err != nil {
		return cs, err
	}
	if err := p.expect(kLParen); err != nil {
		return cs, err
	}
	args, err := p.parseArgs()
	if err != nil {
		return cs, err
	}
	for p.tok.kind == kAttrID || (p.tok.kind == kIdent && isFuncAttr(p.tok.s)) {
		if p.tok.kind == kAttrID {
			p.next()
			continue
		}
		if err := p.skipFuncAttr(); err != nil {
			return cs, err
		}
	}
	cs.callee = callee
	cs.args = args
	cs.ret = typ
	if fnty != nil {
		cs.ret = fnty.RetType
	}
	return cs, nil
}

func (p *parser) parseCall(block *ir.Block, ident ir.LocalIdent, name string) error {
	cs, err := p.parseCallSite()
	if err != nil {
		return err
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	inst := &ir.InstCall{
		LocalIdent:    ident,
		Callee:        cs.callee,
		Args:          cs.args,
		Typ:           cs.ret,
		FastMathFlags: cs.fast,
	}
	block.Insts = append(block.Insts, inst)
	if !types.IsVoid(cs.ret) {
		p.bind(name, inst)
	}
	return nil
}

func (p *parser) parseInvoke(block *ir.Block, ident ir.LocalIdent, name string) error {
	cs, err := p.parseCallSite()
	if err != nil {
		return err
	}
	if err := p.wantIdent("to"); err != nil {
		return err
	}
	if err := p.wantIdent("label"); err != nil {
		return err
	}
	normal, err := p.blockRef()
	if err != nil {
		return err
	}
	if err := p.wantIdent("unwind"); err != nil {
		return err
	}
	if err := p.wantIdent("label"); err != nil {
		return err
	}
	unwind, err := p.blockRef()
	if err != nil {
		return err
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	term := &ir.TermInvoke{
		LocalIdent:         ident,
		Invokee:            cs.callee,
		Args:               cs.args,
		NormalRetTarget:    normal,
		ExceptionRetTarget: unwind,
		Typ:                cs.ret,
	}
	block.Term = term
	if !types.IsVoid(cs.ret) {
		p.bind(name, term)
	}
	return nil
}

func (p *parser) parseLandingPad(block *ir.Block, ident ir.LocalIdent, name string) error {
	typ, err := p.parseType()
	if err != nil {
		return err
	}
	cleanup := false
	var clauses []*ir.Clause
	for p.isIdent("cleanup") || p.isIdent("catch") || p.isIdent("filter") {
		which := p.tok.s
		p.next()
		switch which {
		case "cleanup":
			cleanup = true
		case "catch", "filter":
			x, err := p.parseTypedValue()
			if err != nil {
				return err
			}
			kind := enum.ClauseTypeCatch
			if which == "filter" {
				kind = enum.ClauseTypeFilter
			}
			clauses = append(clauses, ir.NewClause(kind, x))
		}
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	inst := ir.NewLandingPad(typ, clauses...)
	inst.Cleanup = cleanup
	inst.LocalIdent = ident
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) parseResume(block *ir.Block) error {
	x, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	block.Term = ir.NewResume(x)
	return nil
}

func (p *parser) parseFreeze(block *ir.Block, ident ir.LocalIdent, name string) error {
	x, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	inst := ir.NewInstFreeze(x)
	inst.LocalIdent = ident
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) blockRef() (*ir.Block, error) {
	if p.tok.kind != kLocal {
		return nil, p.errorf("expected block label")
	}
	name := p.tok.s
	p.next()
	b := p.blocks[name]
	if b == nil {
		return nil, p.errorf("unknown block %%%s", name)
	}
	return b, nil
}

func (p *parser) parseCallee(ret types.Type) (value.Value, error) {
	switch p.tok.kind {
	case kGlobal:
		return p.refAt(p.tok.s)
	case kLocal:
		name := p.tok.s
		p.next()
		return p.lookupLocal(name, p.ptr), nil
	case kIdent:
		if p.tok.s == "asm" {
			return p.parseInlineAsm(ret)
		}
	}
	return nil, p.errorf("expected callee, got %s", p.tok)
}

func (p *parser) parseInlineAsm(ret types.Type) (*ir.InlineAsm, error) {
	p.next() // asm
	a := ir.NewInlineAsm(types.NewPointer(types.NewFunc(ret)), "", "")
	for p.tok.kind == kIdent {
		switch p.tok.s {
		case "sideeffect":
			a.SideEffect = true
			p.next()
		case "alignstack":
			a.AlignStack = true
			p.next()
		case "inteldialect":
			a.IntelDialect = true
			p.next()
		case "unwind":
			p.next()
		default:
			return nil, p.errorf("unexpected asm flag %q", p.tok.s)
		}
	}
	if p.tok.kind != kString {
		return nil, p.errorf("expected asm string")
	}
	a.Asm = p.tok.s
	p.next()
	if err := p.expect(kComma); err != nil {
		return nil, err
	}
	if p.tok.kind != kString {
		return nil, p.errorf("expected asm constraint")
	}
	a.Constraint = p.tok.s
	p.next()
	return a, nil
}

func (p *parser) parseArgs() ([]value.Value, error) {
	var args []value.Value
	if p.tok.kind == kRParen {
		p.next()
		return args, nil
	}
	for {
		a, err := p.parseTypedValue()
		if err != nil {
			return nil, err
		}
		args = append(args, a)
		if p.tok.kind != kComma {
			break
		}
		p.next()
	}
	if err := p.expect(kRParen); err != nil {
		return nil, err
	}
	return args, nil
}

func (p *parser) parseRet(block *ir.Block) error {
	if p.isIdent("void") {
		p.next()
		block.Term = ir.NewRet(nil)
		if err := p.skipInstMD(); err != nil {
			return err
		}
		return nil
	}
	v, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	block.Term = ir.NewRet(v)
	if err := p.skipInstMD(); err != nil {
		return err
	}
	return nil
}

func (p *parser) parseSwitch(block *ir.Block) error {
	x, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	if err := p.wantIdent("label"); err != nil {
		return err
	}
	if p.tok.kind != kLocal {
		return p.errorf("expected default label")
	}
	def := p.blocks[p.tok.s]
	p.next()
	if def == nil {
		return p.errorf("unknown default block")
	}
	if err := p.expect(kLBrack); err != nil {
		return err
	}
	var cases []*ir.Case
	for p.tok.kind != kRBrack && !p.done() {
		cv, err := p.parseTypedValue()
		if err != nil {
			return err
		}
		c, ok := cv.(constant.Constant)
		if !ok {
			return p.errorf("switch case is not a constant")
		}
		if err := p.expect(kComma); err != nil {
			return err
		}
		if err := p.wantIdent("label"); err != nil {
			return err
		}
		if p.tok.kind != kLocal {
			return p.errorf("expected case label")
		}
		tgt := p.blocks[p.tok.s]
		p.next()
		if tgt == nil {
			return p.errorf("unknown switch target")
		}
		cases = append(cases, ir.NewCase(c, tgt))
	}
	if err := p.expect(kRBrack); err != nil {
		return err
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	block.Term = ir.NewSwitch(x, def, cases...)
	return nil
}

func (p *parser) parseBr(block *ir.Block) error {
	if p.isIdent("label") {
		p.next()
		if p.tok.kind != kLocal {
			return p.errorf("expected label")
		}
		name := p.tok.s
		p.next()
		tgt := p.blocks[name]
		if tgt == nil {
			return p.errorf("unknown block %%%s", name)
		}
		block.Term = ir.NewBr(tgt)
		if err := p.skipInstMD(); err != nil {
			return err
		}
		return nil
	}
	cond, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	if err := p.wantIdent("label"); err != nil {
		return err
	}
	if p.tok.kind != kLocal {
		return p.errorf("expected label")
	}
	tname := p.tok.s
	p.next()
	if err := p.expect(kComma); err != nil {
		return err
	}
	if err := p.wantIdent("label"); err != nil {
		return err
	}
	if p.tok.kind != kLocal {
		return p.errorf("expected label")
	}
	fname := p.tok.s
	p.next()
	tt, ff := p.blocks[tname], p.blocks[fname]
	if tt == nil || ff == nil {
		return p.errorf("unknown branch target")
	}
	block.Term = ir.NewCondBr(cond, tt, ff)
	if err := p.skipInstMD(); err != nil {
		return err
	}
	return nil
}

func (p *parser) parseICmp(block *ir.Block, ident ir.LocalIdent, name string) error {
	if p.isIdent("samesign") {
		p.next()
	}
	if p.tok.kind != kIdent {
		return p.errorf("expected icmp predicate")
	}
	pred, ok := ipred(p.tok.s)
	if !ok {
		return p.errorf("unknown icmp pred %q", p.tok.s)
	}
	p.next()
	x, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	y, err := p.parseValue(x.Type())
	if err != nil {
		return err
	}
	inst := ir.NewICmp(pred, x, y)
	inst.LocalIdent = ident
	if err := p.skipInstMD(); err != nil {
		return err
	}
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) parseFastMath() []enum.FastMathFlag {
	var flags []enum.FastMathFlag
	for p.tok.kind == kIdent {
		var f enum.FastMathFlag
		switch p.tok.s {
		case "nnan":
			f = enum.FastMathFlagNNaN
		case "ninf":
			f = enum.FastMathFlagNInf
		case "nsz":
			f = enum.FastMathFlagNSZ
		case "arcp":
			f = enum.FastMathFlagARcp
		case "contract":
			f = enum.FastMathFlagContract
		case "afn":
			f = enum.FastMathFlagAFn
		case "reassoc":
			f = enum.FastMathFlagReassoc
		case "fast":
			f = enum.FastMathFlagFast
		default:
			return flags
		}
		flags = append(flags, f)
		p.next()
	}
	return flags
}

func (p *parser) skipFastMath() {
	_ = p.parseFastMath()
}

func (p *parser) parseFNeg(block *ir.Block, ident ir.LocalIdent, name string) error {
	p.skipFastMath()
	x, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	inst := ir.NewFNeg(x)
	inst.LocalIdent = ident
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) parseFCmp(block *ir.Block, ident ir.LocalIdent, name string) error {
	p.skipFastMath()
	if p.tok.kind != kIdent {
		return p.errorf("expected fcmp predicate")
	}
	pred, ok := fpred(p.tok.s)
	if !ok {
		return p.errorf("unknown fcmp pred %q", p.tok.s)
	}
	p.next()
	x, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	y, err := p.parseValue(x.Type())
	if err != nil {
		return err
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	inst := ir.NewFCmp(pred, x, y)
	inst.LocalIdent = ident
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) parseBin(block *ir.Block, ident ir.LocalIdent, name, op string) error {
	p.skipFastMath()
	var flags []enum.OverflowFlag
	for p.isIdent("nsw") || p.isIdent("nuw") {
		if p.tok.s == "nsw" {
			flags = append(flags, enum.OverflowFlagNSW)
		} else {
			flags = append(flags, enum.OverflowFlagNUW)
		}
		p.next()
	}
	if p.isIdent("exact") || p.isIdent("nneg") || p.isIdent("disjoint") {
		p.next()
	}
	x, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	y, err := p.parseValue(x.Type())
	if err != nil {
		return err
	}
	var inst ir.Instruction
	switch op {
	case "add":
		n := ir.NewAdd(x, y)
		n.LocalIdent = ident
		n.OverflowFlags = flags
		inst = n
		p.bind(name, n)
	case "sub":
		n := ir.NewSub(x, y)
		n.LocalIdent = ident
		n.OverflowFlags = flags
		inst = n
		p.bind(name, n)
	case "mul":
		n := ir.NewMul(x, y)
		n.LocalIdent = ident
		n.OverflowFlags = flags
		inst = n
		p.bind(name, n)
	case "udiv":
		n := ir.NewUDiv(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "sdiv":
		n := ir.NewSDiv(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "urem":
		n := ir.NewURem(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "srem":
		n := ir.NewSRem(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "and":
		n := ir.NewAnd(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "or":
		n := ir.NewOr(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "xor":
		n := ir.NewXor(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "shl":
		n := ir.NewShl(x, y)
		n.LocalIdent = ident
		n.OverflowFlags = flags
		inst = n
		p.bind(name, n)
	case "lshr":
		n := ir.NewLShr(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "ashr":
		n := ir.NewAShr(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "fadd":
		n := ir.NewFAdd(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "fsub":
		n := ir.NewFSub(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "fmul":
		n := ir.NewFMul(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "fdiv":
		n := ir.NewFDiv(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "frem":
		n := ir.NewFRem(x, y)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	block.Insts = append(block.Insts, inst)
	return nil
}

func (p *parser) parseCast(block *ir.Block, ident ir.LocalIdent, name, op string) error {
	for p.isIdent("nsw") || p.isIdent("nuw") || p.isIdent("nneg") || p.isIdent("exact") {
		p.next()
	}
	from, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.wantIdent("to"); err != nil {
		return err
	}
	to, err := p.parseType()
	if err != nil {
		return err
	}
	var inst ir.Instruction
	switch op {
	case "zext":
		n := ir.NewZExt(from, to)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "sext":
		n := ir.NewSExt(from, to)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "trunc":
		n := ir.NewTrunc(from, to)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "ptrtoint":
		n := ir.NewPtrToInt(from, to)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "inttoptr":
		n := ir.NewIntToPtr(from, to)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "bitcast":
		n := ir.NewBitCast(from, to)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "addrspacecast":
		n := ir.NewAddrSpaceCast(from, to)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "fptoui":
		n := ir.NewFPToUI(from, to)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "fptosi":
		n := ir.NewFPToSI(from, to)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "uitofp":
		n := ir.NewUIToFP(from, to)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "sitofp":
		n := ir.NewSIToFP(from, to)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "fpext":
		n := ir.NewFPExt(from, to)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	case "fptrunc":
		n := ir.NewFPTrunc(from, to)
		n.LocalIdent = ident
		inst = n
		p.bind(name, n)
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	block.Insts = append(block.Insts, inst)
	return nil
}

func (p *parser) parsePhi(block *ir.Block, ident ir.LocalIdent, name string) error {
	typ, err := p.parseType()
	if err != nil {
		return err
	}
	inst := &ir.InstPhi{LocalIdent: ident, Typ: typ}
	var incs []pendingInc
	for {
		if err := p.expect(kLBrack); err != nil {
			return err
		}
		inc, err := p.parsePhiInc(typ)
		if err != nil {
			return err
		}
		incs = append(incs, inc)
		if err := p.expect(kRBrack); err != nil {
			return err
		}
		if p.tok.kind != kComma {
			break
		}
		p.next()
		if p.tok.kind == kMetaID || p.tok.kind == kMetaName {
			p.skipMDNode()
			break
		}
	}
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	p.phis = append(p.phis, pendingPhi{inst: inst, incs: incs})
	return nil
}

func (p *parser) parsePhiInc(typ types.Type) (pendingInc, error) {
	if p.tok.kind == kLocal {
		v := p.tok.s
		p.next()
		if err := p.expect(kComma); err != nil {
			return pendingInc{}, err
		}
		if p.tok.kind != kLocal {
			return pendingInc{}, p.errorf("expected block")
		}
		bb := p.tok.s
		p.next()
		return pendingInc{val: v, bb: bb}, nil
	}
	c, err := p.parseValue(typ)
	if err != nil {
		return pendingInc{}, err
	}
	if err := p.expect(kComma); err != nil {
		return pendingInc{}, err
	}
	if p.tok.kind != kLocal {
		return pendingInc{}, p.errorf("expected block")
	}
	bb := p.tok.s
	p.next()
	return pendingInc{c: c, bb: bb}, nil
}

func (p *parser) parseSelect(block *ir.Block, ident ir.LocalIdent, name string) error {
	cond, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	t, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	f, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	inst := ir.NewSelect(cond, t, f)
	inst.LocalIdent = ident
	if err := p.skipInstMD(); err != nil {
		return err
	}
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) parseExtractValue(block *ir.Block, ident ir.LocalIdent, name string) error {
	x, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	idxs, err := p.parseAggregateIndices()
	if err != nil {
		return err
	}
	inst := ir.NewExtractValue(x, idxs...)
	inst.LocalIdent = ident
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) parseInsertValue(block *ir.Block, ident ir.LocalIdent, name string) error {
	x, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	elem, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	idxs, err := p.parseAggregateIndices()
	if err != nil {
		return err
	}
	inst := &ir.InstInsertValue{LocalIdent: ident, X: x, Elem: elem, Indices: idxs, Typ: x.Type()}
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

// parseAggregateIndices reads ', 0, 1' then ', !dbg !N'. A comma that is
// not followed by an integer is metadata; skipInstMD consumes the rest.
func (p *parser) parseAggregateIndices() ([]uint64, error) {
	var idxs []uint64
	for p.tok.kind == kComma {
		p.next()
		if p.tok.kind != kInt {
			if err := p.skipMDNode(); err != nil {
				return nil, err
			}
			break
		}
		idxs = append(idxs, uint64(p.tok.i))
		p.next()
	}
	if err := p.skipInstMD(); err != nil {
		return nil, err
	}
	return idxs, nil
}

func (p *parser) parseExtractElement(block *ir.Block, ident ir.LocalIdent, name string) error {
	x, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	idx, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	inst := ir.NewExtractElement(x, idx)
	inst.LocalIdent = ident
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) parseInsertElement(block *ir.Block, ident ir.LocalIdent, name string) error {
	x, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	elem, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	idx, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	inst := ir.NewInsertElement(x, elem, idx)
	inst.LocalIdent = ident
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) parseShuffleVector(block *ir.Block, ident ir.LocalIdent, name string) error {
	x, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	y, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.expect(kComma); err != nil {
		return err
	}
	mask, err := p.parseTypedValue()
	if err != nil {
		return err
	}
	if err := p.skipInstMD(); err != nil {
		return err
	}
	inst := ir.NewShuffleVector(x, y, mask)
	inst.LocalIdent = ident
	block.Insts = append(block.Insts, inst)
	p.bind(name, inst)
	return nil
}

func (p *parser) parseAlignAndMD(align *ir.Align) error {
	for p.tok.kind == kComma {
		p.next()
		if p.isIdent("align") {
			p.next()
			if p.tok.kind == kInt {
				if align != nil {
					*align = ir.Align(p.tok.i)
				}
				p.next()
			}
			continue
		}
		if p.tok.kind == kMetaID || p.tok.kind == kMetaName {
			if err := p.skipMDNode(); err != nil {
				return err
			}
			continue
		}
		// unknown suffix
		if p.tok.kind == kIdent {
			p.next()
			if p.tok.kind == kInt {
				p.next()
			}
		}
	}
	return nil
}

func (p *parser) skipDbgRecord() error {
	if p.tok.kind != kLParen {
		return p.errorf("expected dbg record (")
	}
	return p.skipBalanced(kLParen, kRParen)
}

func (p *parser) skipInstMD() error {
	for p.tok.kind == kComma {
		p.next()
		if err := p.skipMDNode(); err != nil {
			return err
		}
	}
	return nil
}

// --- types / values ----------------------------------------------------------

func (p *parser) parseType() (types.Type, error) {
	t, err := p.parseTypePrimary()
	if err != nil {
		return nil, err
	}
	for p.tok.kind == kStar {
		p.next()
		t = types.NewPointer(t)
	}
	if p.tok.kind == kLParen {
		t, err = p.finishFuncType(t)
		if err != nil {
			return nil, err
		}
		for p.tok.kind == kStar {
			p.next()
			t = types.NewPointer(t)
		}
	}
	return t, nil
}

func (p *parser) finishFuncType(ret types.Type) (*types.FuncType, error) {
	p.next() // (
	var params []types.Type
	variadic := false
	if p.tok.kind != kRParen {
		for {
			if p.tok.kind == kDots {
				variadic = true
				p.next()
				break
			}
			pt, err := p.parseType()
			if err != nil {
				return nil, err
			}
			params = append(params, pt)
			if p.tok.kind != kComma {
				break
			}
			p.next()
		}
	}
	if err := p.expect(kRParen); err != nil {
		return nil, err
	}
	ft := types.NewFunc(ret, params...)
	ft.Variadic = variadic
	return ft, nil
}

func (p *parser) parseTypePrimary() (types.Type, error) {
	switch p.tok.kind {
	case kTypeInt:
		n := p.tok.i
		p.next()
		return intType(uint64(n)), nil
	case kIdent:
		switch p.tok.s {
		case "void":
			p.next()
			return types.Void, nil
		case "ptr":
			p.next()
			if p.isIdent("addrspace") {
				as, err := p.parseAddrSpace()
				if err != nil {
					return nil, err
				}
				return types.NewOpaquePointer(as), nil
			}
			return p.ptr, nil
		case "float":
			p.next()
			return types.Float, nil
		case "double":
			p.next()
			return types.Double, nil
		case "half":
			p.next()
			return types.Half, nil
		case "label":
			p.next()
			return types.Label, nil
		case "token":
			p.next()
			return types.Token, nil
		case "metadata":
			p.next()
			return types.Metadata, nil
		case "x86_mmx":
			p.next()
			return types.MMX, nil
		case "x86_fp80":
			p.next()
			return types.X86_FP80, nil
		case "fp128":
			p.next()
			return types.FP128, nil
		default:
			return nil, p.errorf("unknown type %q", p.tok.s)
		}
	case kLocal:
		name := p.tok.s
		p.next()
		if t, ok := p.types[name]; ok {
			return t, nil
		}
		return p.namedStruct(name), nil
	case kLBrace:
		return p.parseStructType(false)
	case kLt:
		// packed struct <{...}> or vector <N x ty>
		if p.peekKind() == kLBrace {
			p.next()
			return p.parseStructType(true)
		}
		return p.parseVectorType()
	case kLBrack:
		return p.parseArrayType()
	default:
		return nil, p.errorf("expected type, got %s", p.tok)
	}
}

func (p *parser) parseStructType(packed bool) (*types.StructType, error) {
	if err := p.expect(kLBrace); err != nil {
		return nil, err
	}
	var fields []types.Type
	if p.tok.kind != kRBrace {
		for {
			t, err := p.parseType()
			if err != nil {
				return nil, err
			}
			fields = append(fields, t)
			if p.tok.kind != kComma {
				break
			}
			p.next()
		}
	}
	if err := p.expect(kRBrace); err != nil {
		return nil, err
	}
	if packed {
		if err := p.expect(kGt); err != nil {
			return nil, err
		}
	}
	st := types.NewStruct(fields...)
	st.Packed = packed
	return st, nil
}

func (p *parser) parseArrayType() (*types.ArrayType, error) {
	if err := p.expect(kLBrack); err != nil {
		return nil, err
	}
	if p.tok.kind != kInt {
		return nil, p.errorf("expected array length")
	}
	n := uint64(p.tok.i)
	p.next()
	if err := p.wantIdent("x"); err != nil {
		return nil, err
	}
	elem, err := p.parseType()
	if err != nil {
		return nil, err
	}
	if err := p.expect(kRBrack); err != nil {
		return nil, err
	}
	return types.NewArray(n, elem), nil
}

func (p *parser) parseVectorType() (*types.VectorType, error) {
	if err := p.expect(kLt); err != nil {
		return nil, err
	}
	if p.isIdent("vscale") {
		p.next()
		if err := p.wantIdent("x"); err != nil {
			return nil, err
		}
	}
	if p.tok.kind != kInt {
		return nil, p.errorf("expected vector length")
	}
	n := uint64(p.tok.i)
	p.next()
	if err := p.wantIdent("x"); err != nil {
		return nil, err
	}
	elem, err := p.parseType()
	if err != nil {
		return nil, err
	}
	if err := p.expect(kGt); err != nil {
		return nil, err
	}
	return types.NewVector(n, elem), nil
}

func (p *parser) parseAddrSpace() (types.AddrSpace, error) {
	p.next() // addrspace
	if err := p.expect(kLParen); err != nil {
		return 0, err
	}
	if p.tok.kind != kInt {
		return 0, p.errorf("expected addrspace integer")
	}
	n := types.AddrSpace(p.tok.i)
	p.next()
	if err := p.expect(kRParen); err != nil {
		return 0, err
	}
	return n, nil
}

func (p *parser) parseTypedValue() (value.Value, error) {
	typ, err := p.parseType()
	if err != nil {
		return nil, err
	}
	if err := p.skipParamAttrs(); err != nil {
		return nil, err
	}
	return p.parseValue(typ)
}

func (p *parser) lookupLocal(name string, typ types.Type) value.Value {
	if v, ok := p.locals[name]; ok {
		return v
	}
	f := &fwdRef{name: name, typ: typ}
	p.locals[name] = f
	return f
}

func (p *parser) parseValue(typ types.Type) (value.Value, error) {
	switch p.tok.kind {
	case kLocal:
		name := p.tok.s
		p.next()
		return p.lookupLocal(name, typ), nil
	case kGlobal:
		return p.refAt(p.tok.s)
	case kMetaID, kMetaName, kBang:
		_ = p.skipMDNode()
		return constant.NewUndef(types.Metadata), nil
	default:
		if p.isConstExpr() {
			return p.parseConstExpr(typ)
		}
		return p.parseConst(typ)
	}
}

func (p *parser) refAt(name string) (value.Value, error) {
	p.next()
	if f, ok := p.funcs[name]; ok {
		return f, nil
	}
	if g, ok := p.globals[name]; ok {
		return g, nil
	}
	if a, ok := p.aliases[name]; ok {
		return a, nil
	}
	return p.ensureGlobal(name), nil
}

// --- attributes / metadata ---------------------------------------------------

func (p *parser) skipComdatDef() error {
	p.next()
	if err := p.expect(kEq); err != nil {
		return err
	}
	if err := p.wantIdent("comdat"); err != nil {
		return err
	}
	if p.tok.kind == kIdent {
		p.next()
	}
	return nil
}

func newParam(name string, typ types.Type) *ir.Param {
	if name == "" {
		return ir.NewParam("", typ)
	}
	if id, err := strconv.ParseInt(name, 10, 64); err == nil {
		p := &ir.Param{Typ: typ}
		p.SetID(id)
		return p
	}
	return ir.NewParam(name, typ)
}

func localKey(param *ir.Param) string {
	if param.IsUnnamed() {
		return strconv.FormatInt(param.LocalID, 10)
	}
	return param.LocalName
}

func (p *parser) parseAttrGroup() error {
	p.next() // attributes
	if p.tok.kind != kAttrID {
		return p.errorf("expected attribute group id")
	}
	id := p.tok.i
	p.next()
	if err := p.expect(kEq); err != nil {
		return err
	}
	if err := p.expect(kLBrace); err != nil {
		return err
	}
	for p.tok.kind != kRBrace && !p.done() {
		if err := p.skipFuncAttr(); err != nil {
			return err
		}
	}
	if err := p.expect(kRBrace); err != nil {
		return err
	}
	p.m.AttrGroupDefs = append(p.m.AttrGroupDefs, &ir.AttrGroupDef{ID: id})
	return nil
}

func (p *parser) skipFuncAttr() error {
	if p.tok.kind == kString {
		p.next()
		if p.tok.kind == kEq {
			p.next()
			if p.tok.kind == kString {
				p.next()
			}
		}
		return nil
	}
	if p.tok.kind != kIdent {
		p.next()
		return nil
	}
	p.next()
	if p.tok.kind == kLParen {
		return p.skipBalanced(kLParen, kRParen)
	}
	return nil
}

func (p *parser) skipParamAttrs() error {
	for {
		if p.tok.kind == kIdent && isParamAttr(p.tok.s) {
			p.next()
			if p.tok.kind == kLParen {
				if err := p.skipBalanced(kLParen, kRParen); err != nil {
					return err
				}
			}
			continue
		}
		if p.isIdent("align") {
			p.next()
			if p.tok.kind == kInt {
				p.next()
			} else if p.tok.kind == kLParen {
				if err := p.skipBalanced(kLParen, kRParen); err != nil {
					return err
				}
			}
			continue
		}
		if p.isIdent("dereferenceable") || p.isIdent("dereferenceable_or_null") ||
			p.isIdent("alignstack") || p.isIdent("byval") || p.isIdent("preallocated") ||
			p.isIdent("sret") || p.isIdent("elementtype") || p.isIdent("captures") ||
			p.isIdent("range") || p.isIdent("initializes") || p.isIdent("allockind") ||
			p.isIdent("allocsize") {
			p.next()
			if p.tok.kind == kLParen {
				if err := p.skipBalanced(kLParen, kRParen); err != nil {
					return err
				}
			}
			continue
		}
		if p.isIdent("noalias") || p.isIdent("readonly") || p.isIdent("writeonly") ||
			p.isIdent("readnone") || p.isIdent("writable") || p.isIdent("dead_on_unwind") ||
			p.isIdent("dead_on_return") ||
			p.isIdent("allocptr") || p.isIdent("allocalign") || p.isIdent("nocapture") ||
			p.isIdent("nofree") || p.isIdent("nonnull") {
			p.next()
			continue
		}
		break
	}
	return nil
}

func (p *parser) skipCallingConv() {
	if p.tok.kind != kIdent {
		return
	}
	switch p.tok.s {
	case "fastcc", "ccc", "coldcc", "tailcc", "webkit_jscc", "anyregcc",
		"preserve_mostcc", "preserve_allcc", "swiftcc", "cxx_fast_tlscc":
		p.next()
	}
}

func (p *parser) skipRetAttr() bool {
	if p.tok.kind == kIdent && isRetAttr(p.tok.s) {
		_, _ = p.parseRetAttr()
		return true
	}
	if p.isIdent("align") {
		p.next()
		if p.tok.kind == kInt {
			p.next()
		}
		return true
	}
	if p.isIdent("dereferenceable") || p.isIdent("dereferenceable_or_null") || p.isIdent("range") {
		p.next()
		if p.tok.kind == kLParen {
			_ = p.skipBalanced(kLParen, kRParen)
		}
		return true
	}
	if p.isIdent("noalias") {
		p.next()
		return true
	}
	return false
}

func (p *parser) parseRetAttr() (ir.ReturnAttribute, error) {
	s := p.tok.s
	p.next()
	switch s {
	case "zeroext":
		return enum.ReturnAttrZeroExt, nil
	case "signext":
		return enum.ReturnAttrSignExt, nil
	case "noundef":
		return enum.ReturnAttrNoUndef, nil
	case "noalias":
		return enum.ReturnAttrNoAlias, nil
	case "nonnull":
		return enum.ReturnAttrNonNull, nil
	case "inreg":
		return enum.ReturnAttrInReg, nil
	default:
		if p.tok.kind == kLParen {
			if err := p.skipBalanced(kLParen, kRParen); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
}

func (p *parser) skipMetadataDef() error {
	p.next() // !id or !name
	if p.tok.kind == kEq {
		p.next()
		if p.isIdent("distinct") {
			p.next()
		}
		return p.skipMDNode()
	}
	return nil
}

func (p *parser) skipMDNode() error {
	switch p.tok.kind {
	case kMetaID, kMetaName:
		p.next()
		if p.tok.kind == kLParen {
			return p.skipBalanced(kLParen, kRParen)
		}
		// attachment: !llvm.loop !6
		if p.tok.kind == kMetaID || p.tok.kind == kMetaName || p.tok.kind == kBang {
			return p.skipMDNode()
		}
		return nil
	case kBang:
		p.next()
		if p.tok.kind == kLBrace {
			return p.skipBalanced(kLBrace, kRBrace)
		}
		if p.tok.kind == kLParen {
			return p.skipBalanced(kLParen, kRParen)
		}
		return nil
	case kLBrace:
		return p.skipBalanced(kLBrace, kRBrace)
	case kIdent:
		p.next()
		if p.tok.kind == kLParen {
			return p.skipBalanced(kLParen, kRParen)
		}
		return nil
	default:
		p.next()
		return nil
	}
}

func (p *parser) skipBalanced(open, close kind) error {
	if p.tok.kind != open {
		return p.errorf("expected %s", token{kind: open})
	}
	depth := 0
	for !p.done() {
		switch p.tok.kind {
		case open:
			depth++
		case close:
			depth--
			p.next()
			if depth == 0 {
				return nil
			}
			continue
		}
		p.next()
	}
	return p.errorf("unbalanced %s", token{kind: open})
}

// --- token helpers -----------------------------------------------------------

func (p *parser) done() bool {
	return p.i >= len(p.toks) || p.toks[p.i].kind == kEOF
}

func (p *parser) tokAt() token {
	if p.i >= len(p.toks) {
		return token{kind: kEOF}
	}
	return p.toks[p.i]
}

func (p *parser) next() token {
	t := p.tokAt()
	if p.i < len(p.toks) {
		p.i++
	}
	p.tok = p.tokAt()
	return t
}

func (p *parser) expect(k kind) error {
	if p.tok.kind != k {
		return p.errorf("expected %s, got %s", token{kind: k}, p.tok)
	}
	p.next()
	return nil
}

func (p *parser) wantIdent(s string) error {
	if !p.isIdent(s) {
		return p.errorf("expected %q, got %s", s, p.tok)
	}
	p.next()
	return nil
}

func (p *parser) isIdent(s string) bool {
	return p.tok.kind == kIdent && p.tok.s == s
}

func (p *parser) ident() (string, error) {
	if p.tok.kind != kIdent {
		return "", p.errorf("expected identifier, got %s", p.tok)
	}
	s := p.tok.s
	p.next()
	return s, nil
}

func (p *parser) stringLit() (string, error) {
	if p.tok.kind != kString {
		return "", p.errorf("expected string, got %s", p.tok)
	}
	s := p.tok.s
	p.next()
	return s, nil
}

func (p *parser) peekKind() kind {
	j := p.i + 1
	if j >= len(p.toks) {
		return kEOF
	}
	return p.toks[j].kind
}

func (p *parser) errorf(format string, args ...any) error {
	t := p.tok
	msg := fmt.Sprintf(format, args...)
	if p.path != "" {
		return fmt.Errorf("%s:%d:%d: %s", p.path, t.line, t.col, msg)
	}
	return fmt.Errorf("%d:%d: %s", t.line, t.col, msg)
}

func intType(n uint64) *types.IntType {
	switch n {
	case 1:
		return types.I1
	case 8:
		return types.I8
	case 16:
		return types.I16
	case 32:
		return types.I32
	case 64:
		return types.I64
	case 128:
		return types.I128
	case 256:
		return types.I256
	default:
		return types.NewInt(n)
	}
}

func isRetAttr(s string) bool {
	switch s {
	case "zeroext", "signext", "inreg", "noalias", "noundef", "nonnull",
		"dereferenceable", "dereferenceable_or_null":
		return true
	}
	return false
}

func isParamAttr(s string) bool {
	switch s {
	case "zeroext", "signext", "inreg", "byval", "sret", "noalias", "nocapture",
		"nofree", "noundef", "nonnull", "readonly", "readnone", "writeonly",
		"immarg", "returned", "swiftself", "swifterror", "nest", "nomerge",
		"inalloca", "preallocated", "byref", "elementtype", "no_cfi",
		"dead_on_return", "dead_on_unwind":
		return true
	}
	return false
}

func isFuncAttr(s string) bool {
	switch s {
	case "alwaysinline", "argmemonly", "builtin", "cold", "convergent",
		"inaccessiblememonly", "inaccessiblemem_or_argmemonly", "inlinehint",
		"jumptable", "minsize", "naked", "nobuiltin", "nocf_check", "noduplicate",
		"nofree", "noimplicitfloat", "noinline", "nomerge", "nonlazybind",
		"norecurse", "noredzone", "noreturn", "nosync", "nounwind",
		"null_pointer_is_valid", "optforfuzzing", "optnone", "optsize",
		"readnone", "readonly", "returns_twice", "safestack", "sanitize_address",
		"sanitize_hwaddress", "sanitize_memory", "sanitize_memtag",
		"sanitize_thread", "shadowcallstack", "speculatable",
		"speculative_load_hardening", "ssp", "sspreq", "sspstrong", "strictfp",
		"uwtable", "willreturn", "writeonly", "nocallback", "mustprogress",
		"nocreateundeforpoison", "memory", "allockind", "allocsize", "alignstack":
		return true
	}
	return false
}

func fpred(s string) (enum.FPred, bool) {
	switch s {
	case "false":
		return enum.FPredFalse, true
	case "oeq":
		return enum.FPredOEQ, true
	case "ogt":
		return enum.FPredOGT, true
	case "oge":
		return enum.FPredOGE, true
	case "olt":
		return enum.FPredOLT, true
	case "ole":
		return enum.FPredOLE, true
	case "one":
		return enum.FPredONE, true
	case "ord":
		return enum.FPredORD, true
	case "ueq":
		return enum.FPredUEQ, true
	case "ugt":
		return enum.FPredUGT, true
	case "uge":
		return enum.FPredUGE, true
	case "ult":
		return enum.FPredULT, true
	case "ule":
		return enum.FPredULE, true
	case "une":
		return enum.FPredUNE, true
	case "uno":
		return enum.FPredUNO, true
	case "true":
		return enum.FPredTrue, true
	}
	return 0, false
}

func ipred(s string) (enum.IPred, bool) {
	switch s {
	case "eq":
		return enum.IPredEQ, true
	case "ne":
		return enum.IPredNE, true
	case "ugt":
		return enum.IPredUGT, true
	case "uge":
		return enum.IPredUGE, true
	case "ult":
		return enum.IPredULT, true
	case "ule":
		return enum.IPredULE, true
	case "sgt":
		return enum.IPredSGT, true
	case "sge":
		return enum.IPredSGE, true
	case "slt":
		return enum.IPredSLT, true
	case "sle":
		return enum.IPredSLE, true
	}
	return 0, false
}

func atomicOp(s string) (enum.AtomicOp, error) {
	switch s {
	case "xchg":
		return enum.AtomicOpXChg, nil
	case "add":
		return enum.AtomicOpAdd, nil
	case "sub":
		return enum.AtomicOpSub, nil
	case "and":
		return enum.AtomicOpAnd, nil
	case "nand":
		return enum.AtomicOpNAnd, nil
	case "or":
		return enum.AtomicOpOr, nil
	case "xor":
		return enum.AtomicOpXor, nil
	case "max":
		return enum.AtomicOpMax, nil
	case "min":
		return enum.AtomicOpMin, nil
	case "umax":
		return enum.AtomicOpUMax, nil
	case "umin":
		return enum.AtomicOpUMin, nil
	case "fadd":
		return enum.AtomicOpFAdd, nil
	case "fsub":
		return enum.AtomicOpFSub, nil
	default:
		return 0, fmt.Errorf("unknown atomicrmw op %q", s)
	}
}

func atomicOrdering(s string) (enum.AtomicOrdering, bool) {
	switch s {
	case "unordered":
		return enum.AtomicOrderingUnordered, true
	case "monotonic":
		return enum.AtomicOrderingMonotonic, true
	case "acquire":
		return enum.AtomicOrderingAcquire, true
	case "release":
		return enum.AtomicOrderingRelease, true
	case "acq_rel":
		return enum.AtomicOrderingAcqRel, true
	case "seq_cst":
		return enum.AtomicOrderingSeqCst, true
	}
	return 0, false
}
