package parser

import (
	l "example.com/suzidb/lexer"
	m "example.com/suzidb/meta"
)

// AstKind enum, describes all possible statements supported by the parser.
type AstKind uint

const (
	SelectKind AstKind = iota
	CreateTableKind
	InsertKind
)

// Describes a Statement. Currently, it is not implemented in a very Go-like way, and
// should use interfaces in the future instead.
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

type SelectStatement2 struct {
	SelectItems *[]Expression
}

// type SelectStatement struct {
// 	SelectItems *[]l.Token
// 	From        *FromType
// }

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
