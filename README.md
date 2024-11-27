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
```

