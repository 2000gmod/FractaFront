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
	getExprType() ast.Type
}

type symbolBase struct {
	pkg string
}

type functionSymbol struct {
	symbolBase
	fType *ast.FunctionType
}

func (functionSymbol) getSymbolKind() symbolKind {
	return symbolFunction
}

func (s *functionSymbol) getSymbolBase() *symbolBase {
	return &s.symbolBase
}

func (s *functionSymbol) getExprType() ast.Type {
	return s.fType
}
