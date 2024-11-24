from ..lexer import Lexer
from ..token_types import TokenType
from ..tokens import Token

def test():
    lexer = None
    with open("./tests/test_lexer.boa", "r") as source:
        lexer = Lexer(source.readlines())
    tokens = lexer.scan_tokens()
    actual_tokens = [
            Token(TokenType.LEFT_PAREN, lexeme="(", literal=None, line=1),
            Token(TokenType.RIGHT_PAREN, lexeme=")", literal=None, line=1),
            Token(TokenType.PLUS, lexeme="+", literal=None, line=1),
            Token(TokenType.MINUS, lexeme="-", literal=None, line=1),
            Token(TokenType.LEFT_PAREN, lexeme="(", literal=None, line=1),
            Token(TokenType.RIGHT_PAREN,lexeme=")", literal=None, line=1),
            Token(TokenType.GREATER_EQUAL, lexeme=">=", literal=None, line=1),
            Token(TokenType.LESS_EQUAL,lexeme="<=", literal=None, line=1),
            Token(TokenType.BANG_EQUAL,lexeme="!=", literal=None, line=1),
            Token(TokenType.LEFT_PAREN, lexeme="(", literal=None, line=1),
            Token(TokenType.RIGHT_PAREN,lexeme=")", literal=None, line=1),
            Token(TokenType.PLUS,lexeme="+", literal=None, line=1),
            Token(TokenType.MINUS,lexeme="-", literal=None, line=1),
            Token(TokenType.PLUS,lexeme="+", literal=None, line=1),
            Token(TokenType.STRING,lexeme='"hello"', literal=None, line=1),
            Token(TokenType.STRING,lexeme='"say"', literal=None, line=1),
            Token(TokenType.STRING,lexeme='"you won\'t let go"', literal=None, line=1),
            Token(TokenType.NUMBER,lexeme='1.222', literal=None, line=1),
            Token(TokenType.NUMBER,lexeme='0.00001', literal=None, line=1),
            Token(TokenType.EOF, lexeme="", literal=None, line=1),
    ]
    for actual, got in zip(actual_tokens, tokens):
        assert actual.type.value == got.type.value
        assert actual.lexeme == got.lexeme

if __name__ == "__main__":
    test()

