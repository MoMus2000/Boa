#!/bin/bash
set -e

# Find all .boa files in the tests directory
files=$(find ./tests -type f -name "*.boa" ! -name "*lexer*" ! -name "*parser*" ! -name "*interpreter*")

# Count the number of files
count=$(echo "$files" | wc -l)

# Print and process the files
echo "Total .boa files found: $count"
for file in $files; do
  go run . $file > /dev/null
  echo -e "\033[32mSuccess: Processed $file\033[0m"
done
