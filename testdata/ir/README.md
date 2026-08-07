# IR fixtures (`testdata/ir/<name>/`)

Each folder is one regression case for `TestIRSanity`.

| File | Role |
|------|------|
| `source.c` | Human-readable C repro (when possible). Document intent + clang flags. |
| `input.ll` | **Oracle** IR fed to leaven (clang-14 typed pointers). May be hand-curated when clang will not emit the pattern. |
| `expect.json` | Table of `contains` / `not_contains` on generated Go; optional `run` builds+executes package main. |

## Regenerate IR (when source.c is authoritative)

```bash
clang-14 -S -emit-llvm -fno-discard-value-names -std=gnu11 -O0 \
  -o testdata/ir/<name>/input.ll testdata/ir/<name>/source.c
```

Do **not** blindly regenerate if the comment on `source.c` says `input.ll` is hand-curated.

## Run

```bash
go test -run TestIR ./...
```
