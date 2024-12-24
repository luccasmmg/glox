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

func (p *Parser) parse() ([]Stmt, []error) {
	var statements = []Stmt{}
  var errors = []error{}
	for !p.isAtEnd() {
    stmt, err := p.declaration()
    if err != nil {
      errors = append(errors, err)
    } else {
      statements = append(statements, stmt)
    }
	}
	return statements, errors
}

func (p *Parser) declaration() (Stmt, error) {
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

func (p *Parser) varDeclaration() (Stmt, error) {
  var token, err = p.consume(IDENTIFIER, "Expect variable name.")
  if err != nil {
    return nil, err
  }
  var initializer Expr;
  if (p.match(EQUAL)) {
    initializer, err = p.expression()
    if err != nil {
      return nil, err
    }
  }
  if _, err := p.consume(SEMICOLON, "Expect ';' after variable declaration."); err != nil {
    return nil, err
  }
  return StmtVarDeclaration{Name: token, Initializer: initializer}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
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

func (p *Parser) expressionStatement() (Stmt, error) {
	var value, err = p.expression()
  if err != nil {
    return nil, err
  }
	if _, err := p.consume(SEMICOLON, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return StmtExpression{Expression: value}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	var statements = []Stmt{}
  for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
    var value, err = p.declaration()
    if err != nil {
      return nil, err
    }
    statements = append(statements, value)
  }
  if _, err := p.consume(RIGHT_BRACE, "Expect '}' after block"); err != nil {
    return nil, err
  }
  return statements, nil
}

func (p *Parser) assignment() (Expr, error) {
  var expr = p.equality()
  if p.match(EQUAL) {
    var equals = p.previous()
    var value, _ = p.assignment()
    if variable, ok := expr.(ExprVariable); ok {
      var name = variable.Name
      return ExprAssign{Name: name, Value: value}, nil
    }
    fmt.Println(p.error(equals, "Invalid assignment target."))
  }
  return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) equality() Expr {
	var expr = p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		var operator Token = p.previous()
		var right Expr = p.comparison()
		expr = ExprBinary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	var expr = p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		var operator Token = p.previous()
		var right Expr = p.term()
		expr = ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) term() Expr {
	var expr = p.factor()
	for p.match(MINUS, PLUS) {
		var operator Token = p.previous()
		var right Expr = p.factor()
		expr = ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	var expr, _ = p.unary()
	for p.match(SLASH, STAR) {
		var operator Token = p.previous()
		var right, _ = p.unary()
		expr = ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		var operator Token = p.previous()
		var right, _ = p.unary()
		return ExprUnary{Operator: operator, Right: right}, nil
	}
	var expr, _ = p.primary()
  return expr, nil
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
			panic(err)
		}
		return ExprGrouping{Expression: expr}, nil
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
