import sys
from lexer import Lexer
from parser import Parser
from expr import AstPrinter
from interpreter import Interpreter

class Boa:
    def __init__(self):
        self.had_error = False
        self.printer = AstPrinter()
        self.interpreter = Interpreter()

    def main(self):
        args = sys.argv
        if len(args) > 2:
            print("Usage: Boa [script]")
        elif len(args) == 2:
            self.run_file(args[1])
        else:
            self.run_prompt()

    def run_file(self, path):
        with open(path, "r") as source_file:
            print(f"Loaded: {path}")
            source_code = source_file.readlines()
            self.run(source_code)
            if self.had_error:
                sys.exit(65)

    def run_prompt(self):
        while True:
            try:
                command = input("> ")
                if command == "":
                    break
                self.run([command])
                self.had_error = False
            except Exception as e:
                print(e)

    def run(self, source_code):
        scanner = Lexer(source_code)
        parser = Parser(scanner.scan_tokens())
        statements = parser.parse()
        self.interpreter.interpret(statements)

    def error(self, line, message):
        self.report(line, " ", message)

    def report(self, line, where, message):
        self.had_error = True
        error = f"""
        [line {line} ] Error {where}: {message}
        """.strip()
        print(error)

if __name__ == "__main__":
    boa = Boa()
    boa.main()
