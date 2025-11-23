package parser

import "fracta/internal/ast"

func (p *Parser) Parse() (*ast.FileSourceNode, error) {
	statements := make([]ast.Statement, 0)

	return &ast.FileSourceNode{
		Filename:   p.filename,
		Statements: statements,
	}, nil
}
