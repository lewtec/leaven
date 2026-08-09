; ModuleID = 'testdata/ir/c_subtree_realloc_lose_children/source.c'
source_filename = "testdata/ir/c_subtree_realloc_lose_children/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

%struct.SubtreeArray = type { ptr, i32, i32 }
%union.Subtree = type { ptr }
%struct.SubtreeHeapData = type { i32, i32, i16, i16, i32, [64 x i8] }
%struct.SubtreeInlineData = type { i8, i8, i16, i8 }

@.str = private unnamed_addr constant [38 x i8] c"sizeof(Subtree)=%zu sizeof(Heap)=%zu\0A\00", align 1
@.str.1 = private unnamed_addr constant [23 x i8] c"leaf raw=%p inline=%d\0A\00", align 1
@.str.2 = private unnamed_addr constant [19 x i8] c"ok child_count=%u\0A\00", align 1
@.str.3 = private unnamed_addr constant [66 x i8] c"before realloc: size=%u cap=%u need=%zu have=%zu will_realloc=%d\0A\00", align 1
@.str.4 = private unnamed_addr constant [28 x i8] c"child[%u] raw=%p inline=%d\0A\00", align 1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %arr = alloca %struct.SubtreeArray, align 8
  %tmp = alloca %union.Subtree, align 8
  %p = alloca %union.Subtree, align 8
  store i32 0, ptr %retval, align 4
  %call = call i32 (ptr, ...) @printf(ptr noundef @.str, i64 noundef 8, i64 noundef 80)
  call void @llvm.memset.p0.i64(ptr align 8 %arr, i8 0, i64 16, i1 false)
  %capacity = getelementptr inbounds nuw %struct.SubtreeArray, ptr %arr, i32 0, i32 2
  store i32 1, ptr %capacity, align 4
  %capacity1 = getelementptr inbounds nuw %struct.SubtreeArray, ptr %arr, i32 0, i32 2
  %0 = load i32, ptr %capacity1, align 4
  %conv = zext i32 %0 to i64
  %call2 = call noalias ptr @calloc(i64 noundef %conv, i64 noundef 8) #8
  %contents = getelementptr inbounds nuw %struct.SubtreeArray, ptr %arr, i32 0, i32 0
  store ptr %call2, ptr %contents, align 8
  %contents3 = getelementptr inbounds nuw %struct.SubtreeArray, ptr %arr, i32 0, i32 0
  %1 = load ptr, ptr %contents3, align 8
  %arrayidx = getelementptr inbounds %union.Subtree, ptr %1, i64 0
  %call4 = call ptr @make_inline(i8 noundef zeroext 7)
  %coerce.dive = getelementptr inbounds nuw %union.Subtree, ptr %tmp, i32 0, i32 0
  store ptr %call4, ptr %coerce.dive, align 8
  call void @llvm.memcpy.p0.p0.i64(ptr align 8 %arrayidx, ptr align 8 %tmp, i64 8, i1 false)
  %size = getelementptr inbounds nuw %struct.SubtreeArray, ptr %arr, i32 0, i32 1
  store i32 1, ptr %size, align 8
  %contents5 = getelementptr inbounds nuw %struct.SubtreeArray, ptr %arr, i32 0, i32 0
  %2 = load ptr, ptr %contents5, align 8
  %arrayidx6 = getelementptr inbounds %union.Subtree, ptr %2, i64 0
  %3 = load ptr, ptr %arrayidx6, align 8
  %4 = ptrtoint ptr %3 to i64
  %5 = inttoptr i64 %4 to ptr
  %contents7 = getelementptr inbounds nuw %struct.SubtreeArray, ptr %arr, i32 0, i32 0
  %6 = load ptr, ptr %contents7, align 8
  %arrayidx8 = getelementptr inbounds %union.Subtree, ptr %6, i64 0
  %bf.load = load i8, ptr %arrayidx8, align 8
  %bf.clear = and i8 %bf.load, 1
  %bf.cast = trunc i8 %bf.clear to i1
  %conv9 = zext i1 %bf.cast to i32
  %call10 = call i32 (ptr, ...) @printf(ptr noundef @.str.1, ptr noundef %5, i32 noundef %conv9)
  %call11 = call ptr @new_node(ptr noundef %arr)
  %coerce.dive12 = getelementptr inbounds nuw %union.Subtree, ptr %p, i32 0, i32 0
  store ptr %call11, ptr %coerce.dive12, align 8
  %7 = load ptr, ptr %p, align 8
  %child_count = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %7, i32 0, i32 1
  %8 = load i32, ptr %child_count, align 4
  %call13 = call i32 (ptr, ...) @printf(ptr noundef @.str.2, i32 noundef %8)
  %contents14 = getelementptr inbounds nuw %struct.SubtreeArray, ptr %arr, i32 0, i32 0
  %9 = load ptr, ptr %contents14, align 8
  call void @free(ptr noundef %9) #9
  ret i32 0
}

