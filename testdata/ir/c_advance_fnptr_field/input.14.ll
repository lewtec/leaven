; ModuleID = 'testdata/ir/advance_fnptr_field/source.c'
source_filename = "testdata/ir/advance_fnptr_field/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

%struct.TSLexer = type { i32, void (%struct.TSLexer*, i32)* }

@__const.main.lx = private unnamed_addr constant %struct.TSLexer { i32 0, void (%struct.TSLexer*, i32)* @real_advance }, align 8

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @scan(%struct.TSLexer* noundef %lexer) #0 {
entry:
  %lexer.addr = alloca %struct.TSLexer*, align 8
  store %struct.TSLexer* %lexer, %struct.TSLexer** %lexer.addr, align 8
  %0 = load %struct.TSLexer*, %struct.TSLexer** %lexer.addr, align 8
  call void @advance(%struct.TSLexer* noundef %0)
  %1 = load %struct.TSLexer*, %struct.TSLexer** %lexer.addr, align 8
  %lookahead = getelementptr inbounds %struct.TSLexer, %struct.TSLexer* %1, i32 0, i32 0
  %2 = load i32, i32* %lookahead, align 8
  ret i32 %2
}

; Function Attrs: noinline nounwind optnone uwtable
define internal void @advance(%struct.TSLexer* noundef %lexer) #0 {
entry:
  %lexer.addr = alloca %struct.TSLexer*, align 8
  store %struct.TSLexer* %lexer, %struct.TSLexer** %lexer.addr, align 8
  %0 = load %struct.TSLexer*, %struct.TSLexer** %lexer.addr, align 8
  %advance = getelementptr inbounds %struct.TSLexer, %struct.TSLexer* %0, i32 0, i32 1
  %1 = load void (%struct.TSLexer*, i32)*, void (%struct.TSLexer*, i32)** %advance, align 8
  %2 = load %struct.TSLexer*, %struct.TSLexer** %lexer.addr, align 8
  call void %1(%struct.TSLexer* noundef %2, i32 noundef 0)
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %lx = alloca %struct.TSLexer, align 8
  store i32 0, i32* %retval, align 4
  %0 = bitcast %struct.TSLexer* %lx to i8*
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 8 %0, i8* align 8 bitcast (%struct.TSLexer* @__const.main.lx to i8*), i64 16, i1 false)
  %call = call i32 @scan(%struct.TSLexer* noundef %lx)
  %cmp = icmp eq i32 %call, 1
  %1 = zext i1 %cmp to i64
  %cond = select i1 %cmp, i32 0, i32 1
  ret i32 %cond
}

; Function Attrs: noinline nounwind optnone uwtable
define internal void @real_advance(%struct.TSLexer* noundef %l, i32 noundef %skip) #0 {
entry:
  %l.addr = alloca %struct.TSLexer*, align 8
  %skip.addr = alloca i32, align 4
  store %struct.TSLexer* %l, %struct.TSLexer** %l.addr, align 8
  store i32 %skip, i32* %skip.addr, align 4
  %0 = load i32, i32* %skip.addr, align 4
  %1 = load %struct.TSLexer*, %struct.TSLexer** %l.addr, align 8
  %lookahead = getelementptr inbounds %struct.TSLexer, %struct.TSLexer* %1, i32 0, i32 0
  %2 = load i32, i32* %lookahead, align 8
  %inc = add nsw i32 %2, 1
  store i32 %inc, i32* %lookahead, align 8
  ret void
}

; Function Attrs: argmemonly nofree nounwind willreturn
declare void @llvm.memcpy.p0i8.p0i8.i64(i8* noalias nocapture writeonly, i8* noalias nocapture readonly, i64, i1 immarg) #1

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { argmemonly nofree nounwind willreturn }

!llvm.module.flags = !{!0, !1, !2}
!llvm.ident = !{!3}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 7, !"uwtable", i32 1}
!2 = !{i32 7, !"frame-pointer", i32 2}
!3 = !{!"clang version 14.0.6 (https://github.com/conda-forge/clangdev-feedstock ceeebe884c3cfd7160cf5a43e147f94439fafee3)"}
