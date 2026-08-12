; ModuleID = 'testdata/ir/rust_vec_bitcast/source.rs'
source_filename = "testdata/ir/rust_vec_bitcast/source.rs"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define i32 @main() {
entry:
  %bytes = add <16 x i8> zeroinitializer, <i8 1, i8 0, i8 0, i8 0, i8 0, i8 0, i8 0, i8 0, i8 2, i8 0, i8 0, i8 0, i8 0, i8 0, i8 0, i8 0>
  %words = bitcast <16 x i8> %bytes to <2 x i64>
  %lo = extractelement <2 x i64> %words, i64 0
  %hi = extractelement <2 x i64> %words, i64 1
  %ok1 = icmp eq i64 %lo, 1
  %ok2 = icmp eq i64 %hi, 2
  %both = and i1 %ok1, %ok2
  br i1 %both, label %back, label %bad

back:
  %round = bitcast <2 x i64> %words to <16 x i8>
  %b0 = extractelement <16 x i8> %round, i64 0
  %b8 = extractelement <16 x i8> %round, i64 8
  %ok3 = icmp eq i8 %b0, 1
  %ok4 = icmp eq i8 %b8, 2
  %ok = and i1 %ok3, %ok4
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
