package v22

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type kind uint8

const (
	kEOF kind = iota
	kIdent
	kLocal
	kGlobal
	kAttrID
	kMetaID
	kMetaName
	kInt
	kFloat
	kString
	kCString
	kTypeInt
	kEq
	kComma
	kColon
	kStar
	kLParen
	kRParen
	kLBrace
	kRBrace
	kLBrack
	kRBrack
	kLt
	kGt
	kDots
	kBang
)

type token struct {
	kind kind
	s    string
	i    int64
	b    []byte
	line int
	col  int
}

func (t token) String() string {
	switch t.kind {
	case kEOF:
		return "EOF"
	case kIdent:
		return t.s
	case kLocal:
		return "%" + t.s
	case kGlobal:
		return "@" + t.s
	case kAttrID:
		return fmt.Sprintf("#%d", t.i)
	case kMetaID:
		return fmt.Sprintf("!%d", t.i)
	case kMetaName:
		return "!" + t.s
	case kInt:
		return fmt.Sprintf("%d", t.i)
	case kFloat:
		return t.s
	case kString:
		return fmt.Sprintf("%q", t.s)
	case kCString:
		return fmt.Sprintf("c%q", t.b)
	case kTypeInt:
		return fmt.Sprintf("i%d", t.i)
	case kEq:
		return "="
	case kComma:
		return ","
	case kColon:
		return ":"
	case kStar:
		return "*"
	case kLParen:
		return "("
	case kRParen:
		return ")"
	case kLBrace:
		return "{"
	case kRBrace:
		return "}"
	case kLBrack:
		return "["
	case kRBrack:
		return "]"
	case kLt:
		return "<"
	case kGt:
		return ">"
	case kDots:
		return "..."
	case kBang:
		return "!"
	default:
		return fmt.Sprintf("token(%d)", t.kind)
	}
}

func lex(src string) ([]token, error) {
	l := &lexer{src: src, line: 1, col: 1}
	var toks []token
	for {
		t, err := l.next()
		if err != nil {
			return nil, err
		}
		toks = append(toks, t)
		if t.kind == kEOF {
			return toks, nil
		}
	}
}

type lexer struct {
	src  string
	i    int
	line int
	col  int
}

