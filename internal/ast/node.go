package ast

type ASTNode interface {
	node()
}

type Expression interface {
	ASTNode
	exprNode()
}

type Statement interface {
	ASTNode
	stmtNode()
}

type Type interface {
	ASTNode
	typeNode()
}
