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

func (*LiteralParser) Precedence() float32 {
	return 0.0
}

type IdentifierParser struct{}

func (*IdentifierParser) Parse(p *Parser, tok token.Token) (ast.Expression, error) {
	return &ast.Identifier{
		ExprBase: ast.ExprBase{Line: tok.Line},
		Ident:    tok,
	}, nil
}

func (*IdentifierParser) Precedence() float32 {
	return 0.0
}

type GroupingParser struct{}

func (*GroupingParser) Parse(p *Parser, tok token.Token) (ast.Expression, error) {
	expr, err := p.parseExpression(0.0)

	if err != nil {
		return nil, err
	}

	p.consume(token.TokCloseParen, "expected ')'")

	return expr, nil
}

func (*GroupingParser) Precedence() float32 {
	return 0.0
}

type PrefixOperatorParser struct {
	rbp float32
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

func (o *PrefixOperatorParser) Precedence() float32 {
	return o.rbp
}

type BinaryOperatorParser struct {
	lbp float32
	rbp float32
}

func (o *BinaryOperatorParser) Parse(p *Parser, left ast.Expression, tok token.Token) (ast.Expression, error) {
	right, err := p.parseExpression(o.rbp)

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

func (o *BinaryOperatorParser) Lbp() float32 {
	return o.lbp
}

type PostfixOperatorParser struct {
	precedence float32
}

func (o *PostfixOperatorParser) Parse(p *Parser, left ast.Expression, tok token.Token) (ast.Expression, error) {
	return &ast.Unary{
		ExprBase: ast.ExprBase{Line: tok.Line},
		Op:       tok,
		SubExpr:  left,
	}, nil
}

func (o *PostfixOperatorParser) Precedence() float32 {
	return o.precedence
}

type CallParser struct {
	precedence float32
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

func (c *CallParser) Precedence() float32 {
	return c.precedence
}

type IndexParser struct {
	precedence float32
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

func (c *IndexParser) Precedence() float32 {
	return c.precedence
}
