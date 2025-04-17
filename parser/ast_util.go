package parser

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
