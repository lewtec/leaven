; ModuleID = 'testdata/ir/c_pthread_stackaddr/source.c'
source_filename = "testdata/ir/c_pthread_stackaddr/source.c"

define i32 @main() {
entry:
  %t = call i64 @pthread_self()
  %p = call ptr @pthread_get_stackaddr_np(i64 %t)
  %n = call i64 @pthread_get_stacksize_np(i64 %t)
  %okp = icmp ne ptr %p, null
  %okn = icmp uge i64 %n, 1
  %ok = and i1 %okp, %okn
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}

declare i64 @pthread_self()
declare ptr @pthread_get_stackaddr_np(i64)
declare i64 @pthread_get_stacksize_np(i64)
