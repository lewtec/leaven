; ModuleID = 'testdata/ir/c_subtree_children_before_heap/source.c'
source_filename = "testdata/ir/c_subtree_children_before_heap/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

%union.Subtree = type { ptr }
%struct.SubtreeHeapData = type { i32, i32, i16, i16, i32 }
%struct.SubtreeInlineData = type { i8, i8, i16, i8 }

@.str = private unnamed_addr constant [40 x i8] c"sizeof(Subtree)=%zu (want 8 on 64-bit)\0A\00", align 1
@.str.1 = private unnamed_addr constant [26 x i8] c"leaf raw=%p is_inline=%d\0A\00", align 1
@.str.2 = private unnamed_addr constant [34 x i8] c"parent child_count=%u depends=%d\0A\00", align 1
@stderr = external global ptr, align 8
@.str.3 = private unnamed_addr constant [31 x i8] c"child[%u] raw=%p is_inline=%d\0A\00", align 1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %leaf = alloca %union.Subtree, align 8
  %kids = alloca [1 x %union.Subtree], align 8
  %parent = alloca %union.Subtree, align 8
  store i32 0, ptr %retval, align 4
  %call = call i32 (ptr, ...) @printf(ptr noundef @.str, i64 noundef 8)
  %call1 = call ptr @make_inline_leaf(i8 noundef zeroext 7)
  %coerce.dive = getelementptr inbounds nuw %union.Subtree, ptr %leaf, i32 0, i32 0
  store ptr %call1, ptr %coerce.dive, align 8
  %0 = load ptr, ptr %leaf, align 8
  %1 = ptrtoint ptr %0 to i64
  %2 = inttoptr i64 %1 to ptr
  %bf.load = load i8, ptr %leaf, align 8
  %bf.clear = and i8 %bf.load, 1
  %bf.cast = trunc i8 %bf.clear to i1
  %conv = zext i1 %bf.cast to i32
  %call2 = call i32 (ptr, ...) @printf(ptr noundef @.str.1, ptr noundef %2, i32 noundef %conv)
  call void @llvm.memcpy.p0.p0.i64(ptr align 8 %kids, ptr align 8 %leaf, i64 8, i1 false)
  %arraydecay = getelementptr inbounds [1 x %union.Subtree], ptr %kids, i64 0, i64 0
  %call3 = call ptr @new_node_with_children(ptr noundef %arraydecay, i32 noundef 1)
  %coerce.dive4 = getelementptr inbounds nuw %union.Subtree, ptr %parent, i32 0, i32 0
  store ptr %call3, ptr %coerce.dive4, align 8
  %3 = load ptr, ptr %parent, align 8
  %child_count = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %3, i32 0, i32 1
  %4 = load i32, ptr %child_count, align 4
  %5 = load ptr, ptr %parent, align 8
  %flags = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %5, i32 0, i32 3
  %6 = load i16, ptr %flags, align 2
  %conv5 = zext i16 %6 to i32
  %and = and i32 %conv5, 256
  %cmp = icmp ne i32 %and, 0
  %conv6 = zext i1 %cmp to i32
  %call7 = call i32 (ptr, ...) @printf(ptr noundef @.str.2, i32 noundef %4, i32 noundef %conv6)
  %7 = load ptr, ptr %parent, align 8
  %8 = load ptr, ptr %parent, align 8
  %child_count8 = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %8, i32 0, i32 1
  %9 = load i32, ptr %child_count8, align 4
  %idx.ext = zext i32 %9 to i64
  %idx.neg = sub i64 0, %idx.ext
  %add.ptr = getelementptr inbounds %union.Subtree, ptr %7, i64 %idx.neg
  call void @free(ptr noundef %add.ptr) #7
  ret i32 0
}

declare i32 @printf(ptr noundef, ...) #1

; Function Attrs: noinline nounwind optnone uwtable
define internal ptr @make_inline_leaf(i8 noundef zeroext %symbol) #0 {
entry:
  %retval = alloca %union.Subtree, align 8
  %symbol.addr = alloca i8, align 1
  store i8 %symbol, ptr %symbol.addr, align 1
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
  %0 = load i8, ptr %symbol.addr, align 1
  %symbol7 = getelementptr inbounds nuw %struct.SubtreeInlineData, ptr %retval, i32 0, i32 1
  store i8 %0, ptr %symbol7, align 1
  %parse_state = getelementptr inbounds nuw %struct.SubtreeInlineData, ptr %retval, i32 0, i32 2
  store i16 1, ptr %parse_state, align 2
  %size_bytes = getelementptr inbounds nuw %struct.SubtreeInlineData, ptr %retval, i32 0, i32 3
  store i8 1, ptr %size_bytes, align 4
  %coerce.dive = getelementptr inbounds nuw %union.Subtree, ptr %retval, i32 0, i32 0
  %1 = load ptr, ptr %coerce.dive, align 8
  ret ptr %1
}

; Function Attrs: nocallback nofree nounwind willreturn memory(argmem: readwrite)
declare void @llvm.memcpy.p0.p0.i64(ptr noalias writeonly captures(none), ptr noalias readonly captures(none), i64, i1 immarg) #2

