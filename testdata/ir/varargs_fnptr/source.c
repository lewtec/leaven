/* varargs function pointer field vs definition (ts_lexer__log pattern) */
typedef struct L {
  void (*log)(struct L *, const char *, ...);
} L;

static void logf(struct L *self, const char *fmt, ...) {
  (void)self;
  (void)fmt;
}

void init_log(L *l) {
  l->log = logf;
}
