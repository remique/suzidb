package parser

import (
	l "example.com/suzidb/lexer"
)

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

// Used to parse From clause.
type JoinFrom struct {
	Left      FromType
	Right     FromType
	Kind      JoinKind
	Predicate *Expression
}

// Used to parse From clause.
// TODO: Add alias support. Then we can use separate TableFrom { name string, alias string }.
type FromType struct {
	Join  *JoinFrom
	Table *l.Token
	Kind  FromKind
}
