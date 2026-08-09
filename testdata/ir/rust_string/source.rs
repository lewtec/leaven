// rust_string — std::string heap via String::from + push_str.
// leaven:opt-level=2  (O0 leaves alloc::raw_vec as declare-only)
extern "C" {
    fn printf(fmt: *const i8, ...) -> i32;
}

#[no_mangle]
pub extern "C" fn main() -> i32 {
    let mut s = String::from("foo");
    s.push_str("bar");
    unsafe {
        printf(b"%zu\n\0".as_ptr() as *const i8, s.len());
    }
    0
}
