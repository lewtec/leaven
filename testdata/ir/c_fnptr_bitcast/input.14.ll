; ModuleID = 'testdata/ir/fnptr_bitcast/source.c'
source_filename = "testdata/ir/fnptr_bitcast/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

%struct.Hooks = type { i8* ()*, void (i8*)*, i1 (i8*, %struct.TSLexer*, i8*)* }
%struct.TSLexer = type opaque
%struct.Payload = type { i32 }

@h = dso_local global %struct.Hooks { i8* ()* @create, void (i8*)* bitcast (void (%struct.Payload*)* @destroy to void (i8*)*), i1 (i8*, %struct.TSLexer*, i8*)* bitcast (i1 (%struct.Payload*, %struct.TSLexer*, i8*)* @scan to i1 (i8*, %struct.TSLexer*, i8*)*) }, align 8

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i8* @create() #0 {
entry:
  ret i8* null
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @destroy(%struct.Payload* noundef %p) #0 {
entry:
  %p.addr = alloca %struct.Payload*, align 8
  store %struct.Payload* %p, %struct.Payload** %p.addr, align 8
  %0 = load %struct.Payload*, %struct.Payload** %p.addr, align 8
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local zeroext i1 @scan(%struct.Payload* noundef %p, %struct.TSLexer* noundef %l, i8* noundef %s) #0 {
entry:
  %p.addr = alloca %struct.Payload*, align 8
  %l.addr = alloca %struct.TSLexer*, align 8
  %s.addr = alloca i8*, align 8
  store %struct.Payload* %p, %struct.Payload** %p.addr, align 8
  store %struct.TSLexer* %l, %struct.TSLexer** %l.addr, align 8
  store i8* %s, i8** %s.addr, align 8
  %0 = load %struct.Payload*, %struct.Payload** %p.addr, align 8
  %1 = load %struct.TSLexer*, %struct.TSLexer** %l.addr, align 8
  %2 = load i8*, i8** %s.addr, align 8
  ret i1 false
}

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }

!llvm.module.flags = !{!0, !1, !2}
!llvm.ident = !{!3}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 7, !"uwtable", i32 1}
!2 = !{i32 7, !"frame-pointer", i32 2}
!3 = !{!"clang version 14.0.6 (https://github.com/conda-forge/clangdev-feedstock ceeebe884c3cfd7160cf5a43e147f94439fafee3)"}
