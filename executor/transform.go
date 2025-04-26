package executor

import (
	"fmt"

	"example.com/suzidb/evaluator"
	"example.com/suzidb/meta"
	"example.com/suzidb/parser"
)

type Transformer struct{}

func NewTransformer() *Transformer {
	return &Transformer{}
}

func (t *Transformer) Project(input []meta.Row, predicates *[]parser.Expression) ([]meta.Row, error) {
	var finalRows []meta.Row
	for _, row := range input {
		res, err := t.projectSingle(row, predicates)
		if err != nil {
			return nil, err
		}

		finalRows = append(finalRows, res)
	}

	return finalRows, nil
}

func (t *Transformer) projectSingle(input meta.Row, predicates *[]parser.Expression) (meta.Row, error) {
	finalRow := make(meta.Row)
	for _, predicate := range *predicates {
		evaluated, err := evaluator.NewEval(&predicate).
			Evaluate(evaluator.WithRow(input))
		if err != nil {
			return nil, err
		}

		asNative, err := evaluator.ValueToNative(evaluated)
		if err != nil {
			return nil, err
		}

		val, ok := input[asNative.(string)]
		if ok {
			finalRow[asNative.(string)] = val
		}
	}

	return finalRow, nil
}
