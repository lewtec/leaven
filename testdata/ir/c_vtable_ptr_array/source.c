/* [N x ptr] slots: string, typeinfo, function, null.
 * leaven:hand-ir
 *
 * clang on C does not emit this mix; C++ vtables do.
 */
void *slot(void);
void *slot(void) { return 0; }
