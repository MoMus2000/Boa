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

func (p *Parser) parse() ([]Statement, error) {
  var statements []Statement = make([]Statement, 0)
  for !p.is_at_end(){
    statement, err := p.declaration()
    if err != nil {
      return nil, err
    }
    statements = append(statements, statement)
  }
  return statements, nil
}

func (p *Parser) declaration() (Statement, error) {
  if p.match(VAR){
    v, err := p.var_declaration()
    if err != nil {
      return nil, err
    }
    return v, nil
  }
  if p.match(FUN){
    return p.func_declaration()
  }
  return p.statement()
}

func (p *Parser) var_declaration() (Statement, error){
  ident, err := p.consume(IDENTIFIER, "Expected an identifier")
  if err != nil {
    return nil, err
  }
  var expr Expression
  if p.match(EQUAL){ expr = p.expression() }
  p.consume(SEMICOLON, "Expected ;")
  return &VarStatement{
    ident: *ident,
    value: expr,
  }, nil
}

func (p *Parser) statement() (Statement, error) {
  if p.match(DEBUG){
     return p.debug_statement(), nil
  }
  if p.match(IF){
    return p.if_statement()
  }
  if p.match(LEFT_BRACE){
    return  p.block_statement()
  }
  if p.match(WHILE){
    return p.while_statement()
  }
  if p.match(FOR){
    return p.for_loop_statement()
  }
  if p.match(RETURN){
    return p.return_statement(), nil
  }
  return p.expression_statement(), nil
}

func (p *Parser) return_statement() Statement {
  ident := p.previous()
  expr  := p.expression()
  p.consume(SEMICOLON, "Expected ;")
  return &ReturnStatement{
    ident: *ident,
    val  : expr,
  }
}

func (p *Parser) func_declaration() (Statement, error) {
  ident, err := p.consume(IDENTIFIER, "Expected Function Name")
  if err != nil {
    return nil, err
  }
  p.consume(LEFT_PAREN, "Expected (")
  args := make([]string , 0)
  if !p.check(RIGHT_PAREN){
    for {
      arg, err := p.consume(IDENTIFIER, "Expected a function arg")
      if err != nil {
        return nil, err
      }
      args = append(args, arg.Lexeme.(string))
      if !p.match(COMMA){
        break
      }
    }
  }
  p.consume(RIGHT_PAREN, "Expected )")
  body, err := p.statement()
  return &FunctionStatement{
    ident: *ident,
    args : args,
    body : body.(*BlockStatement),
  }, nil
}

func (p *Parser) for_loop_statement() (Statement, error){
  p.consume(LEFT_PAREN, "Expected (")
  start, err      := p.declaration()
  if err != nil{
    return nil, err
  }
  predicate := p.expression()
  p.consume(SEMICOLON, "Expected ;")
  increment := p.expression()
  p.consume(RIGHT_PAREN, "Expected )")
  block_statement , err := p.statement()
  return &ForStatement{
    start,
    predicate,
    increment,
    block_statement.(*BlockStatement),
  }, err

}

func (p *Parser) while_statement() (Statement, error) {
  p.consume(LEFT_PAREN, "Expected (")
  predicate := p.expression()
  p.consume(RIGHT_PAREN, "Expected )")
  block_statement, err := p.statement()
  return &WhileStatement{
    predicate: predicate,
    inner_statements: block_statement.(*BlockStatement),
  }, err
}

func (p *Parser) if_statement() (Statement, error){
  p.consume(LEFT_PAREN, "Expected (")
  predicate := p.expression()
  p.consume(RIGHT_PAREN, "Expected )")
  if_block, err := p.statement()
  if err != nil {
    return nil, err
  }
  if p.match(ELSE){
    else_condition, err := p.statement()
    if err != nil {
      return nil, err
    }
    return &IfStatement{
      predicate: predicate,
      if_condition: if_block.(*BlockStatement),
      else_condition: else_condition.(*BlockStatement),
    }, err
  }
  return &IfStatement{
    predicate: predicate,
    if_condition: if_block.(*BlockStatement),
    else_condition: nil,
  }, err
}

func (p *Parser) block_statement() (Statement, error) {
  statements := make([]Statement, 0)
  for !p.match(RIGHT_BRACE) && !p.is_at_end(){
    statement, err := p.declaration()
    if err != nil {
      return nil, err
    }
    statements = append(statements, statement)
  }
  return &BlockStatement{statements: statements}, nil
}

func (p *Parser) debug_statement() Statement{
  expr := p.expression()
  p.consume(SEMICOLON, "Expected ;")
  return &DebugStatement{
    expr: expr,
  }
}

func (p *Parser) assign() Expression {
  expr := p.logical()
  if p.match(EQUAL){
    value := p.assign()
    token := expr.(*VarExpression).ident
    return &AssignExpression{
      ident: token,
      value: value,
    }
  }
  return expr
}

func (p *Parser) expression() Expression{
    expr:= p.assign() // Point where the expression parsing begins
    return expr
}

func (p *Parser) expression_statement() Statement{
  expression := p.expression()
  p.consume(SEMICOLON, "Expected ;")
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
      op   : *op,
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
      op : *op,
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
      op : *op,
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
      op : *op,
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
      op    : *op,
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
      op: *op,
      right: right,
    }
  }
  return p.call()
}

func (p *Parser) call() Expression{
  expr := p.primary()
  for p.match(LEFT_PAREN){
    expr = p.finish_call(expr)
  }
  return expr
}

func (p *Parser) finish_call(expr Expression) Expression{
  ident := expr.(*VarExpression).ident
  args := make([]Expression, 0)
  if !p.check(RIGHT_PAREN){
    for {
      args = append(args, p.expression())
      if !p.match(COMMA){
        break
      }
    }
  }
  p.consume(RIGHT_PAREN, "Expected ( after args")
  return &FuncCallExpression{
    ident: ident,
    args: args,
  }
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
      return &VarExpression{
        ident: *p.previous(),
      }
    default:
      panic(fmt.Sprint("Unexpected primary value encountered ", token.Type))
  }
}

// func (p *Parser) call(ident Token) Expression{
//   args := make([]Expression, 0)
//   if !p.check(RIGHT_PAREN){
//     for{
//       e := p.expression()
//       args = append(args, e)
//       if !p.match(COMMA){
//         break
//       }
//     }
//   }
//   p.consume(RIGHT_PAREN, "Expected )")
//   return &FuncCallExpression{
//     ident: ident,
//     args: args,
//   }
// }

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


func (p *Parser) consume(ttype TokenType, message string) (*Token, error) {
  consumed := p.match(ttype)
  if !consumed {
    return nil, error(fmt.Errorf(message))
  }
  return p.previous(), nil
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
  return *p.previous()
}

func (p *Parser) previous() *Token{
  return &p.tokens[p.current-1]
}

func (p *Parser) peek() Token {
  return p.tokens[p.current]
}

func (p *Parser) is_at_end() bool {
  return p.peek().Type == EOF
}

