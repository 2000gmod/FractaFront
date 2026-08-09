package codegen

import (
	"fmt"
)

type OutputType int

const (
	OutputNone OutputType = iota
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

func CodegenPanicHandler(e *error) {
	if r := recover(); r != nil {
		switch h := r.(type) {
		case GenerationPanic:
			*e = fmt.Errorf("codegen panic: %s", h.Msg)
		default:
			panic(r)
		}
	}
}
