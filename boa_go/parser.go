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
  for !p.is_at_end(){
    // statements = append(statements, )
  }
  return statements
}

//  program     -> declaration* eof
//  
//  declaration -> varDecl
//              | statement
//  
//  varDecl     -> "var" IDENTIFIER ( "=" expression )? ";"
//  
//  statement   -> exprStmt
//              | forStmt
//              | ifStmt
//              | printStmt
//              | whileStmt
//              | blockStmt
//  
//  exprStmt    -> expression ";"
//  
//  forStmt     -> "for" "(" (varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement
//  
//  ifStmt      -> "if" "(" expression ")" statement ( "else" statement )?
//  
//  blockStmt   -> "{" declaration* "}"
//  
//  whileStmt   -> "while" "(" expression ")" statement
//  
//  printStmt   -> "print" expression ";"
//  
//  expression  -> assignment
//  
//  assignment  -> identifier ( "=" assignment )?
//              | logic_or
//  
//  logic_or    -> logic_and ( "or " logic_and )*
//  
//  logic_and   -> equality ( "and" equality )*
//  
//  equality    -> comparison ( ( "!=" | "==" ) comparison )*
//  
//  comparison  -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
//  
//  term        -> factor ( ( "-" | "+" ) factor )*
//  
//  factor      -> unary ( ( "/" | "*" ) unary )*
//  
//  unary       -> ( "!" | "-" ) unary
//              | primary
//  
//  primary     -> NUMBER | STRING
//              | "false" | "true" | "nil"
//              | "(" expression ")"
//              | IDENTIFIER


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

