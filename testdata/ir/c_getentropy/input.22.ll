; ModuleID = 'testdata/ir/c_getentropy/source.c'
source_filename = "testdata/ir/c_getentropy/source.c"

@buf = global [16 x i8] zeroinitializer

define i32 @main() {
entry:
  %r = call i32 @getentropy(ptr @buf, i64 16)
  ret i32 %r
}

declare i32 @getentropy(ptr, i64)
