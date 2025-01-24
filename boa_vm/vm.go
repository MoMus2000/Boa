package main

import (
	"errors"
	"fmt"
)

type InterpretResult int;


const DEBUG_TRACE_EXECUTION = 1;
const STACK_MAX = 256

const (
  INTERPRET_OK              InterpretResult = iota
  INTERPRET_RUNTIME_ERROR
  INTERPRET_COMPILE_ERROR
)

type VM struct {
  chunk    *Chunk
  ip       int
  stackTop Value
  stack    []Value
  compiler Compiler
}

func NewVM() VM {
  return VM{
    compiler : NewCompiler(),
  }
}

func (v *VM) resetStack() {
  v.stackTop = NumberVal(0)
  v.stack    = make([]Value, 0)
}

func (v *VM) FreeVM() {
}

func (v *VM) interpret(source []byte) InterpretResult{
  chunk := NewChunck()
  if ! v.compiler.compile(source, &chunk) {
    return INTERPRET_COMPILE_ERROR
  }
  v.chunk = &chunk
  // DumpByteCode(v.chunk)
  v.ip    = 0
  result := v.run()
  chunk.FreeChunk()
  return result
}

func (v *VM) push(vl Value){
  v.stack = append(v.stack, vl)
  v.stackTop.number ++
}

func (v *VM) pop() *Value {
  last := v.stack[len(v.stack)-1]
  v.stack = v.stack[:len(v.stack)-1]
  v.stackTop.number --
  return &last
}

func (v *VM) read_byte() {
  v.ip ++
}

func (v *VM) run () InterpretResult{
  for {
    if DEBUG_TRACE_EXECUTION == 1 {
      fmt.Printf("Stack Trace: ")
      for i := range len(v.stack) {
        fmt.Printf("[")
        switch v.stack[i].valType {
          case VAL_BOOL: {
            fmt.Printf("'%v'", v.stack[i].AsBoolean())
          }
          case VAL_NUMBER: {
            fmt.Printf("'%v'", v.stack[i].AsNumber())
          }
          case VAL_NIL: {
            fmt.Printf("'%v'", nil)
          }
        }
        fmt.Printf("]")
      }
      fmt.Printf("\n")
      DisassembleInstruction(v.chunk, v.ip)
    }
    ins := v.chunk.code[v.ip]
    v.read_byte()
    switch ins{
      case OpReturn: {
        c := v.pop()
        printValue(*c)
        fmt.Printf("\n")
        return INTERPRET_OK
      }
      case OpConstant: {
        c := v.read_constant()
        v.push(c)
        fmt.Printf("Constant: ")
        printValue(c)
        fmt.Printf("\n")
        break
      }
      case OpNegate: {
        if !v.peek(0).IsNumber(){
          // runtimeError("Operand must be a number")
          return INTERPRET_RUNTIME_ERROR
        }
        c := v.pop().AsNumber()
        d := NumberVal(NumberVal(-1).number * c)
        v.push(d)
        break
      }
      case OpAdd : {
        err := v.binary_op("+"); if err != nil {return INTERPRET_RUNTIME_ERROR }
        break
      }
      case OpMul: {
        err := v.binary_op("*"); if err != nil {return INTERPRET_RUNTIME_ERROR }
        break
      }
      case OpSub : {
        err := v.binary_op("-"); if err != nil {return INTERPRET_RUNTIME_ERROR }
        break
      }
      case OpDiv : {
        err := v.binary_op("/"); if err != nil { return INTERPRET_RUNTIME_ERROR }
        break
      }
      case OpFalse : {
        v.push(BoolVal(false))
      }
      case OpTrue: {
        v.push(BoolVal(true))
      }
      case OpNil: {
        v.push(NilVal())
      }
      default:
    }
  }
}

func (v *VM) peek(distance int) *Value {
  return &v.stack[len(v.stack)-1-distance]
}

func (v *VM) binary_op(op string) (error){
  if !v.peek(0).IsNumber() || !v.peek(1).IsNumber(){
    return errors.New("Runtime Error")
  }
  a := v.pop().AsNumber()
  b := v.pop().AsNumber()
  switch op{
    case "-": {
      v.push(NumberVal(b - a))
    }
    case "+": {
      v.push(NumberVal(b + a))
    }
    case "*": {
      v.push(NumberVal(b * a))
    }
    case "/": {
      v.push(NumberVal(b / a))
    }
    default: {
      return errors.New("Runtime Error")
    }
  }
  return nil
}

func (v *VM) read_constant() Value {
  index := v.chunk.code[v.ip]
  c := v.chunk.constants.values[index]
  v.read_byte()
  return c
}

