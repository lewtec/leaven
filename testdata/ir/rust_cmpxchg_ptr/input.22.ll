; ModuleID = 'testdata/ir/rust_cmpxchg_ptr/source.rs'
source_filename = "testdata/ir/rust_cmpxchg_ptr/source.rs"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@slot = global ptr null

define i32 @main() {
entry:
  %r = cmpxchg ptr @slot, ptr null, ptr @slot acquire monotonic
  %ok = extractvalue { ptr, i1 } %r, 1
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
