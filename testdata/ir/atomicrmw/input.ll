; ModuleID = '/tmp/atomicrmw-repro.c'
source_filename = "/tmp/atomicrmw-repro.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @atomic_inc(i32* noundef %p) #0 {
entry:
  %p.addr = alloca i32*, align 8
  %.atomictmp = alloca i32, align 4
  %atomic-temp = alloca i32, align 4
  store i32* %p, i32** %p.addr, align 8
  %0 = load i32*, i32** %p.addr, align 8
  store i32 1, i32* %.atomictmp, align 4
  %1 = load i32, i32* %.atomictmp, align 4
  %2 = atomicrmw add i32* %0, i32 %1 seq_cst, align 4
  store i32 %2, i32* %atomic-temp, align 4
  %3 = load i32, i32* %atomic-temp, align 4
  %add = add nsw i32 %3, 1
  ret i32 %add
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @atomic_dec(i32* noundef %p) #0 {
entry:
  %p.addr = alloca i32*, align 8
  %.atomictmp = alloca i32, align 4
  %atomic-temp = alloca i32, align 4
  store i32* %p, i32** %p.addr, align 8
  %0 = load i32*, i32** %p.addr, align 8
  store i32 1, i32* %.atomictmp, align 4
  %1 = load i32, i32* %.atomictmp, align 4
  %2 = atomicrmw sub i32* %0, i32 %1 seq_cst, align 4
  store i32 %2, i32* %atomic-temp, align 4
  %3 = load i32, i32* %atomic-temp, align 4
  %sub = sub nsw i32 %3, 1
  ret i32 %sub
}

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }

!llvm.module.flags = !{!0, !1, !2}
!llvm.ident = !{!3}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 7, !"uwtable", i32 1}
!2 = !{i32 7, !"frame-pointer", i32 2}
!3 = !{!"clang version 14.0.6 (https://github.com/conda-forge/clangdev-feedstock ceeebe884c3cfd7160cf5a43e147f94439fafee3)"}