declare i32 @printf(ptr noundef, ...) #1

; Function Attrs: nocallback nofree nounwind willreturn memory(argmem: write)
declare void @llvm.memset.p0.i64(ptr writeonly captures(none), i8, i64, i1 immarg) #2

; Function Attrs: nounwind allocsize(0,1)
declare noalias ptr @calloc(i64 noundef, i64 noundef) #3

; Function Attrs: noinline nounwind optnone uwtable
define internal ptr @make_inline(i8 noundef zeroext %sym) #0 {
entry:
  %retval = alloca %union.Subtree, align 8
  %sym.addr = alloca i8, align 1
  store i8 %sym, ptr %sym.addr, align 1
  call void @llvm.memset.p0.i64(ptr align 8 %retval, i8 0, i64 8, i1 false)
  %bf.load = load i8, ptr %retval, align 8
  %bf.clear = and i8 %bf.load, -2
  %bf.set = or i8 %bf.clear, 1
  store i8 %bf.set, ptr %retval, align 8
  %bf.load1 = load i8, ptr %retval, align 8
  %bf.clear2 = and i8 %bf.load1, -3
  %bf.set3 = or i8 %bf.clear2, 2
  store i8 %bf.set3, ptr %retval, align 8
  %bf.load4 = load i8, ptr %retval, align 8
  %bf.clear5 = and i8 %bf.load4, -5
  %bf.set6 = or i8 %bf.clear5, 4
  store i8 %bf.set6, ptr %retval, align 8
  %0 = load i8, ptr %sym.addr, align 1
  %symbol = getelementptr inbounds nuw %struct.SubtreeInlineData, ptr %retval, i32 0, i32 1
  store i8 %0, ptr %symbol, align 1
  %parse_state = getelementptr inbounds nuw %struct.SubtreeInlineData, ptr %retval, i32 0, i32 2
  store i16 1, ptr %parse_state, align 2
  %size_bytes = getelementptr inbounds nuw %struct.SubtreeInlineData, ptr %retval, i32 0, i32 3
  store i8 1, ptr %size_bytes, align 4
  %coerce.dive = getelementptr inbounds nuw %union.Subtree, ptr %retval, i32 0, i32 0
  %1 = load ptr, ptr %coerce.dive, align 8
  ret ptr %1
}

; Function Attrs: nocallback nofree nounwind willreturn memory(argmem: readwrite)
declare void @llvm.memcpy.p0.p0.i64(ptr noalias writeonly captures(none), ptr noalias readonly captures(none), i64, i1 immarg) #4

