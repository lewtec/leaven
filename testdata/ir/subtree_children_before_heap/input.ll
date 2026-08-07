; ModuleID = '/tmp/subtree-children-before-heap.c'
source_filename = "/tmp/subtree-children-before-heap.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

%struct._IO_FILE = type { i32, i8*, i8*, i8*, i8*, i8*, i8*, i8*, i8*, i8*, i8*, i8*, %struct._IO_marker*, %struct._IO_FILE*, i32, i32, i64, i16, i8, [1 x i8], i8*, i64, %struct._IO_codecvt*, %struct._IO_wide_data*, %struct._IO_FILE*, i8*, i64, i32, [20 x i8] }
%struct._IO_marker = type opaque
%struct._IO_codecvt = type opaque
%struct._IO_wide_data = type opaque
%union.Subtree = type { %struct.SubtreeHeapData* }
%struct.SubtreeHeapData = type { i32, i32, i16, i16, i32 }
%struct.SubtreeInlineData = type { i8, i8, i16, i8 }

@.str = private unnamed_addr constant [40 x i8] c"sizeof(Subtree)=%zu (want 8 on 64-bit)\0A\00", align 1
@.str.1 = private unnamed_addr constant [26 x i8] c"leaf raw=%p is_inline=%d\0A\00", align 1
@.str.2 = private unnamed_addr constant [34 x i8] c"parent child_count=%u depends=%d\0A\00", align 1
@stderr = external dso_local global %struct._IO_FILE*, align 8
@.str.3 = private unnamed_addr constant [31 x i8] c"child[%u] raw=%p is_inline=%d\0A\00", align 1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %leaf = alloca %union.Subtree, align 8
  %kids = alloca [1 x %union.Subtree], align 8
  %parent = alloca %union.Subtree, align 8
  store i32 0, i32* %retval, align 4
  %call = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([40 x i8], [40 x i8]* @.str, i64 0, i64 0), i64 noundef 8)
  %call1 = call %struct.SubtreeHeapData* @make_inline_leaf(i8 noundef zeroext 7)
  %coerce.dive = getelementptr inbounds %union.Subtree, %union.Subtree* %leaf, i32 0, i32 0
  store %struct.SubtreeHeapData* %call1, %struct.SubtreeHeapData** %coerce.dive, align 8
  %ptr = bitcast %union.Subtree* %leaf to %struct.SubtreeHeapData**
  %0 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr, align 8
  %1 = ptrtoint %struct.SubtreeHeapData* %0 to i64
  %2 = inttoptr i64 %1 to i8*
  %data = bitcast %union.Subtree* %leaf to %struct.SubtreeInlineData*
  %3 = bitcast %struct.SubtreeInlineData* %data to i8*
  %bf.load = load i8, i8* %3, align 8
  %bf.clear = and i8 %bf.load, 1
  %bf.cast = trunc i8 %bf.clear to i1
  %conv = zext i1 %bf.cast to i32
  %call2 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([26 x i8], [26 x i8]* @.str.1, i64 0, i64 0), i8* noundef %2, i32 noundef %conv)
  %arrayinit.begin = getelementptr inbounds [1 x %union.Subtree], [1 x %union.Subtree]* %kids, i64 0, i64 0
  %4 = bitcast %union.Subtree* %arrayinit.begin to i8*
  %5 = bitcast %union.Subtree* %leaf to i8*
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 8 %4, i8* align 8 %5, i64 8, i1 false)
  %arraydecay = getelementptr inbounds [1 x %union.Subtree], [1 x %union.Subtree]* %kids, i64 0, i64 0
  %call3 = call %struct.SubtreeHeapData* @new_node_with_children(%union.Subtree* noundef %arraydecay, i32 noundef 1)
  %coerce.dive4 = getelementptr inbounds %union.Subtree, %union.Subtree* %parent, i32 0, i32 0
  store %struct.SubtreeHeapData* %call3, %struct.SubtreeHeapData** %coerce.dive4, align 8
  %ptr5 = bitcast %union.Subtree* %parent to %struct.SubtreeHeapData**
  %6 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr5, align 8
  %child_count = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %6, i32 0, i32 1
  %7 = load i32, i32* %child_count, align 4
  %ptr6 = bitcast %union.Subtree* %parent to %struct.SubtreeHeapData**
  %8 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr6, align 8
  %flags = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %8, i32 0, i32 3
  %9 = load i16, i16* %flags, align 2
  %conv7 = zext i16 %9 to i32
  %and = and i32 %conv7, 256
  %cmp = icmp ne i32 %and, 0
  %conv8 = zext i1 %cmp to i32
  %call9 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([34 x i8], [34 x i8]* @.str.2, i64 0, i64 0), i32 noundef %7, i32 noundef %conv8)
  %ptr10 = bitcast %union.Subtree* %parent to %struct.SubtreeHeapData**
  %10 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr10, align 8
  %11 = bitcast %struct.SubtreeHeapData* %10 to %union.Subtree*
  %ptr11 = bitcast %union.Subtree* %parent to %struct.SubtreeHeapData**
  %12 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr11, align 8
  %child_count12 = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %12, i32 0, i32 1
  %13 = load i32, i32* %child_count12, align 4
  %idx.ext = zext i32 %13 to i64
  %idx.neg = sub i64 0, %idx.ext
  %add.ptr = getelementptr inbounds %union.Subtree, %union.Subtree* %11, i64 %idx.neg
  %14 = bitcast %union.Subtree* %add.ptr to i8*
  call void @free(i8* noundef %14) #6
  ret i32 0
}

