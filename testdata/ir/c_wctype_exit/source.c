#include <stdlib.h>
#include <wctype.h>

int check(wint_t c) {
  if (iswcntrl(c)) return 1;
  if (iswxdigit(c)) return 2;
  return (int)towupper(c) + (int)towlower(c);
}

void die(void) { exit(1); }

int main(void) {
  if (check(L'a') == 0) die();
  return 0;
}
