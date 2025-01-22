package main

import (
	"fmt"
	"os"
)

func DisassembleChunk(chunck *Chunk, message string){
  fmt.Printf(" == %s == \n", message)
  offset := 0
  for offset < len(chunck.code){
    offset = DisassembleInstruction(chunck, offset)
  }
}

func DumpByteCode(c *Chunk){
  file, err := os.Create("./bytecode/main.boac")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write the struct to the file
	_, err = fmt.Fprintf(file, "%+v\n", c)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func DisassembleInstruction(c *Chunk, offset int) int{
  fmt.Printf("%04d ", offset)
  if offset > 0 && c.lines[offset] == c.lines[offset - 1]{
    fmt.Printf("   | ")
  }else{
    fmt.Printf("%04d ", c.lines[offset])
  }
  instruction := c.code[offset]
  switch instruction {
    case OpNegate: {
      return SimpleInstruction("OP_NEGATE", offset)
    }
    case OpReturn : {
      return SimpleInstruction("OP_RETURN", offset)
    }
    case OpConstant : {
      return ConstantInstruction("OP_CONSTANT", c, offset)
    }
    case OpAdd : {
      return SimpleInstruction("OP_ADD", offset)
    }
    case OpSub: {
      return SimpleInstruction("OP_SUB", offset)
    }
    case OpMul: {
      return SimpleInstruction("OP_MUL", offset)
    }
    case OpDiv: {
      return SimpleInstruction("OP_DIV", offset)
    }
    default: {
      fmt.Printf("Unknown OpCode %d\n", instruction)
      return int(offset + 1)
    }
  }
}

func ConstantInstruction(ins string, chunk *Chunk, offset int) int {
  constant := chunk.code[offset + 1]
  fmt.Printf("%-16s %4d", ins, constant)
  printValue(chunk.constants.values[constant])
  fmt.Println("")
  return offset + 2
}

func printValue(v Value){
  fmt.Printf(" '%v'\n", v)
}

func SimpleInstruction(ins string, offset int) int {
  fmt.Println(ins)
  return offset + 1
}

