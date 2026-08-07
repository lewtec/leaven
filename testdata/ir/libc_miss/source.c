#include <stdio.h>
#include <stdlib.h>

void die(const char *msg) {
  fprintf(stderr, "error: %s\n", msg);
  abort();
}

void *grow(void *p, size_t n) {
  void *q = realloc(p, n);
  if (!q)
    die("oom");
  return q;
}

void release(void *p) {
  free(p);
}
