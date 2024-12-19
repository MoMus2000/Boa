import io
import sys
from environment import Environment
from expr import ExprVisitor
from statement import StmtVisitor
from token_types import TokenType

def clock():
    from datetime import datetime
    return f"{datetime.now()}"

class Callable:
    def __init__(self, func, arity):
        self.func  = func
        self.arity = arity

    def call(self):
        return self.func()

class CallableFunc:
    def __init__(self, decl, arity):
        self.decl = decl
        self.arity = arity

    def call(self, interpreter, args):
        env = Environment(interpreter.env)

        for param, arg in zip(self.decl.params, args):
            env.define(param.lexeme, arg)

        interpreter.execute_block(self.decl.body, env)
        return None

class Interpreter(StmtVisitor, ExprVisitor):
    def __init__(self, debug_mode=False):
        self.statements = []
        self.env = Environment()
        self.env.define("clock", Callable(clock, 0))
        self.output = io.StringIO()
        self.saved_stdout = sys.stdout
        self.debug_mode = debug_mode

    def __enter__(self):
        if self.debug_mode:
            sys.stdout = self.output
        return self

    def __exit__(self):
        if self.debug_mode:
            sys.stdout = self.saved_stdout

    def get_output(self):
        return self.output.getvalue()

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

    def visit_if_statement(self, ifstmt):
        pred = self.evaluate(ifstmt.predicate)
        else_block = ifstmt.else_block
        if pred:
            self.visit_block_statement(ifstmt.block)
        if else_block and not pred:
            self.visit_block_statement(ifstmt.else_block)

    def visit_while_statement(self, whilestmt):
        while self.evaluate(whilestmt.predicate) == True:
            self.visit_block_statement(whilestmt.block)

    def visit_loop_statement(self, forstmt):
        self.visit_var_statement(forstmt.start)
        while self.evaluate(forstmt.predicate) == True:
            self.visit_block_statement(forstmt.block)
            self.evaluate(forstmt.incrementer)

    def visit_var_statement(self, stmt):
        identifier = stmt.ident
        if stmt.expression != None:
            val = self.evaluate(stmt.expression)
            self.env.define(identifier.lexeme, val)
        else:
            self.env.define(identifier.lexeme, None)
        return self.env.get(identifier.lexeme)

    def visit_block_statement(self, block):
        return self.execute_block(block, Environment(self.env))

    def visit_logical_expression(self, logicalexpr):
        left = self.evaluate(logicalexpr.left)

        if logicalexpr.op.type == TokenType.OR:
            if self.is_truthy(left):
                return left
            elif not self.is_truthy(left):
                return left

        return self.evaluate(logicalexpr.right)

    def execute_block(self, block, env):
        prev = env
        self.env = env
        res = []
        for statement in block.statements:
            res.append(self.evaluate(statement))
        self.env = prev
        return res

    def visit_assign_expression(self, expr):
        self.env.assign(expr.ident, self.evaluate(expr.value))
        return self.env.get(expr.ident.lexeme)

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

    def visit_call_expression(self, expr):
        if len(expr.args) >= 128:
            raise Exception("To many arguments have been provided")
        func = self.evaluate(expr.callee)
        if func is None or not (isinstance(func, Callable) or isinstance(func,
                                                                         CallableFunc)):
            raise Exception(f"Error with defined func")
        if func.arity != len(expr.args):
            raise Exception(f"Error with the defined number of args")

        if isinstance(func, CallableFunc):
            return func.call(self, expr.args)
        return func.call()
    
    def visit_func_statement(self, visitor):
        func = CallableFunc(visitor, len(visitor.params))
        self.env.define(visitor.token_name.lexeme, func)
        return None

    def is_truthy(self, expr):
        if expr == None:
            return False
        if isinstance(expr, bool):
            return bool(expr)
        return True
    
    def evaluate(self, expr):
        return expr.accept(self)

