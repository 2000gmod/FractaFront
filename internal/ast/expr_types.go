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

func (e *Literal) node()               {}
func (e *Literal) ExprNode() *ExprBase { return &e.ExprBase }

type Identifier struct {
	ExprBase
	Ident token.Token
}

func (e *Identifier) node()               {}
func (e *Identifier) ExprNode() *ExprBase { return &e.ExprBase }

type Unary struct {
	ExprBase
	Op      token.Token
	SubExpr Expression
}

func (e *Unary) node()               {}
func (e *Unary) ExprNode() *ExprBase { return &e.ExprBase }

type Binary struct {
	ExprBase
	Op    token.Token
	Left  Expression
	Right Expression
}

func (e *Binary) node()               {}
func (e *Binary) ExprNode() *ExprBase { return &e.ExprBase }

type Call struct {
	ExprBase
	Callee Expression
	Args   []Expression
}

func (e *Call) node()               {}
func (e *Call) ExprNode() *ExprBase { return &e.ExprBase }

type Indexed struct {
	ExprBase
	Indexee Expression
	Indices []Expression
}

func (e *Indexed) node()               {}
func (e *Indexed) ExprNode() *ExprBase { return &e.ExprBase }