; Function Attrs: noinline nounwind optnone uwtable
define internal ptr @new_node(ptr noundef %children) #0 {
entry:
  %retval = alloca %union.Subtree, align 8
  %children.addr = alloca ptr, align 8
  %nbytes = alloca i64, align 8
  %have = alloca i64, align 8
  %data = alloca ptr, align 8
  %ch = alloca ptr, align 8
  %i = alloca i32, align 4
  %child = alloca %union.Subtree, align 8
  store ptr %children, ptr %children.addr, align 8
  %0 = load ptr, ptr %children.addr, align 8
  %size = getelementptr inbounds nuw %struct.SubtreeArray, ptr %0, i32 0, i32 1
  %1 = load i32, ptr %size, align 8
  %call = call i64 @alloc_size(i32 noundef %1)
  store i64 %call, ptr %nbytes, align 8
  %2 = load ptr, ptr %children.addr, align 8
  %capacity = getelementptr inbounds nuw %struct.SubtreeArray, ptr %2, i32 0, i32 2
  %3 = load i32, ptr %capacity, align 4
  %conv = zext i32 %3 to i64
  %mul = mul i64 %conv, 8
  store i64 %mul, ptr %have, align 8
  %4 = load ptr, ptr %children.addr, align 8
  %size1 = getelementptr inbounds nuw %struct.SubtreeArray, ptr %4, i32 0, i32 1
  %5 = load i32, ptr %size1, align 8
  %6 = load ptr, ptr %children.addr, align 8
  %capacity2 = getelementptr inbounds nuw %struct.SubtreeArray, ptr %6, i32 0, i32 2
  %7 = load i32, ptr %capacity2, align 4
  %8 = load i64, ptr %nbytes, align 8
  %9 = load i64, ptr %have, align 8
  %10 = load i64, ptr %have, align 8
  %11 = load i64, ptr %nbytes, align 8
  %cmp = icmp ult i64 %10, %11
  %conv3 = zext i1 %cmp to i32
  %call4 = call i32 (ptr, ...) @printf(ptr noundef @.str.3, i32 noundef %5, i32 noundef %7, i64 noundef %8, i64 noundef %9, i32 noundef %conv3)
  %12 = load i64, ptr %have, align 8
  %13 = load i64, ptr %nbytes, align 8
  %cmp5 = icmp ult i64 %12, %13
  br i1 %cmp5, label %if.then, label %if.end

if.then:                                          ; preds = %entry
  %14 = load ptr, ptr %children.addr, align 8
  %contents = getelementptr inbounds nuw %struct.SubtreeArray, ptr %14, i32 0, i32 0
  %15 = load ptr, ptr %contents, align 8
  %16 = load i64, ptr %nbytes, align 8
  %call7 = call ptr @ts_realloc(ptr noundef %15, i64 noundef %16)
  %17 = load ptr, ptr %children.addr, align 8
  %contents8 = getelementptr inbounds nuw %struct.SubtreeArray, ptr %17, i32 0, i32 0
  store ptr %call7, ptr %contents8, align 8
  %18 = load i64, ptr %nbytes, align 8
  %div = udiv i64 %18, 8
  %conv9 = trunc i64 %div to i32
  %19 = load ptr, ptr %children.addr, align 8
  %capacity10 = getelementptr inbounds nuw %struct.SubtreeArray, ptr %19, i32 0, i32 2
  store i32 %conv9, ptr %capacity10, align 4
  br label %if.end

