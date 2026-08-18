; ModuleID = 'testdata/ir/c_close_nocancel/source.c'
source_filename = "testdata/ir/c_close_nocancel/source.c"

define i32 @main() {
entry:
  %r = call i32 @"close$NOCANCEL"(i32 -1)
  ret i32 %r
}

declare i32 @"close$NOCANCEL"(i32)
