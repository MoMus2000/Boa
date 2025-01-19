package main

import (
  "fmt"
)

func DisassembleChunk(chunck *Chunck, message string){
  fmt.Printf(" == %s == \n", message)
  offset := 0
  for offset < len(chunck.code){
    offset = DisassembleInstruction(chunck, offset)
  }
}

func DisassembleInstruction(c *Chunck, offset int) int{
  fmt.Printf("%04d ", offset)
  instruction := c.code[offset]
  switch instruction {
    case OpReturn : {
      return SimpleInstruction("OP_RETURN", offset)
    }
    case OpConstant : {
      return ConstantInstruction("OP_CONSTANT", c, offset)
    }
    default: {
      fmt.Printf("Unknown OpCode %d\n", instruction)
      return int(offset + 1)
    }
  }
}

func ConstantInstruction(ins string, chunk *Chunck, offset int) int {
  constant := chunk.code[offset + 1]
  fmt.Printf("%-16s %4d", ins, constant)
  printValue(chunk.constants.values[constant])
  fmt.Println("")
  return offset + 2
}

func printValue(v Value){
  fmt.Printf("C: %v\n", v)
}

func SimpleInstruction(ins string, offset int) int {
  fmt.Println(ins)
  return offset + 1
}

