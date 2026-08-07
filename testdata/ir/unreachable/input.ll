; LLVM unreachable as noreturn / dead branch (clang-style)
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define dso_local void @die() {
entry:
  unreachable
}

define dso_local i32 @f(i32 noundef %x) {
entry:
  %cmp = icmp slt i32 %x, 0
  br i1 %cmp, label %if.then, label %if.end

if.then:
  unreachable

if.end:
  %add = add nsw i32 %x, 1
  ret i32 %add
}
