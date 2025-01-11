package executor

import (
	"encoding/json"
	"fmt"
	"strings"

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

func (e *Executor) InsertRow(table m.Table, row m.Row) error {
	key := fmt.Sprintf("%s:%s", table.Name, row[table.PrimaryKey])
	serializedRow, err := json.Marshal(row)
	if err != nil {
		return err
	}

	err = e.Storage.Set(key, string(serializedRow))
	if err != nil {
		return err
	}

	return nil
}

// Queries entire table with its contents.
func (e *Executor) QueryTable(tableName string) ([]m.Row, error) {
	var rows []m.Row

	allKeys := e.Storage.ScanKeys()
	for _, key := range allKeys {
		if strings.Contains(key, tableName) && !strings.Contains(key, "meta") {
			res := e.Storage.Get(key)

			var row m.Row
			err := json.Unmarshal([]byte(res), &row)
			if err != nil {
				return rows, err
			}

			rows = append(rows, row)
		}
	}

	return rows, nil
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
		for _, row := range p.Rows {
			err := e.InsertRow(p.Table, row)
			if err != nil {
				return "", err
			}
		}

		return "inserted", nil
	case *p.QueryTablePlan:
		rowsRes, err := e.QueryTable(p.TableName)
		if err != nil {
			return "", err
		}

		rowBytes, err := json.Marshal(rowsRes)
		if err != nil {
			return "", err
		}

		return "Query:" + string(rowBytes), nil
	default:
		return "blabla", nil
	}
}
