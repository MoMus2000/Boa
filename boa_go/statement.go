package main

type Statement interface{
  Accept(visitor StatementVisitor)
}

type StatementVisitor interface{
  visit_expression_statement(es *ExpressionStatement) 
}

type ExpressionStatement struct {
  expr Expression
}

func (es *ExpressionStatement) Accept(visitor StatementVisitor){
  visitor.visit_expression_statement(es)
}

