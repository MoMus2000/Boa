package main

import (
	"errors"
	"fmt"
	"unsafe"
)

type InterpretResult int

const DEBUG_TRACE_EXECUTION = 0
const STACK_MAX = 256
const FRAME_MAX = 64

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_RUNTIME_ERROR
	INTERPRET_COMPILE_ERROR
)

type VM struct {
	stackTop     Value
	stack        []Value
	compiler     Compiler
	table        *Table
	frameCount   int
	frames       []CallFrame
	currentFrame *CallFrame
}

type CallFrame struct {
	function   *ObjectFunc
	ip         int
  slots      []Value
	code       []Opcode
}

func NewVM() VM {
	return VM{
		compiler:   NewCompiler(),
		table:      initMap(),
		frames:     make([]CallFrame, FRAME_MAX),
		frameCount: 0,
	}
}

func (v *VM) resetStack() {
	v.stackTop = NumberVal(0)
	v.stack = make([]Value, 0)
}

func (v *VM) FreeVM() {
}

func (v *VM) interpret(source []byte) InterpretResult {
	compiled_code := v.compiler.compile(source)
	if compiled_code == nil {
		return INTERPRET_COMPILE_ERROR
	}
	obj := (*Object)(unsafe.Pointer(compiled_code))

	va := ObjVal(obj)
	v.push(va)

	v.call(compiled_code, 0)

	return v.run()
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

func (v *VM) run() InterpretResult {
	for {
		currentFrame := &v.frames[v.frameCount-1]
		v.currentFrame = currentFrame
		ins := v.currentFrame.code[currentFrame.ip]
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
        result := v.pop()

        v.frameCount --

        if v.frameCount == 0 {
          return INTERPRET_OK
        }

        currentFrame := &v.frames[v.frameCount-1]
        v.push(*result)
        v.currentFrame = currentFrame

        break
			}
		case OpConstant:
			{
				c := v.read_constant()
				v.push(*c) // Push HERE into the stack
				// fmt.Printf("\n")
				// fmt.Println("------------")
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
					fmt.Printf("\n")
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
        slot := v.currentFrame.code[v.currentFrame.ip]
				v.currentFrame.ip++

				v.push((v.currentFrame.slots)[slot])
				break
			}
		case OpSetLocal:
			{
				slot := v.currentFrame.code[v.currentFrame.ip]
				v.currentFrame.ip++
				(v.currentFrame.slots)[slot] = *v.peek(0)
				break
			}
		case OpJumpIfFalse:
			{
				offset := v.read_short()
				if v.isFalsy(v.peek(0)) {
					v.currentFrame.ip += int(offset)
				}
				break
			}
		case OpJump:
			{
				offset := v.read_short()
				v.currentFrame.ip += int(offset)
				break
			}
		case OpLoop:
			{
				offset := v.read_short()
				v.currentFrame.ip -= int(offset)
				break
			}
		case OpCall:
			{
				argCount := int(v.currentFrame.function.chunk.code[v.currentFrame.ip])
				v.currentFrame.ip++
				if !v.callValue(v.peek(argCount), argCount) {
					return INTERPRET_RUNTIME_ERROR
				}
        v.pop()
				break
			}
		default:
			fmt.Println("UnIdentified OpCode: ", ins)
		}
	}
}

func (v *VM) callValue(val *Value, count int) bool {
	if val.IsObj() {
		switch val.obj.objType {
		case OBJ_FUNC:
			return v.call(val.asFunc(), count)
		default:
			break // Non-callable object type.
		}
	}
	v.compiler.error("Can only call functions and classes.")
	return false
}

func (v *VM) call(obj *ObjectFunc, count int) bool {
  if obj.arity != count {
    fmt.Println("Args don't match up")
  }

  fmt.Println("-----FuncOp-----")
  obj.chunk.printOpCode()
  fmt.Println("------------")

  fmt.Println("--Current Stack--")
  fmt.Printf("-------\n")
  for _, val := range v.stack{
    if val.isString(){
      fmt.Printf("| str  | %v\n", val.asString().chars)
    } else if val.isFunc(){
      fmt.Printf("| func | %v\n", val.asFunc().name.chars)
    } else if val.IsNumber(){
      fmt.Printf("| Num  |\n")
    }else{
      fmt.Printf("| NaN  |\n")
    }
  }
  fmt.Printf("-------\n")

	frame := &v.frames[v.frameCount]
	v.currentFrame = frame
	v.frameCount++
	v.currentFrame.function = obj
	v.currentFrame.ip = 0
	v.currentFrame.code = obj.chunk.code
  v.currentFrame.slots = v.stack
	return true
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

func (v *VM) read_byte() {
	v.currentFrame.ip++
}

func (v *VM) read_constant() *Value {
	index := v.currentFrame.function.chunk.code[v.currentFrame.ip]
	c := v.currentFrame.function.chunk.constants.values[index]
	v.read_byte()
	return &c
}

func (v *VM) read_short() uint16 {
	v.currentFrame.ip += 2
	v1 := uint16(v.currentFrame.function.chunk.code[v.currentFrame.ip-1])
	v2 := uint16(v.currentFrame.function.chunk.code[v.currentFrame.ip-2]) << 8
	short := (uint16)(v2 | v1)
	return short
}

func (v *VM) read_string() *ObjectString {
	return v.read_constant().asString()
}

