package main

type Expression interface{
  // Has to be implemented by each of the expression to route to the correct "visitor"
  // For example LiteralExpression Object should route to visit_literal_expression
  // visit_literal_expression has the expression statement given to it as the arg
  // Going back to the implementation which exists inside the Callee which is the
  // interpreter
  Accept(visitor ExpressionVisitor) interface{}
}

// Has to be implemented by the interpreter
type ExpressionVisitor interface {
  visit_literal_expression  (l *LiteralExpression ) interface{}
  visit_unary_expression    (e *UnaryExpression   ) interface{}
  visit_binary_expression   (e *BinaryExpression  ) interface{}
  visit_grouping_expression (e *GroupingExpression) interface{}
  visit_logical_expression  (e *LogicalExpression)  interface{}
  visit_var_expression      (e *VarExpression)      interface{}
  visit_assign_expression   (e *AssignExpression)   interface{}
  visit_func_call_expression(e *FuncCallExpression)   interface{}
}

type UnaryExpression struct {
  op    Token
  right Expression
}

func (u *UnaryExpression) Accept(visitor ExpressionVisitor) interface{}{
  return visitor.visit_unary_expression(u)
}

type LiteralExpression struct {
  value interface{}
}

func (l *LiteralExpression) Accept(visitor ExpressionVisitor) interface{} {
  return visitor.visit_literal_expression(l)
}

type BinaryExpression struct {
  left  Expression
  right Expression
  op    Token
}

func (be *BinaryExpression) Accept(visitor ExpressionVisitor) interface{} {
  return visitor.visit_binary_expression(be)
}

type GroupingExpression struct {
  expr Expression
}

func (ge *GroupingExpression) Accept(visitor ExpressionVisitor) interface{} {
  return visitor.visit_grouping_expression(ge)
}

type LogicalExpression struct {
  op    Token
  left  Expression
  right Expression
}

func (le *LogicalExpression) Accept(visitor ExpressionVisitor) interface{} {
  return visitor.visit_logical_expression(le)
}

type VarExpression struct {
  ident Token
}

func (le *VarExpression) Accept(visitor ExpressionVisitor) interface{} {
  return visitor.visit_var_expression(le)
}

type AssignExpression struct {
  ident Token
  value Expression
}

func (as *AssignExpression) Accept(visitor ExpressionVisitor) interface{}{
  return visitor.visit_assign_expression(as)
}

type FuncCallExpression struct {
  ident Token
  args  []Expression
}

func (fs *FuncCallExpression) Accept(visitor ExpressionVisitor) interface{}{
  return visitor.visit_func_call_expression(fs)
}

