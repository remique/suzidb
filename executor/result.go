package executor

import (
	"encoding/json"
	"example.com/suzidb/evaluator"
	"example.com/suzidb/meta"
	"example.com/suzidb/parser"
	"example.com/suzidb/planner"
	"example.com/suzidb/storage"
	"fmt"
	"strings"
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

type QueryExecutor interface {
	Next() (*meta.Row, error)
}

type ScanExecutor struct {
	Storage storage.Storage
	Catalog storage.Catalog

	// See comments under NewScanExecutor why it is not optimal
	Keys  []string
	Table meta.Table

	cursor int
}

type NestedLoopJoinExecutor struct {
	Storage storage.Storage
	Catalog storage.Catalog

	Left  QueryExecutor
	Right QueryExecutor

	Predicate *parser.Expression

	cursor int
}

func NewScanExecutor(s storage.Storage, c storage.Catalog, table meta.Table) *ScanExecutor {
	// This does not really represent Iterator model, as we have to eagerly fetch
	// all the keys from the Storage. It should have `NextKey()` method that we can
	// fetch as we go through the iterator. For now I am leaving it for the sake of
	// finishing initial implementation.
	var keys []string

	allKeys := s.ScanKeys()
	for _, key := range allKeys {
		if strings.Contains(key, table.Name) && !strings.Contains(key, "meta") {
			keys = append(keys, key)
		}
	}

	return &ScanExecutor{
		Storage: s,
		Catalog: c,
		Keys:    keys,
		Table:   table,
		cursor:  0,
	}
}

func NewNestedLoopJoinExecutor(s storage.Storage, c storage.Catalog,
	left, right QueryExecutor, predicate *parser.Expression) *NestedLoopJoinExecutor {
	return &NestedLoopJoinExecutor{
		Storage:   s,
		Catalog:   c,
		cursor:    0,
		Left:      left,
		Right:     right,
		Predicate: predicate,
	}
}

func (nlj *NestedLoopJoinExecutor) Next() (*meta.Row, error) {
	leftRes, err := nlj.Left.Next()
	if err != nil {
		return nil, err
	}
	rightRes, err := nlj.Right.Next()
	if err != nil {
		return nil, err
	}

	if leftRes != nil {
		// now merge rows and evaluate
		merged := meta.MergeRows(
			meta.WithMergeRow(*leftRes, nlj.Left.(*ScanExecutor).Table.Name),
			meta.WithMergeRow(*rightRes, nlj.Right.(*ScanExecutor).Table.Name),
		)

		match, err := evaluator.NewEval(nlj.Predicate).Evaluate(evaluator.WithRow(merged))
		if err != nil {
			return nil, err
		}

		matchAsNativeVal, err := evaluator.ValueToNative(match)
		if err != nil {
			return nil, err
		}

		if matchAsNativeVal == true {
			return &merged, nil
		}
	}

	return nil, nil
}

func (se *ScanExecutor) Next() (*meta.Row, error) {
	if se.cursor < len(se.Keys) {
		key := se.Keys[se.cursor]
		se.cursor++

		value := se.Storage.Get(key)
		var row meta.Row
		err := json.Unmarshal([]byte(value), &row)
		if err != nil {
			return nil, fmt.Errorf("Error while unmarshalling: %s", err.Error())
		}

		return &row, nil
	}

	return nil, fmt.Errorf("Cursor out of bounds")
}

func (e *Executor) queryExecutorBuilder(node planner.NodeQuery) (QueryExecutor, error) {
	switch n := node.(type) {
	case *planner.NodeScan:
		{
			return NewScanExecutor(e.Storage, e.Catalog, n.Table), nil
		}
	case *planner.NestedLoopJoin:
		{
			left, err := e.queryExecutorBuilder(n.Left)
			if err != nil {
				return nil, err
			}
			right, err := e.queryExecutorBuilder(n.Right)
			if err != nil {
				return nil, err
			}

			return NewNestedLoopJoinExecutor(e.Storage, e.Catalog, left, right, n.Predicate), nil
		}
	default:
		return nil, fmt.Errorf("Unsupported query")
	}

}

type SelectResult struct {
	Rows    []meta.Row
	Columns []meta.Column
}

func (ctr *CreateTableResult) Result() {}
func (ir *InsertResult) Result()       {}
func (sr *SelectResult) Result()       {}
