package evaluator

import (
	"example.com/suzidb/meta"
)

type EvalOptions struct {
	row        meta.Row
	columnName string
}

type EvalOpts func(*EvalOptions) error

func WithRow(row meta.Row) EvalOpts {
	return func(o *EvalOptions) error {
		o.row = row

		return nil
	}
}
