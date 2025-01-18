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

func (p *Parser) expectCurrToken(expectKind l.TokenType) bool {
	if p.currentToken.TokenType == expectKind {
		return true
	}

	return false
}

func (p *Parser) expectPeekToken(expectKind l.TokenType) bool {
	if p.peekToken.TokenType == expectKind {
		return true
	}

	return false
}

func (p *Parser) parseStatement() (*Statement, error) {
	switch p.currentToken.TokenType {
	case "SELECT":
		{
			// return p.parseSelectStatement()
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

func (p *Parser) parseExpression() {
	//
}

// func (p *Parser) parseSelectStatement() (*Statement, error) {
// 	// Skip 'SELECT' and then parseSelectItems
// }

func (p *Parser) parseSelectItems() ([]l.Token, error) {
	var items []l.Token

	// Parse while we peek 'FROM' or semicolon
	// TODO: Expect commas
	for !(p.expectCurrToken(l.SEMICOLON) || p.expectCurrToken(l.WHERE)) {
		if p.expectCurrToken(l.IDENTIFIER) {
			items = append(items, p.currentToken)
		}

		p.nextToken()
	}

	return items, nil
}

// func (p *Parser) parseExpression() {
// 	// todo
// }
