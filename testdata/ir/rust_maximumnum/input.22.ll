; ModuleID = 'testdata/ir/rust_maximumnum/source.rs'
source_filename = "testdata/ir/rust_maximumnum/source.rs"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

declare double @llvm.maximumnum.f64(double, double)

define double @f(double %a, double %b) {
entry:
  %r = tail call nsz double @llvm.maximumnum.f64(double %a, double %b)
  ret double %r
}

define i32 @main() {
entry:
  %m = call double @f(double 1.000000e+00, double 2.000000e+00)
  %ok = fcmp oeq double %m, 2.000000e+00
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
