package sema

import (
	"fmt"
	"fracta/internal/ast"
	"fracta/internal/diag"
)

func NewAnalyzer(pkgName string, packageAsts ...*ast.FileSourceNode) (*SemanticAnalyzer, error) {
	if len(packageAsts) == 0 {
		return nil, fmt.Errorf("no asts")
	}
	for _, v := range packageAsts {
		if v == nil {
			return nil, fmt.Errorf("got a nil ast")
		}
	}

	a := &SemanticAnalyzer{
		packageName: pkgName,
		packageAsts: packageAsts,
		errors:      make([]*diag.ErrorContainer, 0),
	}
	a.pkgScope = newScope(nil, a)
	a.currentScope = a.pkgScope

	return a, nil
}