declare dso_local i32 @printf(i8* noundef, ...) #1

; Function Attrs: noinline nounwind optnone uwtable
define internal %struct.SubtreeHeapData* @make_inline_leaf(i8 noundef zeroext %symbol) #0 {
entry:
  %retval = alloca %union.Subtree, align 8
  %symbol.addr = alloca i8, align 1
  store i8 %symbol, i8* %symbol.addr, align 1
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
  %4 = load i8, i8* %symbol.addr, align 1
  %data9 = bitcast %union.Subtree* %retval to %struct.SubtreeInlineData*
  %symbol10 = getelementptr inbounds %struct.SubtreeInlineData, %struct.SubtreeInlineData* %data9, i32 0, i32 1
  store i8 %4, i8* %symbol10, align 1
  %data11 = bitcast %union.Subtree* %retval to %struct.SubtreeInlineData*
  %parse_state = getelementptr inbounds %struct.SubtreeInlineData, %struct.SubtreeInlineData* %data11, i32 0, i32 2
  store i16 1, i16* %parse_state, align 2
  %data12 = bitcast %union.Subtree* %retval to %struct.SubtreeInlineData*
  %size_bytes = getelementptr inbounds %struct.SubtreeInlineData, %struct.SubtreeInlineData* %data12, i32 0, i32 3
  store i8 1, i8* %size_bytes, align 4
  %coerce.dive = getelementptr inbounds %union.Subtree, %union.Subtree* %retval, i32 0, i32 0
  %5 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %coerce.dive, align 8
  ret %struct.SubtreeHeapData* %5
}

; Function Attrs: argmemonly nofree nounwind willreturn
declare void @llvm.memcpy.p0i8.p0i8.i64(i8* noalias nocapture writeonly, i8* noalias nocapture readonly, i64, i1 immarg) #2

; Function Attrs: noinline nounwind optnone uwtable
define internal %struct.SubtreeHeapData* @new_node_with_children(%union.Subtree* noundef %kids, i32 noundef %n) #0 {
entry:
  %retval = alloca %union.Subtree, align 8
  %kids.addr = alloca %union.Subtree*, align 8
  %n.addr = alloca i32, align 4
  %block = alloca i8*, align 8
  %slot = alloca %union.Subtree*, align 8
  %i = alloca i32, align 4
  %data = alloca %struct.SubtreeHeapData*, align 8
  %ch = alloca %union.Subtree*, align 8
  %i7 = alloca i32, align 4
  %child = alloca %union.Subtree, align 8
  store %union.Subtree* %kids, %union.Subtree** %kids.addr, align 8
  store i32 %n, i32* %n.addr, align 4
  %0 = load i32, i32* %n.addr, align 4
  %call = call i64 @alloc_size(i32 noundef %0)
  %call1 = call noalias i8* @calloc(i64 noundef 1, i64 noundef %call) #6
  store i8* %call1, i8** %block, align 8
  %1 = load i8*, i8** %block, align 8
  %tobool = icmp ne i8* %1, null
  br i1 %tobool, label %if.end, label %if.then

if.then:                                          ; preds = %entry
  call void @abort() #7
  unreachable

if.end:                                           ; preds = %entry
  %2 = load i8*, i8** %block, align 8
  %3 = bitcast i8* %2 to %union.Subtree*
  store %union.Subtree* %3, %union.Subtree** %slot, align 8
  store i32 0, i32* %i, align 4
  br label %for.cond

