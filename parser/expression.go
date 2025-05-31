package parser

import (
	"example.com/suzidb/lexer"
	"fmt"
)

type ExpressionKind uint

const (
	LiteralKind ExpressionKind = iota
	ColumnKind
	BinaryKind
	ConstExprKind
	AllColumnsKind
)

// A column reference, with optionally qualified with tableName.
// eg. cars.car_id
type ColumnExpression struct {
	TableName  string
	ColumnName string
}

type BinaryExpression struct {
	Left     *Expression
	Right    *Expression
	Operator *lexer.Token
}

// This structure does not contain anything, as it is only a marker
// that we want to query all columns returned.
type AllExpression struct{}

type Precedence uint

const (
	LowestPrecedence Precedence = iota
	DotPrecedence
	EqualsPrecedence
)

// TODO: Refactor this into interfaces.
type Expression struct {
	LiteralExpression    *lexer.Token
	ConstExpression      *ConstExpression
	AllColumnsExpression *AllExpression
	ColumnExpression     *ColumnExpression
	BinaryExpression     *BinaryExpression
	Kind                 ExpressionKind
}

type ConstKind uint

const (
	IntKind ConstKind = iota
	StringKind
	NullKind
)

type ConstExpression struct {
	Int  *lexer.Token
	Null *lexer.Token
	Kind ConstKind
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

func (p *Parser) parseExpressionAtom() (*Expression, error) {
	switch p.currentToken.TokenType {
	case lexer.NULL:
		{
			token := p.currentToken

			return &Expression{
				Kind: ConstExprKind,
				ConstExpression: &ConstExpression{
					Null: &token,
					Kind: NullKind,
				},
			}, nil
		}
	case lexer.INT:
		{
			token := p.currentToken

			return &Expression{
				Kind: ConstExprKind,
				ConstExpression: &ConstExpression{
					Int:  &token,
					Kind: IntKind,
				},
			}, nil
		}
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

			if p.expectPeekToken(lexer.DOT) {
				p.nextToken()
				p.nextToken()

				return &Expression{
					Kind: ColumnKind,
					ColumnExpression: &ColumnExpression{
						TableName:  token.Literal,
						ColumnName: p.currentToken.Literal,
					},
				}, nil
			} else {
				return &Expression{
					Kind: ColumnKind,
					ColumnExpression: &ColumnExpression{
						TableName: token.Literal,
					},
				}, nil

			}

		}
	case lexer.STAR:
		{
			return &Expression{
				Kind:                 AllColumnsKind,
				AllColumnsExpression: &AllExpression{},
			}, nil
		}
	default:
		{
			return nil, fmt.Errorf("Invalid expression atom")
		}
	}
}
