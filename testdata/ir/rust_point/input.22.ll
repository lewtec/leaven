; ModuleID = 'source.ddb5a00486c70246-cgu.0'
source_filename = "source.ddb5a00486c70246-cgu.0"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@alloc_61247b90e1706a3f65e71312b599d3d1 = private unnamed_addr constant [4 x i8] c"\C0\01\0A\00", align 1

; Function Attrs: nounwind nonlazybind uwtable
define noundef i32 @main() unnamed_addr #0 {
start:
  %args = alloca [16 x i8], align 8
  %n = alloca [4 x i8], align 4
  call void @llvm.lifetime.start.p0(ptr nonnull %n)
  store i32 7, ptr %n, align 4
  call void @llvm.lifetime.start.p0(ptr nonnull %args)
  store ptr %n, ptr %args, align 8
  %_8.sroa.4.0..sroa_idx = getelementptr inbounds nuw i8, ptr %args, i64 8
  store ptr @_RNvXs9_NtNtNtCs4NRVxsYgnAr_4core3fmt3num3implNtB9_7Display3fmt, ptr %_8.sroa.4.0..sroa_idx, align 8
; call std::io::stdio::_print
  call void @_RNvNtNtCs2AWtUsOyxgP_3std2io5stdio6__print(ptr noundef nonnull @alloc_61247b90e1706a3f65e71312b599d3d1, ptr noundef nonnull %args) #2
  call void @llvm.lifetime.end.p0(ptr nonnull %args)
  call void @llvm.lifetime.end.p0(ptr nonnull %n)
  ret i32 0
}

; Function Attrs: mustprogress nocallback nofree nosync nounwind willreturn memory(argmem: readwrite)
declare void @llvm.lifetime.start.p0(ptr captures(none)) #1

; <i32 as core::fmt::Display>::fmt
; Function Attrs: nounwind nonlazybind uwtable
declare noundef zeroext i1 @_RNvXs9_NtNtNtCs4NRVxsYgnAr_4core3fmt3num3implNtB9_7Display3fmt(ptr noalias noundef readonly align 4 captures(address, read_provenance) dereferenceable(4), ptr noalias noundef align 8 dereferenceable(24)) unnamed_addr #0

; Function Attrs: mustprogress nocallback nofree nosync nounwind willreturn memory(argmem: readwrite)
declare void @llvm.lifetime.end.p0(ptr captures(none)) #1

; std::io::stdio::_print
; Function Attrs: nounwind nonlazybind uwtable
declare void @_RNvNtNtCs2AWtUsOyxgP_3std2io5stdio6__print(ptr noundef nonnull, ptr noundef nonnull) unnamed_addr #0

attributes #0 = { nounwind nonlazybind uwtable "probe-stack"="inline-asm" "target-cpu"="x86-64" }
attributes #1 = { mustprogress nocallback nofree nosync nounwind willreturn memory(argmem: readwrite) }
attributes #2 = { nounwind }

!llvm.module.flags = !{!0, !1}
!llvm.ident = !{!2}

!0 = !{i32 8, !"PIC Level", i32 2}
!1 = !{i32 2, !"RtLibUseGOT", i32 1}
!2 = !{!"rustc version 1.97.1 (8bab26f4f 2026-07-14)"}
