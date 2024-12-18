package main 

type Expr interface {
  accept(v visitor) interface{}
}

type visitor interface {
  visitBinaryExpr(expr ExprBinary) interface{}
  visitGroupingExpr(expr ExprGrouping) interface{}
  visitLiteralExpr(expr ExprLiteral) interface{}
  visitUnaryExpr(expr ExprUnary) interface{}
  visitVariableExpr(expr ExprVariable) interface{}
  visitAssignExpr(expr ExprAssign) interface{}
  visitLogicalExpr(expr ExprLogical) interface{}
  //visitCallExpr(expr ExprCall) interface{}
}

type ExprCall struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

type ExprBinary struct {
	Operator Token
	Left     Expr
	Right    Expr
}

type ExprLogical struct {
	Operator Token
	Left     Expr
	Right    Expr
}

type ExprGrouping struct {
	Expression Expr
}

type ExprAssign struct {
	Name  Token
	Value Expr
}

type ExprLiteral struct {
	Value interface{}
}

type ExprVariable struct {
	Name Token
}

type ExprUnary struct {
	Operator Token
	Right    Expr
}

func (e ExprBinary) accept(v visitor) interface{} {
  return v.visitBinaryExpr(e)
}

func (e ExprGrouping) accept(v visitor) interface{} {
  return v.visitGroupingExpr(e)
}

func (e ExprLiteral) accept(v visitor) interface{} {
  return v.visitLiteralExpr(e)
}

func (e ExprUnary) accept(v visitor) interface{} {
  return v.visitUnaryExpr(e)
}

func (e ExprVariable) accept(v visitor) interface{} {
  return v.visitVariableExpr(e)
}

func (e ExprAssign) accept(v visitor) interface{} {
  return v.visitAssignExpr(e)
}

func (e ExprLogical) accept(v visitor) interface{} {
  return v.visitLogicalExpr(e)
}

//func (e ExprCall) accept(v visitor) interface{} {
//  return v.visitCallExpr(e)
//}

