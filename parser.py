from token_types import TokenType
from expr import Binary, Unary, Literal, Grouping, Var as ExprVar, Assign, Logical, Call
from statement import Print, Expression, Var, Block, IfStmt, WhileStmt, ForLoopStmt, FuncStmt
"""
Order of precedence  
expression → equality ;                                     (Lowest precedence)
equality   → comparison ( ( "!=" | "==" ) comparison )* ;           |
comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;           |
term       → factor ( ( "-" | "+" ) factor )* ;                     |
factor     → unary ( ( "/" | "*" ) unary )* ;                       |
unary      → ( "!" | "-" ) unary                                    |
           | primary ;                                              |
primary    → NUMBER | STRING | "true" | "false" | "nil"             |
           | "(" expression ")" ;                           (Highest precedence)

The parser that we use is recursive descent

 _____________________________________________________
|                                                     |
|Grammar notation	Code representation               |
|-----------------------------------------------------|
|Terminal	        Code to match and consume a token |
|Nonterminal	    Call to that rule’s function      |
||	                if or switch statement            |
|* or +	            while or for loop                 |
|?	                if statement                      |
|_____________________________________________________|
"""

def todo(name):
    raise Exception(f"Not Implemented {name}")

class Parser:
    def __init__(self, tokens):
        self.current = 0
        self.tokens  = tokens

    def parse(self):
        """
            Kick Off parsing
        """
        self.statements = []
        while not self.is_at_end():
            self.statements.append(
                self.declaration()
            )
        return self.statements 

    def statement(self):
        if self.match(TokenType.PRINT):
            return self.print_statement()
        if self.match(TokenType.LEFT_BRACE):
            return self.block()

        return self.expression_statement()

    def print_statement(self):
        value = self.expression()
        self.consume(TokenType.SEMICOLON, "Expected ; after value")
        return Print(value)

    def block(self):
        statements = []
        while not self.check(TokenType.RIGHT_BRACE) and not self.is_at_end():
            statement = self.declaration()
            statements.append(statement)
        self.consume(TokenType.RIGHT_BRACE, "Expected '}' ")
        return Block(statements)
    
    def var_statement(self):
        ident = self.consume(TokenType.IDENTIFIER, "Expected Variable Name");
        init = None
        if self.match(TokenType.EQUAL):
            init = self.expression()
        self.consume(TokenType.SEMICOLON, "Expected ; after value");
        return Var(init, ident)

    def if_statement(self):
        self.consume(TokenType.LEFT_PAREN, "Expected ( after if")
        predicate = self.expression()
        self.consume(TokenType.RIGHT_PAREN, "Expected ) after if")
        self.consume(TokenType.LEFT_BRACE, "Expected { after )")
        block = self.block()
        else_block = None
        if self.match(TokenType.ELSE):
            self.consume(TokenType.LEFT_BRACE, "Expected { after else")
            else_block = self.block()
        return IfStmt(predicate, block, else_block)

    def while_statement(self):
        self.consume(TokenType.LEFT_PAREN, "Expected ( after if")
        predicate = self.expression()
        self.consume(TokenType.RIGHT_PAREN, "Expected ) after if")
        self.consume(TokenType.LEFT_BRACE, "Expected { after )")
        block = self.block()
        return WhileStmt(predicate, block)

    def for_loop_statement(self):
        self.consume(TokenType.LEFT_PAREN, "Expected ( after if")
        start = None
        if self.match(TokenType.VAR):
            start = self.var_statement()
        predi = self.expression_statement()
        incre = self.expression_statement()
        self.consume(TokenType.RIGHT_PAREN, "Expected ) after if")
        self.consume(TokenType.LEFT_BRACE, "Expected { after )")
        block = self.block()
        return ForLoopStmt(start, predi, incre, block)

    def declaration(self):
        if self.match(TokenType.FUN):
            return self.define_fun_statement()
        if self.match(TokenType.VAR):
            return self.var_statement()
        if self.match(TokenType.FOR):
            return self.for_loop_statement()
        if self.match(TokenType.IF):
            return self.if_statement()
        if self.match(TokenType.WHILE):
            return self.while_statement()
        if self.match(TokenType.FUN):
            return self.define_fun_statement()
        return self.statement()

    def define_fun_statement(self):
        name = self.consume(TokenType.IDENTIFIER, "Expected Identifier")
        self.consume(TokenType.LEFT_PAREN, "Expected Left Paren")
        args = []
        if not self.check(TokenType.RIGHT_PAREN):
            while True:
                if self.peek().type == TokenType.IDENTIFIER:
                    arg = self.consume(TokenType.IDENTIFIER, "Expected args")
                    args.append(arg)
                if not self.match(TokenType.COMMA):
                    break
        self.consume(TokenType.RIGHT_PAREN, "Expected RIGHT Paren")
        self.consume(TokenType.LEFT_BRACE, "Expected LEFT Brace")
        block = self.block()
        return FuncStmt(
            name,
            args,
            block
        )

    def expression_statement(self):
        value = self.expression()
        self.consume(TokenType.SEMICOLON, "Expected ; after value")
        return Expression(value)

    def and_expr(self):
        expr = self.equality()
        while self.match(TokenType.AND):
            op = self.previous()
            right = self.equality()
            expr = Logical(expr, op, right)

        return expr

    def or_expr(self):
        expr = self.and_expr()
        if self.match(TokenType.OR):
            operator = self.previous()
            right = self.and_expr()
            expr = Logical(expr, operator, right)
        return expr

    def assign(self):
        expr = self.or_expr()

        if self.match(TokenType.EQUAL):
            _ = self.previous()
            value  = self.assign()
            if isinstance(expr, ExprVar):
               return Assign(expr.ident, value)
            else:
                raise Exception("Invalid Expression Type")
        return expr

    def expression(self):
        return self.assign()

    def equality(self):
        expr = self.comparision()

        while self.match(TokenType.BANG_EQUAL, TokenType.EQUAL_EQUAL):
            op  = self.previous()
            right = self.comparision()
            expr = Binary(
                left = expr,
                op = op,
                right = right
            )

        return expr

    def comparision(self):
        expr = self.term()

        while self.match(TokenType.GREATER, TokenType.GREATER_EQUAL, TokenType.LESS,
                         TokenType.LESS_EQUAL):
            op = self.previous()
            right = self.term()
            expr = Binary(
                    left = expr,
                    op = op,
                    right = right)

        return expr

    def term(self):
        expr = self.factor()

        while self.match(TokenType.PLUS, TokenType.MINUS):
            op = self.previous()
            right = self.factor()
            expr = Binary(
                    left = expr,
                    op = op,
                    right = right)

        return expr

    def factor(self):
        expr = self.unary()

        while self.match(TokenType.SLASH, TokenType.STAR):
            op = self.previous()
            right = self.unary()
            expr = Binary(
                    left = expr,
                    op = op,
                    right = right)

        return expr

    def unary(self):
        if self.match(TokenType.MINUS, TokenType.PLUS, TokenType.BANG):
            op = self.previous()
            right = self.unary()
            return Unary(
                    op = op,
                    right = right
            )

        return self.call()

    def call(self):
        expr = self.primary()

        while self.match(TokenType.LEFT_PAREN):
            expr = self.finish_call(expr)

        return expr

    def finish_call(self, callee):
        args = []
        if not self.check(TokenType.RIGHT_PAREN):
            while True:
                args.append(self.expression())
                if not self.match(TokenType.COMMA):
                    break

        paren = self.consume(TokenType.RIGHT_PAREN, "Expected ')' after args")
        return Call(callee, paren, args)

    def primary(self):
        if self.match(TokenType.FALSE):
            return Literal(False)
        if self.match(TokenType.TRUE):
            return Literal(True)
        if self.match(TokenType.STRING, TokenType.NUMBER):
            return Literal(self.previous().lexeme)
        if self.match(TokenType.NIL):
            return Literal(None)
        if self.match(TokenType.LEFT_PAREN):
            expr = self.expression()
            self.consume(TokenType.RIGHT_PAREN, "Expected ')' after expression")
            return Grouping(expr)
        if self.match(TokenType.IDENTIFIER):
            expr = self.previous()
            return ExprVar(expr)
        raise Exception("expected an expression")

    def synchronize(self):
        self.advance()

        while not self.is_at_end():
            if self.previous().type == TokenType.SEMICOLON:
                return
            
            if self.peek().type == TokenType.CLASS:
                return
            if self.peek().type == TokenType.FUN:
                return
            if self.peek().type == TokenType.VAR:
                return
            if self.peek().type == TokenType.FOR:
                return
            if self.peek().type == TokenType.IF:
                return
            if self.peek().type == TokenType.WHILE:
                return
            if self.peek().type == TokenType.PRINT:
                return
            if self.peek().type == TokenType.RETURN:
                return
            self.advance()

    def consume(self, token, message):
        if self.check(token):
            return self.advance()
        raise Exception(message)

    def match(self, *tokens):
        for token in tokens:
            if self.check(token):
                self.advance()
                return True
        return False
    
    def check(self, token):
        if self.is_at_end():
            return False
        return self.peek().type == token

    def is_at_end(self):
        return self.peek().type == TokenType.EOF

    def peek(self):
        return self.tokens[self.current]

    def advance(self):
        if not self.is_at_end():
            self.current += 1
        return self.previous()

    def previous(self):
        return self.tokens[self.current-1]
