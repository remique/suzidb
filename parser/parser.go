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

func (p *Parser) currTokenIsJoin() bool {
	switch p.currentToken.TokenType {
	case l.LEFT:
		return true
	case l.INNER:
		return true
	default:
		return false
	}
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

	if len(cols) > 0 && len(cols) != len(*values) {
		return nil, fmt.Errorf("Got %d columns and %d values", len(cols), len(*values))
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

func (p *Parser) parseInsertValues() (*[]Expression, error) {
	var vals []Expression

	for {
		expr, err := p.ParseExpression(LowestPrecedence)
		if err != nil {
			return nil, err
		}

		fmt.Println("Got expr", expr.Kind)

		vals = append(vals, *expr)

		if !p.expectPeekToken(l.COMMA) {
			break
		}

		// Skip ','
		p.nextToken()

		// Fetch next TYPE
		p.nextToken()
	}

	return &vals, nil
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
	exprList, err := p.parseSelectClause()
	if err != nil {
		return nil, fmt.Errorf("Could not parse select clause")
	}

	p.nextToken()

	from, err := p.parseFromClause()
	if err != nil {
		return nil, fmt.Errorf("Could not parse from clause")
	}

	return &Statement{
		SelectStatement: &SelectStatement{
			SelectItems: exprList,
			From:        from,
		},
	}, nil
}

// Parses SELECT clause. It should parse until it finds a parseable expression and
// returns a list of expressions (which are items that are being selected).
//
// In the future it should also support aliases.
func (p *Parser) parseSelectClause() (*[]Expression, error) {
	var items []Expression

	if !p.expectCurrToken(l.SELECT) {
		return nil, fmt.Errorf("Expected SELECT token")
	}
	p.nextToken()

	// Now parse expression until we do not get any more
	for {
		expr, err := p.ParseExpression(LowestPrecedence)
		if err != nil {
			return nil, err
		}

		items = append(items, *expr)

		if !p.expectPeekToken(l.COMMA) {
			break
		} else {
			// Skip last parsed expression
			p.nextToken()

			// Skip comma
			p.nextToken()

		}
	}

	return &items, nil
}

// Parses FROM clause with single table and multiple join support.
func (p *Parser) parseFromClause() (From, error) {
	var left From

	if !p.expectCurrToken(l.FROM) {
		return nil, fmt.Errorf("Expected FROM token")
	}
	p.nextToken()

	// First we need a table that we can parse.
	if !p.expectCurrToken(l.IDENTIFIER) {
		return nil, fmt.Errorf("Expected IDENTIFIER (table name)")
	}

	left = &TableFrom{TableName: p.currentToken.Literal}

	p.nextToken()

	for p.currTokenIsJoin() {
		// Skip "LEFT" or "INNER", etc.
		p.nextToken()
		// Skip "JOIN"
		p.nextToken()

		// Parse tableName
		if !p.expectCurrToken(l.IDENTIFIER) {
			return nil, fmt.Errorf("Expected IDENTIFIER (table name)")
		}

		right := TableFrom{TableName: p.currentToken.Literal}

		//  Skip tableName
		p.nextToken()

		// Skip "ON"
		p.nextToken()

		expr, err := p.ParseExpression(LowestPrecedence)
		if err != nil {
			return nil, err
		}

		left = &JoinFrom{
			Left:      left,
			Right:     &right,
			JoinKind:  Left,
			Predicate: expr,
		}

		p.nextToken()
	}

	return left, nil
}
