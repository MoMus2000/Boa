from ..lexer import Lexer
from ..parser import Parser
from interpreter import Interpreter
import io
import sys

def test():
    lexer = None
    interpreter = Interpreter()
    tests = [
    ]
    with open("./tests/test_if_statement.boa", "r") as source:
        old_stdout = sys.stdout
        sys.stdout = buffer = io.StringIO()
        source_code = source.readlines()
        lexer  = Lexer(source_code)
        tokens = lexer.scan_tokens()
        parser = Parser(tokens)
        statements = parser.parse()
        result = interpreter.interpret(statements)
        output = buffer.getvalue()
        print(output)

if __name__ == "__main__":
    test()

