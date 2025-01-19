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
	lexer := l.NewLexer("a, b, c frOm")
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

func TestParseSelectItemsExpectedSemicolonOrWhere(t *testing.T) {
	lexer := l.NewLexer("a,c")
	parser := NewParser(*lexer)

	_, err := parser.parseSelectItems()
	if err == nil {
		t.Fatalf("Expected err to not be nil")
	}
}

func TestParseSelectItemsExpectedComma(t *testing.T) {
	lexer := l.NewLexer("a,")
	parser := NewParser(*lexer)

	_, err := parser.parseSelectItems()
	if err == nil {
		t.Fatalf("Expected err to not be nil")
	}
}

func TestParseSelectStatement(t *testing.T) {
	lexer := l.NewLexer("select * from myTable;")
	parser := NewParser(*lexer)

	_, err := parser.parseSelectStatement()
	if err != nil {
		t.Fatalf("err: %q", err)
	}
}

func TestParseSelectStatementNoFrom(t *testing.T) {
	lexer := l.NewLexer("select * myTable;")
	parser := NewParser(*lexer)

	_, err := parser.parseSelectStatement()
	if err == nil {
		t.Fatalf("Expected err: %q", err)
	}
}

func TestParseSelectStatementNoIdentifierAfterFrom(t *testing.T) {
	lexer := l.NewLexer("select * FROM 10;")
	parser := NewParser(*lexer)

	_, err := parser.parseSelectStatement()
	if err == nil {
		t.Fatalf("Expected err: %q", err)
	}
}

func TestParseStatementWithSelect(t *testing.T) {
	lexer := l.NewLexer("select * FROM myTable;")
	parser := NewParser(*lexer)

	from := l.NewToken(l.IDENTIFIER, "mytable")
	selectStmt := SelectStatement{
		SelectItems: &[]l.Token{
			l.NewToken(l.STAR, "*"),
		},
		From: &from,
	}
	expected := &Statement{SelectStatement: &selectStmt, Kind: SelectKind}

	stmtRes, _ := parser.parseStatement()
	if stmtRes.Kind != expected.Kind {
		t.Fatalf("Expected KIND %q to equal %q", stmtRes.Kind, expected.Kind)
	}
	if *stmtRes.SelectStatement.From != *expected.SelectStatement.From {
		t.Fatalf("Expected FROM %q to equal %q", stmtRes.SelectStatement.From,
			expected.SelectStatement.From)
	}
	if stmtRes.SelectStatement.SelectItems != expected.SelectStatement.SelectItems {
		t.Fatalf("Expected ITEMS %q to equal %q", stmtRes.SelectStatement.SelectItems,
			expected.SelectStatement.SelectItems)
	}
	// TODO
}
