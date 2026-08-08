/* anon_dot — clang nested anonymous structs become %struct.anon.1 in IR.
 * leaven:hand-ir
 *
 * TypeName must sanitize "anon.1" → "anon_1" (legal Go identifier).
 * clang 14 -O0 emits an unnamed struct here; input.ll keeps %struct.anon.1.
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
