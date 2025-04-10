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

func (p *Parser) ParseStatement() (*Statement, error) {
	switch p.currentToken.TokenType {
	case "SELECT":
		return p.parseSelectStatement()
	case "INSERT":
		return p.parseInsertStatement()
	case "CREATE":
		switch p.peekToken.TokenType {
		case "TABLE":
			return p.parseCreateTableStatement()
		case "INDEX":
			return nil, fmt.Errorf("Currently unsupported")
		default:
			return nil, fmt.Errorf("Expected TABLE or INDEX")
		}
	}

	return nil, nil
}

func (p *Parser) parseInsertStatement() (*Statement, error) {
	cols := []l.Token{}

	// Consume `INSERT` and `INTO`
	p.nextToken()
	p.nextToken()

	// Now expect table identifier
	if !p.expectCurrToken(l.IDENTIFIER) {
		return nil, fmt.Errorf("Expected identifier")
	}

	tableName := p.currentToken.Literal

	// then optional (...columnList)
	if p.expectPeekToken(l.L_PAREN) {
		p.nextToken()
		p.nextToken()

		// parse columnList
		customCols, err := p.parseInsertColumnList()
		if err != nil {
			return nil, err
		}

		cols = customCols
	} else {
		p.nextToken()
	}

	// Then we parse `VALUES`
	if !p.expectCurrToken(l.VALUES) {
		return nil, fmt.Errorf("Expected VALUES")
	}

	p.nextToken()

	// Then LPAREN

	p.nextToken()

	// Then list of tokens
	values, err := p.parseInsertValues()
	if err != nil {
		return nil, err
	}

	insertStmt := InsertStatement{
		TableName:     tableName,
		CustomColumns: cols,
		Values:        values,
	}

	if len(cols) > 0 && len(cols) != len(values) {
		return nil, fmt.Errorf("Got %d columns and %d values", len(cols), len(values))
	}

	return &Statement{InsertStatement: &insertStmt, Kind: InsertKind}, nil
}

func (p *Parser) parseInsertColumnList() ([]l.Token, error) {
	var columns []l.Token

	for p.expectCurrToken(l.IDENTIFIER) {
		columns = append(columns, p.currentToken)

		// Skip ','
		p.nextToken()

		// Fetch next IDENTIFIER
		p.nextToken()
	}

	return columns, nil
}

func (p *Parser) parseInsertValues() ([]l.Token, error) {
	var vals []l.Token

	for p.expectCurrToken(l.INT) || p.expectCurrToken(l.STRING) {
		vals = append(vals, p.currentToken)

		// Skip ','
		p.nextToken()

		// Fetch next TYPE
		p.nextToken()
	}

	return vals, nil
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

	// Skip '('
	p.nextToken()

	// Go to first item in columns
	p.nextToken()

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

		// Set column to be nullable by default
		col.Nullable = true

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
			fmt.Printf("Current token: %+v, nextToken: %+v", p.currentToken, p.peekToken)
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

			col.Nullable = false

			p.nextToken()
			p.nextToken()

			pk = col.Name
			pkCount += 1
		}

		if p.expectCurrToken(l.NOT) && p.expectPeekToken(l.NULL) {
			p.nextToken()
			p.nextToken()

			col.Nullable = false
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

// TODO: This should support multiple joins
func (p *Parser) parseSelectFrom() (*FromType, error) {
	// SELECT * FROM mytable ...
	// SELECT * FROM mytable LEFT JOIN anothertable ON

	if !p.expectCurrToken(l.FROM) {
		return nil, fmt.Errorf("Expected FROM")
	}

	p.nextToken()

	// Parse table name after FROM
	if !p.expectCurrToken(l.IDENTIFIER) {
		return nil, fmt.Errorf("Expected IDENTIFIER")
	}

	firstTableName := p.currentToken

	p.nextToken()

	// Now we can check if next token is one of the Joins
	// TODO: Parse JoinKind. Now we just assume LEFT
	if !p.expectCurrToken(l.LEFT) {
		return &FromType{
			Table: &firstTableName,
			Kind:  UseTableKind,
		}, nil
	}

	p.nextToken()

	if !p.expectCurrToken(l.JOIN) {
		return nil, fmt.Errorf("Expected JOIN keyword")
	}

	p.nextToken()

	// Now parse another table
	if !p.expectCurrToken(l.IDENTIFIER) {
		return nil, fmt.Errorf("Expected IDENTIFIER")
	}

	secondTableName := p.currentToken

	p.nextToken()

	if !p.expectCurrToken(l.ON) {
		return nil, fmt.Errorf("Expected ON keyword")
	}

	p.nextToken()

	// Now we parse expression
	// TODO: Support ParseExpression()
	expr, err := p.parseBinaryExpression()
	if err != nil {
		return nil, fmt.Errorf("Bad expression")
	}

	return &FromType{
		Kind: UseJoinKind,
		Join: &JoinFrom{
			Left: FromType{
				Table: &firstTableName,
				Kind:  UseTableKind,
			},
			Right: FromType{
				Table: &secondTableName,
				Kind:  UseTableKind,
			},
			Kind:      Left,
			Predicate: expr,
		},
	}, nil
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
