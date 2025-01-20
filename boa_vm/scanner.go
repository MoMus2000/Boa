package main

type Scanner struct {
  start   int
  current int
  line    int
  src     []rune
}

type Token struct {
  tokenType   TokenType
  runes       []rune
  length      int
  line        int
}

func NewScanner(source []byte) Scanner {
  return Scanner{
    src     : []rune(string(source)),
    start   : 0,
    current : 0,
    line    : 1,
  }
}

func (s *Scanner) advance() rune {
  s.current ++
  return s.src[s.current-1]
}

func (s *Scanner) scanToken() Token {
  s.start = s.current

  c := s.advance()

  switch c {
    case '(': return s.makeToken(LEFT_PAREN)
    case ')': return s.makeToken(RIGHT_PAREN)
    case '{': return s.makeToken(LEFT_BRACE)
    case '}': return s.makeToken(RIGHT_BRACE)
    case ';': return s.makeToken(SEMICOLON)
    case ',': return s.makeToken(COMMA)
    case '.': return s.makeToken(DOT)
    case '-': return s.makeToken(MINUS)
    case '+': return s.makeToken(PLUS)
    case '/': return s.makeToken(SLASH)
    case '*': return s.makeToken(STAR)
  }

  if s.isAtEnd() {
    return s.makeToken(EOF)
  }

  return s.tokenError("Unexpected Character")
}

func (s *Scanner) isAtEnd() bool {
  return s.current >= len(s.src)
}

func (s *Scanner) tokenError(msg string) Token {
  return Token{
    runes : []rune(msg),
    length : len(msg),
    line   : s.line,
    tokenType: TOKEN_ERROR,
  }
}

func (s *Scanner) makeToken(tokenType TokenType) Token {
  return Token{
    runes    : s.src[s.start:s.current],
    tokenType: tokenType,
    length   : s.current - s.start,
    line     : s.line,
  }
}

