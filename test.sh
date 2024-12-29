#!/bin/bash

PYTHONPATH=$(pwd) pytest tests -s -v -q

# Find all .boa files in the tests directory
files=$(find ./tests -type f -name "*.boa" ! -name "*lexer*" ! -name "*parser*" ! -name "*interpreter*")


# Count the number of files
count=$(echo "$files" | wc -l)

# Print and process the files
echo "Total .boa files found: $count"
for file in $files; do
  python3 main.py $file > /dev/null
  echo -e "\033[32mSuccess: Processed $file\033[0m"
done
