/* name_clash_func_local — local %is_verbatim vs function @is_verbatim.
 * leaven:hand-ir
 * Valid C cannot call the function once the local shadows it; input.ll is
 * hand IR that models scanners where SSA uses the same base name.
 * Fix: rename locals to local_<name> when they collide with a function.
 */
typedef struct { int string_type; } Interpolation;
static int is_verbatim(Interpolation *i) { return i->string_type & 1; }
int scan(Interpolation *cur) {
  int is_verbatim = 0;
  if (cur->string_type & 1)
    is_verbatim = 1;
  return is_verbatim; /* local only in real C; IR oracle also calls @is_verbatim */
}
int main(void) {
  Interpolation i = {1};
  return scan(&i) ? 0 : 1;
}
