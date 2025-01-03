# Boa

Boa is an interpreter written in Python. It is inspired by Bob Nystrom's implementation of the Lox programming language. The project was originally intended to be written in Rust, but I found myself struggling with the borrow checker rather than enjoying the process of writing code, so I switched to Python for the initial implementation.

# Warning

**Experimental Language:** This is an experimental language that is currently not optimized for performance nor considered stable. In benchmark tests, a simple for loop with computation is approximately **2000% slower** than Python. Please keep this in mind if you plan to use this for your applications.


**Note:** Future optimizations may improve performance, performance should
not be considered on par with more mature languages like Python for simple computational
tasks. 

## Setup


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

## Sample Boa Code

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

## Standard lib

### Arrays
```lua
import arr

var nums = [1, 2, 3];

for(var i=0; i<arr.length(nums); i = i + 1;) {
    print nums[i];
}

print nums[0];

nums[0] = 4;

var nums2 = [4, 5, 6];

var nums3 = nums + nums2;

var a = [[1, 2], [3, 4]];
var b = [[5, 6], [7, 8]];

fun multiply_matrices(a, b) {
  var m = arr.length(a);
  var n = arr.length(a[0]);
  var p = arr.length(b[0]);

  if(n != arr.length(b)) {
    return nil;
  }

  var c = [];

  for(var i =0; i< m; i = i + 1;) {
    var row = [];
    for(var j=0; j <p; j = j +1;) {
      arr.append(row, 0);
    }
    arr.append(c, row);
  }

  for(var i =0; i< m; i = i+1;){
    for(var j =0; j < p; j = j +1;) {
      for(var k =0; k < n; k = k +1;) {
          c[i][j] = c[i][j] + a[i][k]*b[k][j];
      }
    }
  }

  return c;
}

var mul_ab = multiply_matrices(a, b);
var result = [[19.0, 22.0], [43.0, 50.0]];
assert(mul_ab == result, "Output does not match");

```

### Math
```lua
import math

print math.pow(2, 3);
print math.factorial(10);
print math.ceil(6.9);
print math.floor(6.9);

```
### Time
```lua
import time

print "Sleeping for 5 secs";
time.sleep(5);

print "Time Stamp " + time.clock();

```
