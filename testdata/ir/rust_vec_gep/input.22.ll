; ModuleID = 'testdata/ir/rust_vec_gep/source.rs'
source_filename = "testdata/ir/rust_vec_gep/source.rs"
target triple = "x86_64-unknown-linux-gnu"

define i32 @main() {
  %p = alloca [2 x ptr]
  %base = getelementptr [2 x ptr], ptr %p, i64 0, i64 0
  %vec = load <2 x ptr>, ptr %base
  %off = getelementptr i8, <2 x ptr> %vec, i64 0
  store <2 x ptr> %off, ptr %base
  ret i32 0
}
