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

// Describes a Statement. Currently, it is not implemented in a very Go-like way,
// and should use interfaces in the future instead.
type Statement struct {
	SelectStatement      *SelectStatement
	CreateTableStatement *CreateTableStatement
	InsertStatement      *InsertStatement
	Kind                 AstKind
}

// Represents SelectStatement option of Statement. SelectItems is a list of expressions,
// because it supports QualifiedColumn expressions (eg. 'mytable.col') as well as simple Literal
// expressions (eg. 'mytable').
//
// From describes whether select should use single table or join multiple tables.
type SelectStatement struct {
	SelectItems *[]Expression
	From        From
}

// Represents CreateTable statement. This should be refactored pretty soon.
type CreateTableStatement struct {
	TableName  string
	PrimaryKey string
	Columns    *[]m.Column
}

// Represents Insert statement. This should be refactored as well.
type InsertStatement struct {
	TableName string
	// If specific order is specified, push it to this array
	CustomColumns []l.Token

	Values []l.Token
}
