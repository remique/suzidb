package tests

// import (
// 	"testing"

// 	"example.com/suzidb/executor"
// 	"example.com/suzidb/lexer"
// 	"example.com/suzidb/meta"
// 	"example.com/suzidb/parser"
// 	"example.com/suzidb/planner"
// 	"example.com/suzidb/storage"
// 	"github.com/stretchr/testify/assert"
// )

// func executeCommand(s storage.Storage, sm storage.SchemaManager, command string) (executor.ExecutionResult, error) {
// 	l := lexer.NewLexer(command)
// 	p := parser.NewParser(*l)
// 	planner := planner.NewPlanner(&sm)

// 	parsed, err := p.ParseStatement()
// 	if err != nil {
// 		return nil, err
// 	}

// 	plan, err := planner.Build(*parsed)
// 	if err != nil {
// 		return nil, err
// 	}

// 	e := executor.NewExecutor(s, &sm)
// 	res, err := e.ExecutePlan(plan)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res, nil
// }

// func TestSelectWithJoins(t *testing.T) {
// 	s := storage.NewMemStorage()
// 	sm := storage.NewSchemaManager(s)

// 	_, err := executeCommand(s, *sm, "create table categories(id int primary key, categoryname text);")
// 	assert.NoError(t, err)

// 	_, err = executeCommand(s, *sm, "create table products(id int primary key, productname text, categoryid int);")
// 	assert.NoError(t, err)

// 	_, err = executeCommand(s, *sm, "insert into products(id, productname, categoryid) values (1, 'oliveoil', 1);")
// 	assert.NoError(t, err)

// 	_, err = executeCommand(s, *sm, "insert into products(id, productname, categoryid) values (2, 'water', 2);")
// 	assert.NoError(t, err)

// 	_, err = executeCommand(s, *sm, "insert into categories(id, categoryname) values (1, 'oils');")
// 	assert.NoError(t, err)

// 	_, err = executeCommand(s, *sm, "insert into categories(id, categoryname) values (2, 'waters');")
// 	assert.NoError(t, err)

// 	res, err := executeCommand(s, *sm, "select products.id from products left join categories on products.id = categories.id;")

// 	expected := &executor.SelectResult{
// 		Rows: []meta.Row{
// 			meta.Row{"categories.categoryname": "oils", "categories.id": "1", "products.categoryid": "1", "products.id": "1", "products.productname": "oliveoil"},
// 			meta.Row{"categories.categoryname": "waters", "categories.id": "2", "products.categoryid": "2", "products.id": "2", "products.productname": "water"},
// 		},
// 		Columns: []meta.Column{},
// 	}

// 	assert.NoError(t, err)
// 	assert.ElementsMatch(t, expected.Rows, res.(*executor.SelectResult).Rows)
// }
