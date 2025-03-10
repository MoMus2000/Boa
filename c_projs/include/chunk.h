#ifndef clox_chunk_h
#define clox_chunk_h
#include "common.h"
#include "memory.h"

typedef enum {

  OP_RETURN

} OpCode;

typedef struct {
  int count;
  int capacity;
  uint8_t *code;
} Chunk;

void initChunk(Chunk* chunk);
void writeChunk(Chunk* chunk, uint8_t byte);

#endif
