; ModuleID = 'testdata/ir/c_nbody/source.c'
source_filename = "testdata/ir/c_nbody/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

@x = dso_local global [5 x double] zeroinitializer, align 16
@y = dso_local global [5 x double] zeroinitializer, align 16
@z = dso_local global [5 x double] zeroinitializer, align 16
@mass = dso_local global [5 x double] zeroinitializer, align 16
@vx = dso_local global [5 x double] zeroinitializer, align 16
@vy = dso_local global [5 x double] zeroinitializer, align 16
@vz = dso_local global [5 x double] zeroinitializer, align 16
@.str = private unnamed_addr constant [6 x i8] c"%.9f\0A\00", align 1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @advance(i32 noundef %n) #0 {
entry:
  %n.addr = alloca i32, align 4
  %dx = alloca double, align 8
  %x1 = alloca double, align 8
  %y1 = alloca double, align 8
  %z1 = alloca double, align 8
  %dy = alloca double, align 8
  %dz = alloca double, align 8
  %R = alloca double, align 8
  %mag = alloca double, align 8
  %k = alloca i32, align 4
  %i = alloca i32, align 4
  %j = alloca i32, align 4
  %i64 = alloca i32, align 4
  store i32 %n, ptr %n.addr, align 4
  store i32 1, ptr %k, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc86, %entry
  %0 = load i32, ptr %k, align 4
  %1 = load i32, ptr %n.addr, align 4
  %cmp = icmp sle i32 %0, %1
  br i1 %cmp, label %for.body, label %for.end88

for.body:                                         ; preds = %for.cond
  store i32 0, ptr %i, align 4
  br label %for.cond1

for.cond1:                                        ; preds = %for.inc61, %for.body
  %2 = load i32, ptr %i, align 4
  %cmp2 = icmp slt i32 %2, 5
  br i1 %cmp2, label %for.body3, label %for.end63

for.body3:                                        ; preds = %for.cond1
  %3 = load i32, ptr %i, align 4
  %idxprom = sext i32 %3 to i64
  %arrayidx = getelementptr inbounds [5 x double], ptr @x, i64 0, i64 %idxprom
  %4 = load double, ptr %arrayidx, align 8
  store double %4, ptr %x1, align 8
  %5 = load i32, ptr %i, align 4
  %idxprom4 = sext i32 %5 to i64
  %arrayidx5 = getelementptr inbounds [5 x double], ptr @y, i64 0, i64 %idxprom4
  %6 = load double, ptr %arrayidx5, align 8
  store double %6, ptr %y1, align 8
  %7 = load i32, ptr %i, align 4
  %idxprom6 = sext i32 %7 to i64
  %arrayidx7 = getelementptr inbounds [5 x double], ptr @z, i64 0, i64 %idxprom6
  %8 = load double, ptr %arrayidx7, align 8
  store double %8, ptr %z1, align 8
  %9 = load i32, ptr %i, align 4
  %add = add nsw i32 %9, 1
  store i32 %add, ptr %j, align 4
  br label %for.cond8

for.cond8:                                        ; preds = %for.inc, %for.body3
  %10 = load i32, ptr %j, align 4
  %cmp9 = icmp slt i32 %10, 5
  br i1 %cmp9, label %for.body10, label %for.end

