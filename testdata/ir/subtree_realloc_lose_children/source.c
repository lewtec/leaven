/* subtree_realloc_lose_children — realloc must copy old children[]
 * tree-sitter new_node grows the block; kids live at the start of the allocation.
 * clang-14 -S -emit-llvm -fno-discard-value-names -std=gnu11 -O0 -o input.ll source.c
 */
/*
 * REPRO: leaven libc.Realloc does not copy the old block.
 *
 * tree-sitter ts_subtree_new_node (lib/src/subtree.c):
 *   if (capacity * sizeof(Subtree) < n*sizeof(Subtree) + sizeof(HeapData))
 *     contents = realloc(contents, new_size);  // must preserve children
 *   HeapData *data = &contents[size];
 *
 * Native realloc: copies children → OK.
 * leaven Realloc: fresh allocation, no copy → children become 0
 *   → summarize_children → depends_on_column(0) → nil deref (offset ~0x2c).
 *
 * Earlier "children-before-heap" fixtures used one calloc of the final size
 * (no realloc) and therefore did NOT reproduce.
 *
 * Build native:
 *   clang -O0 -o /tmp/srl /tmp/subtree-realloc-lose-children.c && /tmp/srl
 *   # expect: child[0] inline=1, exit 0
 *
 * leaven:
 *   clang -emit-llvm -S -O0 -o /tmp/srl.ll /tmp/subtree-realloc-lose-children.c
 *   go tool leaven -package main /tmp/srl.ll > /tmp/srl.go
 *   # go run with leaven replace → panic or child raw=0x0
 */

#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* Large enough that 1*sizeof(Subtree) + sizeof(Heap) > small capacity. */
typedef struct {
  uint32_t ref_count;
  uint32_t child_count;
  uint16_t symbol;
  uint16_t flags;
  uint32_t size_row;
  uint8_t pad[64]; /* pad so heap header forces realloc when capacity==1 */
} SubtreeHeapData;

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

typedef struct {
  Subtree *contents;
  uint32_t size;
  uint32_t capacity;
} SubtreeArray;

static bool depends_on_column(Subtree self) {
  if (self.data.is_inline)
    return false;
  /* heap path — null/0 ptr crashes here (matches core stack) */
  return (self.ptr->flags & 0x100) != 0;
}

static size_t alloc_size(uint32_t n) {
  return (size_t)n * sizeof(Subtree) + sizeof(SubtreeHeapData);
}

static void *ts_realloc(void *p, size_t n) {
  void *q = realloc(p, n);
  if (!q && n)
    abort();
  return q;
}

static Subtree make_inline(uint8_t sym) {
  Subtree s;
  memset(&s, 0, sizeof s);
  s.data.is_inline = true;
  s.data.visible = true;
  s.data.named = true;
  s.data.symbol = sym;
  s.data.parse_state = 1;
  s.data.size_bytes = 1;
  return s;
}

/* Mirrors ts_subtree_new_node + depends_on_column loop in summarize_children. */
static Subtree new_node(SubtreeArray *children) {
  size_t nbytes = alloc_size(children->size);
  size_t have = (size_t)children->capacity * sizeof(Subtree);

  printf("before realloc: size=%u cap=%u need=%zu have=%zu will_realloc=%d\n",
         children->size, children->capacity, nbytes, have, have < nbytes);

  if (have < nbytes) {
    children->contents = ts_realloc(children->contents, nbytes);
    children->capacity = (uint32_t)(nbytes / sizeof(Subtree));
  }

  SubtreeHeapData *data =
      (SubtreeHeapData *)&children->contents[children->size];
  memset(data, 0, sizeof *data);
  data->ref_count = 1;
  data->child_count = children->size;
  data->symbol = 100;
  data->size_row = 0;

  Subtree self;
  self.ptr = data;

  Subtree *ch = (Subtree *)self.ptr - self.ptr->child_count;
  for (uint32_t i = 0; i < self.ptr->child_count; i++) {
    Subtree child = ch[i];
    printf("child[%u] raw=%p inline=%d\n", i,
           (void *)(uintptr_t)child.ptr, (int)child.data.is_inline);
    if (self.ptr->size_row == 0 && depends_on_column(child))
      self.ptr->flags |= 0x100;
  }
  return self;
}

int main(void) {
  printf("sizeof(Subtree)=%zu sizeof(Heap)=%zu\n",
         sizeof(Subtree), sizeof(SubtreeHeapData));

  SubtreeArray arr = {0};
  /* capacity 1: only room for the leaf, NOT leaf+heap header → must realloc */
  arr.capacity = 1;
  arr.contents = calloc(arr.capacity, sizeof(Subtree));
  arr.contents[0] = make_inline(7);
  arr.size = 1;

  printf("leaf raw=%p inline=%d\n",
         (void *)(uintptr_t)arr.contents[0].ptr,
         (int)arr.contents[0].data.is_inline);

  Subtree p = new_node(&arr);
  printf("ok child_count=%u\n", p.ptr->child_count);
  free(arr.contents);
  return 0;
}
