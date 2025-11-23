package ast

type FileSourceNode struct {
	Filename   string
	Statements []Statement
}

func (*FileSourceNode) node() {}
