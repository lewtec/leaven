; ModuleID = 'testdata/ir/c_srand48/source.c'
source_filename = "testdata/ir/c_srand48/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define i32 @main() {
entry:
  call void @srand48(i64 1)
  %a = call i64 @lrand48()
  %ok = icmp eq i64 %a, 89400484
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}

declare void @srand48(i64)
declare i64 @lrand48()
