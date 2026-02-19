package parser

import (
	"fracta/internal/ast"
	"fracta/internal/token"
)

type LiteralParser struct{}

func (*LiteralParser) Parse(p *Parser, tok token.Token) (ast.Expression, error) {
	return &ast.Literal{
		ExprBase: ast.ExprBase{Line: tok.Line},
		Value:    tok,
	}, nil
}

func (*LiteralParser) Precedence() int {
	return 0
}

type IdentifierParser struct{}

func (*IdentifierParser) Parse(p *Parser, tok token.Token) (ast.Expression, error) {
	return &ast.Identifier{
		ExprBase: ast.ExprBase{Line: tok.Line},
		Ident:    tok,
	}, nil
}

func (*IdentifierParser) Precedence() int {
	return 0
}

type GroupingParser struct{}

func (*GroupingParser) Parse(p *Parser, tok token.Token) (ast.Expression, error) {
	expr, err := p.parseExpression(0)

	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokCloseParen, "expected ')'")

	return expr, err
}

func (*GroupingParser) Precedence() int {
	return 0
}

type PrefixOperatorParser struct {
	rbp int
}

func (o *PrefixOperatorParser) Parse(p *Parser, tok token.Token) (ast.Expression, error) {
	right, err := p.parseExpression(o.rbp)

	if err != nil {
		return nil, err
	}

	return &ast.Unary{
		ExprBase: ast.ExprBase{Line: tok.Line},
		Op:       tok,
		SubExpr:  right,
	}, nil
}

func (o *PrefixOperatorParser) Precedence() int {
	return o.rbp
}

type Assoc int

const (
	AssocLeft Assoc = iota
	AssocRight
)

type BinaryOperatorParser struct {
	precedence int
	assoc      Assoc
}

func (o *BinaryOperatorParser) Parse(p *Parser, left ast.Expression, tok token.Token) (ast.Expression, error) {
	rbp := o.precedence

	if o.assoc == AssocLeft {
		rbp++
	}

	right, err := p.parseExpression(rbp)

	if err != nil {
		return nil, err
	}

	return &ast.Binary{
		ExprBase: ast.ExprBase{Line: tok.Line},
		Op:       tok,
		Left:     left,
		Right:    right,
	}, nil
}

func (o *BinaryOperatorParser) Lbp() int {
	return o.precedence
}

type PostfixOperatorParser struct {
	precedence int
}

func (o *PostfixOperatorParser) Parse(p *Parser, left ast.Expression, tok token.Token) (ast.Expression, error) {
	return &ast.Unary{
		ExprBase: ast.ExprBase{Line: tok.Line},
		Op:       tok,
		SubExpr:  left,
	}, nil
}

func (o *PostfixOperatorParser) Precedence() int {
	return o.precedence
}

type CallParser struct {
	precedence int
}

func (c *CallParser) Parse(p *Parser, left ast.Expression, tok token.Token) (ast.Expression, error) {
	args := make([]ast.Expression, 0)

	if !p.check(token.TokCloseParen) {
		expr, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		args = append(args, expr)

		for p.check(token.TokOpComma) {
			_, err = p.consume(token.TokOpComma, "expected ','")
			if err != nil {
				return nil, err
			}
			expr, err = p.parseExpression(0)
			if err != nil {
				return nil, err
			}

			args = append(args, expr)
		}
	}
	_, err := p.consume(token.TokCloseParen, "expected ')'")
	if err != nil {
		return nil, err
	}
	return &ast.Call{
		ExprBase: ast.ExprBase{Line: tok.Line},
		Callee:   left,
		Args:     args,
	}, nil
}

func (c *CallParser) Precedence() int {
	return c.precedence
}

type IndexParser struct {
	precedence int
}

func (c *IndexParser) Parse(p *Parser, left ast.Expression, tok token.Token) (ast.Expression, error) {
	args := make([]ast.Expression, 0)

	if !p.check(token.TokCloseSquare) {
		expr, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		args = append(args, expr)

		for p.check(token.TokOpComma) {
			_, err = p.consume(token.TokOpComma, "expected ','")
			if err != nil {
				return nil, err
			}
			expr, err = p.parseExpression(0)
			if err != nil {
				return nil, err
			}

			args = append(args, expr)
		}
	}
	_, err := p.consume(token.TokCloseSquare, "expected ']'")
	if err != nil {
		return nil, err
	}
	return &ast.Indexed{
		ExprBase: ast.ExprBase{Line: tok.Line},
		Indexee:  left,
		Indices:  args,
	}, nil
}

func (c *IndexParser) Precedence() int {
	return c.precedence
}
