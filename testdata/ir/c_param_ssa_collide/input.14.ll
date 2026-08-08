; ModuleID = 'testdata/ir/param_ssa_collide/source.c'
source_filename = "testdata/ir/param_ssa_collide/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @renumber(i32* noundef %self, i32 noundef %v1, i32 noundef %v2) #0 {
entry:
  %self.addr = alloca i32*, align 8
  %v1.addr = alloca i32, align 4
  %v2.addr = alloca i32, align 4
  %a = alloca i32, align 4
  %b = alloca i32, align 4
  store i32* %self, i32** %self.addr, align 8
  store i32 %v1, i32* %v1.addr, align 4
  store i32 %v2, i32* %v2.addr, align 4
  %0 = load i32, i32* %v1.addr, align 4
  store i32 %0, i32* %a, align 4
  %1 = load i32, i32* %v2.addr, align 4
  store i32 %1, i32* %b, align 4
  %2 = load i32, i32* %a, align 4
  %3 = load i32, i32* %b, align 4
  %cmp = icmp eq i32 %2, %3
  br i1 %cmp, label %if.then, label %if.else

if.then:                                          ; preds = %entry
  %4 = load i32, i32* %a, align 4
  %5 = load i32*, i32** %self.addr, align 8
  store i32 %4, i32* %5, align 4
  br label %if.end

if.else:                                          ; preds = %entry
  %6 = load i32, i32* %b, align 4
  %7 = load i32, i32* %a, align 4
  %add = add nsw i32 %6, %7
  %8 = load i32*, i32** %self.addr, align 8
  store i32 %add, i32* %8, align 4
  br label %if.end

if.end:                                           ; preds = %if.else, %if.then
  ret void
}

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }

!llvm.module.flags = !{!0, !1, !2}
!llvm.ident = !{!3}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 7, !"uwtable", i32 1}
!2 = !{i32 7, !"frame-pointer", i32 2}
!3 = !{!"clang version 14.0.6 (https://github.com/conda-forge/clangdev-feedstock ceeebe884c3cfd7160cf5a43e147f94439fafee3)"}
