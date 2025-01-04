from abc import ABC, abstractmethod

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
        self.value = ident # Mirroring ident for literal

    def accept(self, visitor):
        return visitor.visit_var_expression(self)

class Assign(Expr):
    def __init__(self, ident, value):
        self.ident = ident
        self.value = value

    def accept(self, visitor):
        return visitor.visit_assign_expression(self)

class Logical(Expr):
    def __init__(self, left, op, right):
        self.left  = left
        self.op    = op
        self.right = right

    def accept(self, visitor):
        return visitor.visit_logical_expression(self)


class Call(Expr):
    def __init__(self, callee, paren, args):
        self.callee = callee
        self.paren  = paren
        self.args   = args

    def accept(self, visitor):
        return visitor.visit_call_expression(self)

class Arr(Expr):
    def __init__(self, elements):
        self.elements = elements

    def accept(self, visitor):
        return visitor.visit_array_expression(self)

