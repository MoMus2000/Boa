package main

import (
	"errors"
	"fmt"
	"unsafe"
)

type InterpretResult int

const DEBUG_TRACE_EXECUTION = 0
const STACK_MAX = 256

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_RUNTIME_ERROR
	INTERPRET_COMPILE_ERROR
)

type VM struct {
	chunk    *Chunk
	ip       int
	stackTop Value
	stack    []Value
	compiler Compiler
	table    *Table
}

func NewVM() VM {
	return VM{
		compiler: NewCompiler(),
		table:    initMap(),
	}
}

func (v *VM) resetStack() {
	v.stackTop = NumberVal(0)
	v.stack = make([]Value, 0)
}

func (v *VM) FreeVM() {
}

func (v *VM) interpret(source []byte) InterpretResult {
	chunk := NewChunck()
	if !v.compiler.compile(source, &chunk) {
		return INTERPRET_COMPILE_ERROR
	}
	v.chunk = &chunk
	v.ip = 0
	result := v.run()
	chunk.FreeChunk()
	return result
}

func (v *VM) push(vl Value) {
	v.stack = append(v.stack, vl)
	v.stackTop.number++
}

func (v *VM) pop() *Value {
	last := v.stack[len(v.stack)-1]
	v.stack = v.stack[:len(v.stack)-1]
	v.stackTop.number--
	return &last
}

func (v *VM) read_byte() {
	v.ip++
}

func (v *VM) run() InterpretResult {
	for {
		ins := v.chunk.code[v.ip]
		// fmt.Println("Remaining OpCodes: ", v.chunk.code[v.ip:])
		// fmt.Println("Current Instruction: ", ins)
		v.read_byte()
		switch ins {
		case OpPrint:
			{
				c := v.pop()
				printValue(*c)
				fmt.Printf("\n")
				break
			}
		case OpReturn:
			{
				return INTERPRET_OK
			}
		case OpConstant:
			{
				c := v.read_constant()
				v.push(*c)
				// fmt.Printf("Constant: ")
				// printValue(*c)
				// fmt.Printf("\n")
				break
			}
		case OpNegate:
			{
				if !v.peek(0).IsNumber() {
					// runtimeError("Operand must be a number")
					return INTERPRET_RUNTIME_ERROR
				}
				c := v.pop().AsNumber()
				d := NumberVal(NumberVal(-1).number * c)
				v.push(d)
				break
			}
		case OpAdd:
			{
				if v.peek(0).isString() && v.peek(1).isString() {
					v.concatenate()
					break
				}
				err := v.binary_op("+")
				if err != nil {
					return INTERPRET_RUNTIME_ERROR
				}
				break
			}
		case OpMul:
			{
				err := v.binary_op("*")
				if err != nil {
					return INTERPRET_RUNTIME_ERROR
				}
				break
			}
		case OpSub:
			{
				err := v.binary_op("-")
				if err != nil {
					return INTERPRET_RUNTIME_ERROR
				}
				break
			}
		case OpDiv:
			{
				err := v.binary_op("/")
				if err != nil {
					return INTERPRET_RUNTIME_ERROR
				}
				break
			}
		case OpFalse:
			{
				v.push(BoolVal(false))
				break
			}
		case OpTrue:
			{
				v.push(BoolVal(true))
				break
			}
		case OpNil:
			{
				v.push(NilVal())
				break
			}
		case OpNot:
			{
				pred := v.pop()
				bool_val := v.isFalsy(pred)
				v.push(BoolVal(bool_val))
				break
			}
		case OpGreater:
			{
				err := v.binary_op(">")
				if err != nil {
					return INTERPRET_RUNTIME_ERROR
				}
				break
			}
		case OpLess:
			{
				err := v.binary_op("<")
				if err != nil {
					return INTERPRET_RUNTIME_ERROR
				}
				break
			}
		case OpEqual:
			{
				v.push(BoolVal(v.valuesEqual(v.pop(), v.pop())))
				break
			}
		case OpPop:
			{
				v.pop()
				break
			}
		case OpDefineGlobal:
			{
				value := v.peek(0)
				name := string(v.read_string().chars)
				v.table.tableSet(name, *value)
				v.pop()
				break
			}
		case OpGetGlobal:
			{
				name := string(v.read_string().chars)
				c := v.table.tableGet(name)
				if c != nil {
					// printValue(*c)
					// fmt.Printf("\n")
					v.push(*c)
				} else {
					v.push(NilVal())
				}
				break
			}
		case OpSetGlobal:
			{
				name := string(v.read_string().chars)
				ok := v.table.tableGet(name)
				if ok == nil {
					return INTERPRET_RUNTIME_ERROR
				}
				value := v.peek(0)
				v.table.tableDelete(name)
				v.table.tableSet(name, *value)
				break
			}
		case OpGetLocal:
			{
				index := v.chunk.code[v.ip]
				v.ip++
				v.push(v.stack[index])
				break
			}
		case OpSetLocal:
			{
				index := v.chunk.code[v.ip]
				v.ip++
				v.stack[index] = *v.peek(0)
				break
			}
		case OpJumpIfFalse:
			{
				offset := v.read_short()
				if v.isFalsy(v.peek(0)) {
					v.ip += int(offset)
				}
				break
			}
		case OpJump:
			{
				offset := v.read_short()
				v.ip += int(offset)
				break
			}
		case OpLoop:
			{
				offset := v.read_short()
				v.ip -= int(offset)
				break
			}
		default:
			fmt.Println("UnIdentified OpCode: ", ins)
		}
	}
}

