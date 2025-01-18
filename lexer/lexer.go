package lexer

import (
	"strings"
)

type Lexer struct {
	// Input string
	input string

	// Denotes current position in Lexer
	currentPosition int

	// Current byte.
	// TODO: Support Unicode
	ch byte
}

// Lexer constructor. We have to consume first iteration before returing the Lexer object.
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:           input,
		currentPosition: 0,
	}
	l.readChar()

	return l
}

// Consumes a character.
func (l *Lexer) readChar() {
	if l.currentPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.currentPosition]
	}

	l.currentPosition += 1
}

// Returns a Token based on the input
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = NewToken(EQUALS, "=")
	case '+':
		tok = NewToken(PLUS, "+")
	case ';':
		tok = NewToken(SEMICOLON, ";")
	case '*':
		tok = NewToken(STAR, "*")
	case '(':
		tok = NewToken(L_PAREN, "(")
	case ')':
		tok = NewToken(R_PAREN, ")")
	case ',':
		tok = NewToken(COMMA, ",")
	case 0:
		tok.literal = ""
		tok.TokenType = EOF
		return tok
	// Handle internal strings
	case '\'':
		{
			// Skip first '
			l.readChar()

			tok.literal = l.readString()
			tok.TokenType = STRING

			// Skip last '
			l.readChar()

			return tok
		}
	default:
		{
			if isLetter(l.ch) {
				tok.literal = strings.ToLower(l.readIdentifier())
				tok.TokenType = lookupIdent(tok.literal)

				return tok
			} else if isDigit(l.ch) {
				tok.TokenType = INT
				tok.literal = l.readNumber()

				return tok
			} else {
				return NewToken(ILLEGAL, "")
			}
		}
	}

	l.readChar()

	return tok
}

// Utility function to skip any whitespace.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' {
		l.readChar()
	}
}

// Reads identifier (or possibly keyword, but that is determined in lookupIdent function).
func (l *Lexer) readIdentifier() string {
	startPos := l.currentPosition - 1

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[startPos : l.currentPosition-1]
}

// Reads a string. It differs from readIdentifier, because the string itself might not contain
// only letters, but other characters as well.
func (l *Lexer) readString() string {
	startPos := l.currentPosition - 1

	for l.ch != '\'' {
		l.readChar()
	}

	return l.input[startPos : l.currentPosition-1]
}

// Reads a number. Currently we only support integers.
func (l *Lexer) readNumber() string {
	startPos := l.currentPosition - 1

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[startPos : l.currentPosition-1]
}

// Looks up an identifier in keyword map. If the keyword exists in the map, it is a keyword.
func lookupIdent(identifier string) TokenType {
	if token, ok := keywords[identifier]; ok {
		return token
	}

	return IDENTIFIER
}

// Utility function that determines wheter a character is a letter.
func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

// Utility function that determines whether a byte is a number.
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
