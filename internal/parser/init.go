package parser

import (
	"fracta/internal/diag"
	"fracta/internal/token"
)

func NewParser(toks []token.Token, filename string) *Parser {
	parser := Parser{}
	parser.filename = filename
	parser.toks = toks
	parser.Errors = make([]*diag.ErrorContainer, 0)

	parser.prefixParsers = map[token.TokenType]prefixParser{
		token.TokI8:     &LiteralParser{},
		token.TokI16:    &LiteralParser{},
		token.TokI32:    &LiteralParser{},
		token.TokI64:    &LiteralParser{},
		token.TokU8:     &LiteralParser{},
		token.TokU16:    &LiteralParser{},
		token.TokU32:    &LiteralParser{},
		token.TokU64:    &LiteralParser{},
		token.TokF32:    &LiteralParser{},
		token.TokF64:    &LiteralParser{},
		token.TokString: &LiteralParser{},
		token.TokChar:   &LiteralParser{},

		token.TokIdentifier: &IdentifierParser{},
		token.TokOpenParen:  &GroupingParser{},

		token.TokOpPlus:  &PrefixOperatorParser{rbp: 30},
		token.TokOpMinus: &PrefixOperatorParser{rbp: 30},
		token.TokOpStar:  &PrefixOperatorParser{rbp: 40},
	}

	parser.infixParsers = map[token.TokenType]infixParser{
		token.TokOpPlus:  &BinaryOperatorParser{lbp: 10, rbp: 11},
		token.TokOpMinus: &BinaryOperatorParser{lbp: 10, rbp: 11},
		token.TokOpStar:  &BinaryOperatorParser{lbp: 20, rbp: 21},
		token.TokOpSlash: &BinaryOperatorParser{lbp: 20, rbp: 21},
		token.TokOpMod:   &BinaryOperatorParser{lbp: 20, rbp: 21},
	}

	parser.postfixParsers = map[token.TokenType]postfixParser{
		token.TokOpenParen:  &CallParser{precedence: 50},
		token.TokOpenSquare: &IndexParser{precedence: 50},
	}

	return &parser
}
