/* Null/undef ptr in a struct slot must not emit uintptr(nil).
 * leaven:hand-ir
 */
struct slot { char *p; };
static struct slot s;
int main(void) {
  s.p = 0;
  return s.p == 0 ? 0 : 1;
}
