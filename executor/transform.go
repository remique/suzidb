package executor

import (
	"example.com/suzidb/evaluator"
	"example.com/suzidb/meta"
	"example.com/suzidb/parser"
	"strings"
)

type Transformer struct{}

func NewTransformer() *Transformer {
	return &Transformer{}
}

func (t *Transformer) Project(input []meta.Row, predicates *[]parser.Expression) ([]meta.Row, error) {
	var finalRows []meta.Row
	for _, row := range input {
		res, err := t.ProjectSingle(row, predicates)
		if err != nil {
			return nil, err
		}

		finalRows = append(finalRows, res)
	}

	return finalRows, nil
}

func (t *Transformer) ProjectSingle(input meta.Row, predicates *[]parser.Expression) (meta.Row, error) {
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

		key := predicate.ColumnExpression.TableName
		if predicate.ColumnExpression.ColumnName != "" {
			key = strings.Join(
				[]string{
					predicate.ColumnExpression.
						TableName,
					".",
					predicate.ColumnExpression.
						ColumnName,
				},
				"")
		}

		finalRow[key] = asNative
	}

	return finalRow, nil
}
