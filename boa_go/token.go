package main

import "fmt"

type Token struct {
  Line    uint
  Type    TokenType
  Literal interface{}
  Lexeme  interface{}
}

func (t Token) String() string{
  return fmt.Sprintf("Token{Type: '%v', Literal: '%v' Lexeme: '%v'}", t.Type, t.Literal, t.Lexeme)
}

