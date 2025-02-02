package executor

import (
	"example.com/suzidb/meta"
)

type ExecutionResult interface {
	Result()
}

type CreateTableResult struct {
	TableName string
}

type InsertResult struct {
	Count int
}

type SelectResult struct {
	Rows    []meta.Row
	Columns []meta.Column
}

func (ctr *CreateTableResult) Result() {}
func (ir *InsertResult) Result()       {}
func (sr *SelectResult) Result()       {}
