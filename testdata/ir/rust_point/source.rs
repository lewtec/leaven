// rust_point — struct+impl, enum match, println!.
// leaven:opt-level=2  (O0 leaves core::fmt as declare-only)
pub struct Point {
    pub x: i32,
    pub y: i32,
}

impl Point {
    pub fn manhattan(&self) -> i32 {
        abs_i32(self.x) + abs_i32(self.y)
    }
}

fn abs_i32(x: i32) -> i32 {
    if x < 0 {
        -x
    } else {
        x
    }
}

pub enum Cell {
    Empty,
    Neg,
    Value(i32),
}

pub fn cell_value(c: Cell) -> i32 {
    match c {
        Cell::Empty => 0,
        Cell::Neg => -1,
        Cell::Value(v) => v,
    }
}

#[no_mangle]
pub extern "C" fn main() -> i32 {
    let p = Point { x: 3, y: -4 };
    let n = p.manhattan() + cell_value(Cell::Value(0));
    println!("{n}");
    0
}
