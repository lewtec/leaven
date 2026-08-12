; ModuleID = 'testdata/ir/c_ostream_insert/source.c'
source_filename = "testdata/ir/c_ostream_insert/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@msg = private unnamed_addr constant [4 x i8] c"/*\0A\00"
@dummy = global [8 x i8] zeroinitializer

define i32 @main() {
entry:
  %p = call ptr @_ZSt16__ostream_insertIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_PKS3_l(ptr @dummy, ptr @msg, i64 3)
  ret i32 0
}

declare ptr @_ZSt16__ostream_insertIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_PKS3_l(ptr, ptr, i64)
