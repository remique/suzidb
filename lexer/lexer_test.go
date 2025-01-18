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

func TestNextTokenSelectQuery(t *testing.T) {
	lexer := NewLexer("select * from myTable where id = 10;")

	tests := []struct {
		expectedToken Token
	}{
		{expectedToken: newToken(SELECT, "select")},
		{expectedToken: newToken(STAR, "*")},
		{expectedToken: newToken(FROM, "from")},
		{expectedToken: newToken(IDENTIFIER, "mytable")},
		{expectedToken: newToken(WHERE, "where")},
		{expectedToken: newToken(IDENTIFIER, "id")},
		{expectedToken: newToken(EQUALS, "=")},
		{expectedToken: newToken(INT, "10")},
		{expectedToken: newToken(SEMICOLON, ";")},
	}

	for _, test := range tests {
		tok := lexer.NextToken()

		if tok != test.expectedToken {
			t.Fatalf("Expected token %q, got %q", test.expectedToken, tok)
		}
	}
}

func TestNextTokenInsertQuery(t *testing.T) {
	lexer := NewLexer("INSERt into MytAbLe(a, b, c) VALUES ('hello', 'world', '!');")

	tests := []struct {
		expectedToken Token
	}{
		{expectedToken: newToken(INSERT, "insert")},
		{expectedToken: newToken(INTO, "into")},
		{expectedToken: newToken(IDENTIFIER, "mytable")},
		{expectedToken: newToken(L_PAREN, "(")},
		{expectedToken: newToken(IDENTIFIER, "a")},
		{expectedToken: newToken(COMMA, ",")},
		{expectedToken: newToken(IDENTIFIER, "b")},
		{expectedToken: newToken(COMMA, ",")},
		{expectedToken: newToken(IDENTIFIER, "c")},
		{expectedToken: newToken(R_PAREN, ")")},
		{expectedToken: newToken(VALUES, "values")},
		{expectedToken: newToken(L_PAREN, "(")},
		{expectedToken: newToken(STRING, "hello")},
		{expectedToken: newToken(COMMA, ",")},
		{expectedToken: newToken(STRING, "world")},
		{expectedToken: newToken(COMMA, ",")},
		{expectedToken: newToken(STRING, "!")},
		{expectedToken: newToken(R_PAREN, ")")},
		{expectedToken: newToken(SEMICOLON, ";")},
	}

	for _, test := range tests {
		tok := lexer.NextToken()

		if tok != test.expectedToken {
			t.Fatalf("Expected token %q, got %q", test.expectedToken, tok)
		}
	}
}
