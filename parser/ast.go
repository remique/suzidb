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

type InsertStatement struct {
	TableName string
	// If specific order is specified, push it to this array
	CustomColumns []l.Token

	Values []l.Token
}

type ExpressionKind uint

const (
	IdentifierKind ExpressionKind = iota
)

type Expression struct {
	IdentifierExpression *l.Token
	Kind                 ExpressionKind
}
