package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"example.com/suzidb/lexer"
	"example.com/suzidb/parser"
)

const PROMPT = "> "

func Repl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.NewParser(*l)

		parsed, err := p.ParseStatement()
		if err != nil {
			fmt.Println("err", err)
		}

		// As a quick workaround we print parsed info as JSON.
		asJson, _ := json.Marshal(parsed)
		fmt.Println("parsed: ", string(asJson))
	}
}

func main() {
	Repl(os.Stdin, os.Stdout)
}
