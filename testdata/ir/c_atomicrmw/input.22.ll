; ModuleID = 'testdata/ir/c_atomicrmw/source.c'
source_filename = "testdata/ir/c_atomicrmw/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @atomic_inc(ptr noundef %p) #0 {
entry:
  %p.addr = alloca ptr, align 8
  %.atomictmp = alloca i32, align 4
  %atomic-temp = alloca i32, align 4
  store ptr %p, ptr %p.addr, align 8
  %0 = load ptr, ptr %p.addr, align 8
  store i32 1, ptr %.atomictmp, align 4
  %1 = load i32, ptr %.atomictmp, align 4
  %2 = atomicrmw add ptr %0, i32 %1 seq_cst, align 4
  store i32 %2, ptr %atomic-temp, align 4
  %3 = load i32, ptr %atomic-temp, align 4
  %add = add nsw i32 %3, 1
  ret i32 %add
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @atomic_dec(ptr noundef %p) #0 {
entry:
  %p.addr = alloca ptr, align 8
  %.atomictmp = alloca i32, align 4
  %atomic-temp = alloca i32, align 4
  store ptr %p, ptr %p.addr, align 8
  %0 = load ptr, ptr %p.addr, align 8
  store i32 1, ptr %.atomictmp, align 4
  %1 = load i32, ptr %.atomictmp, align 4
  %2 = atomicrmw sub ptr %0, i32 %1 seq_cst, align 4
  store i32 %2, ptr %atomic-temp, align 4
  %3 = load i32, ptr %atomic-temp, align 4
  %sub = sub nsw i32 %3, 1
  ret i32 %sub
}

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }

!llvm.module.flags = !{!0, !1, !2, !3, !4}
!llvm.ident = !{!5}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 8, !"PIC Level", i32 2}
!2 = !{i32 7, !"PIE Level", i32 2}
!3 = !{i32 7, !"uwtable", i32 2}
!4 = !{i32 7, !"frame-pointer", i32 2}
!5 = !{!"clang version 22.1.8 (https://github.com/conda-forge/clangdev-feedstock 015bdba1263c0b3ebb3c518ff5947fbd99692bd0)"}
