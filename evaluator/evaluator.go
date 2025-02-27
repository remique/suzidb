package evaluator

import (
	"fmt"

	"example.com/suzidb/parser"
)

type ExpressionEvaluator struct {
	parser.Expression
}

func NewEval(expr *parser.Expression) *ExpressionEvaluator {
	return &ExpressionEvaluator{
		*expr,
	}
}

func (ee *ExpressionEvaluator) Evaluate(opts ...EvalOpts) (Value, error) {
	evalOpts := EvalOptions{}
	for _, opt := range opts {
		opt(&evalOpts)
	}

	// TODO
	switch ee.Kind {
	case parser.LiteralKind:
		{
			return &LiteralValue{Value: ee.LiteralExpression.Literal}, nil
		}
	case parser.BinaryKind:
		{
			return ee.evaluateBinaryExpr()
		}
	default:
		{
			return nil, fmt.Errorf("Unsupported evaluation")
		}
	}
}

func (ee *ExpressionEvaluator) evaluateBinaryExpr(opts ...EvalOpts) (Value, error) {
	left, err := NewEval(ee.BinaryExpression.Left).Evaluate()
	if err != nil {
		return nil, err
	}

	right, err := NewEval(ee.BinaryExpression.Right).Evaluate()
	if err != nil {
		return nil, err
	}

	switch ee.BinaryExpression.Operator.Literal {
	case "=":
		{
			res := left.(*LiteralValue).Value == right.(*LiteralValue).Value
			return &BooleanValue{
				Value: res,
			}, nil
		}
	default:
		{
			return nil, fmt.Errorf("Unsupported evaluation")
		}
	}
}
