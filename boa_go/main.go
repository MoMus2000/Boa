package main

import (
	"bufio"
	"fmt"
	"os"
)

type Boa struct {
}

func (boa *Boa) run_file(file_path string) {
	fmt.Printf("Loading source code from file %v\n", file_path)
	file, err := os.ReadFile(file_path)
	if err != nil {
		fmt.Println(fmt.Errorf("Could not find the source file %v\n", err))
	}
	boa.run(file)
}

func (boa *Boa) run_prompt() {
  reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Printf("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(fmt.Errorf("Error Reading Stdin"))
		}
		fmt.Printf("%v", text)
		text_bytes := []byte(text)
		boa.run(text_bytes)
	}
}

func (boa *Boa) run(source_code []byte) {
	lexer := NewLexer(source_code)
	fmt.Println(lexer.ScanTokens())
}

func main() {
	args := os.Args
	boa := Boa{}
	if len(args) > 2 {
		fmt.Println(fmt.Errorf("Usage: Boa [script]"))
		return
	} else if len(args) == 2 {
		file_path := args[1]
		boa.run_file(file_path)
	} else {
		boa.run_prompt()
	}
}
