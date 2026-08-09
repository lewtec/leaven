; ModuleID = 'testdata/ir/c_subtree_tagged/source.c'
source_filename = "testdata/ir/c_subtree_tagged/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

%union.Subtree = type { ptr }
%struct.SubtreeInlineData = type { i8, i8, i16, i8, i8, i8, i8 }
%struct.SubtreeHeapData = type { i32, i32, i16, i16 }

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
  store i32 0, ptr %retval, align 4
  %call = call ptr @make_inline_leaf(i8 noundef zeroext 7)
  %coerce.dive = getelementptr inbounds nuw %union.Subtree, ptr %leaf, i32 0, i32 0
  store ptr %call, ptr %coerce.dive, align 8
  %call1 = call ptr @make_heap_node(i1 noundef zeroext true)
  %coerce.dive2 = getelementptr inbounds nuw %union.Subtree, ptr %parent, i32 0, i32 0
  store ptr %call1, ptr %coerce.dive2, align 8
  %call3 = call i32 (ptr, ...) @printf(ptr noundef @.str, i64 noundef 8)
  %bf.load = load i8, ptr %leaf, align 8
  %bf.clear = and i8 %bf.load, 1
  %bf.cast = trunc i8 %bf.clear to i1
  %conv = zext i1 %bf.cast to i32
  %0 = load ptr, ptr %leaf, align 8
  %1 = ptrtoint ptr %0 to i64
  %2 = inttoptr i64 %1 to ptr
  %call4 = call i32 (ptr, ...) @printf(ptr noundef @.str.1, i32 noundef %conv, ptr noundef %2)
  %coerce.dive5 = getelementptr inbounds nuw %union.Subtree, ptr %leaf, i32 0, i32 0
  %3 = load ptr, ptr %coerce.dive5, align 8
  %call6 = call zeroext i1 @ts_subtree_depends_on_column(ptr %3)
  %conv7 = zext i1 %call6 to i32
  %call8 = call i32 (ptr, ...) @printf(ptr noundef @.str.2, i32 noundef %conv7)
  %coerce.dive9 = getelementptr inbounds nuw %union.Subtree, ptr %parent, i32 0, i32 0
  %4 = load ptr, ptr %coerce.dive9, align 8
  %call10 = call zeroext i1 @ts_subtree_depends_on_column(ptr %4)
  %conv11 = zext i1 %call10 to i32
  %call12 = call i32 (ptr, ...) @printf(ptr noundef @.str.3, i32 noundef %conv11)
  %5 = load ptr, ptr %parent, align 8
  call void @free(ptr noundef %5) #5
  ret i32 0
}

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
  %size_bytes = getelementptr inbounds nuw %struct.SubtreeInlineData, ptr %retval, i32 0, i32 6
  store i8 1, ptr %size_bytes, align 1
  %coerce.dive = getelementptr inbounds nuw %union.Subtree, ptr %retval, i32 0, i32 0
  %1 = load ptr, ptr %coerce.dive, align 8
  ret ptr %1
}

; Function Attrs: noinline nounwind optnone uwtable
define internal ptr @make_heap_node(i1 noundef zeroext %depends) #0 {
entry:
  %retval = alloca %union.Subtree, align 8
  %depends.addr = alloca i8, align 1
  %h = alloca ptr, align 8
  %storedv = zext i1 %depends to i8
  store i8 %storedv, ptr %depends.addr, align 1
  %call = call noalias ptr @calloc(i64 noundef 1, i64 noundef 12) #6
  store ptr %call, ptr %h, align 8
  %0 = load ptr, ptr %h, align 8
  %ref_count = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %0, i32 0, i32 0
  store i32 1, ptr %ref_count, align 4
  %1 = load ptr, ptr %h, align 8
  %child_count = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %1, i32 0, i32 1
  store i32 2, ptr %child_count, align 4
  %2 = load ptr, ptr %h, align 8
  %symbol = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %2, i32 0, i32 2
  store i16 42, ptr %symbol, align 4
  %3 = load i8, ptr %depends.addr, align 1
  %loadedv = trunc i8 %3 to i1
  br i1 %loadedv, label %if.then, label %if.end

if.then:                                          ; preds = %entry
  %4 = load ptr, ptr %h, align 8
  %flags = getelementptr inbounds nuw %struct.SubtreeHeapData, ptr %4, i32 0, i32 3
  %5 = load i16, ptr %flags, align 2
  %conv = zext i16 %5 to i32
  %or = or i32 %conv, 256
  %conv1 = trunc i32 %or to i16
  store i16 %conv1, ptr %flags, align 2
  br label %if.end

if.end:                                           ; preds = %if.then, %entry
  %6 = load ptr, ptr %h, align 8
  store ptr %6, ptr %retval, align 8
  %coerce.dive = getelementptr inbounds nuw %union.Subtree, ptr %retval, i32 0, i32 0
  %7 = load ptr, ptr %coerce.dive, align 8
  ret ptr %7
}

declare i32 @printf(ptr noundef, ...) #1

; Function Attrs: noinline nounwind optnone uwtable
define internal zeroext i1 @ts_subtree_depends_on_column(ptr %self.coerce) #0 {
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

; Function Attrs: nounwind
declare void @free(ptr noundef) #2

; Function Attrs: nocallback nofree nounwind willreturn memory(argmem: write)
declare void @llvm.memset.p0.i64(ptr writeonly captures(none), i8, i64, i1 immarg) #3

; Function Attrs: nounwind allocsize(0,1)
declare noalias ptr @calloc(i64 noundef, i64 noundef) #4

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #2 = { nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #3 = { nocallback nofree nounwind willreturn memory(argmem: write) }
attributes #4 = { nounwind allocsize(0,1) "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #5 = { nounwind }
attributes #6 = { nounwind allocsize(0,1) }

!llvm.module.flags = !{!0, !1, !2, !3, !4}
!llvm.ident = !{!5}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 8, !"PIC Level", i32 2}
!2 = !{i32 7, !"PIE Level", i32 2}
!3 = !{i32 7, !"uwtable", i32 2}
!4 = !{i32 7, !"frame-pointer", i32 2}
!5 = !{!"clang version 22.1.8 (https://github.com/conda-forge/clangdev-feedstock 015bdba1263c0b3ebb3c518ff5947fbd99692bd0)"}
