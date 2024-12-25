package main

type Stmt interface {
	accept(v StmtVisitor) error
}

type StmtVisitor interface {
	visitStmtPrint(stmt StmtPrint) error
	visitStmtVarDeclaration(expr StmtVarDeclaration) error
	visitStmtExpression(expr StmtExpression) error
  visitStmtAssign(expr StmtAssign) error
  visitStmtBlock(expr StmtBlock) error
  visitStmtIf(expr StmtIf) error
  visitStmtWhile(expr StmtWhile) error
  visitStmtFunction (expr StmtFunction) error
}

type StmtVarDeclaration struct {
	Name        Token
	Initializer Expr
}

type StmtExpression struct {
	Expression Expr
}

type StmtPrint struct {
	Expression Expr
}

type StmtAssign struct {
  Name  Token
  Value Expr
}

type StmtIf struct {
  Condition Expr
  ThenBranch Stmt
  ElseBranch Stmt
}

type StmtWhile struct {
  Condition Expr 
  Body Stmt
}

type StmtBlock struct {
  Statements []Stmt
}

type StmtFunction struct {
  Name Token
  Params []Token
  Body []Stmt
}

func (stmt StmtVarDeclaration) accept(visitor StmtVisitor) error {
	return visitor.visitStmtVarDeclaration(stmt)
}

func (stmt StmtExpression) accept(visitor StmtVisitor) error {
	return visitor.visitStmtExpression(stmt)
}

func (stmt StmtPrint) accept(visitor StmtVisitor) error {
	return visitor.visitStmtPrint(stmt)
}

func (stmt StmtAssign) accept(visitor StmtVisitor) error {
  return visitor.visitStmtAssign(stmt)
}

func (stmt StmtBlock) accept(visitor StmtVisitor) error {
  return visitor.visitStmtBlock(stmt)
}

func (stmt StmtIf) accept(visitor StmtVisitor) error {
  return visitor.visitStmtIf(stmt)
}

func (stmt StmtWhile) accept(visitor StmtVisitor) error {
  return visitor.visitStmtWhile(stmt)
}

func (stmt StmtFunction) accept(visitor StmtVisitor) error {
  return visitor.visitStmtFunction(stmt)
}
