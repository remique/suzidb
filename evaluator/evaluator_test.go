package evaluator

import (
	"testing"

	"example.com/suzidb/lexer"
	"example.com/suzidb/parser"

	"github.com/stretchr/testify/assert"
)

func TestEvaluateLiteral(t *testing.T) {
	input := &parser.Expression{
		Kind: parser.LiteralKind,
		LiteralExpression: &lexer.Token{
			TokenType: lexer.STRING,
			Literal:   "helloski",
		},
	}

	expected := &LiteralValue{Value: "helloski"}

	res, err := NewEval(input).Evaluate()
	assert.NoError(t, err)

	assert.Equal(t, expected, res)
}

func TestEvaluateLiteralEqualsTrue(t *testing.T) {
	input := &parser.Expression{
		Kind: parser.BinaryKind,
		BinaryExpression: &parser.BinaryExpression{
			Left: &parser.Expression{
				Kind:              parser.LiteralKind,
				LiteralExpression: &lexer.Token{TokenType: lexer.STRING, Literal: "some"},
			},
			Right: &parser.Expression{
				Kind:              parser.LiteralKind,
				LiteralExpression: &lexer.Token{TokenType: lexer.STRING, Literal: "some"},
			},
			Operator: &lexer.Token{TokenType: lexer.EQUALS, Literal: "="},
		},
	}

	expected := &BooleanValue{Value: true}

	res, err := NewEval(input).Evaluate()
	assert.NoError(t, err)

	assert.Equal(t, expected, res)
}

func TestEvaluateLiteralEqualsFalse(t *testing.T) {
	input := &parser.Expression{
		Kind: parser.BinaryKind,
		BinaryExpression: &parser.BinaryExpression{
			Left: &parser.Expression{
				Kind:              parser.LiteralKind,
				LiteralExpression: &lexer.Token{TokenType: lexer.STRING, Literal: "some"},
			},
			Right: &parser.Expression{
				Kind:              parser.LiteralKind,
				LiteralExpression: &lexer.Token{TokenType: lexer.STRING, Literal: "thing"},
			},
			Operator: &lexer.Token{TokenType: lexer.EQUALS, Literal: "="},
		},
	}

	expected := &BooleanValue{Value: false}

	res, err := NewEval(input).Evaluate()
	assert.NoError(t, err)

	assert.Equal(t, expected, res)
}
