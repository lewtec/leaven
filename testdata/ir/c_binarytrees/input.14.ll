; ModuleID = 'testdata/ir/c_binarytrees/source.c'
source_filename = "testdata/ir/c_binarytrees/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

%struct.tn = type { %struct.tn*, %struct.tn* }

@.str = private unnamed_addr constant [38 x i8] c"stretch tree of depth %u\09 check: %li\0A\00", align 1
@.str.1 = private unnamed_addr constant [36 x i8] c"%li\09 trees of depth %u\09 check: %li\0A\00", align 1
@.str.2 = private unnamed_addr constant [41 x i8] c"long lived tree of depth %u\09 check: %li\0A\00", align 1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.tn* @NewTreeNode(%struct.tn* noundef %left, %struct.tn* noundef %right) #0 {
entry:
  %left.addr = alloca %struct.tn*, align 8
  %right.addr = alloca %struct.tn*, align 8
  %new = alloca %struct.tn*, align 8
  store %struct.tn* %left, %struct.tn** %left.addr, align 8
  store %struct.tn* %right, %struct.tn** %right.addr, align 8
  %call = call noalias i8* @malloc(i64 noundef 16) #3
  %0 = bitcast i8* %call to %struct.tn*
  store %struct.tn* %0, %struct.tn** %new, align 8
  %1 = load %struct.tn*, %struct.tn** %left.addr, align 8
  %2 = load %struct.tn*, %struct.tn** %new, align 8
  %left1 = getelementptr inbounds %struct.tn, %struct.tn* %2, i32 0, i32 0
  store %struct.tn* %1, %struct.tn** %left1, align 8
  %3 = load %struct.tn*, %struct.tn** %right.addr, align 8
  %4 = load %struct.tn*, %struct.tn** %new, align 8
  %right2 = getelementptr inbounds %struct.tn, %struct.tn* %4, i32 0, i32 1
  store %struct.tn* %3, %struct.tn** %right2, align 8
  %5 = load %struct.tn*, %struct.tn** %new, align 8
  ret %struct.tn* %5
}

; Function Attrs: nounwind
declare dso_local noalias i8* @malloc(i64 noundef) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i64 @ItemCheck(%struct.tn* noundef %tree) #0 {
entry:
  %retval = alloca i64, align 8
  %tree.addr = alloca %struct.tn*, align 8
  store %struct.tn* %tree, %struct.tn** %tree.addr, align 8
  %0 = load %struct.tn*, %struct.tn** %tree.addr, align 8
  %left = getelementptr inbounds %struct.tn, %struct.tn* %0, i32 0, i32 0
  %1 = load %struct.tn*, %struct.tn** %left, align 8
  %cmp = icmp eq %struct.tn* %1, null
  br i1 %cmp, label %if.then, label %if.else

if.then:                                          ; preds = %entry
  store i64 1, i64* %retval, align 8
  br label %return

if.else:                                          ; preds = %entry
  %2 = load %struct.tn*, %struct.tn** %tree.addr, align 8
  %left1 = getelementptr inbounds %struct.tn, %struct.tn* %2, i32 0, i32 0
  %3 = load %struct.tn*, %struct.tn** %left1, align 8
  %call = call i64 @ItemCheck(%struct.tn* noundef %3)
  %add = add nsw i64 1, %call
  %4 = load %struct.tn*, %struct.tn** %tree.addr, align 8
  %right = getelementptr inbounds %struct.tn, %struct.tn* %4, i32 0, i32 1
  %5 = load %struct.tn*, %struct.tn** %right, align 8
  %call2 = call i64 @ItemCheck(%struct.tn* noundef %5)
  %add3 = add nsw i64 %add, %call2
  store i64 %add3, i64* %retval, align 8
  br label %return

return:                                           ; preds = %if.else, %if.then
  %6 = load i64, i64* %retval, align 8
  ret i64 %6
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.tn* @BottomUpTree(i32 noundef %depth) #0 {
entry:
  %retval = alloca %struct.tn*, align 8
  %depth.addr = alloca i32, align 4
  store i32 %depth, i32* %depth.addr, align 4
  %0 = load i32, i32* %depth.addr, align 4
  %cmp = icmp ugt i32 %0, 0
  br i1 %cmp, label %if.then, label %if.else

if.then:                                          ; preds = %entry
  %1 = load i32, i32* %depth.addr, align 4
  %sub = sub i32 %1, 1
  %call = call %struct.tn* @BottomUpTree(i32 noundef %sub)
  %2 = load i32, i32* %depth.addr, align 4
  %sub1 = sub i32 %2, 1
  %call2 = call %struct.tn* @BottomUpTree(i32 noundef %sub1)
  %call3 = call %struct.tn* @NewTreeNode(%struct.tn* noundef %call, %struct.tn* noundef %call2)
  store %struct.tn* %call3, %struct.tn** %retval, align 8
  br label %return

