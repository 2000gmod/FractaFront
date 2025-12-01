package parser

import (
	"fracta/internal/ast"
	"fracta/internal/diag"
	"fracta/internal/token"
)

func (p *Parser) Parse() (*ast.FileSourceNode, error) {
	statements := make([]ast.Statement, 0)

	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			continue
		}
		statements = append(statements, stmt)
	}

	if len(p.errors) != 0 {
		return nil, &diag.ErrorList{Errors: p.errors}
	}

	return &ast.FileSourceNode{
		Filename:   p.filename,
		Statements: statements,
	}, nil
}

func (p *Parser) typeExpr() (ast.Type, error) {
	switch {
	case p.match(token.TokIdentifier):
		return p.namedType()
	default:
		err := p.addError(diag.PInvalidTypeExpression)
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
	name, err := p.consume(token.TokIdentifier, diag.PExpectedIdentifier)

	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokOpenParen, diag.PExpectedFunctionParenthesis)

	if err != nil {
		return nil, err
	}

	args := make([]ast.ArgPair, 0)

	for !p.match(token.TokCloseParen) {
		pname, err := p.consume(token.TokIdentifier, diag.PExpectedParameterIdentifier)

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
			_, err = p.consume(token.TokOpComma, diag.PExpectedFunctionParamComma)
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
		err = p.addError(diag.PUnexpectedToken)
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

		_, err = p.consume(token.TokSemicolon, diag.PExpectedSemicolon)

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

	_, err := p.consume(token.TokCloseBracket, diag.PExpectedCloseBracketAfterBlock)

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

	_, err = p.consume(token.TokSemicolon, diag.PExpectedSemicolon)

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
		err := p.addError(diag.PInvalidExpressionToken)
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
