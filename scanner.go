package main

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(source string) Scanner {
	return Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) advance() byte {
	c := s.source[s.current]
	s.current += 1
	return c
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	s.addToken(NUMBER, s.source[s.start:s.current])
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		Glox{}.reportError(s.line, "Unterminated string.")
		return
	}
	s.advance()
	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)
}

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) identifier() {
	keywords := map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType, nil)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	addToken := func(tokenType TokenType) {
		s.addToken(tokenType, nil)
	}
	switch c {
	case '(':
		addToken(LEFT_PAREN)
	case ')':
		addToken(RIGHT_PAREN)
	case '{':
		addToken(LEFT_BRACE)
	case '}':
		addToken(RIGHT_BRACE)
	case ',':
		addToken(COMMA)
	case '.':
		addToken(DOT)
	case '-':
		addToken(MINUS)
	case '+':
		addToken(PLUS)
	case ';':
		addToken(SEMICOLON)
	case '*':
		addToken(STAR)
	case '!':
		if s.match('=') {
			addToken(BANG_EQUAL)
		} else {
			addToken(BANG)
		}
	case '=':
		if s.match('=') {
			addToken(EQUAL_EQUAL)
		} else {
			addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			addToken(LESS_EQUAL)
		} else {
			addToken(LESS)
		}
	case '>':
		if s.match('=') {
			addToken(GREATER_EQUAL)
		} else {
			addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			addToken(SLASH)
		}
	case '"':
		s.string()
	case ' ', '\r', '\t':
		return
	case '\n':
		s.line++
		return
	default:
		if s.isDigit(c) {
			s.number()
			return
		} else if s.isAlpha(c) {
			s.identifier()
			return
		}
		Glox{}.reportError(s.line, "Unexpected character.")
	}
}
