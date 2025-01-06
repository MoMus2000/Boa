package main

import "fmt"

type Interpreter struct {
  statements []Statement
}

func NewInterpreter() *Interpreter {
  return &Interpreter{
    statements: make([]Statement, 0),
  }
}

func (i *Interpreter) interpret (statements []Statement) {
  for _, statement := range statements {
    statement.Accept(i)
  }
}

func (i *Interpreter) execute_statement(statement Statement) {
   statement.Accept(i)
}

func (i *Interpreter) visit_expression_statement(visitor *ExpressionStatement){
  val := i.evaluate(visitor.expr)
  fmt.Println(val)
}

func (i *Interpreter) evaluate(expr Expression) interface{} {
  return expr.Accept(i)
}

func (i *Interpreter) visit_literal_expression(visitor *LiteralExpression) interface{} {
  return visitor.value
}

func (i *Interpreter) visit_unary_expression(visitor *UnaryExpression) interface {}{
  right := i.evaluate(visitor.right)
  switch v := right.(type){
    case int32, int64, float32, float64 : {
      if visitor.op.Type == MINUS{
        return v.(float64) * -1.0
      }
      if visitor.op.Type == BANG{
        return !is_truthy(v)
      }
    }
    case bool : {
      if visitor.op.Type == BANG{
        return !is_truthy(v)
      }
    }
    default:
      fmt.Println(v)
  }
  return nil
}

func is_truthy(value interface{}) bool{
  switch value.(type){
    case nil :
      return false
    case float64 :
      return true
    case string :
      return true
    case bool :
      return value.(bool)
    default:
      return true
  }
}

