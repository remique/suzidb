package lexer

import (
	"strings"
)

type Lexer struct {
	input           string
	currentPosition int
	ch              byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:           input,
		currentPosition: 0,
	}
	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.currentPosition >= len(l.input) {
		l.ch = '0'
	} else {
		l.ch = l.input[l.currentPosition]
	}

	l.currentPosition += 1
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(EQUALS, "=")
	case '+':
		tok = newToken(PLUS, "+")
	case ';':
		tok = newToken(SEMICOLON, ";")
	case '*':
		tok = newToken(STAR, "*")
	case '(':
		tok = newToken(L_PAREN, "(")
	case ')':
		tok = newToken(R_PAREN, ")")
	case ',':
		tok = newToken(COMMA, ",")
	// Handle internal strings
	case '\'':
		{
			// Skip first '
			l.readChar()

			tok.literal = l.readString()
			tok.tokenType = STRING

			// Skip last '
			l.readChar()

			return tok
		}
	default:
		{
			if isLetter(l.ch) {
				tok.literal = strings.ToLower(l.readIdentifier())
				tok.tokenType = lookupIdent(tok.literal)

				return tok
			} else if isDigit(l.ch) {
				tok.tokenType = INT
				tok.literal = l.readNumber()

				return tok
			} else {
				return newToken(ILLEGAL, "")
			}
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' {
		l.readChar()
	}
}

// ???
func (l *Lexer) readIdentifier() string {
	startPos := l.currentPosition - 1

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[startPos : l.currentPosition-1]
}

func (l *Lexer) readString() string {
	startPos := l.currentPosition - 1

	for l.ch != '\'' {
		l.readChar()
	}

	return l.input[startPos : l.currentPosition-1]
}

func (l *Lexer) readNumber() string {
	startPos := l.currentPosition - 1

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[startPos : l.currentPosition-1]
}

func lookupIdent(identifier string) TokenType {
	if token, ok := keywords[identifier]; ok {
		return token
	}

	return IDENTIFIER
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
