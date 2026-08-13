; ModuleID = 'testdata/ir/c_packed_bit_iter/source.c'
source_filename = "testdata/ir/c_packed_bit_iter/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

%"struct.std::_Bit_iterator_base.base" = type <{ ptr, i32 }>

define i32 @main() {
entry:
  %it = alloca %"struct.std::_Bit_iterator_base.base", align 8
  store i64 1, ptr %it, align 8
  %off = getelementptr inbounds nuw i8, ptr %it, i64 8
  store i32 1, ptr %off, align 4
  %p = call ptr @malloc(i64 noundef 64)
  store i64 0, ptr %p, align 8
  call void @free(ptr noundef %p)
  ret i32 0
}

declare ptr @malloc(i64 noundef)
declare void @free(ptr noundef)
