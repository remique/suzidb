package parser

import (
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

	items, err := parser.parseSelectItems()
	assert.NoError(t, err)

	for i := range tests.expectedItems {
		assert.Equal(t, tests.expectedItems[i], items[i])
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

	items, err := parser.parseSelectItems()
	assert.NoError(t, err)

	for i := range tests.expectedItems {
		assert.Equal(t, tests.expectedItems[i], items[i])
	}
}

func TestParseSelectItemsExpectedSemicolonOrWhere(t *testing.T) {
	lexer := l.NewLexer("a,c")
	parser := NewParser(*lexer)

	// NOTE: This should probably return nil?
	_, err := parser.parseSelectItems()
	assert.Error(t, err)
}

func TestParseSelectItemsExpectedComma(t *testing.T) {
	lexer := l.NewLexer("a,")
	parser := NewParser(*lexer)

	_, err := parser.parseSelectItems()
	assert.Error(t, err)
}

func TestParseSelectStatement(t *testing.T) {
	lexer := l.NewLexer("select * from myTable;")
	parser := NewParser(*lexer)

	expected := Statement{
		Kind: SelectKind,
		SelectStatement: &SelectStatement{
			SelectItems: &[]l.Token{
				l.NewToken(l.STAR, "*"),
			},
			From: &l.Token{TokenType: l.IDENTIFIER, Literal: "mytable"},
		},
	}
	res, err := parser.parseSelectStatement()
	assert.NoError(t, err)
	assert.Equal(t, &expected, res)
}

func TestParseSelectStatementNoFrom(t *testing.T) {
	lexer := l.NewLexer("select * myTable;")
	parser := NewParser(*lexer)

	_, err := parser.parseSelectStatement()
	assert.Error(t, err)
}

func TestParseSelectStatementNoIdentifierAfterFrom(t *testing.T) {
	lexer := l.NewLexer("select * FROM 10;")
	parser := NewParser(*lexer)

	_, err := parser.parseSelectStatement()
	assert.Error(t, err)
}

func TestParseStatementWithSelect(t *testing.T) {
	lexer := l.NewLexer("select * FROM myTable;")
	parser := NewParser(*lexer)

	expected := &Statement{
		Kind: SelectKind,
		SelectStatement: &SelectStatement{
			SelectItems: &[]l.Token{
				l.NewToken(l.STAR, "*"),
			},
			From: &l.Token{TokenType: l.IDENTIFIER, Literal: "mytable"},
		},
	}

	stmtRes, err := parser.ParseStatement()
	assert.NoError(t, err)
	assert.Equal(t, expected, stmtRes)
}

func TestParseStatementWithCreateTableInvalid(t *testing.T) {
	lexer := l.NewLexer("Create TaBle mytable hehe")
	parser := NewParser(*lexer)

	_, err := parser.parseCreateTableStatement()
	assert.Error(t, err)
}

func TestParseTableColumns(t *testing.T) {
	lexer := l.NewLexer("id int primary key, name text NOT NULL, surname text")
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
				Nullable: false,
			},
			{
				Name:     "surname",
				Type:     m.StringType,
				Nullable: true,
			},
		},
	}

	cols, pk, err := parser.parseCreateTableColumns()
	assert.NoError(t, err)
	assert.Equal(t, createTblStmt.Columns, cols)
	assert.Equal(t, createTblStmt.PrimaryKey, *pk)
}

