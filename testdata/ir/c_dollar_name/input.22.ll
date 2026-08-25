; ModuleID = 'testdata/ir/c_dollar_name/source.c'
source_filename = "testdata/ir/c_dollar_name/source.c"

define i32 @"foo$bar"() {
entry:
  ret i32 0
}

define i32 @main() {
entry:
  %x = call i32 @"foo$bar"()
  ret i32 %x
}
