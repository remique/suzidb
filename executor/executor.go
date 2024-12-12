package executor

import (
	m "example.com/suzidb/meta"
	p "example.com/suzidb/planner"
	s "example.com/suzidb/storage"
)

type Executor struct {
	Storage s.Storage
	Catalog s.Catalog
}

func NewExecutor(storage s.Storage, catalog s.Catalog) *Executor {
	return &Executor{Storage: storage, Catalog: catalog}
}

func (e *Executor) CreateTable(table m.Table) error {
	return e.Catalog.CreateTable(table.Name, table)
}

func (e *Executor) GetTable(name string) (*m.Table, error) {
	return e.Catalog.GetTable(name)
}

func (e *Executor) ExecutePlan(plan p.Plan) (string, error) {
	switch p := plan.(type) {
	case *p.CreateTablePlan:
		err := e.CreateTable(p.Table)
		if err != nil {
			return "", err
		}

		return "created", nil
	case *p.InsertPlan:
		return "insert", nil
	default:
		return "blabla", nil
	}
}
