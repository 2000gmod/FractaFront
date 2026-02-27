package codegen

import (
	"fracta/internal/ast"
	"io"
)

// Represents a generalized code generation backend. Should be the last step in the compilation pipeline.
type CodeGenerator interface {
	Generate(ast ast.AST, w io.Writer) error
}

var codegenFactoryMap = map[string]func(string) CodeGenerator{}

func registerBackend(backendId string, gen func(string) CodeGenerator) {
	codegenFactoryMap[backendId] = gen
}

func GetNewCodeGenerator(backendId string) CodeGenerator {
	gen, ok := codegenFactoryMap[backendId]
	if !ok {
		return nil
	}
	return gen(backendId)
}
