package parser

import (
	"fmt"

	l "example.com/suzidb/lexer"
)

type Parser struct {
	lexer        l.Lexer
	currentToken l.Token
	peekToken    l.Token
}

// Parser constructor. We need lexer object, because we parse along with lexing.
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
			return p.parseSelectStatement()
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
	// Skip 'SELECT' and then parseSelectItems
	p.nextToken()

	selectItems, err := p.parseSelectItems()
	if err != nil {
		return nil, err
	}

	// parse 'FROM'
	p.nextToken()
	if !p.expectCurrToken(l.FROM) {
		return nil, fmt.Errorf("Expected FROM")
	}

	p.nextToken()
	if !p.expectCurrToken(l.IDENTIFIER) {
		return nil, fmt.Errorf("Expected IDENTIFIER")
	}

	selectStmt := SelectStatement{
		SelectItems: &selectItems,
		From:        &p.currentToken,
	}

	return &Statement{SelectStatement: &selectStmt, Kind: SelectKind}, nil
}

// Method to parse identifiers in Select statement.
// note(remique): This can probably be done better.
func (p *Parser) parseSelectItems() ([]l.Token, error) {
	var items []l.Token

	// Parse first item
	items = append(items, p.currentToken)

	for p.expectPeekToken(l.COMMA) {
		p.nextToken()

		p.nextToken()
		items = append(items, p.currentToken)
	}

	if !(p.expectPeekToken(l.SEMICOLON) || p.expectPeekToken(l.FROM)) {
		return items, fmt.Errorf("Expected either SEMICOLON or FROM")
	}

	return items, nil
}

// func (p *Parser) parseExpression() {
// 	// todo
// }
