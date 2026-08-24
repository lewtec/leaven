package leaven

import (
	"strings"

	"github.com/ianlancetaylor/demangle"
)

// cxx is a parsed Itanium C++ function symbol.
type cxx struct {
	recv   string // class without template args; empty for free functions
	ident  string // method, operator ("<<"), or free name ("endl")
	ctor   bool
	dtor   bool
	args   []cxxTy
	ret    cxxTy
	std    bool
	libcxx bool   // std::__1
	valV   string // V in __tree<__value_type<K,V>>
}

type cxxTy struct {
	ident string // "char", "unsigned long", "basic_ostream", "fnptr"
	ptr   bool
	ref   bool
}

func parseCxx(name string) (cxx, bool) {
	if !strings.HasPrefix(name, "_Z") && !strings.HasPrefix(name, "__Z") {
		return cxx{}, false
	}
	ast, err := demangle.ToAST(name)
	if err != nil {
		// GNU internal-linkage templates (std::__1::__tree_next) trip
		// substitution errors; the name is still in the AST.
		ast, err = demangle.ToAST(name, demangle.NoParams)
		if err != nil {
			return cxx{}, false
		}
	}
	if cl, ok := ast.(*demangle.Clone); ok {
		ast = cl.Base
	}
	return walkCxx(ast)
}

func walkCxx(ast demangle.AST) (cxx, bool) {
	if t, ok := ast.(*demangle.Typed); ok {
		c := fromCxxName(t.Name)
		c.args = cxxFnArgs(t.Type)
		c.ret = cxxReturn(t.Type)
		if c.ident == "" && !c.ctor && !c.dtor {
			return cxx{}, false
		}
		return c, true
	}
	c := fromCxxName(ast)
	if c.ident == "" && !c.ctor && !c.dtor {
		return cxx{}, false
	}
	return c, true
}

func cxxReturn(a demangle.AST) cxxTy {
	switch t := a.(type) {
	case *demangle.MethodWithQualifiers:
		return cxxReturn(t.Method)
	case *demangle.FunctionType:
		if t.Return == nil {
			return cxxTy{}
		}
		return cxxFromTy(t.Return)
	default:
		return cxxTy{}
	}
}

func fromCxxName(a demangle.AST) cxx {
	switch n := a.(type) {
	case *demangle.TaggedName:
		return fromCxxName(n.Name)
	case *demangle.Template:
		return fromCxxName(n.Name)
	case *demangle.Qualified:
		c := fromCxxName(n.Name)
		if c.recv == "" {
			c.recv = cxxClassOf(n.Scope)
		}
		if c.valV == "" {
			c.valV = cxxValueV(n.Scope)
		}
		c.markScope(n)
		return c
	case *demangle.Operator:
		return cxx{ident: strings.TrimSpace(n.Name)}
	case *demangle.Name:
		return cxx{ident: n.Name}
	case *demangle.Constructor:
		return cxx{ctor: true, ident: cxxLeaf(n.Name)}
	case *demangle.Destructor:
		return cxx{dtor: true, ident: cxxLeaf(n.Name)}
	default:
		return cxx{}
	}
}

func (c *cxx) markScope(a demangle.AST) {
	var path []string
	cxxCollect(a, &path)
	for i, s := range path {
		if s != "std" {
			continue
		}
		c.std = true
		if i+1 < len(path) && path[i+1] == "__1" {
			c.libcxx = true
		}
	}
}

func cxxCollect(a demangle.AST, path *[]string) {
	switch n := a.(type) {
	case *demangle.Qualified:
		cxxCollect(n.Scope, path)
		cxxCollect(n.Name, path)
	case *demangle.Template:
		cxxCollect(n.Name, path)
	case *demangle.TaggedName:
		cxxCollect(n.Name, path)
	case *demangle.Name:
		*path = append(*path, n.Name)
	}
}

func cxxClassOf(a demangle.AST) string {
	switch n := a.(type) {
	case *demangle.TaggedName:
		return cxxClassOf(n.Name)
	case *demangle.Template:
		return cxxClassOf(n.Name)
	case *demangle.Qualified:
		return cxxClassOf(n.Name)
	case *demangle.Name:
		if cxxNS[n.Name] {
			return ""
		}
		return n.Name
	case *demangle.Constructor:
		return cxxClassOf(n.Name)
	default:
		return ""
	}
}

// cxxValueV is V in __value_type<K, V>, used to tell a trivially
// copyable map pair (unsigned) from Effect (owns a vector).
func cxxValueV(a demangle.AST) string {
	switch n := a.(type) {
	case *demangle.Template:
		if cxxClassOf(n.Name) == "__value_type" && len(n.Args) >= 2 {
			return cxxFromTy(n.Args[1]).ident
		}
		for _, arg := range n.Args {
			if s := cxxValueV(arg); s != "" {
				return s
			}
		}
	case *demangle.Qualified:
		if s := cxxValueV(n.Name); s != "" {
			return s
		}
		return cxxValueV(n.Scope)
	case *demangle.TaggedName:
		return cxxValueV(n.Name)
	}
	return ""
}

