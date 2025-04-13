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

type Precedence uint

const (
	LowestPrecedence Precedence = iota
	DotPrecedence
	EqualsPrecedence
)

// TODO: Refactor this into interfaces.
type Expression struct {
	LiteralExpression         *lexer.Token
	QualifiedColumnExpression *QualifiedColumnExpression
	IdentifierExpression      *lexer.Token
	BinaryExpression          *BinaryExpression
	Kind                      ExpressionKind
}

func tokenToPrecedence(token lexer.Token) Precedence {
	switch token.TokenType {
	case lexer.DOT:
		{
			return DotPrecedence
		}
	case lexer.EQUALS:
		{
			return EqualsPrecedence
		}
	default:
		{
			return LowestPrecedence
		}
	}
}

func (p *Parser) peekPrecedence() Precedence {
	return tokenToPrecedence(p.peekToken)
}

func (p *Parser) currentPrecedence() Precedence {
	return tokenToPrecedence(p.currentToken)
}

func (p *Parser) ParseExpression(precedence Precedence) (*Expression, error) {
	prefix, err := p.parseExpressionAtom()
	if err != nil {
		return nil, fmt.Errorf("ParseExpression err: %s", err.Error())
	}

	for p.peekToken.TokenType != lexer.EOF && precedence < p.peekPrecedence() {
		if prefix != nil {
			p.nextToken()

			var infix *Expression
			switch p.currentToken.TokenType {
			case lexer.DOT:
				{
					res, err := p.parseExpressionColumn2(prefix)
					if err != nil {
						return nil, fmt.Errorf("Err: %s", err.Error())
					}
					infix = res
				}
			case lexer.EQUALS:
				{
					res, err := p.parseExpressionInfixEqual(prefix)
					if err != nil {
						return nil, fmt.Errorf("Err: %s", err.Error())
					}
					infix = res
				}
			default:
				{
					infix = nil
				}
			}

			prefix = infix
		}
	}

	return prefix, nil

}

func (p *Parser) parseExpressionInfixEqual(left *Expression) (*Expression, error) {
	p.nextToken()

	right, err := p.ParseExpression(p.currentPrecedence())
	if err != nil {
		return nil, err
	}

	return &Expression{
		Kind: BinaryKind,
		BinaryExpression: &BinaryExpression{
			Left:     left,
			Right:    right,
			Operator: &lexer.Token{TokenType: lexer.EQUALS, Literal: "="},
		},
	}, nil
}

func (p *Parser) parseExpressionColumn2(left *Expression) (*Expression, error) {
	// Skip DOT
	p.nextToken()

	right, err := p.parseExpressionAtom()
	if err != nil {
		return nil, err
	}

	return &Expression{
		Kind: QualifiedColumnKind,
		QualifiedColumnExpression: &QualifiedColumnExpression{
			TableName:  left,
			ColumnName: right,
		},
	}, nil
}

// NOTE: This is going to be parseBinaryExpression
// func (p *Parser) parseInfixExpression(left *Expression) (*Expression, error) {
// 	currentPrecedence := p.currentPrecedence()
// }

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

func (p *Parser) parseStringExpression() (*Expression, error) {
	token := p.currentToken

	p.nextToken()

	return &Expression{
		Kind:              LiteralKind,
		LiteralExpression: &token,
	}, nil
}

func (p *Parser) parseIdentifierExpression() (*Expression, error) {
	token := p.currentToken

	p.nextToken()

	return &Expression{
		Kind:                 IdentifierKind,
		IdentifierExpression: &token,
	}, nil
}

// NOTE: This is deprecated. Rewrite tests
func (p *Parser) parseExpressionAtom() (*Expression, error) {
	switch p.currentToken.TokenType {
	case lexer.STRING:
		{
			token := p.currentToken

			return &Expression{
				Kind:              LiteralKind,
				LiteralExpression: &token,
			}, nil
		}
	case lexer.IDENTIFIER:
		{
			token := p.currentToken

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
