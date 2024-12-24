package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Glox struct {
	hadError        bool
	hadRuntimeError bool
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

func (g Glox) reportError(line int, message string) {
	g.report(line, "", message)
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
  var environment = NewEnvironment()
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
		g.run(line, &environment)
	}
}

func (g Glox) run(source string, env ...*Environment) {
	scanner := NewScanner(source)
	tokens := scanner.scanTokens()
	parser := NewParser(tokens)
	statements, errors := parser.parse()
	if g.hadError {
		fmt.Println("Error parsing expression")
		return
	}
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
		//os.Exit(65)
	} else {
    var environment Environment;
    if len(env) > 0 {
      environment = *env[0]
    } else {
      environment = NewEnvironment()
    }
		interpreter := Interpreter{
      environment: &environment,
    }
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
