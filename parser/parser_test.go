package parser

import (
	"testing"

	l "example.com/suzidb/lexer"
)

func TestParseSelectItemsSemicolon(t *testing.T) {
	lexer := l.NewLexer("a, b, c;")
	parser := NewParser(*lexer)

	tests := struct {
		expectedItems []l.Token
	}{
		expectedItems: []l.Token{
			l.NewToken(l.IDENTIFIER, "a"),
			l.NewToken(l.IDENTIFIER, "b"),
			l.NewToken(l.IDENTIFIER, "c"),
		},
	}

	items, _ := parser.parseSelectItems()
	for i := range tests.expectedItems {
		if items[i] != tests.expectedItems[i] {
			t.Fatalf("Expected item: %q, got: %q", tests.expectedItems[i], items[i])
		}
	}
}

func TestParseSelectItemsWhere(t *testing.T) {
	lexer := l.NewLexer("a, b, c WheRe")
	parser := NewParser(*lexer)

	tests := struct {
		expectedItems []l.Token
	}{
		expectedItems: []l.Token{
			l.NewToken(l.IDENTIFIER, "a"),
			l.NewToken(l.IDENTIFIER, "b"),
			l.NewToken(l.IDENTIFIER, "c"),
		},
	}

	items, _ := parser.parseSelectItems()
	for i := range tests.expectedItems {
		if items[i] != tests.expectedItems[i] {
			t.Fatalf("Expected item: %q, got: %q", tests.expectedItems[i], items[i])
		}
	}
}

// func TestParseSelectItemsInvalid(t *testing.T) {
// 	lexer := l.NewLexer("a+c")
// 	parser := NewParser(*lexer)

// 	_, err := parser.parseSelectItems()
// 	if err == nil {
// 		t.Fatalf("Expected err to not be nil")
// 	}
// }
