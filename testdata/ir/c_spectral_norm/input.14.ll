; ModuleID = 'testdata/ir/c_spectral_norm/source.c'
source_filename = "testdata/ir/c_spectral_norm/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@.str = private unnamed_addr constant [7 x i8] c"%0.9f\0A\00", align 1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local double @eval_A(i32 noundef %i, i32 noundef %j) #0 {
entry:
  %i.addr = alloca i32, align 4
  %j.addr = alloca i32, align 4
  store i32 %i, i32* %i.addr, align 4
  store i32 %j, i32* %j.addr, align 4
  %0 = load i32, i32* %i.addr, align 4
  %1 = load i32, i32* %j.addr, align 4
  %add = add nsw i32 %0, %1
  %2 = load i32, i32* %i.addr, align 4
  %3 = load i32, i32* %j.addr, align 4
  %add1 = add nsw i32 %2, %3
  %add2 = add nsw i32 %add1, 1
  %mul = mul nsw i32 %add, %add2
  %div = sdiv i32 %mul, 2
  %4 = load i32, i32* %i.addr, align 4
  %add3 = add nsw i32 %div, %4
  %add4 = add nsw i32 %add3, 1
  %conv = sitofp i32 %add4 to double
  %div5 = fdiv double 1.000000e+00, %conv
  ret double %div5
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @eval_A_times_u(i32 noundef %N, double* noundef %u, double* noundef %Au) #0 {
entry:
  %N.addr = alloca i32, align 4
  %u.addr = alloca double*, align 8
  %Au.addr = alloca double*, align 8
  %i = alloca i32, align 4
  %j = alloca i32, align 4
  store i32 %N, i32* %N.addr, align 4
  store double* %u, double** %u.addr, align 8
  store double* %Au, double** %Au.addr, align 8
  store i32 0, i32* %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc8, %entry
  %0 = load i32, i32* %i, align 4
  %1 = load i32, i32* %N.addr, align 4
  %cmp = icmp slt i32 %0, %1
  br i1 %cmp, label %for.body, label %for.end10

for.body:                                         ; preds = %for.cond
  %2 = load double*, double** %Au.addr, align 8
  %3 = load i32, i32* %i, align 4
  %idxprom = sext i32 %3 to i64
  %arrayidx = getelementptr inbounds double, double* %2, i64 %idxprom
  store double 0.000000e+00, double* %arrayidx, align 8
  store i32 0, i32* %j, align 4
  br label %for.cond1

for.cond1:                                        ; preds = %for.inc, %for.body
  %4 = load i32, i32* %j, align 4
  %5 = load i32, i32* %N.addr, align 4
  %cmp2 = icmp slt i32 %4, %5
  br i1 %cmp2, label %for.body3, label %for.end

for.body3:                                        ; preds = %for.cond1
  %6 = load i32, i32* %i, align 4
  %7 = load i32, i32* %j, align 4
  %call = call double @eval_A(i32 noundef %6, i32 noundef %7)
  %8 = load double*, double** %u.addr, align 8
  %9 = load i32, i32* %j, align 4
  %idxprom4 = sext i32 %9 to i64
  %arrayidx5 = getelementptr inbounds double, double* %8, i64 %idxprom4
  %10 = load double, double* %arrayidx5, align 8
  %11 = load double*, double** %Au.addr, align 8
  %12 = load i32, i32* %i, align 4
  %idxprom6 = sext i32 %12 to i64
  %arrayidx7 = getelementptr inbounds double, double* %11, i64 %idxprom6
  %13 = load double, double* %arrayidx7, align 8
  %14 = call double @llvm.fmuladd.f64(double %call, double %10, double %13)
  store double %14, double* %arrayidx7, align 8
  br label %for.inc

for.inc:                                          ; preds = %for.body3
  %15 = load i32, i32* %j, align 4
  %inc = add nsw i32 %15, 1
  store i32 %inc, i32* %j, align 4
  br label %for.cond1, !llvm.loop !4

for.end:                                          ; preds = %for.cond1
  br label %for.inc8

for.inc8:                                         ; preds = %for.end
  %16 = load i32, i32* %i, align 4
  %inc9 = add nsw i32 %16, 1
  store i32 %inc9, i32* %i, align 4
  br label %for.cond, !llvm.loop !6

for.end10:                                        ; preds = %for.cond
  ret void
}

