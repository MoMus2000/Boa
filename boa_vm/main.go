package main

func main(){
  c := NewChunck()
  c.WriteChunck(OpReturn)
  constant := c.AddConstant(1.69)
  c.WriteChunck(OpConstant)
  c.WriteChunck(Opcode(constant))
  constant = c.AddConstant(1.69)
  c.WriteChunck(OpConstant)
  c.WriteChunck(Opcode(constant))
  constant = c.AddConstant(1.69)
  c.WriteChunck(OpConstant)
  c.WriteChunck(Opcode(constant))
  DisassembleChunk(&c, "TEST CHUNK")
  c.FreeChunk()
}

