package cmd

import (
	"fracta/internal/parser"
)

func ProgramEntry() {
	testSrc := []byte(
		`func main() {
			return 5;
		}`,
	)

	p := parser.FromString(string(testSrc), "test")
	_, err := p.Parse()

	if err != nil {
		panic(err)
	}
}
