package fir

type CallingConv string

const (
	ConvFracta CallingConv = "fracta"
	ConvC      CallingConv = "c"
)

type Function struct {
	valueBase

	Name string
	Sig  *FunctionType
	Conv CallingConv

	Params []*Parameter
	Blocks []*Block

	External bool
}
