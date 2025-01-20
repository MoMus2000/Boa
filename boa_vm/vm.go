package main

import (
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
}

func NewVM() VM {
  return VM{}
}

func (v *VM) resetStack() {
  v.stackTop = Value(0)
  v.stack    = make([]Value, 0)
}

func (v *VM) FreeVM() {
}

func (v *VM) interpret(c *Chunk) InterpretResult{
  v.chunk = c
  v.ip    = 0
  return v.run()
}

func (v *VM) push(vl Value){
  v.stack = append(v.stack, vl)
  v.stackTop ++
}

func (v *VM) pop() Value {
  last := v.stack[len(v.stack)-1]
  v.stack = v.stack[:len(v.stack)-1]
  v.stackTop --
  return last
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
        fmt.Printf("%v",v.stack[i])
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
        fmt.Printf("%v", c)
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
        c := v.pop()
        c = Value(-1) * c
        v.push(c)
        break
      }
      case OpAdd : {
        v.binary_op("+")
        break
      }
      case OpMul: {
        v.binary_op("*")
        break
      }
      case OpSub : {
        v.binary_op("-")
        break
      }
      case OpDiv : {
        v.binary_op("/")
        break
      }
      default:
    }
  }
}

func (v *VM) binary_op(op string) {
  a := v.pop()
  b := v.pop()
  switch op{
    case "-": {
      v.push(b - a)
    }
    case "+": {
      v.push(b + a)
    }
    case "*": {
      v.push(b * a)
    }
    case "/": {
      v.push(b / a)
    }
    default: {
      panic("Undefined Op")
    }
  }
}

func (v *VM) read_constant() Value {
  index := v.chunk.code[v.ip]
  c := v.chunk.constants.values[index]
  v.read_byte()
  return c
}

