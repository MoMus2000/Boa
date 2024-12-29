import time
import random

start = time.time()


for i in range(0, 2000000):
    rand = random.random()

end = time.time()

print(start, end)

