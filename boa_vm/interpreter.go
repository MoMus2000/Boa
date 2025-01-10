package main

import (
  "fmt"
  "reflect"
  "time"
)
  

type Interpreter struct {
  env        *Env
  statements []Statement
}

type ReturnValue struct {
  value interface{}
}

type Clock struct {
  arity int32
}

func (clock *Clock) call(interpreter *Interpreter, args []interface{}){
  t := time.Now().Unix()
  fmt.Println(t)
}

func NewInterpreter() *Interpreter {
  env := NewEnv(nil)
  env.define("clock", &Clock{0})
  return &Interpreter{
    env       : env,
    statements: make([]Statement, 0),
  }
}

func (i *Interpreter) interpret (statements []Statement) error {
  for _, statement := range statements {
    err := statement.Accept(i)
    if err != nil {
      return err
    }
  }
  return nil
}

func (i *Interpreter) execute_statement(statement Statement) {
   statement.Accept(i)
}

func (i *Interpreter) visit_return_statement(statement *ReturnStatement) error {
  expr := i.evaluate(statement.val)
  panic(ReturnValue{expr})
}

func (i *Interpreter) visit_func_call_expression(statement *FuncCallExpression) (returnval interface{}){
  fun := i.env.get(statement.ident.Lexeme.(string))
  defer func() {
      if r := recover(); r != nil {
        returnval = r.(ReturnValue).value
        return 
      }
  }()
  i_copy := &Interpreter{
    env: i.env,
  }
  f_args := make([]interface{}, 0)
  for _, arg := range statement.args{
    expr := i.evaluate(arg)
    f_args = append(f_args, expr)
  }
  switch f := fun.(type) {
    case *CallableFunc:
        f.call(i_copy, f_args)
    case *Clock:
        f.call(i_copy, f_args)
    default:
        fmt.Println("Unsupported type")
  }
  return nil
}

func (i *Interpreter) visit_func_statement(statement *FunctionStatement){
  cf := &CallableFunc{
    declaration: statement,
  }
  i.env.define(statement.ident.Lexeme.(string), cf)
}

func (i *Interpreter) visit_for_statement(visitor *ForStatement){
  i.execute_statement(visitor.start)
  for i.evaluate(visitor.predicate) == true{
    i.execute_statement(visitor.inner_statements)
    i.evaluate(visitor.incre)
  }
}

func (i *Interpreter) visit_while_statement(visitor *WhileStatement){
  predicate := is_truthy(i.evaluate(visitor.predicate))
  for predicate {
    i.execute_statement(visitor.inner_statements)
    predicate = is_truthy(i.evaluate(visitor.predicate))
  }
}

func (i *Interpreter) visit_if_statement(visitor *IfStatement){
  predicate := is_truthy(i.evaluate(visitor.predicate))
  if predicate {
    i.execute_statement(visitor.if_condition)
  } else if !predicate && visitor.else_condition != nil{
    i.execute_statement(visitor.else_condition)
  }
}

func (i *Interpreter) visit_assign_expression(visitor *AssignExpression) interface{}{
  assigned := i.env.assign(visitor.ident.Lexeme.(string), i.evaluate(visitor.value))
  return assigned
}

func (i *Interpreter) visit_block_statement(visitor *BlockStatement) {
  env := NewEnv(i.env)
  prev := i.env
  i.env = env
  for _, statement := range visitor.statements{
    i.execute_statement(statement)
  }
  i.env = prev
}

func (i *Interpreter) execute_block(statement *BlockStatement, env *Env){
  prev := i.env
  i.env = env
  for _, statement := range statement.statements{
    i.execute_statement(statement)
  }
  i.env = prev
}

func (i *Interpreter) visit_debug_statement(visitor *DebugStatement){
  evaluated_expr := i.evaluate(visitor.expr)
  fmt.Println(evaluated_expr)
}

func (i *Interpreter) visit_var_statement(visitor *VarStatement){
  i.env.define(visitor.ident.Lexeme.(string), i.evaluate(visitor.value))
}

func (i *Interpreter) visit_logical_expression(visitor *LogicalExpression) interface{}{
  left := i.evaluate(visitor.left)
  right := i.evaluate(visitor.right)
  if visitor.op.Type == OR{
    return is_truthy(left) || is_truthy(right)
  }
  if visitor.op.Type == AND{
    return is_truthy(left) && is_truthy(right)
  }
  panic("Unreachable Code")
}

func (i *Interpreter) visit_expression_statement(visitor *ExpressionStatement){
  i.evaluate(visitor.expr)
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

func (i *Interpreter) visit_var_expression(visitor *VarExpression) interface{} {
  return i.env.get(visitor.ident.Lexeme.(string))
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

