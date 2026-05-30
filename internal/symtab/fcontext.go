package symtab

import "fmt"

// Contains context for a single file
type FileContext struct {
	Name    string
	Path    string
	Imports map[string]*Module
}

func (f *FileContext) String() string {
	return fmt.Sprintf("%s/%s", f.Path, f.Name)
}
