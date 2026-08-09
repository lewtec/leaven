; ModuleID = 'testdata/ir/c_varargs_fnptr/source.c'
source_filename = "testdata/ir/c_varargs_fnptr/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

%struct.L = type { ptr }

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @init_log(ptr noundef %l) #0 {
entry:
  %l.addr = alloca ptr, align 8
  store ptr %l, ptr %l.addr, align 8
  %0 = load ptr, ptr %l.addr, align 8
  %log = getelementptr inbounds nuw %struct.L, ptr %0, i32 0, i32 0
  store ptr @logf, ptr %log, align 8
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define internal void @logf(ptr noundef %self, ptr noundef %fmt, ...) #0 {
entry:
  %self.addr = alloca ptr, align 8
  %fmt.addr = alloca ptr, align 8
  store ptr %self, ptr %self.addr, align 8
  store ptr %fmt, ptr %fmt.addr, align 8
  %0 = load ptr, ptr %self.addr, align 8
  %1 = load ptr, ptr %fmt.addr, align 8
  ret void
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
