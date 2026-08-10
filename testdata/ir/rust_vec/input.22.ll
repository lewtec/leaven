; ModuleID = 'source.ddb5a00486c70246-cgu.0'
source_filename = "source.ddb5a00486c70246-cgu.0"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@alloc_61247b90e1706a3f65e71312b599d3d1 = private unnamed_addr constant [4 x i8] c"\C0\01\0A\00", align 1

; <alloc::raw_vec::RawVec<source::Point>>::grow_one
; Function Attrs: cold noinline nounwind nonlazybind uwtable
define void @_RNvMs3_NtCscdodAO9FK5_5alloc7raw_vecINtB5_6RawVecNtCsj29zdfijJxG_6source5PointE8grow_oneBN_(ptr noalias noundef align 8 captures(none) dereferenceable(16) %self) unnamed_addr #0 {
start:
  %self1 = load i64, ptr %self, align 8, !range !3, !noundef !4
; call <alloc::raw_vec::RawVecInner>::grow_amortized
  %0 = tail call fastcc { i64, i64 } @_RNvMs4_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14grow_amortizedCsj29zdfijJxG_6source(ptr noalias noundef align 8 dereferenceable(16) %self, i64 noundef %self1) #12
  %1 = extractvalue { i64, i64 } %0, 0
  %.not = icmp eq i64 %1, -1
  br i1 %.not, label %bb3, label %bb2, !prof !5

bb2:                                              ; preds = %start
  %2 = extractvalue { i64, i64 } %0, 1
; call alloc::raw_vec::handle_error
  tail call void @_RNvNtCscdodAO9FK5_5alloc7raw_vec12handle_error(i64 noundef %1, i64 %2) #13
  unreachable

bb3:                                              ; preds = %start
  ret void
}

; <alloc::raw_vec::RawVecInner>::finish_grow
; Function Attrs: cold nounwind nonlazybind uwtable
define internal fastcc void @_RNvMs4_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner11finish_growCsj29zdfijJxG_6source(ptr dead_on_unwind noalias noundef nonnull writable writeonly sret([24 x i8]) align 8 captures(none) dereferenceable(24) initializes((0, 8)) %_0, ptr noalias noundef nonnull readonly align 8 captures(none) dereferenceable(16) %self, i64 noundef range(i64 0, -1) %cap) unnamed_addr #1 {
start:
  %_34.0 = shl i64 %cap, 3
  %_34.1 = icmp ult i64 %cap, 2305843009213693952
  %_39 = icmp ult i64 %_34.0, 9223372036854775805
  %or.cond = and i1 %_34.1, %_39
  br i1 %or.cond, label %bb15, label %bb11, !prof !6

bb15:                                             ; preds = %start
  %self1.i = load i64, ptr %self, align 8, !range !3, !alias.scope !7, !noalias !10, !noundef !4
  %0 = icmp eq i64 %self1.i, 0
  br i1 %0, label %bb5, label %bb3

bb3:                                              ; preds = %bb15
  %alloc_size.i = shl nuw i64 %self1.i, 3
  %1 = getelementptr inbounds nuw i8, ptr %self, i64 8
  %self3.i = load ptr, ptr %1, align 8, !alias.scope !7, !noalias !10, !nonnull !4, !noundef !4
  %cond.i.i = icmp uge i64 %_34.0, %alloc_size.i
  tail call void @llvm.assume(i1 %cond.i.i)
; call __rustc::__rust_realloc
  %raw_ptr.i.i = tail call noundef align 4 ptr @_RNvCs9wFQrvczXsK_7___rustc14___rust_realloc(ptr noundef nonnull %self3.i, i64 noundef %alloc_size.i, i64 noundef 4, i64 noundef range(i64 0, 9223372036854775805) %_34.0) #12
  br label %bb7

bb5:                                              ; preds = %bb15
  %2 = icmp eq i64 %_34.0, 0
  br i1 %2, label %bb9, label %bb4.i.i

bb4.i.i:                                          ; preds = %bb5
; call __rustc::__rust_no_alloc_shim_is_unstable_v2
  tail call void @_RNvCs9wFQrvczXsK_7___rustc35___rust_no_alloc_shim_is_unstable_v2() #12
; call __rustc::__rust_alloc
  %3 = tail call noundef align 4 ptr @_RNvCs9wFQrvczXsK_7___rustc12___rust_alloc(i64 noundef range(i64 0, 9223372036854775805) %_34.0, i64 noundef 4) #12
  br label %bb7

