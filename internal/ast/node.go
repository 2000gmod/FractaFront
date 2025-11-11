package ast

type ASTNode interface {
	Node() ASTNode
}

type Expr interface {
	ASTNode
	exprNode() Expr
}

type Statement interface {
	ASTNode
	stmtNode() Statement
}

type Type interface {
	ASTNode
	typeNode() Type
}
