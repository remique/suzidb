package parser

import (
	"reflect"
	"testing"

	l "example.com/suzidb/lexer"
	m "example.com/suzidb/meta"
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
	if !reflect.DeepEqual(stmtRes, expected) {
		t.Fatal("Not deeply equal")
	}
}

func TestParseStatementWithCreateTableInvalid(t *testing.T) {
	lexer := l.NewLexer("Create TaBle mytable hehe")
	parser := NewParser(*lexer)

	_, err := parser.parseCreateTableStatement()
	if err == nil {
		t.Fatalf("Expected err: %q", err)
	}
}

func TestParseTableColumns(t *testing.T) {
	lexer := l.NewLexer("id int primary key, name text, surname text")
	parser := NewParser(*lexer)

	columns := []m.Column{
		{
			Name: "id",
			Type: m.IntType,
		},
		{
			Name: "name",
			Type: m.StringType,
		},
		{
			Name: "surname",
			Type: m.StringType,
		},
	}
	createTblStmt := CreateTableStatement{
		TableName:  "mytable",
		PrimaryKey: "id",
		Columns:    &columns,
	}

	cols, pk, err := parser.parseCreateTableColumns()
	if !reflect.DeepEqual(*cols, columns) {
		t.Fatalf("Columns Not deeply equal: %v, %v", *cols, columns)
	}
	if *pk != createTblStmt.PrimaryKey {
		t.Fatalf("pk not equal: %s, %s", *pk, createTblStmt.PrimaryKey)
	}
	if err != nil {
		t.Fatalf("Err: %s", err.Error())
	}
}

// func TestParseStatementWithCreateTable(t *testing.T) {
// 	lexer := l.NewLexer("Create TaBle mytable(id int primary key, name text, surname text);")
// 	parser := NewParser(*lexer)

// 	columns := []m.Column{
// 		{
// 			Name: "id",
// 			Type: m.IntType,
// 		},
// 		{
// 			Name: "name",
// 			Type: m.StringType,
// 		},
// 		{
// 			Name: "surname",
// 			Type: m.StringType,
// 		},
// 	}
// 	createTblStmt := CreateTableStatement{
// 		TableName:  "mytable",
// 		PrimaryKey: "id",
// 		Columns:    &columns,
// 	}

// 	expected := &Statement{CreateTableStatement: &createTblStmt, Kind: CreateTableKind}

// 	stmtRes, _ := parser.parseCreateTableStatement()
// 	if !reflect.DeepEqual(stmtRes, expected) {
// 		t.Fatal("Not deeply equal")
// 	}
// }
