package main

type Statement interface{
  Accept(visitor StatementVisitor)
}

type StatementVisitor interface{
  visit_expression_statement(es *ExpressionStatement) 
  visit_var_statement(vs *VarStatement) 
  visit_debug_statement(ps *DebugStatement) 
}

type ExpressionStatement struct {
  expr Expression
}

func (es *ExpressionStatement) Accept(visitor StatementVisitor){
  visitor.visit_expression_statement(es)
}

type VarStatement struct {
  ident Token
  value Expression
}

func (es *VarStatement) Accept(visitor StatementVisitor){
  visitor.visit_var_statement(es)
}

type DebugStatement struct {
  expr Expression
}

func (ps *DebugStatement) Accept(visitor StatementVisitor){
  visitor.visit_debug_statement(ps)
}

