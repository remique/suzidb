package planner

import (
	m "example.com/suzidb/meta"
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
