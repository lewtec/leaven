; ModuleID = '/home/lucasew/.grok/worktrees/opensource-own-leaven/20260809-llir-v14/testdata/ir/cpp_point/source.cpp'
source_filename = "/home/lucasew/.grok/worktrees/opensource-own-leaven/20260809-llir-v14/testdata/ir/cpp_point/source.cpp"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

%struct.Point = type { i32, i32 }

$_ZN5PointC2Eii = comdat any

$_ZNK5Point9manhattanEv = comdat any

@.str = private unnamed_addr constant [4 x i8] c"%d\0A\00", align 1

; Function Attrs: mustprogress noinline norecurse optnone uwtable
define dso_local noundef i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %p = alloca %struct.Point, align 4
  store i32 0, i32* %retval, align 4
  call void @_ZN5PointC2Eii(%struct.Point* noundef nonnull align 4 dereferenceable(8) %p, i32 noundef 3, i32 noundef -4)
  %call = call noundef i32 @_ZNK5Point9manhattanEv(%struct.Point* noundef nonnull align 4 dereferenceable(8) %p)
  %call1 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([4 x i8], [4 x i8]* @.str, i64 0, i64 0), i32 noundef %call)
  ret i32 0
}

; Function Attrs: noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZN5PointC2Eii(%struct.Point* noundef nonnull align 4 dereferenceable(8) %this, i32 noundef %x, i32 noundef %y) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca %struct.Point*, align 8
  %x.addr = alloca i32, align 4
  %y.addr = alloca i32, align 4
  store %struct.Point* %this, %struct.Point** %this.addr, align 8
  store i32 %x, i32* %x.addr, align 4
  store i32 %y, i32* %y.addr, align 4
  %this1 = load %struct.Point*, %struct.Point** %this.addr, align 8
  %x2 = getelementptr inbounds %struct.Point, %struct.Point* %this1, i32 0, i32 0
  %0 = load i32, i32* %x.addr, align 4
  store i32 %0, i32* %x2, align 4
  %y3 = getelementptr inbounds %struct.Point, %struct.Point* %this1, i32 0, i32 1
  %1 = load i32, i32* %y.addr, align 4
  store i32 %1, i32* %y3, align 4
  ret void
}

declare dso_local i32 @printf(i8* noundef, ...) #2

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef i32 @_ZNK5Point9manhattanEv(%struct.Point* noundef nonnull align 4 dereferenceable(8) %this) #3 comdat align 2 {
entry:
  %this.addr = alloca %struct.Point*, align 8
  %ax = alloca i32, align 4
  %ay = alloca i32, align 4
  store %struct.Point* %this, %struct.Point** %this.addr, align 8
  %this1 = load %struct.Point*, %struct.Point** %this.addr, align 8
  %x = getelementptr inbounds %struct.Point, %struct.Point* %this1, i32 0, i32 0
  %0 = load i32, i32* %x, align 4
  %cmp = icmp slt i32 %0, 0
  br i1 %cmp, label %cond.true, label %cond.false

cond.true:                                        ; preds = %entry
  %x2 = getelementptr inbounds %struct.Point, %struct.Point* %this1, i32 0, i32 0
  %1 = load i32, i32* %x2, align 4
  %sub = sub nsw i32 0, %1
  br label %cond.end

cond.false:                                       ; preds = %entry
  %x3 = getelementptr inbounds %struct.Point, %struct.Point* %this1, i32 0, i32 0
  %2 = load i32, i32* %x3, align 4
  br label %cond.end

cond.end:                                         ; preds = %cond.false, %cond.true
  %cond = phi i32 [ %sub, %cond.true ], [ %2, %cond.false ]
  store i32 %cond, i32* %ax, align 4
  %y = getelementptr inbounds %struct.Point, %struct.Point* %this1, i32 0, i32 1
  %3 = load i32, i32* %y, align 4
  %cmp4 = icmp slt i32 %3, 0
  br i1 %cmp4, label %cond.true5, label %cond.false8

cond.true5:                                       ; preds = %cond.end
  %y6 = getelementptr inbounds %struct.Point, %struct.Point* %this1, i32 0, i32 1
  %4 = load i32, i32* %y6, align 4
  %sub7 = sub nsw i32 0, %4
  br label %cond.end10

cond.false8:                                      ; preds = %cond.end
  %y9 = getelementptr inbounds %struct.Point, %struct.Point* %this1, i32 0, i32 1
  %5 = load i32, i32* %y9, align 4
  br label %cond.end10

cond.end10:                                       ; preds = %cond.false8, %cond.true5
  %cond11 = phi i32 [ %sub7, %cond.true5 ], [ %5, %cond.false8 ]
  store i32 %cond11, i32* %ay, align 4
  %6 = load i32, i32* %ax, align 4
  %7 = load i32, i32* %ay, align 4
  %add = add nsw i32 %6, %7
  ret i32 %add
}

attributes #0 = { mustprogress noinline norecurse optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #2 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #3 = { mustprogress noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }

!llvm.module.flags = !{!0, !1, !2}
!llvm.ident = !{!3}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 7, !"uwtable", i32 1}
!2 = !{i32 7, !"frame-pointer", i32 2}
!3 = !{!"clang version 14.0.6 (https://github.com/conda-forge/clangdev-feedstock ceeebe884c3cfd7160cf5a43e147f94439fafee3)"}
