; ModuleID = '/tmp/subtree-tagged-repro.c'
source_filename = "/tmp/subtree-tagged-repro.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

%union.Subtree = type { %struct.SubtreeHeapData* }
%struct.SubtreeHeapData = type { i32, i32, i16, i16 }
%struct.SubtreeInlineData = type { i8, i8, i16, i8, i8, i8, i8 }

@.str = private unnamed_addr constant [47 x i8] c"sizeof(Subtree)=%zu (should be pointer-sized)\0A\00", align 1
@.str.1 = private unnamed_addr constant [28 x i8] c"inline bit=%d  raw_word=%p\0A\00", align 1
@.str.2 = private unnamed_addr constant [28 x i8] c"inline depends=%d (want 0)\0A\00", align 1
@.str.3 = private unnamed_addr constant [28 x i8] c"heap   depends=%d (want 1)\0A\00", align 1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %leaf = alloca %union.Subtree, align 8
  %parent = alloca %union.Subtree, align 8
  store i32 0, i32* %retval, align 4
  %call = call %struct.SubtreeHeapData* @make_inline_leaf(i8 noundef zeroext 7)
  %coerce.dive = getelementptr inbounds %union.Subtree, %union.Subtree* %leaf, i32 0, i32 0
  store %struct.SubtreeHeapData* %call, %struct.SubtreeHeapData** %coerce.dive, align 8
  %call1 = call %struct.SubtreeHeapData* @make_heap_node(i1 noundef zeroext true)
  %coerce.dive2 = getelementptr inbounds %union.Subtree, %union.Subtree* %parent, i32 0, i32 0
  store %struct.SubtreeHeapData* %call1, %struct.SubtreeHeapData** %coerce.dive2, align 8
  %call3 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([47 x i8], [47 x i8]* @.str, i64 0, i64 0), i64 noundef 8)
  %data = bitcast %union.Subtree* %leaf to %struct.SubtreeInlineData*
  %0 = bitcast %struct.SubtreeInlineData* %data to i8*
  %bf.load = load i8, i8* %0, align 8
  %bf.clear = and i8 %bf.load, 1
  %bf.cast = trunc i8 %bf.clear to i1
  %conv = zext i1 %bf.cast to i32
  %ptr = bitcast %union.Subtree* %leaf to %struct.SubtreeHeapData**
  %1 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr, align 8
  %2 = ptrtoint %struct.SubtreeHeapData* %1 to i64
  %3 = inttoptr i64 %2 to i8*
  %call4 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([28 x i8], [28 x i8]* @.str.1, i64 0, i64 0), i32 noundef %conv, i8* noundef %3)
  %coerce.dive5 = getelementptr inbounds %union.Subtree, %union.Subtree* %leaf, i32 0, i32 0
  %4 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %coerce.dive5, align 8
  %call6 = call zeroext i1 @ts_subtree_depends_on_column(%struct.SubtreeHeapData* %4)
  %conv7 = zext i1 %call6 to i32
  %call8 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([28 x i8], [28 x i8]* @.str.2, i64 0, i64 0), i32 noundef %conv7)
  %coerce.dive9 = getelementptr inbounds %union.Subtree, %union.Subtree* %parent, i32 0, i32 0
  %5 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %coerce.dive9, align 8
  %call10 = call zeroext i1 @ts_subtree_depends_on_column(%struct.SubtreeHeapData* %5)
  %conv11 = zext i1 %call10 to i32
  %call12 = call i32 (i8*, ...) @printf(i8* noundef getelementptr inbounds ([28 x i8], [28 x i8]* @.str.3, i64 0, i64 0), i32 noundef %conv11)
  %ptr13 = bitcast %union.Subtree* %parent to %struct.SubtreeHeapData**
  %6 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %ptr13, align 8
  %7 = bitcast %struct.SubtreeHeapData* %6 to i8*
  call void @free(i8* noundef %7) #4
  ret i32 0
}

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
  %size_bytes = getelementptr inbounds %struct.SubtreeInlineData, %struct.SubtreeInlineData* %data12, i32 0, i32 6
  store i8 1, i8* %size_bytes, align 1
  %coerce.dive = getelementptr inbounds %union.Subtree, %union.Subtree* %retval, i32 0, i32 0
  %5 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %coerce.dive, align 8
  ret %struct.SubtreeHeapData* %5
}