for.body10:                                       ; preds = %for.cond8
  %11 = load double, ptr %x1, align 8
  %12 = load i32, ptr %j, align 4
  %idxprom11 = sext i32 %12 to i64
  %arrayidx12 = getelementptr inbounds [5 x double], ptr @x, i64 0, i64 %idxprom11
  %13 = load double, ptr %arrayidx12, align 8
  %sub = fsub double %11, %13
  store double %sub, ptr %dx, align 8
  %14 = load double, ptr %dx, align 8
  %15 = load double, ptr %dx, align 8
  %mul = fmul double %14, %15
  store double %mul, ptr %R, align 8
  %16 = load double, ptr %y1, align 8
  %17 = load i32, ptr %j, align 4
  %idxprom13 = sext i32 %17 to i64
  %arrayidx14 = getelementptr inbounds [5 x double], ptr @y, i64 0, i64 %idxprom13
  %18 = load double, ptr %arrayidx14, align 8
  %sub15 = fsub double %16, %18
  store double %sub15, ptr %dy, align 8
  %19 = load double, ptr %dy, align 8
  %20 = load double, ptr %dy, align 8
  %21 = load double, ptr %R, align 8
  %22 = call double @llvm.fmuladd.f64(double %19, double %20, double %21)
  store double %22, ptr %R, align 8
  %23 = load double, ptr %z1, align 8
  %24 = load i32, ptr %j, align 4
  %idxprom17 = sext i32 %24 to i64
  %arrayidx18 = getelementptr inbounds [5 x double], ptr @z, i64 0, i64 %idxprom17
  %25 = load double, ptr %arrayidx18, align 8
  %sub19 = fsub double %23, %25
  store double %sub19, ptr %dz, align 8
  %26 = load double, ptr %dz, align 8
  %27 = load double, ptr %dz, align 8
  %28 = load double, ptr %R, align 8
  %29 = call double @llvm.fmuladd.f64(double %26, double %27, double %28)
  store double %29, ptr %R, align 8
  %30 = load double, ptr %R, align 8
  %call = call double @sqrt(double noundef %30) #4
  store double %call, ptr %R, align 8
  %31 = load double, ptr %R, align 8
  %32 = load double, ptr %R, align 8
  %mul21 = fmul double %31, %32
  %33 = load double, ptr %R, align 8
  %mul22 = fmul double %mul21, %33
  %div = fdiv double 1.000000e-02, %mul22
  store double %div, ptr %mag, align 8
  %34 = load double, ptr %dx, align 8
  %35 = load i32, ptr %j, align 4
  %idxprom23 = sext i32 %35 to i64
  %arrayidx24 = getelementptr inbounds [5 x double], ptr @mass, i64 0, i64 %idxprom23
  %36 = load double, ptr %arrayidx24, align 8
  %mul25 = fmul double %34, %36
  %37 = load double, ptr %mag, align 8
  %38 = load i32, ptr %i, align 4
  %idxprom27 = sext i32 %38 to i64
  %arrayidx28 = getelementptr inbounds [5 x double], ptr @vx, i64 0, i64 %idxprom27
  %39 = load double, ptr %arrayidx28, align 8
  %neg = fneg double %mul25
  %40 = call double @llvm.fmuladd.f64(double %neg, double %37, double %39)
  store double %40, ptr %arrayidx28, align 8
  %41 = load double, ptr %dy, align 8
  %42 = load i32, ptr %j, align 4
  %idxprom29 = sext i32 %42 to i64
  %arrayidx30 = getelementptr inbounds [5 x double], ptr @mass, i64 0, i64 %idxprom29
  %43 = load double, ptr %arrayidx30, align 8
  %mul31 = fmul double %41, %43
  %44 = load double, ptr %mag, align 8
  %45 = load i32, ptr %i, align 4
  %idxprom33 = sext i32 %45 to i64
  %arrayidx34 = getelementptr inbounds [5 x double], ptr @vy, i64 0, i64 %idxprom33
  %46 = load double, ptr %arrayidx34, align 8
  %neg35 = fneg double %mul31
  %47 = call double @llvm.fmuladd.f64(double %neg35, double %44, double %46)
  store double %47, ptr %arrayidx34, align 8
  %48 = load double, ptr %dz, align 8
  %49 = load i32, ptr %j, align 4
  %idxprom36 = sext i32 %49 to i64
  %arrayidx37 = getelementptr inbounds [5 x double], ptr @mass, i64 0, i64 %idxprom36
  %50 = load double, ptr %arrayidx37, align 8
  %mul38 = fmul double %48, %50
  %51 = load double, ptr %mag, align 8
  %52 = load i32, ptr %i, align 4
  %idxprom40 = sext i32 %52 to i64
  %arrayidx41 = getelementptr inbounds [5 x double], ptr @vz, i64 0, i64 %idxprom40
  %53 = load double, ptr %arrayidx41, align 8
  %neg42 = fneg double %mul38
  %54 = call double @llvm.fmuladd.f64(double %neg42, double %51, double %53)
  store double %54, ptr %arrayidx41, align 8
  %55 = load double, ptr %dx, align 8
  %56 = load i32, ptr %i, align 4
  %idxprom43 = sext i32 %56 to i64
  %arrayidx44 = getelementptr inbounds [5 x double], ptr @mass, i64 0, i64 %idxprom43
  %57 = load double, ptr %arrayidx44, align 8
  %mul45 = fmul double %55, %57
  %58 = load double, ptr %mag, align 8
  %59 = load i32, ptr %j, align 4
  %idxprom47 = sext i32 %59 to i64
  %arrayidx48 = getelementptr inbounds [5 x double], ptr @vx, i64 0, i64 %idxprom47
  %60 = load double, ptr %arrayidx48, align 8
  %61 = call double @llvm.fmuladd.f64(double %mul45, double %58, double %60)
  store double %61, ptr %arrayidx48, align 8
  %62 = load double, ptr %dy, align 8
  %63 = load i32, ptr %i, align 4
  %idxprom49 = sext i32 %63 to i64
  %arrayidx50 = getelementptr inbounds [5 x double], ptr @mass, i64 0, i64 %idxprom49
  %64 = load double, ptr %arrayidx50, align 8
  %mul51 = fmul double %62, %64
  %65 = load double, ptr %mag, align 8
  %66 = load i32, ptr %j, align 4
  %idxprom53 = sext i32 %66 to i64
  %arrayidx54 = getelementptr inbounds [5 x double], ptr @vy, i64 0, i64 %idxprom53
  %67 = load double, ptr %arrayidx54, align 8
  %68 = call double @llvm.fmuladd.f64(double %mul51, double %65, double %67)
  store double %68, ptr %arrayidx54, align 8
  %69 = load double, ptr %dz, align 8
  %70 = load i32, ptr %i, align 4
  %idxprom55 = sext i32 %70 to i64
  %arrayidx56 = getelementptr inbounds [5 x double], ptr @mass, i64 0, i64 %idxprom55
  %71 = load double, ptr %arrayidx56, align 8
  %mul57 = fmul double %69, %71
  %72 = load double, ptr %mag, align 8
  %73 = load i32, ptr %j, align 4
  %idxprom59 = sext i32 %73 to i64
  %arrayidx60 = getelementptr inbounds [5 x double], ptr @vz, i64 0, i64 %idxprom59
  %74 = load double, ptr %arrayidx60, align 8
  %75 = call double @llvm.fmuladd.f64(double %mul57, double %72, double %74)
  store double %75, ptr %arrayidx60, align 8
  br label %for.inc

