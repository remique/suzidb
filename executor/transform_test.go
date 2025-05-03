package executor

// import (
// 	"testing"

// 	"example.com/suzidb/lexer"
// 	"example.com/suzidb/meta"
// 	"example.com/suzidb/parser"

// 	"github.com/stretchr/testify/assert"
// )

// func TestTransformProjectLiteralOnly(t *testing.T) {
// 	input := []meta.Row{
// 		{
// 			"id":      "1",
// 			"name":    "alice",
// 			"surname": "kowalski",
// 		},
// 		{
// 			"id":      "2",
// 			"name":    "bob",
// 			"surname": "testowsky",
// 		},
// 	}

// 	expected := []meta.Row{
// 		{
// 			"id": "1",
// 		},
// 		{
// 			"id": "2",
// 		},
// 	}

// 	exprList := &[]parser.Expression{
// 		parser.Expression{
// 			Kind:                 parser.IdentifierKind,
// 			IdentifierExpression: &lexer.Token{TokenType: lexer.IDENTIFIER, Literal: "id"},
// 		},
// 	}

// 	res, err := NewTransformer().Project(input, exprList)
// 	assert.NoError(t, err)
// 	assert.Equal(t, expected, res)
// }

// func TestTransformProjectQualifiedColumns(t *testing.T) {
// 	input := []meta.Row{
// 		{
// 			"sometbl.id":      "1",
// 			"sometbl.name":    "alice",
// 			"sometbl.surname": "kowalski",
// 		},
// 		{
// 			"sometbl.id":      "2",
// 			"sometbl.name":    "bob",
// 			"sometbl.surname": "testowsky",
// 		},
// 	}

// 	expected := []meta.Row{
// 		{
// 			"sometbl.id": "1",
// 		},
// 		{
// 			"sometbl.id": "2",
// 		},
// 	}

// 	exprList := &[]parser.Expression{
// 		parser.Expression{
// 			Kind: parser.ColumnKind,
// 			ColumnExpression: &parser.ColumnExpression{
// 				TableName:  "sometbl",
// 				ColumnName: "id",
// 			},
// 		},
// 	}

// 	res, err := NewTransformer().Project(input, exprList)
// 	assert.NoError(t, err)
// 	assert.Equal(t, expected, res)
// }
