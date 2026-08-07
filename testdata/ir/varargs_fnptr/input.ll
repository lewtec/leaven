; ModuleID = '/home/lucasew/.grok/worktrees/opensource-own-leaven/20260806-csmith-fuzz-fixes/testdata/ir/varargs_fnptr/source.c'
source_filename = "/home/lucasew/.grok/worktrees/opensource-own-leaven/20260806-csmith-fuzz-fixes/testdata/ir/varargs_fnptr/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

%struct.L = type { void (%struct.L*, i8*, ...)* }

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @init_log(%struct.L* noundef %l) #0 {
entry:
  %l.addr = alloca %struct.L*, align 8
  store %struct.L* %l, %struct.L** %l.addr, align 8
  %0 = load %struct.L*, %struct.L** %l.addr, align 8
  %log = getelementptr inbounds %struct.L, %struct.L* %0, i32 0, i32 0
  store void (%struct.L*, i8*, ...)* @logf, void (%struct.L*, i8*, ...)** %log, align 8
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define internal void @logf(%struct.L* noundef %self, i8* noundef %fmt, ...) #0 {
entry:
  %self.addr = alloca %struct.L*, align 8
  %fmt.addr = alloca i8*, align 8
  store %struct.L* %self, %struct.L** %self.addr, align 8
  store i8* %fmt, i8** %fmt.addr, align 8
  %0 = load %struct.L*, %struct.L** %self.addr, align 8
  %1 = load i8*, i8** %fmt.addr, align 8
  ret void
}

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }

!llvm.module.flags = !{!0, !1, !2}
!llvm.ident = !{!3}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 7, !"uwtable", i32 1}
!2 = !{i32 7, !"frame-pointer", i32 2}
!3 = !{!"clang version 14.0.6 (https://github.com/conda-forge/clangdev-feedstock ceeebe884c3cfd7160cf5a43e147f94439fafee3)"}
