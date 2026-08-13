; ModuleID = 'testdata/ir/c_fence/source.c'
source_filename = "testdata/ir/c_fence/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define void @f(ptr %p) {
entry:
  fence acquire
  %x = load ptr, ptr %p, align 8
  ret void
}

define i32 @main() {
entry:
  ret i32 0
}
