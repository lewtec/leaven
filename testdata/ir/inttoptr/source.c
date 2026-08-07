#include <stdint.h>

void *int_to_ptr(uint64_t x) {
  return (void *)(uintptr_t)x;
}

uint64_t ptr_to_int(void *p) {
  return (uint64_t)(uintptr_t)p;
}
