package fir

type LinkageType uint8

const (
	_ LinkageType = iota

	LinkInternal
	LinkExternal
	LinkWeak
)

type Global struct {
	Name    string
	Val     Value
	Ty      Type
	Linkage LinkageType
}

func (g *Global) Type() Type {
	return g.Ty
}
