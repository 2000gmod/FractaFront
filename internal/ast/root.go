package ast

type FileSourceNode struct {
	Filename   string
	Statements []Statement
}

func (*FileSourceNode) node() {}

type PackageAST struct {
	Files       []*FileSourceNode
	PackageName string
}

func (*PackageAST) node() {}
