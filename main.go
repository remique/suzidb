package main

import (
	m "example.com/suzidb/meta"
	// p "example.com/suzidb/parser"
	e "example.com/suzidb/executor"
	p "example.com/suzidb/planner"
	"example.com/suzidb/storage"
	"fmt"
)

func main() {

	myPlan := p.CreateTablePlan{
		Table: m.Table{
			Name:       "tableName",
			PrimaryKey: "columnName",
			Columns: []m.Column{
				{
					Name: "columnName",
					Type: m.IntType,
				},
				{
					Name: "columnName2",
					Type: m.StringType,
				},
			},
		},
	}

	myPlan2 := p.InsertPlan{
		Table: m.Table{
			Name:       "tableName",
			PrimaryKey: "columnName",
			Columns: []m.Column{
				{
					Name: "columnName",
					Type: m.IntType,
				},
				{
					Name: "columnName2",
					Type: m.StringType,
				},
			},
		},
		Rows: []m.Row{
			{"columnName": 1, "columnName2": "Alice"},
			{"columnName": 2, "columnName2": "Bob"},
		},
	}

	myPlan3 := p.QueryTablePlan{
		TableName: "tableName",
	}

	stor := storage.NewMemStorage()
	catalog := storage.NewSchemaManager(stor)
	exec := e.NewExecutor(stor, catalog)
	res, err := exec.ExecutePlan(&myPlan)
	if err != nil {
		fmt.Println("ERR", err)
	}

	fmt.Println(res)

	fmt.Println(stor.Get("meta:tableName"))

	get, err := exec.GetTable("tableName")
	if err != nil {
		fmt.Println("ERR2: ", err)
	}
	fmt.Printf("%v\n", get)

	// ----------------

	res, err = exec.ExecutePlan(&myPlan2)
	if err != nil {
		fmt.Println("ERR", err)
	}

	// Now get the table contents
	res, err = exec.ExecutePlan(&myPlan3)
	if err != nil {
		fmt.Println("ERR", err)
	}

	fmt.Println(res)
}
