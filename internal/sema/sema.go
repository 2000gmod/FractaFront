package sema

import (
	"fracta/internal/ast"
	"fracta/internal/diag"
	"fracta/internal/symtab"
)

type SemanticAnalyzer struct {
	moduleAst *ast.ModuleAST
	errors    []*diag.ErrorContainer
	module    *symtab.Module

	currentScope    *symtab.SymbolTable
	currentFile     *symtab.FileContext
	currentFunction *ast.FunctionDeclaration
}
