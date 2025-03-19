package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"

	"example.com/suzidb/executor"
	"example.com/suzidb/lexer"
	"example.com/suzidb/parser"
	"example.com/suzidb/planner"
	"example.com/suzidb/storage"
)

const PROMPT = "> "

func Repl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	s := storage.NewMemStorage()
	sm := storage.NewSchemaManager(s)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.NewParser(*l)
		planner := planner.NewPlanner(sm)

		parsed, err := p.ParseStatement()
		if err != nil {
			fmt.Println("parse error: ", err)
		}

		plan, err := planner.Build(*parsed)
		if err != nil {
			fmt.Println("plan error: ", err)
		}

		e := executor.NewExecutor(s, sm)
		res, err := e.ExecutePlan(plan)
		if err != nil {
			fmt.Println("executor error: ", err)
		}

		switch x := res.(type) {
		case *executor.CreateTableResult:
			fmt.Println("OK")
		case *executor.InsertResult:
			fmt.Println("OK")
		case *executor.SelectResult:
			PrettyPrintSelect(x)
		}
	}
}

func PrettyPrintSelect(res *executor.SelectResult) {
	// Create table
	table := tablewriter.NewWriter(os.Stdout)

	// Extract column names for header
	headers := []string{}
	for _, col := range res.Columns {
		headers = append(headers, col.Name)
	}
	table.SetHeader(headers)

	// Populate rows
	for _, row := range res.Rows {
		rowData := []string{}
		for _, col := range res.Columns {
			val, exists := row[col.Name]
			if !exists {
				rowData = append(rowData, "NULL")
			} else {
				rowData = append(rowData, fmt.Sprintf("%v", val))
			}
		}
		table.Append(rowData)
	}

	// Render table
	table.Render()
}

func main() {
	Repl(os.Stdin, os.Stdout)
}
