/* void_mid_return — mid-function `ret void` (not only at last block).
 * leaven:hand-ir
 *
 * Bug: TermRet with X==nil fell through to FormatValue(nil) after writing
 * `return`, breaking haskell scanners (e.g. take_line_escaped_newline).
 *
 * Clang often merges early returns; the checked-in input.ll is hand IR that
 * reliably has ret void in a non-final block:
 *
 *   define void @f(i32 %x) {
 *   entry:
 *     %cmp = icmp eq i32 %x, 0
 *     br i1 %cmp, label %early, label %rest
 *   early:
 *     ret void
 *   rest:
 *     %add = add i32 %x, 1
 *     ret void
 *   }
 *
 * Approximate C (may not produce the same CFG at all clang versions):
 */
__attribute__((noinline)) void f(int x) {
  if (x == 0)
    return;
  (void)(x + 1);
}
