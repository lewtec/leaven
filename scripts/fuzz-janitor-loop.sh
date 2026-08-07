#!/usr/bin/env bash
# Fuzz then janitor-refactree. Exit 1 on leaven/test FAIL (not skip).
# Usage: scripts/fuzz-janitor-loop.sh [iters]
set -euo pipefail
ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT"
RFT="${JANITOR_RFT:-$HOME/.agents/skills/janitor-refactree/janitor.rft}"
ITERS="${1:-30}"
export CSMITH_ITERS="$ITERS"
export CSMITH_TIMEOUT="${CSMITH_TIMEOUT:-5s}"
# optional: export CSMITH_EXTRA_ARGS='...'

echo "== fuzz (CSMITH_ITERS=$ITERS) =="
set +e
go test -tags=csmith -count=1 -v -run 'TestCsmithRandom|TestCsmithFixedSeeds' -timeout 10m
FUZZ=$?
set -e

echo "== janitor-refactree =="
# no -- : workspaced eats it; targets default to -C
set +e
rft run "$RFT" --format table -C "$ROOT"
JAN=$?
set -e

echo "fuzz_exit=$FUZZ janitor_exit=$JAN"
if [[ $JAN -ne 0 ]]; then
  echo "janitor catalog dirty"
  exit 2
fi
exit "$FUZZ"
