; ModuleID = 'testdata/ir/c_spectral_norm/source.c'
source_filename = "testdata/ir/c_spectral_norm/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

@.str = private unnamed_addr constant [7 x i8] c"%0.9f\0A\00", align 1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local double @eval_A(i32 noundef %i, i32 noundef %j) #0 {
entry:
  %i.addr = alloca i32, align 4
  %j.addr = alloca i32, align 4
  store i32 %i, ptr %i.addr, align 4
  store i32 %j, ptr %j.addr, align 4
  %0 = load i32, ptr %i.addr, align 4
  %1 = load i32, ptr %j.addr, align 4
  %add = add nsw i32 %0, %1
  %2 = load i32, ptr %i.addr, align 4
  %3 = load i32, ptr %j.addr, align 4
  %add1 = add nsw i32 %2, %3
  %add2 = add nsw i32 %add1, 1
  %mul = mul nsw i32 %add, %add2
  %div = sdiv i32 %mul, 2
  %4 = load i32, ptr %i.addr, align 4
  %add3 = add nsw i32 %div, %4
  %add4 = add nsw i32 %add3, 1
  %conv = sitofp i32 %add4 to double
  %div5 = fdiv double 1.000000e+00, %conv
  ret double %div5
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @eval_A_times_u(i32 noundef %N, ptr noundef %u, ptr noundef %Au) #0 {
entry:
  %N.addr = alloca i32, align 4
  %u.addr = alloca ptr, align 8
  %Au.addr = alloca ptr, align 8
  %i = alloca i32, align 4
  %j = alloca i32, align 4
  store i32 %N, ptr %N.addr, align 4
  store ptr %u, ptr %u.addr, align 8
  store ptr %Au, ptr %Au.addr, align 8
  store i32 0, ptr %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc8, %entry
  %0 = load i32, ptr %i, align 4
  %1 = load i32, ptr %N.addr, align 4
  %cmp = icmp slt i32 %0, %1
  br i1 %cmp, label %for.body, label %for.end10

for.body:                                         ; preds = %for.cond
  %2 = load ptr, ptr %Au.addr, align 8
  %3 = load i32, ptr %i, align 4
  %idxprom = sext i32 %3 to i64
  %arrayidx = getelementptr inbounds double, ptr %2, i64 %idxprom
  store double 0.000000e+00, ptr %arrayidx, align 8
  store i32 0, ptr %j, align 4
  br label %for.cond1

for.cond1:                                        ; preds = %for.inc, %for.body
  %4 = load i32, ptr %j, align 4
  %5 = load i32, ptr %N.addr, align 4
  %cmp2 = icmp slt i32 %4, %5
  br i1 %cmp2, label %for.body3, label %for.end

for.body3:                                        ; preds = %for.cond1
  %6 = load i32, ptr %i, align 4
  %7 = load i32, ptr %j, align 4
  %call = call double @eval_A(i32 noundef %6, i32 noundef %7)
  %8 = load ptr, ptr %u.addr, align 8
  %9 = load i32, ptr %j, align 4
  %idxprom4 = sext i32 %9 to i64
  %arrayidx5 = getelementptr inbounds double, ptr %8, i64 %idxprom4
  %10 = load double, ptr %arrayidx5, align 8
  %11 = load ptr, ptr %Au.addr, align 8
  %12 = load i32, ptr %i, align 4
  %idxprom6 = sext i32 %12 to i64
  %arrayidx7 = getelementptr inbounds double, ptr %11, i64 %idxprom6
  %13 = load double, ptr %arrayidx7, align 8
  %14 = call double @llvm.fmuladd.f64(double %call, double %10, double %13)
  store double %14, ptr %arrayidx7, align 8
  br label %for.inc

for.inc:                                          ; preds = %for.body3
  %15 = load i32, ptr %j, align 4
  %inc = add nsw i32 %15, 1
  store i32 %inc, ptr %j, align 4
  br label %for.cond1, !llvm.loop !6

for.end:                                          ; preds = %for.cond1
  br label %for.inc8

for.inc8:                                         ; preds = %for.end
  %16 = load i32, ptr %i, align 4
  %inc9 = add nsw i32 %16, 1
  store i32 %inc9, ptr %i, align 4
  br label %for.cond, !llvm.loop !8

for.end10:                                        ; preds = %for.cond
  ret void
}

; Function Attrs: nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none)
declare double @llvm.fmuladd.f64(double, double, double) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @eval_At_times_u(i32 noundef %N, ptr noundef %u, ptr noundef %Au) #0 {
entry:
  %N.addr = alloca i32, align 4
  %u.addr = alloca ptr, align 8
  %Au.addr = alloca ptr, align 8
  %i = alloca i32, align 4
  %j = alloca i32, align 4
  store i32 %N, ptr %N.addr, align 4
  store ptr %u, ptr %u.addr, align 8
  store ptr %Au, ptr %Au.addr, align 8
  store i32 0, ptr %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc8, %entry
  %0 = load i32, ptr %i, align 4
  %1 = load i32, ptr %N.addr, align 4
  %cmp = icmp slt i32 %0, %1
  br i1 %cmp, label %for.body, label %for.end10