bb7:                                              ; preds = %bb4.i.i, %bb3
  %raw_ptr.i.i.pn = phi ptr [ %raw_ptr.i.i, %bb3 ], [ %3, %bb4.i.i ]
  %4 = icmp eq ptr %raw_ptr.i.i.pn, null
  br i1 %4, label %bb8, label %bb9

bb8:                                              ; preds = %bb7
  %5 = getelementptr inbounds nuw i8, ptr %_0, i64 8
  store i64 4, ptr %5, align 8
  br label %bb11

bb9:                                              ; preds = %bb5, %bb7
  %raw_ptr.i.i.pn18 = phi ptr [ %raw_ptr.i.i.pn, %bb7 ], [ inttoptr (i64 4 to ptr), %bb5 ]
  %6 = getelementptr inbounds nuw i8, ptr %_0, i64 8
  store ptr %raw_ptr.i.i.pn18, ptr %6, align 8
  br label %bb11

bb11:                                             ; preds = %start, %bb9, %bb8
  %.sink19 = phi i64 [ 16, %bb9 ], [ 16, %bb8 ], [ 8, %start ]
  %_34.0.sink = phi i64 [ %_34.0, %bb9 ], [ %_34.0, %bb8 ], [ 0, %start ]
  %storemerge9 = phi i64 [ 0, %bb9 ], [ 1, %bb8 ], [ 1, %start ]
  %7 = getelementptr inbounds nuw i8, ptr %_0, i64 %.sink19
  store i64 %_34.0.sink, ptr %7, align 8
  store i64 %storemerge9, ptr %_0, align 8
  ret void
}

; <alloc::raw_vec::RawVecInner>::grow_amortized
; Function Attrs: cold nounwind nonlazybind uwtable
define internal fastcc { i64, i64 } @_RNvMs4_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14grow_amortizedCsj29zdfijJxG_6source(ptr noalias noundef nonnull align 8 captures(none) dereferenceable(16) %self, i64 noundef range(i64 0, -9223372036854775808) %len) unnamed_addr #1 personality ptr @rust_eh_personality {
start:
  %self3 = alloca [24 x i8], align 8
  %_25.0 = add nuw i64 %len, 1
  %self5 = load i64, ptr %self, align 8, !range !3, !noundef !4
  %v16 = shl nuw i64 %self5, 1
  %..i = tail call noundef range(i64 0, -1) i64 @llvm.umax.i64(i64 range(i64 0, -1) %_25.0, i64 range(i64 0, -1) %v16)
  %..i15 = tail call noundef range(i64 0, -1) i64 @llvm.umax.i64(i64 range(i64 0, -1) %..i, i64 4)
  call void @llvm.lifetime.start.p0(ptr nonnull %self3)
; call <alloc::raw_vec::RawVecInner>::finish_grow
  call fastcc void @_RNvMs4_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner11finish_growCsj29zdfijJxG_6source(ptr noalias noundef sret([24 x i8]) align 8 captures(none) dereferenceable(24) %self3, ptr noalias noundef readonly align 8 captures(address, read_provenance) dereferenceable(16) %self, i64 noundef %..i15) #12
  %_37 = load i64, ptr %self3, align 8, !range !12, !noundef !4
  %0 = trunc nuw i64 %_37 to i1
  %1 = getelementptr inbounds nuw i8, ptr %self3, i64 8
  br i1 %0, label %bb18, label %bb19

bb6:                                              ; preds = %bb18, %bb19
  %_0.sroa.5.0 = phi i64 [ undef, %bb19 ], [ %e.1, %bb18 ]
  %_0.sroa.0.0 = phi i64 [ -1, %bb19 ], [ %e.0, %bb18 ]
  %2 = insertvalue { i64, i64 } poison, i64 %_0.sroa.0.0, 0
  %3 = insertvalue { i64, i64 } %2, i64 %_0.sroa.5.0, 1
  ret { i64, i64 } %3

bb18:                                             ; preds = %start
  %e.0 = load i64, ptr %1, align 8, !range !13, !noundef !4
  %4 = getelementptr inbounds nuw i8, ptr %self3, i64 16
  %e.1 = load i64, ptr %4, align 8
  call void @llvm.lifetime.end.p0(ptr nonnull %self3)
  br label %bb6

