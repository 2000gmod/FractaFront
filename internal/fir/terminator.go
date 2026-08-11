package fir

type Terminator interface {
	terminator()
}

type BrTerm struct {
	Target *Block
}

func (t *BrTerm) terminator() {}

type CondBrTerm struct {
	Cond  Value
	True  *Block
	False *Block
}

func (t *CondBrTerm) terminator() {}

type ReturnTerm struct {
	Value Value
}

func (t *ReturnTerm) terminator() {}

type UnreachableTerm struct{}

func (t *UnreachableTerm) terminator() {}
