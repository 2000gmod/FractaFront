package ast

type ASTNode interface {
	node()
}

type Expression interface {
	ASTNode
	ExprNode() *ExprBase
}

type Statement interface {
	ASTNode
	StmtNode() *StmtBase
}

type Type interface {
	ASTNode
	TypeNode()
	String() string
}
