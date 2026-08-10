; ModuleID = 'testdata/ir/c_call_agg_ret/source.c'
source_filename = "testdata/ir/c_call_agg_ret/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

declare ptr @_Znwm(i64)

define { { ptr, ptr, ptr } } @make_impl() {
entry:
  ret { { ptr, ptr, ptr } } zeroinitializer
}

define ptr @nw(i64 %n) {
entry:
  %p = call ptr @_Znwm(i64 %n)
  ret ptr %p
}

define i32 @main() {
entry:
  %v = call { { ptr, ptr, ptr } } @make_impl()
  %p = call ptr @nw(i64 8)
  ret i32 0
}
