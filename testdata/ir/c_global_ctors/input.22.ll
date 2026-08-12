; ModuleID = 'testdata/ir/c_global_ctors/source.c'
source_filename = "testdata/ir/c_global_ctors/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@g = global i32 0
@llvm.global_ctors = appending global [1 x { i32, ptr, ptr }] [{ i32, ptr, ptr } { i32 65535, ptr @ctor, ptr null }]

define internal void @ctor() {
  store i32 1, ptr @g, align 4
  ret void
}

define i32 @main() {
  %v = load i32, ptr @g, align 4
  %ok = icmp eq i32 %v, 1
  br i1 %ok, label %good, label %bad
good:
  ret i32 0
bad:
  ret i32 1
}
