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

func (s *FunctionDeclaration) node()               {}
func (s *FunctionDeclaration) StmtNode() *StmtBase { return &s.StmtBase }

type ReturnStatement struct {
	StmtBase
	Value Expression
}

func (s *ReturnStatement) node()               {}
func (s *ReturnStatement) StmtNode() *StmtBase { return &s.StmtBase }

type ExpressionStatement struct {
	StmtBase
	Expression Expression
}

func (s *ExpressionStatement) node()               {}
func (s *ExpressionStatement) StmtNode() *StmtBase { return &s.StmtBase }

type BlockStatement struct {
	StmtBase
	Body []Statement
}

func (s *BlockStatement) node()               {}
func (s *BlockStatement) StmtNode() *StmtBase { return &s.StmtBase }
