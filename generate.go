package leaven

// IR fixtures: testdata/ir/*/input.<n>.ll
// Default mise pins (clang 14 / rustc). Then clang 22 beside existing majors.
//
//go:generate mise exec -- go run ./internal/genir
//go:generate mise exec -- go run ./internal/genir -c conda:clang@22.1.8 -cxx conda:clangxx@22.1.8
