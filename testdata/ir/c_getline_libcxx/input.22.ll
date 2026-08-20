; ModuleID = 'testdata/ir/c_getline_libcxx/source.c'
source_filename = "testdata/ir/c_getline_libcxx/source.c"

define i32 @main() {
entry:
  %is = alloca [64 x i8], align 8
  %str = alloca [32 x i8], align 8
  %p = call ptr @_ZNSt3__17getlineB9nqn220108IcNS_11char_traitsIcEENS_9allocatorIcEEEERNS_13basic_istreamIT_T0_EES9_RNS_12basic_stringIS6_S7_T1_EE(ptr %is, ptr %str)
  ret i32 0
}

declare ptr @_ZNSt3__17getlineB9nqn220108IcNS_11char_traitsIcEENS_9allocatorIcEEEERNS_13basic_istreamIT_T0_EES9_RNS_12basic_stringIS6_S7_T1_EE(ptr, ptr)
