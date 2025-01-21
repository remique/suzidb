package parser

import (
	"fmt"

	l "example.com/suzidb/lexer"
	m "example.com/suzidb/meta"
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
			switch p.peekToken.TokenType {
			case "TABLE":
				return p.parseCreateTableStatement()
			case "INDEX":
				return nil, fmt.Errorf("Currently unsupported")
			default:
				return nil, fmt.Errorf("Expected TABLE or INDEX")
			}
		}
	}

	return nil, nil
}

func (p *Parser) parseCreateTableStatement() (*Statement, error) {
	// Consume 'CREATE' and 'TABLE'
	p.nextToken()
	p.nextToken()

	// Now expect table name identifier
	if !p.expectCurrToken(l.IDENTIFIER) {
		return nil, fmt.Errorf("Expected identifier")
	}

	tableName := p.currentToken.Literal

	if !p.expectPeekToken(l.L_PAREN) {
		return nil, fmt.Errorf("Expected L_PAREN")
	}

	// Now parse column definitions
	columns, pk, err := p.parseCreateTableColumns()
	if err != nil {
		return nil, err
	}

	createTableStmt := CreateTableStatement{
		TableName:  tableName,
		PrimaryKey: *pk,
		Columns:    columns,
	}

	return &Statement{CreateTableStatement: &createTableStmt, Kind: CreateTableKind}, nil
}

func (p *Parser) parseCreateTableColumns() (columns *[]m.Column, primaryKey *string, err error) {
	pk := ""
	pkCount := 0
	result := []m.Column{}

	for {
		var col m.Column
		// If EOF or rParen then break
		if p.expectCurrToken(l.R_PAREN) || p.expectCurrToken(l.EOF) {
			break
		}

		if len(result) > 0 {
			// Expect comma
			if !p.expectCurrToken(l.COMMA) {
				return nil, nil, fmt.Errorf("Expected COMMA")
			}
			p.nextToken()
		}

		if !p.expectCurrToken(l.IDENTIFIER) {
			return nil, nil, fmt.Errorf("Expected Identifier, curr tok: %s", p.currentToken.Literal)
		}

		col.Name = p.currentToken.Literal

		p.nextToken()

		if !(p.expectCurrToken(l.INT_TYPE) || p.expectCurrToken(l.TEXT_TYPE)) {
			return nil, nil, fmt.Errorf("Expected TYPE")
		}

		switch p.currentToken.TokenType {
		case "TEXT_TYPE":
			col.Type = m.StringType
		case "INT_TYPE":
			col.Type = m.IntType
		default:
			return nil, nil, fmt.Errorf("Expected TYPE!")
		}

		p.nextToken()

		if p.expectCurrToken(l.PRIMARY) && p.expectPeekToken(l.KEY) {
			if pkCount >= 1 {
				return nil, nil, fmt.Errorf("Only one PK allowed")
			}
			p.nextToken()
			p.nextToken()
			fmt.Printf("Current tok after PK: %s", p.currentToken.Literal)

			pk = col.Name
			pkCount += 1
		}

		result = append(result, col)
	}

	return &result, &pk, nil
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
