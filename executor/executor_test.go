package executor

import (
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
