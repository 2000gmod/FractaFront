package parser

import (
	"fracta/internal/ast/core"
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
	Parse(*Parser, token.Token) (core.Expression, error)
	Precedence() int
}

type infixParser interface {
	Parse(*Parser, core.Expression, token.Token) (core.Expression, error)
	Lbp() int
}

type postfixParser interface {
	Parse(*Parser, core.Expression, token.Token) (core.Expression, error)
	Precedence() int
}
