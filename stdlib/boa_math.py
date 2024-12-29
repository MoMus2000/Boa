import math

class Math:
    def __init__(self):
        pass

    def pi(self):
        return math.pi

    def pow(self, a, b):
        return math.pow(float(a), float(b))

    def factorial(self, a):
        return math.factorial(int(a))

    def ceil(self, a):
        return math.ceil(float(a))

    def min(self, a, b):
        a = float(a)
        b = float(b)
        return min(a, b)

    def max(self, a, b):
        a = float(a)
        b = float(b)
        return max(a, b)

    def sqrt(self, n):
        n = float(n)
        return math.sqrt(n)

    def floor(self, a):
        import math
        return math.floor(float(a))

    def mod(self, a, b):
        return int(a) % int(b)

    def abs(self, a):
        return abs(float(a))

    def is_prime(self, n):
        n = int(n)
        """
        Check if a number is a prime number.

        :param n: The number to check.
        :return: True if n is prime, False otherwise.
        """
        if n <= 1:
            return False
        if n <= 3:
            return True  # 2 and 3 are prime numbers
        if n % 2 == 0 or n % 3 == 0:
            return False
        i = 5
        while i * i <= n:
            if n % i == 0 or n % (i + 2) == 0:
                return False
            i += 6
        return True