func (l *lexer) next() (token, error) {
	for {
		l.skipSpaceAndComment()
		if l.i >= len(l.src) {
			return token{kind: kEOF, line: l.line, col: l.col}, nil
		}
		line, col := l.line, l.col
		c := l.src[l.i]
		switch c {
		case '=':
			l.adv()
			return token{kind: kEq, line: line, col: col}, nil
		case ',':
			l.adv()
			return token{kind: kComma, line: line, col: col}, nil
		case ':':
			l.adv()
			return token{kind: kColon, line: line, col: col}, nil
		case '*':
			l.adv()
			return token{kind: kStar, line: line, col: col}, nil
		case '(':
			l.adv()
			return token{kind: kLParen, line: line, col: col}, nil
		case ')':
			l.adv()
			return token{kind: kRParen, line: line, col: col}, nil
		case '{':
			l.adv()
			return token{kind: kLBrace, line: line, col: col}, nil
		case '}':
			l.adv()
			return token{kind: kRBrace, line: line, col: col}, nil
		case '[':
			l.adv()
			return token{kind: kLBrack, line: line, col: col}, nil
		case ']':
			l.adv()
			return token{kind: kRBrack, line: line, col: col}, nil
		case '<':
			l.adv()
			return token{kind: kLt, line: line, col: col}, nil
		case '>':
			l.adv()
			return token{kind: kGt, line: line, col: col}, nil
		case '%':
			l.adv()
			s, err := l.ident()
			if err != nil {
				return token{}, l.err(err)
			}
			return token{kind: kLocal, s: s, line: line, col: col}, nil
		case '@':
			l.adv()
			s, err := l.ident()
			if err != nil {
				return token{}, l.err(err)
			}
			return token{kind: kGlobal, s: s, line: line, col: col}, nil
		case '#':
			l.adv()
			n, err := l.uint()
			if err != nil {
				return token{}, l.err(err)
			}
			return token{kind: kAttrID, i: n, line: line, col: col}, nil
		case '!':
			l.adv()
			if l.i < len(l.src) && l.src[l.i] == '"' {
				s, err := l.quoted()
				if err != nil {
					return token{}, l.err(err)
				}
				return token{kind: kMetaName, s: s, line: line, col: col}, nil
			}
			if l.i < len(l.src) && isDigit(l.src[l.i]) {
				n, err := l.uint()
				if err != nil {
					return token{}, l.err(err)
				}
				return token{kind: kMetaID, i: n, line: line, col: col}, nil
			}
			if l.i < len(l.src) && isIdentStart(l.src[l.i]) {
				s, err := l.unquotedIdent()
				if err != nil {
					return token{}, l.err(err)
				}
				return token{kind: kMetaName, s: s, line: line, col: col}, nil
			}
			return token{kind: kBang, line: line, col: col}, nil
		case '"':
			s, err := l.quoted()
			if err != nil {
				return token{}, l.err(err)
			}
			return token{kind: kString, s: s, line: line, col: col}, nil
		case '.':
			if l.i+2 < len(l.src) && l.src[l.i+1] == '.' && l.src[l.i+2] == '.' {
				l.adv()
				l.adv()
				l.adv()
				return token{kind: kDots, line: line, col: col}, nil
			}
			s, err := l.unquotedIdent()
			if err != nil {
				return token{}, l.err(err)
			}
			return token{kind: kIdent, s: s, line: line, col: col}, nil
		case '-':
			if l.i+1 < len(l.src) && isDigit(l.src[l.i+1]) {
				return l.number(line, col)
			}
			l.adv()
			return token{kind: kIdent, s: "-", line: line, col: col}, nil
		default:
			if c == 'c' && l.i+1 < len(l.src) && l.src[l.i+1] == '"' {
				l.adv()
				b, err := l.quotedBytes()
				if err != nil {
					return token{}, l.err(err)
				}
				return token{kind: kCString, b: b, line: line, col: col}, nil
			}
			if c == 'i' && l.i+1 < len(l.src) && isDigit(l.src[l.i+1]) {
				j := l.i + 1
				for j < len(l.src) && isDigit(l.src[j]) {
					j++
				}
				// i32 is a type only if not followed by ident chars.
				if j >= len(l.src) || !isIdentCont(l.src[j]) {
					l.i++
					l.col++
					n, err := l.uint()
					if err != nil {
						return token{}, l.err(err)
					}
					return token{kind: kTypeInt, i: n, line: line, col: col}, nil
				}
			}
			if isDigit(c) {
				return l.number(line, col)
			}
			if isIdentStart(c) {
				s, err := l.unquotedIdent()
				if err != nil {
					return token{}, l.err(err)
				}
				return token{kind: kIdent, s: s, line: line, col: col}, nil
			}
			return token{}, fmt.Errorf("%d:%d: unexpected %q", line, col, c)
		}
	}
}

func (l *lexer) skipSpaceAndComment() {
	for l.i < len(l.src) {
		c := l.src[l.i]
		switch c {
		case ' ', '\t', '\r':
			l.adv()
		case '\n':
			l.adv()
			l.line++
			l.col = 1
		case ';':
			for l.i < len(l.src) && l.src[l.i] != '\n' {
				l.adv()
			}
		default:
			return
		}
	}
}

func (l *lexer) ident() (string, error) {
	if l.i >= len(l.src) {
		return "", fmt.Errorf("empty identifier")
	}
	if l.src[l.i] == '"' {
		return l.quoted()
	}
	return l.unquotedIdent()
}

