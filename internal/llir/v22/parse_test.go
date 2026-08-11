package v22

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/lewtec/leaven/internal/llir/ir"
	"github.com/lewtec/leaven/internal/llir/ir/constant"
	"github.com/lewtec/leaven/internal/llir/ir/enum"
	"github.com/lewtec/leaven/internal/llir/ir/types"
)

func TestParseStrlen22(t *testing.T) {
	root := findRepoRoot(t)
	b, err := os.ReadFile(filepath.Join(root, "testdata/ir/c_strlen_map/input.22.ll"))
	if err != nil {
		t.Fatal(err)
	}
	m, err := ParseString("input.22.ll", string(b))
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Funcs) == 0 {
		t.Fatal("no funcs")
	}
	n := m.Funcs[0]
	if len(n.Params) != 1 {
		t.Fatalf("params %d", len(n.Params))
	}
	pt, ok := n.Params[0].Typ.(*types.PointerType)
	if !ok || !pt.IsOpaque() {
		t.Fatalf("param type %v opaque=%v", n.Params[0].Typ, ok && pt.IsOpaque())
	}
}

func TestParseModuleAsmAndAlias(t *testing.T) {
	src := `target triple = "x86_64-unknown-linux-gnu"
module asm ".globl foo"
@real = dso_local unnamed_addr alias void (ptr), ptr @impl
define void @impl(ptr %p) {
entry:
  ret void
}
`
	m, err := ParseString("alias.ll", src)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.ModuleAsms) != 1 || m.ModuleAsms[0] != ".globl foo" {
		t.Fatalf("module asm: %#v", m.ModuleAsms)
	}
	if len(m.Aliases) != 1 || m.Aliases[0].Name() != "real" {
		t.Fatalf("aliases: %+v", m.Aliases)
	}
}

func TestParseRustAdd22(t *testing.T) {
	root := findRepoRoot(t)
	b, err := os.ReadFile(filepath.Join(root, "testdata/ir/rust_add/input.22.ll"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseString("input.22.ll", string(b)); err != nil {
		t.Fatal(err)
	}
}

func TestParseUseBeforeDef(t *testing.T) {
	src := `define i64 @f() {
entry:
  br label %later
early:
  %inc = add i64 %i, 1
  ret i64 %inc
later:
  %i = add i64 0, 0
  br label %early
}
`
	m, err := ParseString("uad.ll", src)
	if err != nil {
		t.Fatal(err)
	}
	add := m.Funcs[0].Blocks[1].Insts[0].(*ir.InstAdd)
	if _, ok := add.X.(*ir.InstAdd); !ok {
		t.Fatalf("%%inc lhs %T, want defining add", add.X)
	}
}

func TestParseRelocConstExpr(t *testing.T) {
	src := `@s = constant [2 x i8] c"x\00"
@tab = constant [1 x i32] [i32 trunc (i64 sub (i64 ptrtoint (ptr @s to i64), i64 ptrtoint (ptr @tab to i64)) to i32)]
`
	m, err := ParseString("reloc.ll", src)
	if err != nil {
		t.Fatal(err)
	}
	arr, ok := m.Globals[1].Init.(*constant.Array)
	if !ok || len(arr.Elems) != 1 {
		t.Fatalf("tab init %T", m.Globals[1].Init)
	}
	tr, ok := arr.Elems[0].(*constant.ExprTrunc)
	if !ok {
		t.Fatalf("elem %T, want trunc", arr.Elems[0])
	}
	if _, ok := tr.From.(*constant.ExprSub); !ok {
		t.Fatalf("trunc from %T, want sub", tr.From)
	}
}

func TestParseShuffleVector(t *testing.T) {
	src := `define <4 x i32> @f(<4 x i32> %a) {
entry:
  %r = shufflevector <4 x i32> %a, <4 x i32> poison, <4 x i32> <i32 3, i32 2, i32 1, i32 0>
  ret <4 x i32> %r
}
`
	m, err := ParseString("shuf.ll", src)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m.Funcs[0].Blocks[0].Insts[0].(*ir.InstShuffleVector); !ok {
		t.Fatalf("inst %T", m.Funcs[0].Blocks[0].Insts[0])
	}
}

func TestParseFMul(t *testing.T) {
	src := `define double @f(double %x) {
entry:
  %m = fmul double %x, 1.000000e+02
  ret double %m
}
`
	m, err := ParseString("fmul.ll", src)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m.Funcs[0].Blocks[0].Insts[0].(*ir.InstFMul); !ok {
		t.Fatalf("inst %T", m.Funcs[0].Blocks[0].Insts[0])
	}
}

func TestParseZextNneg(t *testing.T) {
	src := `define i64 @f(i32 %x) {
entry:
  %z = zext nneg i32 %x to i64
  ret i64 %z
}
`
	if _, err := ParseString("nneg.ll", src); err != nil {
		t.Fatal(err)
	}
}

func TestParseLoadAtomic(t *testing.T) {
	// rustc once_cell: load atomic ptr, ptr @SEEDS acquire, align 8
	src := `@seeds = internal global ptr null
define ptr @f() {
entry:
  %p = load atomic ptr, ptr @seeds acquire, align 8
  store atomic ptr %p, ptr @seeds release, align 8
  ret ptr %p
}
`
	m, err := ParseString("atomic.ll", src)
	if err != nil {
		t.Fatal(err)
	}
	ld, ok := m.Funcs[0].Blocks[0].Insts[0].(*ir.InstLoad)
	if !ok || !ld.Atomic || ld.Ordering != enum.AtomicOrderingAcquire {
		t.Fatalf("load %+v", m.Funcs[0].Blocks[0].Insts[0])
	}
	st, ok := m.Funcs[0].Blocks[0].Insts[1].(*ir.InstStore)
	if !ok || !st.Atomic || st.Ordering != enum.AtomicOrderingRelease {
		t.Fatalf("store %+v", m.Funcs[0].Blocks[0].Insts[1])
	}
}

func TestParseUnreachableDbg(t *testing.T) {
	src := `define void @f() {
entry:
  unreachable, !dbg !0
}
!0 = !{}
`
	if _, err := ParseString("unreach.ll", src); err != nil {
		t.Fatal(err)
	}
}

func TestParseGEPInrange(t *testing.T) {
	src := `@vt = external global [4 x ptr]
define ptr @f(ptr %p) {
entry:
  store ptr getelementptr inbounds nuw inrange(-16, 16) (i8, ptr @vt, i64 16), ptr %p
  ret ptr %p
}
`
	m, err := ParseString("gep.ll", src)
	if err != nil {
		t.Fatal(err)
	}
	st, ok := m.Funcs[0].Blocks[0].Insts[0].(*ir.InstStore)
	if !ok {
		t.Fatalf("inst %T", m.Funcs[0].Blocks[0].Insts[0])
	}
	if _, ok := st.Src.(*constant.ExprGetElementPtr); !ok {
		t.Fatalf("store src %T, want getelementptr", st.Src)
	}
}

func TestParseInvokeLandingPad(t *testing.T) {
	src := `define void @f(ptr %p) personality ptr @eh {
entry:
  invoke void @g(ptr %p)
          to label %ok unwind label %bad
ok:
  ret void
bad:
  %lp = landingpad { ptr, i32 }
          cleanup
  resume { ptr, i32 } %lp
}
declare void @g(ptr)
declare i32 @eh(...)
`
	m, err := ParseString("eh.ll", src)
	if err != nil {
		t.Fatal(err)
	}
	f := m.Funcs[0]
	if _, ok := f.Blocks[0].Term.(*ir.TermInvoke); !ok {
		t.Fatalf("entry term %T, want invoke", f.Blocks[0].Term)
	}
	if _, ok := f.Blocks[1].Term.(*ir.TermRet); !ok {
		t.Fatalf("ok term %T, want ret", f.Blocks[1].Term)
	}
	if _, ok := f.Blocks[2].Insts[0].(*ir.InstLandingPad); !ok {
		t.Fatalf("bad inst %T, want landingpad", f.Blocks[2].Insts[0])
	}
	if _, ok := f.Blocks[2].Term.(*ir.TermResume); !ok {
		t.Fatalf("bad term %T, want resume", f.Blocks[2].Term)
	}
}

func TestParseThreadLocal(t *testing.T) {
	// rustc: @x = internal thread_local unnamed_addr global ...
	src := `@tls = internal thread_local unnamed_addr global i32 0
@ie = thread_local(initialexec) global i32 1
`
	m, err := ParseString("tls.ll", src)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Globals) != 2 {
		t.Fatalf("globals=%d", len(m.Globals))
	}
	if m.Globals[0].TLSModel != enum.TLSModelGeneric {
		t.Fatalf("tls model=%v, want generic", m.Globals[0].TLSModel)
	}
	if m.Globals[1].TLSModel != enum.TLSModelInitialExec {
		t.Fatalf("ie model=%v, want initialexec", m.Globals[1].TLSModel)
	}
}

func TestParseExternalThenGlobal(t *testing.T) {
	src := `@ext = external global [0 x ptr]
@real = private constant i32 1
`
	m, err := ParseString("ext.ll", src)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Globals) != 2 {
		t.Fatalf("globals=%d", len(m.Globals))
	}
}

