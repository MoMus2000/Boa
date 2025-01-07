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
  visit_for_statement       (fs *ForStatement)
  visit_func_statement      (fs *FunctionStatement)
  visit_return_statement    (fs *ReturnStatement) error
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

type ForStatement struct {
  start            Statement
  predicate        Expression
  incre            Expression
  inner_statements *BlockStatement
}

func (fs *ForStatement) Accept(visitor StatementVisitor){
  visitor.visit_for_statement(fs)
}

type FunctionStatement struct {
  ident   Token
  args    []string
  body    *BlockStatement
}

func (fs *FunctionStatement) Accept(visitor StatementVisitor){
  visitor.visit_func_statement(fs)
}

type ReturnStatement struct {
  ident Token
  val   Expression
}

func (fs *ReturnStatement) Accept(visitor StatementVisitor){
  visitor.visit_return_statement(fs)
}