for.inc:                                          ; preds = %for.body10
  %76 = load i32, ptr %j, align 4
  %inc = add nsw i32 %76, 1
  store i32 %inc, ptr %j, align 4
  br label %for.cond8, !llvm.loop !6

for.end:                                          ; preds = %for.cond8
  br label %for.inc61

for.inc61:                                        ; preds = %for.end
  %77 = load i32, ptr %i, align 4
  %inc62 = add nsw i32 %77, 1
  store i32 %inc62, ptr %i, align 4
  br label %for.cond1, !llvm.loop !8

for.end63:                                        ; preds = %for.cond1
  store i32 0, ptr %i64, align 4
  br label %for.cond65

for.cond65:                                       ; preds = %for.inc83, %for.end63
  %78 = load i32, ptr %i64, align 4
  %cmp66 = icmp slt i32 %78, 5
  br i1 %cmp66, label %for.body67, label %for.end85

for.body67:                                       ; preds = %for.cond65
  %79 = load i32, ptr %i64, align 4
  %idxprom68 = sext i32 %79 to i64
  %arrayidx69 = getelementptr inbounds [5 x double], ptr @vx, i64 0, i64 %idxprom68
  %80 = load double, ptr %arrayidx69, align 8
  %81 = load i32, ptr %i64, align 4
  %idxprom71 = sext i32 %81 to i64
  %arrayidx72 = getelementptr inbounds [5 x double], ptr @x, i64 0, i64 %idxprom71
  %82 = load double, ptr %arrayidx72, align 8
  %83 = call double @llvm.fmuladd.f64(double 1.000000e-02, double %80, double %82)
  store double %83, ptr %arrayidx72, align 8
  %84 = load i32, ptr %i64, align 4
  %idxprom73 = sext i32 %84 to i64
  %arrayidx74 = getelementptr inbounds [5 x double], ptr @vy, i64 0, i64 %idxprom73
  %85 = load double, ptr %arrayidx74, align 8
  %86 = load i32, ptr %i64, align 4
  %idxprom76 = sext i32 %86 to i64
  %arrayidx77 = getelementptr inbounds [5 x double], ptr @y, i64 0, i64 %idxprom76
  %87 = load double, ptr %arrayidx77, align 8
  %88 = call double @llvm.fmuladd.f64(double 1.000000e-02, double %85, double %87)
  store double %88, ptr %arrayidx77, align 8
  %89 = load i32, ptr %i64, align 4
  %idxprom78 = sext i32 %89 to i64
  %arrayidx79 = getelementptr inbounds [5 x double], ptr @vz, i64 0, i64 %idxprom78
  %90 = load double, ptr %arrayidx79, align 8
  %91 = load i32, ptr %i64, align 4
  %idxprom81 = sext i32 %91 to i64
  %arrayidx82 = getelementptr inbounds [5 x double], ptr @z, i64 0, i64 %idxprom81
  %92 = load double, ptr %arrayidx82, align 8
  %93 = call double @llvm.fmuladd.f64(double 1.000000e-02, double %90, double %92)
  store double %93, ptr %arrayidx82, align 8
  br label %for.inc83

