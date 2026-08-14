; ModuleID = 'testdata/ir/c_fannkuch_redux/source.c'
source_filename = "testdata/ir/c_fannkuch_redux/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

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
  %x = alloca ptr, align 8
  %y = alloca ptr, align 8
  %c = alloca i32, align 4
  store ptr @t, ptr %x, align 8
  store ptr @s, ptr %y, align 8
  %0 = load i32, ptr @max_n, align 4
  store i32 %0, ptr %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.body, %entry
  %1 = load i32, ptr %i, align 4
  %dec = add nsw i32 %1, -1
  store i32 %dec, ptr %i, align 4
  %tobool = icmp ne i32 %1, 0
  br i1 %tobool, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %2 = load ptr, ptr %y, align 8
  %incdec.ptr = getelementptr inbounds nuw i32, ptr %2, i32 1
  store ptr %incdec.ptr, ptr %y, align 8
  %3 = load i32, ptr %2, align 4
  %4 = load ptr, ptr %x, align 8
  %incdec.ptr1 = getelementptr inbounds nuw i32, ptr %4, i32 1
  store ptr %incdec.ptr1, ptr %x, align 8
  store i32 %3, ptr %4, align 4
  br label %for.cond, !llvm.loop !6

for.end:                                          ; preds = %for.cond
  store i32 1, ptr %i, align 4
  br label %do.body

do.body:                                          ; preds = %do.cond, %for.end
  store ptr @t, ptr %x, align 8
  %5 = load i32, ptr @t, align 16
  %idx.ext = sext i32 %5 to i64
  %add.ptr = getelementptr inbounds i32, ptr @t, i64 %idx.ext
  store ptr %add.ptr, ptr %y, align 8
  br label %for.cond2

for.cond2:                                        ; preds = %for.body3, %do.body
  %6 = load ptr, ptr %x, align 8
  %7 = load ptr, ptr %y, align 8
  %cmp = icmp ult ptr %6, %7
  br i1 %cmp, label %for.body3, label %for.end6

for.body3:                                        ; preds = %for.cond2
  %8 = load ptr, ptr %x, align 8
  %9 = load i32, ptr %8, align 4
  store i32 %9, ptr %c, align 4
  %10 = load ptr, ptr %y, align 8
  %11 = load i32, ptr %10, align 4
  %12 = load ptr, ptr %x, align 8
  %incdec.ptr4 = getelementptr inbounds nuw i32, ptr %12, i32 1
  store ptr %incdec.ptr4, ptr %x, align 8
  store i32 %11, ptr %12, align 4
  %13 = load i32, ptr %c, align 4
  %14 = load ptr, ptr %y, align 8
  %incdec.ptr5 = getelementptr inbounds i32, ptr %14, i32 -1
  store ptr %incdec.ptr5, ptr %y, align 8
  store i32 %13, ptr %14, align 4
  br label %for.cond2, !llvm.loop !8

for.end6:                                         ; preds = %for.cond2
  %15 = load i32, ptr %i, align 4
  %inc = add nsw i32 %15, 1
  store i32 %inc, ptr %i, align 4
  br label %do.cond

do.cond:                                          ; preds = %for.end6
  %16 = load i32, ptr @t, align 16
  %idxprom = sext i32 %16 to i64
  %arrayidx = getelementptr inbounds [16 x i32], ptr @t, i64 0, i64 %idxprom
  %17 = load i32, ptr %arrayidx, align 4
  %tobool7 = icmp ne i32 %17, 0
  br i1 %tobool7, label %do.body, label %do.end, !llvm.loop !9

