package symtab

// Module represents a single compilation unit.
type Module struct {
	Path     string
	Name     string
	Symbols  *SymbolTable
	Exported map[string]*Symbol
}

func NewModule(path, name string) *Module {
	return &Module{
		Path:     path,
		Name:     name,
		Symbols:  NewSymTab(nil),
		Exported: make(map[string]*Symbol),
	}
}
