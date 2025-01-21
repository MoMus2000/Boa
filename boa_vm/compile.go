package main

import (
	"fmt"
	"strconv"
)


type Parser struct {
  current   Token
  previous  Token
  hadError  bool
  panicMode bool
}

type Compiler struct {
  scanner *Scanner
  parser  *Parser
}

func NewCompiler() Compiler {
  return Compiler{}
}

var compilingChunk *Chunk

func (c *Compiler) compile(source []byte, chunk *Chunk) bool{
  parser    := Parser{}
  scanner   := NewScanner(source)
  c.scanner = &scanner
  c.parser  = &parser
  c.parser.hadError  = false
  c.parser.panicMode = false
  compilingChunk = chunk
  c.advance()
  c.expression()
  // c.consume(EOF, "Expected EOF")
  c.endCompiler()
  return !c.parser.hadError
}

func (c *Compiler) endCompiler() {
  c.emitReturn()
}

func (c *Compiler) number() {
  num, err := strconv.ParseFloat(string(c.parser.current.runes), 32)
  if err != nil  {
    err := err.Error()
    c.errorAtCurrent(err)
  }
  c.emitBytes(OpConstant, c.makeConstant(Value(num)))
}

func (c *Compiler) makeConstant(constant Value) Opcode {
  index := currentChunk().AddConstant(constant)
  return Opcode(index)
}

func(c *Compiler) emitReturn() {
  c.emitByteCode(OpReturn)
}

func (c *Compiler) emitBytes(a Opcode, b Opcode) {
  c.emitByteCode(a)
  c.emitByteCode(b)
}

func currentChunk() *Chunk {
  return compilingChunk
}

func (c *Compiler) expression() {
  c.number()
}

func (c *Compiler) advance() {
  c.parser.previous = c.parser.current 
  for { 
    token := c.scanner.scanToken() 
    c.parser.current = token
    if c.parser.current.tokenType != ERROR {
      break
    }
    c.errorAtCurrent(string(c.parser.current.runes))
  }
}

func (c *Compiler) consume(tokenType TokenType, message string) {
  if c.parser.current.tokenType == tokenType {
    c.advance()
  }
  c.errorAtCurrent(message)
}

func (c *Compiler) emitByteCode(code Opcode) {
  currentChunk().WriteChunk(code, c.parser.previous.line)
}

func (c *Compiler) errorAtCurrent(message string){
  c.errorAt(&c.parser.current, message)
}

func (c *Compiler) error(message string){
  c.errorAt(&c.parser.previous, message)
}

func (c *Compiler) errorAt(token *Token, message string){
  if c.parser.panicMode { return }
  c.parser.panicMode = true
  fmt.Printf("[line %d] Error", token.line)
  if token.tokenType == EOF {
    fmt.Printf(" at end")
  } else if token.tokenType == ERROR {

  } else {
    fmt.Printf(" at '%.*s'", token.length, string(token.runes))
  }
  fmt.Printf(": %s\n", message);
  c.parser.hadError = true
}


