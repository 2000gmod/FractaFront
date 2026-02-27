package codegen

import "fracta/internal/codegen/llvmback"

func wrapFactory[T CodeGenerator](gen func(string) T) func(string) CodeGenerator {
	return func(s string) CodeGenerator {
		return gen(s)
	}
}

func RegisterAllBackends() {
	registerBackend("llvm", wrapFactory(llvmback.NewLlvmGenerator))
}
