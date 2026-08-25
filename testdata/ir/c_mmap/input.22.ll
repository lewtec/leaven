; ModuleID = 'testdata/ir/c_mmap/source.c'
source_filename = "testdata/ir/c_mmap/source.c"

define i32 @main() {
entry:
  %p = call ptr @mmap(ptr null, i64 4096, i32 3, i32 34, i32 -1, i64 0)
  %ok = icmp ne ptr %p, inttoptr (i64 -1 to ptr)
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}

declare ptr @mmap(ptr, i64, i32, i32, i32, i64)
