package main

import (
	_ "fracta/internal/codegen/tags"
	"fracta/internal/fir"
)

func main() {
	m := fir.NewModule("test")
	c := fir.NewContext()

	puts := m.NewExternalFunction(
		"puts",
		c.GetFuncType(
			false,
			c.GetVoidType(),
			c.GetPtrType(),
		),
		fir.ConvC,
	)

	fn := m.NewFunction(
		"main",
		c.GetFuncType(
			false,
			c.GetIntType(32),
		),
		fir.ConvFracta,
	)
	b := fir.NewBuilder(c, fn)
	b.NewBlock("entry")

	s := m.GetGlobalString(c, "str.1", "Hello, World\n\x00")
	b.InsertCall(puts, c.GetVoidType(), fir.ConvC, s)
	b.InsertRet(c.VoidValue())

}
