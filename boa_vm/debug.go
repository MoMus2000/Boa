package main

import (
	"fmt"
	"os"
)

func DisassembleChunk(chunck *Chunk, message string) {
	fmt.Printf(" == %s == \n", message)
	offset := 0
	for offset < len(chunck.code) {
		offset = DisassembleInstruction(chunck, offset)
	}
}

func DumpByteCode(c *Chunk) {
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

func DisassembleInstruction(c *Chunk, offset int) int {
	fmt.Printf("%04d ", offset)
	if offset > 0 && c.lines[offset] == c.lines[offset-1] {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%04d ", c.lines[offset])
	}
	instruction := c.code[offset]
	switch instruction {
	case OpNegate:
		{
			return SimpleInstruction("OP_NEGATE", offset)
		}
	case OpReturn:
		{
			return SimpleInstruction("OP_RETURN", offset)
		}
	case OpConstant:
		{
			return ConstantInstruction("OP_CONSTANT", c, offset)
		}
	case OpAdd:
		{
			return SimpleInstruction("OP_ADD", offset)
		}
	case OpSub:
		{
			return SimpleInstruction("OP_SUB", offset)
		}
	case OpMul:
		{
			return SimpleInstruction("OP_MUL", offset)
		}
	case OpDiv:
		{
			return SimpleInstruction("OP_DIV", offset)
		}
	case OpTrue:
		{
			return SimpleInstruction("OP_TRUE", offset)
		}
	case OpFalse:
		{
			return SimpleInstruction("OP_FALSE", offset)
		}
	case OpNil:
		{
			return SimpleInstruction("OP_NIL", offset)
		}
	case OpNot:
		{
			return SimpleInstruction("OP_NOT", offset)
		}
	case OpEqual:
		{
			return SimpleInstruction("OP_EQUAL", offset)
		}
	case OpGreater:
		{
			return SimpleInstruction("OP_GREATER", offset)
		}
	case OpLess:
		{
			return SimpleInstruction("OP_LESS", offset)
		}
	case OpPrint:
		{
			return SimpleInstruction("OP_PRINT", offset)
		}
	case OpPop:
		{
			return SimpleInstruction("OP_POP", offset)
		}
	case OpDefineGlobal:
		{
			return AssignmentInstruction("OP_DEFINE_GLOBAL", c, offset)
		}
	case OpGetGlobal:
		{
			return SimpleInstruction("OP_GET_GLOBAL", offset)
		}
	case OpSetGlobal:
		{
			return AssignmentInstruction("OP_SET_GLOBAL", c, offset)
		}
	case OpSetLocal:
		{
			return ByteInstruction("OP_SET_LOCAL", c, offset)
		}
	case OpGetLocal:
		{
			return ByteInstruction("OP_GET_LOCAL", c, offset)
		}
	case OpJump:
		{
			return JumpInstruction("OP_JUMP", 1, offset, c)
		}
	case OpJumpIfFalse:
		{
			return JumpInstruction("OP_JUMP_IF_FALSE", 1, offset, c)
		}
	case OpLoop:
		{
			return JumpInstruction("OP_LOOP", -1, offset, c)
		}
	default:
		{
			fmt.Printf("Unknown OpCode %d\n", instruction)
			return int(offset + 1)
		}
	}
}

func ByteInstruction(ins string, c *Chunk, offset int) int {
	constant := c.code[offset+1]
	fmt.Printf("%-16s %4d", ins, constant)
	return offset + 2
}

func JumpInstruction(ins string, sign int, offset int, c *Chunk) int {
	j1 := uint16(c.code[offset+1]) << 8
	j2 := uint16(c.code[offset+2])
	j1 |= j2
	fmt.Printf("%-16s %4d -> %d\n", ins, offset, j1)
	return offset + 3
}

func AssignmentInstruction(ins string, chunk *Chunk, offset int) int {
	ident_index := chunk.code[offset+1]
	ident := chunk.constants.values[ident_index].asCString()
	fmt.Printf("%-16s %s", ins, ident)
	fmt.Println("")
	return offset + 2
}

func ConstantInstruction(ins string, chunk *Chunk, offset int) int {
	constant := chunk.code[offset+1]
	fmt.Printf("%-16s %4d", ins, constant)
	printValue(chunk.constants.values[constant])
	fmt.Println("")
	return offset + 2
}

func printObject(v Value) {
	switch v.AsObj().ObjType() {
	case OBJ_STRING:
		{
			fmt.Printf("%v", v.asString().chars)
		}
	default:
		fmt.Println("Print for Object not implemented")

	}
}

func printValue(v Value) {
	switch v.valType {
	case VAL_BOOL:
		{
			fmt.Printf(" '%v'\n", v.AsBoolean())
		}
	case VAL_NUMBER:
		{
			fmt.Printf(" '%v'\n", v.AsNumber())
		}
	case VAL_NIL:
		{
			fmt.Printf(" '%v'\n", nil)
		}
	case VAL_OBJ:
		{
			printObject(v)
		}
	}
}

func SimpleInstruction(ins string, offset int) int {
	fmt.Println(ins)
	return offset + 1
}
