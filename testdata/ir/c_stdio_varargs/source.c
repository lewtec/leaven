/* stdio_varargs — va_start/end, vsnprintf, snprintf, fdopen, fclose */
#include <stdarg.h>
#include <stdio.h>

void logf(const char *fmt, ...) {
  char buf[128];
  va_list ap;
  va_start(ap, fmt);
  vsnprintf(buf, sizeof buf, fmt, ap);
  va_end(ap);
  snprintf(buf, sizeof buf, "%s", fmt);
}

void file_bits(const char *path) {
  FILE *f = fdopen(1, "w");
  if (f)
    fclose(f);
  (void)path;
}
