package main

import (
	"fracta/internal/diag"
	"fracta/internal/pipeline"

	"github.com/alecthomas/kong"
	"github.com/davecgh/go-spew/spew"
)

var CLI struct {
	File string `arg:"" name:"file" default:"test.fr"`
}

func main() {
	spew.Config.Indent = "    "
	spew.Config.DisablePointerAddresses = true

	kong.Parse(&CLI)
	ast, err := pipeline.SingleFileReadingPipeline("test", CLI.File)

	if err != nil {
		if diag.HadErrors() {
			diag.ReportErrors()
		} else {
			panic(err)
		}
		return
	}

	spew.Dump(ast)

}
