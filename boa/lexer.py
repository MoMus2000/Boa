from token_types import TokenType
from tokens import Token

class Lexer:
    def __init__(self, source):
        self.source = "".join(source)
        self.tokens = []
        self.start = 0
        self.current = 0
        self.line = 0
        self.ident_map = {
            "and"    : TokenType.AND,
            "class"  : TokenType.CLASS,
            "else"   : TokenType.ELSE,
            "false"  : TokenType.FALSE,
            "for"    : TokenType.FOR,
            "fun"    : TokenType.FUN,
            "if"     : TokenType.IF,
            "nil"    : TokenType.NIL,
            "or"     : TokenType.OR,
            "print"  : TokenType.PRINT,
            "return" : TokenType.RETURN,
            "super"  : TokenType.SUPER,
            "this"   : TokenType.THIS,
            "true"   : TokenType.TRUE,
            "var"    : TokenType.VAR,
            "while"  : TokenType.WHILE,
            "import" : TokenType.IMPORT,
        }

    def scan_tokens(self):
        while self.is_not_at_end():
            self.start = self.current
            self.scan_token()
        
        self.tokens.append(
            Token(TokenType.EOF, "", None, self.line)
        )

        return self.tokens
    
    def is_not_at_end(self):
        return self.current < len(self.source) 

    def advance(self):
        c = self.source[self.current]
        self.current += 1
        return c
    
    def add_token(self, token_type, literal=None):
        self.tokens.append(
            Token(token_type, self.source[self.start:self.current], literal, self.line)
        )

    def match(self, expected):
        if not self.is_not_at_end():
            return False
        if self.source[self.current] != expected:
            return False
        self.current += 1
        return True
    
    def peek(self):
        if not self.is_not_at_end():
            return '\0'
        return self.source[self.current]
    
    def string(self):
        inner_string = []
        while self.is_not_at_end() and self.peek() != '"':
            if self.peek() == "\n":
                self.line += 1
            c = self.advance()
            inner_string.append(c)
        self.advance()
        return "".join(inner_string)
    
    def is_digit(self, c):
        return c.isnumeric()

    def number(self):
        while self.is_digit(self.peek()):
            self.advance()

        if self.peek() == "." and self.is_digit(self.peek_next()):
            self.advance()

            while self.is_digit(self.peek()):
                self.advance()

    def peek_next(self):
        if self.current +1 >= len(self.source):
            return '\0'
        return self.source[self.current+1]
    
    def is_alpha(self, c):
        return c.isalpha()

    def is_alphanumeric(self, c):
        if c == "_": return True
        if c == ".": return True
        return c.isalnum()
    
    def identifier(self):
        while self.is_alphanumeric(self.peek()):
            self.advance()

        text = self.source[self.start:self.current]

        if len(text.split(".")) > 2:
            raise Exception(
                "Not supporting dot notation, single level traversal only available"
            )

        ttype = self.ident_map.get(text)

        if ttype == None:
            self.add_token(TokenType.IDENTIFIER)
            return

        self.add_token(ttype)
    
    def scan_token(self):
        c = self.advance()
        if   c == "(":
            self.add_token(TokenType.LEFT_PAREN)
        elif c == ")":
            self.add_token(TokenType.RIGHT_PAREN)
        elif c == "{":
            self.add_token(TokenType.LEFT_BRACE)
        elif c == "}":
            self.add_token(TokenType.RIGHT_BRACE)
        elif c == ",":
            self.add_token(TokenType.COMMA)
        elif c == ".":
            self.add_token(TokenType.DOT)
        elif c == "-":
            self.add_token(TokenType.MINUS)
        elif c == "+":
            self.add_token(TokenType.PLUS)
        elif c == ";":
            self.add_token(TokenType.SEMICOLON)
        elif c == "*":
            self.add_token(TokenType.STAR)
        elif c == "[":
            self.add_token(TokenType.LEFT_ANGLE_BRACKET)
        elif c == "]":
            self.add_token(TokenType.RIGHT_ANGLE_BRACKET)
        elif c == ":":
            self.add_token(TokenType.COLON)
        elif c == "!":
            if self.match("="):
                self.add_token(TokenType.BANG_EQUAL)
            else:
                self.add_token(TokenType.BANG)
        elif c == "=":
            if self.match("="):
                self.add_token(TokenType.EQUAL_EQUAL)
            else:
                self.add_token(TokenType.EQUAL)
        elif c == "<":
            if self.match("="):
                self.add_token(TokenType.LESS_EQUAL)
            else:
                self.add_token(TokenType.LESS)
        elif c == ">":
            if self.match("="):
                self.add_token(TokenType.GREATER_EQUAL)
            else:
                self.add_token(TokenType.GREATER)
        elif c == "|":
            if self.match(">"):
                self.add_token(TokenType.PIPE)
            else:
                raise Exception(f"Unexpected character {repr(c)}")
        elif c == " ":
            pass
        elif c == "/":
            if self.match("/"):
                while self.peek() != "\n" and self.is_not_at_end():
                    self.advance()
            else:
                self.add_token(TokenType.SLASH)
        elif c == "\n":
            self.line += 1
        elif c == "\t":
            pass
        elif c == "\r":
            pass
        elif c == '"':
            self.string()
            self.add_token(TokenType.STRING, literal=str)
        elif self.is_digit(c):
            self.number()
            self.add_token(TokenType.NUMBER, literal=float)
        elif self.is_alpha(c):
            self.identifier()
        else:
            raise Exception(f"Unexpected character {repr(c)}")

