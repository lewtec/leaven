; ModuleID = 'testdata/ir/c_main_args/source.c'
source_filename = "testdata/ir/c_main_args/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define i32 @main(i32 %argc, ptr %argv) {
entry:
  %ok = icmp sge i32 %argc, 1
  br i1 %ok, label %chk, label %bad

chk:
  %p = load ptr, ptr %argv
  %nn = icmp ne ptr %p, null
  br i1 %nn, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
