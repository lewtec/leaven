; ModuleID = 'testdata/ir/rust_i1_bitcast/source.rs'
source_filename = "testdata/ir/rust_i1_bitcast/source.rs"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define i32 @main() {
entry:
  %eq = icmp eq <16 x i8> zeroinitializer, zeroinitializer
  %bits = bitcast <16 x i1> %eq to i16
  %ok = icmp eq i16 %bits, -1
  br i1 %ok, label %chk, label %bad

chk:
  %back = bitcast i16 %bits to <16 x i1>
  %lane = extractelement <16 x i1> %back, i64 0
  br i1 %lane, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
