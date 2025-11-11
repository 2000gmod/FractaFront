package lexer

import (
	"bufio"
	"fracta/internal/diag"
	tok "fracta/internal/token"
	"io"
	"os"
	"strings"
)

var punctuations = map[string]tok.TokenType{
	"+":  tok.TokOpPlus,
	"-":  tok.TokOpMinus,
	"*":  tok.TokOpStar,
	"/":  tok.TokOpSlash,
	"%":  tok.TokOpMod,
	"=":  tok.TokOpAssign,
	"==": tok.TokOpEq,
	"!=": tok.TokOpNotEq,
	"<":  tok.TokOpLessThan,
	">":  tok.TokOpGreaterThan,
	"<=": tok.TokOpLessEqual,
	">=": tok.TokOpGreaterEqual,
	"(":  tok.TokOpenParen,
	")":  tok.TokCloseParen,
	"[":  tok.TokOpenSquare,
	"]":  tok.TokCloseSquare,
	"{":  tok.TokOpenBracket,
	"}":  tok.TokCloseBracket,
	".":  tok.TokOpDot,
	":":  tok.TokOpColon,
	"::": tok.TokOpDoubleColon,
	",":  tok.TokOpComma,
	";":  tok.TokSemicolon,
}

var keywords = map[string]tok.TokenType{
	"func":   tok.TokKwFunc,
	"return": tok.TokKwReturn,
}

type matchResult int

const (
	mNone matchResult = iota
	mPartial
	mMatchButLongerPossible
	mFullMatch
)

func matchPunctuation(src string) matchResult {
	var exact, prefix bool

	for p := range punctuations {
		if strings.HasPrefix(p, src) {
			prefix = true
			if p == src {
				exact = true
			}
		}
	}
	switch {
	case exact && prefix:
		return mMatchButLongerPossible
	case exact:
		return mFullMatch
	case prefix:
		return mPartial
	default:
		return mNone
	}
}

// Transforms valid Fracta source into a token stream
type Lexer struct {
	reader      *bufio.Reader
	closer      io.Closer
	currentLine int
	filename    string
	errors      []diag.ErrorContainer
}

func NewLexerFromFile(path string) (*Lexer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &Lexer{
		reader:      bufio.NewReader(f),
		closer:      f,
		currentLine: 1,
		filename:    path,
		errors:      []diag.ErrorContainer{},
	}, nil
}

func NewLexerFromReader(r io.Reader, name string) *Lexer {
	return &Lexer{
		reader:      bufio.NewReader(r),
		currentLine: 1,
		filename:    name,
	}
}

// Closes the file handle, if any
func (l *Lexer) Close() error {
	if l.closer != nil {
		err := l.closer.Close()
		l.closer = nil
		return err
	}
	return nil
}

func (l *Lexer) advance() rune {
	if l.reader == nil {
		return 0
	}

	r, _, err := l.reader.ReadRune()
	if err == io.EOF {
		l.reader = nil
		return 0
	}
	return r
}

func (l *Lexer) peek() rune {
	if l.reader == nil {
		return 0
	}

	r, _, err := l.reader.ReadRune()
	if err == io.EOF {
		l.reader = nil
		return 0
	}
	_ = l.reader.UnreadRune()

	return r
}

// Gets all the tokens, until an EOF is reached
func (l *Lexer) GetAllTokens() []tok.Token {
	panic("TODO")
}

func (l *Lexer) GetToken() tok.Token {
	var out tok.Token

	if len(l.errors) != 0 || l.reader == nil {
		out.Kind = tok.TokError
		out.Line = l.currentLine
	} else if _, err := l.reader.Peek(1); err == io.EOF {
		l.reader = nil
		out.Kind = tok.TokEndOfFile
		out.Line = l.currentLine
	} else {
		l.ScanToken(&out)
	}

	return out
}

func (l *Lexer) ScanToken(t *tok.Token) {
	c := l.advance()

	if c == 0 {
		t.Kind = tok.TokEndOfFile
		t.Line = l.currentLine
		return
	}

	{
		skipFlag := false

		for ok := true; ok; ok = skipFlag {
		sw:
			switch c {
			case '\n':
				l.currentLine++
				fallthrough
			case ' ':
				fallthrough
			case '\t':
				fallthrough
			case '\r':
				skipFlag = true
				break sw
			default:
				skipFlag = false
			}
			if skipFlag {
				c = l.advance()
			}
		}

	}

	if c == 0 {
		t.Kind = tok.TokEndOfFile
		t.Line = l.currentLine
	}

	{
		var proc string = string(c)
		biggestMatch := ""

		res := matchPunctuation(proc)

		if res == mNone {
			goto L1
		}

		switch res {
		case mMatchButLongerPossible:
			biggestMatch = proc
		case mFullMatch:
			t.Kind = punctuations[proc]
			t.Line = l.currentLine
			return
		}

	loop:
		for {
			proc += string(l.advance())
			res = matchPunctuation(proc)

			switch res {
			case mNone:
				t.Line = l.currentLine
				if _, ok := punctuations[biggestMatch]; !ok {
					t.Kind = tok.TokError
					return
				}
				t.Kind = punctuations[biggestMatch]
				return
			case mPartial:
				continue loop
			case mMatchButLongerPossible:
				biggestMatch = proc
				continue loop
			case mFullMatch:
				t.Kind = punctuations[proc]
				t.Line = l.currentLine
				return
			}
		}
	}

L1:
	if isAlpha(c) {
		l.scanKeywordOrIdentifier(t, c)
		return
	}
	if isDigit(c) {
		l.scanNumberLiteral(t, c)
		return
	}
	if c == '"' {
		l.scanStringLiteral(t)
		return
	}
}

func (l *Lexer) scanKeywordOrIdentifier(t *tok.Token, c rune) {

}

func (l *Lexer) scanNumberLiteral(t *tok.Token, c rune) {

}

func (l *Lexer) scanStringLiteral(t *tok.Token) {

}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
func isAlphaNum(r rune) bool {
	return isAlpha(r) || isDigit(r) || r == '_' || r == '$'
}
