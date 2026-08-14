; ModuleID = 'testdata/ir/c_fannkuch_redux/source.c'
source_filename = "testdata/ir/c_fannkuch_redux/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@maxflips = dso_local global i32 0, align 4
@odd = dso_local global i32 0, align 4
@checksum = dso_local global i32 0, align 4
@t = dso_local global [16 x i32] zeroinitializer, align 16
@s = dso_local global [16 x i32] zeroinitializer, align 16
@max_n = dso_local global i32 0, align 4
@.str = private unnamed_addr constant [25 x i8] c"%d\0APfannkuchen(%d) = %d\0A\00", align 1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @flip() #0 {
entry:
  %i = alloca i32, align 4
  %x = alloca i32*, align 8
  %y = alloca i32*, align 8
  %c = alloca i32, align 4
  store i32* getelementptr inbounds ([16 x i32], [16 x i32]* @t, i64 0, i64 0), i32** %x, align 8
  store i32* getelementptr inbounds ([16 x i32], [16 x i32]* @s, i64 0, i64 0), i32** %y, align 8
  %0 = load i32, i32* @max_n, align 4
  store i32 %0, i32* %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.body, %entry
  %1 = load i32, i32* %i, align 4
  %dec = add nsw i32 %1, -1
  store i32 %dec, i32* %i, align 4
  %tobool = icmp ne i32 %1, 0
  br i1 %tobool, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %2 = load i32*, i32** %y, align 8
  %incdec.ptr = getelementptr inbounds i32, i32* %2, i32 1
  store i32* %incdec.ptr, i32** %y, align 8
  %3 = load i32, i32* %2, align 4
  %4 = load i32*, i32** %x, align 8
  %incdec.ptr1 = getelementptr inbounds i32, i32* %4, i32 1
  store i32* %incdec.ptr1, i32** %x, align 8
  store i32 %3, i32* %4, align 4
  br label %for.cond, !llvm.loop !4

for.end:                                          ; preds = %for.cond
  store i32 1, i32* %i, align 4
  br label %do.body

do.body:                                          ; preds = %do.cond, %for.end
  store i32* getelementptr inbounds ([16 x i32], [16 x i32]* @t, i64 0, i64 0), i32** %x, align 8
  %5 = load i32, i32* getelementptr inbounds ([16 x i32], [16 x i32]* @t, i64 0, i64 0), align 16
  %idx.ext = sext i32 %5 to i64
  %add.ptr = getelementptr inbounds i32, i32* getelementptr inbounds ([16 x i32], [16 x i32]* @t, i64 0, i64 0), i64 %idx.ext
  store i32* %add.ptr, i32** %y, align 8
  br label %for.cond2

for.cond2:                                        ; preds = %for.body3, %do.body
  %6 = load i32*, i32** %x, align 8
  %7 = load i32*, i32** %y, align 8
  %cmp = icmp ult i32* %6, %7
  br i1 %cmp, label %for.body3, label %for.end6

for.body3:                                        ; preds = %for.cond2
  %8 = load i32*, i32** %x, align 8
  %9 = load i32, i32* %8, align 4
  store i32 %9, i32* %c, align 4
  %10 = load i32*, i32** %y, align 8
  %11 = load i32, i32* %10, align 4
  %12 = load i32*, i32** %x, align 8
  %incdec.ptr4 = getelementptr inbounds i32, i32* %12, i32 1
  store i32* %incdec.ptr4, i32** %x, align 8
  store i32 %11, i32* %12, align 4
  %13 = load i32, i32* %c, align 4
  %14 = load i32*, i32** %y, align 8
  %incdec.ptr5 = getelementptr inbounds i32, i32* %14, i32 -1
  store i32* %incdec.ptr5, i32** %y, align 8
  store i32 %13, i32* %14, align 4
  br label %for.cond2, !llvm.loop !6

for.end6:                                         ; preds = %for.cond2
  %15 = load i32, i32* %i, align 4
  %inc = add nsw i32 %15, 1
  store i32 %inc, i32* %i, align 4
  br label %do.cond

do.cond:                                          ; preds = %for.end6
  %16 = load i32, i32* getelementptr inbounds ([16 x i32], [16 x i32]* @t, i64 0, i64 0), align 16
  %idxprom = sext i32 %16 to i64
  %arrayidx = getelementptr inbounds [16 x i32], [16 x i32]* @t, i64 0, i64 %idxprom
  %17 = load i32, i32* %arrayidx, align 4
  %tobool7 = icmp ne i32 %17, 0
  br i1 %tobool7, label %do.body, label %do.end, !llvm.loop !7

