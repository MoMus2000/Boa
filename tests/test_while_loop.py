from ..lexer import Lexer
from ..parser import Parser
from interpreter import Interpreter

def test():
    lexer = None
    interpreter = Interpreter(debug_mode=True)
    with open("./tests/test_while_loop.boa", "r") as source:
        source_code = source.readlines()
        lexer       = Lexer(source_code)
        tokens      = lexer.scan_tokens()
        parser      = Parser(tokens)
        statements  = parser.parse()
        result = interpreter.interpret(statements)

if __name__ == "__main__":
    test()

