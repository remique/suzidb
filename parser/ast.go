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
	QualifiedColumnKind
)

// Tablename.ColumnName
type QualifiedColumnExpression struct {
	tableName  string
	columnName string
}

type Expression struct {
	IdentifierExpression      *l.Token
	QualifiedColumnExpression *QualifiedColumnExpression
	Kind                      ExpressionKind
}

type FromKind uint

const (
	UseTableKind FromKind = iota
	UseJoinKind
)

type JoinKind uint

const (
	Left JoinKind = iota
	Right
	Inner
)

type JoinType struct {
	Left  FromType
	Right FromType
	Kind  JoinKind
}

type FromType struct {
	Join  *JoinType
	Table *l.Token
	Kind  FromKind
}

// TODO: Once we support lexing JOINS we can use the following structure to parse them.
// Since for now it would break tests I commented it out.
// From can be either be a single Table or a Join
//
// type SelectStatement struct {
// 	SelectItems *[]l.Token
// 	From        FromType
// }
