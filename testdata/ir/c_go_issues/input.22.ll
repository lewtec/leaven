; ModuleID = 'testdata/ir/c_go_issues/source.c'
source_filename = "testdata/ir/c_go_issues/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

%struct.Outer = type { %struct.anon, ptr }
%struct.anon = type { i8, i8 }

@range = dso_local global i32 0, align 4

; Function Attrs: noinline nounwind optnone uwtable
define dso_local zeroext i1 @use_bool(i32 noundef %x) #0 {
entry:
  %x.addr = alloca i32, align 4
  %a = alloca i8, align 1
  %b = alloca i8, align 1
  store i32 %x, ptr %x.addr, align 4
  %0 = load i32, ptr %x.addr, align 4
  %cmp = icmp ne i32 %0, 0
  %storedv = zext i1 %cmp to i8
  store i8 %storedv, ptr %a, align 1
  %1 = load i8, ptr %a, align 1
  %loadedv = trunc i8 %1 to i1
  br i1 %loadedv, label %land.rhs, label %land.end

land.rhs:                                         ; preds = %entry
  %2 = load i32, ptr %x.addr, align 4
  %and = and i32 %2, 1
  %cmp1 = icmp ne i32 %and, 0
  br label %land.end

land.end:                                         ; preds = %land.rhs, %entry
  %3 = phi i1 [ false, %entry ], [ %cmp1, %land.rhs ]
  %storedv2 = zext i1 %3 to i8
  store i8 %storedv2, ptr %b, align 1
  %4 = load i8, ptr %b, align 1
  %loadedv3 = trunc i8 %4 to i1
  ret i1 %loadedv3
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @use_ptr(ptr noundef %o, i32 noundef %i) #0 {
entry:
  %o.addr = alloca ptr, align 8
  %i.addr = alloca i32, align 4
  store ptr %o, ptr %o.addr, align 8
  store i32 %i, ptr %i.addr, align 4
  %0 = load ptr, ptr %o.addr, align 8
  %ptr = getelementptr inbounds nuw %struct.Outer, ptr %0, i32 0, i32 1
  %1 = load ptr, ptr %ptr, align 8
  %2 = load i32, ptr %i.addr, align 4
  %idxprom = sext i32 %2 to i64
  %arrayidx = getelementptr inbounds i32, ptr %1, i64 %idxprom
  %3 = load i32, ptr %arrayidx, align 4
  %4 = load ptr, ptr %o.addr, align 8
  %nested = getelementptr inbounds nuw %struct.Outer, ptr %4, i32 0, i32 0
  %a = getelementptr inbounds nuw %struct.anon, ptr %nested, i32 0, i32 0
  %5 = load i8, ptr %a, align 8
  %conv = zext i8 %5 to i32
  %add = add nsw i32 %3, %conv
  ret i32 %add
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @use_range(i32 noundef %range) #0 {
entry:
  %range.addr = alloca i32, align 4
  store i32 %range, ptr %range.addr, align 4
  %0 = load i32, ptr %range.addr, align 4
  %add = add nsw i32 %0, 1
  ret i32 %add
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @set_range(i32 noundef %v) #0 {
entry:
  %v.addr = alloca i32, align 4
  store i32 %v, ptr %v.addr, align 4
  %0 = load i32, ptr %v.addr, align 4
  store i32 %0, ptr @range, align 4
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
