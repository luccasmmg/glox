package main

import (
	"fmt"
	"strconv"
)

type Interpreter struct {
	environment *Environment
	globals     *Environment
}

func NewInterpreter() *Interpreter {
	// global env
	global := NewEnvironment(nil)
	global.define("clock", Time{})

	return &Interpreter{
		globals:     global,
		environment: global,
	}
}

type RuntimeError struct {
	token   Token
	message string
}

type Return struct {
	Value interface{}
}

func (e Return) Error() string {
	return ""
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("Error at '%s': %s", e.token.Lexeme, e.message)
}

func (i *Interpreter) visitLiteralExpr(expr ExprLiteral) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) visitLogicalExpr(expr ExprLogical) (interface{}, error) {
	var left, err = i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	if expr.Operator.TokenType == OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) visitGroupingExpr(expr ExprGrouping) (interface{}, error) {
	var value, err = i.evaluate(expr.Expression)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) visitUnaryExpr(expr ExprUnary) (interface{}, error) {
	var right, err = i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.TokenType {
	case MINUS:
		if err := checkNumberOperand(expr.Operator, right); err != nil {
			return nil, err
		}
		return right.(float64) * (-1), nil
	case BANG:
		return !i.isTruthy(right), nil
	}

	return RuntimeError{token: expr.Operator, message: "Unknown operator."}, nil
}

func (i *Interpreter) visitVariableExpr(expr ExprVariable) (interface{}, error) {
	var value, err = i.environment.get(expr.Name)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}
	switch v := obj.(type) {
	case bool:
		return v
	default:
		return true
	}
}

func (i *Interpreter) isEqual(obj_a interface{}, obj_b interface{}) bool {
	if obj_a == nil && obj_b == nil {
		return true
	}
	if obj_a == nil {
		return false
	}
	return obj_a == obj_b
}

func checkNumberOperand(operator Token, operand interface{}) error {
	_, ok := operand.(float64)
	if !ok {
		if str, ok := operand.(string); ok {
			if _, err := strconv.ParseFloat(str, 64); err == nil {
				ok = true
			}
		}
	}
	if !ok {
		return &RuntimeError{token: operator, message: "Operand must be a number."}
	}
	return nil
}

func checkIfNumberIsZero(number float64, operator Token, operand interface{}) error {
	if number == 0 {
		return &RuntimeError{token: operator, message: "Division by zero."}
	}
	return nil
}

func checkNumberOperands(operator Token, left interface{}, right interface{}) error {
	_, leftOk := left.(float64)
	if !leftOk {
		if leftStr, ok := left.(string); ok {
			if _, err := strconv.ParseFloat(leftStr, 64); err == nil {
				leftOk = true
			}
		}
	}

	_, rightOk := right.(float64)
	if !rightOk {
		if rightStr, ok := right.(string); ok {
			if _, err := strconv.ParseFloat(rightStr, 64); err == nil {
				rightOk = true
			}
		}
	}

	if !leftOk {
		return &RuntimeError{token: operator, message: fmt.Sprintf("Left operand must be a number, but got %T.", left)}
	}
	if !rightOk {
		return &RuntimeError{token: operator, message: fmt.Sprintf("Right operand must be a number, but got %T.", right)}
	}

	// Use leftValue and rightValue as the coerced numbers
	return nil
}

func (i *Interpreter) visitBinaryExpr(expr ExprBinary) (interface{}, error) {
	var left, errLeft = i.evaluate(expr.Left)
	if errLeft != nil {
		return nil, errLeft
	}
	var right, errRight = i.evaluate(expr.Right)
	if errRight != nil {
		return nil, errRight
	}
	switch expr.Operator.TokenType {
	case MINUS:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case SLASH:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		if err := checkIfNumberIsZero(right.(float64), expr.Operator, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case STAR:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case PLUS:
		if leftStr, ok := left.(string); ok {
			if rightStr, ok := right.(string); ok {
				return leftStr + rightStr, nil
			}
		}
		if leftFloat, ok := left.(float64); ok {
			if rightFloat, ok := right.(float64); ok {
				return leftFloat + rightFloat, nil
			}
		}
	case GREATER:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case LESS:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case BANG_EQUAL:
		return !i.isEqual(left, right), nil
	case EQUAL_EQUAL:
		return i.isEqual(left, right), nil
	}
	return RuntimeError{token: expr.Operator, message: "Unknown operator."}, nil
}

func (i *Interpreter) visitCallExpr(expr ExprCall) (interface{}, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}
	function, ok := callee.(GloxCallable)
	if !ok {
		return nil, &RuntimeError{
			token:   expr.Paren,
			message: "Can only call functions and classes",
		}
	}
	var arguments []interface{}
	for _, argument := range expr.Arguments {
		_arg, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, _arg)
	}
	if len(arguments) != function.Arity() {
		message := fmt.Sprintf("Expected %d arguments but got %d", function.Arity(), len(arguments))
		return nil, &RuntimeError{
			token:   expr.Paren,
			message: message,
		}
	}
	return function.Call(i, arguments)
}

func (i *Interpreter) visitStmtExpression(stmt StmtExpression) error {
	var _, err = i.evaluate(stmt.Expression)
	if err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) visitStmtFunction(stmt StmtFunction) error {
	var function GloxFunction = GloxFunction{Declaration: stmt, Closure: i.environment}
	i.environment.define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) visitStmtWhile(stmt StmtWhile) error {
	for {
		condition, err := i.evaluate(stmt.Condition)
		if err != nil {
			return err
		}
		if !i.isTruthy(condition) {
			break
		}
		i.execute(stmt.Body)
	}
	return nil
}

func (i *Interpreter) visitStmtIf(stmt StmtIf) error {
	var value, err = i.evaluate(stmt.Condition)
	if err != nil {
		return err
	}
	if i.isTruthy(value) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) visitStmtPrint(stmt StmtPrint) error {
	var value, err = i.evaluate(stmt.Expression)
	if err != nil {
		return err
	}
	fmt.Println(value)
	return nil
}

func (i *Interpreter) visitStmtReturn(stmt StmtReturn) error {
	value, err := i.evaluate(stmt.Value)
	if err != nil {
		return err
	}

	return Return{Value: value}
}

func (i *Interpreter) visitStmtVarDeclaration(stmt StmtVarDeclaration) error {
	var value interface{} = nil
	if stmt.Initializer != nil {
		var val, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return err
		}
		value = val
	}
	i.environment.define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) visitAssignExpr(expr ExprAssign) (interface{}, error) {
	var value, error = i.evaluate(expr.Value)
	if error != nil {
		return nil, error
	}
	i.environment.assign(expr.Name, value)
	return value, nil
}

func (i *Interpreter) visitStmtAssign(stmt StmtAssign) error {
	var value, err = i.evaluate(stmt.Value)
	if err != nil {
		return err
	}
	i.environment.assign(stmt.Name, value)
	return nil
}

func (i *Interpreter) execute(stmt Stmt) error {
	return stmt.accept(i)
}

func (i *Interpreter) visitStmtBlock(stmt StmtBlock) error {
	return i.executeBlock(stmt.Statements, NewEnvironment(i.environment))
}

func (intr *Interpreter) executeBlock(statements []Stmt, environment *Environment) error {
	previous := intr.environment
	intr.environment = environment

	for _, stmt := range statements {
		if err := intr.execute(stmt); err != nil {
			intr.environment = previous
			return err
		}
	}

	intr.environment = previous
	return nil
}

func (i *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.accept(i)
}

func (i *Interpreter) interpret(expr []Stmt) (interface{}, error) {
	for _, stmt := range expr {
		var err = i.execute(stmt)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
