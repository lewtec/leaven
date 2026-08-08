; ModuleID = 'testdata/ir/struct_state_shadow/source.c'
source_filename = "testdata/ir/struct_state_shadow/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

%struct.state = type { i32 }

@state_new.pool = internal global %struct.state zeroinitializer, align 4

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i8* @create() #0 {
entry:
  %state = alloca %struct.state*, align 8
  %call = call %struct.state* @state_new()
  store %struct.state* %call, %struct.state** %state, align 8
  %0 = load %struct.state*, %struct.state** %state, align 8
  %1 = bitcast %struct.state* %0 to i8*
  ret i8* %1
}

; Function Attrs: noinline nounwind optnone uwtable
define internal %struct.state* @state_new() #0 {
entry:
  store i32 0, i32* getelementptr inbounds (%struct.state, %struct.state* @state_new.pool, i32 0, i32 0), align 4
  ret %struct.state* @state_new.pool
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  store i32 0, i32* %retval, align 4
  %call = call i8* @create()
  %cmp = icmp ne i8* %call, null
  %0 = zext i1 %cmp to i64
  %cond = select i1 %cmp, i32 0, i32 1
  ret i32 %cond
}

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }

!llvm.module.flags = !{!0, !1, !2}
!llvm.ident = !{!3}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 7, !"uwtable", i32 1}
!2 = !{i32 7, !"frame-pointer", i32 2}
!3 = !{!"clang version 14.0.6 (https://github.com/conda-forge/clangdev-feedstock ceeebe884c3cfd7160cf5a43e147f94439fafee3)"}
