package fir

type LinkageType uint8

const (
	_ LinkageType = iota

	LinkInternal
	LinkExternal
	LinkWeak
)

type Global struct {
	valueBase

	Name       string
	Val        Value
	ActualType Type
	Linkage    LinkageType
}
