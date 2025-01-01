import sys
from lexer import Lexer
from parser import Parser
from expr import AstPrinter
from interpreter import Interpreter

class UnSupportedOsError(Exception):
    def __init__(self, message):
        self.message = message

def determine_os():
    import os
    os_name = os.name
    if os_name == "posix":
        print("You are on a POSIX-compliant system (e.g., Linux, macOS).")
        return "posix"
    elif os_name == 'nt': 
        print("Boa has been developed on posix")
        print("You are on a Windows system, proceed at your own risk.") 
        return "nt"
    else:
        raise UnSupportedOsError("Boa does not work on this OS")

class _Getch:
    """Gets a single character from standard input.  Does not echo to the
screen."""
    def __init__(self):
        try:
            self.impl = _GetchWindows()
        except ImportError:
            self.impl = _GetchUnix()

    def __call__(self): return self.impl()


class _GetchUnix:
    def __init__(self):
        import tty, sys

    def __call__(self):
        import sys, tty, termios
        fd = sys.stdin.fileno()
        old_settings = termios.tcgetattr(fd)
        try:
            tty.setraw(sys.stdin.fileno())
            ch = sys.stdin.read(1)
        finally:
            termios.tcsetattr(fd, termios.TCSADRAIN, old_settings)
        return ch


class _GetchWindows:
    def __init__(self):
        import msvcrt

    def __call__(self):
        import msvcrt
        return msvcrt.getch()

class Boa:
    def __init__(self):
        self.had_error = False
        self.printer = AstPrinter()
        self.interpreter = Interpreter()
        self.os = determine_os()

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

    def is_up(self, ords):
        arr = [27, 91, 65]
        sub_len = len(arr)
        for i in range(len(ords) - sub_len + 1):
            if ords[i:i + sub_len] == arr:
                return True
        return False

    def run_prompt(self):
        history = []
        counter = 0
        while True:
            collected = []
            print("> ", end="", flush=True)
            while True:
                press = False
                ch = _Getch()()
                if ord(ch) == 13:
                    print()
                    break
                elif ord(ch) == 3 or ord(ch) == 4:
                    raise KeyboardInterrupt
                elif ord(ch) == 27:
                    ch = _Getch()()
                    if ord(ch) == 91:
                        ch = _Getch()()
                        if ord(ch) == 65:
                            press = True
                            if abs(counter) < len(history) and len(history) != 0:
                                counter -= 1
                                collected = [history[counter]]
                        if ord(ch) == 66:
                            press = True
                            if abs(counter) < len(history) and len(history) != 0:
                                counter += 1
                                collected = [history[counter]]
                elif ord(ch) == 127:
                    if collected:
                        collected.pop()
                        print("\b \b", end="", flush=True)
                        # Move the cursor back to the position where the last character was
                        print("\b", end="", flush=True)
                    else:
                        print(" " * 80, end="\r", flush=True)  # Clear the line by printing spaces
                        print("> ", end="", flush=True)
                else:
                    collected.append(ch)
                if len(collected):
                    if press == True:
                        print(" " * 80, end="\r", flush=True)  # Clear the line by printing spaces
                        print(f"> {collected[-1]}", end="\r", flush=True)
                    elif press == False:
                        print(f"{collected[-1]}", end="", flush=True)
            try:
                command = "".join(collected)
                history.append(command)
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
