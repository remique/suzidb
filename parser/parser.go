package parser

import (
	l "example.com/suzidb/lexer"
)

type Parser struct {
	lexer        l.Lexer
	currentToken l.Token
	peekToken    l.Token
}

func NewParser(lexer l.Lexer) *Parser {
	currentToken := lexer.NextToken()
	peekToken := lexer.NextToken()

	return &Parser{
		lexer:        lexer,
		currentToken: currentToken,
		peekToken:    peekToken,
	}
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) parseStatement() (*Statement, error) {
	switch p.currentToken.TokenType {
	case "SELECT":
		{
			// parseSelectStatement
		}
	case "INSERT":
		{
			// parseInsertStatement
		}
	case "CREATE":
		{
			// parseInsertStatement
		}
	}

	return nil, nil
}

func (p *Parser) parseSelectStatement() (*Statement, error) {
	//
}

// func (p *Parser) parseExpression() {
// 	// todo
// }
