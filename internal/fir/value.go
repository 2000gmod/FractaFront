package fir

type Value interface {
	Type() Type
}

type valueId uint64

type valueBase struct {
	T  Type
	id valueId
}

func (v *valueBase) Type() Type {
	return v.T
}

type Parameter struct {
	valueBase
	Name string
}

type Constant interface {
	Value
	constant()
}

type ConstantInt struct {
	valueBase
	Value uint64
}

func (c *ConstantInt) constant() {}

type ConstantFloat struct {
	valueBase
	Value float64
}

func (c *ConstantFloat) constant() {}

type ConstantBool struct {
	valueBase
	Value bool
}

func (c *ConstantBool) constant() {}

type ConstantNull struct {
	valueBase
}

func (c *ConstantNull) constant() {}

type ConstantArray struct {
	valueBase
	Elements []Constant
}

func (c *ConstantArray) constant() {}

type ConstantVoid struct {
	valueBase
}

func (c *ConstantVoid) constant() {}