if.end:                                           ; preds = %if.then, %entry
  %20 = load ptr, ptr %children.addr, align 8
  %contents11 = getelementptr inbounds nuw %struct.SubtreeArray, ptr %20, i32 0, i32 0
  %21 = load ptr, ptr %contents11, align 8
  %22 = load ptr, ptr %children.addr, align 8
  %size12 = getelementptr inbounds nuw %struct.SubtreeArray, ptr %22, i32 0, i32 1
  %23 = load i32, ptr %size12, align 8
  %idxprom = zext i32 %23 to i64
  %arrayidx = getelementptr inbounds nuw %union.Subtree, ptr %21, i64 %idxprom
  store ptr %arrayidx, ptr %data, align 8
  %24 = load ptr, ptr %data, align 8
  call void @llvm.memset.p0.i64(ptr align 4 %24, i8 0, i64 80, i1 false)
  %25 = load ptr, ptr %data, align 8
  %ref_count = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %25, i32 0, i32 0
  store i32 1, ptr %ref_count, align 4
  %26 = load ptr, ptr %children.addr, align 8
  %size13 = getelementptr inbounds nuw %struct.SubtreeArray, ptr %26, i32 0, i32 1
  %27 = load i32, ptr %size13, align 8
  %28 = load ptr, ptr %data, align 8
  %child_count = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %28, i32 0, i32 1
  store i32 %27, ptr %child_count, align 4
  %29 = load ptr, ptr %data, align 8
  %symbol = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %29, i32 0, i32 2
  store i16 100, ptr %symbol, align 4
  %30 = load ptr, ptr %data, align 8
  %size_row = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %30, i32 0, i32 4
  store i32 0, ptr %size_row, align 4
  %31 = load ptr, ptr %data, align 8
  store ptr %31, ptr %retval, align 8
  %32 = load ptr, ptr %retval, align 8
  %33 = load ptr, ptr %retval, align 8
  %child_count14 = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %33, i32 0, i32 1
  %34 = load i32, ptr %child_count14, align 4
  %idx.ext = zext i32 %34 to i64
  %idx.neg = sub i64 0, %idx.ext
  %add.ptr = getelementptr inbounds %union.Subtree, ptr %32, i64 %idx.neg
  store ptr %add.ptr, ptr %ch, align 8
  store i32 0, ptr %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc, %if.end
  %35 = load i32, ptr %i, align 4
  %36 = load ptr, ptr %retval, align 8
  %child_count15 = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %36, i32 0, i32 1
  %37 = load i32, ptr %child_count15, align 4
  %cmp16 = icmp ult i32 %35, %37
  br i1 %cmp16, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %38 = load ptr, ptr %ch, align 8
  %39 = load i32, ptr %i, align 4
  %idxprom18 = zext i32 %39 to i64
  %arrayidx19 = getelementptr inbounds nuw %union.Subtree, ptr %38, i64 %idxprom18
  call void @llvm.memcpy.p0.p0.i64(ptr align 8 %child, ptr align 8 %arrayidx19, i64 8, i1 false)
  %40 = load i32, ptr %i, align 4
  %41 = load ptr, ptr %child, align 8
  %42 = ptrtoint ptr %41 to i64
  %43 = inttoptr i64 %42 to ptr
  %bf.load = load i8, ptr %child, align 8
  %bf.clear = and i8 %bf.load, 1
  %bf.cast = trunc i8 %bf.clear to i1
  %conv20 = zext i1 %bf.cast to i32
  %call21 = call i32 (ptr, ...) @printf(ptr noundef @.str.4, i32 noundef %40, ptr noundef %43, i32 noundef %conv20)
  %44 = load ptr, ptr %retval, align 8
  %size_row22 = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %44, i32 0, i32 4
  %45 = load i32, ptr %size_row22, align 4
  %cmp23 = icmp eq i32 %45, 0
  br i1 %cmp23, label %land.lhs.true, label %if.end30

land.lhs.true:                                    ; preds = %for.body
  %coerce.dive = getelementptr inbounds nuw %union.Subtree, ptr %child, i32 0, i32 0
  %46 = load ptr, ptr %coerce.dive, align 8
  %call25 = call zeroext i1 @depends_on_column(ptr %46)
  br i1 %call25, label %if.then27, label %if.end30

if.then27:                                        ; preds = %land.lhs.true
  %47 = load ptr, ptr %retval, align 8
  %flags = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %47, i32 0, i32 3
  %48 = load i16, ptr %flags, align 2
  %conv28 = zext i16 %48 to i32
  %or = or i32 %conv28, 256
  %conv29 = trunc i32 %or to i16
  store i16 %conv29, ptr %flags, align 2
  br label %if.end30

if.end30:                                         ; preds = %if.then27, %land.lhs.true, %for.body
  br label %for.inc

for.inc:                                          ; preds = %if.end30
  %49 = load i32, ptr %i, align 4
  %inc = add i32 %49, 1
  store i32 %inc, ptr %i, align 4
  br label %for.cond, !llvm.loop !6

for.end:                                          ; preds = %for.cond
  %coerce.dive31 = getelementptr inbounds nuw %union.Subtree, ptr %retval, i32 0, i32 0
  %50 = load ptr, ptr %coerce.dive31, align 8
  ret ptr %50
}

; Function Attrs: nounwind
declare void @free(ptr noundef) #5

