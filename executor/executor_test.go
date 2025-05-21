package executor

import (
	"fmt"
	"testing"

	"example.com/suzidb/lexer"
	"example.com/suzidb/meta"
	"example.com/suzidb/parser"
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

// This test is unecessarily long, in the future change it to plain SQL query or prepare
// mock storage.
func TestExecuteSelectNodeScan(t *testing.T) {
	s := storage.NewMemStorage()
	sm := storage.NewSchemaManager(s)
	e := NewExecutor(s, sm)

	createPlan := planner.CreateTablePlan{
		Table: meta.Table{
			Name:       "a",
			PrimaryKey: "b",
			Columns: []meta.Column{
				{Name: "col1", Type: meta.StringType},
				{Name: "col2", Type: meta.IntType},
			},
		},
	}
	insert1plan := planner.InsertPlan{
		Table: meta.Table{
			Name:       "mytable",
			PrimaryKey: "id",
			Columns: []meta.Column{
				{Name: "id", Type: meta.IntType},
				{Name: "name", Type: meta.StringType},
			},
		},
		Row: meta.Row{"id": "1", "name": "john"},
	}
	insert2plan := planner.InsertPlan{
		Table: meta.Table{
			Name:       "mytable",
			PrimaryKey: "id",
			Columns: []meta.Column{
				{Name: "id", Type: meta.IntType},
				{Name: "name", Type: meta.StringType},
			},
		},
		Row: meta.Row{"id": "2", "name": "emily"},
	}
	selectPlan := planner.SelectPlan{
		Node: &planner.NodeScan{
			Table: meta.Table{
				Name:       "a",
				PrimaryKey: "b",
				Columns: []meta.Column{
					{Name: "col1", Type: meta.StringType},
					{Name: "col2", Type: meta.IntType},
				},
			},
		},
	}

	_, err := e.executeCreateTable(createPlan)
	assert.NoError(t, err)

	_, err = e.executeInsert(insert1plan)
	assert.NoError(t, err)

	_, err = e.executeInsert(insert2plan)
	assert.NoError(t, err)

	expected := &SelectResult{
		Rows: []meta.Row{
			meta.Row{"id": "1", "name": "john"},
			meta.Row{"id": "2", "name": "emily"},
		},
		Columns: []meta.Column{
			meta.Column{Name: "col1", Type: 0, Nullable: false},
			meta.Column{Name: "col2", Type: 1, Nullable: false},
		},
	}

	// And saving the same plan again should result in an error
	res, err := e.executeSelect(selectPlan)
	assert.NoError(t, err)

	assert.ElementsMatch(t, expected.Rows, res.(*SelectResult).Rows)
	assert.ElementsMatch(t, expected.Columns, res.(*SelectResult).Columns)
}

func TestSelectNestedLoopJoinAsSelect(t *testing.T) {
	s := storage.NewMemStorage()
	sm := storage.NewSchemaManager(s)
	e := NewExecutor(s, sm)

	productTable := meta.Table{
		Name:       "products",
		PrimaryKey: "productid",
		Columns: []meta.Column{
			{Name: "productid", Type: meta.IntType},
			{Name: "productname", Type: meta.StringType},
			{Name: "categoryid", Type: meta.StringType},
		},
	}

	categoriesTable := meta.Table{
		Name:       "categories",
		PrimaryKey: "categoryid",
		Columns: []meta.Column{
			{Name: "categoryid", Type: meta.StringType},
			{Name: "categoryname", Type: meta.StringType},
		},
	}

	nodeInput := planner.SelectPlan{
		Node: &planner.NestedLoopJoin{
			Left: &planner.NodeScan{
				Table: productTable,
			},
			Right: &planner.NodeScan{
				Table: categoriesTable,
			},
			Predicate: &parser.Expression{
				Kind: parser.BinaryKind,
				BinaryExpression: &parser.BinaryExpression{
					Left: &parser.Expression{
						Kind: parser.ColumnKind,
						ColumnExpression: &parser.ColumnExpression{
							TableName:  "products",
							ColumnName: "categoryid",
						},
					},
					Right: &parser.Expression{
						Kind: parser.ColumnKind,
						ColumnExpression: &parser.ColumnExpression{
							TableName:  "categories",
							ColumnName: "categoryid",
						},
					},
					Operator: &lexer.Token{TokenType: lexer.EQUALS, Literal: "="},
				},
			},
		},
	}

	create1Plan := planner.CreateTablePlan{
		Table: productTable,
	}

	create2Plan := planner.CreateTablePlan{
		Table: categoriesTable,
	}

	insert1plan := planner.InsertPlan{
		Table: productTable,
		Row:   meta.Row{"productid": 1, "productname": "oliveoil", "categoryid": "1"},
	}

	insert2plan := planner.InsertPlan{
		Table: categoriesTable,
		Row:   meta.Row{"categoryid": "1", "categoryname": "oils"},
	}

	_, err := e.executeCreateTable(create1Plan)
	assert.NoError(t, err)

	_, err = e.executeCreateTable(create2Plan)
	assert.NoError(t, err)

	_, err = e.executeInsert(insert1plan)
	assert.NoError(t, err)

	_, err = e.executeInsert(insert2plan)
	assert.NoError(t, err)

	expected := &SelectResult{
		Rows: []meta.Row{
			{
				"categories.categoryid":   "1",
				"categories.categoryname": "oils",
				"products.categoryid":     "1",
				// TODO: Fix this float!
				"products.productid":   float64(1),
				"products.productname": "oliveoil",
			},
		},
		Columns: []meta.Column{},
	}

	res, err := e.executeSelect(nodeInput)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected.Rows, res.(*SelectResult).Rows)
	assert.ElementsMatch(t, expected.Columns, res.(*SelectResult).Columns)
}

