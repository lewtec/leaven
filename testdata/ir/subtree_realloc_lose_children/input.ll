; ModuleID = '/tmp/subtree-realloc-lose-children.c'
source_filename = "/tmp/subtree-realloc-lose-children.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

%struct.SubtreeArray = type { %union.Subtree*, i32, i32 }
%union.Subtree = type { %struct.SubtreeHeapData* }
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
  store i32 0, i32* %retval, align 4
  %call = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([38 x i8], [38 x i8]* @.str, i64 0, i64 0), i64 noundef 8, i64 noundef 80)
  %0 = bitcast %struct.SubtreeArray* %arr to i8*
  call void @llvm.memset.p0i8.i64(i8* align 8 %0, i8 0, i64 16, i1 false)
  %capacity = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %arr, i32 0, i32 2
  store i32 1, i32* %capacity, align 4
  %capacity1 = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %arr, i32 0, i32 2
  %1 = load i32, i32* %capacity1, align 4
  %conv = zext i32 %1 to i64
  %call2 = call noalias i8* @calloc(i64 noundef %conv, i64 noundef 8) #6
  %2 = bitcast i8* %call2 to %union.Subtree*
  %contents = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %arr, i32 0, i32 0
  store %union.Subtree* %2, %union.Subtree** %contents, align 8
  %contents3 = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %arr, i32 0, i32 0
  %3 = load %union.Subtree*, %union.Subtree** %contents3, align 8
  %arrayidx = getelementptr inbounds %union.Subtree, %union.Subtree* %3, i64 0
  %call4 = call %struct.SubtreeHeapData* @make_inline(i8 noundef zeroext 7)
  %coerce.dive = getelementptr inbounds %union.Subtree, %union.Subtree* %tmp, i32 0, i32 0
  store %struct.SubtreeHeapData* %call4, %struct.SubtreeHeapData** %coerce.dive, align 8
  %4 = bitcast %union.Subtree* %arrayidx to i8*
  %5 = bitcast %union.Subtree* %tmp to i8*
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 8 %4, i8* align 8 %5, i64 8, i1 false)
  %size = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %arr, i32 0, i32 1
  store i32 1, i32* %size, align 8
  %contents5 = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %arr, i32 0, i32 0
  %6 = load %union.Subtree*, %union.Subtree** %contents5, align 8
  %arrayidx6 = getelementptr inbounds %union.Subtree, %union.Subtree* %6, i64 0
  %ptr = bitcast %union.Subtree* %arrayidx6 to %struct.SubtreeHeapData**
  %7 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr, align 8
  %8 = ptrtoint %struct.SubtreeHeapData* %7 to i64
  %9 = inttoptr i64 %8 to i8*
  %contents7 = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %arr, i32 0, i32 0
  %10 = load %union.Subtree*, %union.Subtree** %contents7, align 8
  %arrayidx8 = getelementptr inbounds %union.Subtree, %union.Subtree* %10, i64 0
  %data = bitcast %union.Subtree* %arrayidx8 to %struct.SubtreeInlineData*
  %11 = bitcast %struct.SubtreeInlineData* %data to i8*
  %bf.load = load i8, i8* %11, align 8
  %bf.clear = and i8 %bf.load, 1
  %bf.cast = trunc i8 %bf.clear to i1
  %conv9 = zext i1 %bf.cast to i32
  %call10 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([23 x i8], [23 x i8]* @.str.1, i64 0, i64 0), i8* noundef %9, i32 noundef %conv9)
  %call11 = call %struct.SubtreeHeapData* @new_node(%struct.SubtreeArray* noundef %arr)
  %coerce.dive12 = getelementptr inbounds %union.Subtree, %union.Subtree* %p, i32 0, i32 0
  store %struct.SubtreeHeapData* %call11, %struct.SubtreeHeapData** %coerce.dive12, align 8
  %ptr13 = bitcast %union.Subtree* %p to %struct.SubtreeHeapData**
  %12 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr13, align 8
  %child_count = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %12, i32 0, i32 1
  %13 = load i32, i32* %child_count, align 4
  %call14 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([19 x i8], [19 x i8]* @.str.2, i64 0, i64 0), i32 noundef %13)
  %contents15 = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %arr, i32 0, i32 0
  %14 = load %union.Subtree*, %union.Subtree** %contents15, align 8
  %15 = bitcast %union.Subtree* %14 to i8*
  call void @free(i8* noundef %15) #6
  ret i32 0
}

