// rust_string — String::from + push_str + println!.
// leaven:opt-level=2  (O0 leaves alloc::raw_vec / core::fmt as declare-only)
#[no_mangle]
pub extern "C" fn main() -> i32 {
    let mut s = String::from("foo");
    s.push_str("bar");
    println!("{}", s.len());
    0
}
