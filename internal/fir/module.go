package fir

type Module struct {
	Name string

	Types     []*Type
	Functions []*Function
	Globals   []*Global
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
		fn.Params = append(fn.Params, &Parameter{valueBase: valueBase{T: pt}})
	}

	m.Functions = append(m.Functions, fn)
	return fn
}

func (m *Module) NewExternalFunction(name string, typ *FuncType, conv CallingConv) *Function {
	f := m.NewFunction(name, typ, conv)
	f.External = true
	return f
}

func (m *Module) NewGlobal(ctx *Context, name string, typ Type, val Value, linkage LinkageType) *Global {
	gl := &Global{
		valueBase:  valueBase{T: ctx.GetPtrType()},
		Name:       name,
		ActualType: typ,
		Val:        val,
		Linkage:    linkage,
	}

	m.Globals = append(m.Globals, gl)
	return gl
}

func (m *Module) GetGlobalString(ctx *Context, name, str string) Value {
	s := ctx.ConstString(str)

	g := m.NewGlobal(ctx, name, s.Type(), s, LinkInternal)
	return g
}

func (m *Module) Dump() []byte {
	return nil
}
