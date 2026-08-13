; ModuleID = 'testdata/ir/c_string_sso/source.c'
source_filename = "testdata/ir/c_string_sso/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@s = global { ptr, i64, [16 x i8] } { ptr getelementptr inbounds (i8, ptr @s, i64 16), i64 0, [16 x i8] zeroinitializer }

define i32 @main() {
entry:
  %p = load ptr, ptr @s
  %q = getelementptr inbounds i8, ptr @s, i64 16
  %ok = icmp eq ptr %p, %q
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}
