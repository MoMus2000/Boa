class Environment:
    def __init__(self, enclosing = None):
        self.map = {}
        self.enclosing = enclosing

    def define(self, name, value):
        self.map[name] = value

    def get(self, name):
        assert type(name) == str, "Not string"
        if name in self.map:
            return self.map.get(name)
        if self.enclosing != None:
            return self.enclosing.get(name)

    def assign(self, name, value):
        if name.lexeme in self.map:
            self.map[name.lexeme] = value
            return
        if self.enclosing != None:
            self.enclosing.assign(name, value)
            return
        if name.lexeme not in self.map:
            raise Exception(f'Error: undefined var "{name.lexeme}"')
