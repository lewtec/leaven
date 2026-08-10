; ModuleID = 'source.ddb5a00486c70246-cgu.0'
source_filename = "source.ddb5a00486c70246-cgu.0"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@alloc_3aef7c697f52c5a3e9548cf3d170303c = private unnamed_addr constant [3 x i8] c"foo", align 1
@alloc_d4d3cbd0d4712cfcd6657cf295c08ada = private unnamed_addr constant [3 x i8] c"bar", align 1
@alloc_61247b90e1706a3f65e71312b599d3d1 = private unnamed_addr constant [4 x i8] c"\C0\01\0A\00", align 1

; <alloc::raw_vec::RawVecInner<_>>::reserve::do_reserve_and_handle::<alloc::alloc::Global>
; Function Attrs: cold nounwind nonlazybind uwtable
define internal fastcc void @_RINvNvMs2_NtCscdodAO9FK5_5alloc7raw_vecINtB8_11RawVecInnerpE7reserve21do_reserve_and_handleNtNtBa_5alloc6GlobalECsj29zdfijJxG_6source(ptr noalias noundef nonnull align 8 captures(none) dereferenceable(16) %slf, i64 noundef %len) unnamed_addr #0 personality ptr @rust_eh_personality {
start:
  %self3.i = alloca [24 x i8], align 8
  tail call void @llvm.experimental.noalias.scope.decl(metadata !3)
  %_25.1.i = icmp ugt i64 %len, -4
  br i1 %_25.1.i, label %bb2, label %bb9.i

bb9.i:                                            ; preds = %start
  %_25.0.i = add nuw i64 %len, 3
  %self5.i = load i64, ptr %slf, align 8, !range !6, !alias.scope !3, !noundef !7
  %v16.i = shl nuw i64 %self5.i, 1
  %..i.i = tail call noundef i64 @llvm.umax.i64(i64 %_25.0.i, i64 range(i64 0, -1) %v16.i)
  %..i15.i = tail call noundef i64 @llvm.umax.i64(i64 %..i.i, i64 8)
  call void @llvm.lifetime.start.p0(ptr nonnull %self3.i), !noalias !3
; call <alloc::raw_vec::RawVecInner>::finish_grow
  call fastcc void @_RNvMs4_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner11finish_growCsj29zdfijJxG_6source(ptr noalias noundef sret([24 x i8]) align 8 captures(none) dereferenceable(24) %self3.i, ptr noalias noundef nonnull readonly align 8 captures(address, read_provenance) dereferenceable(16) %slf, i64 noundef %..i15.i) #11
  %_37.i = load i64, ptr %self3.i, align 8, !range !8, !noalias !3, !noundef !7
  %0 = trunc nuw i64 %_37.i to i1
  %1 = getelementptr inbounds nuw i8, ptr %self3.i, i64 8
  br i1 %0, label %bb18.i, label %bb3

bb18.i:                                           ; preds = %bb9.i
  %e.0.i = load i64, ptr %1, align 8, !range !9, !noalias !3, !noundef !7
  %2 = getelementptr inbounds nuw i8, ptr %self3.i, i64 16
  %e.1.i = load i64, ptr %2, align 8, !noalias !3
  call void @llvm.lifetime.end.p0(ptr nonnull %self3.i), !noalias !3
  br label %bb2

bb2:                                              ; preds = %bb18.i, %start
  %_0.sroa.5.0.i.ph = phi i64 [ undef, %start ], [ %e.1.i, %bb18.i ]
  %_0.sroa.0.0.i.ph = phi i64 [ 0, %start ], [ %e.0.i, %bb18.i ]
; call alloc::raw_vec::handle_error
  tail call void @_RNvNtCscdodAO9FK5_5alloc7raw_vec12handle_error(i64 noundef %_0.sroa.0.0.i.ph, i64 %_0.sroa.5.0.i.ph) #12
  unreachable

bb3:                                              ; preds = %bb9.i
  %v.0.i = load ptr, ptr %1, align 8, !noalias !3, !nonnull !7, !noundef !7
  call void @llvm.lifetime.end.p0(ptr nonnull %self3.i), !noalias !3
  %3 = getelementptr inbounds nuw i8, ptr %slf, i64 8
  store ptr %v.0.i, ptr %3, align 8, !alias.scope !3
  %4 = icmp sgt i64 %..i15.i, -1
  tail call void @llvm.assume(i1 %4)
  store i64 %..i15.i, ptr %slf, align 8, !alias.scope !3
  ret void
}

