/* declare-only cxxabi vtable + __cxa_pure_virtual.
 * leaven:hand-ir
 *
 * clang on this C does not emit those symbols; C++ typeinfo does.
 */
int main(void) { return 0; }
