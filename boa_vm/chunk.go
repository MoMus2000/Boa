package main

import "fmt"

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
  "OpMinus1",
}

// String method to print enum name
func (t Opcode) String() string {
	if t < OpConstant || t > OpMinus1 {
		return fmt.Sprintf("Unknown TokenType(%d)", t)
	}
	return opCodeNames[t]
}

type Chunk struct{
  code []Opcode
  constants ValueArray
  lines []int
}

func NewChunck() Chunk{
  v := NewValue()
  return Chunk{
    code : make([]Opcode, 0),
    constants: v,
    lines: make([]int, 0),
  }
}

func (c *Chunk) WriteChunk(b Opcode, line int){
  c.code = append(c.code, b)
  c.lines = append(c.lines, line)
}

func (c *Chunk) AddConstant(value Value) int{
  c.constants.WriteValue(value)
  return len(c.constants.values) - 1
}

func (c *Chunk) FreeChunk(){
  c.code = make([]Opcode, 0)
  c.constants.FreeValue()
  c.constants = NewValue()
  c.lines = make([]int, 0)
}

