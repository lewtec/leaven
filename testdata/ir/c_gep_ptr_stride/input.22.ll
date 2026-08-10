; ModuleID = 'testdata/ir/c_gep_ptr_stride/source.c'
source_filename = "testdata/ir/c_gep_ptr_stride/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@vt = constant [4 x ptr] zeroinitializer
@p = global ptr getelementptr inbounds (ptr, ptr @vt, i64 2)

define ptr @slot() {
entry:
  ret ptr getelementptr inbounds (ptr, ptr @vt, i64 2)
}
