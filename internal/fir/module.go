package fir

type Module struct {
	Name string

	Types     []Type
	Functions []Function
	Globals   []Global
}

func NewModule(name string) *Module {
	return &Module{
		Name:      name,
		Types:     nil,
		Functions: nil,
		Globals:   nil,
	}
}

func (m *Module) NewFunction(name string, typ *FuncType, conv CallingConv) *Function {
	fn := &Function{
		Name: name,
		Sig:  typ,
		Conv: conv,
	}

	for _, pt := range typ.Args {
		fn.Params = append(fn.Params, &Parameter{valueBase: valueBase{typ: pt}})
	}

	m.Functions = append(m.Functions, *fn)
	return fn
}