; <alloc::raw_vec::RawVecInner>::finish_grow
; Function Attrs: cold nounwind nonlazybind uwtable
define internal fastcc void @_RNvMs4_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner11finish_growCsj29zdfijJxG_6source(ptr dead_on_unwind noalias noundef nonnull writable writeonly sret([24 x i8]) align 8 captures(none) dereferenceable(24) initializes((0, 8)) %_0, ptr noalias noundef nonnull readonly align 8 captures(none) dereferenceable(16) %self, i64 noundef %cap) unnamed_addr #0 {
start:
  %_39 = icmp sgt i64 %cap, -1
  br i1 %_39, label %bb15, label %bb11

bb15:                                             ; preds = %start
  %self1.i = load i64, ptr %self, align 8, !range !6, !alias.scope !10, !noalias !13, !noundef !7
  %0 = icmp eq i64 %self1.i, 0
  br i1 %0, label %bb5, label %_RNvXs_NtCscdodAO9FK5_5alloc5allocNtB4_6GlobalNtNtCs4NRVxsYgnAr_4core5alloc9Allocator4grow.exit

_RNvXs_NtCscdodAO9FK5_5alloc5allocNtB4_6GlobalNtNtCs4NRVxsYgnAr_4core5alloc9Allocator4grow.exit: ; preds = %bb15
  %1 = getelementptr inbounds nuw i8, ptr %self, i64 8
  %self3.i = load ptr, ptr %1, align 8, !alias.scope !10, !noalias !13, !nonnull !7, !noundef !7
  %cond.i.i = icmp samesign uge i64 %cap, %self1.i
  tail call void @llvm.assume(i1 %cond.i.i)
; call __rustc::__rust_realloc
  %raw_ptr.i.i = tail call noundef ptr @_RNvCs9wFQrvczXsK_7___rustc14___rust_realloc(ptr noundef nonnull %self3.i, i64 noundef %self1.i, i64 noundef 1, i64 noundef range(i64 0, -9223372036854775808) %cap) #11
  br label %bb7

bb5:                                              ; preds = %bb15
  %2 = icmp eq i64 %cap, 0
  br i1 %2, label %bb9, label %bb4.i.i

bb4.i.i:                                          ; preds = %bb5
; call __rustc::__rust_no_alloc_shim_is_unstable_v2
  tail call void @_RNvCs9wFQrvczXsK_7___rustc35___rust_no_alloc_shim_is_unstable_v2() #11
; call __rustc::__rust_alloc
  %3 = tail call noundef ptr @_RNvCs9wFQrvczXsK_7___rustc12___rust_alloc(i64 noundef range(i64 0, -9223372036854775808) %cap, i64 noundef 1) #11
  br label %bb7

bb7:                                              ; preds = %bb4.i.i, %_RNvXs_NtCscdodAO9FK5_5alloc5allocNtB4_6GlobalNtNtCs4NRVxsYgnAr_4core5alloc9Allocator4grow.exit
  %raw_ptr.i.i.pn = phi ptr [ %raw_ptr.i.i, %_RNvXs_NtCscdodAO9FK5_5alloc5allocNtB4_6GlobalNtNtCs4NRVxsYgnAr_4core5alloc9Allocator4grow.exit ], [ %3, %bb4.i.i ]
  %4 = icmp eq ptr %raw_ptr.i.i.pn, null
  br i1 %4, label %bb8, label %bb9

bb8:                                              ; preds = %bb7
  %5 = getelementptr inbounds nuw i8, ptr %_0, i64 8
  store i64 1, ptr %5, align 8
  br label %bb11

bb9:                                              ; preds = %bb5, %bb7
  %raw_ptr.i.i.pn17 = phi ptr [ %raw_ptr.i.i.pn, %bb7 ], [ inttoptr (i64 1 to ptr), %bb5 ]
  %6 = getelementptr inbounds nuw i8, ptr %_0, i64 8
  store ptr %raw_ptr.i.i.pn17, ptr %6, align 8
  br label %bb11