; Function Attrs: noinline nounwind optnone uwtable
define internal %struct.SubtreeHeapData* @make_heap_node(i1 noundef zeroext %depends) #0 {
entry:
  %retval = alloca %union.Subtree, align 8
  %depends.addr = alloca i8, align 1
  %h = alloca %struct.SubtreeHeapData*, align 8
  %frombool = zext i1 %depends to i8
  store i8 %frombool, i8* %depends.addr, align 1
  %call = call noalias i8* @calloc(i64 noundef 1, i64 noundef 12) #4
  %0 = bitcast i8* %call to %struct.SubtreeHeapData*
  store %struct.SubtreeHeapData* %0, %struct.SubtreeHeapData** %h, align 8
  %1 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %h, align 8
  %ref_count = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %1, i32 0, i32 0
  store i32 1, i32* %ref_count, align 4
  %2 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %h, align 8
  %child_count = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %2, i32 0, i32 1
  store i32 2, i32* %child_count, align 4
  %3 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %h, align 8
  %symbol = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %3, i32 0, i32 2
  store i16 42, i16* %symbol, align 4
  %4 = load i8, i8* %depends.addr, align 1
  %tobool = trunc i8 %4 to i1
  br i1 %tobool, label %if.then, label %if.end

if.then:                                          ; preds = %entry
  %5 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %h, align 8
  %flags = getelementptr inbounds %struct.SubtreeHeapData, %struct.SubtreeHeapData* %5, i32 0, i32 3
  %6 = load i16, i16* %flags, align 2
  %conv = zext i16 %6 to i32
  %or = or i32 %conv, 256
  %conv1 = trunc i32 %or to i16
  store i16 %conv1, i16* %flags, align 2
  br label %if.end

if.end:                                           ; preds = %if.then, %entry
  %7 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %h, align 8
  %ptr = bitcast %union.Subtree* %retval to %struct.SubtreeHeapData**
  store %struct.SubtreeHeapData* %7, %struct.SubtreeHeapData** %ptr, align 8
  %coerce.dive = getelementptr inbounds %union.Subtree, %union.Subtree* %retval, i32 0, i32 0
  %8 = load %struct.SubtreeHeapData*, %struct.SubtreeHeapData** %coerce.dive, align 8
  ret %struct.SubtreeHeapData* %8
}

declare dso_local i32 @printf(i8* noundef, ...) #1

; Function Attrs: noinline nounwind optnone uwtable
define internal zeroext i1 @ts_subtree_depends_on_column(%struct.SubtreeHeapData* %self.coerce) #0 {
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
declare dso_local void @free(i8* noundef) #2

; Function Attrs: argmemonly nofree nounwind willreturn writeonly
declare void @llvm.memset.p0i8.i64(i8* nocapture writeonly, i8, i64, i1 immarg) #3

; Function Attrs: nounwind
declare dso_local noalias i8* @calloc(i64 noundef, i64 noundef) #2

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #2 = { nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #3 = { argmemonly nofree nounwind willreturn writeonly }
attributes #4 = { nounwind }

!llvm.module.flags = !{!0, !1, !2}
!llvm.ident = !{!3}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 7, !"uwtable", i32 1}
!2 = !{i32 7, !"frame-pointer", i32 2}
!3 = !{!"clang version 14.0.6 (https://github.com/conda-forge/clangdev-feedstock ceeebe884c3cfd7160cf5a43e147f94439fafee3)"}
