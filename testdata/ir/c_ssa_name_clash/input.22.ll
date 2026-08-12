; ModuleID = 'testdata/ir/c_ssa_name_clash/source.c'
source_filename = "testdata/ir/c_ssa_name_clash/source.c"
target triple = "x86_64-unknown-linux-gnu"

define i32 @main() {
entry:
  %0 = add i32 1, 2
  %v0 = add i32 3, 4
  %s = add i32 %0, %v0
  ret i32 0
}