do.end:                                           ; preds = %do.cond
  %18 = load i32, i32* %i, align 4
  ret i32 %18
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @rotate(i32 noundef %n) #0 {
entry:
  %n.addr = alloca i32, align 4
  %c = alloca i32, align 4
  %i = alloca i32, align 4
  store i32 %n, i32* %n.addr, align 4
  %0 = load i32, i32* getelementptr inbounds ([16 x i32], [16 x i32]* @s, i64 0, i64 0), align 16
  store i32 %0, i32* %c, align 4
  store i32 1, i32* %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc, %entry
  %1 = load i32, i32* %i, align 4
  %2 = load i32, i32* %n.addr, align 4
  %cmp = icmp sle i32 %1, %2
  br i1 %cmp, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %3 = load i32, i32* %i, align 4
  %idxprom = sext i32 %3 to i64
  %arrayidx = getelementptr inbounds [16 x i32], [16 x i32]* @s, i64 0, i64 %idxprom
  %4 = load i32, i32* %arrayidx, align 4
  %5 = load i32, i32* %i, align 4
  %sub = sub nsw i32 %5, 1
  %idxprom1 = sext i32 %sub to i64
  %arrayidx2 = getelementptr inbounds [16 x i32], [16 x i32]* @s, i64 0, i64 %idxprom1
  store i32 %4, i32* %arrayidx2, align 4
  br label %for.inc

for.inc:                                          ; preds = %for.body
  %6 = load i32, i32* %i, align 4
  %inc = add nsw i32 %6, 1
  store i32 %inc, i32* %i, align 4
  br label %for.cond, !llvm.loop !8

for.end:                                          ; preds = %for.cond
  %7 = load i32, i32* %c, align 4
  %8 = load i32, i32* %n.addr, align 4
  %idxprom3 = sext i32 %8 to i64
  %arrayidx4 = getelementptr inbounds [16 x i32], [16 x i32]* @s, i64 0, i64 %idxprom3
  store i32 %7, i32* %arrayidx4, align 4
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @tk(i32 noundef %n) #0 {
entry:
  %n.addr = alloca i32, align 4
  %i = alloca i32, align 4
  %f = alloca i32, align 4
  %c = alloca [16 x i32], align 16
  store i32 %n, i32* %n.addr, align 4
  store i32 0, i32* %i, align 4
  %0 = bitcast [16 x i32]* %c to i8*
  call void @llvm.memset.p0i8.i64(i8* align 16 %0, i8 0, i64 64, i1 false)
  br label %while.cond

while.cond:                                       ; preds = %if.end19, %if.then, %entry
  %1 = load i32, i32* %i, align 4
  %2 = load i32, i32* %n.addr, align 4
  %cmp = icmp slt i32 %1, %2
  br i1 %cmp, label %while.body, label %while.end

while.body:                                       ; preds = %while.cond
  %3 = load i32, i32* %i, align 4
  call void @rotate(i32 noundef %3)
  %4 = load i32, i32* %i, align 4
  %idxprom = sext i32 %4 to i64
  %arrayidx = getelementptr inbounds [16 x i32], [16 x i32]* %c, i64 0, i64 %idxprom
  %5 = load i32, i32* %arrayidx, align 4
  %6 = load i32, i32* %i, align 4
  %cmp1 = icmp sge i32 %5, %6
  br i1 %cmp1, label %if.then, label %if.end

if.then:                                          ; preds = %while.body
  %7 = load i32, i32* %i, align 4
  %inc = add nsw i32 %7, 1
  store i32 %inc, i32* %i, align 4
  %idxprom2 = sext i32 %7 to i64
  %arrayidx3 = getelementptr inbounds [16 x i32], [16 x i32]* %c, i64 0, i64 %idxprom2
  store i32 0, i32* %arrayidx3, align 4
  br label %while.cond, !llvm.loop !9

if.end:                                           ; preds = %while.body
  %8 = load i32, i32* %i, align 4
  %idxprom4 = sext i32 %8 to i64
  %arrayidx5 = getelementptr inbounds [16 x i32], [16 x i32]* %c, i64 0, i64 %idxprom4
  %9 = load i32, i32* %arrayidx5, align 4
  %inc6 = add nsw i32 %9, 1
  store i32 %inc6, i32* %arrayidx5, align 4
  store i32 1, i32* %i, align 4
  %10 = load i32, i32* @odd, align 4
  %neg = xor i32 %10, -1
  store i32 %neg, i32* @odd, align 4
  %11 = load i32, i32* getelementptr inbounds ([16 x i32], [16 x i32]* @s, i64 0, i64 0), align 16
  %tobool = icmp ne i32 %11, 0
  br i1 %tobool, label %if.then7, label %if.end19

