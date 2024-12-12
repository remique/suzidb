package main

import (
	e "example.com/suzidb/executor"
	m "example.com/suzidb/meta"
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
				},
			},
		},
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
	// secondRes, err := exec.GetTable("tableName")
	// if err != nil {
	// 	fmt.Println("ER", err)
	// }
	// fmt.Println(secondRes)

	get, err := exec.GetTable("tableName")
	if err != nil {
		fmt.Println("ERR2: ", err)
	}
	fmt.Printf("%v\n", get)
}
