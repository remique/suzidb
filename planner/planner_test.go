package planner

import (
	"testing"

	m "example.com/suzidb/meta"
	p "example.com/suzidb/parser"
	s "example.com/suzidb/storage"

	"github.com/stretchr/testify/assert"
)

func TestBuildCreateTable(t *testing.T) {
	storage := s.NewMemStorage()
	sm := s.NewSchemaManager(storage)
	planner := NewPlanner(sm)

	stmt := p.Statement{
		Kind: p.CreateTableKind,
		CreateTableStatement: &p.CreateTableStatement{
			TableName:  "a",
			PrimaryKey: "b",
			Columns: &[]m.Column{
				{Name: "col1", Type: m.StringType},
				{Name: "col2", Type: m.IntType},
			},
		},
	}

	expected := CreateTablePlan{
		Table: m.Table{
			Name:       "a",
			PrimaryKey: "b",
			Columns: []m.Column{
				{Name: "col1", Type: m.StringType},
				{Name: "col2", Type: m.IntType},
			},
		},
	}

	plan, err := planner.buildCreateTable(stmt)
	assert.Equal(t, plan, &expected, "Expected Plan should be the same")
	assert.NoError(t, err)

}
