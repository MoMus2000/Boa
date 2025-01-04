import io
import sys
from stdlib.boa_math import Math
from stdlib.boa_time import Time
from stdlib.boa_arr  import Arr
from stdlib.boa_map  import Map
from environment import Environment
from expr import ExprVisitor, Literal, Var
from statement import StmtVisitor, ArrayStmt
from token_types import TokenType

def assert_eq(a, b, message):
    assert a == b, message

def assert_boa(a, message):
    assert a, message


class Callable:
    def __init__(self, func, arity):
        self.func  = func
        self.arity = arity

    def call(self, args):
        return self.func(*args)

class CallableFunc:
    def __init__(self, decl, arity):
        self.decl = decl
        self.arity = arity

    def call(self, interpreter, args):
        env = Environment(interpreter.globals)
        for param, arg in zip(self.decl.params, args):
            env.define(param.lexeme, arg)
        try:
            interpreter.execute_block(self.decl.body, env)
        except ReturnException as r:
            return r.return_statement
        return None

class Interpreter(StmtVisitor, ExprVisitor):
    def __init__(self, debug_mode=False):
        self.statements = []
        self.globals = Environment()
        self.env = self.globals
        self.globals.define("assert_eq", Callable(assert_eq, 3))
        self.globals.define("assert", Callable(assert_boa, 2))
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
            if isinstance(value, Var):
                value = self.env.get(value.ident.lexeme)
            elif isinstance(value, Literal):
                value = value.value
            return round(float(value), 3), True
        except Exception:
            return value, False

    def visit_expression_statement(self, stmt):
        # return None In reality expression statements produce no values 
        # (We return for testing)
        return self.evaluate(stmt.expression)

    def visit_print_statement(self, stmt):
        val = self.evaluate(stmt.expression)
        if (isinstance(val, bool) or val) and self.parse_to_float(val)[1]:
            try:
                print(round(val, 3))
            except Exception:
                print(val)
        elif val == None:
            print ("nil")
        else:
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

    def evaluate_array(self, array):
        evaluated = []
        for element in array:
            if isinstance(element, list):
                evaluated.append(self.evaluate_array(element))
            elif hasattr(element, "lexeme"):
                if element.type == TokenType.IDENTIFIER:
                    evaluated.append(self.env.get(element.lexeme))
                else:
                    evaluated.append(element.lexeme)
            else:
                evaluated.append(self.evaluate(element))
        return evaluated

    def visit_var_statement(self, stmt):
        from functools import reduce
        identifier = stmt.ident
        if stmt.expression != None:
            if isinstance(stmt.expression, ArrayStmt):
                val = stmt.expression
                if val.elements != None:
                    val = self.evaluate_array(val.elements)
                else:
                    index = [int(self.evaluate(x)) for x in val.index]
                    val = reduce(lambda acc, idx: acc[idx], index, 
                        self.env.get(stmt.expression.ident.lexeme))
                self.env.define(identifier.lexeme, val)
            else:
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
        else:
            if not self.is_truthy(left):
                return left

        return self.evaluate(logicalexpr.right)

    def execute_block(self, block, env):
        prev = self.env
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
            if isinstance(left, list) and isinstance(right, list) and expr.op.type == TokenType.PLUS:
                return left + right
            if expr.op.type == TokenType.PLUS:
                l = str(left).strip('"')
                r = str(right).strip('"')

                l = l.replace("None", "nil")
                r = r.replace("None", "nil")

                f = l + r
                return f'"{f}"'
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
        args = []
        for arg in expr.args:
            args.append(self.evaluate(arg))
        if isinstance(func, CallableFunc):
            import copy
            return func.call(copy.copy(self), args)
        return func.call(args)
    
    def visit_func_statement(self, visitor):
        func = CallableFunc(visitor, len(visitor.params))
        self.env.define(visitor.token_name.lexeme, func)
        return None

    def visit_return_statement(self, visitor):
        val = None
        if visitor.value != None:
            val = self.evaluate(visitor.value)
        raise ReturnException(val)

    def visit_array_statement(self, visitor):
        from functools import reduce
        try:
            getter = self.env.get(visitor.ident.lexeme)
            if isinstance(getter, list):
                if isinstance(visitor.index, Var):
                    index = self.env.get(visitor.index.ident.lexeme)
                else:
                    index = [int(self.evaluate(x)) for x in visitor.index]
                return reduce(lambda acc, idx: acc[idx], index, getter)
            elif isinstance(getter, dict):
                if isinstance(visitor.index, Var):
                    index = self.env.get(visitor.index.ident.lexeme)
                else:
                    index = []
                    for x in visitor.index:
                        x = self.evaluate(x)
                        if type(x) == float:
                            index.append(int(x))
                        else:
                            index.append(x)
                return reduce(lambda acc, idx: acc[idx], index, getter)
        except Exception:
            raise IndexError("Array does not contain index ", visitor.index)

    def visit_array_assign_statement(self, visitor):
        value = self.evaluate(visitor.value)
        array = self.env.get(visitor.ident.lexeme)
        indices = [int(self.evaluate(x)) for x in visitor.index]
        for idx in indices[:-1]:
            array = array[idx]
        # Assign the value at the final index
        array[indices[-1]] = value
        return None

    def visit_hash_map_statement(self, visitor):
        keys   = [self.evaluate(k) for k in visitor.keys]
        values = [self.evaluate(v) for v in visitor.values]
        if None in keys:
            raise Exception("nil values are not allowed as keys in hashmap")
        return dict(zip(keys, values))

    def visit_array_expression(self, visitor):
        return [self.evaluate(v) for v in visitor.elements]
    
    def visit_import_statement(self, visitor):
        allowed_libs = ["math", "arr", "time", "map"]
        if visitor.lib_name.lexeme not in allowed_libs:
            raise Exception(f"import {visitor.lib_name.lexeme} not found")
        if visitor.lib_name.lexeme == "map":
            map = Map()
            self.globals.define("map.insert", Callable(map.insert, 3))
            self.globals.define("map.keys", Callable(map.keys, 1))
            self.globals.define("map.values", Callable(map.values, 1))
            self.globals.define("map.get", Callable(map.get, 2))
        if visitor.lib_name.lexeme == "math":
            math = Math()
            self.globals.define("math.pi", Callable(math.pi, 0))
            self.globals.define("math.pow", Callable(math.pow, 2))
            self.globals.define("math.factorial", Callable(math.factorial, 1))
            self.globals.define("math.ceil", Callable(math.ceil, 1))
            self.globals.define("math.floor", Callable(math.floor, 1))
            self.globals.define("math.is_prime", Callable(math.is_prime, 1))
            self.globals.define("math.min", Callable(math.min, 2))
            self.globals.define("math.max", Callable(math.max, 2))
            self.globals.define("math.sqrt", Callable(math.sqrt, 1))
            self.globals.define("math.abs", Callable(math.abs, 1))
            self.globals.define("math.random_f", Callable(math.random_f, 0))
        if visitor.lib_name.lexeme == "time":
            time = Time()
            self.globals.define("time.sleep", Callable(time.sleep, 1))
            self.globals.define("time.clock", Callable(time.clock, 0))
        if visitor.lib_name.lexeme == "arr":
            arr = Arr()
            self.globals.define("arr.length", Callable(arr.length, 1))
            self.globals.define("arr.append", Callable(arr.append, 2))
            self.globals.define("arr.pop", Callable(arr.pop, 1))
            self.globals.define("arr.modify", Callable(arr.modify, 3))
        return

    def is_truthy(self, expr):
        if expr == None:
            return False
        if isinstance(expr, bool):
            return bool(expr)
        return True
    
    def evaluate(self, expr):
        return expr.accept(self)

class ReturnException(Exception):
    def __init__(self, return_statement):
        self.return_statement = return_statement
