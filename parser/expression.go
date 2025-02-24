package parser

import (
	"example.com/suzidb/lexer"
	"fmt"
)

type ExpressionKind uint

const (
	LiteralKind ExpressionKind = iota
	QualifiedColumnKind
)

// A column reference, with optionally qualified with tableName.
// eg. cars.car_id
type QualifiedColumnExpression struct {
	tableName  string
	columnName string
}

type Expression struct {
	LiteralExpression         *lexer.Token
	QualifiedColumnExpression *QualifiedColumnExpression
	Kind                      ExpressionKind
}

func (p *Parser) parseExpression() (*Expression, error) {
	// switch p.currentToken.TokenType {
	// case l.IDENTIFIER:
	// 	{
	// 		return p.parseColumnExpression()
	// 	}
	// default:
	// 	{
	// 		return nil, fmt.Errorf("Unsupported expression")
	// 	}
	// }

	// Start with easier expressions
	return nil, nil
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
	default:
		{
			return nil, fmt.Errorf("Invalid expression atom")
		}
	}
}

// func (p *Parser) parseColumnExpression() (*Expression, error) {
// 	tableName := p.currentToken.Literal

// 	if !p.expectPeekToken(l.DOT) {
// 		return nil, fmt.Errorf("Expected .")
// 	}

// 	// Skip dot
// 	p.nextToken()

// 	p.nextToken()

// 	if !p.expectCurrToken(l.IDENTIFIER) {
// 		return nil, fmt.Errorf("Expected identifier")
// 	}

// 	columnName := p.currentToken.Literal

// 	return &Expression{
// 		QualifiedColumnExpression: &QualifiedColumnExpression{
// 			tableName:  tableName,
// 			columnName: columnName,
// 		},
// 		Kind: QualifiedColumnKind,
// 	}, nil
// }
