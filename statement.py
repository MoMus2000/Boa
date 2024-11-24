from abc import ABC, abstractmethod

class Statement(ABC):
    def __init__(self):
        pass

    @abstractmethod
    def accept(self, visitor):
        pass

class StmtVisitor(ABC):
    @abstractmethod
    def visit_print_statement(self, stmt) -> object:
        pass

    @abstractmethod
    def visit_expression_statement(self, stmt) -> object:
        pass

class Expression(Statement):
    def __init__(self, expression):
        self.expression = expression

    def accept(self, visitor):
        return visitor.visit_expression_statement(self)

class Print(Statement):
    def __init__(self, expression):
        self.expression = expression

    def accept(self, visitor):
        return visitor.visit_print_statement(self)