declare dso_local i32 @printf(i8* noundef, ...) #1

; Function Attrs: argmemonly nofree nounwind willreturn writeonly
declare void @llvm.memset.p0i8.i64(i8* nocapture writeonly, i8, i64, i1 immarg) #2

; Function Attrs: nounwind
declare dso_local noalias i8* @calloc(i64 noundef, i64 noundef) #3

; Function Attrs: noinline nounwind optnone uwtable
define internal %struct.SubtreeHeapData* @make_inline(i8 noundef zeroext %sym) #0 {
entry:
  %retval = alloca %union.Subtree, align 8
  %sym.addr = alloca i8, align 1
  store i8 %sym, i8* %sym.addr, align 1
  %0 = bitcast %union.Subtree* %retval to i8*
  call void @llvm.memset.p0i8.i64(i8* align 8 %0, i8 0, i64 8, i1 false)
  %data = bitcast %union.Subtree* %retval to %struct.SubtreeInlineData*
  %1 = bitcast %struct.SubtreeInlineData* %data to i8*
  %bf.load = load i8, i8* %1, align 8
  %bf.clear = and i8 %bf.load, -2
  %bf.set = or i8 %bf.clear, 1
  store i8 %bf.set, i8* %1, align 8
  %data1 = bitcast %union.Subtree* %retval to %struct.SubtreeInlineData*
  %2 = bitcast %struct.SubtreeInlineData* %data1 to i8*
  %bf.load2 = load i8, i8* %2, align 8
  %bf.clear3 = and i8 %bf.load2, -3
  %bf.set4 = or i8 %bf.clear3, 2
  store i8 %bf.set4, i8* %2, align 8
  %data5 = bitcast %union.Subtree* %retval to %struct.SubtreeInlineData*
  %3 = bitcast %struct.SubtreeInlineData* %data5 to i8*
  %bf.load6 = load i8, i8* %3, align 8
  %bf.clear7 = and i8 %bf.load6, -5
  %bf.set8 = or i8 %bf.clear7, 4
  store i8 %bf.set8, i8* %3, align 8
  %4 = load i8, i8* %sym.addr, align 1
  %data9 = bitcast %union.Subtree* %retval to %struct.SubtreeInlineData*
  %symbol = getelementptr inbounds %struct.SubtreeInlineData, %struct.SubtreeInlineData* %data9, i32 0, i32 1
  store i8 %4, i8* %symbol, align 1
  %data10 = bitcast %union.Subtree* %retval to %struct.SubtreeInlineData*
  %parse_state = getelementptr inbounds %struct.SubtreeInlineData, %struct.SubtreeInlineData* %data10, i32 0, i32 2
  store i16 1, i16* %parse_state, align 2
  %data11 = bitcast %union.Subtree* %retval to %struct.SubtreeInlineData*
  %size_bytes = getelementptr inbounds %struct.SubtreeInlineData, %struct.SubtreeInlineData* %data11, i32 0, i32 3
  store i8 1, i8* %size_bytes, align 4
  %coerce.dive = getelementptr inbounds %union.Subtree, %union.Subtree* %retval, i32 0, i32 0
  %5 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %coerce.dive, align 8
  ret %struct.SubtreeHeapData* %5
}

; Function Attrs: argmemonly nofree nounwind willreturn
declare void @llvm.memcpy.p0i8.p0i8.i64(i8* noalias nocapture writeonly, i8* noalias nocapture readonly, i64, i1 immarg) #4

