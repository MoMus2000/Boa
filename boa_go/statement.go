package main

type Statement interface{
  Accept(visitor StatementVisitor)
}

type StatementVisitor interface{
  visit_expression_statement(es  *ExpressionStatement)
  visit_var_statement       (vs  *VarStatement)
  visit_debug_statement     (ps  *DebugStatement)
  visit_block_statement     (bs  *BlockStatement)
  visit_if_statement        (ifs *IfStatement)
  visit_while_statement     (ws *WhileStatement)
}

type ExpressionStatement struct {
  expr Expression
}

func (es *ExpressionStatement) Accept(visitor StatementVisitor){
  visitor.visit_expression_statement(es)
}

type DebugStatement struct {
  expr Expression
}

func (ps *DebugStatement) Accept(visitor StatementVisitor){
  visitor.visit_debug_statement(ps)
}

type BlockStatement struct {
  statements []Statement
}

func (bs *BlockStatement) Accept(visitor StatementVisitor){
  visitor.visit_block_statement(bs)
}

type VarStatement struct {
  ident Token     
  value Expression
}

func (v *VarStatement) Accept(visitor StatementVisitor){
  visitor.visit_var_statement(v)
}

type IfStatement struct {
  predicate      Expression     
  if_condition   *BlockStatement
  else_condition *BlockStatement
}

func (ifs *IfStatement) Accept(visitor StatementVisitor){
  visitor.visit_if_statement(ifs)
}

type WhileStatement struct {
  predicate        Expression
  inner_statements *BlockStatement
}

func (ws *WhileStatement) Accept(visitor StatementVisitor){
  visitor.visit_while_statement(ws)
}