bb19:                                             ; preds = %start
  %v.0 = load ptr, ptr %1, align 8, !nonnull !4, !noundef !4
  call void @llvm.lifetime.end.p0(ptr nonnull %self3)
  %5 = getelementptr inbounds nuw i8, ptr %self, i64 8
  store ptr %v.0, ptr %5, align 8
  %6 = icmp sgt i64 %..i15, -1
  tail call void @llvm.assume(i1 %6)
  store i64 %..i15, ptr %self, align 8
  br label %bb6
}

; Function Attrs: cold nounwind nonlazybind uwtable
define noundef i32 @main() unnamed_addr #1 personality ptr @rust_eh_personality {
_RNvMsF_NtCscdodAO9FK5_5alloc3vecINtB5_3VecNtCsj29zdfijJxG_6source5PointE8push_mutBG_.exit:
  %args = alloca [16 x i8], align 8
  %n = alloca [4 x i8], align 4
  %v = alloca [24 x i8], align 8
  call void @llvm.lifetime.start.p0(ptr nonnull %v)
  store i64 0, ptr %v, align 8
  %0 = getelementptr inbounds nuw i8, ptr %v, i64 8
  store ptr inttoptr (i64 4 to ptr), ptr %0, align 8
  %1 = getelementptr inbounds nuw i8, ptr %v, i64 16
  store i64 0, ptr %1, align 8
  tail call void @llvm.experimental.noalias.scope.decl(metadata !14)
; call <alloc::raw_vec::RawVec<source::Point>>::grow_one
  call void @_RNvMs3_NtCscdodAO9FK5_5alloc7raw_vecINtB5_6RawVecNtCsj29zdfijJxG_6source5PointE8grow_oneBN_(ptr noalias noundef nonnull align 8 dereferenceable(24) %v) #12
  %_14.i = load ptr, ptr %0, align 8, !alias.scope !14, !nonnull !4, !noundef !4
  store i32 3, ptr %_14.i, align 4, !noalias !14
  %2 = getelementptr inbounds nuw i8, ptr %_14.i, i64 4
  store i32 -4, ptr %2, align 4, !noalias !14
  store i64 1, ptr %1, align 8, !alias.scope !14
  tail call void @llvm.experimental.noalias.scope.decl(metadata !17)
  %self1.i2 = load i64, ptr %v, align 8, !range !3, !alias.scope !17, !noundef !4
  %_4.i3 = icmp eq i64 %self1.i2, 1
  br i1 %_4.i3, label %bb1.i6, label %bb5.i

bb1.i6:                                           ; preds = %_RNvMsF_NtCscdodAO9FK5_5alloc3vecINtB5_3VecNtCsj29zdfijJxG_6source5PointE8push_mutBG_.exit
; call <alloc::raw_vec::RawVec<source::Point>>::grow_one
  call void @_RNvMs3_NtCscdodAO9FK5_5alloc7raw_vecINtB5_6RawVecNtCsj29zdfijJxG_6source5PointE8grow_oneBN_(ptr noalias noundef nonnull align 8 dereferenceable(24) %v) #12
  %_14.i4.pre = load ptr, ptr %0, align 8, !alias.scope !17
  %_4.i.i.i.pre = load i32, ptr %_14.i4.pre, align 4, !alias.scope !20
  %.phi.trans.insert = getelementptr inbounds nuw i8, ptr %_14.i4.pre, i64 4
  %_6.i.i.i.pre = load i32, ptr %.phi.trans.insert, align 4, !alias.scope !20
  %self1.i.i.i.i.i.pre = load i64, ptr %v, align 8, !range !3, !alias.scope !25, !noalias !36
  %3 = tail call i32 @llvm.abs.i32(i32 %_4.i.i.i.pre, i1 false)
  %4 = tail call i32 @llvm.abs.i32(i32 %_6.i.i.i.pre, i1 false)
  %5 = add i32 %3, %4
  %6 = add i32 %5, 3
  br label %bb5.i

