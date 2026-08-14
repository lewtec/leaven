/* Load a ptr out of [N x ptr]. Slot is uint64; SSA is unsafe.Pointer.
 * leaven:hand-ir
 */
static char msg[] = "ok";
static char *arr[2] = {msg, 0};

int main(void) {
  return arr[0][0] == 'o' ? 0 : 1;
}
