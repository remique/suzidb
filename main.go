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
			Name: "tableName",
			Columns: []m.Column{
				{
					Name:    "columnName",
					Type:    "columnType",
					Primary: false,
				},
			},
		},
	}

	stor := storage.NewMemStorage()
	exec := e.NewExecutor(stor)
	res := exec.ExecutePlan(&myPlan)

	fmt.Println(res)

	fmt.Println(stor.Get("tableName"))
}