func (l *lexer) unquotedIdent() (string, error) {
	if l.i >= len(l.src) || !isIdentStart(l.src[l.i]) && !isDigit(l.src[l.i]) {
		return "", fmt.Errorf("empty identifier")
	}
	start := l.i
	for l.i < len(l.src) && isIdentCont(l.src[l.i]) {
		l.adv()
	}
	return l.src[start:l.i], nil
}

func (l *lexer) quoted() (string, error) {
	b, err := l.quotedBytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (l *lexer) quotedBytes() ([]byte, error) {
	if l.i >= len(l.src) || l.src[l.i] != '"' {
		return nil, fmt.Errorf("expected string")
	}
	l.adv()
	var out []byte
	for l.i < len(l.src) {
		c := l.src[l.i]
		if c == '"' {
			l.adv()
			return out, nil
		}
		if c == '\\' {
			if l.i+2 >= len(l.src) {
				return nil, fmt.Errorf("unterminated escape")
			}
			h1, h2 := l.src[l.i+1], l.src[l.i+2]
			if isHex(h1) && isHex(h2) {
				out = append(out, unhex(h1)<<4|unhex(h2))
				l.adv()
				l.adv()
				l.adv()
				continue
			}
			out = append(out, h1)
			l.adv()
			l.adv()
			continue
		}
		if c == '\n' {
			return nil, fmt.Errorf("newline in string")
		}
		out = append(out, c)
		l.adv()
	}
	return nil, fmt.Errorf("unterminated string")
}

func (l *lexer) number(line, col int) (token, error) {
	start := l.i
	if l.src[l.i] == '-' {
		l.adv()
	}
	for l.i < len(l.src) && isDigit(l.src[l.i]) {
		l.adv()
	}
	// hex int 0x... is not used as a bare token in our fixtures; float 1.0e+2 is.
	if l.i < len(l.src) && (l.src[l.i] == '.' || l.src[l.i] == 'e' || l.src[l.i] == 'E') {
		if l.src[l.i] == '.' {
			l.adv()
			for l.i < len(l.src) && isDigit(l.src[l.i]) {
				l.adv()
			}
		}
		if l.i < len(l.src) && (l.src[l.i] == 'e' || l.src[l.i] == 'E') {
			l.adv()
			if l.i < len(l.src) && (l.src[l.i] == '+' || l.src[l.i] == '-') {
				l.adv()
			}
			for l.i < len(l.src) && isDigit(l.src[l.i]) {
				l.adv()
			}
		}
		return token{kind: kFloat, s: l.src[start:l.i], line: line, col: col}, nil
	}
	s := l.src[start:l.i]
	var n int64
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}
	for i := 0; i < len(s); i++ {
		n = n*10 + int64(s[i]-'0')
	}
	if neg {
		n = -n
	}
	return token{kind: kInt, i: n, s: l.src[start:l.i], line: line, col: col}, nil
}

func (l *lexer) uint() (int64, error) {
	if l.i >= len(l.src) || !isDigit(l.src[l.i]) {
		return 0, fmt.Errorf("expected integer")
	}
	var n int64
	for l.i < len(l.src) && isDigit(l.src[l.i]) {
		n = n*10 + int64(l.src[l.i]-'0')
		l.adv()
	}
	return n, nil
}

func (l *lexer) adv() {
	if l.i < len(l.src) {
		r, w := utf8.DecodeRuneInString(l.src[l.i:])
		l.i += w
		if r != '\n' {
			l.col++
		}
	}
}

func (l *lexer) err(err error) error {
	return fmt.Errorf("%d:%d: %w", l.line, l.col, err)
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func isHex(c byte) bool {
	return isDigit(c) || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

func unhex(c byte) byte {
	switch {
	case isDigit(c):
		return c - '0'
	case c >= 'a':
		return c - 'a' + 10
	default:
		return c - 'A' + 10
	}
}

func isIdentStart(c byte) bool {
	return c == '-' || c == '$' || c == '.' || c == '_' ||
		c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z'
}

func isIdentCont(c byte) bool {
	return isIdentStart(c) || isDigit(c)
}
