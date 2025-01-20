package main

func main(){
  c := NewChunck()
  constant := c.AddConstant(1.2)
  c.WriteChunk(OpConstant, 123)
  c.WriteChunk(Opcode(constant), 123)
  c.WriteChunk(OpReturn, 123)
  DisassembleChunk(&c, "TEST CHUNK")
  c.FreeChunk()
}

