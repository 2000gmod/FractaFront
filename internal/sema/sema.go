package sema

import (
	"fracta/internal/ast"
	"fracta/internal/diag"
)

type SemanticAnalyzer struct {
	moduleAsts []ast.ASTNode
	errors     []*diag.ErrorContainer
}

type symbolTable struct {
}
