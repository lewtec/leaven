; ModuleID = 'testdata/ir/rust_statx/source.rs'
source_filename = "testdata/ir/rust_statx/source.rs"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@path = private unnamed_addr constant [2 x i8] c".\00"
@buf = global [256 x i8] zeroinitializer

define i32 @main() {
entry:
  %r = call i32 @statx(i32 -100, ptr @path, i32 0, i32 2047, ptr @buf)
  %ok = icmp sge i32 %r, 0
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}

declare extern_weak noundef i32 @statx(i32 noundef, ptr noundef, i32 noundef, i32 noundef, ptr noundef) unnamed_addr