bb5.i:                                            ; preds = %bb1.i6, %_RNvMsF_NtCscdodAO9FK5_5alloc3vecINtB5_3VecNtCsj29zdfijJxG_6source5PointE8push_mutBG_.exit
  %self1.i.i.i.i.i = phi i64 [ %self1.i.i.i.i.i.pre, %bb1.i6 ], [ %self1.i2, %_RNvMsF_NtCscdodAO9FK5_5alloc3vecINtB5_3VecNtCsj29zdfijJxG_6source5PointE8push_mutBG_.exit ]
  %_14.i4 = phi ptr [ %_14.i4.pre, %bb1.i6 ], [ %_14.i, %_RNvMsF_NtCscdodAO9FK5_5alloc3vecINtB5_3VecNtCsj29zdfijJxG_6source5PointE8push_mutBG_.exit ]
  %_4.0.i.i.i = phi i32 [ %6, %bb1.i6 ], [ 10, %_RNvMsF_NtCscdodAO9FK5_5alloc3vecINtB5_3VecNtCsj29zdfijJxG_6source5PointE8push_mutBG_.exit ]
  %end.i5 = getelementptr inbounds nuw i8, ptr %_14.i4, i64 8
  store i32 1, ptr %end.i5, align 4, !noalias !17
  %7 = getelementptr inbounds nuw i8, ptr %_14.i4, i64 12
  store i32 2, ptr %7, align 4, !noalias !17
  call void @llvm.lifetime.start.p0(ptr nonnull %n)
  store i32 %_4.0.i.i.i, ptr %n, align 4
  call void @llvm.lifetime.start.p0(ptr nonnull %args)
  store ptr %n, ptr %args, align 8
  %_10.sroa.4.0..sroa_idx = getelementptr inbounds nuw i8, ptr %args, i64 8
  store ptr @_RNvXs9_NtNtNtCs4NRVxsYgnAr_4core3fmt3num3implNtB9_7Display3fmt, ptr %_10.sroa.4.0..sroa_idx, align 8
; call std::io::stdio::_print
  call void @_RNvNtNtCs2AWtUsOyxgP_3std2io5stdio6__print(ptr noundef nonnull @alloc_61247b90e1706a3f65e71312b599d3d1, ptr noundef nonnull %args) #12
  call void @llvm.lifetime.end.p0(ptr nonnull %args)
  call void @llvm.lifetime.end.p0(ptr nonnull %n)
  call void @llvm.experimental.noalias.scope.decl(metadata !38)
  call void @llvm.experimental.noalias.scope.decl(metadata !39)
  call void @llvm.experimental.noalias.scope.decl(metadata !40)
  call void @llvm.experimental.noalias.scope.decl(metadata !41)
  %8 = icmp eq i64 %self1.i.i.i.i.i, 0
  br i1 %8, label %_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueINtNtCscdodAO9FK5_5alloc3vec3VecNtCsj29zdfijJxG_6source5PointEEB19_.exit, label %bb2.i.i.i.i

bb2.i.i.i.i:                                      ; preds = %bb5.i
  %alloc_size.i.i.i.i.i = shl nuw i64 %self1.i.i.i.i.i, 3
; call __rustc::__rust_dealloc
  call void @_RNvCs9wFQrvczXsK_7___rustc14___rust_dealloc(ptr noundef nonnull %_14.i4, i64 noundef %alloc_size.i.i.i.i.i, i64 noundef range(i64 1, -9223372036854775807) 4) #12, !noalias !42
  br label %_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueINtNtCscdodAO9FK5_5alloc3vec3VecNtCsj29zdfijJxG_6source5PointEEB19_.exit

_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueINtNtCscdodAO9FK5_5alloc3vec3VecNtCsj29zdfijJxG_6source5PointEEB19_.exit: ; preds = %bb5.i, %bb2.i.i.i.i
  call void @llvm.lifetime.end.p0(ptr nonnull %v)
  ret i32 0
}

; Function Attrs: nounwind nonlazybind uwtable
declare noundef range(i32 0, 10) i32 @rust_eh_personality(i32 noundef, i32 noundef, i64 noundef, ptr noundef, ptr noundef) unnamed_addr #2

; Function Attrs: mustprogress nocallback nofree nosync nounwind willreturn memory(argmem: readwrite)
declare void @llvm.lifetime.start.p0(ptr captures(none)) #3

; Function Attrs: mustprogress nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write)
declare void @llvm.assume(i1 noundef) #4

; Function Attrs: mustprogress nocallback nofree nosync nounwind willreturn memory(argmem: readwrite)
declare void @llvm.lifetime.end.p0(ptr captures(none)) #3

; __rustc::__rust_dealloc
; Function Attrs: nounwind nonlazybind allockind("free") uwtable
declare void @_RNvCs9wFQrvczXsK_7___rustc14___rust_dealloc(ptr allocptr noundef nonnull captures(address), i64 noundef, i64 noundef range(i64 1, -9223372036854775807)) unnamed_addr #5

