package meta

type ColumnType int32

const (
	StringType ColumnType = iota
	IntType
)

type Column struct {
	Name string
	Type ColumnType
}

type Table struct {
	Name       string
	Columns    []Column
	PrimaryKey string
}

type Value struct {
	Type ColumnType
}

// Represents a row in the database. When given a key in DB, like `myTable:1` we will
// represent this Row with a map[string]interface{} after it has been deserialized from JSON string.
type Row map[string]interface{}
