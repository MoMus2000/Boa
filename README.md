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
fun add_time_together(){
  var start = clock();
  var end   = clock();

  print "Starting at " + start;
  print "Ending   at " + end;

}

add_time_together();

fun add(a, b){
  return a + b;
}

fun sub(a, b){
  return a - b;
}

print add(1, 7);
print add(sub(add(2,7), 3), 4);
print sub(add(20, sub(15, 5)), add(10, 5));
print add(sub(add(5, 10), sub(8, 3)), sub(20, add(5, 5)));
print sub(sub(add(8, 4), add(2, 3)), add(1, sub(5, 3)));
print add(sub(add(3, sub(9, 6)), 2), sub(10, add(4, 2)));

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