; Function Attrs: nofree nosync nounwind readnone speculatable willreturn
declare double @llvm.fmuladd.f64(double, double, double) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @eval_At_times_u(i32 noundef %N, double* noundef %u, double* noundef %Au) #0 {
entry:
  %N.addr = alloca i32, align 4
  %u.addr = alloca double*, align 8
  %Au.addr = alloca double*, align 8
  %i = alloca i32, align 4
  %j = alloca i32, align 4
  store i32 %N, i32* %N.addr, align 4
  store double* %u, double** %u.addr, align 8
  store double* %Au, double** %Au.addr, align 8
  store i32 0, i32* %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc8, %entry
  %0 = load i32, i32* %i, align 4
  %1 = load i32, i32* %N.addr, align 4
  %cmp = icmp slt i32 %0, %1
  br i1 %cmp, label %for.body, label %for.end10

for.body:                                         ; preds = %for.cond
  %2 = load double*, double** %Au.addr, align 8
  %3 = load i32, i32* %i, align 4
  %idxprom = sext i32 %3 to i64
  %arrayidx = getelementptr inbounds double, double* %2, i64 %idxprom
  store double 0.000000e+00, double* %arrayidx, align 8
  store i32 0, i32* %j, align 4
  br label %for.cond1

for.cond1:                                        ; preds = %for.inc, %for.body
  %4 = load i32, i32* %j, align 4
  %5 = load i32, i32* %N.addr, align 4
  %cmp2 = icmp slt i32 %4, %5
  br i1 %cmp2, label %for.body3, label %for.end

for.body3:                                        ; preds = %for.cond1
  %6 = load i32, i32* %j, align 4
  %7 = load i32, i32* %i, align 4
  %call = call double @eval_A(i32 noundef %6, i32 noundef %7)
  %8 = load double*, double** %u.addr, align 8
  %9 = load i32, i32* %j, align 4
  %idxprom4 = sext i32 %9 to i64
  %arrayidx5 = getelementptr inbounds double, double* %8, i64 %idxprom4
  %10 = load double, double* %arrayidx5, align 8
  %11 = load double*, double** %Au.addr, align 8
  %12 = load i32, i32* %i, align 4
  %idxprom6 = sext i32 %12 to i64
  %arrayidx7 = getelementptr inbounds double, double* %11, i64 %idxprom6
  %13 = load double, double* %arrayidx7, align 8
  %14 = call double @llvm.fmuladd.f64(double %call, double %10, double %13)
  store double %14, double* %arrayidx7, align 8
  br label %for.inc

for.inc:                                          ; preds = %for.body3
  %15 = load i32, i32* %j, align 4
  %inc = add nsw i32 %15, 1
  store i32 %inc, i32* %j, align 4
  br label %for.cond1, !llvm.loop !7

for.end:                                          ; preds = %for.cond1
  br label %for.inc8

for.inc8:                                         ; preds = %for.end
  %16 = load i32, i32* %i, align 4
  %inc9 = add nsw i32 %16, 1
  store i32 %inc9, i32* %i, align 4
  br label %for.cond, !llvm.loop !8

for.end10:                                        ; preds = %for.cond
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @eval_AtA_times_u(i32 noundef %N, double* noundef %u, double* noundef %AtAu) #0 {
entry:
  %N.addr = alloca i32, align 4
  %u.addr = alloca double*, align 8
  %AtAu.addr = alloca double*, align 8
  %saved_stack = alloca i8*, align 8
  %__vla_expr0 = alloca i64, align 8
  store i32 %N, i32* %N.addr, align 4
  store double* %u, double** %u.addr, align 8
  store double* %AtAu, double** %AtAu.addr, align 8
  %0 = load i32, i32* %N.addr, align 4
  %1 = zext i32 %0 to i64
  %2 = call i8* @llvm.stacksave()
  store i8* %2, i8** %saved_stack, align 8
  %vla = alloca double, i64 %1, align 16
  store i64 %1, i64* %__vla_expr0, align 8
  %3 = load i32, i32* %N.addr, align 4
  %4 = load double*, double** %u.addr, align 8
  call void @eval_A_times_u(i32 noundef %3, double* noundef %4, double* noundef %vla)
  %5 = load i32, i32* %N.addr, align 4
  %6 = load double*, double** %AtAu.addr, align 8
  call void @eval_At_times_u(i32 noundef %5, double* noundef %vla, double* noundef %6)
  %7 = load i8*, i8** %saved_stack, align 8
  call void @llvm.stackrestore(i8* %7)
  ret void
}

; Function Attrs: nofree nosync nounwind willreturn
declare i8* @llvm.stacksave() #2

; Function Attrs: nofree nosync nounwind willreturn
declare void @llvm.stackrestore(i8*) #2

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %i = alloca i32, align 4
  %N = alloca i32, align 4
  %saved_stack = alloca i8*, align 8
  %vBv = alloca double, align 8
  %vv = alloca double, align 8
  store i32 0, i32* %retval, align 4
  store i32 2000, i32* %N, align 4
  %0 = call i8* @llvm.stacksave()
  store i8* %0, i8** %saved_stack, align 8
  %vla = alloca double, i64 2000, align 16
  %vla1 = alloca double, i64 2000, align 16
  store i32 0, i32* %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc, %entry
  %1 = load i32, i32* %i, align 4
  %cmp = icmp slt i32 %1, 2000
  br i1 %cmp, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %2 = load i32, i32* %i, align 4
  %idxprom = sext i32 %2 to i64
  %arrayidx = getelementptr inbounds double, double* %vla, i64 %idxprom
  store double 1.000000e+00, double* %arrayidx, align 8
  br label %for.inc

