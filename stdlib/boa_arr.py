class Arr:
    def __init__(self):
        pass

    def append(self, a, val):
        if not isinstance(a, list):
            raise NotImplemented("Append only works for list")
        a.append(val)

    def length(self, a):
        if not isinstance(a, list):
            raise NotImplemented("Append only works for list")
        return len(a)
    
    def pop(self, a):
        if not isinstance(a, list):
            raise NotImplemented("Append only works for list")
        a.pop()
        return self.length(a)

    def modify(self, a, val, idx):
        if not isinstance(a, list):
            raise NotImplemented("Append only works for list")

        a[int(idx)] = val
        
