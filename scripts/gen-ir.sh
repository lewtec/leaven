#!/usr/bin/env bash
# Emit testdata/ir/*/input.<LLVM-major>.ll from source.{c,cpp,rs} via mise conda tools.
# Skip fixtures whose source contains leaven:hand-ir.
# Optional in source: leaven:tool=conda:clang@18.1.8  (mise exec spec)
#
# Compilers are gen-only. Default go test reads committed input.<n>.ll.
#
# Usage: scripts/gen-ir.sh
#        scripts/gen-ir.sh --check   # regenerate to temp, fail if stale
set -euo pipefail
ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT"

CHECK=0
if [[ "${1:-}" == "--check" ]]; then
  CHECK=1
elif [[ "${1:-}" != "" ]]; then
  echo "usage: $0 [--check]" >&2
  exit 1
fi

if ! command -v mise >/dev/null 2>&1; then
  echo "mise not found (tools come from mise.toml conda:*)" >&2
  exit 2
fi

find_source() {
  local dir=$1
  local hits=() f
  for f in "$dir"/source.c "$dir"/source.cpp "$dir"/source.rs; do
    [[ -f $f ]] && hits+=("$f")
  done
  if [[ ${#hits[@]} -ne 1 ]]; then
    echo "$dir: need exactly one of source.c, source.cpp, source.rs (found ${#hits[@]})" >&2
    return 1
  fi
  printf '%s\n' "${hits[0]}"
}

tool_override() {
  grep -oE 'leaven:tool=[^[:space:]]+' "$1" 2>/dev/null | head -n1 | sed 's/^leaven:tool=//' || true
}

mise_run() {
  local src=$1
  shift
  local override
  override=$(tool_override "$src")
  if [[ -n $override ]]; then
    mise exec -C "$ROOT" "$override" -- "$@"
  else
    mise exec -C "$ROOT" -- "$@"
  fi
}

compile_one() {
  local src=$1 out=$2
  case $src in
    *.c)
      mise_run "$src" clang -S -emit-llvm -fno-discard-value-names -std=gnu11 -O0 -o "$out" "$src"
      ;;
    *.cpp)
      mise_run "$src" clang++ -S -emit-llvm -fno-discard-value-names -std=c++17 -O0 -o "$out" "$src"
      ;;
    *.rs)
      mise_run "$src" rustc --emit=llvm-ir -C opt-level=0 -C debuginfo=0 --crate-type=lib -o "$out" "$src"
      ;;
    *)
      echo "unknown source: $src" >&2
      return 1
      ;;
  esac
}

# LLVM IR major written into input.<n>.ll (clang major, or rustc's LLVM).
ir_major() {
  local src=$1
  local line major
  case $src in
    *.c)
      line=$(mise_run "$src" clang --version | head -n1)
      major=$(sed -n 's/.*clang version \([0-9][0-9]*\).*/\1/p' <<<"$line")
      ;;
    *.cpp)
      line=$(mise_run "$src" clang++ --version | head -n1)
      major=$(sed -n 's/.*clang version \([0-9][0-9]*\).*/\1/p' <<<"$line")
      ;;
    *.rs)
      line=$(mise_run "$src" rustc -vV | grep '^LLVM version:')
      major=$(sed -n 's/^LLVM version: \([0-9][0-9]*\).*/\1/p' <<<"$line")
      ;;
  esac
  if [[ -z ${major:-} ]]; then
    echo "could not parse LLVM major from compiler for $src" >&2
    return 1
  fi
  printf '%s\n' "$major"
}

need_clang=0
need_cxx=0
need_rust=0
for dir in testdata/ir/*/; do
  dir=${dir%/}
  [[ -d $dir ]] || continue
  src=$(find_source "$dir") || exit 1
  if grep -q 'leaven:hand-ir' "$src"; then
    continue
  fi
  case $src in
    *.c) need_clang=1 ;;
    *.cpp) need_cxx=1 ;;
    *.rs) need_rust=1 ;;
  esac
done

preflight() {
  local bin=$1
  if ! mise exec -C "$ROOT" -- "$bin" --version >/dev/null 2>&1; then
    echo "mise exec -- $bin failed (mise install; conda pin in mise.toml)" >&2
    exit 2
  fi
}
[[ $need_clang -eq 1 ]] && preflight clang
[[ $need_cxx -eq 1 ]] && preflight clang++
[[ $need_rust -eq 1 ]] && preflight rustc

if [[ $need_clang -eq 1 ]]; then
  echo "clang:  $(mise exec -C "$ROOT" -- clang --version | head -n1)"
fi
if [[ $need_cxx -eq 1 ]]; then
  echo "clang++: $(mise exec -C "$ROOT" -- clang++ --version | head -n1)"
fi
if [[ $need_rust -eq 1 ]]; then
  echo "rustc:  $(mise exec -C "$ROOT" -- rustc --version)  LLVM $(mise exec -C "$ROOT" -- rustc -vV | sed -n 's/^LLVM version: //p')"
fi

tmp=
if [[ $CHECK -eq 1 ]]; then
  tmp=$(mktemp -d)
  trap 'rm -rf "$tmp"' EXIT
fi

regen=0
skip=0
stale=0
for dir in testdata/ir/*/; do
  dir=${dir%/}
  [[ -d $dir ]] || continue
  src=$(find_source "$dir")
  name=$(basename "$dir")
  if grep -q 'leaven:hand-ir' "$src"; then
    echo "skip $name (hand-ir)"
    skip=$((skip + 1))
    continue
  fi
  major=$(ir_major "$src")
  dest="$dir/input.$major.ll"
  if [[ $CHECK -eq 1 ]]; then
    out="$tmp/$name.$major.ll"
    compile_one "$src" "$out"
    if [[ ! -f $dest ]] || ! cmp -s "$dest" "$out"; then
      echo "stale $name input.$major.ll (run mise run ir:gen)"
      stale=$((stale + 1))
    else
      echo "ok   $name input.$major.ll"
    fi
  else
    echo "gen  $name input.$major.ll"
    compile_one "$src" "$dest"
    rm -f "$dir/input.ll"
    regen=$((regen + 1))
  fi
done

if [[ $CHECK -eq 1 ]]; then
  echo "stale=$stale skipped=$skip"
  [[ $stale -eq 0 ]]
else
  echo "regenerated=$regen skipped=$skip"
fi