bb11:                                             ; preds = %start, %bb9, %bb8
  %.sink18 = phi i64 [ 16, %bb9 ], [ 16, %bb8 ], [ 8, %start ]
  %cap.sink = phi i64 [ %cap, %bb9 ], [ %cap, %bb8 ], [ 0, %start ]
  %storemerge9 = phi i64 [ 0, %bb9 ], [ 1, %bb8 ], [ 1, %start ]
  %7 = getelementptr inbounds nuw i8, ptr %_0, i64 %.sink18
  store i64 %cap.sink, ptr %7, align 8
  store i64 %storemerge9, ptr %_0, align 8
  ret void
}

; Function Attrs: cold nounwind nonlazybind uwtable
define noundef i32 @main() unnamed_addr #0 personality ptr @rust_eh_personality {
start:
  %args = alloca [16 x i8], align 8
  %_5 = alloca [8 x i8], align 8
  %s = alloca [24 x i8], align 8
  call void @llvm.lifetime.start.p0(ptr nonnull %s)
; call __rustc::__rust_no_alloc_shim_is_unstable_v2
  tail call void @_RNvCs9wFQrvczXsK_7___rustc35___rust_no_alloc_shim_is_unstable_v2() #11, !noalias !15
; call __rustc::__rust_alloc
  %0 = tail call noundef dereferenceable_or_null(3) ptr @_RNvCs9wFQrvczXsK_7___rustc12___rust_alloc(i64 noundef range(i64 0, -9223372036854775808) 3, i64 noundef 1) #11, !noalias !15
  %1 = icmp eq ptr %0, null
  br i1 %1, label %bb7, label %_RNvMs_NtCscdodAO9FK5_5alloc3vecINtB4_3VechE15append_elementsCsj29zdfijJxG_6source.exit, !prof !18

bb7:                                              ; preds = %start
; call alloc::raw_vec::handle_error
  tail call void @_RNvNtCscdodAO9FK5_5alloc7raw_vec12handle_error(i64 noundef 1, i64 3) #12
  unreachable

_RNvMs_NtCscdodAO9FK5_5alloc3vecINtB4_3VechE15append_elementsCsj29zdfijJxG_6source.exit: ; preds = %start
  tail call void @llvm.memcpy.p0.p0.i64(ptr noundef nonnull align 1 dereferenceable(3) %0, ptr noundef nonnull align 1 dereferenceable(3) @alloc_3aef7c697f52c5a3e9548cf3d170303c, i64 3, i1 false)
  store i64 3, ptr %s, align 8
  %_9.sroa.4.0.s.sroa_idx = getelementptr inbounds nuw i8, ptr %s, i64 8
  store ptr %0, ptr %_9.sroa.4.0.s.sroa_idx, align 8
  %_9.sroa.6.0.s.sroa_idx = getelementptr inbounds nuw i8, ptr %s, i64 16
  store i64 3, ptr %_9.sroa.6.0.s.sroa_idx, align 8
  tail call void @llvm.experimental.noalias.scope.decl(metadata !19)
; call <alloc::raw_vec::RawVecInner<_>>::reserve::do_reserve_and_handle::<alloc::alloc::Global>
  call fastcc void @_RINvNvMs2_NtCscdodAO9FK5_5alloc7raw_vecINtB8_11RawVecInnerpE7reserve21do_reserve_and_handleNtNtBa_5alloc6GlobalECsj29zdfijJxG_6source(ptr noalias noundef nonnull align 8 dereferenceable(24) %s, i64 noundef 3) #11
  %len.pre.i = load i64, ptr %_9.sroa.6.0.s.sroa_idx, align 8, !alias.scope !19
  %_10.i = icmp sgt i64 %len.pre.i, -1
  tail call void @llvm.assume(i1 %_10.i)
  %_11.i = load ptr, ptr %_9.sroa.4.0.s.sroa_idx, align 8, !alias.scope !19, !nonnull !7, !noundef !7
  %dst.i = getelementptr inbounds nuw i8, ptr %_11.i, i64 %len.pre.i
  tail call void @llvm.memcpy.p0.p0.i64(ptr noundef nonnull align 1 dereferenceable(3) %dst.i, ptr noundef nonnull align 1 dereferenceable(3) @alloc_d4d3cbd0d4712cfcd6657cf295c08ada, i64 3, i1 false), !noalias !19
  %2 = add nuw nsw i64 %len.pre.i, 3
  call void @llvm.lifetime.start.p0(ptr nonnull %_5)
  store i64 %2, ptr %_5, align 8
  call void @llvm.lifetime.start.p0(ptr nonnull %args)
  store ptr %_5, ptr %args, align 8
  %_7.sroa.4.0..sroa_idx = getelementptr inbounds nuw i8, ptr %args, i64 8
  store ptr @_RNvXsi_NtNtNtCs4NRVxsYgnAr_4core3fmt3num3impjNtB9_7Display3fmt, ptr %_7.sroa.4.0..sroa_idx, align 8
