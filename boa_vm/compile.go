package main

import (
  "fmt"
)

func compile(source []byte){
  scanner := NewScanner(source)
  line := -1
  for {
    token := scanner.scanToken()
    if token.line != line {
      fmt.Printf("%4d", token.line)
      line  = token.line
    } else{
      fmt.Printf("   | ")
    }
    fmt.Printf(" %s %.*s\n", token.tokenType, token.length, string(token.runes[:token.length]))
    if token.tokenType == EOF || token.tokenType == ERROR{
      break
    }
  }
}

