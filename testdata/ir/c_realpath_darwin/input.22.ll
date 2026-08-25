; ModuleID = 'testdata/ir/c_realpath_darwin/source.c'
source_filename = "testdata/ir/c_realpath_darwin/source.c"

@dot = private unnamed_addr constant [2 x i8] c".\00", align 1

define i32 @main() {
entry:
  %p = call ptr @"realpath$DARWIN_EXTSN"(ptr @dot, ptr null)
  %ok = icmp ne ptr %p, null
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}

declare ptr @"realpath$DARWIN_EXTSN"(ptr, ptr)
