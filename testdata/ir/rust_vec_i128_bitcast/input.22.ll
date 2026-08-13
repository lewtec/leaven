; ModuleID = 'testdata/ir/rust_vec_i128_bitcast/source.rs'
source_filename = "testdata/ir/rust_vec_i128_bitcast/source.rs"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define i32 @main() {
entry:
  %words = add <2 x i64> zeroinitializer, <i64 1, i64 2>
  %wide = bitcast <2 x i64> %words to i128
  %back = bitcast i128 %wide to <2 x i64>
  %lo = extractelement <2 x i64> %back, i64 0
  %hi = extractelement <2 x i64> %back, i64 1
  %ok1 = icmp eq i64 %lo, 1
  %ok2 = icmp eq i64 %hi, 2
  %ok = and i1 %ok1, %ok2
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
