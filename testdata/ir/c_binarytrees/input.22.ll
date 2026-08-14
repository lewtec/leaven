; ModuleID = 'testdata/ir/c_binarytrees/source.c'
source_filename = "testdata/ir/c_binarytrees/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

%struct.tn = type { ptr, ptr }

@.str = private unnamed_addr constant [38 x i8] c"stretch tree of depth %u\09 check: %li\0A\00", align 1
@.str.1 = private unnamed_addr constant [36 x i8] c"%li\09 trees of depth %u\09 check: %li\0A\00", align 1
@.str.2 = private unnamed_addr constant [41 x i8] c"long lived tree of depth %u\09 check: %li\0A\00", align 1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local ptr @NewTreeNode(ptr noundef %left, ptr noundef %right) #0 {
entry:
  %left.addr = alloca ptr, align 8
  %right.addr = alloca ptr, align 8
  %new = alloca ptr, align 8
  store ptr %left, ptr %left.addr, align 8
  store ptr %right, ptr %right.addr, align 8
  %call = call noalias ptr @malloc(i64 noundef 16) #4
  store ptr %call, ptr %new, align 8
  %0 = load ptr, ptr %left.addr, align 8
  %1 = load ptr, ptr %new, align 8
  %left1 = getelementptr inbounds nuw %struct.tn, ptr %1, i32 0, i32 0
  store ptr %0, ptr %left1, align 8
  %2 = load ptr, ptr %right.addr, align 8
  %3 = load ptr, ptr %new, align 8
  %right2 = getelementptr inbounds nuw %struct.tn, ptr %3, i32 0, i32 1
  store ptr %2, ptr %right2, align 8
  %4 = load ptr, ptr %new, align 8
  ret ptr %4
}

; Function Attrs: nounwind allocsize(0)
declare noalias ptr @malloc(i64 noundef) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i64 @ItemCheck(ptr noundef %tree) #0 {
entry:
  %retval = alloca i64, align 8
  %tree.addr = alloca ptr, align 8
  store ptr %tree, ptr %tree.addr, align 8
  %0 = load ptr, ptr %tree.addr, align 8
  %left = getelementptr inbounds nuw %struct.tn, ptr %0, i32 0, i32 0
  %1 = load ptr, ptr %left, align 8
  %cmp = icmp eq ptr %1, null
  br i1 %cmp, label %if.then, label %if.else

if.then:                                          ; preds = %entry
  store i64 1, ptr %retval, align 8
  br label %return

if.else:                                          ; preds = %entry
  %2 = load ptr, ptr %tree.addr, align 8
  %left1 = getelementptr inbounds nuw %struct.tn, ptr %2, i32 0, i32 0
  %3 = load ptr, ptr %left1, align 8
  %call = call i64 @ItemCheck(ptr noundef %3)
  %add = add nsw i64 1, %call
  %4 = load ptr, ptr %tree.addr, align 8
  %right = getelementptr inbounds nuw %struct.tn, ptr %4, i32 0, i32 1
  %5 = load ptr, ptr %right, align 8
  %call2 = call i64 @ItemCheck(ptr noundef %5)
  %add3 = add nsw i64 %add, %call2
  store i64 %add3, ptr %retval, align 8
  br label %return

return:                                           ; preds = %if.else, %if.then
  %6 = load i64, ptr %retval, align 8
  ret i64 %6
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local ptr @BottomUpTree(i32 noundef %depth) #0 {
entry:
  %retval = alloca ptr, align 8
  %depth.addr = alloca i32, align 4
  store i32 %depth, ptr %depth.addr, align 4
  %0 = load i32, ptr %depth.addr, align 4
  %cmp = icmp ugt i32 %0, 0
  br i1 %cmp, label %if.then, label %if.else

if.then:                                          ; preds = %entry
  %1 = load i32, ptr %depth.addr, align 4
  %sub = sub i32 %1, 1
  %call = call ptr @BottomUpTree(i32 noundef %sub)
  %2 = load i32, ptr %depth.addr, align 4
  %sub1 = sub i32 %2, 1
  %call2 = call ptr @BottomUpTree(i32 noundef %sub1)
  %call3 = call ptr @NewTreeNode(ptr noundef %call, ptr noundef %call2)
  store ptr %call3, ptr %retval, align 8
  br label %return