func (v *VM) concatenate() {
	s1 := v.pop().asString().chars
	s1 = s1[1 : len(s1)-1]
	s2 := v.pop().asString().chars
	s2 = s2[1 : len(s2)-1]
	s3 := fmt.Sprintf("\"%s\"", s2+s1)
	objstr := ObjectString{
		obj:    Object{objType: OBJ_STRING},
		chars:  s3,
		length: len(s3),
	}
	obj := (*Object)(unsafe.Pointer(&objstr))
	v.push(ObjVal(obj))
}

func (v *VM) valuesEqual(v1 *Value, v2 *Value) bool {
	if v1.valType != v2.valType {
		return false
	}
	switch v1.valType {
	case VAL_NUMBER:
		{
			return v1.AsNumber() == v2.AsNumber()
		}
	case VAL_BOOL:
		{
			return v1.AsBoolean() == v2.AsBoolean()
		}
	case VAL_NIL:
		{
			return true
		}
	case VAL_OBJ:
		{
			str1 := v1.asString().chars
			str2 := v2.asString().chars
			return str1 == str2
		}
	default:
		return false
	}
}

func (v *VM) isFalsy(p *Value) bool {
	return p.IsNil() || (p.IsBool() && !p.AsBoolean())
}

func (v *VM) peek(distance int) *Value {
	return &v.stack[len(v.stack)-1-distance]
}

func (v *VM) binary_op(op string) error {
	if !v.peek(0).IsNumber() || !v.peek(1).IsNumber() {
		return errors.New("Runtime Error")
	}
	a := v.pop().AsNumber()
	b := v.pop().AsNumber()
	switch op {
	case "-":
		{
			v.push(NumberVal(b - a))
		}
	case "+":
		{
			v.push(NumberVal(b + a))
		}
	case "*":
		{
			v.push(NumberVal(b * a))
		}
	case "/":
		{
			v.push(NumberVal(b / a))
		}
	case ">":
		{
			v.push(BoolVal(a < b))
		}
	case "<":
		{
			v.push(BoolVal(a > b))
		}
	default:
		{
			return errors.New("Runtime Error")
		}
	}
	return nil
}

func (v *VM) read_constant() *Value {
	index := v.chunk.code[v.ip]
	c := v.chunk.constants.values[index]
	v.read_byte()
	return &c
}

func (v *VM) read_short() uint16 {
	v.ip += 2
	v1 := uint16(v.chunk.code[v.ip-1])
	v2 := uint16(v.chunk.code[v.ip-2]) << 8
	short := (uint16)(v2 | v1)
	return short
}

func (v *VM) read_string() *ObjectString {
	return v.read_constant().asString()
}
