; ModuleID = 'testdata/ir/c_struct_state_shadow/source.c'
source_filename = "testdata/ir/c_struct_state_shadow/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

%struct.state = type { i32 }

@state_new.pool = internal global %struct.state zeroinitializer, align 4

; Function Attrs: noinline nounwind optnone uwtable
define dso_local ptr @create() #0 {
entry:
  %state = alloca ptr, align 8
  %call = call ptr @state_new()
  store ptr %call, ptr %state, align 8
  %0 = load ptr, ptr %state, align 8
  ret ptr %0
}

; Function Attrs: noinline nounwind optnone uwtable
define internal ptr @state_new() #0 {
entry:
  store i32 0, ptr @state_new.pool, align 4
  ret ptr @state_new.pool
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  store i32 0, ptr %retval, align 4
  %call = call ptr @create()
  %cmp = icmp ne ptr %call, null
  %0 = zext i1 %cmp to i64
  %cond = select i1 %cmp, i32 0, i32 1
  ret i32 %cond
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
