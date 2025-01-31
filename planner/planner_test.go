package planner

import (
	"testing"

	m "example.com/suzidb/meta"
	"example.com/suzidb/mocks"
	p "example.com/suzidb/parser"

	"github.com/stretchr/testify/assert"
)

func TestBuildCreateTable(t *testing.T) {
	mockCatalog := &mocks.MockCatalog{
		GetTableFunc: func(name string) (*m.Table, error) {
			return nil, nil
		},
	}
	planner := NewPlanner(mockCatalog)

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

func TestBuildCreateTableAlreadyExists(t *testing.T) {
	mockCatalog := &mocks.MockCatalog{
		GetTableFunc: func(name string) (*m.Table, error) {
			existingTable := &m.Table{
				Name:       "randomTable",
				Columns:    []m.Column{},
				PrimaryKey: "somePrimary",
			}
			return existingTable, nil
		},
	}
	planner := NewPlanner(mockCatalog)

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

	plan, err := planner.buildCreateTable(stmt)
	assert.Error(t, err)
	assert.Equal(t, nil, plan)
}
