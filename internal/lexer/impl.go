package lexer

import (
	"fmt"
	"fracta/internal/diag"
	tok "fracta/internal/token"
	"io"
	"strconv"
	"strings"
	"unicode"
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

type matchInfo struct {
	exact  bool
	longer bool
}

var punctInfo map[string]matchInfo

func init() {
	punctInfo = make(map[string]matchInfo)

	for p := range punctuations {
		for i := 1; i <= len(p); i++ {
			prefix := p[:i]
			info := punctInfo[prefix]
			if i == len(p) {
				info.exact = true
			} else {
				info.longer = true
			}
			punctInfo[prefix] = info
		}
	}
}

type matchResult int

const (
	mNone matchResult = iota
	mPartial
	mMatchButLongerPossible
	mFullMatch
)

func matchPunctuation(src string) matchResult {
	info, ok := punctInfo[src]

	if !ok {
		return mNone
	}

	switch {
	case info.exact && info.longer:
		return mMatchButLongerPossible
	case info.exact:
		return mFullMatch
	case info.longer:
		return mPartial
	default:
		return mNone
	}
}

func (l *Lexer) addError(f string, v ...any) {
	msg := fmt.Sprintf(f, v...)
	other := diag.CreateError(msg, l.filename, l.currentLine)

	l.Errors = append(l.Errors, other)
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

	if l.peekedValid {
		l.peekedValid = false
		return l.peekedRune
	}

	r, _, err := l.reader.ReadRune()
	if err == io.EOF {
		l.reader = nil
		return 0
	}

	if err != nil {
		l.reader = nil
		return 0
	}
	return r
}

func (l *Lexer) peek() rune {
	if l.reader == nil {
		return 0
	}

	if l.peekedValid {
		return l.peekedRune
	}

	r, _, err := l.reader.ReadRune()
	if err == io.EOF {
		l.reader = nil
		return 0
	}
	if err != nil {
		l.reader = nil
		return 0
	}
	l.peekedRune = r
	l.peekedValid = true
	return r
}

// Gets all the tokens, until an EOF is reached
func (l *Lexer) GetAllTokens() []tok.Token {
	out := make([]tok.Token, 0)
	for {
		nw := l.GetToken()
		out = append(out, nw)
		if nw.Kind == tok.TokEndOfFile || nw.Kind == tok.TokError {
			break
		}
	}
	return out
}

func (l *Lexer) GetToken() tok.Token {
	var out tok.Token

	if l.reader == nil {
		out.Kind = tok.TokEndOfFile
		out.Line = l.currentLine
		return out
	}

	l.ScanToken(&out)
	return out
}

func (l *Lexer) ScanToken(t *tok.Token) {
	c := l.advance()

	if c == 0 {
		t.Kind = tok.TokEndOfFile
		t.Line = l.currentLine
		return
	}

	// Skip whitespace and handle comments
	for {
		switch c {
		case ' ', '\t', '\r':
			c = l.advance()
		case '\n':
			l.currentLine++
			c = l.advance()
		case '/':
			r := l.peek()
			if r == '/' { // line comment
				_ = l.advance()
				for {
					r2 := l.advance()
					if r2 == 0 || r2 == '\n' {
						if r2 == '\n' {
							l.currentLine++
						}
						c = l.advance()
						break
					}
				}
			} else if r == '*' { // block comment
				_ = l.advance()
				for {
					r2 := l.advance()
					if r2 == 0 {
						t.Kind = tok.TokError
						t.Line = l.currentLine
						l.addError("unterminated block comment")
						return
					}
					if r2 == '\n' {
						l.currentLine++
					}
					if r2 == '*' && l.peek() == '/' {
						_ = l.advance()
						break
					}
				}
				c = l.advance()
			} else {
				goto doneSkipping
			}
		default:
			goto doneSkipping
		}
	}
doneSkipping:
	if c == 0 {
		t.Kind = tok.TokEndOfFile
		t.Line = l.currentLine
		return
	}

	if isDigit(c) || (c == '.' && isDigit(l.peek())) {
		l.scanNumberLiteral(t, c)
		return
	}

	if isAlpha(c) {
		l.scanKeywordOrIdentifier(t, c)
		return
	}

	if c == '"' {
		l.scanStringLiteral(t)
		return
	}

	if c == '\'' {
		l.scanCharLiteral(t)
		return
	}

	proc := string(c)
	biggestMatch := ""
	res := matchPunctuation(proc)
	if res == mNone {
		t.Kind = tok.TokError
		t.Line = l.currentLine
		return
	}
	switch res {
	case mMatchButLongerPossible:
		biggestMatch = proc
	case mFullMatch:
		t.Kind = punctuations[proc]
		t.Line = l.currentLine
		return
	}

	for {
		r := l.peek()
		if r == 0 {
			break
		}

		proc += string(r)
		res = matchPunctuation(proc)

		switch res {
		case mNone:
			if _, ok := punctuations[biggestMatch]; !ok {
				t.Kind = tok.TokError
				t.Line = l.currentLine
				return
			}
			t.Kind = punctuations[biggestMatch]
			t.Line = l.currentLine
			return
		case mPartial:
			_ = l.advance()
			continue
		case mMatchButLongerPossible:
			biggestMatch = proc
			_ = l.advance()
			continue
		case mFullMatch:
			_ = l.advance()
			t.Kind = punctuations[proc]
			t.Line = l.currentLine
			return
		}
	}
}

func (l *Lexer) scanKeywordOrIdentifier(t *tok.Token, c rune) {
	var sb strings.Builder
	sb.WriteRune(c)

	for {
		r := l.peek()
		if r == 0 || !isIdentifierPart(r) {
			break
		}
		_ = l.advance()
		sb.WriteRune(r)
	}

	lex := sb.String()

	if k, ok := keywords[lex]; ok {
		t.Kind = k
	} else {
		t.Kind = tok.TokIdentifier
		t.Lexeme = lex
		t.Identifier = lex
	}

	t.Line = l.currentLine
}

func (l *Lexer) scanNumberLiteral(t *tok.Token, first rune) {
	var sb strings.Builder
	sb.WriteRune(first)

	seenDot := first == '.'
	seenExp := false

	for {
		r := l.peek()
		if r == 0 {
			break
		}

		if r == '.' {
			if seenDot {
				break
			}
			_ = l.advance() // consume the dot
			next := l.peek()
			if !isDigit(next) {
				t.Kind = tok.TokError
				t.Line = l.currentLine
				l.addError("invalid number literal: %q", sb.String())
				return
			}
			sb.WriteRune('.')
			seenDot = true
			continue
		}

		if r == 'e' || r == 'E' {
			if seenExp {
				break
			}
			_ = l.advance()
			sb.WriteRune(r)
			seenExp = true
			r2 := l.peek()
			if r2 == '+' || r2 == '-' {
				_ = l.advance()
				sb.WriteRune(r2)
			}
			continue
		}

		if isDigit(r) || isAlpha(r) {
			_ = l.advance()
			sb.WriteRune(r)
			continue
		}

		break
	}

	lit := sb.String()

	if strings.HasSuffix(lit, ".") {
		t.Kind = tok.TokError
		t.Line = l.currentLine
		l.addError("invalid number literal: %q", sb.String())
		return
	}

	kind, val, err := ClassifyNumberLiteral(lit)
	if err != nil {
		t.Kind = tok.TokError
		t.Line = l.currentLine
		l.addError("invalid number literal: %q", sb.String())
		return
	}

	t.Kind = kind
	t.Value = val
	t.Lexeme = lit
	t.Line = l.currentLine
}

func (l *Lexer) scanStringLiteral(t *tok.Token) {
	var sb strings.Builder

	for {
		r := l.advance()
		if r == 0 {
			t.Kind = tok.TokError
			t.Line = l.currentLine
			l.addError("unterminated string literal")
			return
		}

		if r == '\n' {
			l.currentLine++
		}

		if r == '"' {
			break
		}

		sb.WriteRune(r)
	}

	raw := `"` + sb.String() + `"`

	val, err := strconv.Unquote(raw)
	if err != nil {
		t.Kind = tok.TokError
		t.Line = l.currentLine
		l.addError("invalid escape sequence")
		return
	}

	t.Kind = tok.TokString
	t.Value = val
	t.Lexeme = raw
	t.Line = l.currentLine
}

func (l *Lexer) scanCharLiteral(t *tok.Token) {
	var sb strings.Builder
	sb.WriteRune('\'')

	for {
		r := l.advance()
		if r == 0 {
			t.Kind = tok.TokError
			t.Line = l.currentLine
			l.addError("unterminated character literal")
			return
		}

		sb.WriteRune(r)

		if r == '\n' {
			l.currentLine++
		}

		if r == '\'' {
			break
		}

		if r == '\\' {
			r2 := l.advance()
			if r2 == 0 {
				t.Kind = tok.TokError
				t.Line = l.currentLine
				l.addError("unterminated escape sequence in character literal")
				return
			}
			sb.WriteRune(r2)
		}
	}

	raw := sb.String()

	if len(raw) < 2 {
		t.Kind = tok.TokError
		t.Line = l.currentLine
		l.addError("empty character literal")
		return
	}
	content := raw[1 : len(raw)-1]

	if content == "" {
		t.Kind = tok.TokError
		t.Line = l.currentLine
		l.addError("empty character literal")
		return
	}

	value, _, tail, err := strconv.UnquoteChar(content, '\'')
	if err != nil {
		t.Kind = tok.TokError
		t.Line = l.currentLine
		l.addError("invalid character literal: %v", err)
		return
	}
	if tail != "" {
		t.Kind = tok.TokError
		t.Line = l.currentLine
		l.addError("character literal must contain exactly one character")
		return
	}

	t.Kind = tok.TokChar
	t.Value = value
	t.Lexeme = raw
	t.Line = l.currentLine
}

func isAlpha(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
func isAlphaNum(r rune) bool {
	return isAlpha(r) || isDigit(r)
}

func isIdentifierPart(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
