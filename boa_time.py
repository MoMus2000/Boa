import time

class Time:
    def __init__(self):
        pass

    def sleep(self, s):
        s = int(s)
        time.sleep(s)

    def clock(self):
        from datetime import datetime
        return f"{datetime.now()}"

