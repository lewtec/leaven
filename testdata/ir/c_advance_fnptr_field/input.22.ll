; ModuleID = 'testdata/ir/c_advance_fnptr_field/source.c'
source_filename = "testdata/ir/c_advance_fnptr_field/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

%struct.TSLexer = type { i32, ptr }

@__const.main.lx = private unnamed_addr constant { i32, [4 x i8], ptr } { i32 0, [4 x i8] zeroinitializer, ptr @real_advance }, align 8

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @scan(ptr noundef %lexer) #0 {
entry:
  %lexer.addr = alloca ptr, align 8
  store ptr %lexer, ptr %lexer.addr, align 8
  %0 = load ptr, ptr %lexer.addr, align 8
  call void @advance(ptr noundef %0)
  %1 = load ptr, ptr %lexer.addr, align 8
  %lookahead = getelementptr inbounds nuw %struct.TSLexer, ptr %1, i32 0, i32 0
  %2 = load i32, ptr %lookahead, align 8
  ret i32 %2
}

; Function Attrs: noinline nounwind optnone uwtable
define internal void @advance(ptr noundef %lexer) #0 {
entry:
  %lexer.addr = alloca ptr, align 8
  store ptr %lexer, ptr %lexer.addr, align 8
  %0 = load ptr, ptr %lexer.addr, align 8
  %advance = getelementptr inbounds nuw %struct.TSLexer, ptr %0, i32 0, i32 1
  %1 = load ptr, ptr %advance, align 8
  %2 = load ptr, ptr %lexer.addr, align 8
  call void %1(ptr noundef %2, i32 noundef 0)
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %lx = alloca %struct.TSLexer, align 8
  store i32 0, ptr %retval, align 4
  call void @llvm.memcpy.p0.p0.i64(ptr align 8 %lx, ptr align 8 @__const.main.lx, i64 16, i1 false)
  %call = call i32 @scan(ptr noundef %lx)
  %cmp = icmp eq i32 %call, 1
  %0 = zext i1 %cmp to i64
  %cond = select i1 %cmp, i32 0, i32 1
  ret i32 %cond
}

; Function Attrs: noinline nounwind optnone uwtable
define internal void @real_advance(ptr noundef %l, i32 noundef %skip) #0 {
entry:
  %l.addr = alloca ptr, align 8
  %skip.addr = alloca i32, align 4
  store ptr %l, ptr %l.addr, align 8
  store i32 %skip, ptr %skip.addr, align 4
  %0 = load i32, ptr %skip.addr, align 4
  %1 = load ptr, ptr %l.addr, align 8
  %lookahead = getelementptr inbounds nuw %struct.TSLexer, ptr %1, i32 0, i32 0
  %2 = load i32, ptr %lookahead, align 8
  %inc = add nsw i32 %2, 1
  store i32 %inc, ptr %lookahead, align 8
  ret void
}

; Function Attrs: nocallback nofree nounwind willreturn memory(argmem: readwrite)
declare void @llvm.memcpy.p0.p0.i64(ptr noalias writeonly captures(none), ptr noalias readonly captures(none), i64, i1 immarg) #1

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { nocallback nofree nounwind willreturn memory(argmem: readwrite) }

!llvm.module.flags = !{!0, !1, !2, !3, !4}
!llvm.ident = !{!5}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 8, !"PIC Level", i32 2}
!2 = !{i32 7, !"PIE Level", i32 2}
!3 = !{i32 7, !"uwtable", i32 2}
!4 = !{i32 7, !"frame-pointer", i32 2}
!5 = !{!"clang version 22.1.8 (https://github.com/conda-forge/clangdev-feedstock 015bdba1263c0b3ebb3c518ff5947fbd99692bd0)"}
