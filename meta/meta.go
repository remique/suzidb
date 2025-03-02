package meta

type ColumnType int32

const (
	StringType ColumnType = iota
	IntType
)

type Column struct {
	Name     string
	Type     ColumnType
	Nullable bool
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

// Utility function that merges left and right row with prefix (tableName). This is used
// for joins in order to evaluate whether two columns got same values.
// This is not optimal way of doing it for sure, but will suffice for now.
func MergeRows(left, right Row, leftPrefix, rightPrefix string) Row {
	final := make(map[string]interface{})

	for key, value := range left {
		withPrefix := leftPrefix + "." + key
		final[withPrefix] = value
	}

	for key, value := range right {
		withPrefix := rightPrefix + "." + key
		final[withPrefix] = value
	}

	return final
}
