package parser

import (
	"testing"

	l "example.com/suzidb/lexer"
	m "example.com/suzidb/meta"
	"github.com/stretchr/testify/assert"
)

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
	// lexer := l.NewLexer("'john', 'smith', 1, NULL")
	lexer := l.NewLexer("'john', 3, null")
	parser := NewParser(*lexer)

	expected := []Expression{
		{
			Kind:              LiteralKind,
			LiteralExpression: &l.Token{TokenType: l.STRING, Literal: "john"},
		},
		{
			Kind: ConstExprKind,
			ConstExpression: &ConstExpression{
				Kind: IntKind,
				Int:  &l.Token{TokenType: l.INT, Literal: "3"},
			},
		},
		{
			Kind: ConstExprKind,
			ConstExpression: &ConstExpression{
				Kind: NullKind,
				Null: &l.Token{TokenType: l.NULL, Literal: "null"},
			},
		},
	}

	vals, err := parser.parseInsertValues()
	assert.NoError(t, err)
	assert.Equal(t, &expected, vals, "Insert values should be the same")
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
		Values: &[][]Expression{

			{
				{
					Kind: ConstExprKind,
					ConstExpression: &ConstExpression{
						Kind: IntKind,
						Int:  &l.Token{TokenType: l.INT, Literal: "1"},
					},
				},
				{
					Kind: LiteralKind,
					LiteralExpression: &l.Token{
						TokenType: l.STRING,
						Literal:   "john",
					},
				},
				{
					Kind: LiteralKind,
					LiteralExpression: &l.Token{
						TokenType: l.STRING,
						Literal:   "smith",
					},
				},
			},
		},
	}

	expected := &Statement{InsertStatement: &insertStmt, Kind: InsertKind}

	stmt, err := parser.parseInsertStatement()
	if err != nil {
		t.Fatalf("err %s", err.Error())
	}
	assert.Equal(t, expected.InsertStatement, stmt.InsertStatement, "Expected statements to be the same")
}

