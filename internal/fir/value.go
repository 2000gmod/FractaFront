package fir

type Value interface {
	Type() Type
}

type valueId uint64

type valueBase struct {
	typ Type
	id  valueId
}

func (v *valueBase) Type() Type {
	return v.typ
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
	Value uint64
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

type ConstantStruct struct {
	valueBase
	Fields []Constant
}

func (c *ConstantStruct) constant() {}

type ConstantArray struct {
	valueBase
	Elements []Constant
}

func (c *ConstantArray) constant() {}

type ConstantUnion struct {
	valueBase
	Member uint32
	Value  Constant
}

func (c *ConstantUnion) constant() {}

type Parameter struct {
	valueBase
	Name string
}
