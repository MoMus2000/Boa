# Boa

Boa is an interpreter written in Python. It is inspired by Bob Nystrom's implementation of the Lox programming language. The project was originally intended to be written in Rust, but I found myself struggling with the borrow checker rather than enjoying the process of writing code, so I switched to Python for the initial implementation.

# Warning

**Experimental Language:** This is an experimental language that is currently not optimized for performance. In benchmark tests, a simple for loop with computation is approximately **2000% slower** than Python. Please keep this in mind if you plan to use this for performance-critical applications.


**Note:** Future optimizations may improve performance, performance should
not be considered on par with more mature languages like Python for simple computational
tasks. 

# Setup


```bash
# Run setup
./setup.sh (make sure to provide relevant permissions)

# Load source code
boa path_to_file

# Launch interpreter
boa

# Run tests
./tests.sh
```

# Sample Boa Code

```lua
// Fibonacci Numbers via recursion:
fun fib(n) {
  if (n <= 1){
    return n;
  }
  return fib(n - 2) + fib(n - 1);
}

for (var i = 0; i < 20; i = i + 1;) {
  print fib(i);
}

// Fibonacci Numbers via for loops:
var a = 0;
var temp;

for (var b = 1; a < 10000; b = temp + b;) {
  print a;
  temp = a;
  a = b;
}

fun add(a, b){
  return a + b;
}

fun sub(a, b){
  return a - b;
}

print add(1, 7);
print add(sub(add(2,7), 3), 4);

```

Standard lib (Work in Progress)
```lua
import math
import time

print math.pow(2, 3);
print math.factorial(10);
print math.ceil(6.9);
print math.floor(6.9);

time.sleep(1);

print time.clock();

```

