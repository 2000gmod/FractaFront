package parser

import (
	"fracta/internal/diag"
	"fracta/internal/lexer"
	"fracta/internal/token"
	"strings"
)

func FromString(src, filename string) *Parser {
	lex := lexer.NewLexerFromReader(strings.NewReader(src), filename)
	return FromScanner(lex, filename)
}

func FromFile(fpath string) *Parser {
	lex, err := lexer.NewLexerFromFile(fpath)
	if err != nil {
		panic(err)
	}
	return FromScanner(lex, fpath)
}

func FromScanner(lex *lexer.Lexer, filename string) *Parser {
	parser := Parser{}
	parser.toks = make([]token.Token, 0)
	parser.errors = make([]diag.ErrorContainer, 0)

	if !lex.IsOpen() {
		var tok token.Token
		tok.Kind = token.TokEndOfFile
		parser.toks = append(parser.toks, tok)
		return &parser
	}

	toks := lex.GetAllTokens()
	parser.toks = append(parser.toks, toks...)
	return &parser
}
