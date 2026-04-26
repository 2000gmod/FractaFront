package codegen

import (
	"fmt"
	"fracta/internal/ast"
	"io"
)

// Represents a generalized code generation backend. Should be the last step in the compilation pipeline.
type CodeGenerator interface {
	Generate(ast *ast.PackageAST, w io.Writer) error
	GetOutputKind() OutputType
}

type CodegenOptions struct {
	ModuleName string
	OutputKind OutputType
}

var generatorFactories = map[string]func(ops *CodegenOptions) CodeGenerator{}

// Should be called in each generator's init function.
func RegisterCodeGenerator(name string, factory func(ops *CodegenOptions) CodeGenerator) {
	generatorFactories[name] = factory
}

func ListAvailableGenerators() []string {
	out := make([]string, 0, len(generatorFactories))
	for v := range generatorFactories {
		out = append(out, v)
	}
	return out
}

func GetCodegen(name string, ops *CodegenOptions) (CodeGenerator, error) {
	f, ok := generatorFactories[name]
	if !ok {
		return nil, fmt.Errorf("code generator not available: %q", name)
	}

	return f(ops), nil
}
