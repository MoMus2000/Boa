from abc import ABC, abstractmethod
"""
program        → declaration* EOF ;

declaration    → varDecl
               | statement ;

statement      → exprStmt
               | printStmt ;
"""

class Statement(ABC):
    def __init__(self):
        pass

    @abstractmethod
    def accept(self, visitor):
        pass

class StmtVisitor(ABC):
    @abstractmethod
    def visit_print_statement(self, stmt) -> object:
        pass

    @abstractmethod
    def visit_expression_statement(self, stmt) -> object:
        pass

class Expression(Statement):
    def __init__(self, expression):
        self.expression = expression

    def accept(self, visitor):
        return visitor.visit_expression_statement(self)

class Print(Statement):
    def __init__(self, expression):
        self.expression = expression

    def accept(self, visitor):
        return visitor.visit_print_statement(self)

class Var(Statement):
    def __init__(self, expression, ident):
        self.expression = expression
        self.ident = ident

    def accept(self, visitor):
        return visitor.visit_var_statement(self)

class Block(Statement):
    def __init__(self, statements):
        self.statements = statements

    def accept(self, visitor):
        return visitor.visit_block_statement(self)

class IfStmt(Statement):
    def __init__(self, predicate, block, else_block = None):
        self.predicate = predicate
        self.block = block 
        self.else_block = else_block

    def accept(self, visitor):
        return visitor.visit_if_statement(self)

class WhileStmt(Statement):
    def __init__(self, predicate, block):
        self.predicate = predicate
        self.block = block 

    def accept(self, visitor): return visitor.visit_while_statement(self)

class ForLoopStmt(Statement):
    def __init__(self, start, predicate, incrementer, block):
        self.start = start
        self.predicate = predicate
        self.incrementer = incrementer
        self.block = block

    def accept(self, visitor):
        return visitor.visit_loop_statement(self)

