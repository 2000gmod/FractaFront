package sema

import (
	"fmt"
	"fracta/internal/ast"
	"fracta/internal/diag"
	"fracta/internal/symtab"
)

func NewAnalyzer(modName string, moduleAst *ast.ModuleAST) (*SemanticAnalyzer, error) {
	if len(moduleAst.Files) == 0 {
		return nil, fmt.Errorf("no asts")
	}
	for _, v := range moduleAst.Files {
		if v == nil {
			return nil, fmt.Errorf("got a nil ast")
		}
	}

	a := &SemanticAnalyzer{
		moduleAst: moduleAst,
		errors:    make([]*diag.ErrorContainer, 0),
	}
	a.module = symtab.NewModule("", modName)
	a.currentScope = a.module.Symbols

	moduleAst.Module = a.module

	return a, nil
}
