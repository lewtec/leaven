/* unreachable — LLVM `unreachable` terminator (noreturn / dead branch).
 *
 * Leaven must emit panic("unreachable"), not invalid Go.
 *
 *   clang-14 -S -emit-llvm -fno-discard-value-names -std=gnu11 -O0 -o input.ll source.c
 *
 * Clang may not emit `unreachable` for all patterns at -O0; input.ll is the
 * regression oracle (hand-curated if needed).
 */
#include <stdlib.h>

void die(void) {
  abort(); /* often becomes unreachable after noreturn annotation in IR */
}

int f(int x) {
  if (x < 0) {
    __builtin_unreachable();
  }
  return x + 1;
}