; Function Attrs: noinline nounwind optnone uwtable
define internal ptr @new_node_with_children(ptr noundef %kids, i32 noundef %n) #0 {
entry:
  %retval = alloca %union.Subtree, align 8
  %kids.addr = alloca ptr, align 8
  %n.addr = alloca i32, align 4
  %block = alloca ptr, align 8
  %slot = alloca ptr, align 8
  %i = alloca i32, align 4
  %data = alloca ptr, align 8
  %ch = alloca ptr, align 8
  %i7 = alloca i32, align 4
  %child = alloca %union.Subtree, align 8
  store ptr %kids, ptr %kids.addr, align 8
  store i32 %n, ptr %n.addr, align 4
  %0 = load i32, ptr %n.addr, align 4
  %call = call i64 @alloc_size(i32 noundef %0)
  %call1 = call noalias ptr @calloc(i64 noundef 1, i64 noundef %call) #8
  store ptr %call1, ptr %block, align 8
  %1 = load ptr, ptr %block, align 8
  %tobool = icmp ne ptr %1, null
  br i1 %tobool, label %if.end, label %if.then

if.then:                                          ; preds = %entry
  call void @abort() #9
  unreachable

if.end:                                           ; preds = %entry
  %2 = load ptr, ptr %block, align 8
  store ptr %2, ptr %slot, align 8
  store i32 0, ptr %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc, %if.end
  %3 = load i32, ptr %i, align 4
  %4 = load i32, ptr %n.addr, align 4
  %cmp = icmp ult i32 %3, %4
  br i1 %cmp, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %5 = load ptr, ptr %slot, align 8
  %6 = load i32, ptr %i, align 4
  %idxprom = zext i32 %6 to i64
  %arrayidx = getelementptr inbounds nuw %union.Subtree, ptr %5, i64 %idxprom
  %7 = load ptr, ptr %kids.addr, align 8
  %8 = load i32, ptr %i, align 4
  %idxprom2 = zext i32 %8 to i64
  %arrayidx3 = getelementptr inbounds nuw %union.Subtree, ptr %7, i64 %idxprom2
  call void @llvm.memcpy.p0.p0.i64(ptr align 8 %arrayidx, ptr align 8 %arrayidx3, i64 8, i1 false)
  br label %for.inc

for.inc:                                          ; preds = %for.body
  %9 = load i32, ptr %i, align 4
  %inc = add i32 %9, 1
  store i32 %inc, ptr %i, align 4
  br label %for.cond, !llvm.loop !6

for.end:                                          ; preds = %for.cond
  %10 = load ptr, ptr %slot, align 8
  %11 = load i32, ptr %n.addr, align 4
  %idxprom4 = zext i32 %11 to i64
  %arrayidx5 = getelementptr inbounds nuw %union.Subtree, ptr %10, i64 %idxprom4
  store ptr %arrayidx5, ptr %data, align 8
  %12 = load ptr, ptr %data, align 8
  %ref_count = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %12, i32 0, i32 0
  store i32 1, ptr %ref_count, align 4
  %13 = load i32, ptr %n.addr, align 4
  %14 = load ptr, ptr %data, align 8
  %child_count = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %14, i32 0, i32 1
  store i32 %13, ptr %child_count, align 4
  %15 = load ptr, ptr %data, align 8
  %symbol = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %15, i32 0, i32 2
  store i16 100, ptr %symbol, align 4
  %16 = load ptr, ptr %data, align 8
  %flags = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %16, i32 0, i32 3
  store i16 0, ptr %flags, align 2
  %17 = load ptr, ptr %data, align 8
  %size_row = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %17, i32 0, i32 4
  store i32 0, ptr %size_row, align 4
  %18 = load ptr, ptr %data, align 8
  store ptr %18, ptr %retval, align 8
  %coerce.dive = getelementptr inbounds nuw %union.Subtree, ptr %retval, i32 0, i32 0
  %19 = load ptr, ptr %coerce.dive, align 8
  %call6 = call ptr @children_of(ptr %19)
  store ptr %call6, ptr %ch, align 8
  store i32 0, ptr %i7, align 4
  br label %for.cond8

for.cond8:                                        ; preds = %for.inc26, %for.end
  %20 = load i32, ptr %i7, align 4
  %21 = load ptr, ptr %retval, align 8
  %child_count9 = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %21, i32 0, i32 1
  %22 = load i32, ptr %child_count9, align 4
  %cmp10 = icmp ult i32 %20, %22
  br i1 %cmp10, label %for.body11, label %for.end28

