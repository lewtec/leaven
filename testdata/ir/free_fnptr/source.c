#include <stdlib.h>

void (*current_free)(void *) = free;

void release(void *p) {
  current_free(p);
}

void set_free(void (*fn)(void *)) {
  current_free = fn ? fn : free;
}
