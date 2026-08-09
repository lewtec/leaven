/* iswblank — iswblank wctype map */
#include <wctype.h>

int blank_p(wint_t c) { return iswblank(c); }

int main(void) {
  if (!iswblank(L' ')) return 1;
  if (!iswblank(L'\t')) return 2;
  if (iswblank(L'\n')) return 3;
  if (iswblank(L'a')) return 4;
  return 0;
}
