; ModuleID = 'testdata/ir/c_cxa_atexit/source.c'
source_filename = "testdata/ir/c_cxa_atexit/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@g = global i32 0
@__dso_handle = global i8 0

define void @dtor(ptr %p) {
  store i32 1, ptr %p, align 4
  ret void
}

define i32 @main() {
entry:
  %r = call i32 @__cxa_atexit(ptr @dtor, ptr @g, ptr @__dso_handle)
  ret i32 %r
}

declare i32 @__cxa_atexit(ptr, ptr, ptr)
