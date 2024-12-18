package main 

import "fmt"

type Token struct {
  TokenType TokenType
  Lexeme string
  Literal interface{}
  Line int
}

// NewToken is a constructor function that initializes a Token with default values
func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
	}
}

func (t Token) String() string {
  return t.TokenType.String() + " " + t.Lexeme + " " + fmt.Sprint(t.Literal)
}
