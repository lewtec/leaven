; ModuleID = 'testdata/ir/rust_sext_i1vec/source.rs'
source_filename = "testdata/ir/rust_sext_i1vec/source.rs"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define i32 @main() {
entry:
  %eq = icmp eq <8 x i8> zeroinitializer, zeroinitializer
  %sx = sext <8 x i1> %eq to <8 x i8>
  %lane = extractelement <8 x i8> %sx, i64 0
  %ok = icmp eq i8 %lane, -1
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
