; ModuleID = 'testdata/ir/c_store_vec_start/source.c'
source_filename = "testdata/ir/c_store_vec_start/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@AllVars = global { { ptr, ptr, ptr } } zeroinitializer

define void @grow(ptr %p) {
entry:
  store ptr %p, ptr @AllVars
  ret void
}

define ptr @start() {
entry:
  %0 = load ptr, ptr @AllVars
  ret ptr %0
}

define i32 @main() {
entry:
  %a = alloca i32, align 4
  call void @grow(ptr %a)
  %s = call ptr @start()
  %ok = icmp eq ptr %s, %a
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
