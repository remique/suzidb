package evaluator

import (
	"testing"

	"example.com/suzidb/lexer"
	"example.com/suzidb/parser"

	"github.com/stretchr/testify/assert"
)

func TestToValue(t *testing.T) {
	cases := []struct {
		og       interface{}
		expected Value
	}{
		{og: 13, expected: &IntValue{Value: 13}},
		{og: "hello", expected: &LiteralValue{Value: "hello"}},
		{og: true, expected: &BooleanValue{Value: true}},
	}

	for _, c := range cases {
		val, err := toValue(c.og)
		assert.NoError(t, err)
		assert.Equal(t, c.expected, val)
	}
}

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

func TestEvalQualifiedColumnWithPrefix(t *testing.T) {
	expr := &parser.Expression{
		Kind: parser.QualifiedColumnKind,
		QualifiedColumnExpression: &parser.QualifiedColumnExpression{
			TableName: &parser.Expression{
				Kind:                 parser.IdentifierKind,
				IdentifierExpression: &lexer.Token{TokenType: lexer.STRING, Literal: "tbl"},
			},
			ColumnName: &parser.Expression{
				Kind:                 parser.IdentifierKind,
				IdentifierExpression: &lexer.Token{TokenType: lexer.STRING, Literal: "col"},
			},
		},
	}
	row := map[string]interface{}{
		"tbl.col":  1,
		"tbl.col2": "hello",
	}

	expected := &IntValue{Value: 1}
	res, err := NewEval(expr).evaluateQualifiedColumn(row, true)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestEvalQualifiedColumnWithoutPrefix(t *testing.T) {
	expr := &parser.Expression{
		Kind: parser.QualifiedColumnKind,
		QualifiedColumnExpression: &parser.QualifiedColumnExpression{
			TableName: &parser.Expression{
				Kind:                 parser.IdentifierKind,
				IdentifierExpression: &lexer.Token{TokenType: lexer.STRING, Literal: "tbl"},
			},
			ColumnName: &parser.Expression{
				Kind:                 parser.IdentifierKind,
				IdentifierExpression: &lexer.Token{TokenType: lexer.STRING, Literal: "col"},
			},
		},
	}
	row := map[string]interface{}{
		"col":  1,
		"col2": "hello",
	}

	expected := &IntValue{Value: 1}
	res, err := NewEval(expr).evaluateQualifiedColumn(row, false)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

// NOTE: Now that is one huge expression. This is definitely a cue to finish up parsing and simply pass in
// a string instead of this monstrosity.
func TestEvalBinaryColumnExpr(t *testing.T) {
	expr := &parser.Expression{
		Kind: parser.BinaryKind,
		BinaryExpression: &parser.BinaryExpression{
			Left: &parser.Expression{
				Kind: parser.QualifiedColumnKind,
				QualifiedColumnExpression: &parser.QualifiedColumnExpression{
					TableName: &parser.Expression{
						Kind:                 parser.IdentifierKind,
						IdentifierExpression: &lexer.Token{TokenType: lexer.STRING, Literal: "tbl1"},
					},
					ColumnName: &parser.Expression{
						Kind:                 parser.IdentifierKind,
						IdentifierExpression: &lexer.Token{TokenType: lexer.STRING, Literal: "col"},
					},
				},
			},
			Right: &parser.Expression{
				Kind: parser.QualifiedColumnKind,
				QualifiedColumnExpression: &parser.QualifiedColumnExpression{
					TableName: &parser.Expression{
						Kind:                 parser.IdentifierKind,
						IdentifierExpression: &lexer.Token{TokenType: lexer.STRING, Literal: "tbl2"},
					},
					ColumnName: &parser.Expression{
						Kind:                 parser.IdentifierKind,
						IdentifierExpression: &lexer.Token{TokenType: lexer.STRING, Literal: "col"},
					},
				},
			},
			Operator: &lexer.Token{TokenType: lexer.EQUALS, Literal: "="},
		},
	}

	row := map[string]interface{}{
		"tbl1.col":  "hello",
		"tbl1.col2": 1,
		"tbl2.col":  "hello",
		"tbl2.col2": 2,
	}

	expected := &BooleanValue{Value: true}
	res, err := NewEval(expr).Evaluate(WithRow(row))

	assert.NoError(t, err)
	assert.Equal(t, expected, res)

}
