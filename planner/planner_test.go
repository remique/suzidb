package planner

import (
	"testing"

	l "example.com/suzidb/lexer"
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

func TestBuildInsertPlan(t *testing.T) {
	mockCatalog := &mocks.MockCatalog{
		GetTableFunc: func(name string) (*m.Table, error) {
			return &m.Table{
				Name:       "mytable",
				PrimaryKey: "id",
				Columns: []m.Column{
					{Name: "id", Type: m.IntType},
					{Name: "name", Type: m.StringType},
				},
			}, nil
		},
	}
	planner := NewPlanner(mockCatalog)

	stmt := p.Statement{
		Kind: p.InsertKind,
		InsertStatement: &p.InsertStatement{
			TableName: "mytable",
			Values:    []l.Token{l.NewToken(l.INT_TYPE, "10"), l.NewToken(l.TEXT_TYPE, "john")},
		},
	}

	expected := InsertPlan{
		Table: m.Table{
			Name:       "mytable",
			PrimaryKey: "id",
			Columns: []m.Column{
				{Name: "id", Type: m.IntType},
				{Name: "name", Type: m.StringType},
			},
		},
		Row: m.Row{"id": "10", "name": "john"},
	}

	plan, err := planner.buildInsert(stmt)
	assert.Equal(t, plan, &expected, "Expected Plan should be the same")
	assert.NoError(t, err)
}

func TestBuildInsertPlanCustomCols(t *testing.T) {
	mockCatalog := &mocks.MockCatalog{
		GetTableFunc: func(name string) (*m.Table, error) {
			return &m.Table{
				Name:       "mytable",
				PrimaryKey: "id",
				Columns: []m.Column{
					{Name: "id", Type: m.IntType},
					{Name: "name", Type: m.StringType},
				},
			}, nil
		},
	}
	planner := NewPlanner(mockCatalog)

	stmt := p.Statement{
		Kind: p.InsertKind,
		InsertStatement: &p.InsertStatement{
			TableName:     "mytable",
			CustomColumns: []l.Token{l.NewToken(l.IDENTIFIER, "name"), l.NewToken(l.IDENTIFIER, "id")},
			Values:        []l.Token{l.NewToken(l.TEXT_TYPE, "john"), l.NewToken(l.INT_TYPE, "10")},
		},
	}

	expected := InsertPlan{
		Table: m.Table{
			Name:       "mytable",
			PrimaryKey: "id",
			Columns: []m.Column{
				{Name: "id", Type: m.IntType},
				{Name: "name", Type: m.StringType},
			},
		},
		Row: m.Row{"id": "10", "name": "john"},
	}

	plan, err := planner.buildInsert(stmt)
	assert.Equal(t, plan, &expected, "Expected Plan should be the same")
	assert.NoError(t, err)
}

func TestBuildInsertPlanCustomColsWithNullable(t *testing.T) {
	mockCatalog := &mocks.MockCatalog{
		GetTableFunc: func(name string) (*m.Table, error) {
			return &m.Table{
				Name:       "mytable",
				PrimaryKey: "id",
				Columns: []m.Column{
					{Name: "id", Type: m.IntType, Nullable: false},
					{Name: "name", Type: m.StringType, Nullable: true},
				},
			}, nil
		},
	}
	planner := NewPlanner(mockCatalog)

	stmt := p.Statement{
		Kind: p.InsertKind,
		InsertStatement: &p.InsertStatement{
			TableName:     "mytable",
			CustomColumns: []l.Token{l.NewToken(l.IDENTIFIER, "id")},
			Values:        []l.Token{l.NewToken(l.INT_TYPE, "10")},
		},
	}

	expected := InsertPlan{
		Table: m.Table{
			Name:       "mytable",
			PrimaryKey: "id",
			Columns: []m.Column{
				{Name: "id", Type: m.IntType, Nullable: false},
				{Name: "name", Type: m.StringType, Nullable: true},
			},
		},
		Row: m.Row{"id": "10", "name": ""},
	}

	plan, err := planner.buildInsert(stmt)
	assert.Equal(t, plan, &expected, "Expected Plan should be the same")
	assert.NoError(t, err)
}

func TestBuildSelectAllCols(t *testing.T) {
	mockCatalog := &mocks.MockCatalog{
		GetTableFunc: func(name string) (*m.Table, error) {
			return &m.Table{
				Name:       "mytable",
				PrimaryKey: "id",
				Columns: []m.Column{
					{Name: "id", Type: m.IntType, Nullable: false},
					{Name: "name", Type: m.StringType, Nullable: true},
				},
			}, nil
		},
	}
	planner := NewPlanner(mockCatalog)

	stmt := p.Statement{
		Kind: p.SelectKind,
		SelectStatement: &p.SelectStatement{
			SelectItems: &[]p.Expression{
				p.Expression{
					Kind:                 p.AllColumnsKind,
					AllColumnsExpression: &p.AllExpression{},
				},
			},
			From: &p.TableFrom{TableName: "mytable"},
		},
	}

	expected := SelectPlan{
		Node: &NodeScan{
			Table: m.Table{
				Name:       "mytable",
				PrimaryKey: "id",
				Columns: []m.Column{
					{Name: "id", Type: m.IntType, Nullable: false},
					{Name: "name", Type: m.StringType, Nullable: true},
				},
			},
		},
	}

	plan, err := planner.buildSelect(stmt)
	assert.Equal(t, &expected, plan, "Expected Plan should be the same")
	assert.NoError(t, err)
}

func TestBuildSelectProject(t *testing.T) {
	mockCatalog := &mocks.MockCatalog{
		GetTableFunc: func(name string) (*m.Table, error) {
			return &m.Table{
				Name:       "mytable",
				PrimaryKey: "id",
				Columns: []m.Column{
					{Name: "id", Type: m.IntType, Nullable: false},
					{Name: "name", Type: m.StringType, Nullable: true},
				},
			}, nil
		},
	}
	planner := NewPlanner(mockCatalog)

	stmt := p.Statement{
		Kind: p.SelectKind,
		SelectStatement: &p.SelectStatement{
			SelectItems: &[]p.Expression{
				p.Expression{
					Kind: p.ColumnKind,
					ColumnExpression: &p.ColumnExpression{
						TableName: "name",
					},
				},
			},
			From: &p.TableFrom{TableName: "mytable"},
		},
	}

	expected := SelectPlan{
		Node: &NodeProjection{
			Source: &NodeScan{
				Table: m.Table{
					Name:       "mytable",
					PrimaryKey: "id",
					Columns: []m.Column{
						{Name: "id", Type: m.IntType, Nullable: false},
						{Name: "name", Type: m.StringType, Nullable: true},
					},
				},
			},
			Expressions: &[]p.Expression{
				p.Expression{
					Kind: p.ColumnKind,
					ColumnExpression: &p.ColumnExpression{
						TableName: "name",
					},
				},
			},
		},
	}

	plan, err := planner.buildSelect(stmt)
	assert.Equal(t, &expected, plan, "Expected Plan should be the same")
	assert.NoError(t, err)
}
