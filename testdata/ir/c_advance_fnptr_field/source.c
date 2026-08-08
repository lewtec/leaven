typedef struct TSLexer TSLexer;
struct TSLexer {
  int lookahead;
  void (*advance)(TSLexer *, int skip);
};

static void real_advance(TSLexer *l, int skip) {
  (void)skip;
  l->lookahead++;
}

static void advance(TSLexer *lexer) {
  lexer->advance(lexer, 0);
}

int scan(TSLexer *lexer) {
  advance(lexer);
  return lexer->lookahead;
}

int main(void) {
  TSLexer lx = {0, real_advance};
  return scan(&lx) == 1 ? 0 : 1;
}
