package ast

import (
	"fracta/internal/ast/core"
	"fracta/internal/token"
)

type Literal struct {
	core.ExprBase
	Value token.Token
}

func (e *Literal) Node()                    {}
func (e *Literal) ExprNode() *core.ExprBase { return &e.ExprBase }

type Identifier struct {
	core.ExprBase
	Ident token.Token
}

func (e *Identifier) Node()                    {}
func (e *Identifier) ExprNode() *core.ExprBase { return &e.ExprBase }

type QualifiedName struct {
	core.ExprBase
	Parts []*Identifier
}

func (e *QualifiedName) Node()                    {}
func (e *QualifiedName) ExprNode() *core.ExprBase { return &e.ExprBase }

type Unary struct {
	core.ExprBase
	Op      token.Token
	SubExpr core.Expression
}

func (e *Unary) Node()                    {}
func (e *Unary) ExprNode() *core.ExprBase { return &e.ExprBase }

type Binary struct {
	core.ExprBase
	Op    token.Token
	Left  core.Expression
	Right core.Expression
}

func (e *Binary) Node()                    {}
func (e *Binary) ExprNode() *core.ExprBase { return &e.ExprBase }

type Call struct {
	core.ExprBase
	Callee core.Expression
	Args   []core.Expression
}

func (e *Call) Node()                    {}
func (e *Call) ExprNode() *core.ExprBase { return &e.ExprBase }

type Indexed struct {
	core.ExprBase
	Indexee core.Expression
	Indices []core.Expression
}

func (e *Indexed) Node()                    {}
func (e *Indexed) ExprNode() *core.ExprBase { return &e.ExprBase }
