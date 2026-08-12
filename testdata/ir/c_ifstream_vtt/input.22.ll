; ModuleID = 'testdata/ir/c_ifstream_vtt/source.c'
source_filename = "testdata/ir/c_ifstream_vtt/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@_ZTTSt14basic_ifstreamIcSt11char_traitsIcEE = external unnamed_addr constant [4 x ptr]

define i32 @main() {
  %p = load ptr, ptr @_ZTTSt14basic_ifstreamIcSt11char_traitsIcEE
  %q = getelementptr i8, ptr %p, i64 -24
  %o = load i64, ptr %q
  %c = trunc i64 %o to i32
  ret i32 %c
}
