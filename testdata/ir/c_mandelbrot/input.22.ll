; ModuleID = 'testdata/ir/c_mandelbrot/source.c'
source_filename = "testdata/ir/c_mandelbrot/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

@.str = private unnamed_addr constant [10 x i8] c"P4\0A%d %d\0A\00", align 1
@stdout = external global ptr, align 8

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %w = alloca i32, align 4
  %h = alloca i32, align 4
  %bit_num = alloca i32, align 4
  %byte_acc = alloca i8, align 1
  %i = alloca i32, align 4
  %iter = alloca i32, align 4
  %x = alloca double, align 8
  %y = alloca double, align 8
  %limit = alloca double, align 8
  %Zr = alloca double, align 8
  %Zi = alloca double, align 8
  %Cr = alloca double, align 8
  %Ci = alloca double, align 8
  %Tr = alloca double, align 8
  %Ti = alloca double, align 8
  store i32 0, ptr %retval, align 4
  store i32 0, ptr %bit_num, align 4
  store i8 0, ptr %byte_acc, align 1
  store i32 50, ptr %iter, align 4
  store double 2.000000e+00, ptr %limit, align 8
  store i32 200, ptr %h, align 4
  store i32 200, ptr %w, align 4
  %0 = load i32, ptr %w, align 4
  %1 = load i32, ptr %h, align 4
  %call = call i32 (ptr, ...) @printf(ptr noundef @.str, i32 noundef %0, i32 noundef %1)
  store double 0.000000e+00, ptr %y, align 8
  br label %for.cond

for.cond:                                         ; preds = %for.inc55, %entry
  %2 = load double, ptr %y, align 8
  %3 = load i32, ptr %h, align 4
  %conv = sitofp i32 %3 to double
  %cmp = fcmp olt double %2, %conv
  br i1 %cmp, label %for.body, label %for.end57

for.body:                                         ; preds = %for.cond
  store double 0.000000e+00, ptr %x, align 8
  br label %for.cond2

for.cond2:                                        ; preds = %for.inc52, %for.body
  %4 = load double, ptr %x, align 8
  %5 = load i32, ptr %w, align 4
  %conv3 = sitofp i32 %5 to double
  %cmp4 = fcmp olt double %4, %conv3
  br i1 %cmp4, label %for.body6, label %for.end54

for.body6:                                        ; preds = %for.cond2
  store double 0.000000e+00, ptr %Ti, align 8
  store double 0.000000e+00, ptr %Tr, align 8
  store double 0.000000e+00, ptr %Zi, align 8
  store double 0.000000e+00, ptr %Zr, align 8
  %6 = load double, ptr %x, align 8
  %mul = fmul double 2.000000e+00, %6
  %7 = load i32, ptr %w, align 4
  %conv7 = sitofp i32 %7 to double
  %div = fdiv double %mul, %conv7
  %sub = fsub double %div, 1.500000e+00
  store double %sub, ptr %Cr, align 8
  %8 = load double, ptr %y, align 8
  %mul8 = fmul double 2.000000e+00, %8
  %9 = load i32, ptr %h, align 4
  %conv9 = sitofp i32 %9 to double
  %div10 = fdiv double %mul8, %conv9
  %sub11 = fsub double %div10, 1.000000e+00
  store double %sub11, ptr %Ci, align 8
  store i32 0, ptr %i, align 4
  br label %for.cond12

for.cond12:                                       ; preds = %for.inc, %for.body6
  %10 = load i32, ptr %i, align 4
  %11 = load i32, ptr %iter, align 4
  %cmp13 = icmp slt i32 %10, %11
  br i1 %cmp13, label %land.rhs, label %land.end

land.rhs:                                         ; preds = %for.cond12
  %12 = load double, ptr %Tr, align 8
  %13 = load double, ptr %Ti, align 8
  %add = fadd double %12, %13
  %14 = load double, ptr %limit, align 8
  %15 = load double, ptr %limit, align 8
  %mul15 = fmul double %14, %15
  %cmp16 = fcmp ole double %add, %mul15
  br label %land.end

land.end:                                         ; preds = %land.rhs, %for.cond12
  %16 = phi i1 [ false, %for.cond12 ], [ %cmp16, %land.rhs ]
  br i1 %16, label %for.body18, label %for.end

for.body18:                                       ; preds = %land.end
  %17 = load double, ptr %Zr, align 8
  %mul19 = fmul double 2.000000e+00, %17
  %18 = load double, ptr %Zi, align 8
  %19 = load double, ptr %Ci, align 8
  %20 = call double @llvm.fmuladd.f64(double %mul19, double %18, double %19)
  store double %20, ptr %Zi, align 8
  %21 = load double, ptr %Tr, align 8
  %22 = load double, ptr %Ti, align 8
  %sub21 = fsub double %21, %22
  %23 = load double, ptr %Cr, align 8
  %add22 = fadd double %sub21, %23
  store double %add22, ptr %Zr, align 8
  %24 = load double, ptr %Zr, align 8
  %25 = load double, ptr %Zr, align 8
  %mul23 = fmul double %24, %25
  store double %mul23, ptr %Tr, align 8
  %26 = load double, ptr %Zi, align 8
  %27 = load double, ptr %Zi, align 8
  %mul24 = fmul double %26, %27
  store double %mul24, ptr %Ti, align 8
  br label %for.inc

