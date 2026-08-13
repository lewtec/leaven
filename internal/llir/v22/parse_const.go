package v22

import (
	"strings"

	"github.com/lewtec/leaven/internal/llir/ir/constant"
	"github.com/lewtec/leaven/internal/llir/ir/types"
	"github.com/lewtec/leaven/internal/llir/ir/value"
)

func (p *parser) startsConst() bool { return p.startsValue() }

func (p *parser) startsValue() bool {
	switch p.tok.kind {
	case kInt, kFloat, kCString, kString, kLBrace, kLBrack, kLt, kLocal:
		return true
	case kGlobal:
		// `@next = ...` is the next top-level definition, not an initializer.
		return p.peekKind() != kEq
	case kIdent:
		switch p.tok.s {
		case "true", "false", "null", "undef", "poison", "zeroinitializer", "splat":
			return true
		}
		return p.isConstExpr()
	}
	return false
}

func (p *parser) isConstExpr() bool {
	if p.tok.kind != kIdent {
		return false
	}
	switch p.tok.s {
	case "getelementptr",
		"trunc", "zext", "sext", "fptrunc", "fpext", "fptoui", "fptosi",
		"uitofp", "sitofp", "inttoptr", "ptrtoint", "bitcast", "addrspacecast",
		"add", "sub", "mul", "udiv", "sdiv", "urem", "srem",
		"shl", "lshr", "ashr", "and", "or", "xor",
		"fadd", "fsub", "fmul", "fdiv", "frem":
		return true
	}
	return false
}

func (p *parser) asConst(v value.Value, err error) (constant.Constant, error) {
	if err != nil {
		return nil, err
	}
	c, ok := v.(constant.Constant)
	if !ok {
		return nil, p.errorf("expected constant, got %T", v)
	}
	return c, nil
}

func (p *parser) parseConstExpr(typ types.Type) (constant.Constant, error) {
	op := p.tok.s
	p.next()
	var inbounds bool
	if op == "getelementptr" {
		var err error
		inbounds, err = p.skipGEPFlags()
		if err != nil {
			return nil, err
		}
	} else {
		for p.isIdent("nsw") || p.isIdent("nuw") || p.isIdent("exact") || p.isIdent("nneg") {
			p.next()
		}
	}
	if err := p.expect(kLParen); err != nil {
		return nil, err
	}
	switch op {
	case "getelementptr":
		return p.parseGEPExpr(typ, inbounds)
	case "trunc", "zext", "sext", "fptrunc", "fpext", "fptoui", "fptosi",
		"uitofp", "sitofp", "inttoptr", "ptrtoint", "bitcast", "addrspacecast":
		return p.parseCastExpr(op)
	case "add", "sub", "mul", "udiv", "sdiv", "urem", "srem",
		"shl", "lshr", "ashr", "and", "or", "xor",
		"fadd", "fsub", "fmul", "fdiv", "frem":
		return p.parseBinExpr(op)
	default:
		return nil, p.errorf("unsupported constexpr %q", op)
	}
}

func (p *parser) parseBinExpr(op string) (constant.Constant, error) {
	x, err := p.asConst(p.parseTypedValue())
	if err != nil {
		return nil, err
	}
	if err := p.expect(kComma); err != nil {
		return nil, err
	}
	y, err := p.asConst(p.parseTypedValue())
	if err != nil {
		return nil, err
	}
	if err := p.expect(kRParen); err != nil {
		return nil, err
	}
	switch op {
	case "add":
		return constant.NewAdd(x, y), nil
	case "sub":
		return constant.NewSub(x, y), nil
	case "mul":
		return constant.NewMul(x, y), nil
	case "udiv":
		return constant.NewUDiv(x, y), nil
	case "sdiv":
		return constant.NewSDiv(x, y), nil
	case "urem":
		return constant.NewURem(x, y), nil
	case "srem":
		return constant.NewSRem(x, y), nil
	case "shl":
		return constant.NewShl(x, y), nil
	case "lshr":
		return constant.NewLShr(x, y), nil
	case "ashr":
		return constant.NewAShr(x, y), nil
	case "and":
		return constant.NewAnd(x, y), nil
	case "or":
		return constant.NewOr(x, y), nil
	case "xor":
		return constant.NewXor(x, y), nil
	case "fadd":
		return constant.NewFAdd(x, y), nil
	case "fsub":
		return constant.NewFSub(x, y), nil
	case "fmul":
		return constant.NewFMul(x, y), nil
	case "fdiv":
		return constant.NewFDiv(x, y), nil
	case "frem":
		return constant.NewFRem(x, y), nil
	default:
		return nil, p.errorf("unsupported binary constexpr %q", op)
	}
}

