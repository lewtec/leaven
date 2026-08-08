/* stdio-wctype-memmove-repro.c
 * clang-14 -S -emit-llvm -fno-discard-value-names -std=gnu11 -O0 -o x.ll ...
 */
#include <stdio.h>
#include <string.h>
#include <wctype.h>

void write_bits(FILE *f, const char *s, int c) {
  fputs(s, f);
  fputc(c, f);
}

int space_p(wint_t c) {
  return iswspace(c);
}

void move_buf(void *dst, const void *src, size_t n) {
  memmove(dst, src, n);  /* → llvm.memmove.* */
}
