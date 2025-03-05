package evaluator

import (
	"fmt"
)

type Value interface {
	IsValue()
}

type LiteralValue struct {
	Value string
}

type BooleanValue struct {
	Value bool
}

type IntValue struct {
	Value int
}

func (lv *LiteralValue) IsValue() {}
func (bv *BooleanValue) IsValue() {}
func (iv *IntValue) IsValue()     {}

func toValue(i interface{}) (Value, error) {
	switch v := i.(type) {
	case int:
		return &IntValue{Value: v}, nil
	case string:
		return &LiteralValue{Value: v}, nil
	case bool:
		return &BooleanValue{Value: v}, nil
	default:
		return nil, fmt.Errorf("Unsupported value")

	}
}
