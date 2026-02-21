package ast

import (
	"fmt"
	"fracta/internal/token"
	"strings"
)

type UnkownType struct{}

func (UnkownType) node()     {}
func (UnkownType) TypeNode() {}

func (UnkownType) String() string {
	return "<unknown>"
}

type BuiltinType struct {
	Name string
}

func (*BuiltinType) node()     {}
func (*BuiltinType) TypeNode() {}

func (b *BuiltinType) String() string {
	return b.Name
}

type NamedType struct {
	Name token.Token
}

func (*NamedType) node()     {}
func (*NamedType) TypeNode() {}

func (n *NamedType) String() string {
	return n.Name.String()
}

type FunctionType struct {
	ReturnType Type
	ArgTypes   []Type
}

func (*FunctionType) node()     {}
func (*FunctionType) TypeNode() {}

func (f *FunctionType) String() string {
	s := strings.Builder{}
	_, _ = s.WriteString("func(")

	if len(f.ArgTypes) != 0 {
		for i := range len(f.ArgTypes) - 1 {
			_, _ = fmt.Fprintf(&s, "%s, ", f.ArgTypes[i].String())
		}
		_, _ = s.WriteString(f.ArgTypes[len(f.ArgTypes)-1].String())
	}

	_, _ = fmt.Fprintf(&s, ") %s", f.ReturnType.String())

	return s.String()
}
