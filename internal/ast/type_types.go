package ast

import "fracta/internal/token"

type NamedType struct {
	Name token.Token
}

func (*NamedType) node()     {}
func (*NamedType) typeNode() {}
