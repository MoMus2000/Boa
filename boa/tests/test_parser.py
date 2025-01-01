from ..lexer import Lexer
from ..parser import Parser
from ..expr import *

def test():
    lexer = None
    printer = AstPrinter()
    tests = [
        '"hi"',
        "(+ 1 2)",
        "(* (group (+ 4 5)) (group (+ 6 7)))",
        "(* (group (- 5 3)) 2)",
        "(- 2 (group (+ 3 (* 4 5))))",
        "(+ (+ 3 (* 4 5)) 2)",
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

