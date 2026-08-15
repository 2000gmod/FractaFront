package fir

type BlockID uint32

type Block struct {
	ID   BlockID
	Name string
	Func *Function

	Ins  []Instruction
	Term Terminator
}

func (b *Block) insertIns(ins ...Instruction) {
	b.Ins = append(b.Ins, ins...)
}

func (b *Block) insertTerm(term Terminator) {
	b.Term = term
}

func (b *Block) isTerminated() bool {
	return b.Term != nil
}
