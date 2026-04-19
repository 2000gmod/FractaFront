package llvmgen

import (
	"errors"
	"fracta/internal/ast"
	"fracta/internal/codegen"
	"io"

	"tinygo.org/x/go-llvm"
)

func init() {
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()

	err := llvm.InitializeNativeTarget()

	if err != nil {
		panic(err)
	}

	codegen.RegisterCodeGenerator("llvm", NewLlvmGenerator)
}

func NewLlvmGenerator(ops *codegen.CodegenOptions) codegen.CodeGenerator {
	ctx := llvm.NewContext()
	mod := ctx.NewModule(ops.ModuleName)
	bld := ctx.NewBuilder()

	cache := map[string]llvm.Type{
		"i8":  ctx.Int8Type(),
		"i16": ctx.Int16Type(),
		"i32": ctx.Int32Type(),
		"i64": ctx.Int64Type(),

		"u8":  ctx.Int8Type(),
		"u16": ctx.Int16Type(),
		"u32": ctx.Int32Type(),
		"u64": ctx.Int64Type(),

		"f32": ctx.FloatType(),
		"f64": ctx.DoubleType(),

		"bool": ctx.Int8Type(),

		"ptr": llvm.PointerType(ctx.Int8Type(), 0),
	}

	return &llvmGenerator{
		ctx:     ctx,
		module:  mod,
		builder: bld,
		options: *ops,

		symbolStack: make([]map[string]llvm.Value, 0),
		globals:     make(map[string]llvm.Value),

		currentFunction: llvm.Value{},
		currentBlock:    llvm.BasicBlock{},

		loopStack: make([]loopInfo, 0),

		typeCache: cache,
	}
}

func (g *llvmGenerator) GetOutputKind() codegen.OutputType {
	return g.options.OutputKind
}

func (g *llvmGenerator) Generate(ast *ast.PackageAST, w io.Writer) (e error) {
	defer func() {
		if r := recover(); r != nil {
			switch h := r.(type) {
			case codegen.GenerationPanic:
				e = errors.New(h.Msg)
			default:
				panic(r)
			}
		}
	}()

	for _, f := range ast.Files {
		g.genFile(f)
	}

	if g.options.OutputKind == codegen.OutputLLVM_IR {
		w.Write([]byte(g.module.String()))
	} else {
		tgtTriple := llvm.DefaultTargetTriple()
		tgt, err := llvm.GetTargetFromTriple(tgtTriple)

		if err != nil {
			return err
		}

		machine := tgt.CreateTargetMachine(tgtTriple, "generic", "", llvm.CodeGenLevelNone, llvm.RelocPIC, llvm.CodeModelDefault)
		defer machine.Dispose()

		var tp llvm.CodeGenFileType

		switch g.options.OutputKind {
		case codegen.OutputAssembly:
			tp = llvm.AssemblyFile
		case codegen.OutputELFObject:
			tp = llvm.ObjectFile
		default:
			codegen.DoPanic("Unsupported output kind.")
		}

		buf, err := machine.EmitToMemoryBuffer(g.module, tp)

		if err != nil {
			return err
		}

		w.Write(buf.Bytes())
	}

	return nil
}

func (g *llvmGenerator) llvmTypeFromType(t ast.Type) (llvm.Type, bool) {
	var key string
	switch rt := t.(type) {
	case *ast.BuiltinType:
		key = rt.Name
	case *ast.NamedType:
		key = rt.Name.Identifier
	default:
	}

	l, ok := g.typeCache[key]
	return l, ok
}

func (g *llvmGenerator) genFile(f *ast.FileSourceNode) {
	for _, st := range f.Statements {
		g.genStatement(st)
	}
}

func (g *llvmGenerator) genStatement(st ast.Statement) {
	switch s := st.(type) {
	case *ast.FunctionDeclaration:
		g.genFunctionDeclaration(s)
	case *ast.ReturnStatement:
		g.genReturnStatement(s)
	case *ast.BlockStatement:
		g.genBlockStatement(s)
	case *ast.ExpressionStatement:
		g.genExpressionStatement(s)
	default:
		codegen.DoPanic("unhandled statement node type: %T", s)
	}
}

func (g *llvmGenerator) genFunctionDeclaration(f *ast.FunctionDeclaration) {

}

func (g *llvmGenerator) genReturnStatement(r *ast.ReturnStatement) {

}

func (g *llvmGenerator) genBlockStatement(b *ast.BlockStatement) {

}

func (g *llvmGenerator) genExpressionStatement(es *ast.ExpressionStatement) {

}
