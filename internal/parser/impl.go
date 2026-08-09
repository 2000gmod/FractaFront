package parser

import (
	"fmt"
	"fracta/internal/ast"
	"fracta/internal/ast/core"
	"fracta/internal/diag"
	"fracta/internal/token"
)

func (p *Parser) Parse() (*ast.FileSourceNode, error) {
	if p.done {
		return nil, fmt.Errorf("already parsed")
	}

	defer func() { p.done = true }()

	statements := make([]core.Statement, 0)

	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			continue
		}
		statements = append(statements, stmt)
	}

	if len(p.errors) > 0 {
		return nil, diag.ErrorList(p.errors)
	}

	return &ast.FileSourceNode{
		Filename:   p.filename,
		Statements: statements,
	}, nil
}

func (p *Parser) typeExpr() (core.Type, error) {
	switch {
	case p.match(token.TokIdentifier):
		id := p.previous()
		btype, ok := ast.BuiltinTypeNameMap[id.Identifier]

		if !ok {
			return p.namedType()
		} else {
			return btype, nil
		}

	default:
		err := p.addError("invalid type expression")
		return nil, err
	}
}

func (p *Parser) namedType() (core.Type, error) {
	name := p.previous()
	return &ast.NamedType{
		Name: *name,
	}, nil
}

func (p *Parser) statement() (core.Statement, error) {
	var stmt core.Statement
	var err error

	switch {
	case p.match(token.TokKwFunc):
		stmt, err = p.funcDeclStmt()
	case p.match(token.TokKwReturn):
		stmt, err = p.returnStmt()
	default:
		stmt, err = p.exprStmt()
	}

	if err != nil {
		p.synchronize()
	}

	return stmt, err
}

func (p *Parser) funcDeclStmt() (core.Statement, error) {
	line := p.previous().Line
	name, err := p.consume(token.TokIdentifier, "expected identifier")

	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokOpenParen, "expected '('")

	if err != nil {
		return nil, err
	}

	args := make([]ast.ArgPair, 0)

	for !p.match(token.TokCloseParen) {
		pname, err := p.consume(token.TokIdentifier, "expected parameter identifier")

		if err != nil {
			return nil, err
		}

		ptype, err := p.typeExpr()

		if err != nil {
			return nil, err
		}

		args = append(args, ast.ArgPair{
			Type: ptype,
			Name: *pname,
		})

		if p.match(token.TokCloseParen) {
			break
		} else {
			_, err = p.consume(token.TokOpComma, "expected ','")
			if err != nil {
				return nil, err
			}
		}
	}

	var rtp core.Type

	if !p.check(token.TokOpenBracket) {
		rtp, err = p.typeExpr()

		if err != nil {
			return nil, err
		}
	}

	var body core.Statement

	if p.match(token.TokOpenBracket) {
		body, err = p.blockStmt()
		if err != nil {
			return nil, err
		}
	} else if p.match(token.TokSemicolon) {
		body = nil
	} else {
		err = p.addError("unexpected token")
		return nil, err
	}

	return &ast.FunctionDeclaration{
		StmtBase:   core.StmtBase{Line: line},
		Name:       *name,
		Args:       args,
		ReturnType: rtp,
		Body:       body.(*ast.BlockStatement),
	}, nil

}

func (p *Parser) returnStmt() (core.Statement, error) {
	line := p.previous().Line
	var value core.Expression
	var err error

	if p.match(token.TokSemicolon) {
		value = nil
	} else {
		value, err = p.parseExpression(0)

		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.TokSemicolon, "expected ';'")

		if err != nil {
			return nil, err
		}
	}

	return &ast.ReturnStatement{
		StmtBase: core.StmtBase{Line: line},
		Value:    value,
	}, nil
}

func (p *Parser) blockStmt() (core.Statement, error) {
	line := p.previous().Line
	body := make([]core.Statement, 0)

	for !p.check(token.TokCloseBracket) && !p.isAtEnd() {
		stmt, err := p.statement()

		if err != nil {
			if !p.synchronize(token.TokCloseBracket) {
				return nil, err
			}
			continue
		}

		body = append(body, stmt)
	}

	_, err := p.consume(token.TokCloseBracket, "expected '}'")

	if err != nil {
		return nil, err
	}

	return &ast.BlockStatement{
		StmtBase: core.StmtBase{Line: line},
		Body:     body,
	}, nil
}

func (p *Parser) exprStmt() (core.Statement, error) {
	expr, err := p.parseExpression(0)

	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokSemicolon, "expected ';'")

	if err != nil {
		return nil, err
	}

	return &ast.ExpressionStatement{
		StmtBase:   core.StmtBase{Line: expr.ExprNode().Line},
		Expression: expr,
	}, nil
}

func (p *Parser) parseExpression(minBp int) (core.Expression, error) {
	tok := p.advance()

	prefix, ok := p.prefixParsers[tok.Kind]

	if !ok {
		err := p.addError("invalid token in expression: %v", *tok)
		return nil, err
	}

	left, err := prefix.Parse(p, *tok)

	if err != nil {
		return nil, err
	}

	for {
		nextTok := p.peek()

		postfix, ok := p.postfixParsers[nextTok.Kind]

		if ok && postfix.Precedence() >= minBp {
			tok2 := p.advance()
			left, err = postfix.Parse(p, left, *tok2)

			if err != nil {
				return nil, err
			}

			continue
		}

		infix, ok := p.infixParsers[nextTok.Kind]

		if ok && infix.Lbp() >= minBp {
			tok2 := p.advance()
			left, err = infix.Parse(p, left, *tok2)

			if err != nil {
				return nil, err
			}

			continue
		}
		break
	}

	left.ExprNode().Type = &ast.UnkownType{}

	return left, nil
}
