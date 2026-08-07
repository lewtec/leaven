/* subtree_tagged — tree-sitter Subtree tagged union (inline LSB vs heap ptr)
 * clang-14 -S -emit-llvm -fno-discard-value-names -std=gnu11 -O0 -o input.ll source.c
 * Expect: inline depends=0, heap depends=1 (and no GC crash).
 */
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
  uint32_t ref_count;
  uint32_t child_count;
  uint16_t symbol;
  uint16_t flags;
} SubtreeHeapData;

enum {
  FLAG_DEPENDS_ON_COLUMN = 1u << 8,
};

typedef struct {
  bool is_inline : 1;
  bool visible : 1;
  bool named : 1;
  bool extra : 1;
  bool has_changes : 1;
  bool is_missing : 1;
  bool is_keyword : 1;
  uint8_t symbol;
  uint16_t parse_state;
  uint8_t padding_columns;
  uint8_t padding_rows : 4;
  uint8_t lookahead_bytes : 4;
  uint8_t padding_bytes;
  uint8_t size_bytes;
} SubtreeInlineData;

typedef union {
  SubtreeInlineData data;
  const SubtreeHeapData *ptr;
} Subtree;

static bool ts_subtree_depends_on_column(Subtree self) {
  if (self.data.is_inline) {
    return false;
  }
  return (self.ptr->flags & FLAG_DEPENDS_ON_COLUMN) != 0;
}

static Subtree make_inline_leaf(uint8_t symbol) {
  Subtree s;
  memset(&s, 0, sizeof s);
  s.data.is_inline = true;
  s.data.visible = true;
  s.data.named = true;
  s.data.symbol = symbol;
  s.data.parse_state = 1;
  s.data.size_bytes = 1;
  return s;
}

static Subtree make_heap_node(bool depends) {
  SubtreeHeapData *h = calloc(1, sizeof *h);
  h->ref_count = 1;
  h->child_count = 2;
  h->symbol = 42;
  if (depends) {
    h->flags |= FLAG_DEPENDS_ON_COLUMN;
  }
  Subtree s;
  s.ptr = h;
  return s;
}

int main(void) {
  Subtree leaf = make_inline_leaf(7);
  Subtree parent = make_heap_node(true);

  printf("sizeof(Subtree)=%zu (should be pointer-sized)\n", sizeof(Subtree));
  printf("inline bit=%d  raw_word=%p\n",
         leaf.data.is_inline, (void *)(uintptr_t)leaf.ptr);

  printf("inline depends=%d (want 0)\n", ts_subtree_depends_on_column(leaf));
  printf("heap   depends=%d (want 1)\n", ts_subtree_depends_on_column(parent));

  free((void *)parent.ptr);
  return 0;
}
