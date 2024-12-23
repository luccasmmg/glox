package main

import (
  "fmt"
  "strconv"
)

type Interpreter struct{}

type RuntimeError struct {
  token Token
  message string
}

func (e *RuntimeError) Error() string {
  return fmt.Sprintf("Error at '%s': %s", e.token.Lexeme, e.message)
}

func (i *Interpreter) visitLiteralExpr(expr ExprLiteral) (interface{}, error) {
	return expr.Value, nil
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
    fmt.Println(left, right)
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

func (i *Interpreter) evaluate(expr Expr) (interface{}, error) {
  value, err := expr.accept(i)
  return value, err
}

func (i *Interpreter) interpret(expr Expr) (interface{}, error) {
  return i.evaluate(expr)
}
