; ModuleID = 'testdata/ir/rust_i256/source.rs'
source_filename = "testdata/ir/rust_i256/source.rs"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define i32 @main() {
entry:
  %w = zext i128 10000 to i256
  %q = udiv i256 %w, 100
  %r = urem i256 %w, 100
  %s = add i256 %q, %r
  %ok = icmp eq i256 %s, 100
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
