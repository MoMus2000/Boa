class Map:
    def __init__(self):
        pass
    
    def insert(self, map_itself, key, value):
        map_itself[key] = value

    def get(self,map_itself, key):
        return map_itself.get(key)
    
    def keys(self, map_itself):
        return list(map_itself.keys())

    def values(self, map_itself):
        return list(map_itself.values())

