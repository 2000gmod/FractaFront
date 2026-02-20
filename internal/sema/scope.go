package sema

import "fmt"

type scope struct {
	sema    *SemanticAnalyzer
	symbols map[string]symbol
	parent  *scope
}

func newScope(parent *scope, sema *SemanticAnalyzer) *scope {
	return &scope{
		sema:    sema,
		symbols: map[string]symbol{},
		parent:  parent,
	}
}

func (s *scope) newChildScope() *scope {
	return newScope(s, s.sema)
}

func (s *scope) isSymbolPresent(name string) bool {
	_, ok := s.symbols[name]
	if !ok {
		if s.parent == nil {
			return false
		}
		return s.parent.isSymbolPresent(name)
	}
	return true
}

func (s *scope) addSymbol(name string, sym symbol) error {
	if s.isSymbolPresent(name) {
		return fmt.Errorf("redefining symbol")
	}
	s.symbols[name] = sym
	return nil
}
