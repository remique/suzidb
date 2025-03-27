package parser

import (
	"fmt"
	"testing"

	l "example.com/suzidb/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParseExpressionLiteral(t *testing.T) {
	lexer := l.NewLexer("'helloski'")
	parser := NewParser(*lexer)

	expected := &Expression{
		Kind:              LiteralKind,
		LiteralExpression: &l.Token{TokenType: l.STRING, Literal: "helloski"},
	}

	res, err := parser.parseExpressionAtom()
	assert.NoError(t, err)

	assert.Equal(t, expected, res)
}

func TestParseExpressionIdentifier(t *testing.T) {
	lexer := l.NewLexer("someident")
	parser := NewParser(*lexer)

	expected := &Expression{
		Kind:                 IdentifierKind,
		IdentifierExpression: &l.Token{TokenType: l.IDENTIFIER, Literal: "someident"},
	}

	res, err := parser.parseExpressionAtom()
	assert.NoError(t, err)

	assert.Equal(t, expected, res)
}

// func TestParseExpressionQualifiedColumn(t *testing.T) {
// 	lexer := l.NewLexer("sometable.somecol")
// 	parser := NewParser(*lexer)

// 	expected := &Expression{
// 		Kind: QualifiedColumnKind,
// 		QualifiedColumnExpression: &QualifiedColumnExpression{
// 			TableName: &Expression{
// 				Kind:                 IdentifierKind,
// 				IdentifierExpression: &l.Token{TokenType: l.IDENTIFIER, Literal: "sometable"},
// 			},
// 			ColumnName: &Expression{
// 				Kind:                 IdentifierKind,
// 				IdentifierExpression: &l.Token{TokenType: l.IDENTIFIER, Literal: "somecol"},
// 			},
// 		},
// 	}

// 	res, err := parser.parseExpressionColumn()

// 	assert.NoError(t, err)

// 	assert.Equal(t, expected, res)
// }

// func TestParseExpressionBinaryEqual(t *testing.T) {
// 	lexer := l.NewLexer("'some' = 'thing'")
// 	parser := NewParser(*lexer)

// 	expected := &Expression{
// 		Kind: BinaryKind,
// 		BinaryExpression: &BinaryExpression{
// 			Left: &Expression{
// 				Kind:              LiteralKind,
// 				LiteralExpression: &l.Token{TokenType: l.STRING, Literal: "some"},
// 			},
// 			Right: &Expression{
// 				Kind:              LiteralKind,
// 				LiteralExpression: &l.Token{TokenType: l.STRING, Literal: "thing"},
// 			},
// 			Operator: &l.Token{TokenType: l.EQUALS, Literal: "="},
// 		},
// 	}

// 	res, err := parser.parseBinaryExpression()

// 	assert.NoError(t, err)

// 	assert.Equal(t, expected, res)
// }

// func TestParseExpressionIdentifierOnly(t *testing.T) {
// 	lexer := l.NewLexer("sometable")
// 	parser := NewParser(*lexer)

// 	expected := &Expression{
// 		Kind:                 IdentifierKind,
// 		IdentifierExpression: &l.Token{TokenType: l.IDENTIFIER, Literal: "sometable"},
// 	}

// 	res, err := parser.ParseExpression(LowestPrecedence)
// 	assert.NoError(t, err)

// 	assert.Equal(t, expected, res)
// }

// func TestParseExpressionLiteralOnly(t *testing.T) {
// 	lexer := l.NewLexer("'sometable'")
// 	parser := NewParser(*lexer)

// 	expected := &Expression{
// 		Kind:              LiteralKind,
// 		LiteralExpression: &l.Token{TokenType: l.STRING, Literal: "sometable"},
// 	}

// 	res, err := parser.ParseExpression(LowestPrecedence)
// 	assert.NoError(t, err)

// 	assert.Equal(t, expected, res)
// }

func TestParseExpressionWithColumns(t *testing.T) {
	lexer := l.NewLexer("sometable.somecol")
	parser := NewParser(*lexer)

	expected := &Expression{
		Kind: QualifiedColumnKind,
		QualifiedColumnExpression: &QualifiedColumnExpression{
			TableName: &Expression{
				Kind:                 IdentifierKind,
				IdentifierExpression: &l.Token{TokenType: l.IDENTIFIER, Literal: "sometable"},
			},
			ColumnName: &Expression{
				Kind:                 IdentifierKind,
				IdentifierExpression: &l.Token{TokenType: l.IDENTIFIER, Literal: "somecol"},
			},
		},
	}

	res, err := parser.ParseExpression(LowestPrecedence)

	fmt.Println(res)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}
