package parser

import (
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

func TestParseExpressionWithColumns(t *testing.T) {
	lexer := l.NewLexer("sometable.somecol")
	parser := NewParser(*lexer)

	expected := &Expression{
		Kind: ColumnKind,
		ColumnExpression: &ColumnExpression{
			TableName:  "sometable",
			ColumnName: "somecol",
		},
	}

	res, err := parser.ParseExpression(LowestPrecedence)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestParseExpressionWithColumnsEqual(t *testing.T) {
	lexer := l.NewLexer("sometable.somecol = othertable.othercol")
	parser := NewParser(*lexer)

	expected := &Expression{
		Kind: BinaryKind,
		BinaryExpression: &BinaryExpression{
			Left: &Expression{
				Kind: ColumnKind,
				ColumnExpression: &ColumnExpression{
					TableName:  "sometable",
					ColumnName: "somecol",
				},
			},
			Right: &Expression{
				Kind: ColumnKind,
				ColumnExpression: &ColumnExpression{
					TableName:  "othertable",
					ColumnName: "othercol",
				},
			},
			Operator: &l.Token{TokenType: l.EQUALS, Literal: "="},
		},
	}

	res, err := parser.ParseExpression(LowestPrecedence)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestParseExpressionAllColumns(t *testing.T) {
	lexer := l.NewLexer("*")
	parser := NewParser(*lexer)

	expected := &Expression{
		Kind:                 AllColumnsKind,
		AllColumnsExpression: &AllExpression{},
	}

	res, err := parser.ParseExpression(LowestPrecedence)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}
