package parser

import (
	"reflect"
	"testing"

	l "example.com/suzidb/lexer"
	m "example.com/suzidb/meta"
	"github.com/stretchr/testify/assert"
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
	lexer := l.NewLexer("id int primary key, name text NOT NULL, surname text")
	parser := NewParser(*lexer)

	columns := []m.Column{
		{
			Name:     "id",
			Type:     m.IntType,
			Nullable: false,
		},
		{
			Name:     "name",
			Type:     m.StringType,
			Nullable: false,
		},
		{
			Name:     "surname",
			Type:     m.StringType,
			Nullable: true,
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

func TestParseTableColumnsOnePKAllowed(t *testing.T) {
	lexer := l.NewLexer("id int primary key, name text primary key")
	parser := NewParser(*lexer)

	_, _, err := parser.parseCreateTableColumns()
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestParseStatementWithCreateTable(t *testing.T) {
	lexer := l.NewLexer("Create TaBle mytable(id int primary key, name text, surname text);")
	parser := NewParser(*lexer)

	createTblStmt := CreateTableStatement{
		TableName:  "mytable",
		PrimaryKey: "id",
		Columns: &[]m.Column{
			{
				Name:     "id",
				Type:     m.IntType,
				Nullable: false,
			},
			{
				Name:     "name",
				Type:     m.StringType,
				Nullable: true,
			},
			{
				Name:     "surname",
				Type:     m.StringType,
				Nullable: true,
			},
		},
	}

	expected := &Statement{CreateTableStatement: &createTblStmt, Kind: CreateTableKind}

	stmtRes, _ := parser.parseCreateTableStatement()
	assert.Equal(t, expected, stmtRes, "CreateTableStatements should be the same")
}

func TestParseInsertParseColumns(t *testing.T) {
	lexer := l.NewLexer("id, name, surname")
	parser := NewParser(*lexer)

	expected := []l.Token{
		l.NewToken(l.IDENTIFIER, "id"),
		l.NewToken(l.IDENTIFIER, "name"),
		l.NewToken(l.IDENTIFIER, "surname"),
	}

	cols, _ := parser.parseInsertColumnList()
	assert.Equal(t, expected, cols, "Insert column names should be the same")
}

func TestParseInsertParseValues(t *testing.T) {
	lexer := l.NewLexer("'john', 'smith', 1")
	parser := NewParser(*lexer)

	expected := []l.Token{
		l.NewToken(l.STRING, "john"),
		l.NewToken(l.STRING, "smith"),
		l.NewToken(l.INT, "1"),
	}

	vals, _ := parser.parseInsertValues()
	assert.Equal(t, expected, vals, "Insert values should be the same")
}

func TestParseInsertStatementWithCustomCols(t *testing.T) {
	lexer := l.NewLexer("insert into mytable(id, name, surname) values (1, 'john', 'smith');")
	parser := NewParser(*lexer)

	insertStmt := InsertStatement{
		TableName: "mytable",
		CustomColumns: []l.Token{
			l.NewToken(l.IDENTIFIER, "id"),
			l.NewToken(l.IDENTIFIER, "name"),
			l.NewToken(l.IDENTIFIER, "surname"),
		},
		Values: []l.Token{
			l.NewToken(l.INT, "1"),
			l.NewToken(l.STRING, "john"),
			l.NewToken(l.STRING, "smith"),
		},
	}

	expected := &Statement{InsertStatement: &insertStmt, Kind: InsertKind}

	stmt, err := parser.parseInsertStatement()
	if err != nil {
		t.Fatalf("err %s", err.Error())
	}
	assert.Equal(t, expected.InsertStatement, stmt.InsertStatement, "Expected statements to be the same")
}

func TestParseInsertStatementWithCustomColsNotMatching(t *testing.T) {
	lexer := l.NewLexer("insert into mytable(id, name) values (1, 'john', 'smith');")
	parser := NewParser(*lexer)

	// TODO: Use interface and assert.Equal(nil)
	_, err := parser.parseInsertStatement()
	assert.Error(t, err)
}

func TestParseInsertStatementWithoutColumnNames(t *testing.T) {
	lexer := l.NewLexer("insert into mytable values (1, 'john', 'smith');")
	parser := NewParser(*lexer)

	insertStmt := InsertStatement{
		TableName:     "mytable",
		CustomColumns: []l.Token{},
		Values: []l.Token{
			l.NewToken(l.INT, "1"),
			l.NewToken(l.STRING, "john"),
			l.NewToken(l.STRING, "smith"),
		},
	}

	expected := &Statement{InsertStatement: &insertStmt, Kind: InsertKind}

	stmt, err := parser.parseInsertStatement()
	if err != nil {
		t.Fatalf("err %s", err.Error())
	}
	assert.Equal(t, expected.InsertStatement, stmt.InsertStatement, "Expected statements to be the same")
}
