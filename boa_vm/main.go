package main

func main(){
  vm := NewVM();
  c := NewChunck()
  constant := c.AddConstant(1.2)
  c.WriteChunk(OpConstant, 123)
  c.WriteChunk(Opcode(constant), 123)
  constant = c.AddConstant(1.69)
  c.WriteChunk(OpConstant, 124)
  c.WriteChunk(Opcode(constant), 124)
  c.WriteChunk(OpNegate, 125)
  c.WriteChunk(OpAdd, 125)
  constant = c.AddConstant(1.69)
  c.WriteChunk(OpConstant, 124)
  c.WriteChunk(Opcode(constant), 124)
  c.WriteChunk(OpMul, 125)
  c.WriteChunk(OpReturn, 125)
  vm.interpret(&c)
  c.FreeChunk()
  vm.FreeVM()
  c.FreeChunk()
}

