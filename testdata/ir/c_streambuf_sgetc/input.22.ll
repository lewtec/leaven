; ModuleID = 'testdata/ir/c_streambuf_sgetc/source.c'
source_filename = "testdata/ir/c_streambuf_sgetc/source.c"

define i32 @main() {
entry:
  %c = call i32 @_ZNSt3__115basic_streambufIcNS_11char_traitsIcEEE5sgetcEv(ptr null)
  %p = call ptr @_ZNKSt3__115basic_streambufIcNS_11char_traitsIcEEE4gptrEv(ptr null)
  %ok = icmp eq i32 %c, -1
  br i1 %ok, label %good, label %bad

good:
  ret i32 0

bad:
  ret i32 1
}

declare i32 @_ZNSt3__115basic_streambufIcNS_11char_traitsIcEEE5sgetcEv(ptr)
declare ptr @_ZNKSt3__115basic_streambufIcNS_11char_traitsIcEEE4gptrEv(ptr)
