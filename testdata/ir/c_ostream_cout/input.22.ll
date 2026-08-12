; ModuleID = 'testdata/ir/c_ostream_cout/source.c'
source_filename = "testdata/ir/c_ostream_cout/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@_ZSt4cout = external global [64 x i8]
@msg = private unnamed_addr constant [3 x i8] c"/*\00"

define i32 @main() {
entry:
  %vp = load ptr, ptr @_ZSt4cout
  %slot = getelementptr i8, ptr %vp, i64 -24
  %off = load i64, ptr %slot
  %z = icmp eq i64 %off, 0
  br i1 %z, label %write, label %bad

write:
  %p = call ptr @_ZSt16__ostream_insertIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_PKS3_l(ptr @_ZSt4cout, ptr @msg, i64 2)
  ret i32 0

bad:
  ret i32 1
}

declare ptr @_ZSt16__ostream_insertIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_PKS3_l(ptr, ptr, i64)
