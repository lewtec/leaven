/* C main(argc, argv) from os.Args.
 * leaven:hand-ir
 *
 * IR keeps explicit argc/argv uses; clang -O0 would alloca them.
 */
int main(int argc, char **argv) {
  (void)argv;
  return argc >= 1 ? 0 : 1;
}
