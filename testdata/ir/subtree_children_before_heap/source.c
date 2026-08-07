/* subtree_children_before_heap — children[] laid out before SubtreeHeapData
 *   [child0..childN-1][HeapData]  ptr → HeapData, children = (Subtree*)ptr - N
 * clang-14 -S -emit-llvm -fno-discard-value-names -std=gnu11 -O0 -o input.ll source.c
 * Native/leaven: child[0] is_inline=1, no panic in depends_on_column.
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
  uint32_t size_row;
} SubtreeHeapData;

enum { FLAG_DEPENDS_ON_COLUMN = 1u << 8 };

typedef struct {
  bool is_inline : 1;
  bool visible : 1;
  bool named : 1;
  uint8_t symbol;
  uint16_t parse_state;
  uint8_t size_bytes;
} SubtreeInlineData;

typedef union {
  SubtreeInlineData data;
  SubtreeHeapData *ptr;
} Subtree;

static bool depends_on_column(Subtree self) {
  if (self.data.is_inline)
    return false;
  return (self.ptr->flags & FLAG_DEPENDS_ON_COLUMN) != 0;
}

static Subtree *children_of(Subtree self) {
  if (self.data.is_inline)
    return NULL;
  return (Subtree *)self.ptr - self.ptr->child_count;
}

static size_t alloc_size(uint32_t n) {
  return (size_t)n * sizeof(Subtree) + sizeof(SubtreeHeapData);
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

static Subtree new_node_with_children(Subtree *kids, uint32_t n) {
  void *block = calloc(1, alloc_size(n));
  if (!block)
    abort();

  Subtree *slot = (Subtree *)block;
  for (uint32_t i = 0; i < n; i++)
    slot[i] = kids[i];

  SubtreeHeapData *data = (SubtreeHeapData *)&slot[n];
  data->ref_count = 1;
  data->child_count = n;
  data->symbol = 100;
  data->flags = 0;
  data->size_row = 0;

  Subtree parent;
  parent.ptr = data;

  Subtree *ch = children_of(parent);
  for (uint32_t i = 0; i < parent.ptr->child_count; i++) {
    Subtree child = ch[i];
    fprintf(stderr, "child[%u] raw=%p is_inline=%d\n",
            i, (void *)(uintptr_t)child.ptr, (int)child.data.is_inline);
    if (parent.ptr->size_row == 0 && depends_on_column(child))
      parent.ptr->flags |= FLAG_DEPENDS_ON_COLUMN;
  }
  return parent;
}

int main(void) {
  printf("sizeof(Subtree)=%zu (want 8 on 64-bit)\n", sizeof(Subtree));
  Subtree leaf = make_inline_leaf(7);
  printf("leaf raw=%p is_inline=%d\n",
         (void *)(uintptr_t)leaf.ptr, (int)leaf.data.is_inline);
  Subtree kids[1] = {leaf};
  Subtree parent = new_node_with_children(kids, 1);
  printf("parent child_count=%u depends=%d\n",
         parent.ptr->child_count,
         (parent.ptr->flags & FLAG_DEPENDS_ON_COLUMN) != 0);
  free((Subtree *)parent.ptr - parent.ptr->child_count);
  return 0;
}
