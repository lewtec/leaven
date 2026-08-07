typedef struct TSLexer TSLexer;
typedef struct { int d; } Payload;
void *create(void) { return 0; }
void destroy(Payload *p) { (void)p; }
_Bool scan(Payload *p, TSLexer *l, const _Bool *s) { (void)p;(void)l;(void)s; return 0; }

typedef void *(*create_fn)(void);
typedef void (*destroy_fn)(void *);
typedef _Bool (*scan_fn)(void *, TSLexer *, const _Bool *);

struct Hooks {
  create_fn create;
  destroy_fn destroy;
  scan_fn scan;
};

struct Hooks h = {
  create,
  (destroy_fn)destroy,
  (scan_fn)scan,
};