; call std::io::stdio::_print
  call void @_RNvNtNtCs2AWtUsOyxgP_3std2io5stdio6__print(ptr noundef nonnull @alloc_61247b90e1706a3f65e71312b599d3d1, ptr noundef nonnull %args) #11
  call void @llvm.lifetime.end.p0(ptr nonnull %args)
  call void @llvm.lifetime.end.p0(ptr nonnull %_5)
  call void @llvm.experimental.noalias.scope.decl(metadata !22)
  call void @llvm.experimental.noalias.scope.decl(metadata !25)
  call void @llvm.experimental.noalias.scope.decl(metadata !28)
  call void @llvm.experimental.noalias.scope.decl(metadata !31)
  call void @llvm.experimental.noalias.scope.decl(metadata !34)
  %self1.i.i.i.i.i.i = load i64, ptr %s, align 8, !range !6, !alias.scope !37, !noalias !40, !noundef !7
  %3 = icmp eq i64 %self1.i.i.i.i.i.i, 0
  br i1 %3, label %_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueNtNtCscdodAO9FK5_5alloc6string6StringECsj29zdfijJxG_6source.exit, label %_RNvXs_NtCscdodAO9FK5_5alloc5allocNtB4_6GlobalNtNtCs4NRVxsYgnAr_4core5alloc9Allocator10deallocate.exit.i.i.i.i.i

_RNvXs_NtCscdodAO9FK5_5alloc5allocNtB4_6GlobalNtNtCs4NRVxsYgnAr_4core5alloc9Allocator10deallocate.exit.i.i.i.i.i: ; preds = %_RNvMs_NtCscdodAO9FK5_5alloc3vecINtB4_3VechE15append_elementsCsj29zdfijJxG_6source.exit
; call __rustc::__rust_dealloc
  call void @_RNvCs9wFQrvczXsK_7___rustc14___rust_dealloc(ptr noundef nonnull %_11.i, i64 noundef %self1.i.i.i.i.i.i, i64 noundef range(i64 1, -9223372036854775807) 1) #11, !noalias !42
  br label %_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueNtNtCscdodAO9FK5_5alloc6string6StringECsj29zdfijJxG_6source.exit

_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueNtNtCscdodAO9FK5_5alloc6string6StringECsj29zdfijJxG_6source.exit: ; preds = %_RNvMs_NtCscdodAO9FK5_5alloc3vecINtB4_3VechE15append_elementsCsj29zdfijJxG_6source.exit, %_RNvXs_NtCscdodAO9FK5_5alloc5allocNtB4_6GlobalNtNtCs4NRVxsYgnAr_4core5alloc9Allocator10deallocate.exit.i.i.i.i.i
  call void @llvm.lifetime.end.p0(ptr nonnull %s)
  ret i32 0
}

; Function Attrs: mustprogress nocallback nofree nosync nounwind willreturn memory(argmem: readwrite)
declare void @llvm.lifetime.start.p0(ptr captures(none)) #1

; alloc::raw_vec::handle_error
; Function Attrs: cold minsize noreturn nounwind nonlazybind optsize uwtable
declare void @_RNvNtCscdodAO9FK5_5alloc7raw_vec12handle_error(i64 noundef range(i64 0, -9223372036854775807), i64) unnamed_addr #2

; Function Attrs: mustprogress nocallback nofree nosync nounwind willreturn memory(argmem: readwrite)
declare void @llvm.lifetime.end.p0(ptr captures(none)) #1

; Function Attrs: mustprogress nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write)
declare void @llvm.assume(i1 noundef) #3

