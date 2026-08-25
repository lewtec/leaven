; ModuleID = 'testdata/ir/c_ptrmask/source.c'
source_filename = "testdata/ir/c_ptrmask/source.c"

define i32 @main() {
entry:
  %p = inttoptr i64 255 to ptr
  %r = call ptr @llvm.ptrmask.p0.i64(ptr %p, i64 -16)
  %i = ptrtoint ptr %r to i64
  %ok = icmp eq i64 %i, 240
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}

declare ptr @llvm.ptrmask.p0.i64(ptr, i64)
