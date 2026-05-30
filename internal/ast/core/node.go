package core

type ASTNode interface {
	Node()
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

type ExprBase struct {
	Type Type
	Line int
}

type StmtBase struct {
	Line int
}
