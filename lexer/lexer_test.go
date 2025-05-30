package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	lexer := NewLexer("+=abc def select inSerT null")

	tests := []struct {
		expectedToken Token
	}{
		{expectedToken: NewToken(PLUS, "+")},
		{expectedToken: NewToken(EQUALS, "=")},
		{expectedToken: NewToken(IDENTIFIER, "abc")},
		{expectedToken: NewToken(IDENTIFIER, "def")},
		{expectedToken: NewToken(SELECT, "select")},
		{expectedToken: NewToken(INSERT, "insert")},
		{expectedToken: NewToken(NULL, "null")},
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

func TestLexerSelectWithJoin(t *testing.T) {
	lexer := NewLexer(`SELECT ProductID, ProductName, CategoryName FROM Products
		INNER LEFT JOIN Categories ON Products.CategoryID = Categories.CategoryID;`)

	tests := []struct {
		expectedToken Token
	}{
		{expectedToken: NewToken(SELECT, "select")},
		{expectedToken: NewToken(IDENTIFIER, "productid")},
		{expectedToken: NewToken(COMMA, ",")},
		{expectedToken: NewToken(IDENTIFIER, "productname")},
		{expectedToken: NewToken(COMMA, ",")},
		{expectedToken: NewToken(IDENTIFIER, "categoryname")},
		{expectedToken: NewToken(FROM, "from")},
		{expectedToken: NewToken(IDENTIFIER, "products")},
		{expectedToken: NewToken(INNER, "inner")},
		{expectedToken: NewToken(LEFT, "left")},
		{expectedToken: NewToken(JOIN, "join")},
		{expectedToken: NewToken(IDENTIFIER, "categories")},
		{expectedToken: NewToken(ON, "on")},
		{expectedToken: NewToken(IDENTIFIER, "products")},
		{expectedToken: NewToken(DOT, ".")},
		{expectedToken: NewToken(IDENTIFIER, "categoryid")},
		{expectedToken: NewToken(EQUALS, "=")},
		{expectedToken: NewToken(IDENTIFIER, "categories")},
		{expectedToken: NewToken(DOT, ".")},
		{expectedToken: NewToken(IDENTIFIER, "categoryid")},
	}

	for _, test := range tests {
		tok := lexer.NextToken()

		assert.Equal(t, test.expectedToken, tok)
	}
}

func TestLexerInsertMultiple(t *testing.T) {
	lexer := NewLexer(`INSERT INTO products(id, productname) VALUES 
	(1, 'someproduct'),
	(2, 'anotherproduct');`)

	tests := []struct {
		expectedToken Token
	}{
		{expectedToken: NewToken(INSERT, "insert")},
		{expectedToken: NewToken(INTO, "into")},
		{expectedToken: NewToken(IDENTIFIER, "products")},
		{expectedToken: NewToken(L_PAREN, "(")},
		{expectedToken: NewToken(IDENTIFIER, "id")},
		{expectedToken: NewToken(COMMA, ",")},
		{expectedToken: NewToken(IDENTIFIER, "productname")},
		{expectedToken: NewToken(R_PAREN, ")")},
		{expectedToken: NewToken(VALUES, "values")},
		{expectedToken: NewToken(L_PAREN, "(")},
		{expectedToken: NewToken(INT, "1")},
		{expectedToken: NewToken(COMMA, ",")},
		{expectedToken: NewToken(STRING, "someproduct")},
		{expectedToken: NewToken(R_PAREN, ")")},
		{expectedToken: NewToken(COMMA, ",")},
		{expectedToken: NewToken(L_PAREN, "(")},
		{expectedToken: NewToken(INT, "2")},
		{expectedToken: NewToken(COMMA, ",")},
		{expectedToken: NewToken(STRING, "anotherproduct")},
		{expectedToken: NewToken(R_PAREN, ")")},
		{expectedToken: NewToken(SEMICOLON, ";")},
	}

	for _, test := range tests {
		tok := lexer.NextToken()

		assert.Equal(t, test.expectedToken, tok)
	}
}

func TestLexerSelectWithWhere(t *testing.T) {
	lexer := NewLexer(`SELECT * FROM sometbl WHERE x > 5;`)

	tests := []struct {
		expectedToken Token
	}{
		{expectedToken: NewToken(SELECT, "select")},
		{expectedToken: NewToken(STAR, "*")},
		{expectedToken: NewToken(FROM, "from")},
		{expectedToken: NewToken(IDENTIFIER, "sometbl")},
		{expectedToken: NewToken(WHERE, "where")},
		{expectedToken: NewToken(IDENTIFIER, "x")},
		{expectedToken: NewToken(GREATER, ">")},
		{expectedToken: NewToken(INT, "5")},
		{expectedToken: NewToken(SEMICOLON, ";")},
	}

	for _, test := range tests {
		tok := lexer.NextToken()

		assert.Equal(t, test.expectedToken, tok)
	}
}