; __rustc::__rust_realloc
; Function Attrs: nounwind nonlazybind allockind("realloc,aligned") allocsize(3) uwtable
declare noalias noundef ptr @_RNvCs9wFQrvczXsK_7___rustc14___rust_realloc(ptr allocptr noundef nonnull, i64 noundef, i64 allocalign noundef range(i64 1, -9223372036854775807), i64 noundef) unnamed_addr #6

; __rustc::__rust_no_alloc_shim_is_unstable_v2
; Function Attrs: nounwind nonlazybind uwtable
declare void @_RNvCs9wFQrvczXsK_7___rustc35___rust_no_alloc_shim_is_unstable_v2() unnamed_addr #2

; __rustc::__rust_alloc
; Function Attrs: nounwind nonlazybind allockind("alloc,uninitialized,aligned") allocsize(0) uwtable
declare noalias noundef ptr @_RNvCs9wFQrvczXsK_7___rustc12___rust_alloc(i64 noundef, i64 allocalign noundef range(i64 1, -9223372036854775807)) unnamed_addr #7

; alloc::raw_vec::handle_error
; Function Attrs: cold minsize noreturn nounwind nonlazybind optsize uwtable
declare void @_RNvNtCscdodAO9FK5_5alloc7raw_vec12handle_error(i64 noundef range(i64 0, -9223372036854775807), i64) unnamed_addr #8

; <i32 as core::fmt::Display>::fmt
; Function Attrs: nounwind nonlazybind uwtable
declare noundef zeroext i1 @_RNvXs9_NtNtNtCs4NRVxsYgnAr_4core3fmt3num3implNtB9_7Display3fmt(ptr noalias noundef readonly align 4 captures(address, read_provenance) dereferenceable(4), ptr noalias noundef align 8 dereferenceable(24)) unnamed_addr #2

; std::io::stdio::_print
; Function Attrs: nounwind nonlazybind uwtable
declare void @_RNvNtNtCs2AWtUsOyxgP_3std2io5stdio6__print(ptr noundef nonnull, ptr noundef nonnull) unnamed_addr #2

; Function Attrs: nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: readwrite)
declare void @llvm.experimental.noalias.scope.decl(metadata) #9

; Function Attrs: nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none)
declare i64 @llvm.umax.i64(i64, i64) #10

; Function Attrs: nocallback nofree nosync nounwind speculatable willreturn memory(none)
declare i32 @llvm.abs.i32(i32, i1 immarg) #11

attributes #0 = { cold noinline nounwind nonlazybind uwtable "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #1 = { cold nounwind nonlazybind uwtable "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #2 = { nounwind nonlazybind uwtable "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #3 = { mustprogress nocallback nofree nosync nounwind willreturn memory(argmem: readwrite) }
attributes #4 = { mustprogress nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write) }
attributes #5 = { nounwind nonlazybind allockind("free") uwtable "alloc-family"="__rust_alloc" "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #6 = { nounwind nonlazybind allockind("realloc,aligned") allocsize(3) uwtable "alloc-family"="__rust_alloc" "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #7 = { nounwind nonlazybind allockind("alloc,uninitialized,aligned") allocsize(0) uwtable "alloc-family"="__rust_alloc" "alloc-variant-zeroed"="_RNvCs9wFQrvczXsK_7___rustc19___rust_alloc_zeroed" "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #8 = { cold minsize noreturn nounwind nonlazybind optsize uwtable "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #9 = { nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: readwrite) }
attributes #10 = { nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none) }
attributes #11 = { nocallback nofree nosync nounwind speculatable willreturn memory(none) }
attributes #12 = { nounwind }
attributes #13 = { noreturn nounwind }

!llvm.module.flags = !{!0, !1}
!llvm.ident = !{!2}