func (n cxx) trivialTreeValue() bool {
	switch n.valV {
	case "unsigned int", "unsigned long", "unsigned long long",
		"int", "long", "long long", "unsigned short", "bool":
		return true
	default:
		return false
	}
}

func cxxLeaf(a demangle.AST) string {
	switch n := a.(type) {
	case *demangle.TaggedName:
		return cxxLeaf(n.Name)
	case *demangle.Template:
		return cxxLeaf(n.Name)
	case *demangle.Name:
		return n.Name
	default:
		return ""
	}
}

var cxxNS = map[string]bool{
	"std":       true,
	"__1":       true,
	"__cxx11":   true,
	"__gnu_cxx": true,
	"__detail":  true,
}

func cxxFnArgs(a demangle.AST) []cxxTy {
	switch t := a.(type) {
	case *demangle.MethodWithQualifiers:
		return cxxFnArgs(t.Method)
	case *demangle.FunctionType:
		if len(t.Args) == 0 {
			return nil
		}
		out := make([]cxxTy, len(t.Args))
		for i, arg := range t.Args {
			out[i] = cxxFromTy(arg)
		}
		return out
	default:
		return nil
	}
}

func cxxFromTy(a demangle.AST) cxxTy {
	var t cxxTy
	for a != nil {
		switch n := a.(type) {
		case *demangle.ReferenceType, *demangle.RvalueReferenceType:
			t.ref = true
			if r, ok := n.(*demangle.ReferenceType); ok {
				a = r.Base
			} else {
				a = n.(*demangle.RvalueReferenceType).Base
			}
		case *demangle.PointerType:
			if _, ok := cxxUnwrapTy(n.Base).(*demangle.FunctionType); ok {
				t.ident = "fnptr"
				t.ptr = true
				return t
			}
			t.ptr = true
			a = n.Base
		case *demangle.TypeWithQualifiers:
			a = n.Base
		case *demangle.VendorQualifier:
			a = n.Type
		case *demangle.Template:
			t.ident = cxxClassOf(n.Name)
			return t
		case *demangle.Qualified:
			t.ident = cxxClassOf(n)
			return t
		case *demangle.BuiltinType:
			t.ident = n.Name
			return t
		case *demangle.Name:
			t.ident = n.Name
			return t
		default:
			return t
		}
	}
	return t
}

func cxxUnwrapTy(a demangle.AST) demangle.AST {
	for {
		switch n := a.(type) {
		case *demangle.TypeWithQualifiers:
			a = n.Base
		case *demangle.VendorQualifier:
			a = n.Type
		default:
			return a
		}
	}
}

func cxxCanon(s string) string {
	switch s {
	case "ostream":
		return "basic_ostream"
	case "istream":
		return "basic_istream"
	case "ifstream":
		return "basic_ifstream"
	case "ofstream":
		return "basic_ofstream"
	case "ostringstream":
		return "basic_ostringstream"
	case "stringstream":
		return "basic_stringstream"
	case "istringstream":
		return "basic_istringstream"
	case "stringbuf":
		return "basic_stringbuf"
	case "filebuf":
		return "basic_filebuf"
	case "streambuf":
		return "basic_streambuf"
	case "ios":
		return "basic_ios"
	case "string":
		return "basic_string"
	default:
		return s
	}
}

func (c cxx) class() string { return cxxCanon(c.recv) }

func (c cxx) isClass(names ...string) bool {
	got := c.class()
	for _, n := range names {
		if got == n || got == cxxCanon(n) {
			return true
		}
	}
	return false
}

func (t cxxTy) is(name string) bool {
	return t.ident == name || cxxCanon(t.ident) == cxxCanon(name)
}

func (t cxxTy) isCStr() bool {
	return t.ident == "char" && t.ptr && !t.ref
}

func (t cxxTy) isString() bool {
	return cxxCanon(t.ident) == "basic_string"
}

func (t cxxTy) isOstream() bool {
	return cxxCanon(t.ident) == "basic_ostream"
}

func (c cxx) hasCStr() bool {
	for _, a := range c.args {
		if a.isCStr() {
			return true
		}
	}
	return false
}

func (c cxx) hasString() bool {
	for _, a := range c.args {
		if a.isString() {
			return true
		}
	}
	return false
}

// streamArg is the value operand of operator<< on an ostream
// (method: the only arg; free: the arg after the stream).
func (c cxx) streamArg() (cxxTy, bool) {
	if c.ident != "<<" {
		return cxxTy{}, false
	}
	if c.class() == "basic_ostream" {
		if len(c.args) != 1 {
			return cxxTy{}, false
		}
		return c.args[0], true
	}
	if c.recv == "" && len(c.args) >= 2 && c.args[0].isOstream() {
		return c.args[1], true
	}
	return cxxTy{}, false
}
