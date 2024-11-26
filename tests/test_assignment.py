from ..lexer import Lexer
from ..parser import Parser
from interpreter import Interpreter

def test():
    lexer = None
    tests = [
        '5',
        '6',
        '"mustafa"',
        '6',
        '6'
    ]
    interpreter = Interpreter()
    with open("./tests/test_assignment.boa", "r") as source:
        for test, line in zip(tests, source.readlines()):
            lexer  = Lexer(line)
            tokens = lexer.scan_tokens()
            parser = Parser(tokens)
            statements = parser.parse()
            result = interpreter.interpret(statements)        
            print(result[0], test)
            assert result[0] == test

if __name__ == "__main__":
    test()
