package evaluator

import (
	"fmt"
	"strings"

	"example.com/suzidb/meta"
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

// We should separate them based on left and right types.
func (ee *ExpressionEvaluator) evaluateBinaryExpr() (Value, error) {
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

func (ee *ExpressionEvaluator) evaluateQualifiedColumn(row meta.Row, prefix bool) (Value, error) {
	keyStr := ee.QualifiedColumnExpression.ColumnName.IdentifierExpression.Literal
	if prefix {
		keyStr = strings.Join(
			[]string{
				ee.QualifiedColumnExpression.
					TableName.IdentifierExpression.Literal,
					".",
				keyStr,
			},
			"")
	}

	get, ok := row[keyStr]
	if !ok {
		return nil, fmt.Errorf("Could not get key: %s", keyStr)
	}

	asValue, err := toValue(get)
	if err != nil {
		return nil, err
	}

	return asValue, nil
}
