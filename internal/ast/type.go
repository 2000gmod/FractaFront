package ast

import (
	"fmt"
	"fracta/internal/ast/core"
	"fracta/internal/token"
	"strings"
)

type UnkownType struct{}

func (UnkownType) Node()     {}
func (UnkownType) TypeNode() {}

func (UnkownType) String() string {
	return "<unknown>"
}

type BuiltinType struct {
	Name string
}

func (*BuiltinType) Node()     {}
func (*BuiltinType) TypeNode() {}

func (b *BuiltinType) String() string {
	return b.Name
}

type NamedType struct {
	Name token.Token
}

func (*NamedType) Node()     {}
func (*NamedType) TypeNode() {}

func (n *NamedType) String() string {
	return n.Name.Identifier
}

type FunctionType struct {
	ReturnType core.Type
	ArgTypes   []core.Type
}

func (*FunctionType) Node()     {}
func (*FunctionType) TypeNode() {}

func (f *FunctionType) String() string {
	s := strings.Builder{}
	s.WriteString("func(")

	if len(f.ArgTypes) != 0 {
		for i := range len(f.ArgTypes) - 1 {
			fmt.Fprintf(&s, "%s, ", f.ArgTypes[i].String())
		}
		s.WriteString(f.ArgTypes[len(f.ArgTypes)-1].String())
	}

	fmt.Fprintf(&s, ") %s", f.ReturnType.String())

	return s.String()
}
