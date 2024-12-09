package executor

import (
	m "example.com/suzidb/meta"
	p "example.com/suzidb/planner"
	s "example.com/suzidb/storage"
)

type Executor struct {
	Storage s.Storage
}

func NewExecutor(storage s.Storage) *Executor {
	return &Executor{Storage: storage}
}

func (e *Executor) CreateTable(table m.Table) {
	e.Storage.Set(table.Name, "x")
}

func (e *Executor) ExecutePlan(plan p.Plan) string {
	switch p := plan.(type) {
	case *p.CreateTablePlan:
		e.CreateTable(p.Table)

		return "create"
	case *p.InsertPlan:
		return "insert"
	default:
		return "blabla"
	}
}
