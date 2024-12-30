package main

import (
	"fmt"
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          Stack[map[string]bool]
	currentFunction FunctionType
	currentClass    ClassType
}

type ResolverError struct {
	token   Token
	message string
}

type FunctionType string

const (
	NONE_FUNCTION FunctionType = "NONE"
	FUNCTION      FunctionType = "FUNCTION"
	METHOD        FunctionType = "METHOD"
	INITIALIZER   FunctionType = "INITIALIZER"
)

type ClassType string

const (
	NONE_CLASS     ClassType = "NONE"
	CLASS_RESOLVER ClassType = "CLASS"
  SUBCLASS       ClassType = "SUBCLASS"
)

func (e *ResolverError) Error() string {
	return fmt.Sprintf("Error at '%s': %s", e.token.Lexeme, e.message)
}

func NewResolver(interpreter *Interpreter) Resolver {
	return Resolver{
		interpreter:     interpreter,
		scopes:          Stack[map[string]bool]{},
		currentFunction: NONE_FUNCTION,
		currentClass:    NONE_CLASS,
	}
}

func (r *Resolver) visitStmtBlock(stmt StmtBlock) error {
	r.beginScope()
	err := r.resolveStatements(stmt.Statements)
	if err != nil {
		return err
	}
	r.endScope()
	return nil
}

func (r *Resolver) visitStmtClass(stmt StmtClass) error {
	var enclosingClass = r.currentClass
	r.currentClass = "CLASS"
	r.declare(stmt.Name)
	r.define(stmt.Name)
	if stmt.Superclass != nil && stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
		return &ResolverError{
			token:   stmt.Superclass.Name,
			message: "A class cant inherint from itself",
		}
	}
	if stmt.Superclass != nil {
    r.currentClass = "SUBCLASS"
		r.resolveExpr(stmt.Superclass)
	}
	if stmt.Superclass != nil {
		r.beginScope()
		value, err := r.scopes.Peek()
		if err != nil {
			return err
		}
		value["super"] = true
	}
	r.beginScope()
	value, err := r.scopes.Peek()
	if err != nil {
		return err
	}
	value["this"] = true
	for _, method := range stmt.Methods {
		var declaration FunctionType = METHOD
		if (method.(StmtFunction)).Name.Lexeme == "init" {
			declaration = INITIALIZER
		}
		r.resolveFunction(method.(StmtFunction), declaration)
	}
	r.endScope()
  if stmt.Superclass != nil {
    r.endScope()
  }
	r.currentClass = enclosingClass
	return nil
}