for.body:                                         ; preds = %for.cond
  %2 = load ptr, ptr %Au.addr, align 8
  %3 = load i32, ptr %i, align 4
  %idxprom = sext i32 %3 to i64
  %arrayidx = getelementptr inbounds double, ptr %2, i64 %idxprom
  store double 0.000000e+00, ptr %arrayidx, align 8
  store i32 0, ptr %j, align 4
  br label %for.cond1

for.cond1:                                        ; preds = %for.inc, %for.body
  %4 = load i32, ptr %j, align 4
  %5 = load i32, ptr %N.addr, align 4
  %cmp2 = icmp slt i32 %4, %5
  br i1 %cmp2, label %for.body3, label %for.end

for.body3:                                        ; preds = %for.cond1
  %6 = load i32, ptr %j, align 4
  %7 = load i32, ptr %i, align 4
  %call = call double @eval_A(i32 noundef %6, i32 noundef %7)
  %8 = load ptr, ptr %u.addr, align 8
  %9 = load i32, ptr %j, align 4
  %idxprom4 = sext i32 %9 to i64
  %arrayidx5 = getelementptr inbounds double, ptr %8, i64 %idxprom4
  %10 = load double, ptr %arrayidx5, align 8
  %11 = load ptr, ptr %Au.addr, align 8
  %12 = load i32, ptr %i, align 4
  %idxprom6 = sext i32 %12 to i64
  %arrayidx7 = getelementptr inbounds double, ptr %11, i64 %idxprom6
  %13 = load double, ptr %arrayidx7, align 8
  %14 = call double @llvm.fmuladd.f64(double %call, double %10, double %13)
  store double %14, ptr %arrayidx7, align 8
  br label %for.inc

for.inc:                                          ; preds = %for.body3
  %15 = load i32, ptr %j, align 4
  %inc = add nsw i32 %15, 1
  store i32 %inc, ptr %j, align 4
  br label %for.cond1, !llvm.loop !9

for.end:                                          ; preds = %for.cond1
  br label %for.inc8

for.inc8:                                         ; preds = %for.end
  %16 = load i32, ptr %i, align 4
  %inc9 = add nsw i32 %16, 1
  store i32 %inc9, ptr %i, align 4
  br label %for.cond, !llvm.loop !10

for.end10:                                        ; preds = %for.cond
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @eval_AtA_times_u(i32 noundef %N, ptr noundef %u, ptr noundef %AtAu) #0 {
entry:
  %N.addr = alloca i32, align 4
  %u.addr = alloca ptr, align 8
  %AtAu.addr = alloca ptr, align 8
  %saved_stack = alloca ptr, align 8
  %__vla_expr0 = alloca i64, align 8
  store i32 %N, ptr %N.addr, align 4
  store ptr %u, ptr %u.addr, align 8
  store ptr %AtAu, ptr %AtAu.addr, align 8
  %0 = load i32, ptr %N.addr, align 4
  %1 = zext i32 %0 to i64
  %2 = call ptr @llvm.stacksave.p0()
  store ptr %2, ptr %saved_stack, align 8
  %vla = alloca double, i64 %1, align 16
  store i64 %1, ptr %__vla_expr0, align 8
  %3 = load i32, ptr %N.addr, align 4
  %4 = load ptr, ptr %u.addr, align 8
  call void @eval_A_times_u(i32 noundef %3, ptr noundef %4, ptr noundef %vla)
  %5 = load i32, ptr %N.addr, align 4
  %6 = load ptr, ptr %AtAu.addr, align 8
  call void @eval_At_times_u(i32 noundef %5, ptr noundef %vla, ptr noundef %6)
  %7 = load ptr, ptr %saved_stack, align 8
  call void @llvm.stackrestore.p0(ptr %7)
  ret void
}

; Function Attrs: nocallback nofree nosync nounwind willreturn
declare ptr @llvm.stacksave.p0() #2

; Function Attrs: nocallback nofree nosync nounwind willreturn
declare void @llvm.stackrestore.p0(ptr) #2

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %i = alloca i32, align 4
  %N = alloca i32, align 4
  %saved_stack = alloca ptr, align 8
  %vBv = alloca double, align 8
  %vv = alloca double, align 8
  store i32 0, ptr %retval, align 4
  store i32 2000, ptr %N, align 4
  %0 = call ptr @llvm.stacksave.p0()
  store ptr %0, ptr %saved_stack, align 8
  %vla = alloca double, i64 2000, align 16
  %vla1 = alloca double, i64 2000, align 16
  store i32 0, ptr %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc, %entry
  %1 = load i32, ptr %i, align 4
  %cmp = icmp slt i32 %1, 2000
  br i1 %cmp, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %2 = load i32, ptr %i, align 4
  %idxprom = sext i32 %2 to i64
  %arrayidx = getelementptr inbounds double, ptr %vla, i64 %idxprom
  store double 1.000000e+00, ptr %arrayidx, align 8
  br label %for.inc

