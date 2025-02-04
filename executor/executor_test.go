package executor

import (
	"fmt"
	"testing"

	"example.com/suzidb/meta"
	"example.com/suzidb/planner"
	"example.com/suzidb/storage"

	"github.com/stretchr/testify/assert"
)

func TestExecuteCreateTable(t *testing.T) {
	s := storage.NewMemStorage()
	sm := storage.NewSchemaManager(s)
	e := NewExecutor(s, sm)

	plan := planner.CreateTablePlan{
		Table: meta.Table{
			Name:       "a",
			PrimaryKey: "b",
			Columns: []meta.Column{
				{Name: "col1", Type: meta.StringType},
				{Name: "col2", Type: meta.IntType},
			},
		},
	}

	expected := &CreateTableResult{TableName: plan.Table.Name}

	res, err := e.executeCreateTable(plan)
	assert.NoError(t, err)
	assert.Equal(t, expected, res)

	// Assert that storage key has been properly saved. This should be tested separetely.
	storageRes := s.Get("meta:" + plan.Table.Name)
	assert.Greater(t, len(storageRes), 0)
}

func TestExecuteInsertPlan(t *testing.T) {
	s := storage.NewMemStorage()
	sm := storage.NewSchemaManager(s)
	e := NewExecutor(s, sm)

	plan := planner.InsertPlan{
		Table: meta.Table{
			Name:       "mytable",
			PrimaryKey: "id",
			Columns: []meta.Column{
				{Name: "id", Type: meta.IntType},
				{Name: "name", Type: meta.StringType},
			},
		},
		Row: meta.Row{"id": "10", "name": "john"},
	}

	expected := &InsertResult{Count: 1}

	res, err := e.executeInsert(plan)
	assert.NoError(t, err)
	assert.Equal(t, expected, res)

	// Assert that storage key has been properly saved. This should be tested separetely.
	key := fmt.Sprintf("%s:%s", plan.Table.Name, plan.Row[plan.Table.PrimaryKey])
	storageRes := s.Get(key)
	assert.Greater(t, len(storageRes), 0)
}

func TestExecuteInsertPlanPKAlreadyExist(t *testing.T) {
	s := storage.NewMemStorage()
	sm := storage.NewSchemaManager(s)
	e := NewExecutor(s, sm)

	plan := planner.InsertPlan{
		Table: meta.Table{
			Name:       "mytable",
			PrimaryKey: "id",
			Columns: []meta.Column{
				{Name: "id", Type: meta.IntType},
				{Name: "name", Type: meta.StringType},
			},
		},
		Row: meta.Row{"id": "10", "name": "john"},
	}

	expected := &InsertResult{Count: 1}

	res, err := e.executeInsert(plan)
	assert.NoError(t, err)
	assert.Equal(t, expected, res)

	// And saving the same plan again should result in an error
	res2, err := e.executeInsert(plan)
	assert.Error(t, err)
	assert.Equal(t, nil, res2)
}
