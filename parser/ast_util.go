package parser

type FromKind uint

const (
	UseTableKind FromKind = iota
	UseJoinKind
)

// Represents currently "supported" joins. As of now, because of
// rapid development state of this database, only Left joins are
// truly supported and possible to query.
type JoinKind uint

const (
	Left JoinKind = iota
	Right
	Inner
)

// Represents a single Table query.
type TableFrom struct {
	TableName string
}

type From interface {
	isFrom()
}

// Represents From when encountering Joins. It uses left-associative
// tree to parse multiple joins. Inside Left and Right you can have
// either another JoinFrom or simply TableFrom.
type JoinFrom struct {
	Left  From
	Right From

	JoinKind  JoinKind
	Predicate *Expression
}

func (tf *TableFrom) isFrom() {}
func (jf *JoinFrom) isFrom()  {}
