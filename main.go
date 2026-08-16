package main

import (
	"fmt"
	_ "fracta/internal/codegen/tags"
	"fracta/internal/fir"
)

func main() {
	ctx := fir.NewContext()
	mod := ctx.NewModule("fracta_fir_test")

	puts := mod.NewExternalFunction(
		"puts",
		ctx.GetFuncType(
			false,
			ctx.GetVoidType(),
			ctx.GetPtrType(),
		),
		fir.ConvC,
	)

	printf := mod.NewExternalFunction(
		"printf",
		ctx.GetFuncType(
			true,
			ctx.GetIntType(32),
			ctx.GetPtrType(),
		),
		fir.ConvC,
	)

	fn := mod.NewFunction(
		"main",
		ctx.GetFuncType(
			false,
			ctx.GetIntType(32),
		),
		fir.ConvFracta,
	)

	fibo := mod.NewFunction(
		"fibonacci",
		ctx.GetFuncType(
			false,
			ctx.GetIntType(64),
			ctx.GetIntType(64),
		),
		fir.ConvFracta,
	)

	bd := fir.NewBuilder(fn)
	bd.NewBlock("entry")

	s := mod.GetGlobalString("str.1", "Hello, World\n\x00")
	bd.InsertCall(puts, ctx.GetVoidType(), fir.ConvC, s)

	res := bd.InsertCall(fibo, ctx.GetIntType(64), fir.ConvFracta, ctx.ConstInt(64, 7))
	fmtStr := mod.GetGlobalString("str.2", "fibo(%d) = %d\n\x00")
	bd.InsertCall(printf, ctx.GetIntType(32), fir.ConvC, fmtStr, ctx.ConstInt(32, 7), res)
	bd.InsertRet(ctx.VoidValue())

	bd = fir.NewBuilder(fibo)
	entry := bd.NewBlock("entry")
	base := bd.NewBlock("base")
	init := bd.NewBlock("init")
	loop := bd.NewBlock("loop")
	body := bd.NewBlock("body")
	done := bd.NewBlock("done")

	slotA := bd.InsertStackSlot(ctx.GetIntType(64), "a")
	slotB := bd.InsertStackSlot(ctx.GetIntType(64), "b")
	slotI := bd.InsertStackSlot(ctx.GetIntType(64), "i")

	bd.SetInsertionPoint(entry)
	cmp := bd.InsertCmp(fir.CmpUle, fibo.GetParameter(0), ctx.ConstInt(64, 0))
	bd.InsertCondBr(cmp, base, init)

	bd.SetInsertionPoint(base)
	bd.InsertRet(ctx.ConstInt(64, 0))

	bd.SetInsertionPoint(init)
	bd.InsertStore(ctx.ConstInt(64, 0), slotA)
	bd.InsertStore(ctx.ConstInt(64, 1), slotB)
	bd.InsertStore(ctx.ConstInt(64, 0), slotI)
	bd.InsertBr(loop)

	bd.SetInsertionPoint(loop)
	iv := bd.InsertLoad(ctx.GetIntType(64), slotI)
	cmpLoop := bd.InsertCmp(fir.CmpUlt, iv, fibo.GetParameter(0))
	bd.InsertCondBr(cmpLoop, body, done)

	bd.SetInsertionPoint(body)
	av := bd.InsertLoad(ctx.GetIntType(64), slotA)
	bv := bd.InsertLoad(ctx.GetIntType(64), slotB)
	next := bd.InsertBinaryOp(fir.BinAdd, av, bv)
	bd.InsertStore(bv, slotA)
	bd.InsertStore(next, slotB)
	iv2 := bd.InsertLoad(ctx.GetIntType(64), slotI)
	one := ctx.ConstInt(64, 1)
	i2 := bd.InsertBinaryOp(fir.BinAdd, iv2, one)
	bd.InsertStore(i2, slotI)
	bd.InsertBr(loop)

	bd.SetInsertionPoint(done)
	result := bd.InsertLoad(ctx.GetIntType(64), slotA)
	bd.InsertRet(result)

	fmt.Println(string(mod.Dump()))
}
