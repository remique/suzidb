package executor

import (
	"encoding/json"
	"fmt"

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

func (e *Executor) ExecutePlan(plan p.Plan) (ExecutionResult, error) {
	switch p := plan.(type) {
	case *p.CreateTablePlan:
		return e.executeCreateTable(*p)
	case *p.InsertPlan:
		return e.executeInsert(*p)
	case *p.SelectPlan:
		return e.executeSelect(*p)
	default:
		return nil, fmt.Errorf("Invalid query")
	}
}

func (e *Executor) executeCreateTable(createTablePlan p.CreateTablePlan) (ExecutionResult, error) {
	err := e.Catalog.CreateTable(createTablePlan.Table.Name, createTablePlan.Table)
	if err != nil {
		return nil, err
	}

	return &CreateTableResult{TableName: createTablePlan.Table.Name}, nil
}

func (e *Executor) executeInsert(insertPlan p.InsertPlan) (ExecutionResult, error) {
	count := 0
	for _, row := range insertPlan.Rows {
		// Check if already exists
		key := fmt.Sprintf("%s:%s", insertPlan.Table.Name, row[insertPlan.Table.PrimaryKey])
		checkIfExists := e.Storage.Get(key)
		if len(checkIfExists) > 0 {
			return nil, fmt.Errorf("UNIQUE constraint failed: %s",
				row[insertPlan.Table.PrimaryKey])
		}

		serializedRow, err := json.Marshal(row)
		if err != nil {
			return nil, err
		}

		err = e.Storage.Set(key, string(serializedRow))
		if err != nil {
			return nil, err
		}

		count++
	}

	// Hardcoded for now, as we can insert only one row
	return &InsertResult{Count: count}, nil
}

func (e *Executor) executeSelect(selectPlan p.SelectPlan) (ExecutionResult, error) {
	executor, err := e.queryExecutorBuilder(selectPlan.Node)
	if err != nil {
		return nil, err
	}

	var rows []m.Row
	for {
		row, err := executor.Next()
		if row == nil {
			break
		}
		if err != nil {
			return nil, err
		}

		rows = append(rows, *row)
	}

	// TODO: Temporary solution until I find a way to merge columns of NestedLoopJoin
	switch n := selectPlan.Node.(type) {
	case *p.NodeScan:
		return &SelectResult{Rows: rows, Columns: n.Table.Columns}, nil
	case *p.NestedLoopJoin:
		return &SelectResult{Rows: rows, Columns: []m.Column{}}, nil
	case *p.NodeProjection:
		return &SelectResult{Rows: rows, Columns: []m.Column{}}, nil
	default:
		return &SelectResult{Rows: rows, Columns: []m.Column{}}, nil
	}
}
