package planner

import (
	"fmt"

	l "example.com/suzidb/lexer"
	m "example.com/suzidb/meta"
	p "example.com/suzidb/parser"
	s "example.com/suzidb/storage"
)

type Planner struct {
	Catalog s.Catalog
}

func (pl *Planner) Build(statement p.Statement) (Plan, error) {
	switch statement.Kind {
	case p.CreateTableKind:
		return pl.buildCreateTable(statement)
	case p.InsertKind:
		return pl.buildInsert(statement)
	}

	return nil, nil
}

func NewPlanner(c s.Catalog) *Planner {
	return &Planner{Catalog: c}
}

func (pl *Planner) buildCreateTable(statement p.Statement) (Plan, error) {
	tableExists, err := pl.Catalog.GetTable(statement.CreateTableStatement.TableName)
	if err != nil {
		return nil, fmt.Errorf("Error while fetching getTable: %s", err.Error())
	}
	if tableExists != nil {
		return nil, fmt.Errorf("Table already exists")
	}

	table := m.Table{
		Name:       statement.CreateTableStatement.TableName,
		Columns:    *statement.CreateTableStatement.Columns,
		PrimaryKey: statement.CreateTableStatement.PrimaryKey,
	}

	plan := CreateTablePlan{Table: table}

	return &plan, nil
}

func (pl *Planner) buildInsert(stmt p.Statement) (Plan, error) {
	// var row m.Row
	row := make(map[string]interface{})

	// Get table
	tableName := stmt.InsertStatement.TableName
	table, err := pl.Catalog.GetTable(tableName)
	if err != nil {
		return nil, fmt.Errorf("Error while fetching getTable: %s", err.Error())
	}
	if table == nil {
		return nil, fmt.Errorf("Table %s does not exist", tableName)
	}

	if len(stmt.InsertStatement.CustomColumns) != len(stmt.InsertStatement.Values) &&
		len(stmt.InsertStatement.CustomColumns) > 0 {
		return nil, fmt.Errorf("Got %d columns and %d values",
			len(stmt.InsertStatement.CustomColumns), len(stmt.InsertStatement.Values))
	}

	// If len(stmt.CustomColumns) == 0, then we for loop on the values and assign it to row
	// and push the row
	if len(stmt.InsertStatement.CustomColumns) == 0 {
		// NOTE: Or we could not have this if and simply find in customCols the index that we need
		// Then instead of [i] we would assign it. If not found in customCols then we check if
		// The column is nullable. If not, the return error, else all good.
		for i, c := range table.Columns {
			// Check the type
			currTok := stmt.InsertStatement.Values[i]
			if c.Type != tokenToColumnType(currTok) {
				return nil, fmt.Errorf("Expected %d, got %d", c.Type, tokenToColumnType(currTok))
			}

			row[c.Name] = currTok.Literal
		}
	} else {
		for _, c := range table.Columns {
			// Get index of customCols
			idx := getColumnIndex(stmt.InsertStatement.CustomColumns, c.Name)
			currTok := stmt.InsertStatement.Values[idx]
			if c.Type != tokenToColumnType(currTok) {
				return nil, fmt.Errorf("Expected %d, got %d", c.Type, tokenToColumnType(currTok))
			}

			row[c.Name] = currTok.Literal
		}
	}

	return &InsertPlan{Table: *table, Row: row}, nil
}

func getColumnIndex(slice []l.Token, columnName string) int {
	for i := range slice {
		if columnName == slice[i].Literal {
			return i
		}
	}

	return -1
}

func tokenToColumnType(token l.Token) m.ColumnType {
	switch token.TokenType {
	case l.TEXT_TYPE:
		return m.StringType
	case l.INT_TYPE:
		return m.IntType
	default:
		return -1
	}
}
