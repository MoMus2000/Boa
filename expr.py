from abc import ABC, abstractmethod
from token_types import TokenType
from tokens import Token

class Expr(ABC):
    def __init__(self):
        pass
    
    @abstractmethod
    def accept(self, visitor):
        pass

class Visitor(ABC):
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

class AstPrinter(Visitor):
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

if __name__ == "__main__":
    expressions = [ 
        Binary(
            left = Unary(
                Token(
                    token_type=TokenType.MINUS, lexeme="-", literal=None, line=0
                ),
                Literal(12)
            ),
            op = Token(
                    token_type=TokenType.MINUS, lexeme="+", literal=None, line=0
                ),
            right=
                Literal(12)
            ),
        Grouping(
            Unary(
                Token(TokenType.PLUS, lexeme="+", literal=None, line = 0),
                Literal(12)
        )),
        Binary(
        left=Grouping(
            Binary(
            left = Literal(1),
            op = Token(
                    token_type=TokenType.PLUS, lexeme="+", literal=None, line=0
                ),
            right = Literal(2)
            )
        ),
        op = Token(
                    token_type=TokenType.STAR, lexeme="*", literal=None, line=0
                ),
        right=Unary(Token(TokenType.MINUS, lexeme="-", literal=None, line = 0),
            Literal(3)
        )
        ),
        Binary(
            Unary(
                Token(TokenType.MINUS, "-", None, 1),
                Literal(123)
            ),
            Token(TokenType.STAR, "*", None, 1),
            Grouping(
                Literal(45.67)
            )
        )
    ]

    ast_printer = AstPrinter()
    for expr in expressions:
        print(ast_printer.print(expr))

