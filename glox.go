package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var hadError        bool
var	hadRuntimeError bool

type Glox struct {
}

func (g Glox) runFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
  var source = string(bytes)
  g.run(source)
  return nil
}

func reportError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
}

func (g Glox) runtimeError(err RuntimeError) {
	fmt.Printf("%s\n[line %d]\n", err.message, err.token.Line)
	hadRuntimeError = true
	os.Exit(70)
}

func parseError(token Token, message string) {
	if token.TokenType == EOF {
		report(token.Line, "at end", message)
	} else {
		report(token.Line, fmt.Sprintf("at '%v'", token.Lexeme), message)
	}
}

func (g Glox) runPrompt() {
  var environment = NewEnvironment(nil)
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
		g.run(line, environment)
	}
}

func (g Glox) run(source string, env ...*Environment) {
	scanner := NewScanner(source)
	tokens := scanner.scanTokens()
	parser := NewParser(tokens)
	statements, errors := parser.parse()
	if hadError {
		fmt.Println("Error parsing expression")
		return
	}
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
		//os.Exit(65)
	} else {
    //var environment *Environment;
    //if len(env) > 0 {
    //  environment = env[0]
    //} else {
    //  environment = NewEnvironment(nil)
    //}
		interpreter := NewInterpreter()
    if len(statements) == 1 {
      var stmt = statements[0]
      if stmt, ok := stmt.(StmtExpression); ok {
        value, err := interpreter.evaluate(stmt.Expression)
        if err != nil {
          fmt.Println(err)
        }
        fmt.Printf("%v\n", value)
        return
      }
    }
		value, err := interpreter.interpret(statements)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(value)
	}
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