for.cond:                                         ; preds = %for.inc, %if.end
  %4 = load i32, i32* %i, align 4
  %5 = load i32, i32* %n.addr, align 4
  %cmp = icmp ult i32 %4, %5
  br i1 %cmp, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %6 = load %union.Subtree*, %union.Subtree** %slot, align 8
  %7 = load i32, i32* %i, align 4
  %idxprom = zext i32 %7 to i64
  %arrayidx = getelementptr inbounds %union.Subtree, %union.Subtree* %6, i64 %idxprom
  %8 = load %union.Subtree*, %union.Subtree** %kids.addr, align 8
  %9 = load i32, i32* %i, align 4
  %idxprom2 = zext i32 %9 to i64
  %arrayidx3 = getelementptr inbounds %union.Subtree, %union.Subtree* %8, i64 %idxprom2
  %10 = bitcast %union.Subtree* %arrayidx to i8*
  %11 = bitcast %union.Subtree* %arrayidx3 to i8*
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 8 %10, i8* align 8 %11, i64 8, i1 false)
  br label %for.inc

for.inc:                                          ; preds = %for.body
  %12 = load i32, i32* %i, align 4
  %inc = add i32 %12, 1
  store i32 %inc, i32* %i, align 4
  br label %for.cond, !llvm.loop !4

for.end:                                          ; preds = %for.cond
  %13 = load %union.Subtree*, %union.Subtree** %slot, align 8
  %14 = load i32, i32* %n.addr, align 4
  %idxprom4 = zext i32 %14 to i64
  %arrayidx5 = getelementptr inbounds %union.Subtree, %union.Subtree* %13, i64 %idxprom4
  %15 = bitcast %union.Subtree* %arrayidx5 to %struct.SubtreeHeapData*
  store %struct.SubtreeHeapData* %15, %struct.SubtreeHeapData** %data, align 8
  %16 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %data, align 8
  %ref_count = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %16, i32 0, i32 0
  store i32 1, i32* %ref_count, align 4
  %17 = load i32, i32* %n.addr, align 4
  %18 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %data, align 8
  %child_count = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %18, i32 0, i32 1
  store i32 %17, i32* %child_count, align 4
  %19 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %data, align 8
  %symbol = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %19, i32 0, i32 2
  store i16 100, i16* %symbol, align 4
  %20 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %data, align 8
  %flags = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %20, i32 0, i32 3
  store i16 0, i16* %flags, align 2
  %21 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %data, align 8
  %size_row = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %21, i32 0, i32 4
  store i32 0, i32* %size_row, align 4
  %22 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %data, align 8
  %ptr = bitcast %union.Subtree* %retval to %struct.SubtreeHeapData**
  store %struct.SubtreeHeapData* %22, %struct.SubtreeHeapData** %ptr, align 8
  %coerce.dive = getelementptr inbounds %union.Subtree, %union.Subtree* %retval, i32 0, i32 0
  %23 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %coerce.dive, align 8
  %call6 = call %union.Subtree* @children_of(%struct.SubtreeHeapData* %23)
  store %union.Subtree* %call6, %union.Subtree** %ch, align 8
  store i32 0, i32* %i7, align 4
  br label %for.cond8

for.cond8:                                        ; preds = %for.inc31, %for.end
  %24 = load i32, i32* %i7, align 4
  %ptr9 = bitcast %union.Subtree* %retval to %struct.SubtreeHeapData**
  %25 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr9, align 8
  %child_count10 = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %25, i32 0, i32 1
  %26 = load i32, i32* %child_count10, align 4
  %cmp11 = icmp ult i32 %24, %26
  br i1 %cmp11, label %for.body12, label %for.end33

for.body12:                                       ; preds = %for.cond8
  %27 = load %union.Subtree*, %union.Subtree** %ch, align 8
  %28 = load i32, i32* %i7, align 4
  %idxprom13 = zext i32 %28 to i64
  %arrayidx14 = getelementptr inbounds %union.Subtree, %union.Subtree* %27, i64 %idxprom13
  %29 = bitcast %union.Subtree* %child to i8*
  %30 = bitcast %union.Subtree* %arrayidx14 to i8*
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 8 %29, i8* align 8 %30, i64 8, i1 false)
  %31 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %32 = load i32, i32* %i7, align 4
  %ptr15 = bitcast %union.Subtree* %child to %struct.SubtreeHeapData**
  %33 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr15, align 8
  %34 = ptrtoint %struct.SubtreeHeapData* %33 to i64
  %35 = inttoptr i64 %34 to i8*
  %data16 = bitcast %union.Subtree* %child to %struct.SubtreeInlineData*
  %36 = bitcast %struct.SubtreeInlineData* %data16 to i8*
  %bf.load = load i8, i8* %36, align 8
  %bf.clear = and i8 %bf.load, 1
  %bf.cast = trunc i8 %bf.clear to i1
  %conv = zext i1 %bf.cast to i32
  %call17 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* noundef %31, i8* noundef getelementptr inbounds ([31 x i8], [31 x i8]* @.str.3, i64 0, i64 0), i32 noundef %32, i8* noundef %35, i32 noundef %conv)
  %ptr18 = bitcast %union.Subtree* %retval to %struct.SubtreeHeapData**
  %37 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr18, align 8
  %size_row19 = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %37, i32 0, i32 4
  %38 = load i32, i32* %size_row19, align 4
  %cmp20 = icmp eq i32 %38, 0
  br i1 %cmp20, label %land.lhs.true, label %if.end30

