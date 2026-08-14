; ModuleID = 'testdata/ir/c_null_ptr_slot/source.c'
source_filename = "testdata/ir/c_null_ptr_slot/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

%slot = type { ptr }

@g = internal global %slot { ptr undef }

define i32 @main() {
entry:
  store ptr null, ptr @g
  %p = load ptr, ptr @g
  %ok = icmp eq ptr %p, null
  %r = zext i1 %ok to i32
  %inv = xor i32 %r, 1
  ret i32 %inv
}
