package main

import (
  "fmt"
  "reflect"
)

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

func (i *Interpreter) visit_binary_expression(visitor *BinaryExpression) interface{}{
  right := i.evaluate(visitor.right)
  left  := i.evaluate(visitor.left)
  rightType := reflect.TypeOf(right)
  leftType  := reflect.TypeOf(left)
  switch visitor.op.Type {
    case EQUAL_EQUAL: {
      return right == left
    }
    case BANG_EQUAL: {
      return right != left
    }
    case GREATER: {
      if rightType.Kind() == reflect.Float64 && leftType.Kind() == reflect.Float64 {
        return left.(float64) > right.(float64)
      }else if rightType.Kind() == reflect.String && leftType.Kind() == reflect.String {
        return left.(string) > right.(string)
      }
    }
    case GREATER_EQUAL: {
      if rightType.Kind() == reflect.Float64 && leftType.Kind() == reflect.Float64 {
        return left.(float64) >= right.(float64)
      }else if rightType.Kind() == reflect.String && leftType.Kind() == reflect.String {
        return left.(string) >= right.(string)
      }
    }
    case LESS: {
      if rightType.Kind() == reflect.Float64 && leftType.Kind() == reflect.Float64 {
        return left.(float64) < right.(float64)
      }else if rightType.Kind() == reflect.String && leftType.Kind() == reflect.String {
        return left.(string) < right.(string)
      }
    }
    case LESS_EQUAL: {
      if rightType.Kind() == reflect.Float64 && leftType.Kind() == reflect.Float64 {
        return left.(float64) <= right.(float64)
      }else if rightType.Kind() == reflect.String && leftType.Kind() == reflect.String {
        return left.(string) <= right.(string)
      }
    }
    case PLUS: {
      if rightType.Kind() == reflect.Float64 && leftType.Kind() == reflect.Float64 {
        return right.(float64) + left.(float64)
      }else if rightType.Kind() == reflect.String && leftType.Kind() == reflect.String {
        return right.(string) + left.(string)
      }
    }
    case MINUS: {
      if rightType.Kind() == reflect.Float64 && leftType.Kind() == reflect.Float64 {
        return left.(float64) - right.(float64)
      }
    }
    case STAR: {
      if rightType.Kind() == reflect.Float64 && leftType.Kind() == reflect.Float64 {
        return right.(float64) * left.(float64)
      }
    }
    case SLASH: {
      if rightType.Kind() == reflect.Float64 && leftType.Kind() == reflect.Float64 {
        return left.(float64) / right.(float64)
      }
    }
    default:
      fmt.Println(visitor.op.Type)
      panic("Undefined operator on binary expression")
  }
  return nil
}

func (i *Interpreter) visit_grouping_expression(visitor *GroupingExpression) interface{}{
  return i.evaluate(visitor.expr)
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

