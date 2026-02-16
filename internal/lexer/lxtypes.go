package lexer

import (
	"bufio"
	"io"
	"os"
)

// Transforms valid Fracta source into a token stream
type Lexer struct {
	reader      *bufio.Reader
	closer      io.Closer
	currentLine int
	filename    string
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
	}, nil
}

func NewLexerFromReader(r io.Reader, name string) *Lexer {
	return &Lexer{
		reader:      bufio.NewReader(r),
		currentLine: 1,
		filename:    name,
	}
}

func (l *Lexer) IsOpen() bool {
	return l.reader != nil
}