for.inc83:                                        ; preds = %for.body67
  %94 = load i32, ptr %i64, align 4
  %inc84 = add nsw i32 %94, 1
  store i32 %inc84, ptr %i64, align 4
  br label %for.cond65, !llvm.loop !9

for.end85:                                        ; preds = %for.cond65
  br label %for.inc86

for.inc86:                                        ; preds = %for.end85
  %95 = load i32, ptr %k, align 4
  %inc87 = add nsw i32 %95, 1
  store i32 %inc87, ptr %k, align 4
  br label %for.cond, !llvm.loop !10

for.end88:                                        ; preds = %for.cond
  ret void
}

; Function Attrs: nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none)
declare double @llvm.fmuladd.f64(double, double, double) #1

; Function Attrs: nounwind
declare double @sqrt(double noundef) #2

; Function Attrs: noinline nounwind optnone uwtable
define dso_local double @energy() #0 {
entry:
  %e = alloca double, align 8
  %i = alloca i32, align 4
  %j = alloca i32, align 4
  %dx = alloca double, align 8
  %dy = alloca double, align 8
  %dz = alloca double, align 8
  %distance = alloca double, align 8
  store double 0.000000e+00, ptr %e, align 8
  store i32 0, ptr %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc43, %entry
  %0 = load i32, ptr %i, align 4
  %cmp = icmp slt i32 %0, 5
  br i1 %cmp, label %for.body, label %for.end45