func (p *parser) parseGEPExpr(_ types.Type, inbounds bool) (constant.Constant, error) {
	elem, err := p.parseType()
	if err != nil {
		return nil, err
	}
	if err := p.expect(kComma); err != nil {
		return nil, err
	}
	src, err := p.asConst(p.parseTypedValue())
	if err != nil {
		return nil, err
	}
	var idxs []constant.Constant
	for p.tok.kind == kComma {
		p.next()
		inrange := false
		if p.isIdent("inrange") {
			inrange = true
			p.next()
			if p.tok.kind == kLParen {
				if err := p.skipBalanced(kLParen, kRParen); err != nil {
					return nil, err
				}
			}
		}
		idx, err := p.asConst(p.parseTypedValue())
		if err != nil {
			return nil, err
		}
		if inrange {
			i := constant.NewIndex(idx)
			i.InRange = true
			idx = i
		}
		idxs = append(idxs, idx)
	}
	if err := p.expect(kRParen); err != nil {
		return nil, err
	}
	expr := constant.NewGetElementPtr(elem, src, idxs...)
	expr.InBounds = inbounds
	return expr, nil
}

func (p *parser) parseCastExpr(op string) (constant.Constant, error) {
	from, err := p.asConst(p.parseTypedValue())
	if err != nil {
		return nil, err
	}
	if err := p.wantIdent("to"); err != nil {
		return nil, err
	}
	to, err := p.parseType()
	if err != nil {
		return nil, err
	}
	if err := p.expect(kRParen); err != nil {
		return nil, err
	}
	switch op {
	case "trunc":
		return constant.NewTrunc(from, to), nil
	case "zext":
		return constant.NewZExt(from, to), nil
	case "sext":
		return constant.NewSExt(from, to), nil
	case "fptrunc":
		return constant.NewFPTrunc(from, to), nil
	case "fpext":
		return constant.NewFPExt(from, to), nil
	case "fptoui":
		return constant.NewFPToUI(from, to), nil
	case "fptosi":
		return constant.NewFPToSI(from, to), nil
	case "uitofp":
		return constant.NewUIToFP(from, to), nil
	case "sitofp":
		return constant.NewSIToFP(from, to), nil
	case "inttoptr":
		return constant.NewIntToPtr(from, to), nil
	case "ptrtoint":
		return constant.NewPtrToInt(from, to), nil
	case "addrspacecast":
		return constant.NewAddrSpaceCast(from, to), nil
	case "bitcast":
		return constant.NewBitCast(from, to), nil
	default:
		return nil, p.errorf("unsupported cast constexpr %q", op)
	}
}

func (p *parser) parseConst(typ types.Type) (constant.Constant, error) {
	if p.isConstExpr() {
		return p.parseConstExpr(typ)
	}
	switch p.tok.kind {
	case kInt:
		if ft, ok := typ.(*types.FloatType); ok && strings.HasPrefix(strings.ToLower(p.tok.s), "0x") {
			c, err := constant.NewFloatFromString(ft, p.tok.s)
			if err != nil {
				return nil, p.errorf("hex float %q: %v", p.tok.s, err)
			}
			p.next()
			return c, nil
		}
		it, ok := typ.(*types.IntType)
		if !ok {
			return nil, p.errorf("integer constant for non-int type %s", typ)
		}
		c, err := parseIntConst(it, p.tok.s, p.tok.i)
		if err != nil {
			return nil, p.errorf("integer constant %q: %v", p.tok.s, err)
		}
		p.next()
		return c, nil
	case kFloat:
		ft, ok := typ.(*types.FloatType)
		if !ok {
			return nil, p.errorf("float constant for non-float type %s", typ)
		}
		c, err := constant.NewFloatFromString(ft, p.tok.s)
		if err != nil {
			return nil, p.errorf("float %q: %v", p.tok.s, err)
		}
		p.next()
		return c, nil
	case kIdent:
		switch p.tok.s {
		case "true":
			p.next()
			return constant.True, nil
		case "false":
			p.next()
			return constant.False, nil
		case "null":
			p.next()
			pt, ok := typ.(*types.PointerType)
			if !ok {
				pt = p.ptr
			}
			return constant.NewNull(pt), nil
		case "undef":
			p.next()
			return constant.NewUndef(typ), nil
		case "poison":
			p.next()
			return constant.NewPoison(typ), nil
		case "zeroinitializer":
			p.next()
			return constant.NewZeroInitializer(typ), nil
		case "splat":
			return p.parseSplat(typ)
		default:
			return nil, p.errorf("unexpected constant %q", p.tok.s)
		}
	case kCString:
		b := p.tok.b
		p.next()
		return constant.NewCharArray(b), nil
	case kLBrace:
		return p.parseStructConst(typ, false)
	case kLt:
		if p.peekKind() == kLBrace {
			p.next()
			return p.parseStructConst(typ, true)
		}
		return p.parseVectorConst(typ)
	case kLBrack:
		return p.parseArrayConst(typ)
	case kGlobal:
		return p.asConst(p.refAt(p.tok.s))
	default:
		return nil, p.errorf("expected constant, got %s", p.tok)
	}
}

