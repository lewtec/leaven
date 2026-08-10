/* gep ptr, [N x ptr]*, i64 n — C++ vtable / typeinfo form.
 * leaven:hand-ir
 *
 * clang -O0 on C emits gep i8 or typed geeps, not this.
 */
void *slot(void);

void *slot(void) { return 0; }
