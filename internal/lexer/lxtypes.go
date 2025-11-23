package lexer

import (
	"bufio"
	"fracta/internal/diag"
	"io"
	"os"
)

// Transforms valid Fracta source into a token stream
type Lexer struct {
	reader      *bufio.Reader
	closer      io.Closer
	currentLine int
	filename    string
	errors      []*diag.ErrorContainer
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
		errors:      []*diag.ErrorContainer{},
	}, nil
}

func NewLexerFromReader(r io.Reader, name string) *Lexer {
	return &Lexer{
		reader:      bufio.NewReader(r),
		currentLine: 1,
		filename:    name,
		errors:      []*diag.ErrorContainer{},
	}
}

func (l *Lexer) IsOpen() bool {
	return l.reader != nil
}
