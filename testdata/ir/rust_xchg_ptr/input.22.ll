; ModuleID = 'testdata/ir/rust_xchg_ptr/source.rs'
source_filename = "testdata/ir/rust_xchg_ptr/source.rs"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@slot = global ptr null
@val = global i32 0

define i32 @main() {
entry:
  %old = atomicrmw xchg ptr @slot, ptr @val acq_rel
  %ok = icmp eq ptr %old, null
  br i1 %ok, label %chk, label %bad

chk:
  %now = load ptr, ptr @slot
  %hit = icmp eq ptr %now, @val
  br i1 %hit, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