func TestParseTableColumnsOnePKAllowed(t *testing.T) {
	lexer := l.NewLexer("id int primary key, name text primary key")
	parser := NewParser(*lexer)

	_, _, err := parser.parseCreateTableColumns()
	assert.Error(t, err)
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

func TestParseFromTableOnly(t *testing.T) {
	lexer := l.NewLexer("from mytable")
	parser := NewParser(*lexer)

	expected := &FromType{
		Table: &l.Token{Literal: "mytable", TokenType: l.IDENTIFIER},
		Kind:  UseTableKind,
	}

	fromRes, err := parser.parseSelectFrom()
	assert.NoError(t, err)
	assert.Equal(t, expected, fromRes, "Expected Froms to be the same")
}

// TODO: This currently fails because we dont have proper parsing of the expressions.
// func TestParseFromWithSingleJoin(t *testing.T) {
// 	lexer := l.NewLexer("from mytable left join anothertable on mytable.id = anothertable.id")
// 	parser := NewParser(*lexer)

// 	expected := &FromType{
// 		Kind: UseJoinKind,
// 		Join: &JoinFrom{
// 			Left: FromType{
// 				Table: &l.Token{Literal: "mytable", TokenType: l.IDENTIFIER},
// 				Kind:  UseTableKind,
// 			},
// 			Right: FromType{
// 				Table: &l.Token{Literal: "anothertable", TokenType: l.IDENTIFIER},
// 				Kind:  UseTableKind,
// 			},
// 			Kind: Left,

// 			Predicate: &Expression{
// 				Kind: BinaryKind,
// 				BinaryExpression: &BinaryExpression{
// 					Left: &Expression{
// 						Kind: QualifiedColumnKind,
// 						QualifiedColumnExpression: &QualifiedColumnExpression{
// 							TableName: &Expression{
// 								Kind:                 IdentifierKind,
// 								IdentifierExpression: &l.Token{TokenType: l.STRING, Literal: "mytable"},
// 							},
// 							ColumnName: &Expression{
// 								Kind:                 IdentifierKind,
// 								IdentifierExpression: &l.Token{TokenType: l.STRING, Literal: "id"},
// 							},
// 						},
// 					},
// 					Right: &Expression{
// 						Kind: QualifiedColumnKind,
// 						QualifiedColumnExpression: &QualifiedColumnExpression{
// 							TableName: &Expression{
// 								Kind:                 IdentifierKind,
// 								IdentifierExpression: &l.Token{TokenType: l.STRING, Literal: "anothertable"},
// 							},
// 							ColumnName: &Expression{
// 								Kind:                 IdentifierKind,
// 								IdentifierExpression: &l.Token{TokenType: l.STRING, Literal: "id"},
// 							},
// 						},
// 					},
// 					Operator: &l.Token{TokenType: l.EQUALS, Literal: "="},
// 				},
// 			},
// 		},
// 	}

// 	fromRes, err := parser.parseSelectFrom()
// 	assert.NoError(t, err)
// 	assert.Equal(t, expected, fromRes, "Expected Froms to be the same")
// }

func TestParseSelectClause2(t *testing.T) {
	lexer := l.NewLexer("select player.name, team.name, coach.name, withoutcol")
	parser := NewParser(*lexer)

	expected := &[]Expression{
		Expression{
			Kind: QualifiedColumnKind,
			QualifiedColumnExpression: &QualifiedColumnExpression{
				TableName: &Expression{
					Kind:                 IdentifierKind,
					IdentifierExpression: &l.Token{TokenType: l.IDENTIFIER, Literal: "player"},
				},
				ColumnName: &Expression{
					Kind:                 IdentifierKind,
					IdentifierExpression: &l.Token{TokenType: l.IDENTIFIER, Literal: "name"},
				},
			},
		},
		Expression{
			Kind: QualifiedColumnKind,
			QualifiedColumnExpression: &QualifiedColumnExpression{
				TableName: &Expression{
					Kind:                 IdentifierKind,
					IdentifierExpression: &l.Token{TokenType: l.IDENTIFIER, Literal: "team"},
				},
				ColumnName: &Expression{
					Kind:                 IdentifierKind,
					IdentifierExpression: &l.Token{TokenType: l.IDENTIFIER, Literal: "name"},
				},
			},
		},
		Expression{
			Kind: QualifiedColumnKind,
			QualifiedColumnExpression: &QualifiedColumnExpression{
				TableName: &Expression{
					Kind:                 IdentifierKind,
					IdentifierExpression: &l.Token{TokenType: l.IDENTIFIER, Literal: "coach"},
				},
				ColumnName: &Expression{
					Kind:                 IdentifierKind,
					IdentifierExpression: &l.Token{TokenType: l.IDENTIFIER, Literal: "name"},
				},
			},
		},
		Expression{
			Kind:                 IdentifierKind,
			IdentifierExpression: &l.Token{TokenType: l.IDENTIFIER, Literal: "withoutcol"},
		},
	}

	fromRes, err := parser.parseSelectClause2()
	assert.NoError(t, err)
	assert.Equal(t, expected, fromRes)
}
