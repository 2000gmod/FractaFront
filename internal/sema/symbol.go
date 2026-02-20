package sema

import "fracta/internal/ast"

type symbolKind int

const (
	symbolFunction symbolKind = iota
	symbolType
)

type symbol interface {
	getSymbolKind() symbolKind
	getSymbolBase() *symbolBase
}

type symbolBase struct {
	pkg string
}

type functionSymbol struct {
	symbolBase
	decl *ast.FunctionDeclaration
}

func (functionSymbol) getSymbolKind() symbolKind {
	return symbolFunction
}

func (s *functionSymbol) getSymbolBase() *symbolBase {
	return &s.symbolBase
}