do.end:                                           ; preds = %do.cond
  %18 = load i32, ptr %i, align 4
  ret i32 %18
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @rotate(i32 noundef %n) #0 {
entry:
  %n.addr = alloca i32, align 4
  %c = alloca i32, align 4
  %i = alloca i32, align 4
  store i32 %n, ptr %n.addr, align 4
  %0 = load i32, ptr @s, align 16
  store i32 %0, ptr %c, align 4
  store i32 1, ptr %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc, %entry
  %1 = load i32, ptr %i, align 4
  %2 = load i32, ptr %n.addr, align 4
  %cmp = icmp sle i32 %1, %2
  br i1 %cmp, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %3 = load i32, ptr %i, align 4
  %idxprom = sext i32 %3 to i64
  %arrayidx = getelementptr inbounds [16 x i32], ptr @s, i64 0, i64 %idxprom
  %4 = load i32, ptr %arrayidx, align 4
  %5 = load i32, ptr %i, align 4
  %sub = sub nsw i32 %5, 1
  %idxprom1 = sext i32 %sub to i64
  %arrayidx2 = getelementptr inbounds [16 x i32], ptr @s, i64 0, i64 %idxprom1
  store i32 %4, ptr %arrayidx2, align 4
  br label %for.inc

for.inc:                                          ; preds = %for.body
  %6 = load i32, ptr %i, align 4
  %inc = add nsw i32 %6, 1
  store i32 %inc, ptr %i, align 4
  br label %for.cond, !llvm.loop !10

for.end:                                          ; preds = %for.cond
  %7 = load i32, ptr %c, align 4
  %8 = load i32, ptr %n.addr, align 4
  %idxprom3 = sext i32 %8 to i64
  %arrayidx4 = getelementptr inbounds [16 x i32], ptr @s, i64 0, i64 %idxprom3
  store i32 %7, ptr %arrayidx4, align 4
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @tk(i32 noundef %n) #0 {
entry:
  %n.addr = alloca i32, align 4
  %i = alloca i32, align 4
  %f = alloca i32, align 4
  %c = alloca [16 x i32], align 16
  store i32 %n, ptr %n.addr, align 4
  store i32 0, ptr %i, align 4
  call void @llvm.memset.p0.i64(ptr align 16 %c, i8 0, i64 64, i1 false)
  br label %while.cond

while.cond:                                       ; preds = %if.end19, %if.then, %entry
  %0 = load i32, ptr %i, align 4
  %1 = load i32, ptr %n.addr, align 4
  %cmp = icmp slt i32 %0, %1
  br i1 %cmp, label %while.body, label %while.end

while.body:                                       ; preds = %while.cond
  %2 = load i32, ptr %i, align 4
  call void @rotate(i32 noundef %2)
  %3 = load i32, ptr %i, align 4
  %idxprom = sext i32 %3 to i64
  %arrayidx = getelementptr inbounds [16 x i32], ptr %c, i64 0, i64 %idxprom
  %4 = load i32, ptr %arrayidx, align 4
  %5 = load i32, ptr %i, align 4
  %cmp1 = icmp sge i32 %4, %5
  br i1 %cmp1, label %if.then, label %if.end

if.then:                                          ; preds = %while.body
  %6 = load i32, ptr %i, align 4
  %inc = add nsw i32 %6, 1
  store i32 %inc, ptr %i, align 4
  %idxprom2 = sext i32 %6 to i64
  %arrayidx3 = getelementptr inbounds [16 x i32], ptr %c, i64 0, i64 %idxprom2
  store i32 0, ptr %arrayidx3, align 4
  br label %while.cond, !llvm.loop !11

if.end:                                           ; preds = %while.body
  %7 = load i32, ptr %i, align 4
  %idxprom4 = sext i32 %7 to i64
  %arrayidx5 = getelementptr inbounds [16 x i32], ptr %c, i64 0, i64 %idxprom4
  %8 = load i32, ptr %arrayidx5, align 4
  %inc6 = add nsw i32 %8, 1
  store i32 %inc6, ptr %arrayidx5, align 4
  store i32 1, ptr %i, align 4
  %9 = load i32, ptr @odd, align 4
  %not = xor i32 %9, -1
  store i32 %not, ptr @odd, align 4
  %10 = load i32, ptr @s, align 16
  %tobool = icmp ne i32 %10, 0
  br i1 %tobool, label %if.then7, label %if.end19

