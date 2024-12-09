package meta

type Column struct {
	Name    string
	Type    string
	Primary bool
}

type Table struct {
	Name    string
	Columns []Column
}
