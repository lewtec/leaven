; ModuleID = 'testdata/ir/c_locale_ctor/source.c'
source_filename = "testdata/ir/c_locale_ctor/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

define i32 @main() {
entry:
  %obj = alloca [16 x i8], align 8
  call void @_ZNSt6localeC1Ev(ptr %obj)
  ret i32 0
}

declare void @_ZNSt6localeC1Ev(ptr)
