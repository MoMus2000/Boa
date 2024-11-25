from ..lexer import Lexer
from ..parser import Parser
from interpreter import Interpreter

def test():
    lexer = None
    tests = [
        '5',
        '5',
        '6969',
        '"mustafa"',
        '"mustafa"',
        '6969'
    ]
    interpreter = Interpreter()
    with open("./tests/test_var_declaration.boa", "r") as source:
        i = 0
        for test, line in zip(tests, source.readlines()):
            lexer  = Lexer(line)
            tokens = lexer.scan_tokens()
            parser = Parser(tokens)
            statements = parser.parse()
            result = interpreter.interpret(statements)        
            assert result[0] == test
            i+=1

if __name__ == "__main__":
    test()

