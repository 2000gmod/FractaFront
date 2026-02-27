package main

import (
	"fracta/internal/codegen"
	"fracta/internal/diag"
	"fracta/internal/pipeline"

	"github.com/alecthomas/kong"
	"github.com/davecgh/go-spew/spew"
)

var CLI struct {
	File string `arg:"" name:"file" default:"test.fr"`
}

func main() {
	codegen.RegisterAllBackends()
	spew.Config.Indent = "  "
	spew.Config.DisablePointerAddresses = true

	kong.Parse(&CLI)
	ast, err := pipeline.SingleFileReadingPipeline("test", CLI.File)

	if err != nil {
		switch e := err.(type) {
		case diag.ErrorList:
			diag.DiagnoseErrors(e)
			return
		default:
			panic(e)
		}
	}

	gen := codegen.GetNewCodeGenerator("llvm")
	gen.Generate(ast, nil)

	spew.Dump(ast)

}