for.body:                                         ; preds = %for.cond
  %1 = load i32, ptr %i, align 4
  %idxprom = sext i32 %1 to i64
  %arrayidx = getelementptr inbounds [5 x double], ptr @mass, i64 0, i64 %idxprom
  %2 = load double, ptr %arrayidx, align 8
  %mul = fmul double 5.000000e-01, %2
  %3 = load i32, ptr %i, align 4
  %idxprom1 = sext i32 %3 to i64
  %arrayidx2 = getelementptr inbounds [5 x double], ptr @vx, i64 0, i64 %idxprom1
  %4 = load double, ptr %arrayidx2, align 8
  %5 = load i32, ptr %i, align 4
  %idxprom3 = sext i32 %5 to i64
  %arrayidx4 = getelementptr inbounds [5 x double], ptr @vx, i64 0, i64 %idxprom3
  %6 = load double, ptr %arrayidx4, align 8
  %7 = load i32, ptr %i, align 4
  %idxprom6 = sext i32 %7 to i64
  %arrayidx7 = getelementptr inbounds [5 x double], ptr @vy, i64 0, i64 %idxprom6
  %8 = load double, ptr %arrayidx7, align 8
  %9 = load i32, ptr %i, align 4
  %idxprom8 = sext i32 %9 to i64
  %arrayidx9 = getelementptr inbounds [5 x double], ptr @vy, i64 0, i64 %idxprom8
  %10 = load double, ptr %arrayidx9, align 8
  %mul10 = fmul double %8, %10
  %11 = call double @llvm.fmuladd.f64(double %4, double %6, double %mul10)
  %12 = load i32, ptr %i, align 4
  %idxprom11 = sext i32 %12 to i64
  %arrayidx12 = getelementptr inbounds [5 x double], ptr @vz, i64 0, i64 %idxprom11
  %13 = load double, ptr %arrayidx12, align 8
  %14 = load i32, ptr %i, align 4
  %idxprom13 = sext i32 %14 to i64
  %arrayidx14 = getelementptr inbounds [5 x double], ptr @vz, i64 0, i64 %idxprom13
  %15 = load double, ptr %arrayidx14, align 8
  %16 = call double @llvm.fmuladd.f64(double %13, double %15, double %11)
  %17 = load double, ptr %e, align 8
  %18 = call double @llvm.fmuladd.f64(double %mul, double %16, double %17)
  store double %18, ptr %e, align 8
  %19 = load i32, ptr %i, align 4
  %add = add nsw i32 %19, 1
  store i32 %add, ptr %j, align 4
  br label %for.cond17

for.cond17:                                       ; preds = %for.inc, %for.body
  %20 = load i32, ptr %j, align 4
  %cmp18 = icmp slt i32 %20, 5
  br i1 %cmp18, label %for.body19, label %for.end

