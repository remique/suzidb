package planner

import (
	"fmt"

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
