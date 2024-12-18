package main

import (
  "fmt"
)

type Parser struct {
  tokens []Token
  current int
  glox *Glox
}

type ParseError struct {
  token Token
  message string
}

func (e *ParseError) Error() string {
  return fmt.Sprintf("Error at '%s': %s", e.token.Lexeme, e.message)
}

func NewParser(tokens []Token) Parser {
  return Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() Expr {
  return p.expression()
}

func (p *Parser) expression() Expr {
  return p.equality()
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
  var expr = p.unary()
  for p.match(SLASH, STAR) {
    var operator Token = p.previous()
    var right Expr = p.unary()
    expr = ExprBinary{Left: expr, Operator: operator, Right: right}
  }
  return expr
}

func (p *Parser) unary() Expr {
  if p.match(BANG, MINUS) {
    var operator Token = p.previous()
    var right Expr = p.unary()
    return ExprUnary{Operator: operator, Right: right}
  }
  return p.primary()
}

func (p *Parser) primary() Expr {
  if p.match(FALSE) {
    return ExprLiteral{Value: false}
  }
  if p.match(TRUE) {
    return ExprLiteral{Value: true}
  }
  if p.match(NIL) {
    return ExprLiteral{Value: nil}
  }
  if p.match(NUMBER, STRING) {
    return ExprLiteral{Value: p.previous().Literal}
  }
  if p.match(LEFT_PAREN) {
    var expr Expr = p.expression()
    if _, err := p.consume(RIGHT_PAREN, "Expect ')' after expression."); err != nil {
      panic(err)
    }
    return ExprGrouping{Expression: expr}
  }
  //Golang throws without a return
  panic(p.error(p.peek(), "Expect expression."))
}

func (p *Parser) synchronize() {
  p.advance()
  for p.isAtEnd() {
    if p.previous().TokenType == SEMICOLON {
      return;
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
  if t.TokenType == EOF {
    p.glox.report(t.Line, "at end", message)
  } else {
    p.glox.report(t.Line, "at '" + t.Lexeme + "'", message)
  }
  return &ParseError{token: t, message: message}
}

func (p * Parser) match(types ...TokenType) bool {
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
  return p.tokens[p.current - 1]
}
