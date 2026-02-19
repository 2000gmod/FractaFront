package parser

import (
	"fracta/internal/ast"
	"fracta/internal/diag"
	"fracta/internal/token"
)

type Parser struct {
	toks     []token.Token
	current  int
	filename string

	prefixParsers  map[token.TokenType]prefixParser
	infixParsers   map[token.TokenType]infixParser
	postfixParsers map[token.TokenType]postfixParser

	errors []*diag.ErrorContainer
	done   bool
}

type prefixParser interface {
	Parse(*Parser, token.Token) (ast.Expression, error)
	Precedence() int
}

type infixParser interface {
	Parse(*Parser, ast.Expression, token.Token) (ast.Expression, error)
	Lbp() int
}

type postfixParser interface {
	Parse(*Parser, ast.Expression, token.Token) (ast.Expression, error)
	Precedence() int
}
