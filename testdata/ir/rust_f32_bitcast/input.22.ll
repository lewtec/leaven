; ModuleID = 'testdata/ir/rust_f32_bitcast/source.rs'
source_filename = "testdata/ir/rust_f32_bitcast/source.rs"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define i32 @main() {
entry:
  %f = fadd float 1.000000e+00, 0.000000e+00
  %bits = bitcast float %f to i32
  %ok1 = icmp eq i32 %bits, 1065353216
  %back = bitcast i32 %bits to float
  %ok2 = fcmp oeq float %back, 1.000000e+00
  %both = and i1 %ok1, %ok2
  br i1 %both, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
