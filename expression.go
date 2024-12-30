package main 

type Expr interface {
  accept(v ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
  visitBinaryExpr(expr ExprBinary) (interface{}, error)
  visitGroupingExpr(expr ExprGrouping) (interface{}, error)
  visitLiteralExpr(expr ExprLiteral) (interface{}, error)
  visitUnaryExpr(expr ExprUnary) (interface{}, error)
  visitVariableExpr(expr ExprVariable) (interface{}, error)
  visitLogicalExpr(expr ExprLogical) (interface{}, error)
  visitAssignExpr(expr ExprAssign) (interface{}, error)
  visitCallExpr(expr ExprCall) (interface{}, error)
  visitGetExpr(expr ExprGet) (interface{}, error)
  visitSetExpr(expr ExprSet) (interface{}, error)
  visitThisExpr(expr ExprThis) (interface{}, error)
  visitSuperExpr(expr ExprSuper) (interface{}, error)
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

type ExprGet struct {
  Object Expr
  Name Token
}

type ExprSet struct {
  Object Expr
  Name Token
  Value Expr
}

type ExprThis struct {
  Keyword Token
}

type ExprSuper struct {
  Keyword Token
  Method Token
}

func (e ExprBinary) accept(v ExprVisitor) (interface{}, error) {
  value, err := v.visitBinaryExpr(e)
  return value, err
}

func (e ExprGrouping) accept(v ExprVisitor)  (interface{}, error) {
  value, err := v.visitGroupingExpr(e)
  return value, err
}

func (e ExprLiteral) accept(v ExprVisitor)  (interface{}, error) {
  value, err := v.visitLiteralExpr(e)
  return value, err
}

func (e ExprUnary) accept(v ExprVisitor)  (interface{}, error) {
  value, err := v.visitUnaryExpr(e)
  return value, err
}

func (e ExprVariable) accept(v ExprVisitor)  (interface{}, error) {
  return v.visitVariableExpr(e)
}

func (e ExprAssign) accept(v ExprVisitor)  (interface{}, error) {
  return v.visitAssignExpr(e)
}

func (e ExprLogical) accept(v ExprVisitor)  (interface{}, error) {
  return v.visitLogicalExpr(e)
}

func (e ExprCall) accept(v ExprVisitor) (interface{}, error) {
  return v.visitCallExpr(e)
}

func (e ExprGet) accept(v ExprVisitor) (interface{}, error) {
  return v.visitGetExpr(e)
}

func (e ExprSet) accept(v ExprVisitor) (interface{}, error) {
  return v.visitSetExpr(e)
}

func (e ExprThis) accept(v ExprVisitor) (interface{}, error) {
  return v.visitThisExpr(e)
}

func (e ExprSuper) accept(v ExprVisitor) (interface{}, error) {
  return v.visitSuperExpr(e)
}