; Function Attrs: mustprogress nocallback nofree nounwind willreturn memory(argmem: readwrite)
declare void @llvm.memcpy.p0.p0.i64(ptr noalias writeonly captures(none), ptr noalias readonly captures(none), i64, i1 immarg) #4

; __rustc::__rust_dealloc
; Function Attrs: nounwind nonlazybind allockind("free") uwtable
declare void @_RNvCs9wFQrvczXsK_7___rustc14___rust_dealloc(ptr allocptr noundef nonnull captures(address), i64 noundef, i64 noundef range(i64 1, -9223372036854775807)) unnamed_addr #5

; __rustc::__rust_realloc
; Function Attrs: nounwind nonlazybind allockind("realloc,aligned") allocsize(3) uwtable
declare noalias noundef ptr @_RNvCs9wFQrvczXsK_7___rustc14___rust_realloc(ptr allocptr noundef nonnull, i64 noundef, i64 allocalign noundef range(i64 1, -9223372036854775807), i64 noundef) unnamed_addr #6

; __rustc::__rust_no_alloc_shim_is_unstable_v2
; Function Attrs: nounwind nonlazybind uwtable
declare void @_RNvCs9wFQrvczXsK_7___rustc35___rust_no_alloc_shim_is_unstable_v2() unnamed_addr #7

; __rustc::__rust_alloc
; Function Attrs: nounwind nonlazybind allockind("alloc,uninitialized,aligned") allocsize(0) uwtable
declare noalias noundef ptr @_RNvCs9wFQrvczXsK_7___rustc12___rust_alloc(i64 noundef, i64 allocalign noundef range(i64 1, -9223372036854775807)) unnamed_addr #8

; Function Attrs: nounwind nonlazybind uwtable
declare noundef range(i32 0, 10) i32 @rust_eh_personality(i32 noundef, i32 noundef, i64 noundef, ptr noundef, ptr noundef) unnamed_addr #7

; <usize as core::fmt::Display>::fmt
; Function Attrs: nounwind nonlazybind uwtable
declare noundef zeroext i1 @_RNvXsi_NtNtNtCs4NRVxsYgnAr_4core3fmt3num3impjNtB9_7Display3fmt(ptr noalias noundef readonly align 8 captures(address, read_provenance) dereferenceable(8), ptr noalias noundef align 8 dereferenceable(24)) unnamed_addr #7

; std::io::stdio::_print
; Function Attrs: nounwind nonlazybind uwtable
declare void @_RNvNtNtCs2AWtUsOyxgP_3std2io5stdio6__print(ptr noundef nonnull, ptr noundef nonnull) unnamed_addr #7

; Function Attrs: nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: readwrite)
declare void @llvm.experimental.noalias.scope.decl(metadata) #9

; Function Attrs: nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none)
declare i64 @llvm.umax.i64(i64, i64) #10

attributes #0 = { cold nounwind nonlazybind uwtable "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #1 = { mustprogress nocallback nofree nosync nounwind willreturn memory(argmem: readwrite) }
attributes #2 = { cold minsize noreturn nounwind nonlazybind optsize uwtable "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #3 = { mustprogress nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write) }
attributes #4 = { mustprogress nocallback nofree nounwind willreturn memory(argmem: readwrite) }
attributes #5 = { nounwind nonlazybind allockind("free") uwtable "alloc-family"="__rust_alloc" "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #6 = { nounwind nonlazybind allockind("realloc,aligned") allocsize(3) uwtable "alloc-family"="__rust_alloc" "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #7 = { nounwind nonlazybind uwtable "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #8 = { nounwind nonlazybind allockind("alloc,uninitialized,aligned") allocsize(0) uwtable "alloc-family"="__rust_alloc" "alloc-variant-zeroed"="_RNvCs9wFQrvczXsK_7___rustc19___rust_alloc_zeroed" "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #9 = { nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: readwrite) }
attributes #10 = { nocallback nocreateundeforpoison nofree nosync nounwind speculatable willreturn memory(none) }
attributes #11 = { nounwind }
attributes #12 = { noreturn nounwind }

!llvm.module.flags = !{!0, !1}
!llvm.ident = !{!2}

