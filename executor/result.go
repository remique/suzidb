package executor

import (
	"encoding/json"
	"example.com/suzidb/meta"
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

func (se *ScanExecutor) Next() (*meta.Row, error) {
	if se.cursor < len(se.Keys) {
		key := se.Keys[se.cursor]

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

type SelectResult struct {
	Rows    []meta.Row
	Columns []meta.Column
}

func (ctr *CreateTableResult) Result() {}
func (ir *InsertResult) Result()       {}
func (sr *SelectResult) Result()       {}
