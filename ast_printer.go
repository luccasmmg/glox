package main
//
//import (
//	"fmt"
//  "strings"
//)
//
//type Printer struct {}
//
//func (p *Printer) visitBinaryExpr(expr ExprBinary) (interface{}, error) {
//	return fmt.Sprintf("%v", p.parenthize(expr.Operator.Lexeme, expr.Left, expr.Right)), nil
//}
//
//func (p *Printer) visitGroupingExpr(expr ExprGrouping) (interface{}, error) {
//	return fmt.Sprintf("%v", p.parenthize("group", expr.Expression)), nil
//}
//
//func (p *Printer) visitLiteralExpr(expr ExprLiteral) (interface{}, error) {
//  return fmt.Sprintf("%v", expr.Value), nil
//}
//
//func (p *Printer) visitUnaryExpr(expr ExprUnary) (interface{}, error) {
//  return fmt.Sprintf("%v", p.parenthize(expr.Operator.Lexeme, expr.Right)), nil
//}
//
//func (p *Printer) visitVariableExpr(expr ExprVariable) (interface{}, error) {
//  return fmt.Sprintf("%v", expr.Name.Lexeme), nil
//}
//
//func (p *Printer) visitAssignExpr(expr ExprAssign) (interface{}, error) {
//  return fmt.Sprintf("%v", p.parenthize("= " + expr.Name.Lexeme, expr.Value)), nil
//}
//
//func (p *Printer) visitLogicalExpr(expr ExprLogical) (interface{}, error) {
//  return fmt.Sprintf("%v", p.parenthize(expr.Operator.Lexeme, expr.Left, expr.Right)), nil
//}
//
//func (p *Printer) visitCallExpr(expr ExprCall) (interface{}, error) {
//  return fmt.Sprintf("%v", "f"), nil
//}
//
//func (p *Printer) print(expr Expr) (string, error) {
//  value, err := expr.accept(p)
//  return value.(string), err
//}
//
//func (p *Printer) parenthize(name string, exprs ...Expr) string {
//	var builder strings.Builder
//  builder.WriteString("(")
//  builder.WriteString(name)
//  for _, expr := range exprs {
//    builder.WriteString(" ")
//    var str, error = p.print(expr)
//    if error != nil {
//      panic(error)
//    }
//    builder.WriteString(str)
//  }
//  builder.WriteString(")")
//  return builder.String()
//}
