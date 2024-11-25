from abc import ABC, abstractmethod
from token_types import TokenType
from tokens import Token

class Expr(ABC):
    def __init__(self):
        pass
    
    @abstractmethod
    def accept(self, visitor):
        pass

class ExprVisitor(ABC):
    @abstractmethod
    def visit_binary_expression(self, expr) -> object:
        pass

    @abstractmethod
    def visit_literal_expression(self, expr) -> object:
        pass

    @abstractmethod
    def visit_grouping_expression(self, expr) -> object:
        pass

    @abstractmethod
    def visit_unary_expression(self, expr) -> object:
        pass

class AstPrinter(ExprVisitor):
    def visit_binary_expression(self, expr):
        return f"({expr.op.lexeme} {expr.left.accept(self)} {expr.right.accept(self)})"

    def visit_literal_expression(self, expr):
        return str(expr.value)

    def visit_grouping_expression(self, expr):
        return f"(group {expr.expression.accept(self)})"

    def visit_unary_expression(self, expr):
        return f"({expr.op.lexeme}{expr.right.accept(self)})"

    def print(self, expr):
        return expr.accept(self)

class Binary(Expr):
    def __init__(self, left, op, right):
        self.left  = left
        self.op    = op
        self.right = right

    def accept(self, visitor):
        return visitor.visit_binary_expression(self)

class Literal(Expr):
    def __init__(self, value):
        self.value = value

    def accept(self, visitor):
        return visitor.visit_literal_expression(self)

class Grouping(Expr):
    def __init__(self, expression):
        self.expression = expression

    def accept(self, visitor):
        return visitor.visit_grouping_expression(self)

class Unary(Expr):
    def __init__(self, op, right):
        self.op    = op
        self.right = right

    def accept(self, visitor):
        return visitor.visit_unary_expression(self)

class Var(Expr):
    def __init__(self, ident):
        self.ident = ident

    def accept(self, visitor):
        return visitor.visit_var_expression(self)

