package executor

import (
	"encoding/json"
	"fmt"
	"strings"

	"example.com/suzidb/evaluator"
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
	// Check if already exists
	key := fmt.Sprintf("%s:%s", insertPlan.Table.Name, insertPlan.Row[insertPlan.Table.PrimaryKey])
	checkIfExists := e.Storage.Get(key)
	if len(checkIfExists) > 0 {
		return nil, fmt.Errorf("UNIQUE constraint failed: %s",
			insertPlan.Row[insertPlan.Table.PrimaryKey])
	}

	serializedRow, err := json.Marshal(insertPlan.Row)
	if err != nil {
		return nil, err
	}

	err = e.Storage.Set(key, string(serializedRow))
	if err != nil {
		return nil, err
	}

	// Hardcoded for now, as we can insert only one row
	return &InsertResult{Count: 1}, nil
}

func (e *Executor) executeSelect(selectPlan p.SelectPlan) (ExecutionResult, error) {
	switch internal := selectPlan.Node.(type) {
	case *p.NodeScan:
		return e.executeNodeScan(*internal)
	default:
		return nil, fmt.Errorf("Invalid Node query")
	}
}

func (e *Executor) executeNodeScan(node p.NodeScan) (ExecutionResult, error) {
	var rows []m.Row

	allKeys := e.Storage.ScanKeys()
	for _, key := range allKeys {
		if strings.Contains(key, node.Table.Name) && !strings.Contains(key, "meta") {
			res := e.Storage.Get(key)

			var row m.Row
			err := json.Unmarshal([]byte(res), &row)
			if err != nil {
				return nil, err
			}

			rows = append(rows, row)
		}
	}

	return &SelectResult{Rows: rows, Columns: node.Table.Columns}, nil
}

func (e *Executor) executeNestedLoopJoin(node p.NestedLoopJoin) (ExecutionResult, error) {
	leftAsScan, ok := node.Left.(*p.NodeScan)
	if !ok {
		return nil, fmt.Errorf("expected NodeScan but got %T", node.Left)
	}
	rightAsScan, ok := node.Right.(*p.NodeScan)
	if !ok {
		return nil, fmt.Errorf("expected NodeScan but got %T", node.Left)
	}

	left, err := e.executeNodeScan(*leftAsScan)
	if err != nil {
		return nil, err
	}

	right, err := e.executeNodeScan(*rightAsScan)
	if err != nil {
		return nil, err
	}

	for _, leftItem := range left.(*SelectResult).Rows {
		for _, rightItem := range right.(*SelectResult).Rows {
			// merge rows
			merged := m.MergeRows(
				leftItem,
				rightItem,
				leftAsScan.Table.Name,
				rightAsScan.Table.Name,
			)

			match, err := evaluator.NewEval(node.Predicate).Evaluate(evaluator.WithRow(merged))
			if err != nil {
				return nil, err
			}

			//
		}
	}

	return nil, nil
}