; Function Attrs: noinline nounwind optnone uwtable
define internal %struct.SubtreeHeapData* @new_node(%struct.SubtreeArray* noundef %children) #0 {
entry:
  %retval = alloca %union.Subtree, align 8
  %children.addr = alloca %struct.SubtreeArray*, align 8
  %nbytes = alloca i64, align 8
  %have = alloca i64, align 8
  %data = alloca %struct.SubtreeHeapData*, align 8
  %ch = alloca %union.Subtree*, align 8
  %i = alloca i32, align 4
  %child = alloca %union.Subtree, align 8
  store %struct.SubtreeArray* %children, %struct.SubtreeArray** %children.addr, align 8
  %0 = load %struct.SubtreeArray*, %struct.SubtreeArray** %children.addr, align 8
  %size = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %0, i32 0, i32 1
  %1 = load i32, i32* %size, align 8
  %call = call i64 @alloc_size(i32 noundef %1)
  store i64 %call, i64* %nbytes, align 8
  %2 = load %struct.SubtreeArray*, %struct.SubtreeArray** %children.addr, align 8
  %capacity = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %2, i32 0, i32 2
  %3 = load i32, i32* %capacity, align 4
  %conv = zext i32 %3 to i64
  %mul = mul i64 %conv, 8
  store i64 %mul, i64* %have, align 8
  %4 = load %struct.SubtreeArray*, %struct.SubtreeArray** %children.addr, align 8
  %size1 = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %4, i32 0, i32 1
  %5 = load i32, i32* %size1, align 8
  %6 = load %struct.SubtreeArray*, %struct.SubtreeArray** %children.addr, align 8
  %capacity2 = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %6, i32 0, i32 2
  %7 = load i32, i32* %capacity2, align 4
  %8 = load i64, i64* %nbytes, align 8
  %9 = load i64, i64* %have, align 8
  %10 = load i64, i64* %have, align 8
  %11 = load i64, i64* %nbytes, align 8
  %cmp = icmp ult i64 %10, %11
  %conv3 = zext i1 %cmp to i32
  %call4 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([66 x i8], [66 x i8]* @.str.3, i64 0, i64 0), i32 noundef %5, i32 noundef %7, i64 noundef %8, i64 noundef %9, i32 noundef %conv3)
  %12 = load i64, i64* %have, align 8
  %13 = load i64, i64* %nbytes, align 8
  %cmp5 = icmp ult i64 %12, %13
  br i1 %cmp5, label %if.then, label %if.end

if.then:                                          ; preds = %entry
  %14 = load %struct.SubtreeArray*, %struct.SubtreeArray** %children.addr, align 8
  %contents = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %14, i32 0, i32 0
  %15 = load %union.Subtree*, %union.Subtree** %contents, align 8
  %16 = bitcast %union.Subtree* %15 to i8*
  %17 = load i64, i64* %nbytes, align 8
  %call7 = call i8* @ts_realloc(i8* noundef %16, i64 noundef %17)
  %18 = bitcast i8* %call7 to %union.Subtree*
  %19 = load %struct.SubtreeArray*, %struct.SubtreeArray** %children.addr, align 8
  %contents8 = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %19, i32 0, i32 0
  store %union.Subtree* %18, %union.Subtree** %contents8, align 8
  %20 = load i64, i64* %nbytes, align 8
  %div = udiv i64 %20, 8
  %conv9 = trunc i64 %div to i32
  %21 = load %struct.SubtreeArray*, %struct.SubtreeArray** %children.addr, align 8
  %capacity10 = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %21, i32 0, i32 2
  store i32 %conv9, i32* %capacity10, align 4
  br label %if.end

