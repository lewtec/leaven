// rust_point — struct+impl, enum match, C ABI main. No println! (pulls libstd fmt).
extern "C" {
    fn printf(fmt: *const i8, ...) -> i32;
}

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
    unsafe {
        printf(b"%d\n\0".as_ptr() as *const i8, n);
    }
    0
}
