package ast

import "fracta/internal/token"

type StmtBase struct {
	Line int
}

type ArgPair struct {
	Type Type
	Name token.Token
}

type FunctionDeclaration struct {
	StmtBase
	Name       token.Token
	Args       []ArgPair
	ReturnType Type
	Body       Statement
}

func (*FunctionDeclaration) node()     {}
func (*FunctionDeclaration) stmtNode() {}

type ReturnStatement struct {
	StmtBase
	Value Expression
}

func (*ReturnStatement) node()     {}
func (*ReturnStatement) stmtNode() {}

type ExpressionStatement struct {
	StmtBase
	Expression Expression
}

func (*ExpressionStatement) node()     {}
func (*ExpressionStatement) stmtNode() {}

type BlockStatement struct {
	StmtBase
	Body []Statement
}

func (*BlockStatement) node()     {}
func (*BlockStatement) stmtNode() {}