!0 = !{i32 8, !"PIC Level", i32 2}
!1 = !{i32 2, !"RtLibUseGOT", i32 1}
!2 = !{!"rustc version 1.97.1 (8bab26f4f 2026-07-14)"}
!3 = !{i64 0, i64 -9223372036854775808}
!4 = !{}
!5 = !{!"branch_weights", !"expected", i32 2000, i32 1}
!6 = !{!"branch_weights", i32 2000, i32 2002}
!7 = !{!8}
!8 = distinct !{!8, !9, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14current_memoryCsj29zdfijJxG_6source: %self"}
!9 = distinct !{!9, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14current_memoryCsj29zdfijJxG_6source"}
!10 = !{!11}
!11 = distinct !{!11, !9, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14current_memoryCsj29zdfijJxG_6source: %_0"}
!12 = !{i64 0, i64 2}
!13 = !{i64 0, i64 -9223372036854775807}
!14 = !{!15}
!15 = distinct !{!15, !16, !"_RNvMsF_NtCscdodAO9FK5_5alloc3vecINtB5_3VecNtCsj29zdfijJxG_6source5PointE8push_mutBG_: %self"}
!16 = distinct !{!16, !"_RNvMsF_NtCscdodAO9FK5_5alloc3vecINtB5_3VecNtCsj29zdfijJxG_6source5PointE8push_mutBG_"}
!17 = !{!18}
!18 = distinct !{!18, !19, !"_RNvMsF_NtCscdodAO9FK5_5alloc3vecINtB5_3VecNtCsj29zdfijJxG_6source5PointE8push_mutBG_: %self"}
!19 = distinct !{!19, !"_RNvMsF_NtCscdodAO9FK5_5alloc3vecINtB5_3VecNtCsj29zdfijJxG_6source5PointE8push_mutBG_"}
!20 = !{!21, !23}
!21 = distinct !{!21, !22, !"_RNCNvCsj29zdfijJxG_6source4main0B3_: %p"}
!22 = distinct !{!22, !"_RNCNvCsj29zdfijJxG_6source4main0B3_"}
!23 = distinct !{!23, !24, !"_RNCINvNtNtNtCs4NRVxsYgnAr_4core4iter8adapters3map8map_foldRNtCsj29zdfijJxG_6source5PointllNCNvBX_4main0NCINvXsa_NtNtB8_6traits5accumlNtB1M_3Sum3sumINtB4_3MapINtNtNtBa_5slice4iter4IterBV_EB1q_EE0E0BX_: %elt"}
!24 = distinct !{!24, !"_RNCINvNtNtNtCs4NRVxsYgnAr_4core4iter8adapters3map8map_foldRNtCsj29zdfijJxG_6source5PointllNCNvBX_4main0NCINvXsa_NtNtB8_6traits5accumlNtB1M_3Sum3sumINtB4_3MapINtNtNtBa_5slice4iter4IterBV_EB1q_EE0E0BX_"}
!25 = !{!26, !28, !30, !32, !34}
!26 = distinct !{!26, !27, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14current_memoryCsj29zdfijJxG_6source: %self"}
!27 = distinct !{!27, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14current_memoryCsj29zdfijJxG_6source"}
!28 = distinct !{!28, !29, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner10deallocateCsj29zdfijJxG_6source: %self"}
!29 = distinct !{!29, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner10deallocateCsj29zdfijJxG_6source"}
!30 = distinct !{!30, !31, !"_RNvXs1_NtCscdodAO9FK5_5alloc7raw_vecINtB5_6RawVecNtCsj29zdfijJxG_6source5PointENtNtNtCs4NRVxsYgnAr_4core3ops4drop4Drop4dropBN_: %self"}
!31 = distinct !{!31, !"_RNvXs1_NtCscdodAO9FK5_5alloc7raw_vecINtB5_6RawVecNtCsj29zdfijJxG_6source5PointENtNtNtCs4NRVxsYgnAr_4core3ops4drop4Drop4dropBN_"}
!32 = distinct !{!32, !33, !"_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueINtNtCscdodAO9FK5_5alloc7raw_vec6RawVecNtCsj29zdfijJxG_6source5PointEEB1g_: %_1"}
!33 = distinct !{!33, !"_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueINtNtCscdodAO9FK5_5alloc7raw_vec6RawVecNtCsj29zdfijJxG_6source5PointEEB1g_"}
!34 = distinct !{!34, !35, !"_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueINtNtCscdodAO9FK5_5alloc3vec3VecNtCsj29zdfijJxG_6source5PointEEB19_: %_1"}
!35 = distinct !{!35, !"_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueINtNtCscdodAO9FK5_5alloc3vec3VecNtCsj29zdfijJxG_6source5PointEEB19_"}
!36 = !{!37}
!37 = distinct !{!37, !27, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14current_memoryCsj29zdfijJxG_6source: %_0"}
!38 = !{!34}
!39 = !{!32}
!40 = !{!30}
!41 = !{!28}
!42 = !{!28, !30, !32, !34}
