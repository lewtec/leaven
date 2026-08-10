; ModuleID = 'testdata/ir/c_llvm_intrinsics/source.c'
source_filename = "testdata/ir/c_llvm_intrinsics/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@blob = private unnamed_addr constant [8 x i8] c"\04\00\00\00XXXX"

declare i64 @llvm.umin.i64(i64, i64)
declare i64 @llvm.umax.i64(i64, i64)
declare i64 @llvm.smin.i64(i64, i64)
declare i64 @llvm.smax.i64(i64, i64)
declare i32 @llvm.umin.i32(i32, i32)
declare i32 @llvm.umax.i32(i32, i32)
declare i32 @llvm.smin.i32(i32, i32)
declare i32 @llvm.smax.i32(i32, i32)
declare void @llvm.trap()
declare double @llvm.ceil.f64(double)
declare i32 @llvm.vector.reduce.add.v4i32(<4 x i32>)
declare ptr @llvm.load.relative.i64(ptr, i64)

define void @die() {
entry:
  call void @llvm.trap()
  ret void
}

define i32 @main() {
entry:
  %u = call i64 @llvm.umin.i64(i64 -1, i64 2)
  %ok1 = icmp eq i64 %u, 2
  br i1 %ok1, label %c2, label %bad

c2:
  %s = call i64 @llvm.smin.i64(i64 -2, i64 5)
  %ok2 = icmp eq i64 %s, -2
  br i1 %ok2, label %c3, label %bad

c3:
  %m = call i32 @llvm.umax.i32(i32 1, i32 9)
  %ok3 = icmp eq i32 %m, 9
  br i1 %ok3, label %c4, label %bad

c4:
  %ce = call double @llvm.ceil.f64(double 1.500000e+00)
  %ci = fptosi double %ce to i32
  %ok4 = icmp eq i32 %ci, 2
  br i1 %ok4, label %c5, label %bad

c5:
  %sum = call i32 @llvm.vector.reduce.add.v4i32(<4 x i32> <i32 1, i32 2, i32 3, i32 4>)
  %ok5 = icmp eq i32 %sum, 10
  br i1 %ok5, label %c6, label %bad

c6:
  %p = getelementptr inbounds [8 x i8], ptr @blob, i64 0, i64 0
  %r = call ptr @llvm.load.relative.i64(ptr %p, i64 0)
  %q = getelementptr inbounds i8, ptr %p, i64 4
  %ok6 = icmp eq ptr %r, %q
  br i1 %ok6, label %ok, label %bad

ok:
  ret i32 0

bad:
  ret i32 1
}
