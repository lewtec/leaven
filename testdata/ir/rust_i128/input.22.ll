; ModuleID = 'testdata/ir/rust_i128/source.rs'
source_filename = "testdata/ir/rust_i128/source.rs"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define i32 @main() {
entry:
  %_23.i = add i128 1, 2
  %t = add i128 %_23.i, -1
  %ok = icmp eq i128 %t, 2
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
