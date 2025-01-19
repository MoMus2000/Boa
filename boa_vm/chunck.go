package main

type Opcode uint8

const (
  OpConstant Opcode = iota
  OpReturn   Opcode = iota
)

type Chunck struct{
  code []Opcode
  constants ValueArray
}

func NewChunck() Chunck{
  v := NewValue()
  return Chunck{
    code : make([]Opcode, 0),
    constants: v,
  }
}

func (c *Chunck) WriteChunck(b Opcode){
  c.code = append(c.code, b)
}

func (c *Chunck) AddConstant(value Value) int{
  c.constants.WriteValue(value)
  return len(c.constants.values) - 1
}

func (c *Chunck) FreeChunk(){
  c.code = make([]Opcode, 0)
  c.constants.FreeValue()
  c.constants = NewValue()
}

