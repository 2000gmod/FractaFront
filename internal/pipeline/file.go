package pipeline

import (
	"fmt"
	"fracta/internal/ast"
	"fracta/internal/diag"
	"fracta/internal/lexer"
	"fracta/internal/parser"
)

// Does a single-source pass from file to AST (no AST analysis)
func SingleFileReadingPipeline(fname string) (*ast.FileSourceNode, error) {
	lex, err := lexer.NewLexerFromFile(fname)

	if err != nil {
		return nil, err
	}

	toks := lex.GetAllTokens()

	if len(lex.Errors) > 0 {
		diag.AppendError(lex.Errors...)
		return nil, fmt.Errorf("had lexing errors")
	}

	parser := parser.NewParser(toks, fname)
	fsn := parser.Parse()

	if len(parser.Errors) > 0 {
		diag.AppendError(parser.Errors...)
		return nil, fmt.Errorf("had parsing errors")
	}

	return fsn, nil
}
