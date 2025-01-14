package main

type Statement interface{
  Accept(visitor StatementVisitor) error
}

type StatementVisitor interface{
  visit_expression_statement(es  *ExpressionStatement) error
  visit_var_statement       (vs  *VarStatement) error
  visit_debug_statement     (ps  *DebugStatement) error
  visit_block_statement     (bs  *BlockStatement) error
  visit_if_statement        (ifs *IfStatement) error
  visit_while_statement     (ws *WhileStatement) error
  visit_for_statement       (fs *ForStatement) error
  visit_func_statement      (fs *FunctionStatement) error
  visit_return_statement    (fs *ReturnStatement) error
  visit_enhanced_loop_statement (es *EnhancedLoop) error
}

type ExpressionStatement struct {
  expr Expression
}

func (es *ExpressionStatement) Accept(visitor StatementVisitor) error {
  visitor.visit_expression_statement(es)
  return nil
}

type DebugStatement struct {
  expr Expression
}

func (ps *DebugStatement) Accept(visitor StatementVisitor) error {
  visitor.visit_debug_statement(ps)
  return nil
}

type BlockStatement struct {
  statements []Statement
}

func (bs *BlockStatement) Accept(visitor StatementVisitor) error {
  visitor.visit_block_statement(bs)
  return nil
}

type VarStatement struct {
  ident Token     
  value Expression
}

func (v *VarStatement) Accept(visitor StatementVisitor) error {
  visitor.visit_var_statement(v)
  return nil
}

type IfStatement struct {
  predicate      Expression     
  if_condition   *BlockStatement
  else_condition *BlockStatement
}

func (ifs *IfStatement) Accept(visitor StatementVisitor) error {
  visitor.visit_if_statement(ifs)
  return nil
}

type WhileStatement struct {
  predicate        Expression
  inner_statements *BlockStatement
}

func (ws *WhileStatement) Accept(visitor StatementVisitor) error {
  visitor.visit_while_statement(ws)
  return nil
}

type ForStatement struct {
  start            Statement
  predicate        Expression
  incre            Expression
  inner_statements *BlockStatement
}

func (fs *ForStatement) Accept(visitor StatementVisitor) error {
  visitor.visit_for_statement(fs)
  return nil
}

type FunctionStatement struct {
  ident   Token
  args    []string
  body    *BlockStatement
}

func (fs *FunctionStatement) Accept(visitor StatementVisitor) error {
  visitor.visit_func_statement(fs)
  return nil
}

type ReturnStatement struct {
  ident Token
  val   Expression
}

func (fs *ReturnStatement) Accept(visitor StatementVisitor) error {
  visitor.visit_return_statement(fs)
  return nil
}

type EnhancedLoop struct {
  iteratee Statement
  iterator Iterator
}

func (es *EnhancedLoop) Accept(visitor StatementVisitor) error {
  visitor.visit_enhanced_loop_statement(es)
  return nil
}

