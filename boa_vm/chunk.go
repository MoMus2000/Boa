package main

import (
	"fmt"
	"strings"
)

type Opcode uint8

const (
	OpConstant Opcode = iota
	OpReturn
	OpNegate
	OpAdd
	OpSub
	OpDiv
	OpMul
	OpNil
	OpTrue
	OpFalse
	OpNot
	OpEqual
	OpLess
	OpGreater
	OpPrint
	OpPop
	OpDefineGlobal
	OpGetGlobal
	OpSetGlobal
	OpSetLocal
	OpGetLocal
	OpJumpIfFalse
	OpJump
	OpLoop
	OpCall
	OpMinus1
)

var opCodeNames = [...]string{
	"OpConstant",
	"OpReturn",
	"OpNegate",
	"OpAdd",
	"OpSub",
	"OpDiv",
	"OpMul",
	"OpNil",
	"OpTrue",
	"OpFalse",
	"OpNot",
	"OpEqual",
	"OpLess",
	"OpGreater",
	"OpPrint",
	"OpPop",
	"OpDefineGlobal",
	"OpGetGlobal",
	"OpSetGlobal",
	"OpSetLocal",
	"OpGetLocal",
	"OpJumpIfFalse",
	"OpJump",
	"OpLoop",
	"OpCall",
	"OpMinus1",
}

// String method to print enum name
func (t Opcode) String() string {
	if t < OpConstant || t > OpMinus1 {
		return fmt.Sprintf("Unknown TokenType(%d)", t)
	}
	return opCodeNames[t]
}

type Chunk struct {
	code      []Opcode
	constants ValueArray
	lines     []int
}

func NewChunck() Chunk {
	v := NewValue()
	return Chunk{
		code:      make([]Opcode, 0),
		constants: v,
		lines:     make([]int, 0),
	}
}

func (c *Chunk) printOpCode() {
  var result []string

  for i := 0; i < len(c.code); i++ {
    opcode := c.code[i]
    entry := fmt.Sprintf("%v", opcode)

    // Check if the opcode has an operand (requires next index)
    switch opcode {
    case OpConstant, OpDefineGlobal, OpGetGlobal, OpSetGlobal, OpSetLocal, OpGetLocal, OpGreater, OpLess, OpCall:
      if i+1 < len(c.code) {
        entry += fmt.Sprintf(" index: %d", int(c.code[i+1]))
        i++ // Skip the operand index
      }
    }

    result = append(result, entry)
  }

  // Print the entire result array in one line
  fmt.Println("[", strings.Join(result, ", "), "]")
}


func (c *Chunk) WriteChunk(b Opcode, line int) {
	c.code = append(c.code, b)
	c.lines = append(c.lines, line)
}

func (c *Chunk) AddConstant(value Value) int {
	c.constants.WriteValue(value)
	return len(c.constants.values) - 1
}

func (c *Chunk) FreeChunk() {
	c.code = make([]Opcode, 0)
	c.constants.FreeValue()
	c.constants = NewValue()
	c.lines = make([]int, 0)
}
