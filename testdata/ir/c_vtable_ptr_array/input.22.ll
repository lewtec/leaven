; ModuleID = 'testdata/ir/c_vtable_ptr_array/source.c'
source_filename = "testdata/ir/c_vtable_ptr_array/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@str = private unnamed_addr constant [2 x i8] c"x\00"
@typeinfo = constant { ptr, ptr } zeroinitializer

define void @dtor() {
entry:
  ret void
}

@vt = constant [4 x ptr] [ptr @str, ptr @typeinfo, ptr @dtor, ptr null]

define i32 @main() {
entry:
  ret i32 0
}
