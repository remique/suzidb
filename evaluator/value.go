package evaluator

type Value interface {
	IsValue()
}

type LiteralValue struct {
	Value string
}

type BooleanValue struct {
	Value bool
}

func (lv *LiteralValue) IsValue() {}
func (bv *BooleanValue) IsValue() {}
