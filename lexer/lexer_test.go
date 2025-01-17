package lexer

import (
	"testing"
)

func TestReadChar(t *testing.T) {
	lexer := NewLexer("abc")

	tests := []struct {
		expectedChar byte
	}{
		{'a'},
		{'b'},
		{'c'},
		{'0'},
		{'0'},
	}

	for _, test := range tests {
		if lexer.ch != test.expectedChar {
			t.Fatalf("Expected %q, got %q", test.expectedChar, lexer.ch)
		}

		lexer.readChar()
	}
}

func TestNextToken(t *testing.T) {
	lexer := NewLexer("+=abc def select inSerT")

	tests := []struct {
		expectedToken Token
	}{
		{expectedToken: newToken(PLUS, "+")},
		{expectedToken: newToken(EQUALS, "=")},
		{expectedToken: newToken(IDENTIFIER, "abc")},
		{expectedToken: newToken(IDENTIFIER, "def")},
		{expectedToken: newToken(SELECT, "select")},
		{expectedToken: newToken(INSERT, "insert")},
	}

	for _, test := range tests {
		tok := lexer.NextToken()

		if tok != test.expectedToken {
			t.Fatalf("Expected token %q, got %q", test.expectedToken, tok)
		}
	}

}
