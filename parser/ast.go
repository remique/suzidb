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
	SelectItems *[]Expression
	From        FromInterface
}

type TableFrom struct {
	TableName string
}

type FromInterface interface {
	isFrom()
}

type JoinFrom struct {
	Left  FromInterface
	Right FromInterface

	JoinKind  JoinKind
	Predicate *Expression
}

func (tf *TableFrom) isFrom() {}
func (jf *JoinFrom) isFrom()  {}

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
