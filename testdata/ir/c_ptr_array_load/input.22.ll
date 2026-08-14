; ModuleID = 'testdata/ir/c_ptr_array_load/source.c'
source_filename = "testdata/ir/c_ptr_array_load/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@msg = private unnamed_addr constant [3 x i8] c"ok\00"
@arr = internal global [2 x ptr] [ptr @msg, ptr null]

define i32 @main() {
entry:
  %p = load ptr, ptr getelementptr inbounds ([2 x ptr], ptr @arr, i64 0, i64 0)
  %c = load i8, ptr %p
  %ok = icmp eq i8 %c, 111
  %r = zext i1 %ok to i32
  %inv = xor i32 %r, 1
  ret i32 %inv
}
