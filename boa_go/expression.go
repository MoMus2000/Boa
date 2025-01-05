package main

type Expression interface{
  // Has to be implemented by each of the expression to route to the correct "visitor"
  // For example LiteralExpression Object should route to visit_literal_expression
  // visit_literal_expression has the expression statement given to it as the arg
  // Going back to the implementation which exists inside the Callee which is the
  // interpreter
  Accept(visitor ExpressionVisitor)
}

// Has to be implemented by the interpreter
type ExpressionVisitor interface {
  visit_literal_expression(e Expression)
  visit_binary_expression(e Expression)
}

type LiteralExpression struct {
  value interface{}
}

type BinaryExpression struct {
  value interface{}
}

func (l *LiteralExpression) Accept(visitor ExpressionVisitor){
  visitor.visit_literal_expression(l)
}

func (l *BinaryExpression) Accept(visitor ExpressionVisitor){
  visitor.visit_binary_expression(l)
}