if.else:                                          ; preds = %entry
  %call4 = call %struct.tn* @NewTreeNode(%struct.tn* noundef null, %struct.tn* noundef null)
  store %struct.tn* %call4, %struct.tn** %retval, align 8
  br label %return

return:                                           ; preds = %if.else, %if.then
  %3 = load %struct.tn*, %struct.tn** %retval, align 8
  ret %struct.tn* %3
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @DeleteTree(%struct.tn* noundef %tree) #0 {
entry:
  %tree.addr = alloca %struct.tn*, align 8
  store %struct.tn* %tree, %struct.tn** %tree.addr, align 8
  %0 = load %struct.tn*, %struct.tn** %tree.addr, align 8
  %left = getelementptr inbounds %struct.tn, %struct.tn* %0, i32 0, i32 0
  %1 = load %struct.tn*, %struct.tn** %left, align 8
  %cmp = icmp ne %struct.tn* %1, null
  br i1 %cmp, label %if.then, label %if.end

if.then:                                          ; preds = %entry
  %2 = load %struct.tn*, %struct.tn** %tree.addr, align 8
  %left1 = getelementptr inbounds %struct.tn, %struct.tn* %2, i32 0, i32 0
  %3 = load %struct.tn*, %struct.tn** %left1, align 8
  call void @DeleteTree(%struct.tn* noundef %3)
  %4 = load %struct.tn*, %struct.tn** %tree.addr, align 8
  %right = getelementptr inbounds %struct.tn, %struct.tn* %4, i32 0, i32 1
  %5 = load %struct.tn*, %struct.tn** %right, align 8
  call void @DeleteTree(%struct.tn* noundef %5)
  br label %if.end

if.end:                                           ; preds = %if.then, %entry
  %6 = load %struct.tn*, %struct.tn** %tree.addr, align 8
  %7 = bitcast %struct.tn* %6 to i8*
  call void @free(i8* noundef %7) #3
  ret void
}

; Function Attrs: nounwind
declare dso_local void @free(i8* noundef) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %N = alloca i32, align 4
  %depth = alloca i32, align 4
  %minDepth = alloca i32, align 4
  %maxDepth = alloca i32, align 4
  %stretchDepth = alloca i32, align 4
  %stretchTree = alloca %struct.tn*, align 8
  %longLivedTree = alloca %struct.tn*, align 8
  %tempTree = alloca %struct.tn*, align 8
  %i = alloca i64, align 8
  %iterations = alloca i64, align 8
  %check = alloca i64, align 8
  %i7 = alloca i32, align 4
  store i32 0, i32* %retval, align 4
  store i32 10, i32* %N, align 4
  store i32 4, i32* %minDepth, align 4
  %0 = load i32, i32* %minDepth, align 4
  %add = add i32 %0, 2
  %1 = load i32, i32* %N, align 4
  %cmp = icmp ugt i32 %add, %1
  br i1 %cmp, label %if.then, label %if.else

if.then:                                          ; preds = %entry
  %2 = load i32, i32* %minDepth, align 4
  %add1 = add i32 %2, 2
  store i32 %add1, i32* %maxDepth, align 4
  br label %if.end

if.else:                                          ; preds = %entry
  %3 = load i32, i32* %N, align 4
  store i32 %3, i32* %maxDepth, align 4
  br label %if.end

if.end:                                           ; preds = %if.else, %if.then
  %4 = load i32, i32* %maxDepth, align 4
  %add2 = add i32 %4, 1
  store i32 %add2, i32* %stretchDepth, align 4
  %5 = load i32, i32* %stretchDepth, align 4
  %call = call %struct.tn* @BottomUpTree(i32 noundef %5)
  store %struct.tn* %call, %struct.tn** %stretchTree, align 8
  %6 = load i32, i32* %stretchDepth, align 4
  %7 = load %struct.tn*, %struct.tn** %stretchTree, align 8
  %call3 = call i64 @ItemCheck(%struct.tn* noundef %7)
  %call4 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([38 x i8], [38 x i8]* @.str, i64 0, i64 0), i32 noundef %6, i64 noundef %call3)
  %8 = load %struct.tn*, %struct.tn** %stretchTree, align 8
  call void @DeleteTree(%struct.tn* noundef %8)
  %9 = load i32, i32* %maxDepth, align 4
  %call5 = call %struct.tn* @BottomUpTree(i32 noundef %9)
  store %struct.tn* %call5, %struct.tn** %longLivedTree, align 8
  %10 = load i32, i32* %minDepth, align 4
  store i32 %10, i32* %depth, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc22, %if.end
  %11 = load i32, i32* %depth, align 4
  %12 = load i32, i32* %maxDepth, align 4
  %cmp6 = icmp ule i32 %11, %12
  br i1 %cmp6, label %for.body, label %for.end24

