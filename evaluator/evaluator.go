package evaluator

import (
	"example.com/suzidb/parser"
)

type ExpressionEvaluator struct {
	parser.Expression
}

func NewEval(expr parser.Expression) *ExpressionEvaluator {
	return &ExpressionEvaluator{
		expr,
	}
}

func (ee *ExpressionEvaluator) Evaluate(opts ...EvalOpts) {
	evalOpts := EvalOptions{}
	for _, opt := range opts {
		opt(&evalOpts)
	}

	// TODO
}
