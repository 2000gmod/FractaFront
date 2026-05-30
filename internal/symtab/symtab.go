package symtab

import "fmt"

// Symbol Table represents all symbols inside a single scope
type SymbolTable struct {
	symbols map[string]*Symbol
	parent  *SymbolTable
}

func NewSymTab(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		symbols: make(map[string]*Symbol),
		parent:  parent,
	}
}

func (s *SymbolTable) NewChildTable() *SymbolTable {
	return NewSymTab(s)
}

func (s *SymbolTable) GetSymbol(name string) (sym *Symbol, ok bool) {
	val, ok := s.symbols[name]
	if !ok {
		if s.parent == nil {
			return nil, false
		}
		return s.parent.GetSymbol(name)
	}
	return val, true
}

func (s *SymbolTable) AddSymbol(name string, sym *Symbol) error {
	_, ok := s.GetSymbol(name)
	if ok {
		return fmt.Errorf("redefining symbol %q", name)
	}
	s.symbols[name] = sym
	return nil
}

func (s *SymbolTable) GetParent() *SymbolTable {
	return s.parent
}
