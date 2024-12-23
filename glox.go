package main

import (
	"fmt"
	"io"
	"os"
	"strings"
  "bufio"
)

type Glox struct {
  hadError bool
  hadRuntimeError bool
}

func (g Glox) runFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (g Glox) reportError(line int, message string) {
    g.report(line, "", message);
}

func (g Glox) report(line int, where string, message string) {
    fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
}

func (g Glox) runtimeError(err RuntimeError) {
  fmt.Printf("%s\n[line %d]\n", err.message, err.token.Line)
  g.hadRuntimeError = true
  os.Exit(70)
}

func (g Glox) runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}
    line = strings.TrimSpace(line)
		if line == "" {
			break
		}
    g.run(line)
	}
}

func (g Glox) run(source string) {
  scanner := NewScanner(source)
  tokens := scanner.scanTokens()
  parser := NewParser(tokens)
  expression := parser.Parse()
  if g.hadError {
    fmt.Println("Error parsing expression")
    return
  }
  interpreter := Interpreter{}
  value, err := interpreter.interpret(expression)
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(value)
}

func main() {
	argCount := len(os.Args)
  g := Glox{}
	if argCount > 2 {
		fmt.Println("Usage: glox [script]")
	} else if argCount == 2 {
		g.runFile(os.Args[1])
	} else {
		g.runPrompt()
	}
}

