package llvmgen

import (
	"fracta/internal/codegen"

	"tinygo.org/x/go-llvm"
)

type llvmGenerator struct {
	ctx     llvm.Context
	module  llvm.Module
	builder llvm.Builder
	options codegen.CodegenOptions

	symbolStack []map[string]llvm.Value
	globals     map[string]llvm.Value

	currentFunction llvm.Value
	currentBlock    llvm.BasicBlock

	loopStack []loopInfo

	typeCache map[string]llvm.Type
}

type loopInfo struct {
	breakBlock    llvm.BasicBlock
	continueBlock llvm.BasicBlock
}