func TestParseForwardGlobal(t *testing.T) {
	// C++ vtables mention @_ZTI* before that global is defined.
	// Typeinfo inits use constant getelementptr.
	src := `@_ZTV = constant { [2 x ptr] } { [2 x ptr] [ptr null, ptr @_ZTI] }
@_ZTI = constant { ptr, ptr } { ptr getelementptr inbounds (ptr, ptr @_ZTV, i64 2), ptr @_ZTS }
@_ZTS = constant [4 x i8] c"Foo\00"
@small = constant i8 trunc (i32 300 to i8)
declare void @__cxa_pure_virtual()
`
	m, err := ParseString("fwd.ll", src)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Globals) < 3 {
		t.Fatalf("globals=%d", len(m.Globals))
	}
	zti := m.Globals[1]
	st, ok := zti.Init.(*constant.Struct)
	if !ok || len(st.Fields) < 1 {
		t.Fatalf("ZTI init %T", zti.Init)
	}
	if _, ok := st.Fields[0].(*constant.ExprGetElementPtr); !ok {
		t.Fatalf("ZTI field0 %T, want getelementptr", st.Fields[0])
	}
	if _, ok := m.Globals[3].Init.(*constant.ExprTrunc); !ok {
		t.Fatalf("trunc init %T", m.Globals[3].Init)
	}
}

func TestParseAll22Fixtures(t *testing.T) {
	root := findRepoRoot(t)
	dir := filepath.Join(root, "testdata", "ir")
	ents, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range ents {
		if !e.IsDir() {
			continue
		}
		path := filepath.Join(dir, e.Name(), "input.22.ll")
		b, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			t.Fatal(err)
		}
		t.Run(e.Name(), func(t *testing.T) {
			m, err := ParseString(path, string(b))
			if err != nil {
				t.Fatal(err)
			}
			if len(m.Funcs) == 0 {
				t.Fatal("no functions")
			}
			for _, f := range m.Funcs {
				if len(f.Blocks) == 0 {
					continue
				}
				for _, b := range f.Blocks {
					if b.Term == nil {
						t.Fatalf("%s: block %s has no terminator", f.Name(), b.Name())
					}
				}
			}
		})
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("no go.mod")
	return ""
}
