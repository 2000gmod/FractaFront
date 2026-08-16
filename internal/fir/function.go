package fir

type CallingConv uint8

const (
	_ CallingConv = iota

	ConvFracta
	ConvC
)

type Function struct {
	valueBase

	Module *Module

	Name string
	Sig  *FuncType
	Conv CallingConv

	Params     []*Parameter
	Blocks     []*Block
	Stackslots []StackSlot

	Linkage  LinkageType
	External bool

	nextId valueId
}

func (f *Function) getNextId() valueId {
	id := f.nextId
	f.nextId++
	return id
}

func (f *Function) GetParameter(index int) Value {
	return f.Params[index]
}

type StackSlot struct {
	valueBase
	Name      string
	InnerType Type
}
