# IR fixtures (`testdata/ir/<lang>_<name>/`)

Folder name starts with the language: `c_`, `cpp_`, `rust_`. Example: `c_strlen_map`, `cpp_add`, `rust_add`.

Each folder is one regression case for `TestIRSanity`. Default `go test` reads committed IR only. Clang / clang++ / rustc are gen-only (`go generate .`).

| File | Role |
|------|------|
| `source.c` / `source.cpp` / `source.rs` | Exactly one. Extension picks the compiler and must match the dir prefix. |
| `input.<n>.ll` | IR dialect **n** (LLVM major). Example: clang 14 → `input.14.ll`, rustc LLVM 22 → `input.22.ll`. |
| `expect.json` | `contains` / `not_contains` on generated Go; optional `run`; `"parse": false` skips llir. |

A folder may keep several majors (`input.14.ll` and `input.18.ll`). `TestIRSanity` runs each file.

Compilers are mise conda pins in `mise.toml`:

| Source | Tool | Default pin | Typical `input.<n>.ll` |
|--------|------|-------------|------------------------|
| `source.c` | `clang` | `conda:clang` 14.0.6 | `input.14.ll` |
| `source.cpp` | `clang++` | `conda:clangxx` 14.0.6 | `input.14.ll` |
| `source.rs` | `rustc` | `conda:rust` 1.97.1 | `input.22.ll` |

C/C++ 14 is typed pointers (v14 parser). rustc and clang 22 emit opaque `ptr`.

`expect.json` `"parse"` may be a bool (all majors) or an object of major→bool. Missing majors default to true. Example: `"parse": {"22": false}` runs 14 and skips 22 until the v22 frontend exists.

Another clang major: `go generate .` also runs clang 22, or put `leaven:tool=conda:clang@18.1.8` in the source. Gen writes `input.<n>.ll` and leaves other majors in place.

## Generate (needs compilers)

```bash
mise install
go generate .        # clang 14 + clang 22 + rustc
mise run ir:gen      # same
mise run ir:gen22    # clang 22 only
mise run ir:check    # default pins only
```

Skip generation when the compiler cannot emit the pattern. Put `leaven:hand-ir` in the source. Current hand IR: `c_anon_dot`, `c_name_clash_func_local`, `c_unreachable`, `c_void_mid_return`.

## Test (no clang)

```bash
go test ./...
mise run test:live    # optional C vs Go; needs clang
```
