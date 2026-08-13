; ModuleID = 'testdata/ir/c_vtable_init_cycle/source.c'
source_filename = "testdata/ir/c_vtable_init_cycle/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@vt = global [2 x ptr] [ptr null, ptr @dtor]

define void @dtor(ptr %this) {
entry:
  store ptr getelementptr inbounds (ptr, ptr @vt, i64 2), ptr %this
  ret void
}

define i32 @main() {
entry:
  ret i32 0
}
