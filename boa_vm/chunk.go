package main

type Opcode uint8

const (
  OpConstant Opcode = iota
  OpReturn
  OpNegate
  OpAdd
  OpSub
  OpDiv
  OpMul
)

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

