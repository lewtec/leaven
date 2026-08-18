; ModuleID = 'testdata/ir/c_fstat/source.c'
source_filename = "testdata/ir/c_fstat/source.c"

@st = global [256 x i8] zeroinitializer

define i32 @main() {
entry:
  %r = call i32 @fstat(i32 1, ptr @st)
  ret i32 %r
}

declare i32 @fstat(i32, ptr)
