package ast

import "fracta/internal/token"

type ExprBase struct {
	Type Type
	Line int
}

type Literal struct {
	ExprBase
	Value token.Token
}

func (*Literal) node()     {}
func (*Literal) exprNode() {}

type Identifier struct {
	ExprBase
	Ident token.Token
}

func (*Identifier) node()     {}
func (*Identifier) exprNode() {}

type Unary struct {
	ExprBase
	Op      token.Token
	SubExpr Expression
}

func (*Unary) node()     {}
func (*Unary) exprNode() {}

type Binary struct {
	ExprBase
	Op    token.Token
	Left  Expression
	Right Expression
}

func (*Binary) node()     {}
func (*Binary) exprNode() {}

type Call struct {
	ExprBase
	Callee Expression
	Args   []Expression
}

func (*Call) node()     {}
func (*Call) exprNode() {}

type Indexed struct {
	ExprBase
	Indexee Expression
	Indices []Expression
}

func (*Indexed) node()     {}
func (*Indexed) exprNode() {}