for.body11:                                       ; preds = %for.cond8
  %23 = load ptr, ptr %ch, align 8
  %24 = load i32, ptr %i7, align 4
  %idxprom12 = zext i32 %24 to i64
  %arrayidx13 = getelementptr inbounds nuw %union.Subtree, ptr %23, i64 %idxprom12
  call void @llvm.memcpy.p0.p0.i64(ptr align 8 %child, ptr align 8 %arrayidx13, i64 8, i1 false)
  %25 = load ptr, ptr @stderr, align 8
  %26 = load i32, ptr %i7, align 4
  %27 = load ptr, ptr %child, align 8
  %28 = ptrtoint ptr %27 to i64
  %29 = inttoptr i64 %28 to ptr
  %bf.load = load i8, ptr %child, align 8
  %bf.clear = and i8 %bf.load, 1
  %bf.cast = trunc i8 %bf.clear to i1
  %conv = zext i1 %bf.cast to i32
  %call14 = call i32 (ptr, ptr, ...) @fprintf(ptr noundef %25, ptr noundef @.str.3, i32 noundef %26, ptr noundef %29, i32 noundef %conv) #7
  %30 = load ptr, ptr %retval, align 8
  %size_row15 = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %30, i32 0, i32 4
  %31 = load i32, ptr %size_row15, align 4
  %cmp16 = icmp eq i32 %31, 0
  br i1 %cmp16, label %land.lhs.true, label %if.end25

land.lhs.true:                                    ; preds = %for.body11
  %coerce.dive18 = getelementptr inbounds nuw %union.Subtree, ptr %child, i32 0, i32 0
  %32 = load ptr, ptr %coerce.dive18, align 8
  %call19 = call zeroext i1 @depends_on_column(ptr %32)
  br i1 %call19, label %if.then21, label %if.end25

if.then21:                                        ; preds = %land.lhs.true
  %33 = load ptr, ptr %retval, align 8
  %flags22 = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %33, i32 0, i32 3
  %34 = load i16, ptr %flags22, align 2
  %conv23 = zext i16 %34 to i32
  %or = or i32 %conv23, 256
  %conv24 = trunc i32 %or to i16
  store i16 %conv24, ptr %flags22, align 2
  br label %if.end25

if.end25:                                         ; preds = %if.then21, %land.lhs.true, %for.body11
  br label %for.inc26

for.inc26:                                        ; preds = %if.end25
  %35 = load i32, ptr %i7, align 4
  %inc27 = add i32 %35, 1
  store i32 %inc27, ptr %i7, align 4
  br label %for.cond8, !llvm.loop !8

for.end28:                                        ; preds = %for.cond8
  %coerce.dive29 = getelementptr inbounds nuw %union.Subtree, ptr %retval, i32 0, i32 0
  %36 = load ptr, ptr %coerce.dive29, align 8
  ret ptr %36
}

; Function Attrs: nounwind
declare void @free(ptr noundef) #3

; Function Attrs: nocallback nofree nounwind willreturn memory(argmem: write)
declare void @llvm.memset.p0.i64(ptr writeonly captures(none), i8, i64, i1 immarg) #4

; Function Attrs: nounwind allocsize(0,1)
declare noalias ptr @calloc(i64 noundef, i64 noundef) #5

; Function Attrs: noinline nounwind optnone uwtable
define internal i64 @alloc_size(i32 noundef %n) #0 {
entry:
  %n.addr = alloca i32, align 4
  store i32 %n, ptr %n.addr, align 4
  %0 = load i32, ptr %n.addr, align 4
  %conv = zext i32 %0 to i64
  %mul = mul i64 %conv, 8
  %add = add i64 %mul, 16
  ret i64 %add
}

; Function Attrs: noreturn nounwind
declare void @abort() #6

; Function Attrs: noinline nounwind optnone uwtable
define internal ptr @children_of(ptr %self.coerce) #0 {
entry:
  %retval = alloca ptr, align 8
  %self = alloca %union.Subtree, align 8
  %coerce.dive = getelementptr inbounds nuw %union.Subtree, ptr %self, i32 0, i32 0
  store ptr %self.coerce, ptr %coerce.dive, align 8
  %bf.load = load i8, ptr %self, align 8
  %bf.clear = and i8 %bf.load, 1
  %bf.cast = trunc i8 %bf.clear to i1
  br i1 %bf.cast, label %if.then, label %if.end

if.then:                                          ; preds = %entry
  store ptr null, ptr %retval, align 8
  br label %return

if.end:                                           ; preds = %entry
  %0 = load ptr, ptr %self, align 8
  %1 = load ptr, ptr %self, align 8
  %child_count = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %1, i32 0, i32 1
  %2 = load i32, ptr %child_count, align 4
  %idx.ext = zext i32 %2 to i64
  %idx.neg = sub i64 0, %idx.ext
  %add.ptr = getelementptr inbounds %union.Subtree, ptr %0, i64 %idx.neg
  store ptr %add.ptr, ptr %retval, align 8
  br label %return

return:                                           ; preds = %if.end, %if.then
  %3 = load ptr, ptr %retval, align 8
  ret ptr %3
}

; Function Attrs: nounwind
declare i32 @fprintf(ptr noundef, ptr noundef, ...) #3

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

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #2 = { nocallback nofree nounwind willreturn memory(argmem: readwrite) }
attributes #3 = { nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #4 = { nocallback nofree nounwind willreturn memory(argmem: write) }
attributes #5 = { nounwind allocsize(0,1) "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #6 = { noreturn nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #7 = { nounwind }
attributes #8 = { nounwind allocsize(0,1) }
attributes #9 = { noreturn nounwind }

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
!8 = distinct !{!8, !7}
