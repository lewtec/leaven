/* atomicrmw — clang atomicrmw with align; strip align before llir parse */
#include <stdatomic.h>

int atomic_inc(int *p) {
  return atomic_fetch_add_explicit((_Atomic int *)p, 1, memory_order_seq_cst) + 1;
}

int atomic_dec(int *p) {
  return atomic_fetch_sub_explicit((_Atomic int *)p, 1, memory_order_seq_cst) - 1;
}
