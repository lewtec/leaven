; ModuleID = 'testdata/ir/rust_atomic_i8/source.rs'
source_filename = "testdata/ir/rust_atomic_i8/source.rs"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define i32 @main() {
entry:
  %p = alloca i8, align 1
  store i8 3, ptr %p, align 1
  %old = atomicrmw sub ptr %p, i8 1 acquire
  %ok = icmp eq i8 %old, 3
  br i1 %ok, label %chk, label %bad

chk:
  %now = load i8, ptr %p, align 1
  %ok2 = icmp eq i8 %now, 2
  br i1 %ok2, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
