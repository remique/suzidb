package parser

import (
	l "example.com/suzidb/lexer"
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
	//
}

type InsertStatement struct {
	//
}

type ExpressionKind uint

const (
	IdentifierKind ExpressionKind = iota
)

type Expression struct {
	IdentifierExpression *l.Token
	Kind                 ExpressionKind
}
