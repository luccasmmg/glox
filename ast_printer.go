package main

import (
	"fmt"
  "strings"
)

type Printer struct {}

func (p *Printer) visitBinaryExpr(expr ExprBinary) interface{} {
	return fmt.Sprintf("%v", p.parenthize(expr.Operator.Lexeme, expr.Left, expr.Right))
}

func (p *Printer) visitGroupingExpr(expr ExprGrouping) interface{} {
	return fmt.Sprintf("%v", p.parenthize("group", expr.Expression))
}

func (p *Printer) visitLiteralExpr(expr ExprLiteral) interface{} {
  return fmt.Sprintf("%v", expr.Value)
}

func (p *Printer) visitUnaryExpr(expr ExprUnary) interface{} {
  return fmt.Sprintf("%v", p.parenthize(expr.Operator.Lexeme, expr.Right))
}

func (p *Printer) visitVariableExpr(expr ExprVariable) interface{} {
  return fmt.Sprintf("%v", expr.Name.Lexeme)
}

func (p *Printer) visitAssignExpr(expr ExprAssign) interface{} {
  return fmt.Sprintf("%v", p.parenthize("= " + expr.Name.Lexeme, expr.Value))
}

func (p *Printer) visitLogicalExpr(expr ExprLogical) interface{} {
  return fmt.Sprintf("%v", p.parenthize(expr.Operator.Lexeme, expr.Left, expr.Right))
}

func (p *Printer) print(expr Expr) string {
  return expr.accept(p).(string)
}

func (p *Printer) parenthize(name string, exprs ...Expr) string {
	var builder strings.Builder
  builder.WriteString("(")
  builder.WriteString(name)
  for _, expr := range exprs {
    builder.WriteString(" ")
    builder.WriteString(p.print(expr))
  }
  builder.WriteString(")")
  return builder.String()
}
