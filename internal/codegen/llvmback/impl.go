package llvmback

import (
	"errors"
	"fracta/internal/ast"
	"io"

	"tinygo.org/x/go-llvm"
)

func NewLlvmGenerator(modname string) *llvmGenerator {
	ctx := llvm.NewContext()
	mod := ctx.NewModule(modname)
	bld := ctx.NewBuilder()

	return &llvmGenerator{
		ctx: ctx,
		mod: mod,
		bld: bld,
	}
}

func (g *llvmGenerator) Generate(ast ast.AST, w io.Writer) (e error) {
	defer func() {
		if r := recover(); r != nil {
			switch h := r.(type) {
			case generationPanic:
				e = errors.New(h.msg)
				return
			default:
				panic(r)
			}
		}
	}()
	panic("Not implemented")
}
