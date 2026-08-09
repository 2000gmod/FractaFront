package ast

import (
	"fracta/internal/ast/core"
	"fracta/internal/symtab"
	"fracta/internal/token"
)

type FunctionDeclaration struct {
	core.StmtBase
	Name       token.Token
	Args       []ArgPair
	ReturnType core.Type
	Body       *BlockStatement
	Symbol     *symtab.Symbol
}

func (s *FunctionDeclaration) Node()                    {}
func (s *FunctionDeclaration) StmtNode() *core.StmtBase { return &s.StmtBase }

func (s *FunctionDeclaration) GetSigType() core.Type {
	argTypes := make([]core.Type, 0, len(s.Args))

	for _, v := range s.Args {
		argTypes = append(argTypes, v.Type)
	}

	return &FunctionType{
		ReturnType: s.ReturnType,
		ArgTypes:   argTypes,
	}
}

type ReturnStatement struct {
	core.StmtBase
	Value core.Expression
}

func (s *ReturnStatement) Node()                    {}
func (s *ReturnStatement) StmtNode() *core.StmtBase { return &s.StmtBase }

type ExpressionStatement struct {
	core.StmtBase
	Expression core.Expression
}

func (s *ExpressionStatement) Node()                    {}
func (s *ExpressionStatement) StmtNode() *core.StmtBase { return &s.StmtBase }

type BlockStatement struct {
	core.StmtBase
	Body []core.Statement
}

func (s *BlockStatement) Node()                    {}
func (s *BlockStatement) StmtNode() *core.StmtBase { return &s.StmtBase }