for.inc:                                          ; preds = %for.body
  %3 = load i32, i32* %i, align 4
  %inc = add nsw i32 %3, 1
  store i32 %inc, i32* %i, align 4
  br label %for.cond, !llvm.loop !9

for.end:                                          ; preds = %for.cond
  store i32 0, i32* %i, align 4
  br label %for.cond2

for.cond2:                                        ; preds = %for.inc5, %for.end
  %4 = load i32, i32* %i, align 4
  %cmp3 = icmp slt i32 %4, 10
  br i1 %cmp3, label %for.body4, label %for.end7

for.body4:                                        ; preds = %for.cond2
  call void @eval_AtA_times_u(i32 noundef 2000, double* noundef %vla, double* noundef %vla1)
  call void @eval_AtA_times_u(i32 noundef 2000, double* noundef %vla1, double* noundef %vla)
  br label %for.inc5

for.inc5:                                         ; preds = %for.body4
  %5 = load i32, i32* %i, align 4
  %inc6 = add nsw i32 %5, 1
  store i32 %inc6, i32* %i, align 4
  br label %for.cond2, !llvm.loop !10

for.end7:                                         ; preds = %for.cond2
  store double 0.000000e+00, double* %vv, align 8
  store double 0.000000e+00, double* %vBv, align 8
  store i32 0, i32* %i, align 4
  br label %for.cond8

for.cond8:                                        ; preds = %for.inc19, %for.end7
  %6 = load i32, i32* %i, align 4
  %cmp9 = icmp slt i32 %6, 2000
  br i1 %cmp9, label %for.body10, label %for.end21

for.body10:                                       ; preds = %for.cond8
  %7 = load i32, i32* %i, align 4
  %idxprom11 = sext i32 %7 to i64
  %arrayidx12 = getelementptr inbounds double, double* %vla, i64 %idxprom11
  %8 = load double, double* %arrayidx12, align 8
  %9 = load i32, i32* %i, align 4
  %idxprom13 = sext i32 %9 to i64
  %arrayidx14 = getelementptr inbounds double, double* %vla1, i64 %idxprom13
  %10 = load double, double* %arrayidx14, align 8
  %11 = load double, double* %vBv, align 8
  %12 = call double @llvm.fmuladd.f64(double %8, double %10, double %11)
  store double %12, double* %vBv, align 8
  %13 = load i32, i32* %i, align 4
  %idxprom15 = sext i32 %13 to i64
  %arrayidx16 = getelementptr inbounds double, double* %vla1, i64 %idxprom15
  %14 = load double, double* %arrayidx16, align 8
  %15 = load i32, i32* %i, align 4
  %idxprom17 = sext i32 %15 to i64
  %arrayidx18 = getelementptr inbounds double, double* %vla1, i64 %idxprom17
  %16 = load double, double* %arrayidx18, align 8
  %17 = load double, double* %vv, align 8
  %18 = call double @llvm.fmuladd.f64(double %14, double %16, double %17)
  store double %18, double* %vv, align 8
  br label %for.inc19

for.inc19:                                        ; preds = %for.body10
  %19 = load i32, i32* %i, align 4
  %inc20 = add nsw i32 %19, 1
  store i32 %inc20, i32* %i, align 4
  br label %for.cond8, !llvm.loop !11

for.end21:                                        ; preds = %for.cond8
  %20 = load double, double* %vBv, align 8
  %21 = load double, double* %vv, align 8
  %div = fdiv double %20, %21
  %call = call double @sqrt(double noundef %div) #5
  %call22 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([7 x i8], [7 x i8]* @.str, i64 0, i64 0), double noundef %call)
  store i32 0, i32* %retval, align 4
  %22 = load i8*, i8** %saved_stack, align 8
  call void @llvm.stackrestore(i8* %22)
  %23 = load i32, i32* %retval, align 4
  ret i32 %23
}

declare dso_local i32 @printf(i8* noundef, ...) #3

; Function Attrs: nounwind
declare dso_local double @sqrt(double noundef) #4

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { nofree nosync nounwind readnone speculatable willreturn }
attributes #2 = { nofree nosync nounwind willreturn }
attributes #3 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #4 = { nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #5 = { nounwind }

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
!11 = distinct !{!11, !5}
