package pipeline

import (
	"fracta/internal/ast"
	"fracta/internal/lexer"
	"fracta/internal/parser"
	"fracta/internal/sema"
)

// Does a single-source pass from file to AST
func SingleFileReadingPipeline(pkgName, fname string) (*ast.ModuleAST, error) {
	lex, err := lexer.NewLexerFromFile(fname)

	if err != nil {
		return nil, err
	}

	toks, err := lex.GetAllTokens()

	if err != nil {
		return nil, err
	}

	parser := parser.NewParser(toks, fname)
	fsn, err := parser.Parse()

	if err != nil {
		return nil, err
	}

	pkgAst := ast.ModuleAST{
		Files:       []*ast.FileSourceNode{fsn},
		PackageName: pkgName,
	}

	sm, err := sema.NewAnalyzer(pkgName, &pkgAst)

	if err != nil {
		return nil, err
	}

	pfsn, err := sm.Analyze()

	if err != nil {
		return nil, err
	}

	return pfsn, nil
}
