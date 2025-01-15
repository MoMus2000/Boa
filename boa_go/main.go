package main

import (
	"bufio"
	"fmt"
	"os"
)

// boa/boa_go [main] $ go run -ldflags="-X 'main.Ip=localhost:8000'" .
var Ip string;

type Boa struct {
  interpreter Interpreter
}

func (boa *Boa) run_file(file_path string, telem *Client) {
	fmt.Printf("Loading source code from file %v\n", file_path)
	file, err := os.ReadFile(file_path)
	if err != nil {
		fmt.Println(fmt.Errorf("Could not find the source file %v\n", err))
	}
  if telem != nil {
    telem.write_to_server(string(file))
  }
	boa.run(file)
}

func (boa *Boa) run_prompt(telem *Client) {
  reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Printf("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(fmt.Errorf("Error Reading Stdin"))
		}
		text_bytes := []byte(text)
    if telem != nil {
      telem.write_to_server(text)
    }
		boa.run(text_bytes)
	}
}

func (boa *Boa) run(source_code []byte) error {
  parser      := NewParser(source_code)
  statements, err  := parser.parse()
  if err != nil {
    fmt.Println(err)
    return err
  }
  err = boa.interpreter.interpret(statements)
  if err != nil {
    fmt.Println(err)
    return err
  }
  return nil
}

func main() {
  fmt.Println("Boa V0.1 - by momus2000")
  telem := init_telem(Ip)
	args := os.Args
  interpreter := NewInterpreter()
	boa := Boa{
    interpreter: *interpreter,
  }
	if len(args) > 2 {
		fmt.Println(fmt.Errorf("Usage: Boa [script]"))
		return
	} else if len(args) == 2 {
		file_path := args[1]
		boa.run_file(file_path, telem)
	} else {
		boa.run_prompt(telem)
	}
}

