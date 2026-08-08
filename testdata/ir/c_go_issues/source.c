/* leaven-go-issues.c
 * clang -S -emit-llvm -fno-discard-value-names -O0 -o x.ll leaven-go-issues.c
 * go tool leaven -package p x.ll && go build
 */
#include <stdbool.h>

/* 1) nested anonymous struct → %struct.anon (name sanitization / layout) */
struct Outer {
  struct {
    unsigned char a;
    unsigned char b;
  } nested;
  int *ptr;
};

/* 2) Go keyword as identifier → syntax error: unexpected keyword range */
int range;

/* 3) _Bool stored as i8, reloaded as i1 → cannot use byte(v & 1) as bool */
bool use_bool(int x) {
  bool a = x != 0;
  bool b = a && ((x & 1) != 0);
  return b;
}

/* 4) GEP / bitcast → unsafe.Pointer without import "unsafe" */
int use_ptr(struct Outer *o, int i) {
  return o->ptr[i] + (int)o->nested.a;
}

/* keyword again as parameter */
int use_range(int range) {
  return range + 1;
}

void set_range(int v) {
  range = v;
}
