; ModuleID = '/tmp/free-fnptr-repro.c'
source_filename = "/tmp/free-fnptr-repro.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@current_free = dso_local global void (i8*)* @free, align 8

; Function Attrs: nounwind
declare dso_local void @free(i8* noundef) #0

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @release(i8* noundef %p) #1 {
entry:
  %p.addr = alloca i8*, align 8
  store i8* %p, i8** %p.addr, align 8
  %0 = load void (i8*)*, void (i8*)** @current_free, align 8
  %1 = load i8*, i8** %p.addr, align 8
  call void %0(i8* noundef %1)
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @set_free(void (i8*)* noundef %fn) #1 {
entry:
  %fn.addr = alloca void (i8*)*, align 8
  store void (i8*)* %fn, void (i8*)** %fn.addr, align 8
  %0 = load void (i8*)*, void (i8*)** %fn.addr, align 8
  %tobool = icmp ne void (i8*)* %0, null
  br i1 %tobool, label %cond.true, label %cond.false

cond.true:                                        ; preds = %entry
  %1 = load void (i8*)*, void (i8*)** %fn.addr, align 8
  br label %cond.end

cond.false:                                       ; preds = %entry
  br label %cond.end

cond.end:                                         ; preds = %cond.false, %cond.true
  %cond = phi void (i8*)* [ %1, %cond.true ], [ @free, %cond.false ]
  store void (i8*)* %cond, void (i8*)** @current_free, align 8
  ret void
}

attributes #0 = { nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }

!llvm.module.flags = !{!0, !1, !2}
!llvm.ident = !{!3}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 7, !"uwtable", i32 1}
!2 = !{i32 7, !"frame-pointer", i32 2}
!3 = !{!"clang version 14.0.6 (https://github.com/conda-forge/clangdev-feedstock ceeebe884c3cfd7160cf5a43e147f94439fafee3)"}