func TestSelectNestedLoopJoinAsSelectWithProjection(t *testing.T) {
	s := storage.NewMemStorage()
	sm := storage.NewSchemaManager(s)
	e := NewExecutor(s, sm)

	productTable := meta.Table{
		Name:       "products",
		PrimaryKey: "productid",
		Columns: []meta.Column{
			{Name: "productid", Type: meta.IntType},
			{Name: "productname", Type: meta.StringType},
			{Name: "categoryid", Type: meta.StringType},
		},
	}

	categoriesTable := meta.Table{
		Name:       "categories",
		PrimaryKey: "categoryid",
		Columns: []meta.Column{
			{Name: "categoryid", Type: meta.StringType},
			{Name: "categoryname", Type: meta.StringType},
		},
	}

	nodeInput := planner.SelectPlan{
		Node: &planner.NodeProjection{
			Source: &planner.NestedLoopJoin{
				Left: &planner.NodeScan{
					Table: productTable,
				},
				Right: &planner.NodeScan{
					Table: categoriesTable,
				},
				Predicate: &parser.Expression{
					Kind: parser.BinaryKind,
					BinaryExpression: &parser.BinaryExpression{
						Left: &parser.Expression{
							Kind: parser.ColumnKind,
							ColumnExpression: &parser.ColumnExpression{
								TableName:  "products",
								ColumnName: "categoryid",
							},
						},
						Right: &parser.Expression{
							Kind: parser.ColumnKind,
							ColumnExpression: &parser.ColumnExpression{
								TableName:  "categories",
								ColumnName: "categoryid",
							},
						},
						Operator: &lexer.Token{TokenType: lexer.EQUALS, Literal: "="},
					},
				},
			},

			Expressions: &[]parser.Expression{
				parser.Expression{
					Kind: parser.ColumnKind,
					ColumnExpression: &parser.ColumnExpression{
						TableName:  "products",
						ColumnName: "productname",
					},
				},
				parser.Expression{
					Kind: parser.ColumnKind,
					ColumnExpression: &parser.ColumnExpression{
						TableName:  "categories",
						ColumnName: "categoryname",
					},
				},
			},
		},
	}

	create1Plan := planner.CreateTablePlan{
		Table: productTable,
	}

	create2Plan := planner.CreateTablePlan{
		Table: categoriesTable,
	}

	insert1plan := planner.InsertPlan{
		Table: productTable,
		Row:   meta.Row{"productid": 1, "productname": "oliveoil", "categoryid": "1"},
	}

	insert2plan := planner.InsertPlan{
		Table: categoriesTable,
		Row:   meta.Row{"categoryid": "1", "categoryname": "oils"},
	}

	insert3plan := planner.InsertPlan{
		Table: productTable,
		Row:   meta.Row{"productid": 2, "productname": "shampoo", "categoryid": "2"},
	}

	insert4plan := planner.InsertPlan{
		Table: categoriesTable,
		Row:   meta.Row{"categoryid": "2", "categoryname": "hair"},
	}

	_, err := e.executeCreateTable(create1Plan)
	assert.NoError(t, err)

	_, err = e.executeCreateTable(create2Plan)
	assert.NoError(t, err)

	_, err = e.executeInsert(insert1plan)
	assert.NoError(t, err)

	_, err = e.executeInsert(insert2plan)
	assert.NoError(t, err)

	_, err = e.executeInsert(insert3plan)
	assert.NoError(t, err)

	_, err = e.executeInsert(insert4plan)
	assert.NoError(t, err)

	expected := &SelectResult{
		Rows: []meta.Row{
			{
				"products.productname":    "oliveoil",
				"categories.categoryname": "oils",
			},
			{
				"products.productname":    "shampoo",
				"categories.categoryname": "hair",
			},
		},
		Columns: []meta.Column{},
	}

	res, err := e.executeSelect(nodeInput)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected.Rows, res.(*SelectResult).Rows)
	assert.ElementsMatch(t, expected.Columns, res.(*SelectResult).Columns)
}
