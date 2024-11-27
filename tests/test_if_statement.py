from ..lexer import Lexer
from ..parser import Parser
from interpreter import Interpreter

def test():
    lexer = None
    interpreter = Interpreter()
    tests = [
    ]
    with open("./tests/test_if_statement.boa", "r") as source:
            source_code = source.readlines()
            lexer  = Lexer(source_code)
            tokens = lexer.scan_tokens()
            parser = Parser(tokens)
            statements = parser.parse()
            result = interpreter.interpret(statements)
            print(result)
            # for test, r in zip(tests, result):
            #     assert r == test

if __name__ == "__main__":
    test()

