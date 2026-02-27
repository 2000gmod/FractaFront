package llvmback

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

type llvmGenerator struct {
	ctx llvm.Context
	mod llvm.Module
	bld llvm.Builder
}

type generationPanic struct {
	msg string
}

func genPanic(f string, v ...any) generationPanic {
	return generationPanic{
		msg: fmt.Sprintf(f, v...),
	}
}
