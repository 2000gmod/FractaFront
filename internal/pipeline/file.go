package pipeline

import (
	"fracta/internal/ast"
	"fracta/internal/lexer"
	"fracta/internal/parser"
)

// Does a single-source pass from file to AST (no AST analysis)
func SingleFileReadingPipeline(fname string) (*ast.FileSourceNode, error) {
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

	return fsn, nil
}
