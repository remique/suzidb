package evaluator

import (
	"fmt"
)

type Value interface {
	IsValue()
	Type() string
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

// Should be an enum
func (lv *LiteralValue) Type() string { return "literal" }
func (bv *BooleanValue) Type() string { return "bool" }
func (iv *IntValue) Type() string     { return "int" }

// Converts primitive type to specific Value struct.
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

// Converts evaluated Value to native Golang primitive type.
func ValueToNative(i interface{}) (interface{}, error) {
	switch i.(type) {
	case *BooleanValue:
		return i.(*BooleanValue).Value, nil
	case *IntValue:
		return i.(*IntValue).Value, nil
	case *LiteralValue:
		return i.(*LiteralValue).Value, nil
	default:
		return nil, fmt.Errorf("Unrecognized Value: %s", i)
	}
}
