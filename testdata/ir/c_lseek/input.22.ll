; ModuleID = 'testdata/ir/c_lseek/source.c'
source_filename = "testdata/ir/c_lseek/source.c"

define i32 @main() {
entry:
  %r = call i64 @lseek(i32 1, i64 0, i32 1)
  %ok = icmp sge i64 %r, 0
  %z = zext i1 %ok to i32
  ret i32 %z
}

declare i64 @lseek(i32, i64, i32)
