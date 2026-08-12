; ModuleID = 'testdata/ir/c_ios_base_ctor/source.c'
source_filename = "testdata/ir/c_ios_base_ctor/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@msg = private unnamed_addr constant [3 x i8] c"ok\00"

define i32 @main() {
entry:
  %obj = alloca [272 x i8], align 8
  call void @_ZNSt8ios_baseC2Ev(ptr %obj)
  %vp = load ptr, ptr %obj
  %nullvp = icmp eq ptr %vp, null
  br i1 %nullvp, label %bad, label %chk

chk:
  %slot = getelementptr i8, ptr %vp, i64 -24
  %off = load i64, ptr %slot
  %z = icmp eq i64 %off, 0
  %ctp = getelementptr i8, ptr %obj, i64 240
  %ct = load ptr, ptr %ctp
  %nullct = icmp eq ptr %ct, null
  %ok = xor i1 %nullct, true
  %both = and i1 %z, %ok
  br i1 %both, label %good, label %bad

good:
  %p = call i32 @puts(ptr @msg)
  ret i32 0

bad:
  ret i32 1
}

declare void @_ZNSt8ios_baseC2Ev(ptr)
declare i32 @puts(ptr)
