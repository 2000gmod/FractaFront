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
	errors   []diag.ErrorContainer

	prefixParsers  map[token.TokenType]prefixParser
	infixParsers   map[token.TokenType]infixParser
	postfixParsers map[token.TokenType]postfixParser
}

type prefixParser interface {
	Parse(*Parser, token.Token) ast.Expr
	Precedence() float32
}

type infixParser interface {
	Parse(*Parser, ast.Expr, token.Token) ast.Expr
	Lbp() float32
}

type postfixParser interface {
	Parse(*Parser, ast.Expr, token.Token) ast.Expr
	Precedence() float32
}
