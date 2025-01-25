package parser

import (
	l "example.com/suzidb/lexer"
	m "example.com/suzidb/meta"
)

type AstKind uint

const (
	SelectKind AstKind = iota
	CreateTableKind
	InsertKind
)

type Statement struct {
	SelectStatement      *SelectStatement
	CreateTableStatement *CreateTableStatement
	InsertStatement      *InsertStatement
	Kind                 AstKind
}

type SelectStatement struct {
	SelectItems *[]l.Token
	From        *l.Token
}

type CreateTableStatement struct {
	TableName  string
	PrimaryKey string
	Columns    *[]m.Column
}

type InsertValue struct {
	value      l.Token
	columnName string
}

type InsertStatement struct {
	TableName string
	// If specific order is specified, push it to this array
	customColumns []l.Token

	values []InsertValue
}

/*
Insert statement needs to have an array of structure
struct InsertValue {
	value Token
	columnName string
}

If the columnList is empty, then we assume original
sequence.

I guess we would have to change the plan so that
it does not have []m.Row but []InsertValue

and then in executor
row := m.Row
for value := range insertValues {
	// Check for type
	row[value.columnName] = value.value
}

And then iterate over columns, find if if has a key
if it does not, check if column is nullable.

*/

type ExpressionKind uint

const (
	IdentifierKind ExpressionKind = iota
)

type Expression struct {
	IdentifierExpression *l.Token
	Kind                 ExpressionKind
}