if.else:                                          ; preds = %entry
  %call4 = call ptr @NewTreeNode(ptr noundef null, ptr noundef null)
  store ptr %call4, ptr %retval, align 8
  br label %return

return:                                           ; preds = %if.else, %if.then
  %3 = load ptr, ptr %retval, align 8
  ret ptr %3
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @DeleteTree(ptr noundef %tree) #0 {
entry:
  %tree.addr = alloca ptr, align 8
  store ptr %tree, ptr %tree.addr, align 8
  %0 = load ptr, ptr %tree.addr, align 8
  %left = getelementptr inbounds nuw %struct.tn, ptr %0, i32 0, i32 0
  %1 = load ptr, ptr %left, align 8
  %cmp = icmp ne ptr %1, null
  br i1 %cmp, label %if.then, label %if.end

if.then:                                          ; preds = %entry
  %2 = load ptr, ptr %tree.addr, align 8
  %left1 = getelementptr inbounds nuw %struct.tn, ptr %2, i32 0, i32 0
  %3 = load ptr, ptr %left1, align 8
  call void @DeleteTree(ptr noundef %3)
  %4 = load ptr, ptr %tree.addr, align 8
  %right = getelementptr inbounds nuw %struct.tn, ptr %4, i32 0, i32 1
  %5 = load ptr, ptr %right, align 8
  call void @DeleteTree(ptr noundef %5)
  br label %if.end

if.end:                                           ; preds = %if.then, %entry
  %6 = load ptr, ptr %tree.addr, align 8
  call void @free(ptr noundef %6) #5
  ret void
}

; Function Attrs: nounwind
declare void @free(ptr noundef) #2

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %N = alloca i32, align 4
  %depth = alloca i32, align 4
  %minDepth = alloca i32, align 4
  %maxDepth = alloca i32, align 4
  %stretchDepth = alloca i32, align 4
  %stretchTree = alloca ptr, align 8
  %longLivedTree = alloca ptr, align 8
  %tempTree = alloca ptr, align 8
  %i = alloca i64, align 8
  %iterations = alloca i64, align 8
  %check = alloca i64, align 8
  %i7 = alloca i32, align 4
  store i32 0, ptr %retval, align 4
  store i32 10, ptr %N, align 4
  store i32 4, ptr %minDepth, align 4
  %0 = load i32, ptr %minDepth, align 4
  %add = add i32 %0, 2
  %1 = load i32, ptr %N, align 4
  %cmp = icmp ugt i32 %add, %1
  br i1 %cmp, label %if.then, label %if.else

if.then:                                          ; preds = %entry
  %2 = load i32, ptr %minDepth, align 4
  %add1 = add i32 %2, 2
  store i32 %add1, ptr %maxDepth, align 4
  br label %if.end

if.else:                                          ; preds = %entry
  %3 = load i32, ptr %N, align 4
  store i32 %3, ptr %maxDepth, align 4
  br label %if.end

if.end:                                           ; preds = %if.else, %if.then
  %4 = load i32, ptr %maxDepth, align 4
  %add2 = add i32 %4, 1
  store i32 %add2, ptr %stretchDepth, align 4
  %5 = load i32, ptr %stretchDepth, align 4
  %call = call ptr @BottomUpTree(i32 noundef %5)
  store ptr %call, ptr %stretchTree, align 8
  %6 = load i32, ptr %stretchDepth, align 4
  %7 = load ptr, ptr %stretchTree, align 8
  %call3 = call i64 @ItemCheck(ptr noundef %7)
  %call4 = call i32 (ptr, ...) @printf(ptr noundef @.str, i32 noundef %6, i64 noundef %call3)
  %8 = load ptr, ptr %stretchTree, align 8
  call void @DeleteTree(ptr noundef %8)
  %9 = load i32, ptr %maxDepth, align 4
  %call5 = call ptr @BottomUpTree(i32 noundef %9)
  store ptr %call5, ptr %longLivedTree, align 8
  %10 = load i32, ptr %minDepth, align 4
  store i32 %10, ptr %depth, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc22, %if.end
  %11 = load i32, ptr %depth, align 4
  %12 = load i32, ptr %maxDepth, align 4
  %cmp6 = icmp ule i32 %11, %12
  br i1 %cmp6, label %for.body, label %for.end24