for.inc:                                          ; preds = %for.body
  %3 = load i32, ptr %i, align 4
  %inc = add nsw i32 %3, 1
  store i32 %inc, ptr %i, align 4
  br label %for.cond, !llvm.loop !11

for.end:                                          ; preds = %for.cond
  store i32 0, ptr %i, align 4
  br label %for.cond2

for.cond2:                                        ; preds = %for.inc5, %for.end
  %4 = load i32, ptr %i, align 4
  %cmp3 = icmp slt i32 %4, 10
  br i1 %cmp3, label %for.body4, label %for.end7

for.body4:                                        ; preds = %for.cond2
  call void @eval_AtA_times_u(i32 noundef 2000, ptr noundef %vla, ptr noundef %vla1)
  call void @eval_AtA_times_u(i32 noundef 2000, ptr noundef %vla1, ptr noundef %vla)
  br label %for.inc5

for.inc5:                                         ; preds = %for.body4
  %5 = load i32, ptr %i, align 4
  %inc6 = add nsw i32 %5, 1
  store i32 %inc6, ptr %i, align 4
  br label %for.cond2, !llvm.loop !12

for.end7:                                         ; preds = %for.cond2
  store double 0.000000e+00, ptr %vv, align 8
  store double 0.000000e+00, ptr %vBv, align 8
  store i32 0, ptr %i, align 4
  br label %for.cond8

for.cond8:                                        ; preds = %for.inc19, %for.end7
  %6 = load i32, ptr %i, align 4
  %cmp9 = icmp slt i32 %6, 2000
  br i1 %cmp9, label %for.body10, label %for.end21

for.body10:                                       ; preds = %for.cond8
  %7 = load i32, ptr %i, align 4
  %idxprom11 = sext i32 %7 to i64
  %arrayidx12 = getelementptr inbounds double, ptr %vla, i64 %idxprom11
  %8 = load double, ptr %arrayidx12, align 8
  %9 = load i32, ptr %i, align 4
  %idxprom13 = sext i32 %9 to i64
  %arrayidx14 = getelementptr inbounds double, ptr %vla1, i64 %idxprom13
  %10 = load double, ptr %arrayidx14, align 8
  %11 = load double, ptr %vBv, align 8
  %12 = call double @llvm.fmuladd.f64(double %8, double %10, double %11)
  store double %12, ptr %vBv, align 8
  %13 = load i32, ptr %i, align 4
  %idxprom15 = sext i32 %13 to i64
  %arrayidx16 = getelementptr inbounds double, ptr %vla1, i64 %idxprom15
  %14 = load double, ptr %arrayidx16, align 8
  %15 = load i32, ptr %i, align 4
  %idxprom17 = sext i32 %15 to i64
  %arrayidx18 = getelementptr inbounds double, ptr %vla1, i64 %idxprom17
  %16 = load double, ptr %arrayidx18, align 8
  %17 = load double, ptr %vv, align 8
  %18 = call double @llvm.fmuladd.f64(double %14, double %16, double %17)
  store double %18, ptr %vv, align 8
  br label %for.inc19

for.inc19:                                        ; preds = %for.body10
  %19 = load i32, ptr %i, align 4
  %inc20 = add nsw i32 %19, 1
  store i32 %inc20, ptr %i, align 4
  br label %for.cond8, !llvm.loop !13

for.end21:                                        ; preds = %for.cond8
  %20 = load double, ptr %vBv, align 8
  %21 = load double, ptr %vv, align 8
  %div = fdiv double %20, %21
  %call = call double @sqrt(double noundef %div) #5
  %call22 = call i32 (ptr, ...) @printf(ptr noundef @.str, double noundef %call)
  store i32 0, ptr %retval, align 4
  %22 = load ptr, ptr %saved_stack, align 8
  call void @llvm.stackrestore.p0(ptr %22)
  %23 = load i32, ptr %retval, align 4
  ret i32 %23
}

declare i32 @printf(ptr noundef, ...) #3

; Function Attrs: nounwind
declare double @sqrt(double noundef) #4

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none) }
attributes #2 = { nocallback nofree nosync nounwind willreturn }
attributes #3 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #4 = { nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
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
!10 = distinct !{!10, !7}
!11 = distinct !{!11, !7}
!12 = distinct !{!12, !7}
!13 = distinct !{!13, !7}
