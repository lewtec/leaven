; ModuleID = 'testdata/ir/c_mmap_fixed/source.c'
source_filename = "testdata/ir/c_mmap_fixed/source.c"

define i32 @main() {
entry:
  %p = inttoptr i64 123456789 to ptr
  %r = call ptr @mmap64(ptr %p, i64 4096, i32 3, i32 18, i32 -1, i64 0)
  %ok = icmp eq ptr %r, %p
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}

declare ptr @mmap64(ptr, i64, i32, i32, i32, i64)
