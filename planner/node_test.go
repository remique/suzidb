package planner

// import (
// 	"testing"

// 	"example.com/suzidb/lexer"
// 	"example.com/suzidb/meta"
// 	"example.com/suzidb/mocks"
// 	"example.com/suzidb/parser"

// 	"github.com/stretchr/testify/assert"
// )

// func TestBuildNodeScan(t *testing.T) {
// 	mockCatalog := &mocks.MockCatalog{
// 		GetTableFunc: func(name string) (*meta.Table, error) {
// 			return &meta.Table{
// 				Name:       "mytable",
// 				PrimaryKey: "id",
// 				Columns: []meta.Column{
// 					{Name: "id", Type: meta.IntType, Nullable: false},
// 					{Name: "name", Type: meta.StringType, Nullable: true},
// 				},
// 			}, nil
// 		},
// 	}

// 	nb := NewNodeBuilder(mockCatalog)

// 	statement := parser.Statement{
// 		Kind: parser.SelectKind,
// 		SelectStatement: &parser.SelectStatement{
// 			SelectItems: &[]lexer.Token{{TokenType: lexer.STAR, Literal: "*"}},
// 			From:        &lexer.Token{TokenType: lexer.IDENTIFIER, Literal: "mytable"},
// 		},
// 	}

// 	expected := &NodeScan{
// 		Table: meta.Table{
// 			Name:       "mytable",
// 			PrimaryKey: "id",
// 			Columns: []meta.Column{
// 				{Name: "id", Type: meta.IntType, Nullable: false},
// 				{Name: "name", Type: meta.StringType, Nullable: true},
// 			},
// 		},
// 	}

// 	res, err := nb.buildNodeScan(statement)
// 	assert.NoError(t, err)
// 	assert.Equal(t, expected, res)
// }
