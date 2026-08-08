/* ctype_wctype_dup — __ctype_b_loc, iswalnum, dup */
#include <ctype.h>
#include <unistd.h>
#include <wctype.h>

int print_p(int c) {
  return (c > 0 && c < 128 && ((*__ctype_b_loc())[c] & 16384) != 0);
}

int alnum_p(wint_t c) {
  return iswalnum(c);
}

int copy_fd(int fd) {
  return dup(fd);
}