if.end:                                           ; preds = %if.then, %entry
  %22 = load %struct.SubtreeArray*, %struct.SubtreeArray** %children.addr, align 8
  %contents11 = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %22, i32 0, i32 0
  %23 = load %union.Subtree*, %union.Subtree** %contents11, align 8
  %24 = load %struct.SubtreeArray*, %struct.SubtreeArray** %children.addr, align 8
  %size12 = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %24, i32 0, i32 1
  %25 = load i32, i32* %size12, align 8
  %idxprom = zext i32 %25 to i64
  %arrayidx = getelementptr inbounds %union.Subtree, %union.Subtree* %23, i64 %idxprom
  %26 = bitcast %union.Subtree* %arrayidx to %struct.SubtreeHeapData*
  store %struct.SubtreeHeapData* %26, %struct.SubtreeHeapData** %data, align 8
  %27 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %data, align 8
  %28 = bitcast %struct.SubtreeHeapData* %27 to i8*
  call void @llvm.memset.p0i8.i64(i8* align 4 %28, i8 0, i64 80, i1 false)
  %29 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %data, align 8
  %ref_count = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %29, i32 0, i32 0
  store i32 1, i32* %ref_count, align 4
  %30 = load %struct.SubtreeArray*, %struct.SubtreeArray** %children.addr, align 8
  %size13 = getelementptr inbounds %struct.SubtreeArray, %struct.SubtreeArray* %30, i32 0, i32 1
  %31 = load i32, i32* %size13, align 8
  %32 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %data, align 8
  %child_count = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %32, i32 0, i32 1
  store i32 %31, i32* %child_count, align 4
  %33 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %data, align 8
  %symbol = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %33, i32 0, i32 2
  store i16 100, i16* %symbol, align 4
  %34 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %data, align 8
  %size_row = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %34, i32 0, i32 4
  store i32 0, i32* %size_row, align 4
  %35 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %data, align 8
  %ptr = bitcast %union.Subtree* %retval to %struct.SubtreeHeapData**
  store %struct.SubtreeHeapData* %35, %struct.SubtreeHeapData** %ptr, align 8
  %ptr14 = bitcast %union.Subtree* %retval to %struct.SubtreeHeapData**
  %36 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr14, align 8
  %37 = bitcast %struct.SubtreeHeapData* %36 to %union.Subtree*
  %ptr15 = bitcast %union.Subtree* %retval to %struct.SubtreeHeapData**
  %38 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr15, align 8
  %child_count16 = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %38, i32 0, i32 1
  %39 = load i32, i32* %child_count16, align 4
  %idx.ext = zext i32 %39 to i64
  %idx.neg = sub i64 0, %idx.ext
  %add.ptr = getelementptr inbounds %union.Subtree, %union.Subtree* %37, i64 %idx.neg
  store %union.Subtree* %add.ptr, %union.Subtree** %ch, align 8
  store i32 0, i32* %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc, %if.end
  %40 = load i32, i32* %i, align 4
  %ptr17 = bitcast %union.Subtree* %retval to %struct.SubtreeHeapData**
  %41 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr17, align 8
  %child_count18 = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %41, i32 0, i32 1
  %42 = load i32, i32* %child_count18, align 4
  %cmp19 = icmp ult i32 %40, %42
  br i1 %cmp19, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %43 = load %union.Subtree*, %union.Subtree** %ch, align 8
  %44 = load i32, i32* %i, align 4
  %idxprom21 = zext i32 %44 to i64
  %arrayidx22 = getelementptr inbounds %union.Subtree, %union.Subtree* %43, i64 %idxprom21
  %45 = bitcast %union.Subtree* %child to i8*
  %46 = bitcast %union.Subtree* %arrayidx22 to i8*
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 8 %45, i8* align 8 %46, i64 8, i1 false)
  %47 = load i32, i32* %i, align 4
  %ptr23 = bitcast %union.Subtree* %child to %struct.SubtreeHeapData**
  %48 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr23, align 8
  %49 = ptrtoint %struct.SubtreeHeapData* %48 to i64
  %50 = inttoptr i64 %49 to i8*
  %data24 = bitcast %union.Subtree* %child to %struct.SubtreeInlineData*
  %51 = bitcast %struct.SubtreeInlineData* %data24 to i8*
  %bf.load = load i8, i8* %51, align 8
  %bf.clear = and i8 %bf.load, 1
  %bf.cast = trunc i8 %bf.clear to i1
  %conv25 = zext i1 %bf.cast to i32
  %call26 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([28 x i8], [28 x i8]* @.str.4, i64 0, i64 0), i32 noundef %47, i8* noundef %50, i32 noundef %conv25)
  %ptr27 = bitcast %union.Subtree* %retval to %struct.SubtreeHeapData**
  %52 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr27, align 8
  %size_row28 = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %52, i32 0, i32 4
  %53 = load i32, i32* %size_row28, align 4
  %cmp29 = icmp eq i32 %53, 0
  br i1 %cmp29, label %land.lhs.true, label %if.end37

land.lhs.true:                                    ; preds = %for.body
  %coerce.dive = getelementptr inbounds %union.Subtree, %union.Subtree* %child, i32 0, i32 0
  %54 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %coerce.dive, align 8
  %call31 = call zeroext i1 @depends_on_column(%struct.SubtreeHeapData* %54)
  br i1 %call31, label %if.then33, label %if.end37

if.then33:                                        ; preds = %land.lhs.true
  %ptr34 = bitcast %union.Subtree* %retval to %struct.SubtreeHeapData**
  %55 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr34, align 8
  %flags = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %55, i32 0, i32 3
  %56 = load i16, i16* %flags, align 2
  %conv35 = zext i16 %56 to i32
  %or = or i32 %conv35, 256
  %conv36 = trunc i32 %or to i16
  store i16 %conv36, i16* %flags, align 2
  br label %if.end37

if.end37:                                         ; preds = %if.then33, %land.lhs.true, %for.body
  br label %for.inc

for.inc:                                          ; preds = %if.end37
  %57 = load i32, i32* %i, align 4
  %inc = add i32 %57, 1
  store i32 %inc, i32* %i, align 4
  br label %for.cond, !llvm.loop !4

