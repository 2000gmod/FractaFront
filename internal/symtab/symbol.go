package symtab

import (
	"fracta/internal/ast/core"
	"strings"
)

type SymbolKind byte

const (
	KindVar SymbolKind = iota
	KindFunc
	KindType
)

// Symbol is a named entity in a scope.
type Symbol struct {
	Name   string       // Identifier name for this symbol.
	Kind   SymbolKind   // Symbol kind.
	Const  bool         // Is this symbol const?
	Public bool         // Is this symbol visible from outside the module? (file scope symbols only).
	Node   core.ASTNode // The AST node that defined this symbol.
	Type   core.Type    // The AST type that this symbol has.
	Module *Module      // Module where this symbol was defined (global symbols only).

	OwnerChain []*Symbol          // If this symbol is a member of some other entity, contains the chain from top-to-bottom.
	Members    map[string]*Symbol // Type-definition symbols only: the members of the current type.
}

func (s *Symbol) GetMangledName() string {
	b := strings.Builder{}
	b.WriteString("fr::")
	b.WriteString(s.Module.Name)
	b.WriteString("::")

	if s.OwnerChain != nil {
		for _, v := range s.OwnerChain {
			b.WriteString(v.Name)
			b.WriteString("::")
		}
	}
	b.WriteString(s.Name)
	return b.String()
}
