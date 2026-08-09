; ModuleID = 'source.ddb5a00486c70246-cgu.0'
source_filename = "source.ddb5a00486c70246-cgu.0"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@alloc_32a1fd3b71af0e58a3b40a4f8c2e3027 = private unnamed_addr constant [10 x i8] c"source.rs\00", align 1
@alloc_3361e9e22ad822271a31aaae6159a1a5 = private unnamed_addr constant <{ ptr, [16 x i8] }> <{ ptr @alloc_32a1fd3b71af0e58a3b40a4f8c2e3027, [16 x i8] c"\09\00\00\00\00\00\00\00\13\00\00\00\09\00\00\00" }>, align 8
@alloc_767dbc325ab8b67804bea82b9d99bb69 = private unnamed_addr constant <{ ptr, [16 x i8] }> <{ ptr @alloc_32a1fd3b71af0e58a3b40a4f8c2e3027, [16 x i8] c"\09\00\00\00\00\00\00\00\0D\00\00\00\09\00\00\00" }>, align 8
@alloc_91a735d78e039264b4ebc0bd119a8c13 = private unnamed_addr constant <{ ptr, [16 x i8] }> <{ ptr @alloc_32a1fd3b71af0e58a3b40a4f8c2e3027, [16 x i8] c"\09\00\00\00\00\00\00\00*\00\00\00\0D\00\00\00" }>, align 8
@alloc_87551382a9de3243abbfdbda2f0b586b = private unnamed_addr constant [4 x i8] c"%d\0A\00", align 1

; source::cell_value
; Function Attrs: nounwind nonlazybind uwtable
define i32 @_RNvCsj29zdfijJxG_6source10cell_value(i32 %0, i32 %1) unnamed_addr #0 {
start:
  %_0 = alloca [4 x i8], align 4
  %c = alloca [8 x i8], align 4
  store i32 %0, ptr %c, align 4
  %2 = getelementptr inbounds i8, ptr %c, i64 4
  store i32 %1, ptr %2, align 4
  %3 = load i32, ptr %c, align 4
  %4 = getelementptr inbounds i8, ptr %c, i64 4
  %5 = load i32, ptr %4, align 4
  %_2 = zext i32 %3 to i64
  switch i64 %_2, label %bb1 [
    i64 0, label %bb4
    i64 1, label %bb3
    i64 2, label %bb2
  ]

bb1:                                              ; preds = %start
  unreachable

bb4:                                              ; preds = %start
  store i32 0, ptr %_0, align 4
  br label %bb5

bb3:                                              ; preds = %start
  store i32 -1, ptr %_0, align 4
  br label %bb5

bb2:                                              ; preds = %start
  %6 = getelementptr inbounds i8, ptr %c, i64 4
  %v = load i32, ptr %6, align 4
  store i32 %v, ptr %_0, align 4
  br label %bb5

bb5:                                              ; preds = %bb2, %bb3, %bb4
  %7 = load i32, ptr %_0, align 4
  ret i32 %7
}

; source::abs_i32
; Function Attrs: nounwind nonlazybind uwtable
define internal i32 @_RNvCsj29zdfijJxG_6source7abs_i32(i32 %x) unnamed_addr #0 {
start:
  %_0 = alloca [4 x i8], align 4
  %_2 = icmp slt i32 %x, 0
  br i1 %_2, label %bb1, label %bb3

bb3:                                              ; preds = %start
  store i32 %x, ptr %_0, align 4
  br label %bb4

bb1:                                              ; preds = %start
  %_3 = icmp eq i32 %x, -2147483648
  br i1 %_3, label %panic, label %bb2

bb4:                                              ; preds = %bb2, %bb3
  %0 = load i32, ptr %_0, align 4
  ret i32 %0

bb2:                                              ; preds = %bb1
  %1 = sub i32 0, %x
  store i32 %1, ptr %_0, align 4
  br label %bb4

panic:                                            ; preds = %bb1
; call core::panicking::panic_const::panic_const_neg_overflow
  call void @_RNvNtNtCs4NRVxsYgnAr_4core9panicking11panic_const24panic_const_neg_overflow(ptr align 8 @alloc_3361e9e22ad822271a31aaae6159a1a5) #3
  unreachable
}

