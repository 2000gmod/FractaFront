package codegen

import "fmt"

type OutputType int

const (
	_ OutputType = iota
	OutputNone
	OutputELFObject
	OutputLLVM_IR
	OutputAssembly
)

type GenerationPanic struct {
	Msg string
}

func DoPanic(f string, v ...any) {
	info := GenerationPanic{
		Msg: fmt.Sprintf(f, v...),
	}

	panic(info)
}
