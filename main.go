package main

import (
	"fracta/internal/diag"
	"fracta/internal/parser"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	testSrc := []byte(
		`func main() {
			a
			return 5;
		}`,
	)

	p := parser.FromString(string(testSrc), "test.fr")
	ast := p.Parse()

	if diag.HadErrors() {
		diag.ReportErrors()
		os.Exit(1)
	}

	spew.Config.Indent = "    "
	spew.Dump(ast)
}
