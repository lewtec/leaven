; ModuleID = 'testdata/ir/c_fcntl/source.c'
source_filename = "testdata/ir/c_fcntl/source.c"

define i32 @main() {
entry:
  %r = call i32 (i32, i32, ...) @fcntl(i32 1, i32 3)
  %ok = icmp sge i32 %r, 0
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}

declare i32 @fcntl(i32, i32, ...)
