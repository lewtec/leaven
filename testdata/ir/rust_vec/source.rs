// rust_vec — Vec<Point> grow + iterator map/sum + println!.
// leaven:opt-level=2  (O0 leaves RawVec / fmt as declare-only)
pub struct Point {
    pub x: i32,
    pub y: i32,
}

impl Point {
    pub fn manhattan(&self) -> i32 {
        self.x.abs() + self.y.abs()
    }
}

#[no_mangle]
pub extern "C" fn main() -> i32 {
    let mut v = Vec::new();
    v.push(Point { x: 3, y: -4 });
    v.push(Point { x: 1, y: 2 });
    let n: i32 = v.iter().map(|p| p.manhattan()).sum();
    println!("{n}");
    0
}