func TestParseInsertStatementMultiple(t *testing.T) {
	lexer := l.NewLexer("insert into mytable(id, name, surname) values (1, 'john', 'smith'), (2, 'alice', 'kovalsky');")
	parser := NewParser(*lexer)

	insertStmt := InsertStatement{
		TableName: "mytable",
		CustomColumns: []l.Token{
			l.NewToken(l.IDENTIFIER, "id"),
			l.NewToken(l.IDENTIFIER, "name"),
			l.NewToken(l.IDENTIFIER, "surname"),
		},
		Values: &[][]Expression{
			{
				{
					Kind: ConstExprKind,
					ConstExpression: &ConstExpression{
						Kind: IntKind,
						Int:  &l.Token{TokenType: l.INT, Literal: "1"},
					},
				},
				{
					Kind: LiteralKind,
					LiteralExpression: &l.Token{
						TokenType: l.STRING,
						Literal:   "john",
					},
				},
				{
					Kind: LiteralKind,
					LiteralExpression: &l.Token{
						TokenType: l.STRING,
						Literal:   "smith",
					},
				},
			},
			{
				{
					Kind: ConstExprKind,
					ConstExpression: &ConstExpression{
						Kind: IntKind,
						Int:  &l.Token{TokenType: l.INT, Literal: "2"},
					},
				},
				{
					Kind: LiteralKind,
					LiteralExpression: &l.Token{
						TokenType: l.STRING,
						Literal:   "alice",
					},
				},
				{
					Kind: LiteralKind,
					LiteralExpression: &l.Token{
						TokenType: l.STRING,
						Literal:   "kovalsky",
					},
				},
			},
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
		Values: &[][]Expression{
			{
				{
					Kind: ConstExprKind,
					ConstExpression: &ConstExpression{
						Kind: IntKind,
						Int:  &l.Token{TokenType: l.INT, Literal: "1"},
					},
				},
				{
					Kind: LiteralKind,
					LiteralExpression: &l.Token{
						TokenType: l.STRING,
						Literal:   "john",
					},
				},
				{
					Kind: LiteralKind,
					LiteralExpression: &l.Token{
						TokenType: l.STRING,
						Literal:   "smith",
					},
				},
			},
		},
	}

	expected := &Statement{InsertStatement: &insertStmt, Kind: InsertKind}

	stmt, err := parser.parseInsertStatement()
	if err != nil {
		t.Fatalf("err %s", err.Error())
	}
	assert.Equal(t, expected.InsertStatement, stmt.InsertStatement, "Expected statements to be the same")
}

func TestParseSelectClause2(t *testing.T) {
	lexer := l.NewLexer("select player.name, team.name, coach.name, withoutcol")
	parser := NewParser(*lexer)

	expected := &[]Expression{
		Expression{
			Kind: ColumnKind,
			ColumnExpression: &ColumnExpression{
				TableName:  "player",
				ColumnName: "name",
			},
		},
		Expression{
			Kind: ColumnKind,
			ColumnExpression: &ColumnExpression{
				TableName:  "team",
				ColumnName: "name",
			},
		},
		Expression{
			Kind: ColumnKind,
			ColumnExpression: &ColumnExpression{
				TableName:  "coach",
				ColumnName: "name",
			},
		},
		Expression{
			Kind: ColumnKind,
			ColumnExpression: &ColumnExpression{
				TableName: "withoutcol",
			},
		},
	}

	fromRes, err := parser.parseSelectClause()
	assert.NoError(t, err)
	assert.Equal(t, expected, fromRes)
}

func TestParseFromClauseTableOnly(t *testing.T) {
	lexer := l.NewLexer("from sometbl")
	parser := NewParser(*lexer)

	expected := &TableFrom{TableName: "sometbl"}

	fromRes, err := parser.parseFromClause()
	assert.NoError(t, err)
	assert.Equal(t, expected, fromRes)
}

func TestParseFromClauseWithSingleJoin(t *testing.T) {
	lexer := l.NewLexer("from sometbl left join othertbl on sometbl.x = othertbl.y")
	parser := NewParser(*lexer)

	expected := &JoinFrom{
		Left:     &TableFrom{TableName: "sometbl"},
		Right:    &TableFrom{TableName: "othertbl"},
		JoinKind: Left,
		Predicate: &Expression{
			Kind: BinaryKind,
			BinaryExpression: &BinaryExpression{
				Left: &Expression{
					Kind: ColumnKind,
					ColumnExpression: &ColumnExpression{
						TableName:  "sometbl",
						ColumnName: "x",
					},
				},
				Right: &Expression{
					Kind: ColumnKind,
					ColumnExpression: &ColumnExpression{
						TableName:  "othertbl",
						ColumnName: "y",
					},
				},
				Operator: &l.Token{TokenType: l.EQUALS, Literal: "="},
			},
		},
	}

	fromRes, err := parser.parseFromClause()
	assert.NoError(t, err)
	assert.Equal(t, expected, fromRes)
}

func TestParseFromClauseWitDoubleJoin(t *testing.T) {
	lexer := l.NewLexer("from sometbl left join othertbl on sometbl.x = othertbl.y left join anotherone on othertbl.y = anotherone.z")
	parser := NewParser(*lexer)

	left := &JoinFrom{
		Left:     &TableFrom{TableName: "sometbl"},
		Right:    &TableFrom{TableName: "othertbl"},
		JoinKind: Left,
		Predicate: &Expression{
			Kind: BinaryKind,
			BinaryExpression: &BinaryExpression{
				Left: &Expression{
					Kind: ColumnKind,
					ColumnExpression: &ColumnExpression{
						TableName:  "sometbl",
						ColumnName: "x",
					},
				},
				Right: &Expression{
					Kind: ColumnKind,
					ColumnExpression: &ColumnExpression{
						TableName:  "othertbl",
						ColumnName: "y",
					},
				},
				Operator: &l.Token{TokenType: l.EQUALS, Literal: "="},
			},
		},
	}

	expected := &JoinFrom{
		Left:     left,
		Right:    &TableFrom{TableName: "anotherone"},
		JoinKind: Left,
		Predicate: &Expression{
			Kind: BinaryKind,
			BinaryExpression: &BinaryExpression{
				Left: &Expression{
					Kind: ColumnKind,
					ColumnExpression: &ColumnExpression{
						TableName:  "othertbl",
						ColumnName: "y",
					},
				},
				Right: &Expression{
					Kind: ColumnKind,
					ColumnExpression: &ColumnExpression{
						TableName:  "anotherone",
						ColumnName: "z",
					},
				},
				Operator: &l.Token{TokenType: l.EQUALS, Literal: "="},
			},
		},
	}

	fromRes, err := parser.parseFromClause()

	assert.NoError(t, err)
	assert.Equal(t, expected, fromRes)
}

func TestParseSelectStatementFull(t *testing.T) {
	lexer := l.NewLexer("select sometbl.x, othertbl.x from sometbl left join othertbl on sometbl.x = othertbl.y left join anotherone on othertbl.y = anotherone.z")
	parser := NewParser(*lexer)

	items := &[]Expression{
		Expression{
			Kind: ColumnKind,
			ColumnExpression: &ColumnExpression{
				TableName:  "sometbl",
				ColumnName: "x",
			},
		},
		Expression{
			Kind: ColumnKind,
			ColumnExpression: &ColumnExpression{
				TableName:  "othertbl",
				ColumnName: "x",
			},
		},
	}

	from := &JoinFrom{
		Left: &JoinFrom{
			Left:     &TableFrom{TableName: "sometbl"},
			Right:    &TableFrom{TableName: "othertbl"},
			JoinKind: Left,
			Predicate: &Expression{
				Kind: BinaryKind,
				BinaryExpression: &BinaryExpression{
					Left: &Expression{
						Kind: ColumnKind,
						ColumnExpression: &ColumnExpression{
							TableName:  "sometbl",
							ColumnName: "x",
						},
					},
					Right: &Expression{
						Kind: ColumnKind,
						ColumnExpression: &ColumnExpression{
							TableName:  "othertbl",
							ColumnName: "y",
						},
					},
					Operator: &l.Token{TokenType: l.EQUALS, Literal: "="},
				},
			},
		},
		Right:    &TableFrom{TableName: "anotherone"},
		JoinKind: Left,
		Predicate: &Expression{
			Kind: BinaryKind,
			BinaryExpression: &BinaryExpression{
				Left: &Expression{
					Kind: ColumnKind,
					ColumnExpression: &ColumnExpression{
						TableName:  "othertbl",
						ColumnName: "y",
					},
				},
				Right: &Expression{
					Kind: ColumnKind,
					ColumnExpression: &ColumnExpression{
						TableName:  "anotherone",
						ColumnName: "z",
					},
				},
				Operator: &l.Token{TokenType: l.EQUALS, Literal: "="},
			},
		},
	}

	expected := &Statement{
		SelectStatement: &SelectStatement{
			SelectItems: items,
			From:        from,
		},
	}

	fromRes, err := parser.ParseStatement()

	assert.NoError(t, err)
	assert.Equal(t, expected, fromRes)
}
