package parser

import (
	"fracta/internal/ast"
	"fracta/internal/ast/core"
	"fracta/internal/token"
)

type LiteralParser struct{}

func (*LiteralParser) Parse(p *Parser, tok token.Token) (core.Expression, error) {
	return &ast.Literal{
		ExprBase: core.ExprBase{Line: tok.Line},
		Value:    tok,
	}, nil
}

func (*LiteralParser) Precedence() int {
	return 0
}

type IdentifierParser struct{}

func (*IdentifierParser) Parse(p *Parser, tok token.Token) (core.Expression, error) {
	ident := &ast.Identifier{
		ExprBase: core.ExprBase{Line: tok.Line},
		Ident:    tok,
	}

	if p.match(token.TokOpDoubleColon) {
		out := &ast.QualifiedName{
			ExprBase: core.ExprBase{Line: tok.Line},
			Parts:    []*ast.Identifier{ident},
		}

		for {
			next, err := p.consume(token.TokIdentifier, "expected identifier")
			if err != nil {
				return nil, err
			}
			out.Parts = append(
				out.Parts,
				&ast.Identifier{
					ExprBase: core.ExprBase{Line: next.Line},
					Ident:    *next,
				},
			)

			if !p.match(token.TokOpDoubleColon) {
				break
			}
		}
		return out, nil
	} else {
		return ident, nil
	}
}

func (*IdentifierParser) Precedence() int {
	return 0
}

type GroupingParser struct{}

func (*GroupingParser) Parse(p *Parser, tok token.Token) (core.Expression, error) {
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

func (o *PrefixOperatorParser) Parse(p *Parser, tok token.Token) (core.Expression, error) {
	right, err := p.parseExpression(o.rbp)

	if err != nil {
		return nil, err
	}

	return &ast.Unary{
		ExprBase: core.ExprBase{Line: tok.Line},
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

func (o *BinaryOperatorParser) Parse(p *Parser, left core.Expression, tok token.Token) (core.Expression, error) {
	rbp := o.precedence

	if o.assoc == AssocLeft {
		rbp++
	}

	right, err := p.parseExpression(rbp)

	if err != nil {
		return nil, err
	}

	return &ast.Binary{
		ExprBase: core.ExprBase{Line: tok.Line},
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

func (o *PostfixOperatorParser) Parse(p *Parser, left core.Expression, tok token.Token) (core.Expression, error) {
	return &ast.Unary{
		ExprBase: core.ExprBase{Line: tok.Line},
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

func (c *CallParser) Parse(p *Parser, left core.Expression, tok token.Token) (core.Expression, error) {
	args := make([]core.Expression, 0)

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
		ExprBase: core.ExprBase{Line: tok.Line},
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

func (c *IndexParser) Parse(p *Parser, left core.Expression, tok token.Token) (core.Expression, error) {
	args := make([]core.Expression, 0)

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
		ExprBase: core.ExprBase{Line: tok.Line},
		Indexee:  left,
		Indices:  args,
	}, nil
}

func (c *IndexParser) Precedence() int {
	return c.precedence
}