if.then7:                                         ; preds = %if.end
  %12 = load i32, i32* getelementptr inbounds ([16 x i32], [16 x i32]* @s, i64 0, i64 0), align 16
  %idxprom8 = sext i32 %12 to i64
  %arrayidx9 = getelementptr inbounds [16 x i32], [16 x i32]* @s, i64 0, i64 %idxprom8
  %13 = load i32, i32* %arrayidx9, align 4
  %tobool10 = icmp ne i32 %13, 0
  br i1 %tobool10, label %cond.true, label %cond.false

cond.true:                                        ; preds = %if.then7
  %call = call i32 @flip()
  br label %cond.end

cond.false:                                       ; preds = %if.then7
  br label %cond.end

cond.end:                                         ; preds = %cond.false, %cond.true
  %cond = phi i32 [ %call, %cond.true ], [ 1, %cond.false ]
  store i32 %cond, i32* %f, align 4
  %14 = load i32, i32* %f, align 4
  %15 = load i32, i32* @maxflips, align 4
  %cmp11 = icmp sgt i32 %14, %15
  br i1 %cmp11, label %if.then12, label %if.end13

if.then12:                                        ; preds = %cond.end
  %16 = load i32, i32* %f, align 4
  store i32 %16, i32* @maxflips, align 4
  br label %if.end13

if.end13:                                         ; preds = %if.then12, %cond.end
  %17 = load i32, i32* @odd, align 4
  %tobool14 = icmp ne i32 %17, 0
  br i1 %tobool14, label %cond.true15, label %cond.false16

cond.true15:                                      ; preds = %if.end13
  %18 = load i32, i32* %f, align 4
  %sub = sub nsw i32 0, %18
  br label %cond.end17

cond.false16:                                     ; preds = %if.end13
  %19 = load i32, i32* %f, align 4
  br label %cond.end17

cond.end17:                                       ; preds = %cond.false16, %cond.true15
  %cond18 = phi i32 [ %sub, %cond.true15 ], [ %19, %cond.false16 ]
  %20 = load i32, i32* @checksum, align 4
  %add = add nsw i32 %20, %cond18
  store i32 %add, i32* @checksum, align 4
  br label %if.end19

if.end19:                                         ; preds = %cond.end17, %if.end
  br label %while.cond, !llvm.loop !9

while.end:                                        ; preds = %while.cond
  ret void
}

; Function Attrs: argmemonly nofree nounwind willreturn writeonly
declare void @llvm.memset.p0i8.i64(i8* nocapture writeonly, i8, i64, i1 immarg) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %i = alloca i32, align 4
  store i32 0, i32* %retval, align 4
  store i32 10, i32* @max_n, align 4
  store i32 0, i32* %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc, %entry
  %0 = load i32, i32* %i, align 4
  %1 = load i32, i32* @max_n, align 4
  %cmp = icmp slt i32 %0, %1
  br i1 %cmp, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %2 = load i32, i32* %i, align 4
  %3 = load i32, i32* %i, align 4
  %idxprom = sext i32 %3 to i64
  %arrayidx = getelementptr inbounds [16 x i32], [16 x i32]* @s, i64 0, i64 %idxprom
  store i32 %2, i32* %arrayidx, align 4
  br label %for.inc

for.inc:                                          ; preds = %for.body
  %4 = load i32, i32* %i, align 4
  %inc = add nsw i32 %4, 1
  store i32 %inc, i32* %i, align 4
  br label %for.cond, !llvm.loop !10

for.end:                                          ; preds = %for.cond
  %5 = load i32, i32* @max_n, align 4
  call void @tk(i32 noundef %5)
  %6 = load i32, i32* @checksum, align 4
  %7 = load i32, i32* @max_n, align 4
  %8 = load i32, i32* @maxflips, align 4
  %call = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([25 x i8], [25 x i8]* @.str, i64 0, i64 0), i32 noundef %6, i32 noundef %7, i32 noundef %8)
  ret i32 0
}

declare dso_local i32 @printf(i8* noundef, ...) #2

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { argmemonly nofree nounwind willreturn writeonly }
attributes #2 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }

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
!8 = distinct !{!8, !5}
!9 = distinct !{!9, !5}
!10 = distinct !{!10, !5}