for.body:                                         ; preds = %for.cond
  store i64 1, ptr %iterations, align 8
  store i32 0, ptr %i7, align 4
  br label %for.cond8

for.cond8:                                        ; preds = %for.inc, %for.body
  %13 = load i32, ptr %i7, align 4
  %14 = load i32, ptr %maxDepth, align 4
  %15 = load i32, ptr %depth, align 4
  %sub = sub i32 %14, %15
  %16 = load i32, ptr %minDepth, align 4
  %add9 = add i32 %sub, %16
  %cmp10 = icmp ult i32 %13, %add9
  br i1 %cmp10, label %for.body11, label %for.end

for.body11:                                       ; preds = %for.cond8
  %17 = load i64, ptr %iterations, align 8
  %mul = mul nsw i64 %17, 2
  store i64 %mul, ptr %iterations, align 8
  br label %for.inc

for.inc:                                          ; preds = %for.body11
  %18 = load i32, ptr %i7, align 4
  %inc = add nsw i32 %18, 1
  store i32 %inc, ptr %i7, align 4
  br label %for.cond8, !llvm.loop !6

for.end:                                          ; preds = %for.cond8
  store i64 0, ptr %check, align 8
  store i64 1, ptr %i, align 8
  br label %for.cond12

for.cond12:                                       ; preds = %for.inc18, %for.end
  %19 = load i64, ptr %i, align 8
  %20 = load i64, ptr %iterations, align 8
  %cmp13 = icmp sle i64 %19, %20
  br i1 %cmp13, label %for.body14, label %for.end20

for.body14:                                       ; preds = %for.cond12
  %21 = load i32, ptr %depth, align 4
  %call15 = call ptr @BottomUpTree(i32 noundef %21)
  store ptr %call15, ptr %tempTree, align 8
  %22 = load ptr, ptr %tempTree, align 8
  %call16 = call i64 @ItemCheck(ptr noundef %22)
  %23 = load i64, ptr %check, align 8
  %add17 = add nsw i64 %23, %call16
  store i64 %add17, ptr %check, align 8
  %24 = load ptr, ptr %tempTree, align 8
  call void @DeleteTree(ptr noundef %24)
  br label %for.inc18

for.inc18:                                        ; preds = %for.body14
  %25 = load i64, ptr %i, align 8
  %inc19 = add nsw i64 %25, 1
  store i64 %inc19, ptr %i, align 8
  br label %for.cond12, !llvm.loop !8

for.end20:                                        ; preds = %for.cond12
  %26 = load i64, ptr %iterations, align 8
  %27 = load i32, ptr %depth, align 4
  %28 = load i64, ptr %check, align 8
  %call21 = call i32 (ptr, ...) @printf(ptr noundef @.str.1, i64 noundef %26, i32 noundef %27, i64 noundef %28)
  br label %for.inc22

for.inc22:                                        ; preds = %for.end20
  %29 = load i32, ptr %depth, align 4
  %add23 = add i32 %29, 2
  store i32 %add23, ptr %depth, align 4
  br label %for.cond, !llvm.loop !9

for.end24:                                        ; preds = %for.cond
  %30 = load i32, ptr %maxDepth, align 4
  %31 = load ptr, ptr %longLivedTree, align 8
  %call25 = call i64 @ItemCheck(ptr noundef %31)
  %call26 = call i32 (ptr, ...) @printf(ptr noundef @.str.2, i32 noundef %30, i64 noundef %call25)
  ret i32 0
}

declare i32 @printf(ptr noundef, ...) #3

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { nounwind allocsize(0) "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #2 = { nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #3 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #4 = { nounwind allocsize(0) }
attributes #5 = { nounwind }

!llvm.module.flags = !{!0, !1, !2, !3, !4}
!llvm.ident = !{!5}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 8, !"PIC Level", i32 2}
!2 = !{i32 7, !"PIE Level", i32 2}
!3 = !{i32 7, !"uwtable", i32 2}
!4 = !{i32 7, !"frame-pointer", i32 2}
!5 = !{!"clang version 22.1.8 (https://github.com/conda-forge/clangdev-feedstock 015bdba1263c0b3ebb3c518ff5947fbd99692bd0)"}
!6 = distinct !{!6, !7}
!7 = !{!"llvm.loop.mustprogress"}
!8 = distinct !{!8, !7}
!9 = distinct !{!9, !7}
