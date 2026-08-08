struct state {
  int depth;
};

static struct state *state_new(void) {
  static struct state pool;
  pool.depth = 0;
  return &pool;
}

void *create(void) {
  struct state *state = state_new();
  return state;
}

int main(void) {
  return create() != 0 ? 0 : 1;
}