for.body:                                         ; preds = %for.cond
  store i64 1, i64* %iterations, align 8
  store i32 0, i32* %i7, align 4
  br label %for.cond8

for.cond8:                                        ; preds = %for.inc, %for.body
  %13 = load i32, i32* %i7, align 4
  %14 = load i32, i32* %maxDepth, align 4
  %15 = load i32, i32* %depth, align 4
  %sub = sub i32 %14, %15
  %16 = load i32, i32* %minDepth, align 4
  %add9 = add i32 %sub, %16
  %cmp10 = icmp ult i32 %13, %add9
  br i1 %cmp10, label %for.body11, label %for.end

for.body11:                                       ; preds = %for.cond8
  %17 = load i64, i64* %iterations, align 8
  %mul = mul nsw i64 %17, 2
  store i64 %mul, i64* %iterations, align 8
  br label %for.inc

for.inc:                                          ; preds = %for.body11
  %18 = load i32, i32* %i7, align 4
  %inc = add nsw i32 %18, 1
  store i32 %inc, i32* %i7, align 4
  br label %for.cond8, !llvm.loop !4

for.end:                                          ; preds = %for.cond8
  store i64 0, i64* %check, align 8
  store i64 1, i64* %i, align 8
  br label %for.cond12

for.cond12:                                       ; preds = %for.inc18, %for.end
  %19 = load i64, i64* %i, align 8
  %20 = load i64, i64* %iterations, align 8
  %cmp13 = icmp sle i64 %19, %20
  br i1 %cmp13, label %for.body14, label %for.end20

for.body14:                                       ; preds = %for.cond12
  %21 = load i32, i32* %depth, align 4
  %call15 = call %struct.tn* @BottomUpTree(i32 noundef %21)
  store %struct.tn* %call15, %struct.tn** %tempTree, align 8
  %22 = load %struct.tn*, %struct.tn** %tempTree, align 8
  %call16 = call i64 @ItemCheck(%struct.tn* noundef %22)
  %23 = load i64, i64* %check, align 8
  %add17 = add nsw i64 %23, %call16
  store i64 %add17, i64* %check, align 8
  %24 = load %struct.tn*, %struct.tn** %tempTree, align 8
  call void @DeleteTree(%struct.tn* noundef %24)
  br label %for.inc18

for.inc18:                                        ; preds = %for.body14
  %25 = load i64, i64* %i, align 8
  %inc19 = add nsw i64 %25, 1
  store i64 %inc19, i64* %i, align 8
  br label %for.cond12, !llvm.loop !6

for.end20:                                        ; preds = %for.cond12
  %26 = load i64, i64* %iterations, align 8
  %27 = load i32, i32* %depth, align 4
  %28 = load i64, i64* %check, align 8
  %call21 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([36 x i8], [36 x i8]* @.str.1, i64 0, i64 0), i64 noundef %26, i32 noundef %27, i64 noundef %28)
  br label %for.inc22

for.inc22:                                        ; preds = %for.end20
  %29 = load i32, i32* %depth, align 4
  %add23 = add i32 %29, 2
  store i32 %add23, i32* %depth, align 4
  br label %for.cond, !llvm.loop !7

for.end24:                                        ; preds = %for.cond
  %30 = load i32, i32* %maxDepth, align 4
  %31 = load %struct.tn*, %struct.tn** %longLivedTree, align 8
  %call25 = call i64 @ItemCheck(%struct.tn* noundef %31)
  %call26 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([41 x i8], [41 x i8]* @.str.2, i64 0, i64 0), i32 noundef %30, i64 noundef %call25)
  ret i32 0
}

declare dso_local i32 @printf(i8* noundef, ...) #2

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #2 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #3 = { nounwind }

!llvm.module.flags = !{!0, !1, !2}
!llvm.ident = !{!3}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 7, !"uwtable", i32 1}
!2 = !{i32 7, !"frame-pointer", i32 2}
!3 = !{!"clang version 14.0.6 (https://github.com/conda-forge/clangdev-feedstock ceeebe884c3cfd7160cf5a43e147f94439fafee3)"}
!4 = distinct !{!4, !5}
!5 = !{!"llvm.loop.mustprogress"}
!6 = distinct !{!6, !5}
!7 = distinct !{!7, !5}
