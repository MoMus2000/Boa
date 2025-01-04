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

class FuncStmt(Statement):
    def __init__(self, token_name, params, body):
        self.token_name = token_name
        self.params = params
        self.body = body

    def accept(self, visitor):
        return visitor.visit_func_statement(self)

class ReturnStmt(Statement):
    def __init__(self, keyword, value):
        self.keyword = keyword
        self.value = value

    def accept(self, visitor):
        return visitor.visit_return_statement(self)

class ImportStmt(Statement):
    def __init__(self, lib_name):
        self.lib_name = lib_name

    def accept(self, visitor):
        return visitor.visit_import_statement(self)

class ArrayStmt(Statement):
    def __init__(self, ident, elements=None, index=None):
        if index is not None and elements is not None:
            raise Exception("Collision Error")
        self.ident = ident
        self.elements = elements
        self.index = index

    def accept(self, visitor):
        return visitor.visit_array_statement(self)

class ArrayAssignStmt(Statement):
    def __init__(self, ident, index, value):
        self.ident = ident
        self.index = index
        self.value = value

    def accept(self, visitor):
        return visitor.visit_array_assign_statement(self)

class HashMapStatement(Statement):
    def __init__(self, keys, values):
        self.keys  = keys
        self.values = values

    def accept(self, visitor):
        return visitor.visit_hash_map_statement(self)

