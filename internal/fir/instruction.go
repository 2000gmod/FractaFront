package fir

type instructionBase struct {
	parent *Block
}

func (i *instructionBase) Parent() *Block {
	return i.parent
}

type Instruction interface {
	Parent() *Block
}

type ValueInstruction interface {
	Instruction
	Value
}

type OperationKind uint8

const (
	none OperationKind = iota
	Load
	Store
	ElemPtr

	Add
	Sub
	Mul
	SDiv
	UDiv
	SRem
	URem

	FAdd
	FSub
	FMul
	FDiv
	FRem

	And
	Or
	Xor
	Shl
	LShr
	AShr

	Neg
	FNeg
	Not

	Cmp
	Cast

	Call
	ICall

	Br
	CondBr
	Ret
	Unreachable
)
