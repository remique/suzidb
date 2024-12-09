package planner

import (
	"example.com/suzidb/meta"
)

type Plan interface {
	Plan()
}

type CreateTablePlan struct {
	Table meta.Table
}

type InsertPlan struct{}

func (ctp *CreateTablePlan) Plan() {}
func (ip *InsertPlan) Plan()       {}

type Planner struct{}

// TODO: Once we have AST
// func (p *Planner) build() Plan {
//     switch
// }
