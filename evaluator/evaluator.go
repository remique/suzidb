package evaluator

import (
	"fmt"
	"strings"

	"example.com/suzidb/meta"
	"example.com/suzidb/parser"
)

type ExpressionEvaluator struct {
	expr parser.Expression
	opts EvalOptions
}

func NewEval(expr *parser.Expression) *ExpressionEvaluator {
	return &ExpressionEvaluator{
		expr: *expr,
		opts: EvalOptions{},
	}
}

func (ee *ExpressionEvaluator) Evaluate(opts ...EvalOpts) (Value, error) {
	for _, opt := range opts {
		opt(&ee.opts)
	}

	switch ee.expr.Kind {
	case parser.LiteralKind:
		return &LiteralValue{Value: ee.expr.LiteralExpression.Literal}, nil
	case parser.BinaryKind:
		return ee.evaluateBinaryExpr()
	case parser.ColumnKind:
		return ee.evaluateColumn(ee.opts.row)
	default:
		return nil, fmt.Errorf("Unsupported evaluation")
	}
}

// We should separate them based on left and right types.
func (ee *ExpressionEvaluator) evaluateBinaryExpr() (Value, error) {
	left, err := NewEval(ee.expr.BinaryExpression.Left).Evaluate(WithRow(ee.opts.row))
	if err != nil {
		return nil, err
	}

	right, err := NewEval(ee.expr.BinaryExpression.Right).Evaluate(WithRow(ee.opts.row))
	if err != nil {
		return nil, err
	}

	switch ee.expr.BinaryExpression.Operator.Literal {
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

func (ee *ExpressionEvaluator) evaluateColumn(row meta.Row) (Value, error) {
	keyStr := ee.expr.ColumnExpression.TableName
	if ee.expr.ColumnExpression.ColumnName != "" {
		fmt.Println("Column: ", ee.expr.ColumnExpression.ColumnName)
		keyStr = strings.Join(
			[]string{
				ee.expr.ColumnExpression.
					TableName,
				".",
				ee.expr.ColumnExpression.
					ColumnName,
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
