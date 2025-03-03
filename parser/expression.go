package parser

import (
	"example.com/suzidb/lexer"
	"fmt"
)

type ExpressionKind uint

const (
	LiteralKind ExpressionKind = iota
	QualifiedColumnKind
	IdentifierKind
	BinaryKind
)

// A column reference, with optionally qualified with tableName.
// eg. cars.car_id
type QualifiedColumnExpression struct {
	TableName  *Expression
	ColumnName *Expression
}

type BinaryExpression struct {
	Left     *Expression
	Right    *Expression
	Operator *lexer.Token
}

// TODO: Refactor this into interfaces.
type Expression struct {
	LiteralExpression         *lexer.Token
	QualifiedColumnExpression *QualifiedColumnExpression
	IdentifierExpression      *lexer.Token
	BinaryExpression          *BinaryExpression
	Kind                      ExpressionKind
}

func (p *Parser) parseBinaryExpression() (*Expression, error) {
	left, err := p.parseExpressionAtom()
	if err != nil {
		return nil, err
	}

	operator := p.currentToken
	p.nextToken()

	right, err := p.parseExpressionAtom()
	if err != nil {
		return nil, err
	}

	return &Expression{
		Kind: BinaryKind,
		BinaryExpression: &BinaryExpression{
			Left:     left,
			Right:    right,
			Operator: &operator,
		},
	}, nil
}

func (p *Parser) parseExpressionColumn() (*Expression, error) {
	tableName, err := p.parseExpressionAtom()
	if err != nil {
		return nil, err
	}

	if !p.expectCurrToken(lexer.DOT) {
		return nil, fmt.Errorf("Expected '.'")
	}

	// Skip DOT
	p.nextToken()

	columnName, err := p.parseExpressionAtom()
	if err != nil {
		return nil, err
	}

	return &Expression{
		Kind: QualifiedColumnKind,
		QualifiedColumnExpression: &QualifiedColumnExpression{
			TableName:  tableName,
			ColumnName: columnName,
		},
	}, nil
}

func (p *Parser) parseExpressionAtom() (*Expression, error) {
	switch p.currentToken.TokenType {
	case lexer.STRING:
		{
			token := p.currentToken

			p.nextToken()

			return &Expression{
				Kind:              LiteralKind,
				LiteralExpression: &token,
			}, nil
		}
	case lexer.IDENTIFIER:
		{
			token := p.currentToken

			p.nextToken()

			return &Expression{
				Kind:                 IdentifierKind,
				IdentifierExpression: &token,
			}, nil
		}
	default:
		{
			return nil, fmt.Errorf("Invalid expression atom")
		}
	}
}
