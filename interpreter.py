from token_types import TokenType
from expr import (
    ExprVisitor, 
)
from statement import (
    StmtVisitor
)
from environment import Environment

class Interpreter(StmtVisitor, ExprVisitor):
    def __init__(self):
        self.statements = []
        self.env = Environment()

    def interpret(self, statements):
        result = []
        for statement in statements:
            r = self.execute_statement(statement)
            result.append(r)
        return result

    def execute_statement(self, stmt):
        return stmt.accept(self)

    def parse_to_float(self, value):
        try:
            return float(value), True
        except ValueError:
            return value, False

    def visit_expression_statement(self, stmt):
        # return None In reality expression statements produce no values 
        # (We return for testing)
        return self.evaluate(stmt.expression)

    def visit_print_statement(self, stmt):
        val = self.evaluate(stmt.expression)
        print(val)
        return None

    def visit_var_statement(self, stmt):
        identifier = stmt.ident
        print(self.env.map)
        if stmt.expression != None:
            val = self.evaluate(stmt.expression)
            self.env.define(identifier.lexeme, val)
        else:
            self.env.define(identifier.lexeme, None)
        print(self.env.map)

    def visit_binary_expression(self, expr):
        left  = self.evaluate(expr.left)
        right = self.evaluate(expr.right)
        left,  l_parsed = self.parse_to_float(left)
        right, r_parsed = self.parse_to_float(right)

        if l_parsed and r_parsed:
            if expr.op.type == TokenType.PLUS:
                return left + right
            if expr.op.type == TokenType.MINUS:
                return left - right
            if expr.op.type == TokenType.SLASH:
                return left / right
            if expr.op.type == TokenType.STAR:
                return left * right
            if expr.op.type == TokenType.EQUAL_EQUAL:
                return left == right
            if expr.op.type == TokenType.BANG_EQUAL:
                return left != right
            if expr.op.type == TokenType.GREATER:
                return left > right
            if expr.op.type == TokenType.LESS:
                return left < right
            if expr.op.type == TokenType.LESS_EQUAL:
                return left <= right
            if expr.op.type == TokenType.GREATER_EQUAL:
                return left >= right
        else:
            if expr.op.type == TokenType.PLUS:
                return str(left).strip('"') + str(right).strip('"') 
            if expr.op.type == TokenType.BANG_EQUAL:
                return left != right
            if expr.op.type == TokenType.EQUAL_EQUAL:
                return left == right
        return None

    def visit_literal_expression(self, expr):
        return expr.value

    def visit_grouping_expression(self, expr):
        return self.evaluate(expr.expression)

    def visit_unary_expression(self, expr):
        right = self.evaluate(expr.right)
        if expr.op.type == TokenType.MINUS:
            return float(right) * -1
        if expr.op.type == TokenType.BANG:
            return not self.is_truthy(right)
        return None

    def visit_var_expression(self, expr):
        return self.env.get(expr.ident.lexeme)

    def is_truthy(self, expr):
        if expr == None:
            return False
        if isinstance(expr, bool):
            return bool(expr)
        return True
    
    def evaluate(self, expr):
        return expr.accept(self)

