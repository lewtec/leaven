/* unreachable — LLVM `unreachable` terminator (noreturn / dead branch).
 * leaven:hand-ir
 *
 * Leaven must emit panic("unreachable"), not invalid Go.
 * clang -O0 emits call @abort, not `unreachable`; input.ll is the oracle.
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