for.inc:                                          ; preds = %for.body18
  %28 = load i32, ptr %i, align 4
  %inc = add nsw i32 %28, 1
  store i32 %inc, ptr %i, align 4
  br label %for.cond12, !llvm.loop !6

for.end:                                          ; preds = %land.end
  %29 = load i8, ptr %byte_acc, align 1
  %conv25 = sext i8 %29 to i32
  %shl = shl i32 %conv25, 1
  %conv26 = trunc i32 %shl to i8
  store i8 %conv26, ptr %byte_acc, align 1
  %30 = load double, ptr %Tr, align 8
  %31 = load double, ptr %Ti, align 8
  %add27 = fadd double %30, %31
  %32 = load double, ptr %limit, align 8
  %33 = load double, ptr %limit, align 8
  %mul28 = fmul double %32, %33
  %cmp29 = fcmp ole double %add27, %mul28
  br i1 %cmp29, label %if.then, label %if.end

if.then:                                          ; preds = %for.end
  %34 = load i8, ptr %byte_acc, align 1
  %conv31 = sext i8 %34 to i32
  %or = or i32 %conv31, 1
  %conv32 = trunc i32 %or to i8
  store i8 %conv32, ptr %byte_acc, align 1
  br label %if.end

if.end:                                           ; preds = %if.then, %for.end
  %35 = load i32, ptr %bit_num, align 4
  %inc33 = add nsw i32 %35, 1
  store i32 %inc33, ptr %bit_num, align 4
  %36 = load i32, ptr %bit_num, align 4
  %cmp34 = icmp eq i32 %36, 8
  br i1 %cmp34, label %if.then36, label %if.else

if.then36:                                        ; preds = %if.end
  %37 = load i8, ptr %byte_acc, align 1
  %conv37 = sext i8 %37 to i32
  %38 = load ptr, ptr @stdout, align 8
  %call38 = call i32 @putc(i32 noundef %conv37, ptr noundef %38)
  store i8 0, ptr %byte_acc, align 1
  store i32 0, ptr %bit_num, align 4
  br label %if.end51

if.else:                                          ; preds = %if.end
  %39 = load double, ptr %x, align 8
  %40 = load i32, ptr %w, align 4
  %sub39 = sub nsw i32 %40, 1
  %conv40 = sitofp i32 %sub39 to double
  %cmp41 = fcmp oeq double %39, %conv40
  br i1 %cmp41, label %if.then43, label %if.end50

if.then43:                                        ; preds = %if.else
  %41 = load i32, ptr %w, align 4
  %rem = srem i32 %41, 8
  %sub44 = sub nsw i32 8, %rem
  %42 = load i8, ptr %byte_acc, align 1
  %conv45 = sext i8 %42 to i32
  %shl46 = shl i32 %conv45, %sub44
  %conv47 = trunc i32 %shl46 to i8
  store i8 %conv47, ptr %byte_acc, align 1
  %43 = load i8, ptr %byte_acc, align 1
  %conv48 = sext i8 %43 to i32
  %44 = load ptr, ptr @stdout, align 8
  %call49 = call i32 @putc(i32 noundef %conv48, ptr noundef %44)
  store i8 0, ptr %byte_acc, align 1
  store i32 0, ptr %bit_num, align 4
  br label %if.end50

if.end50:                                         ; preds = %if.then43, %if.else
  br label %if.end51

if.end51:                                         ; preds = %if.end50, %if.then36
  br label %for.inc52

for.inc52:                                        ; preds = %if.end51
  %45 = load double, ptr %x, align 8
  %inc53 = fadd double %45, 1.000000e+00
  store double %inc53, ptr %x, align 8
  br label %for.cond2, !llvm.loop !8

for.end54:                                        ; preds = %for.cond2
  br label %for.inc55

for.inc55:                                        ; preds = %for.end54
  %46 = load double, ptr %y, align 8
  %inc56 = fadd double %46, 1.000000e+00
  store double %inc56, ptr %y, align 8
  br label %for.cond, !llvm.loop !9

for.end57:                                        ; preds = %for.cond
  ret i32 0
}

declare i32 @printf(ptr noundef, ...) #1

; Function Attrs: nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none)
declare double @llvm.fmuladd.f64(double, double, double) #2

declare i32 @putc(i32 noundef, ptr noundef) #1

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #2 = { nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none) }

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
