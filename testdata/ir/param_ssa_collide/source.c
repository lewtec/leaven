/* params named v1/v2 like tree-sitter stack APIs */
void renumber(int *self, int v1, int v2) {
  int a = v1;
  int b = v2;
  if (a == b)
    *self = a;
  else
    *self = b + a;
}
