package planner

import (
	"testing"

	"example.com/suzidb/lexer"
	"example.com/suzidb/meta"
	"example.com/suzidb/mocks"
	"example.com/suzidb/parser"

	"github.com/stretchr/testify/assert"
)

func TestBuildNodeScan(t *testing.T) {
	mockCatalog := &mocks.MockCatalog{
		GetTableFunc: func(name string) (*meta.Table, error) {
			return &meta.Table{
				Name:       "mytable",
				PrimaryKey: "id",
				Columns: []meta.Column{
					{Name: "id", Type: meta.IntType, Nullable: false},
					{Name: "name", Type: meta.StringType, Nullable: true},
				},
			}, nil
		},
	}

	nb := NewNodeBuilder(mockCatalog)

	stmt := &parser.Statement{
		SelectStatement: &parser.SelectStatement{
			SelectItems: &[]parser.Expression{
				parser.Expression{
					Kind: parser.ColumnKind,
					ColumnExpression: &parser.ColumnExpression{
						TableName: "withoutcol",
					},
				},
			},
			From: &parser.TableFrom{TableName: "mytable"},
		},
	}

	expected := &NodeScan{
		Table: meta.Table{
			Name:       "mytable",
			PrimaryKey: "id",
			Columns: []meta.Column{
				{Name: "id", Type: meta.IntType, Nullable: false},
				{Name: "name", Type: meta.StringType, Nullable: true},
			},
		},
	}

	res, err := nb.BuildNode(stmt.SelectStatement.From)
	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestBuildJoin(t *testing.T) {
	mockCatalog := &mocks.MockCatalog{
		GetTableFunc: func(name string) (*meta.Table, error) {
			return &meta.Table{
				Name:       name,
				PrimaryKey: "id",
				Columns: []meta.Column{
					{Name: "id", Type: meta.IntType, Nullable: false},
					{Name: "name", Type: meta.StringType, Nullable: true},
				},
			}, nil
		},
	}

	nb := NewNodeBuilder(mockCatalog)

	firstPredicate := &parser.Expression{
		Kind: parser.BinaryKind,
		BinaryExpression: &parser.BinaryExpression{
			Left: &parser.Expression{
				Kind: parser.ColumnKind,
				ColumnExpression: &parser.ColumnExpression{
					TableName:  "sometbl",
					ColumnName: "id",
				},
			},
			Right: &parser.Expression{
				Kind: parser.ColumnKind,
				ColumnExpression: &parser.ColumnExpression{
					TableName:  "othertbl",
					ColumnName: "id",
				},
			},
			Operator: &lexer.Token{TokenType: lexer.EQUALS, Literal: "="},
		},
	}

	secondPredicate := &parser.Expression{
		Kind: parser.BinaryKind,
		BinaryExpression: &parser.BinaryExpression{
			Left: &parser.Expression{
				Kind: parser.ColumnKind,
				ColumnExpression: &parser.ColumnExpression{
					TableName:  "othertbl",
					ColumnName: "id",
				},
			},
			Right: &parser.Expression{
				Kind: parser.ColumnKind,
				ColumnExpression: &parser.ColumnExpression{
					TableName:  "anotherone",
					ColumnName: "id",
				},
			},
			Operator: &lexer.Token{TokenType: lexer.EQUALS, Literal: "="},
		},
	}

	from := &parser.JoinFrom{
		Left: &parser.JoinFrom{
			Left:      &parser.TableFrom{TableName: "sometbl"},
			Right:     &parser.TableFrom{TableName: "othertbl"},
			JoinKind:  parser.Left,
			Predicate: firstPredicate,
		},
		Right:     &parser.TableFrom{TableName: "anotherone"},
		JoinKind:  parser.Left,
		Predicate: secondPredicate,
	}

	expected := &NestedLoopJoin{
		Left: &NestedLoopJoin{
			Left: &NodeScan{
				Table: meta.Table{
					Name:       "sometbl",
					PrimaryKey: "id",
					Columns: []meta.Column{
						{Name: "id", Type: meta.IntType, Nullable: false},
						{Name: "name", Type: meta.StringType, Nullable: true},
					},
				},
			},
			Right: &NodeScan{
				Table: meta.Table{
					Name:       "othertbl",
					PrimaryKey: "id",
					Columns: []meta.Column{
						{Name: "id", Type: meta.IntType, Nullable: false},
						{Name: "name", Type: meta.StringType, Nullable: true},
					},
				},
			},
			Predicate: firstPredicate,
		},
		Right: &NodeScan{
			Table: meta.Table{
				Name:       "anotherone",
				PrimaryKey: "id",
				Columns: []meta.Column{
					{Name: "id", Type: meta.IntType, Nullable: false},
					{Name: "name", Type: meta.StringType, Nullable: true},
				},
			},
		},
		Predicate: secondPredicate,
	}

	res, err := nb.BuildNode(from)
	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}
