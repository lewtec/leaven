; ModuleID = 'source.ddb5a00486c70246-cgu.0'
source_filename = "source.ddb5a00486c70246-cgu.0"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@alloc_fae6be5744a40d7485552c134bd2e362 = private unnamed_addr constant [31 x i8] c"testdata/ir/rust_add/source.rs\00", align 1
@alloc_9a73f36d4b9165f707b56159f9b1052a = private unnamed_addr constant <{ ptr, [16 x i8] }> <{ ptr @alloc_fae6be5744a40d7485552c134bd2e362, [16 x i8] c"\1E\00\00\00\00\00\00\00\03\00\00\00\05\00\00\00" }>, align 8

; source::add
; Function Attrs: nonlazybind uwtable
define i32 @_RNvCsj29zdfijJxG_6source3add(i32 %a, i32 %b) unnamed_addr #0 {
start:
  %0 = call { i32, i1 } @llvm.sadd.with.overflow.i32(i32 %a, i32 %b)
  %_3.0 = extractvalue { i32, i1 } %0, 0
  %_3.1 = extractvalue { i32, i1 } %0, 1
  br i1 %_3.1, label %panic, label %bb1

bb1:                                              ; preds = %start
  ret i32 %_3.0

panic:                                            ; preds = %start
; call core::panicking::panic_const::panic_const_add_overflow
  call void @_RNvNtNtCs4NRVxsYgnAr_4core9panicking11panic_const24panic_const_add_overflow(ptr align 8 @alloc_9a73f36d4b9165f707b56159f9b1052a) #3
  unreachable
}

; Function Attrs: nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none)
declare { i32, i1 } @llvm.sadd.with.overflow.i32(i32, i32) #1

; core::panicking::panic_const::panic_const_add_overflow
; Function Attrs: cold noinline noreturn nonlazybind uwtable
declare void @_RNvNtNtCs4NRVxsYgnAr_4core9panicking11panic_const24panic_const_add_overflow(ptr align 8) unnamed_addr #2

attributes #0 = { nonlazybind uwtable "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #1 = { nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none) }
attributes #2 = { cold noinline noreturn nonlazybind uwtable "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #3 = { noreturn }

!llvm.module.flags = !{!0, !1}
!llvm.ident = !{!2}

!0 = !{i32 8, !"PIC Level", i32 2}
!1 = !{i32 2, !"RtLibUseGOT", i32 1}
!2 = !{!"rustc version 1.97.1 (8bab26f4f 2026-07-14)"}
