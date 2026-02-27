package sema

import (
	"fracta/internal/ast"
	"fracta/internal/diag"
)

type SemanticAnalyzer struct {
	packageName string
	packageAsts ast.AST
	errors      []*diag.ErrorContainer
	pkgScope    *scope

	currentScope    *scope
	currentFile     string
	currentFunction *ast.FunctionDeclaration
}
