package ast

import (
	"fracta/internal/ast/core"
	"fracta/internal/symtab"
)

type FileSourceNode struct {
	Context    *symtab.FileContext
	Filename   string
	Statements []core.Statement
}

func (*FileSourceNode) node() {}

type ModuleAST struct {
	Files       []*FileSourceNode
	Module      *symtab.Module
	PackageName string
}

func (*ModuleAST) node() {}
