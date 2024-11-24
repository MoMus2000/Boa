from ..lexer import Lexer
from ..parser import Parser
from interpreter import Interpreter

def test():
    lexer = None
    tests = [
        3,
        "hello kitty jungle party",
        True,
        False,
        True
    ]
    with open("./tests/test_interpreter.boa", "r") as source:
        for test, line in zip(tests, source.readlines()):
            lexer  = Lexer(line)
            tokens = lexer.scan_tokens()
            parser = Parser(tokens)
            statements = parser.parse()
            result = Interpreter().interpret(statements)
            assert result[0] == test

if __name__ == "__main__":
    test()

