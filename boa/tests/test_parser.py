from ..lexer import Lexer
from ..parser import Parser
from ..expr import *

def test():
    lexer = None
    printer = AstPrinter()
    tests = [
        '"hi"',
        "(+ 1.0 2.0)",
        "(* (group (+ 4.0 5.0)) (group (+ 6.0 7.0)))",
        "(* (group (- 5.0 3.0)) 2.0)",
        "(- 2.0 (group (+ 3.0 (* 4.0 5.0))))",
        "(+ (+ 3.0 (* 4.0 5.0)) 2.0)",
    ]
    with open("./tests/test_parser.boa", "r") as source:
        for test, line in zip(tests, source.readlines()):
            lexer  = Lexer(line)
            tokens = lexer.scan_tokens()
            parser = Parser(tokens)
            statements = parser.parse()
            for statement in statements:
                expression = printer.print(statement.expression)
                assert test == expression

if __name__ == "__main__":
    test()

