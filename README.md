# Boa

Creating a language.

```bash
# Load source code
python3 main.py examples/demo.boa

# Launch interpreter
python3 main.py

# Run tests
./tests.sh
```

Sample Boa Code

```lua

fun print_hello(){
  print "Hello";
}

fun add_time_together(){
  var start = clock();
  var end   = clock();

  print "Starting at " + start;
  print "Ending   at " + end;

}

print_hello();
print_hello();

add_time_together();

```

```lua
var hello = "Hello";
var world = "World";

print hello + world;

{
    var boa = 1;
    boa = boa + 1;
    print boa;
}

if ( 1 == 1 ) {
    print "One is equal to One";
} else {
    print "Not expected to get to this point";
}

var a = 1;
while ( a <= 5 ) {
    a = a + 1;
}

print a;

// Fibonacci Numbers via for loops:
var a = 0;
var temp;

for (var b = 1; a < 10000; b = temp + b;) {
  print a;
  temp = a;
  a = b;
}

```

