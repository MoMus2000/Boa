class Environment:
    def __init__(self):
        self.map = {}

    def define(self, name, value):
        self.map[name] = value

    def get(self, name):
        return self.map.get(name)

