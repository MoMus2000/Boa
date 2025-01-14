package main

import (
  "unicode"
)

type Lexer struct {
  source    string
  tokens    []Token
  start     uint
  current   uint
  line      uint
  ident_map map[string]TokenType
}

func NewLexer(source_code []byte) *Lexer {
  ident_map := map[string]TokenType{
    "and"    : AND,
    "class"  : CLASS,
    "else"   : ELSE,
    "false"  : FALSE,
    "for"    : FOR,
    "fun"    : FUN,
    "if"     : IF,
    "nil"    : NIL,
    "or"     : OR,
    "dbg"    : DEBUG,
    "return" : RETURN,
    "super"  : SUPER,
    "this"   : THIS,
    "true"   : TRUE,
    "var"    : VAR,
    "while"  : WHILE,
    "import" : IMPORT,
    "in"     : IN,
    "range"  : RANGE,
  }
	return &Lexer{
    source:    string(source_code),
    tokens:    make([]Token, 0),
    start :    0,
    current :  0,
    line :     0,
    ident_map: ident_map,
  }
}

func (lexer *Lexer) ScanTokens() []Token{
  for lexer.isNotAtEnd(){
    lexer.start = lexer.current
    lexer.scanToken()
  }
  lexer.tokens = append(lexer.tokens, Token{
    Type: EOF,
    Line: lexer.line,
    Lexeme: "",
    Literal: nil,
  })
  return lexer.tokens
}

func (lexer *Lexer) isNotAtEnd() bool{
  if int(lexer.current) >= len(lexer.source){
    return false
  }
  return true
}

func (lexer *Lexer) add_token(ttype TokenType) {
  lexer.tokens = append(lexer.tokens, Token{
    Type:   ttype,
    Line:   lexer.line,
    Lexeme: lexer.source[lexer.start : lexer.current],
    Literal: nil,
  })
}

func (lexer *Lexer) add_token_literal(ttype TokenType, literal string) {
  lexer.tokens = append(lexer.tokens, Token{
    Type:   ttype,
    Line:   lexer.line,
    Lexeme: lexer.source[lexer.start : lexer.current],
    Literal: literal,
  })
}

func (lexer *Lexer) scanToken() {
  c := lexer.advance()
  switch c {
  case '(' :
    lexer.add_token(LEFT_PAREN)
  case ')' : 
    lexer.add_token(RIGHT_PAREN)
  case '{' :
    lexer.add_token(LEFT_BRACE)
  case '}' : 
    lexer.add_token(RIGHT_BRACE)
  case ',' :
    lexer.add_token(COMMA)
  case '.' :
    lexer.add_token(DOT)
  case '-' : 
    lexer.add_token(MINUS)
  case '+' :
    lexer.add_token(PLUS)
  case ';' :
    lexer.add_token(SEMICOLON)
  case '*' :
    lexer.add_token(STAR)
  case '[' : 
    lexer.add_token(LEFT_ANGLE_BRACKET)
  case ']' :
    lexer.add_token(RIGHT_ANGLE_BRACKET)
  case ':' :
    lexer.add_token(COLON)
  case '!' :
    if lexer.match('=') {
      lexer.add_token(BANG_EQUAL)
    }else{
      lexer.add_token(BANG)
    }
  case '=' :
    if lexer.match('=') {
      lexer.add_token(EQUAL_EQUAL)
    }else{
      lexer.add_token(EQUAL)
    }
  case '<' :
    if lexer.match('=') {
      lexer.add_token(LESS_EQUAL)
    }else{
      lexer.add_token(LESS)
    }
  case '>' :
    if lexer.match('=') {
      lexer.add_token(GREATER_EQUAL)
    }else{
      lexer.add_token(GREATER)
    }
  case ' ':
    // do nothing
  case '/':
    if lexer.match('/'){
      for lexer.peek() != '\x00' && lexer.isNotAtEnd(){
        lexer.advance()
      }
    } else {
      lexer.add_token(SLASH)
    }
  case '\n':
    // do nothing
    lexer.line += 1
  case '\t':
    // do nothing
  case '\r':
    // do nothing
  case '"':{
    lexer.lex_string()
    lexer.add_token_literal(STRING, "string")
  }
  default:
    if unicode.IsDigit(rune(c)){
      lexer.lex_number()
      lexer.add_token_literal(NUMBER, "number")
    }else if unicode.IsLetter(rune(c)){
      lexer.lex_ident()
    }
  }
}

func (lexer *Lexer) lex_ident() {
  for unicode.IsDigit(rune(lexer.peek())) || unicode.IsLetter(rune(lexer.peek())) || 
    lexer.peek() == '_' {
    lexer.advance()
  }

  ttype, found := lexer.ident_map[lexer.source[lexer.start : lexer.current]]
  if !found {
    lexer.add_token(IDENTIFIER)
  } else{
    lexer.add_token(ttype)
  }
}

func (lexer *Lexer) lex_number(){
  for unicode.IsDigit(rune(lexer.peek())){
    lexer.advance()
  }
  if lexer.peek() == '.' && unicode.IsDigit(rune(lexer.peek_next())){
    lexer.advance()
    for unicode.IsDigit(rune(lexer.peek())){
      lexer.advance()
    }
  }
}

func (lexer *Lexer) peek_next() byte{
  if int(lexer.current + 1) >= len(lexer.source){
    return '\x00'
  }
  return lexer.source[lexer.current+1]
}

func (lexer *Lexer) lex_string() {
  for lexer.peek() != '"' && lexer.isNotAtEnd() {
    if lexer.peek() == '\n'{
      lexer.line += 1
    }
    lexer.advance()
  }
  lexer.advance()
}

func (lexer *Lexer) peek() byte {
  if !lexer.isNotAtEnd(){
    return '\x00'
  }
  return lexer.source[lexer.current]
}

func (lexer *Lexer) match(c byte) bool{
  if !lexer.isNotAtEnd() {
    return false
  }
  if lexer.source[lexer.current] != c{
    return false
  }
  lexer.current += 1
  return true
}

func (lexer *Lexer) advance() byte {
  c := lexer.source[lexer.current]
  lexer.current += 1
  return c
}