// parseSplat is LLVM splat (TY C): a vector of typ.Len copies of C.
func (p *parser) parseSplat(typ types.Type) (constant.Constant, error) {
	p.next() // splat
	if err := p.expect(kLParen); err != nil {
		return nil, err
	}
	et, err := p.parseType()
	if err != nil {
		return nil, err
	}
	elem, err := p.parseConst(et)
	if err != nil {
		return nil, err
	}
	if err := p.expect(kRParen); err != nil {
		return nil, err
	}
	vt, ok := typ.(*types.VectorType)
	if !ok {
		return nil, p.errorf("splat for non-vector type %s", typ)
	}
	elems := make([]constant.Constant, vt.Len)
	for i := range elems {
		elems[i] = elem
	}
	return constant.NewVector(vt, elems...), nil
}

func (p *parser) parseStructConst(typ types.Type, packed bool) (constant.Constant, error) {
	if err := p.expect(kLBrace); err != nil {
		return nil, err
	}
	var fields []constant.Constant
	if p.tok.kind != kRBrace {
		for {
			ft, err := p.parseType()
			if err != nil {
				return nil, err
			}
			if err := p.skipParamAttrs(); err != nil {
				return nil, err
			}
			f, err := p.parseConst(ft)
			if err != nil {
				return nil, err
			}
			fields = append(fields, f)
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
	st, _ := typ.(*types.StructType)
	if st == nil {
		st = types.NewStruct()
		st.Packed = packed
		for _, f := range fields {
			st.Fields = append(st.Fields, f.Type())
		}
	}
	return constant.NewStruct(st, fields...), nil
}

func (p *parser) parseArrayConst(typ types.Type) (constant.Constant, error) {
	if err := p.expect(kLBrack); err != nil {
		return nil, err
	}
	var elems []constant.Constant
	if p.tok.kind != kRBrack {
		for {
			et, err := p.parseType()
			if err != nil {
				return nil, err
			}
			e, err := p.parseConst(et)
			if err != nil {
				return nil, err
			}
			elems = append(elems, e)
			if p.tok.kind != kComma {
				break
			}
			p.next()
		}
	}
	if err := p.expect(kRBrack); err != nil {
		return nil, err
	}
	at, _ := typ.(*types.ArrayType)
	return constant.NewArray(at, elems...), nil
}

func (p *parser) parseVectorConst(typ types.Type) (constant.Constant, error) {
	if err := p.expect(kLt); err != nil {
		return nil, err
	}
	var elems []constant.Constant
	if p.tok.kind != kGt {
		for {
			et, err := p.parseType()
			if err != nil {
				return nil, err
			}
			e, err := p.parseConst(et)
			if err != nil {
				return nil, err
			}
			elems = append(elems, e)
			if p.tok.kind != kComma {
				break
			}
			p.next()
		}
	}
	if err := p.expect(kGt); err != nil {
		return nil, err
	}
	vt, _ := typ.(*types.VectorType)
	return constant.NewVector(vt, elems...), nil
}

// parseIntConst keeps the full lexeme so i128 (and other wide) literals
// are not truncated to the lexer's int64 field.
func parseIntConst(typ *types.IntType, s string, fallback int64) (*constant.Int, error) {
	if s == "" {
		return constant.NewInt(typ, fallback), nil
	}
	if len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		s = "u0x" + s[2:]
	}
	c, err := constant.NewIntFromString(typ, s)
	if err != nil {
		return nil, err
	}
	return c, nil
}
