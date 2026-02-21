package sema

import "fmt"

type scope struct {
	symbols map[string]symbol
	parent  *scope
}

func newScope(parent *scope) *scope {
	return &scope{
		symbols: map[string]symbol{},
		parent:  parent,
	}
}

func (s *scope) newChildScope() *scope {
	return newScope(s)
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

func (s *scope) getSymbol(name string) (symbol, bool) {
	val, ok := s.symbols[name]
	if !ok {
		if s.parent == nil {
			return nil, false
		}
		return s.parent.getSymbol(name)
	}
	return val, true
}
