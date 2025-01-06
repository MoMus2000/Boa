package main

import "fmt"

type TokenType int

const (

  LEFT_PAREN TokenType = iota
  RIGHT_PAREN 

  LEFT_BRACE
  RIGHT_BRACE

  LEFT_ANGLE_BRACKET   
  RIGHT_ANGLE_BRACKET  
                       
  PIPE                 
  COLON                
                       
  COMMA                
  DOT                  
  MINUS                
  PLUS                 
  SEMICOLON            
  SLASH                
  STAR                 
                       
  BANG                 
  BANG_EQUAL           
  EQUAL                
  EQUAL_EQUAL          
  GREATER              
  GREATER_EQUAL        
  LESS                 
  LESS_EQUAL           
                       
  IDENTIFIER           
  STRING               
  NUMBER               
                       
  AND                  
  CLASS                
  ELSE                 
  FALSE                
  FUN                  
  FOR                  
  IF                   
  NIL                  
  OR                   
  DEBUG
  RETURN               
  SUPER                
  THIS                 
  TRUE                 
  VAR                  
  WHILE                
  IMPORT               
                       
  EOF                  

)

var tokenNames = [...]string{
	"LEFT_PAREN", "RIGHT_PAREN", "LEFT_BRACE", "RIGHT_BRACE", "LEFT_ANGLE_BRACKET", 
	"RIGHT_ANGLE_BRACKET", "PIPE", "COLON", "COMMA", "DOT", "MINUS", "PLUS", 
	"SEMICOLON", "SLASH", "STAR", "BANG", "BANG_EQUAL", "EQUAL", "EQUAL_EQUAL", 
	"GREATER", "GREATER_EQUAL", "LESS", "LESS_EQUAL", "IDENTIFIER", "STRING", 
	"NUMBER", "AND", "CLASS", "ELSE", "FALSE", "FUN", "FOR", "IF", "NIL", 
	"OR", "DEBUG", "RETURN", "SUPER", "THIS", "TRUE", "VAR", "WHILE", 
	"IMPORT", "EOF",
}

// String method to print enum name
func (t TokenType) String() string {
	if t < LEFT_PAREN || t > EOF {
		return fmt.Sprintf("Unknown TokenType(%d)", t)
	}
	return tokenNames[t]
}
