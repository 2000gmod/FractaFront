package main

import (
	"bytes"
	"fmt"
	"fracta/internal/codegen"
	"fracta/internal/diag"
	"fracta/internal/pipeline"
	"time"

	_ "fracta/internal/codegen/tags"

	"github.com/alecthomas/kong"
	"github.com/davecgh/go-spew/spew"
)

var CLI struct {
	File string `arg:"" name:"file" default:"test.fr"`
}

var start time.Time

func init() {
	start = time.Now()
}

func main() {
	defer func() {
		fmt.Printf("took %v\n", time.Since(start))
	}()

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

	gen, err := codegen.GetCodegen("llvm", &codegen.CodegenOptions{
		ModuleName: "test",
		OutputKind: codegen.OutputLLVM_IR,
	})

	if err != nil {
		panic(err)
	}

	buf := bytes.Buffer{}
	err = gen.Generate(ast, &buf)

	if err != nil {
		return
	}

	fmt.Println(buf.String())

	//spew.Dump(ast)

}
