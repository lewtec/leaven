; ModuleID = 'testdata/ir/c_i1_blank/source.c'
source_filename = "testdata/ir/c_i1_blank/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define i1 @zero_i1() {
entry:
  ret i1 undef
}

define i32 @use_blank() {
entry:
  %_ = add i32 1, 2
  ret i32 %_
}

define i32 @main() {
entry:
  %z = call i1 @zero_i1()
  %n = call i32 @use_blank()
  br i1 %z, label %bad, label %chk

chk:
  %ok = icmp eq i32 %n, 3
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
