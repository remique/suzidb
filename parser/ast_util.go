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