func (r *Resolver) visitStmtVarDeclaration(stmt StmtVarDeclaration) error {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		_, err := r.resolveExpr(stmt.Initializer)
		if err != nil {
			return err
		}
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) visitVariableExpr(expr ExprVariable) (interface{}, error) {
	var value, err = r.scopes.Peek()
	if err != nil {
		return err, nil
	}
	if value[expr.Name.Lexeme] == false {
		return nil, &ResolverError{
			expr.Name,
			"Can't read local variable in its own initializer.",
		}
	}
  fmt.Printf("Resolver: Memory address: %p\n", expr)
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) visitThisExpr(expr ExprThis) (interface{}, error) {
	if r.currentClass == NONE_CLASS {
		return nil, &ResolverError{
			expr.Keyword,
			"Can't use 'this' outside of a class.",
		}
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) visitSuperExpr(expr ExprSuper) (interface{}, error) {
  if r.currentClass == "NONE_CLASS" {
		return nil, &ResolverError{
			expr.Keyword,
			"Can't use 'super' outside of a class.",
		}
  } else if r.currentClass != "SUBCLASS" {
		return nil, &ResolverError{
			expr.Keyword,
			"Can't use 'super' in a class with no subclass.",
		}

  }
	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) visitAssignExpr(expr ExprAssign) (interface{}, error) {
	_, err := r.resolveExpr(*expr.Value)
	if err != nil {
		return err, nil
	}
	err = r.resolveLocal(expr, expr.Name)
	if err != nil {
		return err, nil
	}
	return nil, nil
}

func (r *Resolver) visitStmtFunction(stmt StmtFunction) error {
	err := r.declare(stmt.Name)
	if err != nil {
		return err
	}
	err = r.define(stmt.Name)
	if err != nil {
		return err
	}
	err = r.resolveFunction(stmt, FUNCTION)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resolver) visitStmtExpression(stmt StmtExpression) error {
	_, err := r.resolveExpr(stmt.Expression)
	return err
}

func (r *Resolver) visitStmtIf(stmt StmtIf) error {
	_, err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	err = r.resolveStmt(stmt.ThenBranch)
	if err != nil {
		return err
	}
	if stmt.ElseBranch != nil {
		err = r.resolveStmt(stmt.ElseBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) visitStmtWhile(stmt StmtWhile) error {
	_, err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	err = r.resolveStmt(stmt.Body)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resolver) visitStmtPrint(stmt StmtPrint) error {
	_, err := r.resolveExpr(stmt.Expression)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resolver) visitStmtReturn(stmt StmtReturn) error {
	if r.currentFunction == NONE_FUNCTION {
		return &ResolverError{
			stmt.Keyword,
			"Can't return from top-level code.",
		}
	}
	if stmt.Value != nil {
		if r.currentFunction == INITIALIZER {
			return &ResolverError{
				stmt.Keyword,
				"Can't return a value from an initializer.",
			}
		}
		var _, err = r.resolveExpr(stmt.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) visitBinaryExpr(expr ExprBinary) (interface{}, error) {
	_, err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, nil
	}
	_, err = r.resolveExpr(expr.Right)
	if err != nil {
		return nil, nil
	}
	return nil, nil
}

func (r *Resolver) visitGroupingExpr(expr ExprGrouping) (interface{}, error) {
	return r.resolveExpr(expr.Expression)
}

func (r *Resolver) visitLiteralExpr(expr ExprLiteral) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) visitUnaryExpr(expr ExprUnary) (interface{}, error) {
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) visitCallExpr(expr ExprCall) (interface{}, error) {
	_, err := r.resolveExpr(*expr.Callee)
	if err != nil {
		return err, nil
	}
	for _, arg := range expr.Arguments {
		_, err = r.resolveExpr(*arg)
		if err != nil {
			return err, nil
		}
	}
	return nil, nil
}

func (r *Resolver) visitGetExpr(expr ExprGet) (interface{}, error) {
	_, err := r.resolveExpr(expr.Object)
	if err != nil {
		return err, nil
	}
	return nil, nil
}

func (r *Resolver) visitLogicalExpr(expr ExprLogical) (interface{}, error) {
	_, err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, nil
	}
	_, err = r.resolveExpr(expr.Right)
	if err != nil {
		return nil, nil
	}
	return nil, nil
}

func (r *Resolver) visitSetExpr(expr ExprSet) (interface{}, error) {
	_, err := r.resolveExpr(expr.Value)
	if err != nil {
		return err, nil
	}
	_, err = r.resolveExpr(expr.Object)
	if err != nil {
		return err, nil
	}
	return nil, nil
}

func (r *Resolver) resolveFunction(stmt StmtFunction, _type FunctionType) error {
	var enclosingFunction = r.currentFunction
	r.currentFunction = _type
	r.beginScope()
	for _, param := range stmt.Params {
		err := r.declare(param)
		if err != nil {
			return err
		}
		err = r.define(param)
		if err != nil {
			return err
		}
	}
	r.resolveStatements(stmt.Body)
	r.endScope()
	r.currentFunction = enclosingFunction
	return nil
}

func (r *Resolver) resolveStmt(stmt Stmt) error {
	return stmt.accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) (interface{}, error) {
	return expr.accept(r)
}

func (r *Resolver) resolveStatements(statements []Stmt) error {
	for _, stmt := range statements {
		err := r.resolveStmt(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveLocal(expr Expr, name Token) error {
	for i := len(r.scopes.elements) - 1; i >= 0; i-- {
		scope := r.scopes.elements[i]
		if _, ok := scope[name.Lexeme]; ok {
			r.interpreter.resolve(expr, r.scopes.Size()-1-i)
			return nil
		}
	}
	return nil
}

func (r *Resolver) declare(name Token) error {
	if r.scopes.IsEmpty() {
		return nil
	}
	var scope, err = r.scopes.Peek()
	if err != nil {
		return err
	}
	if _, ok := scope[name.Lexeme]; ok {
		return &ResolverError{
			name,
			"Variable with this name already declared in this scope.",
		}
	}
	scope[name.Lexeme] = false
	return nil
}

func (r *Resolver) define(name Token) error {
	if r.scopes.IsEmpty() {
		return nil
	}
	var scope, err = r.scopes.Peek()
	if err != nil {
		return err
	}
	scope[name.Lexeme] = true
	return nil
}

func (r *Resolver) beginScope() error {
	var scope = make(map[string]bool)
	r.scopes.Push(scope)
	return nil
}

func (r *Resolver) endScope() error {
	var _, err = r.scopes.Pop()
	return err
}
