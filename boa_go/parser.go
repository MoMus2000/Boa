package main

import "fmt"

type Parser struct {
  current uint
  tokens  []Token
}

func NewParser(source_code []byte) *Parser {
	lexer  := NewLexer(source_code)
  tokens := lexer.ScanTokens()
  return &Parser{
    current: 0,
    tokens : tokens,
  }
}

func (p *Parser) parse() []Statement {
  var statements []Statement = make([]Statement, 0)
  le := LiteralExpression{}
  le.Accept(p)
  return statements
}

func (p *Parser) visit_binary_expression(bi Expression){

}

func (p *Parser) visit_literal_expression(le Expression){
  fmt.Println("Inside the Literal Expression")
}

func (p *Parser) match(ttype ...TokenType) bool {

  for _, token := range ttype {
    if p.check(token){
      p.advance()
      return true
    }
  }

  return false
}

func (p *Parser) consume(ttype TokenType, message string) TokenType{
  consumed := p.match(ttype)
  if !consumed {
    panic(fmt.Errorf(p.peek().String(), message))
  }
  return ttype
}

func (p *Parser) check(ttype TokenType) bool{
  if p.is_at_end(){
    return false
  }
  return p.peek().Type == ttype
}

func (p *Parser) advance() Token{
  if !p.is_at_end() {
    p.current += 1
  }
  return p.previous()
}

func (p *Parser) previous() Token{
  return p.tokens[p.current-1]
}

func (p *Parser) peek() Token {
  return p.tokens[p.current]
}

func (p *Parser) is_at_end() bool {
  return p.peek().Type == EOF
}

