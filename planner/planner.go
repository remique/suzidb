package planner

import (
	"fmt"

	m "example.com/suzidb/meta"
	p "example.com/suzidb/parser"
	s "example.com/suzidb/storage"
)

// A Plan is created after AST, then is passed to the Executor.
//
// A Plan interface is a "marker trait". For a plan to be a proper Plan
// it must have an empty implementation of Plan() method. Eg.
// ```go
//
//	func (sap *SomeArbitraryPlan) Plan() {}
//
// ```
type Plan interface {
	Plan()
}

type Planner struct {
	Catalog s.Catalog
}

// A Plan to create new Table in the database.
type CreateTablePlan struct {
	Table m.Table
}

type InsertPlan struct {
	Table m.Table
	Rows  []m.Row
}

// Temporary plan, before actual query plan.
type QueryTablePlan struct {
	TableName string
}

func (ctp *CreateTablePlan) Plan() {}
func (ip *InsertPlan) Plan()       {}
func (qtp *QueryTablePlan) Plan()  {}

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
