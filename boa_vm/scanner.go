package main

import (
  _"fmt"
)

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

func (s *Scanner) skipWhiteSpace() {
  for {
    c := s.peek()
    switch c {
      case ' ' : {
        s.advance()
        break
      }
      case '\r': {
        s.advance()
        break
      }
      case '\t': {
        s.advance()
        break
      }
      case '\n': {
        s.line ++
        s.advance()
        break
      }
      case '/': {
        if s.peekNext() == '/' {
          for s.peek() != '\n' && !s.isAtEnd() { s.advance() }
        } else {
          return
        }
      }
      default:
        return
    }
  }
}

func (s *Scanner) peekNext() rune {
  if s.isAtEnd() { return '\u0000' }
  return s.src[s.current+1]
}

func (s *Scanner) peek() (r rune) {
  if s.isAtEnd() { return }
  return s.src[s.current]
}

func (s *Scanner) scanToken() Token {
  s.skipWhiteSpace()
  s.start = s.current

  if s.isAtEnd() {
    return s.makeToken(EOF)
  }

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
    case '!': {
      if s.match('=') {
        return s.makeToken(BANG_EQUAL)
      } else {
        return s.makeToken(BANG)
      }
    }
    case '=': {
      if s.match('=') {
        return s.makeToken(EQUAL_EQUAL)
      } else {
        return s.makeToken(EQUAL)
      }
    }
    case '<': {
      if s.match('=') {
        return s.makeToken(LESS_EQUAL)
      } else {
        return s.makeToken(LESS)
      }
    }
    case '>': {
      if s.match('=') {
        return s.makeToken(GREATER_EQUAL)
      } else {
        return s.makeToken(GREATER)
      }
    }
    case '"': {
      return s.makeStringToken()
    }
  }

  return s.tokenError("Unexpected Character")
}

func (s *Scanner) makeStringToken() Token {
  for s.peek() != '"' && !s.isAtEnd() {
    if s.peek() == '\n' { s.line ++ }
    s.advance()
  }
  if s.isAtEnd() { return s.tokenError("Unterminated String") }
  s.advance()
  return s.makeToken(TOKEN_STRING)
}

func (s *Scanner) match(lexeme rune) bool {
  if s.isAtEnd(){return false}
  if s.src[s.current] != lexeme {
    return false
  }
  s.current ++
  return true
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