for.end:                                          ; preds = %for.cond
  %coerce.dive38 = getelementptr inbounds %union.Subtree, %union.Subtree* %retval, i32 0, i32 0
  %58 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %coerce.dive38, align 8
  ret %struct.SubtreeHeapData* %58
}

; Function Attrs: nounwind
declare dso_local void @free(i8* noundef) #3

; Function Attrs: noinline nounwind optnone uwtable
define internal i64 @alloc_size(i32 noundef %n) #0 {
entry:
  %n.addr = alloca i32, align 4
  store i32 %n, i32* %n.addr, align 4
  %0 = load i32, i32* %n.addr, align 4
  %conv = zext i32 %0 to i64
  %mul = mul i64 %conv, 8
  %add = add i64 %mul, 80
  ret i64 %add
}

; Function Attrs: noinline nounwind optnone uwtable
define internal i8* @ts_realloc(i8* noundef %p, i64 noundef %n) #0 {
entry:
  %p.addr = alloca i8*, align 8
  %n.addr = alloca i64, align 8
  %q = alloca i8*, align 8
  store i8* %p, i8** %p.addr, align 8
  store i64 %n, i64* %n.addr, align 8
  %0 = load i8*, i8** %p.addr, align 8
  %1 = load i64, i64* %n.addr, align 8
  %call = call i8* @realloc(i8* noundef %0, i64 noundef %1) #6
  store i8* %call, i8** %q, align 8
  %2 = load i8*, i8** %q, align 8
  %tobool = icmp ne i8* %2, null
  br i1 %tobool, label %if.end, label %land.lhs.true

land.lhs.true:                                    ; preds = %entry
  %3 = load i64, i64* %n.addr, align 8
  %tobool1 = icmp ne i64 %3, 0
  br i1 %tobool1, label %if.then, label %if.end

if.then:                                          ; preds = %land.lhs.true
  call void @abort() #7
  unreachable

if.end:                                           ; preds = %land.lhs.true, %entry
  %4 = load i8*, i8** %q, align 8
  ret i8* %4
}

; Function Attrs: noinline nounwind optnone uwtable
define internal zeroext i1 @depends_on_column(%struct.SubtreeHeapData* %self.coerce) #0 {
entry:
  %retval = alloca i1, align 1
  %self = alloca %union.Subtree, align 8
  %coerce.dive = getelementptr inbounds %union.Subtree, %union.Subtree* %self, i32 0, i32 0
  store %struct.SubtreeHeapData* %self.coerce, %struct.SubtreeHeapData** %coerce.dive, align 8
  %data = bitcast %union.Subtree* %self to %struct.SubtreeInlineData*
  %0 = bitcast %struct.SubtreeInlineData* %data to i8*
  %bf.load = load i8, i8* %0, align 8
  %bf.clear = and i8 %bf.load, 1
  %bf.cast = trunc i8 %bf.clear to i1
  br i1 %bf.cast, label %if.then, label %if.end

if.then:                                          ; preds = %entry
  store i1 false, i1* %retval, align 1
  br label %return

if.end:                                           ; preds = %entry
  %ptr = bitcast %union.Subtree* %self to %struct.SubtreeHeapData**
  %1 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr, align 8
  %flags = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %1, i32 0, i32 3
  %2 = load i16, i16* %flags, align 2
  %conv = zext i16 %2 to i32
  %and = and i32 %conv, 256
  %cmp = icmp ne i32 %and, 0
  store i1 %cmp, i1* %retval, align 1
  br label %return

return:                                           ; preds = %if.end, %if.then
  %3 = load i1, i1* %retval, align 1
  ret i1 %3
}

; Function Attrs: nounwind
declare dso_local i8* @realloc(i8* noundef, i64 noundef) #3

; Function Attrs: noreturn nounwind
declare dso_local void @abort() #5

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #2 = { argmemonly nofree nounwind willreturn writeonly }
attributes #3 = { nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #4 = { argmemonly nofree nounwind willreturn }
attributes #5 = { noreturn nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #6 = { nounwind }
attributes #7 = { noreturn nounwind }

!llvm.module.flags = !{!0, !1, !2}
!llvm.ident = !{!3}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 7, !"uwtable", i32 1}
!2 = !{i32 7, !"frame-pointer", i32 2}
!3 = !{!"clang version 14.0.6 (https://github.com/conda-forge/clangdev-feedstock ceeebe884c3cfd7160cf5a43e147f94439fafee3)"}
!4 = distinct !{!4, !5}
!5 = !{!"llvm.loop.mustprogress"}
