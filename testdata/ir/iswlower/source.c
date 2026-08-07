#include <wctype.h>
int lo(wint_t c) { return iswlower(c); }
int up(wint_t c) { return iswupper(c); }