for.body19:                                       ; preds = %for.cond17
  %21 = load i32, ptr %i, align 4
  %idxprom20 = sext i32 %21 to i64
  %arrayidx21 = getelementptr inbounds [5 x double], ptr @x, i64 0, i64 %idxprom20
  %22 = load double, ptr %arrayidx21, align 8
  %23 = load i32, ptr %j, align 4
  %idxprom22 = sext i32 %23 to i64
  %arrayidx23 = getelementptr inbounds [5 x double], ptr @x, i64 0, i64 %idxprom22
  %24 = load double, ptr %arrayidx23, align 8
  %sub = fsub double %22, %24
  store double %sub, ptr %dx, align 8
  %25 = load i32, ptr %i, align 4
  %idxprom24 = sext i32 %25 to i64
  %arrayidx25 = getelementptr inbounds [5 x double], ptr @y, i64 0, i64 %idxprom24
  %26 = load double, ptr %arrayidx25, align 8
  %27 = load i32, ptr %j, align 4
  %idxprom26 = sext i32 %27 to i64
  %arrayidx27 = getelementptr inbounds [5 x double], ptr @y, i64 0, i64 %idxprom26
  %28 = load double, ptr %arrayidx27, align 8
  %sub28 = fsub double %26, %28
  store double %sub28, ptr %dy, align 8
  %29 = load i32, ptr %i, align 4
  %idxprom29 = sext i32 %29 to i64
  %arrayidx30 = getelementptr inbounds [5 x double], ptr @z, i64 0, i64 %idxprom29
  %30 = load double, ptr %arrayidx30, align 8
  %31 = load i32, ptr %j, align 4
  %idxprom31 = sext i32 %31 to i64
  %arrayidx32 = getelementptr inbounds [5 x double], ptr @z, i64 0, i64 %idxprom31
  %32 = load double, ptr %arrayidx32, align 8
  %sub33 = fsub double %30, %32
  store double %sub33, ptr %dz, align 8
  %33 = load double, ptr %dx, align 8
  %34 = load double, ptr %dx, align 8
  %35 = load double, ptr %dy, align 8
  %36 = load double, ptr %dy, align 8
  %mul35 = fmul double %35, %36
  %37 = call double @llvm.fmuladd.f64(double %33, double %34, double %mul35)
  %38 = load double, ptr %dz, align 8
  %39 = load double, ptr %dz, align 8
  %40 = call double @llvm.fmuladd.f64(double %38, double %39, double %37)
  %call = call double @sqrt(double noundef %40) #4
  store double %call, ptr %distance, align 8
  %41 = load i32, ptr %i, align 4
  %idxprom37 = sext i32 %41 to i64
  %arrayidx38 = getelementptr inbounds [5 x double], ptr @mass, i64 0, i64 %idxprom37
  %42 = load double, ptr %arrayidx38, align 8
  %43 = load i32, ptr %j, align 4
  %idxprom39 = sext i32 %43 to i64
  %arrayidx40 = getelementptr inbounds [5 x double], ptr @mass, i64 0, i64 %idxprom39
  %44 = load double, ptr %arrayidx40, align 8
  %mul41 = fmul double %42, %44
  %45 = load double, ptr %distance, align 8
  %div = fdiv double %mul41, %45
  %46 = load double, ptr %e, align 8
  %sub42 = fsub double %46, %div
  store double %sub42, ptr %e, align 8
  br label %for.inc

for.inc:                                          ; preds = %for.body19
  %47 = load i32, ptr %j, align 4
  %inc = add nsw i32 %47, 1
  store i32 %inc, ptr %j, align 4
  br label %for.cond17, !llvm.loop !11

for.end:                                          ; preds = %for.cond17
  br label %for.inc43

for.inc43:                                        ; preds = %for.end
  %48 = load i32, ptr %i, align 4
  %inc44 = add nsw i32 %48, 1
  store i32 %inc44, ptr %i, align 4
  br label %for.cond, !llvm.loop !12

