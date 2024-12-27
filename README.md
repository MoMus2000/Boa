# Boa

Creating a language.

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

Sample Boa Code

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

