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
		{0},
		{0},
	}

	for _, test := range tests {
		if lexer.ch != test.expectedChar {
			t.Fatalf("Expected %q, got %q", test.expectedChar, lexer.ch)
		}

		lexer.readChar()
	}
}

func TestNextToken2(t *testing.T) {
	lexer := NewLexer("a, b, c;")

	tests := []struct {
		expectedToken Token
	}{
		{expectedToken: NewToken(IDENTIFIER, "a")},
		{expectedToken: NewToken(COMMA, ",")},
		{expectedToken: NewToken(IDENTIFIER, "b")},
		{expectedToken: NewToken(COMMA, ",")},
		{expectedToken: NewToken(IDENTIFIER, "c")},
		{expectedToken: NewToken(SEMICOLON, ";")},
		{expectedToken: NewToken(EOF, "")},
	}

	for _, test := range tests {
		tok := lexer.NextToken()

		if tok != test.expectedToken {
			t.Fatalf("Expected token %q, got %q", test.expectedToken, tok)
		}
	}
}

func TestNextToken(t *testing.T) {
	lexer := NewLexer("+=abc def select inSerT")

	tests := []struct {
		expectedToken Token
	}{
		{expectedToken: NewToken(PLUS, "+")},
		{expectedToken: NewToken(EQUALS, "=")},
		{expectedToken: NewToken(IDENTIFIER, "abc")},
		{expectedToken: NewToken(IDENTIFIER, "def")},
		{expectedToken: NewToken(SELECT, "select")},
		{expectedToken: NewToken(INSERT, "insert")},
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
		{expectedToken: NewToken(SELECT, "select")},
		{expectedToken: NewToken(STAR, "*")},
		{expectedToken: NewToken(FROM, "from")},
		{expectedToken: NewToken(IDENTIFIER, "mytable")},
		{expectedToken: NewToken(WHERE, "where")},
		{expectedToken: NewToken(IDENTIFIER, "id")},
		{expectedToken: NewToken(EQUALS, "=")},
		{expectedToken: NewToken(INT, "10")},
		{expectedToken: NewToken(SEMICOLON, ";")},
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
		{expectedToken: NewToken(INSERT, "insert")},
		{expectedToken: NewToken(INTO, "into")},
		{expectedToken: NewToken(IDENTIFIER, "mytable")},
		{expectedToken: NewToken(L_PAREN, "(")},
		{expectedToken: NewToken(IDENTIFIER, "a")},
		{expectedToken: NewToken(COMMA, ",")},
		{expectedToken: NewToken(IDENTIFIER, "b")},
		{expectedToken: NewToken(COMMA, ",")},
		{expectedToken: NewToken(IDENTIFIER, "c")},
		{expectedToken: NewToken(R_PAREN, ")")},
		{expectedToken: NewToken(VALUES, "values")},
		{expectedToken: NewToken(L_PAREN, "(")},
		{expectedToken: NewToken(STRING, "hello")},
		{expectedToken: NewToken(COMMA, ",")},
		{expectedToken: NewToken(STRING, "world")},
		{expectedToken: NewToken(COMMA, ",")},
		{expectedToken: NewToken(STRING, "!")},
		{expectedToken: NewToken(R_PAREN, ")")},
		{expectedToken: NewToken(SEMICOLON, ";")},
	}

	for _, test := range tests {
		tok := lexer.NextToken()

		if tok != test.expectedToken {
			t.Fatalf("Expected token %q, got %q", test.expectedToken, tok)
		}
	}
}

func TestNextTokenCreateTableQuery(t *testing.T) {
	lexer := NewLexer("CREATE TABLE cars(brand TEXT primary kEy, year INT not NuLL);")

	tests := []struct {
		expectedToken Token
	}{
		{expectedToken: NewToken(CREATE, "create")},
		{expectedToken: NewToken(TABLE, "table")},
		{expectedToken: NewToken(IDENTIFIER, "cars")},
		{expectedToken: NewToken(L_PAREN, "(")},
		{expectedToken: NewToken(IDENTIFIER, "brand")},
		{expectedToken: NewToken(TEXT_TYPE, "text")},
		{expectedToken: NewToken(PRIMARY, "primary")},
		{expectedToken: NewToken(KEY, "key")},
		{expectedToken: NewToken(COMMA, ",")},
		{expectedToken: NewToken(IDENTIFIER, "year")},
		{expectedToken: NewToken(INT_TYPE, "int")},
		{expectedToken: NewToken(NOT, "not")},
		{expectedToken: NewToken(NULL, "null")},
		{expectedToken: NewToken(R_PAREN, ")")},
		{expectedToken: NewToken(SEMICOLON, ";")},
	}

	for _, test := range tests {
		tok := lexer.NextToken()

		if tok != test.expectedToken {
			t.Fatalf("Expected token %q, got %q", test.expectedToken, tok)
		}
	}
}
