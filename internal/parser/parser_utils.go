package parser

import (
	"fmt"
	"fracta/internal/diag"
	"fracta/internal/token"
	"slices"
)

func (p *Parser) isAtEnd() bool {
	return p.peek().Kind == token.TokEndOfFile
}

func (p *Parser) peek() *token.Token {
	return &p.toks[p.current]
}

func (p *Parser) previous() *token.Token {
	if p.current <= 0 {
		return nil
	}
	return &p.toks[p.current-1]
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) consume(tt token.TokenType, f string, v ...any) (*token.Token, error) {
	if p.check(tt) {
		return p.advance(), nil
	}
	err := p.addError(f, v...)
	return nil, err
}

func (p *Parser) check(tts ...token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return slices.Contains(tts, p.peek().Kind)
}

func (p *Parser) match(tts ...token.TokenType) bool {
	if p.check(tts...) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) addError(f string, v ...any) *diag.ErrorContainer {
	msg := fmt.Sprintf(f, v...)
	o := diag.CreateError(msg, p.filename, p.previous().Line)

	diag.AppendError(o)
	return o
}
