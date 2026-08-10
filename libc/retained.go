package libc

// retained pins RustAlloc blocks. Master rustalloc.go still appends here after
// the C path moved to allocs. PR #13 switches RustAlloc onto allocs; delete
// this var once that lands.
var retained []any
