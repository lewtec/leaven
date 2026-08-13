; ModuleID = 'testdata/ir/c_sscanf/source.c'
source_filename = "testdata/ir/c_sscanf/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@src = private unnamed_addr constant [2 x i8] c"1\00"
@fmt = private unnamed_addr constant [4 x i8] c"%lu\00"
@v = global i64 0

define i32 @main() {
entry:
  %r = call i32 (ptr, ptr, ...) @__isoc23_sscanf(ptr @src, ptr @fmt, ptr @v)
  %ok = icmp eq i32 %r, 1
  br i1 %ok, label %chk, label %bad

chk:
  %n = load i64, ptr @v
  %eq = icmp eq i64 %n, 1
  br i1 %eq, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}

declare i32 @__isoc23_sscanf(ptr, ptr, ...)