land.lhs.true:                                    ; preds = %for.body12
  %coerce.dive22 = getelementptr inbounds %union.Subtree, %union.Subtree* %child, i32 0, i32 0
  %39 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %coerce.dive22, align 8
  %call23 = call zeroext i1 @depends_on_column(%struct.SubtreeHeapData* %39)
  br i1 %call23, label %if.then25, label %if.end30

if.then25:                                        ; preds = %land.lhs.true
  %ptr26 = bitcast %union.Subtree* %retval to %struct.SubtreeHeapData**
  %40 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr26, align 8
  %flags27 = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %40, i32 0, i32 3
  %41 = load i16, i16* %flags27, align 2
  %conv28 = zext i16 %41 to i32
  %or = or i32 %conv28, 256
  %conv29 = trunc i32 %or to i16
  store i16 %conv29, i16* %flags27, align 2
  br label %if.end30

if.end30:                                         ; preds = %if.then25, %land.lhs.true, %for.body12
  br label %for.inc31

for.inc31:                                        ; preds = %if.end30
  %42 = load i32, i32* %i7, align 4
  %inc32 = add i32 %42, 1
  store i32 %inc32, i32* %i7, align 4
  br label %for.cond8, !llvm.loop !6

for.end33:                                        ; preds = %for.cond8
  %coerce.dive34 = getelementptr inbounds %union.Subtree, %union.Subtree* %retval, i32 0, i32 0
  %43 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %coerce.dive34, align 8
  ret %struct.SubtreeHeapData* %43
}

; Function Attrs: nounwind
declare dso_local void @free(i8* noundef) #3

; Function Attrs: argmemonly nofree nounwind willreturn writeonly
declare void @llvm.memset.p0i8.i64(i8* nocapture writeonly, i8, i64, i1 immarg) #4

; Function Attrs: nounwind
declare dso_local noalias i8* @calloc(i64 noundef, i64 noundef) #3

; Function Attrs: noinline nounwind optnone uwtable
define internal i64 @alloc_size(i32 noundef %n) #0 {
entry:
  %n.addr = alloca i32, align 4
  store i32 %n, i32* %n.addr, align 4
  %0 = load i32, i32* %n.addr, align 4
  %conv = zext i32 %0 to i64
  %mul = mul i64 %conv, 8
  %add = add i64 %mul, 16
  ret i64 %add
}

; Function Attrs: noreturn nounwind
declare dso_local void @abort() #5

; Function Attrs: noinline nounwind optnone uwtable
define internal %union.Subtree* @children_of(%struct.SubtreeHeapData* %self.coerce) #0 {
entry:
  %retval = alloca %union.Subtree*, align 8
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
  store %union.Subtree* null, %union.Subtree** %retval, align 8
  br label %return

if.end:                                           ; preds = %entry
  %ptr = bitcast %union.Subtree* %self to %struct.SubtreeHeapData**
  %1 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr, align 8
  %2 = bitcast %struct.SubtreeHeapData* %1 to %union.Subtree*
  %ptr1 = bitcast %union.Subtree* %self to %struct.SubtreeHeapData**
  %3 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr1, align 8
  %child_count = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %3, i32 0, i32 1
  %4 = load i32, i32* %child_count, align 4
  %idx.ext = zext i32 %4 to i64
  %idx.neg = sub i64 0, %idx.ext
  %add.ptr = getelementptr inbounds %union.Subtree, %union.Subtree* %2, i64 %idx.neg
  store %union.Subtree* %add.ptr, %union.Subtree** %retval, align 8
  br label %return

return:                                           ; preds = %if.end, %if.then
  %5 = load %union.Subtree*, %union.Subtree** %retval, align 8
  ret %union.Subtree* %5
}

declare dso_local i32 @fprintf(%struct._IO_FILE* noundef, i8* noundef, ...) #1

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

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #2 = { argmemonly nofree nounwind willreturn }
attributes #3 = { nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #4 = { argmemonly nofree nounwind willreturn writeonly }
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
!6 = distinct !{!6, !5}
