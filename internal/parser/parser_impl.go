package parser

import (
	"fracta/internal/ast"
	"fracta/internal/diag"
	"fracta/internal/token"
)

func (p *Parser) Parse() *ast.FileSourceNode {
	statements := make([]ast.Statement, 0)

	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			continue
		}
		statements = append(statements, stmt)
	}

	if len(diag.GlobalErrors) != 0 {
		return nil
	}

	return &ast.FileSourceNode{
		Filename:   p.filename,
		Statements: statements,
	}
}

func (p *Parser) typeExpr() (ast.Type, error) {
	switch {
	case p.match(token.TokIdentifier):
		return p.namedType()
	default:
		err := p.addError("invalid type expression")
		return nil, err
	}
}

func (p *Parser) namedType() (ast.Type, error) {
	name := p.previous()
	return &ast.NamedType{
		Name: *name,
	}, nil
}

func (p *Parser) statement() (ast.Statement, error) {
	var stmt ast.Statement = nil
	var err error = nil

	switch {
	case p.match(token.TokKwFunc):
		stmt, err = p.funcDeclStmt()
	case p.match(token.TokKwReturn):
		stmt, err = p.returnStmt()
	default:
		stmt, err = p.exprStmt()
	}

	return stmt, err
}

func (p *Parser) funcDeclStmt() (ast.Statement, error) {
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

	var rtp ast.Type

	if !p.check(token.TokOpenBracket) {
		rtp, err = p.typeExpr()

		if err != nil {
			return nil, err
		}
	}

	var body ast.Statement

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
		StmtBase:   ast.StmtBase{Line: line},
		Name:       *name,
		Args:       args,
		ReturnType: rtp,
		Body:       body,
	}, nil
}

func (p *Parser) returnStmt() (ast.Statement, error) {
	line := p.previous().Line
	var value ast.Expression
	var err error

	if p.match(token.TokSemicolon) {
		value = nil
	} else {
		value, err = p.parseExpression(0.0)

		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.TokSemicolon, "expected ';'")

		if err != nil {
			return nil, err
		}
	}

	return &ast.ReturnStatement{
		StmtBase: ast.StmtBase{Line: line},
		Value:    value,
	}, nil
}

func (p *Parser) blockStmt() (ast.Statement, error) {
	line := p.previous().Line
	body := make([]ast.Statement, 0)

	for !p.check(token.TokCloseBracket) {
		stmt, err := p.statement()

		if err != nil {
			return nil, err
		}

		body = append(body, stmt)
	}

	_, err := p.consume(token.TokCloseBracket, "expected '}'")

	if err != nil {
		return nil, err
	}

	return &ast.BlockStatement{
		StmtBase: ast.StmtBase{Line: line},
		Body:     body,
	}, nil
}

func (p *Parser) exprStmt() (ast.Statement, error) {
	expr, err := p.parseExpression(0.0)

	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokSemicolon, "expected ';'")

	if err != nil {
		return nil, err
	}

	return &ast.ExpressionStatement{
		Expression: expr,
	}, nil
}

func (p *Parser) parseExpression(minBp float32) (ast.Expression, error) {
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
	return left, nil
}
