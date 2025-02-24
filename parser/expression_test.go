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