!0 = !{i32 8, !"PIC Level", i32 2}
!1 = !{i32 2, !"RtLibUseGOT", i32 1}
!2 = !{!"rustc version 1.97.1 (8bab26f4f 2026-07-14)"}
!3 = !{!4}
!4 = distinct !{!4, !5, !"_RNvMs4_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14grow_amortizedCsj29zdfijJxG_6source: %self"}
!5 = distinct !{!5, !"_RNvMs4_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14grow_amortizedCsj29zdfijJxG_6source"}
!6 = !{i64 0, i64 -9223372036854775808}
!7 = !{}
!8 = !{i64 0, i64 2}
!9 = !{i64 0, i64 -9223372036854775807}
!10 = !{!11}
!11 = distinct !{!11, !12, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14current_memoryCsj29zdfijJxG_6source: %self"}
!12 = distinct !{!12, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14current_memoryCsj29zdfijJxG_6source"}
!13 = !{!14}
!14 = distinct !{!14, !12, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14current_memoryCsj29zdfijJxG_6source: %_0"}
!15 = !{!16}
!16 = distinct !{!16, !17, !"_RNvMs4_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner15try_allocate_inCsj29zdfijJxG_6source: %_0"}
!17 = distinct !{!17, !"_RNvMs4_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner15try_allocate_inCsj29zdfijJxG_6source"}
!18 = !{!"branch_weights", !"expected", i32 1, i32 2000}
!19 = !{!20}
!20 = distinct !{!20, !21, !"_RNvMs_NtCscdodAO9FK5_5alloc3vecINtB4_3VechE15append_elementsCsj29zdfijJxG_6source: %self"}
!21 = distinct !{!21, !"_RNvMs_NtCscdodAO9FK5_5alloc3vecINtB4_3VechE15append_elementsCsj29zdfijJxG_6source"}
!22 = !{!23}
!23 = distinct !{!23, !24, !"_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueNtNtCscdodAO9FK5_5alloc6string6StringECsj29zdfijJxG_6source: %_1"}
!24 = distinct !{!24, !"_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueNtNtCscdodAO9FK5_5alloc6string6StringECsj29zdfijJxG_6source"}
!25 = !{!26}
!26 = distinct !{!26, !27, !"_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueINtNtCscdodAO9FK5_5alloc3vec3VechEECsj29zdfijJxG_6source: %_1"}
!27 = distinct !{!27, !"_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueINtNtCscdodAO9FK5_5alloc3vec3VechEECsj29zdfijJxG_6source"}
!28 = !{!29}
!29 = distinct !{!29, !30, !"_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueINtNtCscdodAO9FK5_5alloc7raw_vec6RawVechEECsj29zdfijJxG_6source: %_1"}
!30 = distinct !{!30, !"_RINvNtCs4NRVxsYgnAr_4core3ptr9drop_glueINtNtCscdodAO9FK5_5alloc7raw_vec6RawVechEECsj29zdfijJxG_6source"}
!31 = !{!32}
!32 = distinct !{!32, !33, !"_RNvXs1_NtCscdodAO9FK5_5alloc7raw_vecINtB5_6RawVechENtNtNtCs4NRVxsYgnAr_4core3ops4drop4Drop4dropCsj29zdfijJxG_6source: %self"}
!33 = distinct !{!33, !"_RNvXs1_NtCscdodAO9FK5_5alloc7raw_vecINtB5_6RawVechENtNtNtCs4NRVxsYgnAr_4core3ops4drop4Drop4dropCsj29zdfijJxG_6source"}
!34 = !{!35}
!35 = distinct !{!35, !36, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner10deallocateCsj29zdfijJxG_6source: %self"}
!36 = distinct !{!36, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner10deallocateCsj29zdfijJxG_6source"}
!37 = !{!38, !35, !32, !29, !26, !23}
!38 = distinct !{!38, !39, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14current_memoryCsj29zdfijJxG_6source: %self"}
!39 = distinct !{!39, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14current_memoryCsj29zdfijJxG_6source"}
!40 = !{!41}
!41 = distinct !{!41, !39, !"_RNvMs2_NtCscdodAO9FK5_5alloc7raw_vecNtB5_11RawVecInner14current_memoryCsj29zdfijJxG_6source: %_0"}
!42 = !{!35, !32, !29, !26, !23}