for.end45:                                        ; preds = %for.cond
  %49 = load double, ptr %e, align 8
  ret double %49
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @offset_momentum() #0 {
entry:
  %px = alloca double, align 8
  %py = alloca double, align 8
  %pz = alloca double, align 8
  %i = alloca i32, align 4
  store double 0.000000e+00, ptr %px, align 8
  store double 0.000000e+00, ptr %py, align 8
  store double 0.000000e+00, ptr %pz, align 8
  store i32 0, ptr %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc, %entry
  %0 = load i32, ptr %i, align 4
  %cmp = icmp slt i32 %0, 5
  br i1 %cmp, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %1 = load i32, ptr %i, align 4
  %idxprom = sext i32 %1 to i64
  %arrayidx = getelementptr inbounds [5 x double], ptr @vx, i64 0, i64 %idxprom
  %2 = load double, ptr %arrayidx, align 8
  %3 = load i32, ptr %i, align 4
  %idxprom1 = sext i32 %3 to i64
  %arrayidx2 = getelementptr inbounds [5 x double], ptr @mass, i64 0, i64 %idxprom1
  %4 = load double, ptr %arrayidx2, align 8
  %5 = load double, ptr %px, align 8
  %6 = call double @llvm.fmuladd.f64(double %2, double %4, double %5)
  store double %6, ptr %px, align 8
  %7 = load i32, ptr %i, align 4
  %idxprom3 = sext i32 %7 to i64
  %arrayidx4 = getelementptr inbounds [5 x double], ptr @vy, i64 0, i64 %idxprom3
  %8 = load double, ptr %arrayidx4, align 8
  %9 = load i32, ptr %i, align 4
  %idxprom5 = sext i32 %9 to i64
  %arrayidx6 = getelementptr inbounds [5 x double], ptr @mass, i64 0, i64 %idxprom5
  %10 = load double, ptr %arrayidx6, align 8
  %11 = load double, ptr %py, align 8
  %12 = call double @llvm.fmuladd.f64(double %8, double %10, double %11)
  store double %12, ptr %py, align 8
  %13 = load i32, ptr %i, align 4
  %idxprom7 = sext i32 %13 to i64
  %arrayidx8 = getelementptr inbounds [5 x double], ptr @vz, i64 0, i64 %idxprom7
  %14 = load double, ptr %arrayidx8, align 8
  %15 = load i32, ptr %i, align 4
  %idxprom9 = sext i32 %15 to i64
  %arrayidx10 = getelementptr inbounds [5 x double], ptr @mass, i64 0, i64 %idxprom9
  %16 = load double, ptr %arrayidx10, align 8
  %17 = load double, ptr %pz, align 8
  %18 = call double @llvm.fmuladd.f64(double %14, double %16, double %17)
  store double %18, ptr %pz, align 8
  br label %for.inc

for.inc:                                          ; preds = %for.body
  %19 = load i32, ptr %i, align 4
  %inc = add nsw i32 %19, 1
  store i32 %inc, ptr %i, align 4
  br label %for.cond, !llvm.loop !13