; Function Attrs: noinline nounwind optnone uwtable
define internal i64 @alloc_size(i32 noundef %n) #0 {
entry:
  %n.addr = alloca i32, align 4
  store i32 %n, ptr %n.addr, align 4
  %0 = load i32, ptr %n.addr, align 4
  %conv = zext i32 %0 to i64
  %mul = mul i64 %conv, 8
  %add = add i64 %mul, 80
  ret i64 %add
}

; Function Attrs: noinline nounwind optnone uwtable
define internal ptr @ts_realloc(ptr noundef %p, i64 noundef %n) #0 {
entry:
  %p.addr = alloca ptr, align 8
  %n.addr = alloca i64, align 8
  %q = alloca ptr, align 8
  store ptr %p, ptr %p.addr, align 8
  store i64 %n, ptr %n.addr, align 8
  %0 = load ptr, ptr %p.addr, align 8
  %1 = load i64, ptr %n.addr, align 8
  %call = call ptr @realloc(ptr noundef %0, i64 noundef %1) #10
  store ptr %call, ptr %q, align 8
  %2 = load ptr, ptr %q, align 8
  %tobool = icmp ne ptr %2, null
  br i1 %tobool, label %if.end, label %land.lhs.true

land.lhs.true:                                    ; preds = %entry
  %3 = load i64, ptr %n.addr, align 8
  %tobool1 = icmp ne i64 %3, 0
  br i1 %tobool1, label %if.then, label %if.end

if.then:                                          ; preds = %land.lhs.true
  call void @abort() #11
  unreachable

if.end:                                           ; preds = %land.lhs.true, %entry
  %4 = load ptr, ptr %q, align 8
  ret ptr %4
}

; Function Attrs: noinline nounwind optnone uwtable
define internal zeroext i1 @depends_on_column(ptr %self.coerce) #0 {
entry:
  %retval = alloca i1, align 1
  %self = alloca %union.Subtree, align 8
  %coerce.dive = getelementptr inbounds nuw %union.Subtree, ptr %self, i32 0, i32 0
  store ptr %self.coerce, ptr %coerce.dive, align 8
  %bf.load = load i8, ptr %self, align 8
  %bf.clear = and i8 %bf.load, 1
  %bf.cast = trunc i8 %bf.clear to i1
  br i1 %bf.cast, label %if.then, label %if.end

if.then:                                          ; preds = %entry
  store i1 false, ptr %retval, align 1
  br label %return

if.end:                                           ; preds = %entry
  %0 = load ptr, ptr %self, align 8
  %flags = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %0, i32 0, i32 3
  %1 = load i16, ptr %flags, align 2
  %conv = zext i16 %1 to i32
  %and = and i32 %conv, 256
  %cmp = icmp ne i32 %and, 0
  store i1 %cmp, ptr %retval, align 1
  br label %return

return:                                           ; preds = %if.end, %if.then
  %2 = load i1, ptr %retval, align 1
  ret i1 %2
}

; Function Attrs: nounwind allocsize(1)
declare ptr @realloc(ptr noundef, i64 noundef) #6

; Function Attrs: noreturn nounwind
declare void @abort() #7

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #2 = { nocallback nofree nounwind willreturn memory(argmem: write) }
attributes #3 = { nounwind allocsize(0,1) "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #4 = { nocallback nofree nounwind willreturn memory(argmem: readwrite) }
attributes #5 = { nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #6 = { nounwind allocsize(1) "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #7 = { noreturn nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #8 = { nounwind allocsize(0,1) }
attributes #9 = { nounwind }
attributes #10 = { nounwind allocsize(1) }
attributes #11 = { noreturn nounwind }

!llvm.module.flags = !{!0, !1, !2, !3, !4}
!llvm.ident = !{!5}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 8, !"PIC Level", i32 2}
!2 = !{i32 7, !"PIE Level", i32 2}
!3 = !{i32 7, !"uwtable", i32 2}
!4 = !{i32 7, !"frame-pointer", i32 2}
!5 = !{!"clang version 22.1.8 (https://github.com/conda-forge/clangdev-feedstock 015bdba1263c0b3ebb3c518ff5947fbd99692bd0)"}
!6 = distinct !{!6, !7}
!7 = !{!"llvm.loop.mustprogress"}
