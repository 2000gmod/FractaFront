package fir

type Terminator interface {
	Instruction
	terminator()
}

type BrTerm struct {
	instructionBase
	Target *Block
}

func (t *BrTerm) terminator() {}

type CondBrTerm struct {
	instructionBase
	Cond  Value
	True  *Block
	False *Block
}

func (t *CondBrTerm) terminator() {}

type ReturnTerm struct {
	instructionBase
	Value Value
}

func (t *ReturnTerm) terminator() {}

type UnreachableTerm struct {
	instructionBase
}

func (t *UnreachableTerm) terminator() {}