if.then7:                                         ; preds = %if.end
  %11 = load i32, ptr @s, align 16
  %idxprom8 = sext i32 %11 to i64
  %arrayidx9 = getelementptr inbounds [16 x i32], ptr @s, i64 0, i64 %idxprom8
  %12 = load i32, ptr %arrayidx9, align 4
  %tobool10 = icmp ne i32 %12, 0
  br i1 %tobool10, label %cond.true, label %cond.false

cond.true:                                        ; preds = %if.then7
  %call = call i32 @flip()
  br label %cond.end

cond.false:                                       ; preds = %if.then7
  br label %cond.end

cond.end:                                         ; preds = %cond.false, %cond.true
  %cond = phi i32 [ %call, %cond.true ], [ 1, %cond.false ]
  store i32 %cond, ptr %f, align 4
  %13 = load i32, ptr %f, align 4
  %14 = load i32, ptr @maxflips, align 4
  %cmp11 = icmp sgt i32 %13, %14
  br i1 %cmp11, label %if.then12, label %if.end13

if.then12:                                        ; preds = %cond.end
  %15 = load i32, ptr %f, align 4
  store i32 %15, ptr @maxflips, align 4
  br label %if.end13

if.end13:                                         ; preds = %if.then12, %cond.end
  %16 = load i32, ptr @odd, align 4
  %tobool14 = icmp ne i32 %16, 0
  br i1 %tobool14, label %cond.true15, label %cond.false16

cond.true15:                                      ; preds = %if.end13
  %17 = load i32, ptr %f, align 4
  %sub = sub nsw i32 0, %17
  br label %cond.end17

cond.false16:                                     ; preds = %if.end13
  %18 = load i32, ptr %f, align 4
  br label %cond.end17

cond.end17:                                       ; preds = %cond.false16, %cond.true15
  %cond18 = phi i32 [ %sub, %cond.true15 ], [ %18, %cond.false16 ]
  %19 = load i32, ptr @checksum, align 4
  %add = add nsw i32 %19, %cond18
  store i32 %add, ptr @checksum, align 4
  br label %if.end19

if.end19:                                         ; preds = %cond.end17, %if.end
  br label %while.cond, !llvm.loop !11

while.end:                                        ; preds = %while.cond
  ret void
}

; Function Attrs: nocallback nofree nounwind willreturn memory(argmem: write)
declare void @llvm.memset.p0.i64(ptr writeonly captures(none), i8, i64, i1 immarg) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %i = alloca i32, align 4
  store i32 0, ptr %retval, align 4
  store i32 10, ptr @max_n, align 4
  store i32 0, ptr %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc, %entry
  %0 = load i32, ptr %i, align 4
  %1 = load i32, ptr @max_n, align 4
  %cmp = icmp slt i32 %0, %1
  br i1 %cmp, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %2 = load i32, ptr %i, align 4
  %3 = load i32, ptr %i, align 4
  %idxprom = sext i32 %3 to i64
  %arrayidx = getelementptr inbounds [16 x i32], ptr @s, i64 0, i64 %idxprom
  store i32 %2, ptr %arrayidx, align 4
  br label %for.inc

for.inc:                                          ; preds = %for.body
  %4 = load i32, ptr %i, align 4
  %inc = add nsw i32 %4, 1
  store i32 %inc, ptr %i, align 4
  br label %for.cond, !llvm.loop !12

for.end:                                          ; preds = %for.cond
  %5 = load i32, ptr @max_n, align 4
  call void @tk(i32 noundef %5)
  %6 = load i32, ptr @checksum, align 4
  %7 = load i32, ptr @max_n, align 4
  %8 = load i32, ptr @maxflips, align 4
  %call = call i32 (ptr, ...) @printf(ptr noundef @.str, i32 noundef %6, i32 noundef %7, i32 noundef %8)
  ret i32 0
}

declare i32 @printf(ptr noundef, ...) #2

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { nocallback nofree nounwind willreturn memory(argmem: write) }
attributes #2 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }

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
!10 = distinct !{!10, !7}
!11 = distinct !{!11, !7}
!12 = distinct !{!12, !7}