for.end:                                          ; preds = %for.cond
  %20 = load double, ptr %px, align 8
  %fneg = fneg double %20
  %div = fdiv double %fneg, 0x4043BD3CC9BE45DE
  store double %div, ptr @vx, align 16
  %21 = load double, ptr %py, align 8
  %fneg11 = fneg double %21
  %div12 = fdiv double %fneg11, 0x4043BD3CC9BE45DE
  store double %div12, ptr @vy, align 16
  %22 = load double, ptr %pz, align 8
  %fneg13 = fneg double %22
  %div14 = fdiv double %fneg13, 0x4043BD3CC9BE45DE
  store double %div14, ptr @vz, align 16
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @init() #0 {
entry:
  store double 0.000000e+00, ptr @x, align 16
  store double 0.000000e+00, ptr @y, align 16
  store double 0.000000e+00, ptr @z, align 16
  store double 0.000000e+00, ptr @vx, align 16
  store double 0.000000e+00, ptr @vy, align 16
  store double 0.000000e+00, ptr @vz, align 16
  store double 0x4043BD3CC9BE45DE, ptr @mass, align 16
  store double 0x40135DA0343CD92C, ptr getelementptr inbounds ([5 x double], ptr @x, i64 0, i64 1), align 8
  store double 0xBFF290ABC01FDB7C, ptr getelementptr inbounds ([5 x double], ptr @y, i64 0, i64 1), align 8
  store double 0xBFBA86F96C25EBF0, ptr getelementptr inbounds ([5 x double], ptr @z, i64 0, i64 1), align 8
  store double 0x3FE367069B93CCBC, ptr getelementptr inbounds ([5 x double], ptr @vx, i64 0, i64 1), align 8
  store double 0x40067EF2F57D949B, ptr getelementptr inbounds ([5 x double], ptr @vy, i64 0, i64 1), align 8
  store double 0xBF99D2D79A5A0715, ptr getelementptr inbounds ([5 x double], ptr @vz, i64 0, i64 1), align 8
  store double 0x3FA34C95D9AB33D8, ptr getelementptr inbounds ([5 x double], ptr @mass, i64 0, i64 1), align 8
  store double 0x4020AFCDC332CA67, ptr getelementptr inbounds ([5 x double], ptr @x, i64 0, i64 2), align 16
  store double 0x40107FCB31DE01B0, ptr getelementptr inbounds ([5 x double], ptr @y, i64 0, i64 2), align 16
  store double 0xBFD9D353E1EB467C, ptr getelementptr inbounds ([5 x double], ptr @z, i64 0, i64 2), align 16
  store double 0xBFF02C21B8879442, ptr getelementptr inbounds ([5 x double], ptr @vx, i64 0, i64 2), align 16
  store double 0x3FFD35E9BF1F8F13, ptr getelementptr inbounds ([5 x double], ptr @vy, i64 0, i64 2), align 16
  store double 0x3F813C485F1123B4, ptr getelementptr inbounds ([5 x double], ptr @vz, i64 0, i64 2), align 16
  store double 0x3F871D490D07C637, ptr getelementptr inbounds ([5 x double], ptr @mass, i64 0, i64 2), align 16
  store double 0x4029C9EACEA7D9CF, ptr getelementptr inbounds ([5 x double], ptr @x, i64 0, i64 3), align 8
  store double 0xC02E38E8D626667E, ptr getelementptr inbounds ([5 x double], ptr @y, i64 0, i64 3), align 8
  store double 0xBFCC9557BE257DA0, ptr getelementptr inbounds ([5 x double], ptr @z, i64 0, i64 3), align 8
  store double 0x3FF1531CA9911BEF, ptr getelementptr inbounds ([5 x double], ptr @vx, i64 0, i64 3), align 8
  store double 0x3FEBCC7F3E54BBC5, ptr getelementptr inbounds ([5 x double], ptr @vy, i64 0, i64 3), align 8
  store double 0xBF862F6BFAF23E7C, ptr getelementptr inbounds ([5 x double], ptr @vz, i64 0, i64 3), align 8
  store double 0x3F5C3DD29CF41EB3, ptr getelementptr inbounds ([5 x double], ptr @mass, i64 0, i64 3), align 8
  store double 0x402EC267A905572A, ptr getelementptr inbounds ([5 x double], ptr @x, i64 0, i64 4), align 16
  store double 0xC039EB5833C8A220, ptr getelementptr inbounds ([5 x double], ptr @y, i64 0, i64 4), align 16
  store double 0x3FC6F1F393ABE540, ptr getelementptr inbounds ([5 x double], ptr @z, i64 0, i64 4), align 16
  store double 0x3FEF54B61659BC4A, ptr getelementptr inbounds ([5 x double], ptr @vx, i64 0, i64 4), align 16
  store double 0x3FE307C631C4FBA3, ptr getelementptr inbounds ([5 x double], ptr @vy, i64 0, i64 4), align 16
  store double 0xBFA1CB88587665F6, ptr getelementptr inbounds ([5 x double], ptr @vz, i64 0, i64 4), align 16
  store double 0x3F60A8F3531799AC, ptr getelementptr inbounds ([5 x double], ptr @mass, i64 0, i64 4), align 16
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %n = alloca i32, align 4
  store i32 0, ptr %retval, align 4
  store i32 5000000, ptr %n, align 4
  call void @init()
  call void @offset_momentum()
  %call = call double @energy()
  %call1 = call i32 (ptr, ...) @printf(ptr noundef @.str, double noundef %call)
  %0 = load i32, ptr %n, align 4
  call void @advance(i32 noundef %0)
  %call2 = call double @energy()
  %call3 = call i32 (ptr, ...) @printf(ptr noundef @.str, double noundef %call2)
  ret i32 0
}

declare i32 @printf(ptr noundef, ...) #3

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none) }
attributes #2 = { nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #3 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #4 = { nounwind }

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
