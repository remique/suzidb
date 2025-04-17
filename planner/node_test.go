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
					Kind:                 parser.IdentifierKind,
					IdentifierExpression: &lexer.Token{TokenType: lexer.IDENTIFIER, Literal: "withoutcol"},
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
