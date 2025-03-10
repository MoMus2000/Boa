#include <stdio.h>
#include "../include/common.h"


int main(int argc, char* argv[]){
  if (argc != 2) {
    printf("ERROR: Usage: %s <filename>", argv[0]);
    return -1;
  }
  char *fileName = argv[1];

  printf("Loading fileName: %s\n", fileName);

  return 0;
}

