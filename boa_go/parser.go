package main

import (
  "fmt"
  "strconv"
)

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
    statement := p.statement()
    statements = append(statements, statement)
  }
  return statements
}

func (p *Parser) statement() Statement {
  if p.match(DEBUG){
     return p.debug_statement()
  }
  return p.expression_statement()
}

func (p *Parser) debug_statement() Statement{
  expr := p.expression()
  return &DebugStatement{
    expr: expr,
  }
}

func (p *Parser) expression() Expression{
  return p.logical() // Point where the expression parsing begins
}

func (p *Parser) expression_statement() Statement{
  expression := p.expression()
  return &ExpressionStatement{
    expr: expression,
  }
}

func (p *Parser) logical() Expression{
  expr := p.equality()
  for p.match(OR, AND){
    op    := p.previous()
    right := p.logical()
    return &LogicalExpression{
      op   : op,
      right: right,
      left : expr,
    }
  }
  return expr
}

func (p *Parser) equality() Expression{
  expr := p.comparision()
  for p.match(EQUAL_EQUAL, BANG_EQUAL){
    op    := p.previous()
    right := p.comparision()
    return &BinaryExpression{
      op : op,
      right : right,
      left  : expr,
    }
  }
  return expr
}

func (p *Parser) comparision() Expression {
  expr := p.term()
  for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL){
    op    := p.previous()
    right := p.term()
    return &BinaryExpression{
      op : op,
      right : right,
      left  : expr,
    }
  }
  return expr
}

func (p *Parser) term() Expression {
  expr := p.factor()
  for p.match(PLUS, MINUS){
    op    := p.previous()
    right := p.factor()
    return &BinaryExpression{
      op : op,
      right : right,
      left  : expr,
    }
  }
  return expr
}

func (p *Parser) factor() Expression {
  expr := p.unary()
  for p.match(SLASH, STAR){
    op    := p.previous()
    right := p.unary()
    return &BinaryExpression{
      op    : op,
      right : right,
      left  : expr,
    }
  }
  return expr
}

func (p *Parser) unary() Expression {
  if p.match(BANG, MINUS){
    op := p.previous()
    right := p.unary()
    return &UnaryExpression{
      op: op,
      right: right,
    }
  }
  return p.primary()
}

func (p *Parser) primary() Expression {
  token := p.advance()
  switch token.Type {
    case NUMBER : 
      floatValue, err := strconv.ParseFloat(token.Lexeme.(string), 64)
      if err != nil {
        panic("Valid number not provided")
      }
      return &LiteralExpression{
        value: floatValue,
      }
    case STRING : 
      return &LiteralExpression{
        value: token.Lexeme,
      }
    case FALSE  : 
      return &LiteralExpression{
        value: false,
      }
    case TRUE   : 
      return &LiteralExpression{
        value: true,
      }
    case NIL :
      return &LiteralExpression{
        value: nil,
      }
    case LEFT_PAREN:
      ge := &GroupingExpression{
        expr: p.expression(),
      }
      p.consume(RIGHT_PAREN, "Expect )")
      return ge
    case IDENTIFIER:
      fmt.Println("Found an Identifier")
    default:
      panic(fmt.Sprint("Unexpected primary value encountered ", token.Type))
  }
  panic("Unreachable Code")
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
    fmt.Println(message)
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

