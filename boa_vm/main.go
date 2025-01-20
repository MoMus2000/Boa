package main

import (
  "bufio"
  "fmt"
  "os"
)

func repl(vm VM) {
  reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(fmt.Errorf("Error Reading Stdin"))
		}
		text_bytes := []byte(text)
    vm.interpret(text_bytes)
	}
}

func runFile(vm VM, filePath string) {
  fmt.Printf("Loading source code from file %v\n", filePath)
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(fmt.Errorf("Could not find the source file %v\n", err))
	}
  result := vm.interpret(file)
  if result == INTERPRET_RUNTIME_ERROR || result == INTERPRET_COMPILE_ERROR {
    os.Exit(69)
  }
}

func main(){
  vm := NewVM();
  args := os.Args
  if len(args) == 1 {
    repl(vm)
  } else if len(args) == 2 {
    filePath := args[1]
    runFile(vm, filePath)
  } else {
    fmt.Println("Usage: boa [path]")
    os.Exit(69)
  }
  vm.FreeVM()
}

