/* free_assert — free + __assert_fail libc maps */
#include <assert.h>
#include <stdlib.h>

void release(void *p) {
  free(p);
}

int need(int x) {
  assert(x > 0);   /* → __assert_fail in IR */
  return x;
}
