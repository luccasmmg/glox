package main

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
	glox    *Glox
}

type ParseError struct {
	token   Token
	message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Error at '%s': %s", e.token.Lexeme, e.message)
}

func NewParser(tokens []Token) Parser {
	return Parser{tokens: tokens, current: 0}
}

func (p *Parser) parse() ([]*Stmt, []error) {
	var statements = []*Stmt{}
	var errors = []error{}
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			errors = append(errors, err)
		} else {
			statements = append(statements, &stmt)
		}
	}
	return statements, errors
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(CLASS) {
		var value, err = p.classDeclaration()
		if err != nil {
			p.synchronize()
			return nil, nil
		} else {
			return value, nil
		}
	}
	if p.match(FUN) {
		var value, err = p.function("function")
		if err != nil {
			p.synchronize()
			return nil, nil
		} else {
			return value, nil
		}
	}
	if p.match(VAR) {
		var value, err = p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil, nil
		} else {
			return value, nil
		}
	} else {
		return p.statement()
	}
}

func (p *Parser) classDeclaration() (Stmt, error) {
	var name, err = p.consume(IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}
	var superclass ExprVariable
	if p.match(LESS) {
		var _, err = p.consume(IDENTIFIER, "Expect superclass name.")
		if err != nil {
			return nil, err
		}
		superclass = ExprVariable{
			Name: p.previous(),
		}
	}
	_, err = p.consume(LEFT_BRACE, "Expect { before class body.")
	if err != nil {
		return nil, err
	}
	methods := []Stmt{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		function, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, function)
	}
	_, err = p.consume(RIGHT_BRACE, "Expect } before class body.")
	if err != nil {
		return nil, err
	}
	fmt.Println(superclass)
	return StmtClass{
		Name:    name,
		Methods: methods,
		Superclass: func() *ExprVariable {
			if superclass != (ExprVariable{}) {
				return &superclass
			} else {
				return nil
			}
		}(),
	}, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	var token, err = p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer Expr
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(SEMICOLON, "Expect ';' after variable declaration."); err != nil {
		return nil, err
	}
	return StmtVarDeclaration{Name: token, Initializer: &initializer}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(LEFT_BRACE) {
		var value, err = p.block()
		if err != nil {
			return nil, err
		}
		return StmtBlock{Statements: value}, nil
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() (Stmt, error) {
	var initializer Stmt
	if _, err := p.consume(LEFT_PAREN, "Expect '(' at start of for loop."); err != nil {
		return nil, err
	}
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		var _initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
		initializer = _initializer
	} else {
		var _initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
		initializer = _initializer
	}
	var condition Expr = nil
	if !p.check(SEMICOLON) {
		var _condition, err = p.expression()
		if err != nil {
			return nil, err
		}
		condition = _condition
	}
	if _, err := p.consume(SEMICOLON, "Expect ';' after 'loop condition'."); err != nil {
		return nil, err
	}
	var increment Expr = nil
	if !p.check(RIGHT_PAREN) {
		var _increment, err = p.expression()
		if err != nil {
			return nil, err
		}
		increment = _increment
	}
	if _, err := p.consume(RIGHT_PAREN, "Expect ')' after 'for clauses'."); err != nil {
		return nil, err
	}
	var body, err = p.statement()
	if err != nil {
		return nil, err
	}
	if increment != nil {
    var statements []Stmt
    statements = append(statements, StmtExpression{Expression: &increment})
    stmtPointers := make([]*Stmt, len(statements))
    for i := range statements {
      stmtPointers[i] = &statements[i]
    }
		body = StmtBlock{
			Statements: stmtPointers,
		}
	}
	if condition == nil {
		condition = ExprLiteral{Value: true}
	}
	body = StmtWhile{
		Condition: condition,
		Body:      body,
	}
	if initializer != nil {
    var statements []Stmt
    statements = append([]Stmt{initializer})
    stmtPointers := make([]*Stmt, len(statements))
    for i := range statements {
      stmtPointers[i] = &statements[i]
    }
		body = StmtBlock{
			Statements: stmtPointers,
		}
	}
	return body, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	if _, err := p.consume(LEFT_PAREN, "Expect '(' after 'while'."); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(RIGHT_PAREN, "Expect ')' after while condition."); err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return StmtWhile{Condition: condition, Body: body}, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	if _, err := p.consume(LEFT_PAREN, "Expect '(' after 'if'."); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(RIGHT_PAREN, "Expect ')' after if condition."); err != nil {
		return nil, err
	}
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch Stmt = nil
	if p.match(ELSE) {
		branch, err := p.statement()
		if err != nil {
			return nil, err
		}
		elseBranch = branch
	}
	return StmtIf{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

func (p *Parser) printStatement() (Stmt, error) {
	var value, err = p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(SEMICOLON, "Expect ';' after value."); err != nil {
		return nil, err
	}
	return StmtPrint{Expression: value}, nil
}

func (p *Parser) returnStatement() (Stmt, error) {
	var keyword = p.previous()
	var value Expr = nil
	if !p.check(SEMICOLON) {
		var err error
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(SEMICOLON, "Expect ';' after return value."); err != nil {
		return nil, err
	}
	return StmtReturn{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	var value, err = p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(SEMICOLON, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return StmtExpression{Expression: &value}, nil
}

func (p *Parser) function(kind string) (Stmt, error) {
	var name, err = p.consume(IDENTIFIER, "Expect "+kind+" name")
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(LEFT_PAREN, "Expect '(' after "+kind+" name"); err != nil {
		return nil, err
	}
	var params []Token
	if !p.check(RIGHT_PAREN) {
		for {
			if len(params) >= 255 {
				parseError(p.peek(), "Cant have more than 255 parameters")
			}
			if identifier, err := p.consume(IDENTIFIER, "Expect parameter name"); err != nil {
				return nil, err
			} else {
				params = append(params, identifier)
			}
			if !p.match(COMMA) {
				break
			}
		}
	}
	if _, err := p.consume(RIGHT_PAREN, "Expect ')' after params"); err != nil {
		return nil, err
	}
	if _, err := p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body"); err != nil {
		return nil, err
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return StmtFunction{
		Name:   name,
		Params: params,
		Body:   body,
	}, nil
}

func (p *Parser) block() ([]*Stmt, error) {
	var statements = []*Stmt{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		var value, err = p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, &value)
	}
	if _, err := p.consume(RIGHT_BRACE, "Expect '}' after block"); err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	var expr, err = p.or()
	if err != nil {
		return nil, err
	}
	if p.match(EQUAL) {
		var equals = p.previous()
		var value, _ = p.assignment()
		if variable, ok := expr.(ExprVariable); ok {
			var name = variable.Name
			return ExprAssign{Name: name, Value: &value}, nil
			//TODO Reread this part of the book
		} else if get, ok := expr.(ExprGet); ok {
			return ExprSet{
				Object: get.Object,
				Name:   get.Name,
				Value:  &value,
			}, nil
		}
		fmt.Println(p.error(equals, "Invalid assignment target."))
	}
	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	var expr, err = p.and()
	if err != nil {
		return nil, err
	}
	for p.match(OR) {
		var operator = p.previous()
		var right, err = p.and()
		if err != nil {
			return nil, err
		}
		expr = ExprLogical{
			Operator: operator,
			Right:    &right,
			Left:     &expr,
		}
	}
	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	var expr, err = p.equality()
	if err != nil {
		return nil, err
	}
	for p.match(AND) {
		var operator = p.previous()
		var right, err = p.equality()
		if err != nil {
			return nil, err
		}
		expr = ExprLogical{
			Operator: operator,
			Right:    &right,
			Left:     &expr,
		}
	}
	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	var expr, err = p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		var operator Token = p.previous()
		var right, err = p.comparison()
		if err != nil {
			return nil, err
		}
		expr = ExprBinary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	var expr, err = p.term()
	if err != nil {
		return nil, err
	}
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		var operator Token = p.previous()
		var right, err = p.term()
		if err != nil {
			return nil, err
		}
		expr = ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	var expr, err = p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(MINUS, PLUS) {
		var operator Token = p.previous()
		var right, err = p.factor()
		if err != nil {
			return nil, err
		}
		expr = ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	var expr, err = p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(SLASH, STAR) {
		var operator Token = p.previous()
		var right, _ = p.unary()
		expr = ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		var operator Token = p.previous()
		var right, _ = p.unary()
		return ExprUnary{Operator: operator, Right: &right}, nil
	}
	var expr, err = p.call()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) call() (Expr, error) {
	var expr, err = p.primary()
	if err != nil {
		return nil, err
	}
	for {
		if p.match(LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(DOT) {
			name, err := p.consume(IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = ExprGet{
				Object: &expr,
				Name:   name,
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	var arguments []*Expr
	if !p.check(RIGHT_PAREN) {
		for {
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, &arg)
			if len(arguments) >= 255 {
				fmt.Println("Cant have more than 255 arguments")
			}
			if !p.match(COMMA) {
				break
			}
		}
	}
	paren, err := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return ExprCall{Callee: &callee, Paren: paren, Arguments: arguments}, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return ExprLiteral{Value: false}, nil
	}
	if p.match(TRUE) {
		return ExprLiteral{Value: true}, nil
	}
	if p.match(IDENTIFIER) {
		return ExprVariable{Name: p.previous()}, nil
	}
	if p.match(THIS) {
		return ExprThis{Keyword: p.previous()}, nil
	}
	if p.match(SUPER) {
		keyword := p.previous()
		if _, err := p.consume(DOT, "Expect '.' after super."); err != nil {
			return nil, err
		}
		method, err := p.consume(IDENTIFIER, "Expect superclass method name.")
		if err != nil {
			return nil, err
		}
		return ExprSuper{Keyword: keyword, Method: method}, nil
	}
	if p.match(NIL) {
		return ExprLiteral{Value: nil}, nil
	}
	if p.match(NUMBER, STRING) {
		return ExprLiteral{Value: p.previous().Literal}, nil
	}
	if p.match(LEFT_PAREN) {
		var expr, err = p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			return nil, err
		}
		return ExprGrouping{Expression: &expr}, nil
	}
	//Golang throws without a return
	panic(p.error(p.peek(), "Expect expression."))
}

func (p *Parser) synchronize() {
	p.advance()
	for p.isAtEnd() {
		if p.previous().TokenType == SEMICOLON {
			return
		}
		switch p.peek().TokenType {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) consume(t TokenType, message string) (Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	return Token{}, p.error(p.peek(), message)
}

func (p *Parser) error(t Token, message string) error {
	//if t.TokenType == EOF {
	//	p.glox.report(t.Line, "at end", message)
	//} else {
	//	p.glox.report(t.Line, "at '"+t.Lexeme+"'", message)
	//}
	return &ParseError{token: t, message: message}
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == t
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}