; <source::Point>::manhattan
; Function Attrs: nounwind nonlazybind uwtable
define i32 @_RNvMCsj29zdfijJxG_6sourceNtB2_5Point9manhattan(ptr align 4 %self) unnamed_addr #0 {
start:
  %_3 = load i32, ptr %self, align 4
; call source::abs_i32
  %_2 = call i32 @_RNvCsj29zdfijJxG_6source7abs_i32(i32 %_3) #4
  %0 = getelementptr inbounds i8, ptr %self, i64 4
  %_5 = load i32, ptr %0, align 4
; call source::abs_i32
  %_4 = call i32 @_RNvCsj29zdfijJxG_6source7abs_i32(i32 %_5) #4
  %1 = call { i32, i1 } @llvm.sadd.with.overflow.i32(i32 %_2, i32 %_4)
  %_6.0 = extractvalue { i32, i1 } %1, 0
  %_6.1 = extractvalue { i32, i1 } %1, 1
  br i1 %_6.1, label %panic, label %bb3

bb3:                                              ; preds = %start
  ret i32 %_6.0

panic:                                            ; preds = %start
; call core::panicking::panic_const::panic_const_add_overflow
  call void @_RNvNtNtCs4NRVxsYgnAr_4core9panicking11panic_const24panic_const_add_overflow(ptr align 8 @alloc_767dbc325ab8b67804bea82b9d99bb69) #3
  unreachable
}

; Function Attrs: nounwind nonlazybind uwtable
define i32 @main() unnamed_addr #0 {
start:
  %p = alloca [8 x i8], align 4
  store i32 3, ptr %p, align 4
  %0 = getelementptr inbounds i8, ptr %p, i64 4
  store i32 -4, ptr %0, align 4
; call <source::Point>::manhattan
  %_3 = call i32 @_RNvMCsj29zdfijJxG_6sourceNtB2_5Point9manhattan(ptr align 4 %p) #4
; call source::cell_value
  %_5 = call i32 @_RNvCsj29zdfijJxG_6source10cell_value(i32 2, i32 0) #4
  %1 = call { i32, i1 } @llvm.sadd.with.overflow.i32(i32 %_3, i32 %_5)
  %_7.0 = extractvalue { i32, i1 } %1, 0
  %_7.1 = extractvalue { i32, i1 } %1, 1
  br i1 %_7.1, label %panic, label %bb3

bb3:                                              ; preds = %start
  %_8 = call i32 (ptr, ...) @printf(ptr @alloc_87551382a9de3243abbfdbda2f0b586b, i32 %_7.0) #4
  ret i32 0

panic:                                            ; preds = %start
; call core::panicking::panic_const::panic_const_add_overflow
  call void @_RNvNtNtCs4NRVxsYgnAr_4core9panicking11panic_const24panic_const_add_overflow(ptr align 8 @alloc_91a735d78e039264b4ebc0bd119a8c13) #3
  unreachable
}

; core::panicking::panic_const::panic_const_neg_overflow
; Function Attrs: cold noinline noreturn nounwind nonlazybind uwtable
declare void @_RNvNtNtCs4NRVxsYgnAr_4core9panicking11panic_const24panic_const_neg_overflow(ptr align 8) unnamed_addr #1

; Function Attrs: nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none)
declare { i32, i1 } @llvm.sadd.with.overflow.i32(i32, i32) #2

; core::panicking::panic_const::panic_const_add_overflow
; Function Attrs: cold noinline noreturn nounwind nonlazybind uwtable
declare void @_RNvNtNtCs4NRVxsYgnAr_4core9panicking11panic_const24panic_const_add_overflow(ptr align 8) unnamed_addr #1

; Function Attrs: nounwind nonlazybind uwtable
declare i32 @printf(ptr, ...) unnamed_addr #0

attributes #0 = { nounwind nonlazybind uwtable "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #1 = { cold noinline noreturn nounwind nonlazybind uwtable "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #2 = { nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none) }
attributes #3 = { noreturn nounwind }
attributes #4 = { nounwind }

!llvm.module.flags = !{!0, !1}
!llvm.ident = !{!2}

!0 = !{i32 8, !"PIC Level", i32 2}
!1 = !{i32 2, !"RtLibUseGOT", i32 1}
!2 = !{!"rustc version 1.97.1 (8bab26f4f 2026-07-14)"}
