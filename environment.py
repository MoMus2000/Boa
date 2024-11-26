class Environment:
    def __init__(self):
        self.map = {}

    def define(self, name, value):
        self.map[name] = value

    def get(self, name):
        return self.map.get(name)

    def assign(self, name, value):
        if name.lexeme not in self.map:
            raise Exception(f'Error: undefined var "{name.lexeme}"')
        self.map[name.lexeme] = value

