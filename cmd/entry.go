package cmd

import (
	"bytes"
	"fmt"
	"fracta/internal/lexer"
)

func ProgramEntry() {
	testSrc := []byte(
		`hmm 3+3 3e-5`,
	)

	lex := lexer.NewLexerFromReader(bytes.NewReader(testSrc), "test")

	toks := lex.GetAllTokens()

	for _, t := range toks {
		fmt.Printf("%v\n", t)
	}
}
