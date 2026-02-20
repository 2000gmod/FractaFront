package pipeline

import (
	"fracta/internal/ast"
	"fracta/internal/lexer"
	"fracta/internal/parser"
	"fracta/internal/sema"
)

// Does a single-source pass from file to AST (no AST analysis)
func SingleFileReadingPipeline(pkgName, fname string) ([]*ast.FileSourceNode, error) {
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

	sm, err := sema.NewAnalyzer(pkgName, fsn)

	if err != nil {
		return nil, err
	}

	pfsn, err := sm.Analyze()

	if err != nil {
		return nil, err
	}

	return pfsn, nil
}
