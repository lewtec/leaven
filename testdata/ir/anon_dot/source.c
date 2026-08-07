/* anon_dot — clang nested anonymous structs become %struct.anon.1 in IR.
 *
 * TypeName must sanitize "anon.1" → "anon_1" (legal Go identifier).
 *
 *   clang-14 -S -emit-llvm -fno-discard-value-names -std=gnu11 -O0 -o input.ll source.c
 *
 * Note: clang may name the type slightly differently by version; the checked-in
 * input.ll is the regression oracle if IR differs.
 */
struct Wrapper {
  struct {
    char a;
    char b;
  }; /* anonymous nested — often lowered as anon.N */
};

static const struct {
  char a;
  char b;
} g = {1, 2};

char get(void) {
  return g.a;
}
