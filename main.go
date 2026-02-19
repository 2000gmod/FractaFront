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
	kong.Parse(&CLI)
	ast, err := pipeline.SingleFileReadingPipeline(CLI.File)

	if err != nil {
		diag.ReportErrors()
		return
	}

	spew.Config.Indent = "    "
	spew.Config.DisablePointerAddresses = true
	spew.Dump(ast)

}
